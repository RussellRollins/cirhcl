package convert

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pkg/errors"
)

type StepLiteral struct {
	content string
}

func (sl *StepLiteral) JobStepYAML(indent int) string {
	return indentLines([]string{sl.content}, indent)
}

func (sl *StepLiteral) Decode(*hcl.Block) error {
	return nil
}

type StepCheckout struct{}

func (sc *StepCheckout) JobStepYAML(indent int) string {
	return indentLines([]string{"- checkout"}, indent)
}

func (sc *StepCheckout) Decode(*hcl.Block) error {
	return nil
}

type StepRun struct {
	Command string `hcl:"command"`
}

func (sr *StepRun) JobStepYAML(indent int) string {
	return indentLines([]string{fmt.Sprintf("- run: %s", sr.Command)}, indent)
}

func (sr *StepRun) Decode(block *hcl.Block) error {
	diags := gohcl.DecodeBody(block.Body, nil, sr)
	if diags.HasErrors() {
		return errors.Wrap(diags, "error in StepRun Decode")
	}
	return nil
}

type StepPersist struct {
	Root  string   `hcl:"root"`
	Paths []string `hcl:"paths,optional"`
}

func (sp *StepPersist) JobStepYAML(indent int) string {
	parts := []string{
		"- persist_to_workspace:",
		fmt.Sprintf("    root: %s", sp.Root),
	}
	if len(sp.Paths) != 0 {
		parts = append(parts, "    paths:")
	}
	for _, p := range sp.Paths {
		parts = append(parts, fmt.Sprintf("      - %s", p))
	}

	return indentLines(parts, indent)
}

func (sp *StepPersist) Decode(block *hcl.Block) error {
	diags := gohcl.DecodeBody(block.Body, nil, sp)
	if diags.HasErrors() {
		return errors.Wrap(diags, "error in StepPersist Decode")
	}
	return nil
}

type StepAttach struct {
	At string `hcl:"at"`
}

func (sa *StepAttach) JobStepYAML(indent int) string {
	parts := []string{
		"- attach_workspace:",
		fmt.Sprintf("    at: %s", sa.At),
	}
	return indentLines(parts, indent)
}

func (sa *StepAttach) Decode(block *hcl.Block) error {
	diags := gohcl.DecodeBody(block.Body, nil, sa)
	if diags.HasErrors() {
		return errors.Wrap(diags, "error in StepAttach Decode")
	}
	return nil
}

func indentLines(lines []string, indent int) string {
	ret := ""
	spacer := strings.Repeat(" ", indent)
	for _, l := range lines {
		ret = fmt.Sprintf("%s%s%s\n", ret, spacer, l)
	}
	ret = strings.TrimRight(ret, "\n")
	return ret
}
