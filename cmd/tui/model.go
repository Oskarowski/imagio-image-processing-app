package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type clearErrorMsg struct{}
type clearSuccessMsg struct{}

type appView int

const (
	FILE_PICKER_VIEW appView = iota
	IMAGE_PREVIEW_VIEW
	COMMAND_SELECTION_VIEW
	COMMAND_EXECUTION_VIEW
)

type UIStyle struct {
	filepickerViewStyle       lipgloss.Style
	commandExecutionViewStyle lipgloss.Style
}

type UIState struct {
	err            error
	successMessage string
	imagePreview   string
}

type CommandState struct {
	selectedCommand string
	commandArgs     map[string]string
	inputs          []textinput.Model
}
