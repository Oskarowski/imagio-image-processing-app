package tui

import "image-processing/analysis"

type clearErrorMsg struct{}
type clearSuccessMsg struct{}

type appView int

const (
	FILE_PICKER_VIEW appView = iota
	IMAGE_PREVIEW_VIEW
	COMMAND_SELECTION_VIEW
	COMMAND_EXECUTION_VIEW
)

type UIState struct {
	err                    error
	successMessage         string
	imageComparisonResults []analysis.CharacteristicsEntry
}

type CommandState struct {
	selectedCommand string
	commandArgs     map[string]string
}
