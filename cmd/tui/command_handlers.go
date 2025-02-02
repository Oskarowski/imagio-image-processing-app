package tui

import (
	"errors"
	"image-processing/imageio"
	"image-processing/orthogonal_transforms"
)

func handleBandpassCommand(imgPath string, lowCut, highCut int, withSpectrumImgGenerated bool) (successMsgString string, err error) {
	if lowCut < 0 || highCut < 0 {
		return "", errors.New("low and high cut frequencies must be positive")
	}

	if lowCut >= highCut {
		return "", errors.New("low cut frequency must be lower than high cut frequency")
	}

	img, err := imageio.OpenBmpImage(imgPath)
	if err != nil {
		return "", err
	}

	if lowCut >= img.Bounds().Dx()/2 || highCut >= img.Bounds().Dx()/2 {
		return "", errors.New("low and high cut frequencies must be lower than half the image width")
	}

	imgFileName := imageio.GetPureFileName(imgPath)

	handlerResult := orthogonal_transforms.HandleBandpassFiltering(img, imgFileName, lowCut, highCut, withSpectrumImgGenerated)

	for _, resultImg := range handlerResult {
		err = imageio.SaveBmpImage(&resultImg.Img, resultImg.Name)
		if err != nil {
			return "", err
		}
	}

	return "Bandpass filtering successful", nil
}

func handleLowpassCommand(imgPath string, cutoff int, withSpectrumImgGenerated bool) (successMsgString string, err error) {
	if cutoff < 0 {
		return "", errors.New("cutoff frequency must be positive")
	}

	img, err := imageio.OpenBmpImage(imgPath)
	if err != nil {
		return "", err
	}

	if cutoff >= img.Bounds().Dx()/2 || cutoff >= img.Bounds().Dy()/2 {
		return "", errors.New("cutoff frequency must be lower than half the image width")
	}

	imgFileName := imageio.GetPureFileName(imgPath)

	handlerResult := orthogonal_transforms.HandleLowpassFiltering(img, imgFileName, cutoff, withSpectrumImgGenerated)

	for _, resultImg := range handlerResult {
		err = imageio.SaveBmpImage(&resultImg.Img, resultImg.Name)
		if err != nil {
			return "", err
		}
	}

	return "Lowpass filtering successful", nil
}

func handleHighpassCommand(imgPath string, cutoff int, withSpectrumImgGenerated bool) (successMsgString string, err error) {
	if cutoff < 0 {
		return "", errors.New("cutoff frequency must be positive")
	}

	img, err := imageio.OpenBmpImage(imgPath)
	if err != nil {
		return "", err
	}

	if cutoff >= img.Bounds().Dx()/2 || cutoff >= img.Bounds().Dy()/2 {
		return "", errors.New("cutoff frequency must be lower than half the image width")
	}

	imgFileName := imageio.GetPureFileName(imgPath)

	handlerResult := orthogonal_transforms.HandleHighpassFiltering(img, imgFileName, cutoff, withSpectrumImgGenerated)

	for _, resultImg := range handlerResult {
		err = imageio.SaveBmpImage(&resultImg.Img, resultImg.Name)
		if err != nil {
			return "", err
		}
	}

	return "Highpass filtering successful", nil
}

func handleBandcutCommand(imgPath string, lowCut, highCut int, withSpectrumImgGenerated bool) (successMsgString string, err error) {
	if lowCut < 0 || highCut < 0 {
		return "", errors.New("low and high cut frequencies must be positive")
	}

	if lowCut >= highCut {
		return "", errors.New("low cut frequency must be lower than high cut frequency")
	}

	img, err := imageio.OpenBmpImage(imgPath)
	if err != nil {
		return "", err
	}

	if lowCut >= img.Bounds().Dx()/2 || highCut >= img.Bounds().Dx()/2 {
		return "", errors.New("low and high cut frequencies must be lower than half the image width")
	}

	imgFileName := imageio.GetPureFileName(imgPath)

	handlerResult := orthogonal_transforms.HandleBandcutFiltering(img, imgFileName, lowCut, highCut, withSpectrumImgGenerated)

	for _, resultImg := range handlerResult {
		err = imageio.SaveBmpImage(&resultImg.Img, resultImg.Name)
		if err != nil {
			return "", err
		}
	}

	return "Bandcut filtering successful", nil
}

func handlePhasemodCommand(imgPath string, k, l int) (successMsgString string, err error) {
	if k < 0 || l < 0 {
		return "", errors.New("k and l must be positive")
	}

	img, err := imageio.OpenBmpImage(imgPath)
	if err != nil {
		return "", err
	}

	imgFileName := imageio.GetPureFileName(imgPath)

	handlerResult := orthogonal_transforms.HandlePhaseModification(img, imgFileName, k, l)

	for _, resultImg := range handlerResult {
		err = imageio.SaveBmpImage(&resultImg.Img, resultImg.Name)
		if err != nil {
			return "", err
		}
	}

	return "Phase modified successful", nil
}

func handleMaskpassCommand(imgPath, maskPath string, withSpectrumImgGenerated bool) (successMsgString string, err error) {

	maskImg, err := imageio.OpenBmpImage(maskPath)
	if err != nil {
		return "", err
	}

	img, err := imageio.OpenBmpImage(imgPath)
	if err != nil {
		return "", err
	}

	imgFileName := imageio.GetPureFileName(imgPath)

	handlerResult := orthogonal_transforms.HandleMaskpassFiltering(img, imgFileName, maskImg, withSpectrumImgGenerated)

	for _, resultImg := range handlerResult {
		err = imageio.SaveBmpImage(&resultImg.Img, resultImg.Name)
		if err != nil {
			return "", err
		}
	}

	return "Maskpass filter successful applied", nil
}
