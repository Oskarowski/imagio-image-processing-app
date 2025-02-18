package cmd

import "image"

type ResultImage interface {
	GetImage() *image.RGBA
	GetName() string
}

type BasicImgResult struct {
	Img  *image.RGBA
	Name string
}

func (r BasicImgResult) GetImage() *image.RGBA { return r.Img }
func (r BasicImgResult) GetName() string       { return r.Name }

func ToResultImage[T ResultImage](imgs []T) []ResultImage {
	result := make([]ResultImage, len(imgs))
	for i, img := range imgs {
		result[i] = img
	}
	return result
}
