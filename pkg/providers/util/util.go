package util

import (
	"meow/pkg/files"
	"meow/pkg/types"
)

func GetDimensionsOrNil(videoUrl string, check bool) (types.SimplifiedVideo, error) {
	if !check {
		return types.SimplifiedVideo{}, nil
	}

	width, height, err := files.GetVideoDimensions(videoUrl)
	if err != nil {
		return types.SimplifiedVideo{}, err
	}

	return types.SimplifiedVideo{
		Width:  width,
		Height: height,
		Url:    videoUrl,
	}, nil
}
