package gui

import (
	"fmt"
	"image-processing/cmd"
	"log"
	"os/exec"
	"strings"

	"github.com/rivo/tview"
)

var commandBuilder struct {
	File                    string
	FilePath                string
	ComparisonImage         string
	CommandSet              []string
	RequiresComparisonImage bool
}

func getCommandManagerGrid(AppState *AppState) *tview.Grid {
	grid := tview.NewGrid().
		SetRows(1, 1, 0, 20, 0, 0, 0).
		SetColumns(0).
		SetBorders(true)

	commandManagerHeader := tview.NewTextView().
		SetText("Command Execution Manager - Select and Run Commands").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	fileNameHeader := tview.NewTextView().
		SetText("File: [cyan]<No file selected>[-]").
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetWrap(false)

	loadedFilesList := tview.NewList().
		ShowSecondaryText(false)
	loadedFilesList.SetTitle(" File System ").SetBorder(true)

	var updateLoadedImagesList func()
	updateLoadedImagesList = func() {
		loadedFilesList.Clear()

		for i, img := range AppState.ImgStore.Images {
			var displayItemName string
			if img.Loaded {
				displayItemName = "[green]📷 " + img.Filename + "[-]"
			} else {
				displayItemName = "📷 " + img.Filename
			}
			idx := i
			loadedFilesList.AddItem(displayItemName, "", 0, func() {
				for j := range AppState.ImgStore.Images {
					AppState.ImgStore.Images[j].Loaded = (j == idx)
				}
				fileNameHeader.SetText("File: [cyan]" + AppState.ImgStore.Images[idx].Filename + "[-]")
				commandBuilder.File = AppState.ImgStore.Images[idx].Filename
				commandBuilder.FilePath = AppState.ImgStore.Images[idx].Filepath
				updateLoadedImagesList()
			})
		}
	}
	AppState.ImgStore.SubscribeToChanges(updateLoadedImagesList)

	commandList := tview.NewList()
	commandDetails := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)

	parameterForm := tview.NewForm()

	executionLog := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetScrollable(true)

	var runCommand func(cmd cmd.CommandInfo)

	updateCommandDetails := func(cmd cmd.CommandInfo) {
		commandDetails.SetText(fmt.Sprintf("[yellow]%s[-]\n\n[white]%s[-]\n\n[cyan]Parameters:[-]\n%s",
			cmd.Name,
			cmd.Description,
			strings.Join(cmd.Arguments, "\n"),
		))

		if cmd.Name == "mse" {
			commandBuilder.RequiresComparisonImage = true
		}

		parameterForm.Clear(true)

		for _, arg := range cmd.Arguments {
			parts := strings.SplitN(arg, ":", 2)
			if len(parts) > 1 {
				paramName := strings.Trim(parts[0], "-")
				parameterForm.AddInputField(paramName, "", 20, nil, nil)
			}
		}

		parameterForm.AddButton("Run", func() {
			runCommand(cmd)
		})
	}

	// TODO add feature to chain commands
	// TODO add option to select comparison image
	runCommand = func(cmd cmd.CommandInfo) {
		if commandBuilder.RequiresComparisonImage && commandBuilder.ComparisonImage == "" {
			executionLog.SetText("Please load a comparison image first")
			log.Println("Error: A comparison image is required for this command.")
			return
		}

		args := []string{"go", "run", "main.go", "--" + cmd.Name}

		for i := 0; i < parameterForm.GetFormItemCount(); i++ {
			input := parameterForm.GetFormItem(i).(*tview.InputField)
			trimmedLabel := strings.Split(input.GetLabel(), "=")[0]
			value := input.GetText()
			if value != "" {
				args = append(args, fmt.Sprintf("-%s=%s", trimmedLabel, value))
			}
		}

		if commandBuilder.File != "" {
			args = append(args, commandBuilder.FilePath)
		} else {
			executionLog.SetText("[red]Error: No file selected for processing[-]")
			log.Println("Error: No file selected for processing.")
			return
		}

		fullCommand := strings.Join(args, " ")

		executionLog.SetText(fullCommand)
		log.Printf("Executing command: %s\n", fullCommand)

		execCmd := exec.Command(args[0], args[1:]...)
		output, err := execCmd.CombinedOutput()
		if err != nil {
			executionLog.SetText(fmt.Sprintf("[red]Error:[-] %v\n[red]Output:[-] %s", err, string(output)))
			log.Printf("Command execution failed: %v\n", err)
			return
		}

		log.Printf("Command executed successfully: %s\n", string(output))
		executionLog.SetText(fmt.Sprintf("[green]Success![-]\n%s", string(output)))
	}

	for _, cmd := range cmd.AvailableCommands {
		command := cmd
		commandList.AddItem(cmd.Name, cmd.Description, 0, func() {
			updateCommandDetails(command)
		})
	}

	grid.AddItem(commandManagerHeader, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(fileNameHeader, 1, 0, 1, 1, 0, 0, false)
	grid.AddItem(loadedFilesList, 2, 0, 1, 1, 0, 0, true)
	grid.AddItem(commandList, 3, 0, 1, 1, 0, 0, true)
	grid.AddItem(commandDetails, 4, 0, 1, 1, 0, 0, false)
	grid.AddItem(parameterForm, 5, 0, 1, 1, 0, 0, false)
	grid.AddItem(executionLog, 6, 0, 1, 1, 0, 0, false)

	return grid
}
