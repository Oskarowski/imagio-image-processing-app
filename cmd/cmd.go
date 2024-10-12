package cmd

import (
	"path/filepath"
	"strings"
)

func IsImagePath(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".bmp"
}

type Command struct {
	Name string
	Args map[string]string
}

func ParseCommands(args []string) []Command {
	var commands []Command
	var currentCommand *Command

	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			if currentCommand != nil {
				commands = append(commands, *currentCommand)
			}

			currentCommand = &Command{
				Name: strings.TrimPrefix(arg, "--"),
				Args: make(map[string]string),
			}
		} else if strings.HasPrefix(arg, "-") && currentCommand != nil {
			// -argument=value), split it by '='

			parts := strings.SplitN(arg, "=", 2)

			if len(parts) == 2 {
				key := strings.TrimPrefix(parts[0], "-")
				currentCommand.Args[key] = parts[1]
			} else {
				currentCommand.Args[parts[0]] = ""
			}
		}
	}

	if currentCommand != nil {
		commands = append(commands, *currentCommand)
	}

	return commands
}
