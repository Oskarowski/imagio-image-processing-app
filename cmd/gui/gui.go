package gui

import (
	"fmt"
	"image"
	"image-processing/imageio"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LoadedImage struct {
	Filename  string
	Filepath  string
	ImageData image.Image
	Loaded    bool
}

type LoadedImagesStore struct {
	Images            []LoadedImage
	ChangeSubscribers []func()
}

func (store *LoadedImagesStore) notifyChangeSubscribers() {
	for _, sub := range store.ChangeSubscribers {
		sub()
	}
}

func (store *LoadedImagesStore) SubscribeToChanges(sub func()) {
	store.ChangeSubscribers = append(store.ChangeSubscribers, sub)
}

func (store *LoadedImagesStore) loadImageData(index int) error {
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

func (store *LoadedImagesStore) removeImage(index int) {
	if index < 0 || index >= len(store.Images) {
		return
	}
	store.Images = append(store.Images[:index], store.Images[index+1:]...)
	store.notifyChangeSubscribers()
}

func (store *LoadedImagesStore) addImage(filename string, filepath string) {
	store.Images = append(store.Images, LoadedImage{Filename: filename, Filepath: filepath, Loaded: false})
	store.notifyChangeSubscribers()
}

type AppState struct {
	App      *tview.Application
	Pages    *tview.Pages
	ImgStore *LoadedImagesStore
}

func RunAsTUIApp() {
	appState := &AppState{
		App:      tview.NewApplication(),
		Pages:    tview.NewPages(),
		ImgStore: &LoadedImagesStore{Images: []LoadedImage{}},
	}

	imgLoaderManagerGrid := getImageLoaderGrid(appState)
	commandManagerGrid := getCommandManagerGrid(appState)

	appPages := tview.NewPages()

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

	appState.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
			appState.App.Stop()
		}

		return event
	})

	appPages.AddPage("Image Loader Manager", wrapWithNavBar(imgLoaderManagerGrid), true, true)
	appPages.AddPage("Command Execution Manager", wrapWithNavBar(commandManagerGrid), true, false)

	if err := appState.App.SetRoot(appPages, true).SetFocus(appPages).EnableMouse(true).Run(); err != nil {
		log.Fatalf("Failed to start TUI app: %v", err)
	}
}
