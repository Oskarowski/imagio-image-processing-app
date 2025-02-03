package tui

import (
	"fmt"
	"image-processing/cmd/tui/executioner"
	"image-processing/imageio"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	UIState      UIState
	CommandState CommandState
	UIStyles     UIStyle
	filepicker   filepicker.Model
	selectedFile string
	currentView  appView
	quitting     bool
	imagePreview string
	terminalSize terminalSize
	commandsList list.Model
	cursor       int
}

func (m Model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m *Model) initializeTextInputs() {
	m.CommandState.inputs = nil
	m.CommandState.commandArgs = make(map[string]string)

	for _, cmd := range executioner.CommandDefinitions {
		if cmd.Name == m.CommandState.selectedCommand {
			m.CommandState.inputs = make([]textinput.Model, len(cmd.Args))
			for j, arg := range cmd.Args {
				input := textinput.New()
				input.Placeholder = arg
				m.CommandState.commandArgs[arg] = ""
				m.CommandState.inputs[j] = input
			}

			m.cursor = -1
			m.CommandState.inputs[0].Focus()
			break
		}
	}
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	switch m.currentView {
	case FILE_PICKER_VIEW:

		titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)

		title := titleStyle.Render("Select Files (Press 'Tab' to switch views, 'Enter' to select):")

		var errorMessage string
		if m.UIState.err != nil {
			errorMessage = errorStyle.Render("Error: " + m.UIState.err.Error())
		}

		content := lipgloss.JoinVertical(lipgloss.Top,
			title,
			errorMessage,
			m.filepicker.View(),
		)

		return m.UIStyles.filepickerViewStyle.Render(content)

	case IMAGE_PREVIEW_VIEW:
		s.WriteString("\n  Image Preview (Press 'Tab' to view selected files):\n")
		s.WriteString(m.imagePreview)

	case COMMAND_SELECTION_VIEW:
		s.WriteString("\n  Select a command (Press 'Tab' to view selected files, 'Enter' to enter for details):\n")
		s.WriteString(m.commandsList.View())

	case COMMAND_EXECUTION_VIEW:
		titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
		labelStyle := lipgloss.NewStyle().Bold(true)

		s.WriteString(titleStyle.Render("Execute Command") + "\n\n")

		if m.CommandState.selectedCommand != "" {
			s.WriteString(labelStyle.Render("Command: ") + m.CommandState.selectedCommand + "\n\n")
		} else {
			s.WriteString(labelStyle.Render("No command selected!") + "\n\n")
		}

		if m.selectedFile != "" {
			fileName := imageio.GetPureFileName(m.selectedFile)
			s.WriteString(labelStyle.Render("File: ") + fileName + "\n\n")
		} else {
			s.WriteString(labelStyle.Render("No file selected!") + "\n\n")
		}

		for _, input := range m.CommandState.inputs {
			s.WriteString(input.View() + "\n")
		}

		buttonStyle := lipgloss.NewStyle().
			Padding(0, 4).
			Background(lipgloss.Color("205")).
			Foreground(lipgloss.Color("0")).
			Bold(true)

		if m.CommandState.selectedCommand != "" && m.selectedFile != "" {
			if m.cursor == len(m.CommandState.inputs) {
				s.WriteString("\n\n" + buttonStyle.Render("[ Execute ]"))
			} else {
				s.WriteString("\n\n" + lipgloss.NewStyle().Render("[ Execute ]"))
			}
		}

		if m.UIState.err != nil {
			errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
			s.WriteString("\n\n" + errorStyle.Render("Error: "+m.UIState.err.Error()))
		} else if m.UIState.successMessage != "" {
			successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
			s.WriteString("\n\n" + successStyle.Render("Success: "+m.UIState.successMessage))
		}

		return m.UIStyles.commandExecutionViewStyle.Render(s.String())
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
