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
	{"negative", "Apply negative transformation to the image.", []string{}},
	{"flip_horizontally", "Flip the image horizontally.", []string{}},
	{"flip_vertically", "Flip the image vertically.", []string{}},
	{"flip_diagonally", "Flip the image diagonally.", []string{}},
	{"shrink", "Shrink the image by given factor.", []string{"shrinkFactor"}},
	{"enlarge", "Enlarge the image by given factor.", []string{"enlargeFactor"}},
	{"adaptive_filter_denoising", "Apply adaptive median noise removal filter to the image.", []string{"minWindowSize", "maxWindowSize"}},
	{"min_filter_denoising", "Apply min noise removal filter to the image.", []string{"minWindowSize"}},
	{"max_filter_denoising", "Apply max noise removal filter to the image.", []string{"maxWindowSize"}},
	{"img_comparison_commands", "Compare the image with another image.", []string{"dummy"}},
	{"generate_img_histogram", "Generate and save a graphical representation of the histogram of the image.", []string{"dummy"}},
	{"histogram_img_characteristics", "Calculate image characteristics based on it's histogram", []string{"dummy"}},
	{"rayleigh_transform", "Apply Rayleigh transform to the image.", []string{"lowCut", "highCut", "alphaValue"}},
	{"mask_edge_sharpening", "Apply edge sharpening mask to the image.", []string{"maskName"}},
	{"kirsh_edge_detection", "Apply Kirsh edge detection to the image.", []string{"dummy"}},

	{"bandpass", "Apply bandpass filtering to the image.", []string{"lowCut", "highCut", "withSpectrumImgGenerated"}},
	{"lowpass", "Apply lowpass filtering to the image.", []string{"cutoff", "withSpectrumImgGenerated"}},
	{"highpass", "Apply highpass filtering to the image.", []string{"cutoff", "withSpectrumImgGenerated"}},
	{"bandcut", "Apply bandcut filtering to the image.", []string{"lowCut", "highCut", "withSpectrumImgGenerated"}},
	{"phasemod", "Apply phase modification to the image.", []string{"k", "l"}},
	{"maskpass", "Apply maskpass filtering to the image.", []string{"maskName", "withSpectrumImgGenerated"}},
}
