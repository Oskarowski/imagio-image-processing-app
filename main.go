package main

import (
	"fmt"
	"image"
	"image-processing/cmd"
	"image-processing/imageio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--help" {
		cmd.PrintHelp()
		return
	}

	if len(os.Args) > 1 {
		cmd.RunAsCliApp()
	} else {
		RunAsTUIApp()
	}
}

type LoadedImage struct {
	Filename  string
	Filepath  string
	ImageData image.Image
	Loaded    bool
}

type LoadedImagesStore struct {
	Images []LoadedImage
}

func (store *LoadedImagesStore) AddImage(filename string, filepath string) {
	store.Images = append(store.Images, LoadedImage{Filename: filename, Filepath: filepath, Loaded: false})
}

func (store *LoadedImagesStore) LoadImageData(index int) error {
	if index < 0 || index >= len(store.Images) {
		return fmt.Errorf("index out of range")
	}

	img, err := imageio.OpenBmpImage(store.Images[index].Filepath)
	if err != nil {
		return fmt.Errorf("error opening BMP image file: %v", err)
	}

	store.Images[index].ImageData = img
	store.Images[index].Loaded = true
	return nil
}

func (store *LoadedImagesStore) RemoveImage(index int) {
	if index < 0 || index >= len(store.Images) {
		return
	}

	store.Images = append(store.Images[:index], store.Images[index+1:]...)
}

func RunAsTUIApp() {
	logFile, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	app := tview.NewApplication()

	appPages := tview.NewPages()

	initialDir := "./imgs"

	imageStore := &LoadedImagesStore{Images: []LoadedImage{}}

	var updateFileSystemList func(dir string)
	var updateLoadedImagesList func()

	fileSystemList := tview.NewList()
	fileSystemList.SetTitle(" File System ").SetBorder(true)

	updateFileSystemList = func(dir string) {
		fileSystemList.Clear()

		if dir != "./" {
			fileSystemList.AddItem("../", "Go up one level", 0, func() {
				updateFileSystemList(filepath.Dir(dir))
			})
		}

		files, err := os.ReadDir(dir)
		if err != nil {
			log.Fatalf("Error reading directory '%s': %v", dir, err)
		}

		for _, file := range files {
			itemName := file.Name()
			var displayItemName string
			fullPath := filepath.Join(dir, itemName)

			if file.IsDir() {
				displayItemName = "[blue]📁 " + itemName + "[-]"
				fileSystemList.AddItem(displayItemName+"/", "", 0, func(targetDir string) func() {
					return func() {
						updateFileSystemList(fullPath)
					}
				}(file.Name()))
			} else {
				displayItemName = "[green]📄 " + itemName + "[-]"
				fileSystemList.AddItem(displayItemName, "", 0, func() {
					imageStore.AddImage(itemName, fullPath)
					imageStore.LoadImageData(len(imageStore.Images) - 1)
					updateLoadedImagesList()
				})
			}
		}
	}

	updateFileSystemList(initialDir)

	loadedImagesList := tview.NewList()
	loadedImagesList.SetTitle(" Loaded Images ").SetBorder(true)

	updateLoadedImagesList = func() {
		loadedImagesList.Clear()

		var displayItemName string

		for i, img := range imageStore.Images {
			if img.Loaded {
				displayItemName = "[green]📷 " + img.Filename + "[-]"
				idx := i
				loadedImagesList.AddItem(displayItemName, "", 0, func() {
					imageStore.RemoveImage(idx)
					updateLoadedImagesList()
				})
			}
		}
	}

	imgLoaderManagerGrid := tview.NewGrid().
		SetRows(1, 0).
		SetColumns(0, 0).
		SetBorders(true)

	header := tview.NewTextView().
		SetText("Image Loader Manager - Navigate and Load Images").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	imgLoaderManagerGrid.AddItem(header, 0, 0, 1, 2, 0, 0, false)
	imgLoaderManagerGrid.AddItem(fileSystemList, 1, 0, 1, 1, 0, 0, true)
	imgLoaderManagerGrid.AddItem(loadedImagesList, 1, 1, 1, 1, 0, 0, false)

	currentFocus := "left"

	imgLoaderManagerGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			if currentFocus == "left" {
				currentFocus = "right"
				app.SetFocus(loadedImagesList)
				fileSystemList.SetBorderColor(tcell.ColorWhite)
				loadedImagesList.SetBorderColor(tcell.ColorYellow)
			} else {
				currentFocus = "left"
				app.SetFocus(fileSystemList)
				fileSystemList.SetBorderColor(tcell.ColorYellow)
				loadedImagesList.SetBorderColor(tcell.ColorWhite)
			}
		}
		return event
	})

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

	commandManagerGrid := tview.NewGrid().
		SetRows(1, 0, 0).
		SetColumns(0, 0).
		SetBorders(true)

	commandManagerHeader := tview.NewTextView().
		SetText("Command Execution Manager - Select and Run Commands").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	commandManagerGrid.AddItem(commandManagerHeader, 0, 0, 1, 2, 0, 0, false)
	commandManagerGrid.AddItem(commandList, 1, 0, 1, 1, 0, 0, true)
	commandManagerGrid.AddItem(commandDetails, 1, 1, 1, 1, 0, 0, false)
	commandManagerGrid.AddItem(parameterForm, 2, 0, 1, 1, 0, 0, false)
	commandManagerGrid.AddItem(executionLog, 2, 1, 1, 1, 0, 0, false)

	navBar := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]Image Loader Manager[-] [white]Command Execution Manager[-]").
		SetTextAlign(tview.AlignCenter)

	wrapWithNavBar := func(content tview.Primitive) *tview.Grid {
		grid := tview.NewGrid().
			SetRows(1, 0).
			SetColumns(0).
			SetBorders(false)
		grid.AddItem(navBar, 0, 0, 1, 1, 0, 0, false)
		grid.AddItem(content, 1, 0, 1, 1, 0, 0, true)
		return grid
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'l':
			appPages.SwitchToPage("Image Loader Manager")
			navBar.SetText("[yellow]Image Loader Manager[-] [white]Command Execution Manager[-]")
			return nil
		case 'c':
			appPages.SwitchToPage("Command Execution Manager")
			navBar.SetText("[white]Image Loader Manager[-] [yellow]Command Execution Manager[-]")
			return nil
		}

		switch event.Key() {
		case tcell.KeyESC:
			app.Stop()
		}

		return event
	})

	appPages.AddPage("Image Loader Manager", wrapWithNavBar(imgLoaderManagerGrid), true, true)
	appPages.AddPage("Command Execution Manager", wrapWithNavBar(commandManagerGrid), true, false)

	if err := app.SetRoot(appPages, true).SetFocus(appPages).EnableMouse(true).Run(); err != nil {
		log.Fatalf("Failed to start TUI app: %v", err)
	}
}
