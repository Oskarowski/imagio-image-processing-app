package tui

import (
	"time"

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
