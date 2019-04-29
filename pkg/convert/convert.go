package convert

import (
	"fmt"
	"io"
	"io/ioutil"
	"text/template"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
)

const (
	DefaultVersion         = "2"
	DefaultWorkflowVersion = "2"
)

type CircleConfig struct {
	Version         string            `hcl:"version,optional"`
	Jobs            []*CircleJob      `hcl:"job,block"`
	Workflows       []*CircleWorkflow `hcl:"workflow,block"`
	WorkflowVersion string            `hcl:"workflow_version,optional"`
}

type CircleJob struct {
	Name     string              `hcl:"name,label"`
	Docker   *CircleDockerConfig `hcl:"docker,block"`
	StepBody hcl.Body            `hcl:",remain"`
	Steps    []CircleJobStepper
}

type CircleJobStepper interface {
	JobStepYAML(int) string
	Decode(*hcl.Block) error
}

type CircleDockerConfig struct {
	Image string `hcl:"image"`
}

type CircleWorkflow struct {
	Name         string               `hcl:"name,label"`
	WorkflowJobs []*CircleWorkflowJob `hcl:"workflow_job,block"`
}

type CircleWorkflowJob struct {
	Name         string         `hcl:"name,label"`
	RequiresExpr hcl.Expression `hcl:"requires"`
	Requires     []string
}

var stepSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type: "checkout",
		},
		{
			Type: "run",
		},
		{
			Type: "persist_to_workspace",
		},
		{
			Type: "attach_workspace",
		},
	},
}

func Convert(r io.Reader, filename string, w io.Writer) error {
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "error in Convert opening input")
	}

	parser := hclparse.NewParser()
	f, diag := parser.ParseHCL(src, filename)
	if diag.HasErrors() {
		return errors.Wrap(diag, "error in Convert parsing HCL configuration")
	}

	var config CircleConfig
	if diag := gohcl.DecodeBody(f.Body, nil, &config); diag.HasErrors() {
		return errors.Wrap(diag, "error in Convert decoding HCL configuration")
	}

	if config.Version == "" {
		config.Version = DefaultVersion
	}
	if config.WorkflowVersion == "" {
		config.WorkflowVersion = DefaultWorkflowVersion
	}

	workflowJobVars := map[string]cty.Value{}
	for _, w := range config.Workflows {
		for _, wj := range w.WorkflowJobs {
			if _, ok := workflowJobVars[wj.Name]; ok {
				return errors.Errorf("multiple workflow_job steps named %s", wj.Name)
			}
			workflowJobVars[wj.Name] = cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal(wj.Name),
			})
		}
	}

	ctx := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"workflow_job": cty.ObjectVal(workflowJobVars),
		},
	}

	for _, w := range config.Workflows {
		for _, wj := range w.WorkflowJobs {
			diag := gohcl.DecodeExpression(wj.RequiresExpr, ctx, &wj.Requires)
			if diag.HasErrors() {
				return errors.Wrap(diag, "error in Convert decoding HCL configuration")
			}
		}
	}

	// TODO: Kinda iffy on this whole bit, although it _does_ work.
	// First, cast our body interface to a hclsyntax.Body
	for _, job := range config.Jobs {
		cBody, ok := job.StepBody.(*hclsyntax.Body)
		if !ok {
			return errors.New("error in Convert interface cast to Body failed")
		}

		// Next, extract the BodyContent from the hclsyntax.Body, using the schema defined in this file
		pBody, diags := cBody.Content(stepSchema)
		if diags.HasErrors() {
			return errors.Wrap(diags, "error in Convert dynamically decoding steps")
		}

		// Iterate over blocks, switch on the type, then decode into a known CircleJobStepper type
		for _, b := range pBody.Blocks {
			switch b.Type {
			case "checkout":
				sc := &StepCheckout{}
				if err := sc.Decode(b); err != nil {
					return errors.Wrap(err, "error in Convert dynamically decoding checkout step")
				}
				job.Steps = append(job.Steps, sc)
			case "run":
				sr := &StepRun{}
				if err := sr.Decode(b); err != nil {
					return errors.Wrap(err, "error in Convert dynamically decoding run step")
				}
				job.Steps = append(job.Steps, sr)
			case "persist_to_workspace":
				sp := &StepPersist{}
				if err := sp.Decode(b); err != nil {
					return errors.Wrap(err, "error in Convert dynamically decoding persist_to_workspace")
				}
				job.Steps = append(job.Steps, sp)
			case "attach_workspace":
				sa := &StepAttach{}
				if err := sa.Decode(b); err != nil {
					return errors.Wrap(err, "error in Convert dynamically decoding attach_workspace")
				}
				job.Steps = append(job.Steps, sa)
			default:
				// Should never happen since the above cases should always cover
				// all of the block types in our schema.
				panic(fmt.Errorf("unhandled block type %q", b.Type))
			}
		}
	}

	yamlTmpl, err := template.New("yaml").Parse(yamlTemplate)
	if err != nil {
		return errors.Wrap(err, "error in Convert parsing YAML template")
	}
	if err := yamlTmpl.Execute(w, config); err != nil {
		return errors.Wrap(err, "error in Convert executing YAML template")
	}

	return nil
}

const (
	yamlTemplate = `---
version: {{.Version}}
{{- if .Jobs}}
jobs:
  {{- range .Jobs}}
  {{.Name}}:
  {{- if .Docker}}
    docker:
      - image: {{.Docker.Image}}
  {{- end -}}
  {{- if .Steps}}
    steps:
    {{- range .Steps}}
{{.JobStepYAML 6}}
    {{- end -}}
  {{- end -}}
  {{- end -}}
{{- end -}}
{{- if .Workflows}}
workflows:
  version: {{.Version}}
  {{- range .Workflows}}
  {{.Name}}:
    jobs:
      {{- range .WorkflowJobs}}
      - {{.Name}}{{if .Requires}}:{{end}}
        {{- if .Requires}}
          requires:
          {{- range .Requires}}
            - {{.}}
          {{- end -}}
        {{- end -}}
      {{- end -}}
  {{- end -}}
{{- end -}}
`
)
