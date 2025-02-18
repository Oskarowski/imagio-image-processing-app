package tui

import (
	"fmt"
	"imagio/manipulations"
	"imagio/morphological"
	"imagio/orthogonal_transforms"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

func (m *Model) buildCommandForm() error {
	var (
		lowCut, highCut, cutoff, k, l, maskName, brightness, contrast, shrinkFactor, enlargeFactor, minWindowSize, maxWindowSize, comparisonImagePath, alpha, structureElementName, foregroundStructureElementName, backgroundStructureElementName, seedPointsStr, thresholdStr string
		selectedComparisonCommands, selectedHistogramCharacteristicsCommands                                                                                                                                                                                                    []string
		withSpectrum                                                                                                                                                                                                                                                            bool
		distanceMetric                                                                                                                                                                                                                                                          int
	)

	customKM := huh.NewDefaultKeyMap()
	customKM.Input.Next = key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "Next field"))
	customKM.Input.Prev = key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "Previous field"))
	customKM.Confirm.Next = key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "Next field"))
	customKM.Confirm.Prev = key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "Previous field"))

	customKM.Select.Next = key.NewBinding(key.WithKeys("p", "j"), key.WithHelp("p", "Select"))
	customKM.Select.Prev = key.NewBinding(key.WithKeys("k"), key.WithHelp("k", "Previous"))
	customKM.Select.Down = key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "down"))
	customKM.Select.Up = key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "up"))

	customKM.Note.Submit.SetEnabled(false)
	customKM.Note.Next.SetEnabled(false)
	customKM.Note.Prev.SetEnabled(false)

	customKM.Confirm.Submit.SetEnabled(false)
	customKM.Confirm.Accept.SetEnabled(false)
	customKM.Confirm.Reject.SetEnabled(false)

	customKM.FilePicker.Open = key.NewBinding(key.WithKeys("p", "right"), key.WithHelp("p", "Pick file"))
	customKM.FilePicker.Select = key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "Pick file"))
	customKM.FilePicker.Up = key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "up"))
	customKM.FilePicker.Down = key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "down"))
	customKM.FilePicker.Next = key.NewBinding(key.WithKeys("j"), key.WithHelp("j", "Next field"))
	customKM.FilePicker.Prev = key.NewBinding(key.WithKeys("k"), key.WithHelp("k", "Previous field"))

	customKM.MultiSelect.Next = key.NewBinding(key.WithKeys("j"), key.WithHelp("j", "Next field"))
	customKM.MultiSelect.Prev = key.NewBinding(key.WithKeys("k"), key.WithHelp("k", "Previous field"))
	customKM.MultiSelect.Up = key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "up"))
	customKM.MultiSelect.Down = key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "down"))

	cmd := m.CommandState.selectedCommand

	var form *huh.Form

	switch cmd {
	case "brightness":

		inputBrightness := huh.NewInput().
			Title("Brightness percentage, can be negative to decrease brightness").
			Placeholder("Enter brightness percentage").
			Value(&brightness)

		form = huh.NewForm(huh.NewGroup(inputBrightness)).WithTheme(huh.ThemeCatppuccin())

	case "contrast":

		inputContrast := huh.NewInput().
			Title("Contrast adjustment value").
			Placeholder("Enter contrast adjustment value").
			Value(&contrast)

		form = huh.NewForm(huh.NewGroup(inputContrast)).WithTheme(huh.ThemeCatppuccin())

	case "negative", "flip_horizontally", "flip_vertically", "flip_diagonally":

		dummyNote := huh.NewNote().Title("No arguments required for this command")

		form = huh.NewForm(huh.NewGroup(dummyNote)).WithTheme(huh.ThemeCatppuccin())

	case "shrink":

		inputShrinkFactor := huh.NewInput().
			Title("Shrink factor").
			Placeholder("Enter shrink factor").
			Value(&shrinkFactor)

		form = huh.NewForm(huh.NewGroup(inputShrinkFactor)).WithTheme(huh.ThemeCatppuccin())

	case "enlarge":

		inputEnlargeFactor := huh.NewInput().
			Title("Enlarge factor").
			Placeholder("Enter enlarge factor").
			Value(&enlargeFactor)

		form = huh.NewForm(huh.NewGroup(inputEnlargeFactor)).WithTheme(huh.ThemeCatppuccin())

	case "adaptive_filter_denoising":

		inputMinWindowSize := huh.NewInput().
			Title("Minimal size of window size for filter").
			Placeholder("Enter minimal size of window size for filter").
			Value(&minWindowSize)

		inputMaxWindowSize := huh.NewInput().
			Title("Maximal size of window size for filter").
			Placeholder("Enter maximal size of window size for filter").
			Value(&maxWindowSize)

		form = huh.NewForm(huh.NewGroup(inputMinWindowSize, inputMaxWindowSize)).WithTheme(huh.ThemeCatppuccin())

	case "min_filter_denoising":

		inputMinWindowSize := huh.NewInput().
			Title("Minimal size of window size for filter").
			Placeholder("Enter minimal size of window size for filter").
			Value(&minWindowSize)

		form = huh.NewForm(huh.NewGroup(inputMinWindowSize)).WithTheme(huh.ThemeCatppuccin())

	case "max_filter_denoising":

		inputMaxWindowSize := huh.NewInput().
			Title("Maximal size of window size for filter").
			Placeholder("Enter maximal size of window size for filter").
			Value(&maxWindowSize)

		form = huh.NewForm(huh.NewGroup(inputMaxWindowSize)).WithTheme(huh.ThemeCatppuccin())

	case "img_comparison_commands":
		wd, _ := os.Getwd()

		fpComparison := huh.NewFilePicker().
			Title("Select comparison image").
			AllowedTypes([]string{".bmp", ".png"}).
			Value(&comparisonImagePath).
			CurrentDirectory(wd)

		comparisonOptions := []huh.Option[string]{
			huh.NewOption("MSE (Mean Square Error)", "MSE"),
			huh.NewOption("PMSE (Peak Mean Square Error)", "PMSE"),
			huh.NewOption("SNR (Signal-to-Noise Ratio)", "SNR"),
			huh.NewOption("PSNR (Peak Signal-to-Noise Ratio)", "PSNR"),
			huh.NewOption("MD (Max Difference)", "MD"),
		}

		msComparison := huh.NewMultiSelect[string]().
			Title("Select comparison commands").
			Options(comparisonOptions...).
			Value(&selectedComparisonCommands)

		form = huh.NewForm(huh.NewGroup(fpComparison, msComparison)).WithTheme(huh.ThemeCatppuccin())

	case "generate_img_histogram":

		dummyNote := huh.NewNote().Title("No arguments required for this command")

		form = huh.NewForm(huh.NewGroup(dummyNote)).WithTheme(huh.ThemeCatppuccin())

	case "histogram_img_characteristics":

		characteristicsOptions := []huh.Option[string]{
			huh.NewOption("Mean", "cmean"),
			huh.NewOption("Variance intensity", "cvariance"),
			huh.NewOption("Standard deviation", "cstdev"),
			huh.NewOption("Coefficient of variation I", "cvarcoi"),
			huh.NewOption("Asymmetry coefficient", "casyco"),
			huh.NewOption("Flattening coefficient", "cflatco"),
			huh.NewOption("Coefficient of variation II", "cvarcoii"),
			huh.NewOption("Entropy", "centropy"),
		}

		msCharacteristics := huh.NewMultiSelect[string]().
			Title("Select characteristics to calculate").
			Options(characteristicsOptions...).
			Value(&selectedHistogramCharacteristicsCommands)

		form = huh.NewForm(huh.NewGroup(msCharacteristics)).WithTheme(huh.ThemeCatppuccin())

	case "rayleigh_transform":

		inputMinBrightness := huh.NewInput().
			Title("Minimum value in the range [0, 255]").
			Placeholder("Enter minimum value").
			Value(&lowCut)

		inputMaxBrightness := huh.NewInput().
			Title("Maximum value in the range [0, 255], must be greater than min").
			Placeholder("Enter maximum value").
			Value(&highCut)

		inputAlpha := huh.NewInput().
			Title("Alpha value for transformation").
			Placeholder("Enter alpha value").
			Value(&alpha)

		form = huh.NewForm(huh.NewGroup(inputMinBrightness, inputMaxBrightness, inputAlpha)).WithTheme(huh.ThemeCatppuccin())

	case "kirsh_edge_detection", "thinning":

		dummyNote := huh.NewNote().Title("No arguments required for this command")

		form = huh.NewForm(huh.NewGroup(dummyNote)).WithTheme(huh.ThemeCatppuccin())

	case "mask_edge_sharpening":

		availableMasks, err := manipulations.GetAvailableEdgeSharpeningMasksNames()
		if err != nil {
			return fmt.Errorf("failed to get available edge sharpening masks: %w", err)
		}

		maskOptions := huh.NewOptions(availableMasks...)

		selectMask := huh.NewSelect[string]().
			Title("Mask Name").
			Options(maskOptions...).
			Value(&maskName)

		form = huh.NewForm(huh.NewGroup(selectMask)).WithTheme(huh.ThemeCatppuccin())

	case "dilation", "erosion", "opening", "closing":

		availableStructuringElements, err := morphological.GetAvailableStructureElementsNames()
		if err != nil {
			return fmt.Errorf("failed to get available structuring elements: %w", err)
		}

		seOptions := huh.NewOptions(availableStructuringElements...)

		selectSE := huh.NewSelect[string]().
			Title("Structuring Element Name").
			Options(seOptions...).
			Value(&structureElementName)

		form = huh.NewForm(huh.NewGroup(selectSE)).WithTheme(huh.ThemeCatppuccin())

	case "hit_or_miss":

		availableStructuringElements, err := morphological.GetAvailableStructureElementsNames()
		if err != nil {
			return fmt.Errorf("failed to get available structuring elements: %w", err)
		}
		seOptions := huh.NewOptions(availableStructuringElements...)

		selectForegroundSE := huh.NewSelect[string]().
			Title("Foreground Structuring Element Name").
			Options(seOptions...).
			Value(&foregroundStructureElementName)

		selectBackgroundSE := huh.NewSelect[string]().
			Title("Background Structuring Element Name").
			Options(seOptions...).
			Value(&backgroundStructureElementName)

		form = huh.NewForm(huh.NewGroup(selectForegroundSE, selectBackgroundSE)).WithTheme(huh.ThemeCatppuccin())

	case "region_grow":

		//TODO enhance this part of adding seed points
		seedPointInput := huh.NewInput().
			Title("Seed points").
			Placeholder("Enter seeds as [x,y][x,y]...").
			Value(&seedPointsStr)

		distanceMetricOptions := []huh.Option[int]{
			huh.NewOption("Euclidean (0)", 0),
			huh.NewOption("Manhattan (1)", 1),
			huh.NewOption("Chebyshev (2)", 2),
		}

		distanceMetricSelect := huh.NewSelect[int]().
			Title("Distance metric").
			Options(distanceMetricOptions...).
			Value(&distanceMetric)

		thresholdInput := huh.NewInput().
			Title("Threshold").
			Placeholder("Enter similarity threshold (e.g. 20.0)").
			Value(&thresholdStr).
			Validate(func(s string) error {
				t, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return fmt.Errorf("failed to parse threshold value: %w", err)
				}
				if t < 0 {
					return fmt.Errorf("threshold value must be greater than 0")
				}
				return nil
			})

		form = huh.NewForm(huh.NewGroup(seedPointInput, distanceMetricSelect, thresholdInput)).WithTheme(huh.ThemeCatppuccin())

	case "bandpass", "bandcut":

		inputLowCut := huh.NewInput().
			Title("Low cut frequency").
			Placeholder("Enter low cut frequency value").
			Value(&lowCut)

		inputHighCut := huh.NewInput().
			Title("High Cut").
			Placeholder("Enter high cut frequency value").
			Value(&highCut)

		confirmSpectrum := huh.NewConfirm().
			Title("Generate Spectrum Image?").
			Affirmative("Yes").
			Negative("No").
			Value(&withSpectrum)

		form = huh.NewForm(huh.NewGroup(inputLowCut, inputHighCut, confirmSpectrum)).
			WithTheme(huh.ThemeCatppuccin())

	case "lowpass", "highpass":

		inputCutoff := huh.NewInput().
			Title("Cutoff").
			Placeholder("Enter cutoff value").
			Value(&cutoff)
		confirmSpectrum := huh.NewConfirm().
			Title("Generate Spectrum Image?").
			Affirmative("Yes").
			Negative("No").
			Value(&withSpectrum)

		form = huh.NewForm(huh.NewGroup(inputCutoff, confirmSpectrum)).
			WithTheme(huh.ThemeCatppuccin())

	case "phasemod":

		inputK := huh.NewInput().
			Title("k").
			Placeholder("Enter value for k").
			Value(&k)
		inputL := huh.NewInput().
			Title("l").
			Placeholder("Enter value for l").
			Value(&l)

		form = huh.NewForm(huh.NewGroup(inputK, inputL)).
			WithTheme(huh.ThemeCatppuccin())

	case "maskpass":

		availableMasks, err := orthogonal_transforms.GetAvailableSpectrumMasks()

		if err != nil {
			return fmt.Errorf("failed to get available masks: %w", err)
		}

		options := huh.NewOptions(availableMasks...)

		selectMask := huh.NewSelect[string]().
			Title("Mask Name").
			Options(options...).
			Value(&maskName)
		confirmSpectrum := huh.NewConfirm().
			Title("Generate Spectrum Image?").
			Affirmative("Yes").
			Negative("No").
			Value(&withSpectrum)

		form = huh.NewForm(huh.NewGroup(selectMask, confirmSpectrum)).
			WithTheme(huh.ThemeCatppuccin())

	default:
		return fmt.Errorf("unsupported command: %s", cmd)
	}

	form.WithKeyMap(customKM)
	form.Init()
	m.form = form

	m.formGetter = func() map[string]string {
		args := make(map[string]string)
		switch cmd {
		case "brightness":
			args["brightness"] = brightness
		case "contrast":
			args["contrast"] = contrast
		case "negative", "flip_horizontally", "flip_vertically", "flip_diagonally":
			// not the most elegant solution, but it is what it is
			args["dummy"] = "dummy"
		case "shrink":
			args["shrinkFactor"] = shrinkFactor
		case "enlarge":
			args["enlargeFactor"] = enlargeFactor
		case "adaptive_filter_denoising":
			args["minWindowSize"] = minWindowSize
			args["maxWindowSize"] = maxWindowSize
		case "min_filter_denoising":
			args["minWindowSize"] = minWindowSize
		case "max_filter_denoising":
			args["maxWindowSize"] = maxWindowSize
		case "img_comparison_commands":
			args["comparisonImagePath"] = comparisonImagePath
			// i love this totally not hacky type safe solution
			args["selectedComparisonCommands"] = strings.Join(selectedComparisonCommands, "|")
		case "generate_img_histogram":
			args["dummy"] = "dummy"
		case "histogram_img_characteristics":
			args["selectedHistogramCharacteristicsCommands"] = strings.Join(selectedHistogramCharacteristicsCommands, "|")
		case "rayleigh_transform":
			args["lowCut"] = lowCut
			args["highCut"] = highCut
			args["alpha"] = alpha
		case "mask_edge_sharpening":
			args["maskName"] = maskName
		case "kirsh_edge_detection":
			args["dummy"] = "dummy"
		case "dilation", "erosion", "opening", "closing":
			args["structureElementName"] = structureElementName
		case "hit_or_miss":
			args["foregroundStructureElementName"] = foregroundStructureElementName
			args["backgroundStructureElementName"] = backgroundStructureElementName
		case "thinning":
			args["dummy"] = "dummy"
		case "region_grow":
			args["seedPoints"] = seedPointsStr
			args["distanceMetric"] = strconv.Itoa(distanceMetric)
			args["threshold"] = thresholdStr
		case "bandpass", "bandcut":
			args["lowCut"] = lowCut
			args["highCut"] = highCut
			args["withSpectrumImgGenerated"] = strconv.FormatBool(withSpectrum)
		case "lowpass", "highpass":
			args["cutoff"] = cutoff
			args["withSpectrumImgGenerated"] = strconv.FormatBool(withSpectrum)
		case "phasemod":
			args["k"] = k
			args["l"] = l
		case "maskpass":
			args["maskName"] = maskName
			args["withSpectrumImgGenerated"] = strconv.FormatBool(withSpectrum)
		}
		return args
	}

	return nil
}
