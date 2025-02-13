package tui

import (
	"errors"
	"image-processing/analysis"
	"image-processing/cmd/tui/executioner"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
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

	if m.form != nil {
		var formCmd tea.Cmd
		form, formCmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			cmd = tea.Batch(cmd, formCmd)
		}

		return m, cmd
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

					formErr := m.buildCommandForm()
					if formErr != nil {
						m.UIState.err = formErr
						return m, clearErrorAfter(3 * time.Second)
					}

					m.currentView = COMMAND_EXECUTION_VIEW
				}

			case COMMAND_EXECUTION_VIEW:

				if m.form != nil {
					m.CommandState.commandArgs = m.formGetter()

					for key, value := range m.CommandState.commandArgs {
						if strings.Trim(value, " ") == "" {
							m.UIState.err = errors.New("missing value for " + key)
							return m, clearErrorAfter(2 * time.Second)
						}

						result := executioner.ExecuteCommand(m.selectedFile, m.CommandState.selectedCommand, m.CommandState.commandArgs)

						if result.Err != nil {
							m.UIState.err = result.Err
							return m, clearErrorAfter(3 * time.Second)
						}

						if entries, ok := result.Output.([]analysis.CharacteristicsEntry); ok && len(entries) > 0 {
							m.UIState.imageComparisonResults = entries
						} else {
							m.UIState.imageComparisonResults = nil
						}

						m.UIState.successMessage = result.Message
						return m, clearSuccessAfter(3 * time.Second)
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
