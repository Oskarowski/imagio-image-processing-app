package main

import (
	"imagio/cmd"
	"imagio/cmd/tui"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--help" {
		cmd.PrintHelp()
		return
	}

	logFile, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	if len(os.Args) > 1 {
		cmd.RunAsCliApp()
	} else {
		tui.RunAsTUIApp()
	}
}
