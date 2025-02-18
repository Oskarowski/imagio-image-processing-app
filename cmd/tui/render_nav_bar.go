package tui

import "github.com/charmbracelet/lipgloss"

func (m Model) renderNavBar() string {
	navItems := []struct {
		label string
		view  appView
	}{
		{"File Picker", FILE_PICKER_VIEW},
		{"Image Preview", IMAGE_PREVIEW_VIEW},
		{"Command Selection", COMMAND_SELECTION_VIEW},
		{"Command Execution", COMMAND_EXECUTION_VIEW},
	}

	var renderedItems []string
	for _, item := range navItems {
		style := lipgloss.NewStyle().Padding(1, 1)

		if m.currentView == item.view {
			style = style.Foreground(lipgloss.Color("205")).Bold(true).Underline(true)
		} else {
			style = style.Foreground(lipgloss.Color("240"))
		}

		renderedItems = append(renderedItems, style.Render(item.label))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedItems...)
}
