package executioner

import (
	"errors"
	"fmt"
	"image-processing/imageio"
	"image-processing/orthogonal_transforms"
)

type handlingCommandOptions struct {
	imgPath                  string
	maskPath                 string
	lowCut, highCut          int
	cutoff, k, l             int
	withSpectrumImgGenerated bool
}

func validateFrequencyRange(lowCut, highCut, imgWidth int) error {
	if lowCut < 0 || highCut < 0 {
		return errors.New("low and high cut frequencies must be positive")
	}

	if lowCut >= highCut {
		return errors.New("low cut frequency must be lower than high cut frequency")
	}

	if lowCut >= imgWidth/2 || highCut >= imgWidth/2 {
		return errors.New("low and high cut frequencies must be lower than half the image width")
	}

	return nil
}

func saveFilteringResults(handlerResult []orthogonal_transforms.SpectrumImage) error {
	for _, resultImg := range handlerResult {
		err := imageio.SaveBmpImage(&resultImg.Img, resultImg.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleBandpassCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	if err := validateFrequencyRange(opts.lowCut, opts.highCut, img.Bounds().Dx()); err != nil {
		return "", err
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)

	handlerResult := orthogonal_transforms.HandleBandpassFiltering(img, imgFileName, opts.lowCut, opts.highCut, opts.withSpectrumImgGenerated)

	if err := saveFilteringResults(handlerResult); err != nil {
		return "", err
	}

	return "Bandpass filtering successful", nil
}

func handleLowpassCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	if opts.cutoff < 0 || opts.cutoff >= img.Bounds().Dx()/2 || opts.cutoff >= img.Bounds().Dy()/2 {
		return "", errors.New("cutoff frequency must be positive and less than half the image width/height")
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)

	handlerResult := orthogonal_transforms.HandleLowpassFiltering(img, imgFileName, opts.cutoff, opts.withSpectrumImgGenerated)

	if err := saveFilteringResults(handlerResult); err != nil {
		return "", err
	}

	return "Lowpass filtering successful", nil
}

func handleHighpassCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	if opts.cutoff < 0 || opts.cutoff >= img.Bounds().Dx()/2 || opts.cutoff >= img.Bounds().Dy()/2 {
		return "", errors.New("cutoff frequency must be positive and less than half the image width/height")
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)

	handlerResult := orthogonal_transforms.HandleHighpassFiltering(img, imgFileName, opts.cutoff, opts.withSpectrumImgGenerated)

	if err := saveFilteringResults(handlerResult); err != nil {
		return "", err
	}

	return "Highpass filtering successful", nil
}

func handleBandcutCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	if err := validateFrequencyRange(opts.lowCut, opts.highCut, img.Bounds().Dx()); err != nil {
		return "", err
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)

	handlerResult := orthogonal_transforms.HandleBandcutFiltering(img, imgFileName, opts.lowCut, opts.highCut, opts.withSpectrumImgGenerated)

	if err := saveFilteringResults(handlerResult); err != nil {
		return "", err
	}

	return "Bandcut filtering successful", nil
}

func handlePhasemodCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	if opts.k < 0 || opts.l < 0 {
		return "", errors.New("k and l must be positive")
	}

	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)

	handlerResult := orthogonal_transforms.HandlePhaseModification(img, imgFileName, opts.k, opts.l)

	if err := saveFilteringResults(handlerResult); err != nil {
		return "", err
	}

	return "Phase modified successful", nil
}

func handleMaskpassCommand(opts handlingCommandOptions) (successMsgString string, err error) {

	maskImg, err := imageio.OpenBmpImage(opts.maskPath)
	if err != nil {
		return "", fmt.Errorf("failed to open mask image: %w", err)
	}

	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)

	handlerResult := orthogonal_transforms.HandleMaskpassFiltering(img, imgFileName, maskImg, opts.withSpectrumImgGenerated)

	if err := saveFilteringResults(handlerResult); err != nil {
		return "", err
	}

	return "Maskpass filter successful applied", nil
}
