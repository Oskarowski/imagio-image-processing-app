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
