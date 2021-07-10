package main

import (
	"fmt"
	"os"

	"github.com/hlcfan/langtool/langtool"
)

func main() {
	if err := langtool.Check(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}
