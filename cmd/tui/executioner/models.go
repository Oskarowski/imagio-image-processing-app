package executioner

type CommandDefinition struct {
	Name string
	Desc string
	Args []string
}

func (i CommandDefinition) Title() string       { return i.Name }
func (i CommandDefinition) Description() string { return i.Desc }
func (i CommandDefinition) FilterValue() string { return i.Name }

var CommandDefinitions = []CommandDefinition{
	{"brightness", "Adjust brightness of the image by given percentage. ", []string{"brightness"}},
	{"contrast", "Adjust contrast of the image by given value", []string{"contrast"}},
	{"bandpass", "Apply bandpass filtering to the image.", []string{"lowCut", "highCut", "withSpectrumImgGenerated"}},
	{"lowpass", "Apply lowpass filtering to the image.", []string{"cutoff", "withSpectrumImgGenerated"}},
	{"highpass", "Apply highpass filtering to the image.", []string{"cutoff", "withSpectrumImgGenerated"}},
	{"bandcut", "Apply bandcut filtering to the image.", []string{"lowCut", "highCut", "withSpectrumImgGenerated"}},
	{"phasemod", "Apply phase modification to the image.", []string{"k", "l"}},
	{"maskpass", "Apply maskpass filtering to the image.", []string{"maskName", "withSpectrumImgGenerated"}},
}
