package convert

import (
	"testing"
)

func TestStepLiteral(t *testing.T) {
	tcs := []struct {
		name    string
		content string
    indent int
    want string
	}{
		{
			name:    "Test Empty String",
			content: "",
      indent: 0,
      want: "",
		},
		{
			name:    "Test Content String",
			content: "Hello, World!",
      indent: 0,
      want: "Hello, World!",
		},
    {
      name: "Test Indent",
      content: "- test",
      indent: 4,
      want: "    - test",
    },
	}

	for _, tc := range tcs {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sl := &StepLiteral{content: tc.content}
			want := tc.want
			got := sl.JobStepYAML(tc.indent)
			if got != want {
				t.Errorf("want: `%s`\ngot: `%s`\n", want, got)
			}
		})
	}
}
