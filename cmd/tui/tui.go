package tui

import (
	"fmt"
	"image-processing/cmd/executioner"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	UIState      UIState
	CommandState CommandState
	filepicker   filepicker.Model
	selectedFile string
	currentView  appView
	quitting     bool
	imagePreview string
	terminalSize terminalSize
	commandsList list.Model
	form         *huh.Form
	formGetter   func() map[string]string
}

func (m Model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m Model) View() string {
	if m.quitting {
		return "Goodbye!"
	}

	navBar := m.renderNavBar()

	var content string
	switch m.currentView {
	case FILE_PICKER_VIEW:
		content = m.filePickerView()
	case IMAGE_PREVIEW_VIEW:
		content = m.imagePreviewView()
	case COMMAND_SELECTION_VIEW:
		content = m.commandSelectionView()
	case COMMAND_EXECUTION_VIEW:
		content = m.commandExecutionView()
	}

	fullContent := lipgloss.JoinVertical(lipgloss.Left, navBar, content)

	containerStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Width(m.terminalSize.width - 2).
		Height(m.terminalSize.height - 2).
		Border(lipgloss.RoundedBorder())

	return containerStyle.Render(fullContent)
}

func RunAsTUIApp() {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".png", ".bmp"}
	fp.ShowHidden = false
	fp.ShowSize = false
	fp.ShowPermissions = false
	fp.CurrentDirectory, _ = os.Getwd()

	var commandItems []list.Item
	for _, cmd := range executioner.CommandDefinitions {
		commandItems = append(commandItems, cmd)
	}
	commandList := list.New(commandItems, list.NewDefaultDelegate(), 0, 0)
	commandList.Title = "Available Commands"

	m := Model{
		filepicker:   fp,
		currentView:  FILE_PICKER_VIEW,
		commandsList: commandList,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nSee you next time " + m.filepicker.Styles.Selected.Render("👋"))
}
