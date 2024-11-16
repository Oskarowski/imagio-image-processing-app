package main

import (
	"fmt"
	"image-processing/cmd"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--help" {
		cmd.PrintHelp()
		return
	}

	if len(os.Args) > 1 {
		cmd.RunAsCliApp()
	} else {
		fmt.Println("Should run as TUI app")
	}
}
