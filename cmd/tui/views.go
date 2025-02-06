package tui

import (
	"image-processing/imageio"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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

	return s.String()

}
