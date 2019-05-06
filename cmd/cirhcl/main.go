package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/russellrollins/cirhcl/pkg/convert"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(errors.Wrap(err, "error running cirhcl"))
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return errors.New("insufficient command line arguments, please specify the file to be converted")
	}
	if len(os.Args) > 2 {
		return errors.Errorf("unexpected command line arguments [%s]", strings.Join(os.Args[2:], " "))
	}

	name := os.Args[1]
	input, err := os.Open(name)
	if err != nil {
		return errors.Wrapf(err, "unable to open %s", name)
	}
	defer input.Close()

	noExtension := strings.TrimSuffix(input.Name(), filepath.Ext(input.Name()))
	outName := fmt.Sprintf("%s.yaml", noExtension)

	output, err := os.Create(outName)
	if err != nil {
		return errors.Wrapf(err, "unable to open %s", outName)
	}
	defer output.Close()

	return convert.Convert(input, name, output)

}
