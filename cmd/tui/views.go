package tui

import (
	"imagio/imageio"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) currentContentView() string {
	switch m.currentView {
	case FILE_PICKER_VIEW:
		return m.filePickerView()
	case IMAGE_PREVIEW_VIEW:
		return m.imagePreviewView()
	case COMMAND_SELECTION_VIEW:
		return m.commandSelectionView()
	case COMMAND_EXECUTION_VIEW:
		return m.commandExecutionView()
	default:
		return "Unknown view"
	}
}

func (m Model) filePickerView() string {
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

	return content
}

func (m Model) imagePreviewView() string {
	var s strings.Builder
	s.WriteString(m.imagePreview)
	return s.String()
}

func (m Model) commandSelectionView() string {
	var s strings.Builder
	s.WriteString(m.commandsList.View())
	return s.String()
}

func (m Model) commandExecutionView() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	labelStyle := lipgloss.NewStyle().Bold(true)

	var s strings.Builder

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

	if m.form == nil {
		s.WriteString("No form defined.\n")
		return s.String()
	}

	formContainer := lipgloss.NewStyle().
		Render(m.form.View())

	s.WriteString(formContainer)

	if m.CommandState.selectedCommand != "" && m.selectedFile != "" {
		executeMsgStyle := lipgloss.NewStyle().
			Padding(0, 4).
			Background(lipgloss.Color("205")).
			Foreground(lipgloss.Color("0")).
			Bold(true)
		s.WriteString("\n\n" + executeMsgStyle.Render("Press Enter to Execute Command"))
	} else {
		disabledStyle := lipgloss.NewStyle().
			Padding(0, 4).
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("250")).
			Bold(true)
		s.WriteString("\n\n" + disabledStyle.Render("[ Execute ]"))
	}

	if m.UIState.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
		s.WriteString("\n\n" + errorStyle.Render("Error: "+m.UIState.err.Error()))
	} else if m.UIState.successMessage != "" {
		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
		s.WriteString("\n\n" + successStyle.Render("Success: "+m.UIState.successMessage))
	}

	if len(m.UIState.imageComparisonResults) > 0 {
		s.WriteString("\n\n" + m.renderComparisonResults(m.UIState.imageComparisonResults))
	}

	return s.String()

}
