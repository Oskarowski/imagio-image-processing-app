package tui

import (
	"fmt"
	"image-processing/orthogonal_transforms"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

func (m *Model) buildCommandForm() error {
	var (
		lowCut, highCut, cutoff, k, l, maskName, brightness, contrast, shrinkFactor, enlargeFactor, minWindowSize, maxWindowSize string
		withSpectrum                                                                                                             bool
	)

	customKM := huh.NewDefaultKeyMap()
	customKM.Input.Next = key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "Next field"))
	customKM.Input.Prev = key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "Previous field"))
	customKM.Confirm.Next = key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "Next field"))
	customKM.Confirm.Prev = key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "Previous field"))
	customKM.Select.Next = key.NewBinding(key.WithKeys("p", "j"), key.WithHelp("p", "Select"))
	customKM.Select.Prev = key.NewBinding(key.WithKeys("k"), key.WithHelp("k", "Previous"))

	customKM.Note.Submit.SetEnabled(false)
	customKM.Note.Next.SetEnabled(false)
	customKM.Note.Prev.SetEnabled(false)

	customKM.Confirm.Submit.SetEnabled(false)
	customKM.Confirm.Accept.SetEnabled(false)
	customKM.Confirm.Reject.SetEnabled(false)

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
