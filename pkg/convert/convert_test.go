package convert

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

"github.com/stretchr/testify/assert"
)

const (
	examplePath = "testdata/examples"
)

type exampleTest struct {
	name   string
	input  string
	output string
}

func TestConvertExamples(t *testing.T) {
	tcs := getExamples(t)
	for _, tc := range tcs {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			in, err := os.Open(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			defer in.Close()

			var gotBuffer bytes.Buffer
			err = Convert(in, tc.input, &gotBuffer)
			if err != nil {
				t.Fatal(err)
			}

			out, err := os.Open(tc.output)
			if err != nil {
				t.Fatal(err)
			}
			defer out.Close()

			wantBytes, err := ioutil.ReadAll(out)
			if err != nil {
				t.Fatal(err)
			}

			want := string(wantBytes)
			got := gotBuffer.String()

      assert.Equal(t, want, got, "The converted YAML should equal the expected YAML.")
		})
	}
}

func getExamples(t *testing.T) []exampleTest {
	outRe := regexp.MustCompile(`examples/(?P<test>.*)-out\.yml`)
	inRe := regexp.MustCompile(`examples/(?P<test>.*)-in\.hcl`)

	var inputs []string
	var outputs []string

	err := filepath.Walk(examplePath, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			t.Fatalf("error walking path %q: %v", path, err)
		}

		match := inRe.FindStringSubmatch(path)
		if len(match) != 0 {
			result := make(map[string]string)
			for i, name := range inRe.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}
			inputs = append(inputs, result["test"])
		}

		match = outRe.FindStringSubmatch(path)
		if len(match) != 0 {
			result := make(map[string]string)
			for i, name := range outRe.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}
			outputs = append(outputs, result["test"])
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(inputs) != len(outputs) {
		t.Fatalf("example inputs do not match example outputs")
	}

	if len(inputs) == 0 {
		t.Fatalf("no examples found in testdata/examples")
	}

	for _, i := range inputs {
		found := false
		for _, o := range outputs {
			if i == o {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("example inputs do not match example outputs")
		}
	}

	var tcs []exampleTest
	for _, i := range inputs {
		tcs = append(tcs, exampleTest{
			name:   i,
			input:  fmt.Sprintf("%s/%s-in.hcl", examplePath, i),
			output: fmt.Sprintf("%s/%s-out.yml", examplePath, i),
		})
	}

	return tcs
}
