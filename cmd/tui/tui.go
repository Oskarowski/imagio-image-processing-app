package tui

import (
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type view int

type clearErrorMsg struct{}
type clearSuccessMsg struct{}

const (
	filePickerView view = iota
	imagePreviewView
	commandsSelectionView
	commandDetailView
	commandExecutionView
)

type model struct {
	filepicker             filepicker.Model
	selectedFile           string
	currentView            view
	quitting               bool
	err                    error
	successMessage         string
	imagePreview           string
	loadedImage            image.Image
	terminalSize           terminalSize
	commandsList           list.Model
	selectedCommand        string
	selectedCommandArgs    map[string]string
	inputs                 []textinput.Model
	cursor                 int
	commandDetailViewStyle lipgloss.Style
	filePickerViewStyle    lipgloss.Style
}

type commandDefinition struct {
	name        string
	description string
	args        []string
}

func (i commandDefinition) Title() string       { return i.name }
func (i commandDefinition) Description() string { return i.description }
func (i commandDefinition) FilterValue() string { return i.name }

var commandDefinitions = []commandDefinition{
	{"bandpass", "Apply bandpass filtering to the image.", []string{"lowCut", "highCut", "withSpectrumImgGenerated"}},
	{"lowpass", "Apply lowpass filtering to the image.", []string{"cutoff", "withSpectrumImgGenerated"}},
	{"highpass", "Apply highpass filtering to the image.", []string{"cutoff", "withSpectrumImgGenerated"}},
	{"bandcut", "Apply bandcut filtering to the image.", []string{"lowCut", "highCut", "withSpectrumImgGenerated"}},
	{"phasemod", "Apply phase modification to the image.", []string{"k", "l"}},
	{"maskpass", "Apply maskpass filtering to the image.", []string{"maskName", "withSpectrumImgGenerated"}},
}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func clearSuccessAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearSuccessMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m *model) initializeTextInputs() {
	m.inputs = nil
	m.selectedCommandArgs = make(map[string]string)

	for _, cmd := range commandDefinitions {
		if cmd.name == m.selectedCommand {
			m.inputs = make([]textinput.Model, len(cmd.args))
			for j, arg := range cmd.args {
				input := textinput.New()
				input.Placeholder = arg
				m.selectedCommandArgs[arg] = ""
				m.inputs[j] = input
			}

			m.cursor = 0
			m.inputs[0].Focus()
			break
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.currentView == filePickerView {
		m.filepicker, cmd = m.filepicker.Update(msg)
	} else if m.currentView == commandsSelectionView {
		m.commandsList, cmd = m.commandsList.Update(msg)
	} else if m.currentView == commandDetailView {
		for i, input := range m.inputs {
			m.inputs[i], cmd = input.Update(msg)
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "tab":
			m.currentView = (m.currentView + 1) % 5
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
				if selectedItem, ok := m.commandsList.SelectedItem().(commandDefinition); ok {
					m.selectedCommand = selectedItem.name
					m.initializeTextInputs()
					m.currentView = commandDetailView
				}

			}

			if m.currentView == commandDetailView {

				if m.cursor == len(m.inputs) {

					for _, input := range m.inputs {
						argName := input.Placeholder
						argValue := input.Value()
						m.selectedCommandArgs[argName] = argValue
					}

					for key, value := range m.selectedCommandArgs {
						if strings.TrimSpace(value) == "" {
							m.err = fmt.Errorf("missing value for %s", key)
							return m, clearErrorAfter(2 * time.Second)
						}
					}

					result, err := executeCommand(m.selectedFile, m.selectedCommand, m.selectedCommandArgs)

					log.Default().Printf("Command result: %s", result)

					if err != nil {
						log.Default().Printf("Error executing command: %v", err)
						m.err = err
						return m, clearErrorAfter(3 * time.Second)
					}

					m.successMessage = result
					return m, clearSuccessAfter(3 * time.Second)

				} else {
					m.cursor = (m.cursor + 1) % (len(m.inputs) + 1)

					for i := range m.inputs {
						if i == m.cursor {
							m.inputs[i].Focus()
						} else {
							m.inputs[i].Blur()
						}
					}
				}
			}

		case "up":
			if m.currentView == commandDetailView {
				if m.cursor > 0 {
					m.cursor--
				}
				m.cursor = m.cursor % len(m.inputs)
				for i := range m.inputs {
					if i == m.cursor {
						m.inputs[i].Focus()
					} else {
						m.inputs[i].Blur()
					}
				}
			}

		case "down":
			if m.currentView == commandDetailView {
				m.cursor = (m.cursor + 1) % (len(m.inputs) + 1)
				for i := range m.inputs {
					if i == m.cursor && i < len(m.inputs) {
						m.inputs[i].Focus()
					} else {
						m.inputs[i].Blur()
					}
				}
			}

		}

	case tea.WindowSizeMsg:
		m.terminalSize.width = msg.Width
		m.terminalSize.height = msg.Height
		m.commandsList.SetSize(msg.Width, msg.Height)

		m.commandDetailViewStyle = lipgloss.NewStyle().
			Width(msg.Width-2).
			Height(msg.Height-2).
			Padding(1, 2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205"))

		m.filePickerViewStyle = lipgloss.NewStyle().
			Width(msg.Width - 2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205"))

		return m, nil

	case clearErrorMsg:
		m.err = nil

	case clearSuccessMsg:
		m.successMessage = ""
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

		titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)

		title := titleStyle.Render("Select Files (Press 'Tab' to switch views, 'Enter' to select):")

		var errorMessage string
		if m.err != nil {
			errorMessage = errorStyle.Render("Error: " + m.err.Error())
		}

		content := lipgloss.JoinVertical(lipgloss.Top,
			title,
			errorMessage,
			m.filepicker.View(),
		)

		return m.filePickerViewStyle.Render(content)

	case imagePreviewView:
		s.WriteString("\n  Image Preview (Press 'Tab' to view selected files):\n")
		s.WriteString(m.imagePreview)

	case commandsSelectionView:
		s.WriteString("\n  Select a command (Press 'Tab' to view selected files, 'Enter' to enter for details):\n")
		s.WriteString(m.commandsList.View())

	case commandDetailView:
		s.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render("Command Details"))
		s.WriteString("\n\n")

		s.WriteString(lipgloss.NewStyle().Bold(true).Render("Command: "))
		s.WriteString(m.selectedCommand + "\n\n")

		for _, input := range m.inputs {
			s.WriteString("\n" + input.View())
		}

		buttonStyle := lipgloss.NewStyle().
			Padding(0, 4).
			Background(lipgloss.Color("205")).
			Foreground(lipgloss.Color("0")).
			Bold(true)

		if m.cursor == len(m.inputs) {
			s.WriteString("\n\n" + buttonStyle.Render("[ Submit ]"))
		} else {
			s.WriteString("\n\n" + lipgloss.NewStyle().Render("[ Submit ]"))
		}

		if m.err != nil {
			errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
			s.WriteString("\n\n" + errorStyle.Render("Error: "+m.err.Error()))
		} else if m.successMessage != "" {
			successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
			s.WriteString("\n\n" + successStyle.Render("Success: "+m.successMessage))
		}

		return m.commandDetailViewStyle.Render(s.String())
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
	for _, cmd := range commandDefinitions {
		commandItems = append(commandItems, cmd)
	}
	commandList := list.New(commandItems, list.NewDefaultDelegate(), 0, 0)
	commandList.Title = "Available Commands"

	var commandDetailStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))

	var filePickerViewStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))

	m := model{
		filepicker:             fp,
		currentView:            filePickerView,
		commandsList:           commandList,
		selectedCommandArgs:    make(map[string]string),
		commandDetailViewStyle: commandDetailStyle,
		filePickerViewStyle:    filePickerViewStyle,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nSee you next time " + m.filepicker.Styles.Selected.Render("👋"))
}
