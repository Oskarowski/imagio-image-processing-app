package tui

import (
	"image-processing/cmd/executioner"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

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

func buildCommandListItems(cmds []executioner.CommandDefinition) []list.Item {
	items := make([]list.Item, len(cmds))
	for i, cmd := range cmds {
		items[i] = cmd
	}
	return items
}
