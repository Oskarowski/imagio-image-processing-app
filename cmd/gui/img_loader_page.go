package gui

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func getImageLoaderGrid(appState *AppState) *tview.Grid {
	imgLoaderManagerGrid := tview.NewGrid().
		SetRows(1, 0).
		SetColumns(0, 0).
		SetBorders(true)

	header := tview.NewTextView().
		SetText("Image Loader Manager - Navigate and Load Images").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	fileSystemList := tview.NewList()
	fileSystemList.SetTitle(" File System ").SetBorder(true)

	loadedImagesList := tview.NewList()
	loadedImagesList.SetTitle(" Loaded Images ").SetBorder(true)

	var updateFileSystemList func(dir string)
	var updateLoadedImagesList func()

	initialDir := "./imgs"
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
					appState.ImgStore.addImage(itemName, fullPath)
					appState.ImgStore.loadImageData(len(appState.ImgStore.Images) - 1)
					updateLoadedImagesList()
				})
			}
		}
	}

	fileSystemList.SetFocusFunc(func() {
		fileSystemList.SetBorderColor(tcell.Color178)
		loadedImagesList.SetBorderColor(tcell.ColorWhite)
	})

	loadedImagesList.SetFocusFunc(func() {
		loadedImagesList.SetBorderColor(tcell.Color178)
		fileSystemList.SetBorderColor(tcell.ColorWhite)
	})

	updateFileSystemList(initialDir)

	updateLoadedImagesList = func() {
		loadedImagesList.Clear()

		var displayItemName string

		for i, img := range appState.ImgStore.Images {
			if img.Loaded {
				displayItemName = "[green]📷 " + img.Filename + "[-]"
				idx := i
				loadedImagesList.AddItem(displayItemName, "", 0, func() {
					appState.ImgStore.removeImage(idx)
					updateLoadedImagesList()
				})
			}
		}
	}

	imgLoaderManagerGrid.AddItem(header, 0, 0, 1, 2, 0, 0, false)
	imgLoaderManagerGrid.AddItem(fileSystemList, 1, 0, 1, 1, 0, 0, true)
	imgLoaderManagerGrid.AddItem(loadedImagesList, 1, 1, 1, 1, 0, 0, false)

	currentFocus := "left"

	imgLoaderManagerGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			if currentFocus == "left" {
				currentFocus = "right"
				appState.App.SetFocus(loadedImagesList)
			} else {
				currentFocus = "left"
				appState.App.SetFocus(fileSystemList)
			}
		}
		return event
	})

	return imgLoaderManagerGrid
}
