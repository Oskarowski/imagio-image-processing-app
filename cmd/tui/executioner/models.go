package executioner

type CommandDefinition struct {
	Name string
	Desc string
	// TODO make this an array of structs implementing .Value() interface to allow creation of checkboxes and options selectors
	Args []string
}

func (i CommandDefinition) Title() string       { return i.Name }
func (i CommandDefinition) Description() string { return i.Desc }
func (i CommandDefinition) FilterValue() string { return i.Name }

var CommandDefinitions = []CommandDefinition{
	{"bandpass", "Apply bandpass filtering to the image.", []string{"lowCut", "highCut", "withSpectrumImgGenerated"}},
	{"lowpass", "Apply lowpass filtering to the image.", []string{"cutoff", "withSpectrumImgGenerated"}},
	{"highpass", "Apply highpass filtering to the image.", []string{"cutoff", "withSpectrumImgGenerated"}},
	{"bandcut", "Apply bandcut filtering to the image.", []string{"lowCut", "highCut", "withSpectrumImgGenerated"}},
	{"phasemod", "Apply phase modification to the image.", []string{"k", "l"}},
	{"maskpass", "Apply maskpass filtering to the image.", []string{"maskName", "withSpectrumImgGenerated"}},
}
