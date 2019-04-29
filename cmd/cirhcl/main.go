package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/russellrollins/cirhcl/pkg/convert"
)

func main() {
	fmt.Println("Hello, World!")
	convert.Convert(os.Stdin, "foo", os.Stdout)

	err := errors.Errorf("Error: %d", 3)
	fmt.Println(err)
}
