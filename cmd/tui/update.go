package tui

import (
	"errors"
	"fmt"
	"image-processing/cmd/tui/executioner"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) updateFilePickerView(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)
	return m, cmd
}

func (m Model) updateCommandSelectionView(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.commandsList, cmd = m.commandsList.Update(msg)
	return m, cmd
}

func (m Model) updateCommandExecutionView(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	for i, input := range m.CommandState.inputs {
		m.CommandState.inputs[i], cmd = input.Update(msg)
	}
	return m, cmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.currentView {
	case FILE_PICKER_VIEW:
		m, cmd = m.updateFilePickerView(msg)
	case COMMAND_SELECTION_VIEW:
		m, cmd = m.updateCommandSelectionView(msg)
	case COMMAND_EXECUTION_VIEW:
		m, cmd = m.updateCommandExecutionView(msg)
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

		case "shift+tab":
			m.currentView = (m.currentView - 1 + 4) % 4

			return m, nil

		case "enter":
			switch m.currentView {

			case FILE_PICKER_VIEW:

				if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
					m.selectedFile = path
					m.loadImagePreview(path)
					m.currentView = IMAGE_PREVIEW_VIEW

					return m, nil
				}

				if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
					m.UIState.err = errors.New(path + " is not valid.")
					m.selectedFile = ""
					return m, clearErrorAfter(2 * time.Second)
				}

			case COMMAND_SELECTION_VIEW:

				if selectedItem, ok := m.commandsList.SelectedItem().(executioner.CommandDefinition); ok {
					m.CommandState.selectedCommand = selectedItem.Name
					m.initializeTextInputs()
					m.currentView = COMMAND_EXECUTION_VIEW
				}

			case COMMAND_EXECUTION_VIEW:

				if m.cursor == len(m.CommandState.inputs) {
					// mv this to a separate function which will handle this

					for _, formInput := range m.CommandState.inputs {
						argName := formInput.Placeholder
						argValue := formInput.Value()
						m.CommandState.commandArgs[argName] = argValue
					}

					// for future the command arg should have a flag to indicate if it is required and if not provide a default value
					for key, value := range m.CommandState.commandArgs {
						if strings.TrimSpace(value) == "" {
							m.UIState.err = fmt.Errorf("missing value for %s", key)
							return m, clearErrorAfter(2 * time.Second)
						}
					}

					result, err := executioner.ExecuteCommand(m.selectedFile, m.CommandState.selectedCommand, m.CommandState.commandArgs)

					if err != nil {
						m.UIState.err = err
						return m, clearErrorAfter(3 * time.Second)
					}

					m.UIState.successMessage = result
					return m, clearSuccessAfter(3 * time.Second)

				} else {
					m.cursor = (m.cursor + 1) % (len(m.CommandState.inputs) + 1)

					for i := range m.CommandState.inputs {
						if i == m.cursor {
							m.CommandState.inputs[i].Focus()
						} else {
							m.CommandState.inputs[i].Blur()
						}
					}
				}

			}

		case "up":
			switch m.currentView {

			case COMMAND_EXECUTION_VIEW:

				if len(m.CommandState.inputs) > 0 {
					if m.cursor > 0 {
						m.cursor--
					}

					m.cursor = m.cursor % len(m.CommandState.inputs)

					for i := range m.CommandState.inputs {
						if i == m.cursor {
							m.CommandState.inputs[i].Focus()
						} else {
							m.CommandState.inputs[i].Blur()
						}
					}
				}

			}

		case "down":
			switch m.currentView {

			case COMMAND_EXECUTION_VIEW:

				if len(m.CommandState.inputs) > 0 {
					m.cursor = (m.cursor + 1) % (len(m.CommandState.inputs) + 1)
					for i := range m.CommandState.inputs {
						if i == m.cursor && i < len(m.CommandState.inputs) {
							m.CommandState.inputs[i].Focus()
						} else {
							m.CommandState.inputs[i].Blur()
						}
					}
				}
			}

		}

	case tea.WindowSizeMsg:
		m.terminalSize.width = msg.Width
		m.terminalSize.height = msg.Height

		navBarHeight := lipgloss.Height(m.renderNavBar())
		containerPadding := 4
		availableHeight := msg.Height - navBarHeight - containerPadding
		m.commandsList.SetSize(msg.Width-4, availableHeight)

		m.filepicker.Height = availableHeight - 4

		return m, nil

	case clearErrorMsg:
		m.UIState.err = nil

	case clearSuccessMsg:
		m.UIState.successMessage = ""
	}

	return m, cmd
}

// TODO: Implement the following functions
// func (m Model) handleEnterKey(msg tea.Msg) (tea.Model, tea.Cmd) {

// }
