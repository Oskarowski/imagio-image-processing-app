package tui

import (
	"errors"
	"fmt"
	"image"
	"image-processing/imageio"
	"image-processing/internal/ascii_preview"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type view int

type terminalSize struct {
	width  int
	height int
}

type listItem struct {
	title, desc string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.desc }
func (i listItem) FilterValue() string { return i.title }

type clearErrorMsg struct{}

const (
	filePickerView view = iota
	imagePreviewView
	commandsSelectionView
	commandExecutionView
)

type model struct {
	filepicker      filepicker.Model
	selectedFile    string
	currentView     view
	quitting        bool
	err             error
	imagePreview    string
	loadedImage     image.Image
	terminalSize    terminalSize
	commandsList    list.Model
	selectedCommand string
}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m *model) loadImagePreview(path string) {
	file, err := imageio.OpenBmpImage(path)
	if err != nil {
		m.err = fmt.Errorf("failed to open image: %v", err)
		return
	}

	m.loadedImage = file

	availableHeight := m.terminalSize.height

	convertOptions := ascii_preview.DefaultOptions
	convertOptions.FixedWidth = availableHeight * 2
	convertOptions.FixedHeight = availableHeight

	converter := ascii_preview.NewImageConverter()
	converted := converter.Image2ASCIIString(m.loadedImage, &convertOptions)

	m.imagePreview = converted
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.currentView == filePickerView {
		m.filepicker, cmd = m.filepicker.Update(msg)
	} else if m.currentView == commandsSelectionView {
		m.commandsList, cmd = m.commandsList.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "tab":
			m.currentView = (m.currentView + 1) % 4
			return m, nil

		case "enter":
			if m.currentView == filePickerView {
				if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
					m.selectedFile = path
					m.loadImagePreview(path)
					m.currentView = imagePreviewView

					return m, nil
				}

				if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
					m.err = errors.New(path + " is not valid.")
					m.selectedFile = ""
					return m, clearErrorAfter(2 * time.Second)
				}
			}

			if m.currentView == commandsSelectionView {
				selectedItem := m.commandsList.SelectedItem()
				log.Default().Println("Selected item: ", selectedItem)

			}

		case "r":
			if m.currentView == commandExecutionView {
				log.Default().Println("Executing command: ", m.selectedCommand)
			}
		}

	case tea.WindowSizeMsg:
		m.terminalSize.width = msg.Width
		m.terminalSize.height = msg.Height
		m.commandsList.SetSize(msg.Width, msg.Height)
		return m, nil

	case clearErrorMsg:
		m.err = nil
	}

	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	switch m.currentView {
	case filePickerView:
		s.WriteString("\n  Select files (Press 'Tab' to view selected files, 'Enter' to add):\n")
		if m.err != nil {
			s.WriteString("  " + m.filepicker.Styles.DisabledFile.Render(m.err.Error()) + "\n")
		}
		s.WriteString("\n" + m.filepicker.View() + "\n")

	case imagePreviewView:
		s.WriteString("\n  Image Preview (Press 'Tab' to view selected files):\n")
		s.WriteString("\n" + m.imagePreview + "\n")

	case commandsSelectionView:
		s.WriteString("\n  Select a command (Press 'Tab' to view selected files, 'Enter' to enter for details):\n")

		s.WriteString(m.commandsList.View())

	}

	return s.String()
}

func RunAsTUIApp() {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".png", ".bmp"}
	fp.ShowHidden = false
	fp.ShowSize = false
	fp.ShowPermissions = false
	fp.CurrentDirectory, _ = os.Getwd()

	commands := []list.Item{
		listItem{title: "hflip", desc: "Apply horizontal flip to the image."},
		listItem{title: "vflip", desc: "Apply vertical flip to the image."},
		listItem{title: "bandpass", desc: "Apply bandpass filtering to the image."},
	}

	commandList := list.New(commands, list.NewDefaultDelegate(), 0, 0)
	commandList.Title = "Available Commands"

	m := model{
		filepicker:   fp,
		currentView:  filePickerView,
		commandsList: commandList,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nSee you next time " + m.filepicker.Styles.Selected.Render("👋"))
}
