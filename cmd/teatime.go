package main

import (
	"fmt"
	"os"

	"github.com/cspengl/teatime/internal/tui"
)

func main() {

	program := tui.LoadProgram(tui.Interactive, os.Args)
	if err := program.Start(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
