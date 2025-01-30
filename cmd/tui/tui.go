package tui

import (
	"errors"
	"fmt"
	"image"
	"image-processing/imageio"
	"image-processing/internal/ascii_preview"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type view int

type terminalSize struct {
	width  int
	height int
}

const (
	filePickerView view = iota
	imagePreviewView
	commandsSelectionView
	commandExecutionView
)

type model struct {
	filepicker   filepicker.Model
	selectedFile string
	currentView  view
	quitting     bool
	err          error
	imagePreview string
	loadedImage  image.Image
	terminalSize terminalSize
	commandsList viewport.Model
}

type clearErrorMsg struct{}

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
	m.filepicker, cmd = m.filepicker.Update(msg)

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

		}

	case tea.WindowSizeMsg:
		m.terminalSize.width = msg.Width
		m.terminalSize.height = msg.Height
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

	m := model{
		filepicker:  fp,
		currentView: filePickerView,
	}

	tm, _ := tea.NewProgram(&m).Run()
	mm := tm.(model)
	fmt.Println("\nSee you next time " + mm.filepicker.Styles.Selected.Render("👋"))
}
