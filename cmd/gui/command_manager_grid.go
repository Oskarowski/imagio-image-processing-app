package gui

import (
	"fmt"
	"image-processing/cmd"
	"log"
	"strings"

	"github.com/rivo/tview"
)

func getCommandManagerGrid(AppState *AppState) *tview.Grid {
	commandManagerGrid := tview.NewGrid().
		SetRows(1, 0, 0).
		SetColumns(0, 0).
		SetBorders(true)

	commandManagerHeader := tview.NewTextView().
		SetText("Command Execution Manager - Select and Run Commands").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	commandList := tview.NewList()

	commandDetails := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)

	parameterForm := tview.NewForm()

	executionLog := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)

	var runCommand func(cmd cmd.CommandInfo)

	updateCommandDetails := func(cmd cmd.CommandInfo) {
		commandDetails.SetText(fmt.Sprintf("[yellow]%s[-]\n\n[white]%s[-]\n\n[cyan]Parameters:[-]\n%s",
			cmd.Name,
			cmd.Description,
			strings.Join(cmd.Arguments, "\n"),
		))

		parameterForm.Clear(true)
		for _, arg := range cmd.Arguments {
			parts := strings.SplitN(arg, ":", 2)
			if len(parts) > 1 {
				paramName := strings.Trim(parts[0], "-()")
				parameterForm.AddInputField("[cyan]"+paramName+"[-]", "", 20, nil, nil)
			}
		}
		parameterForm.AddButton("[green]Run[-]", func() {
			runCommand(cmd)
		})
	}

	runCommand = func(cmd cmd.CommandInfo) {
		args := []string{cmd.Name}
		for i := 0; i < parameterForm.GetFormItemCount()-1; i++ {
			input := parameterForm.GetFormItem(i).(*tview.InputField)
			args = append(args, fmt.Sprintf("-%s=%s", input.GetLabel(), input.GetText()))
		}
		logMessage := fmt.Sprintf("[green]Executing: %s[-]\n", strings.Join(args, " "))
		executionLog.SetText(logMessage)
		log.Println(logMessage)
	}

	for _, cmd := range cmd.AvailableCommands {
		command := cmd
		commandList.AddItem(cmd.Name, cmd.Description, 0, func() {
			updateCommandDetails(command)
		})
	}

	commandManagerGrid.AddItem(commandManagerHeader, 0, 0, 1, 2, 0, 0, false)
	commandManagerGrid.AddItem(commandList, 1, 0, 1, 1, 0, 0, true)
	commandManagerGrid.AddItem(commandDetails, 1, 1, 1, 1, 0, 0, false)
	commandManagerGrid.AddItem(parameterForm, 2, 0, 1, 1, 0, 0, false)
	commandManagerGrid.AddItem(executionLog, 2, 1, 1, 1, 0, 0, false)

	return commandManagerGrid
}
