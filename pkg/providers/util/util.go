package util

import (
	"meow/pkg/files"
	"meow/pkg/util"
)

type VideoDimensions struct {
	Width  string
	Height string
	Err    error
}

func GetDimensionsOrNil(videoUrl string, check bool) VideoDimensions {
	return util.Ternary(check,
		func() VideoDimensions {
			width, height, err := files.GetVideoDimensionsFromUrl(videoUrl)
			return VideoDimensions{Width: width, Height: height, Err: err}
		}(),
		VideoDimensions{})
}
