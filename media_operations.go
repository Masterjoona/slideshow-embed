package main

import (
	"fmt"
	"strings"
)

func FetchImages(links []string, videoId string) error {
	err := DownloadImages(links, videoId)
	if err != nil {
		return err
	}
	return nil
}

func FetchAudio(url string, videoId string) error {
	err := DownloadAudio(url, videoId)
	if err != nil {
		return err
	}
	return nil
}

func FetchTiktokData(videoId string) (SimplifiedData, error) {
	resp, err := PostDetails(videoId)
	if err != nil {
		return SimplifiedData{}, err
	}
	postAweme := resp.AwemeList[0]
	isVideo := !strings.HasSuffix(postAweme.Video.PlayAddr.URLList[0], "mp3")
	imageLinks := []string{}
	if !isVideo {
		imageLinks = GetImageLinks(postAweme)
	}
	return SimplifiedData{
		AuthorName: GetAuthor(postAweme),
		Caption:    postAweme.Desc,
		VideoID:    videoId,
		Details:    GetVideoDetails(postAweme),
		ImageLinks: imageLinks,
		SoundUrl:   postAweme.Music.PlayURL.URI,
		IsVideo:    isVideo,
		Video:      postAweme.Video,
	}, nil
}

func GenerateCollage(videoId string, collageFilename string) error {
	err := MakeCollage(videoId, collageFilename)
	if err != nil {
		return err
	}
	return nil
}

func GenerateVideo(
	videoId string,
	collageFilename string,
	videoFilename string,
	sliding bool,
) (string, string, error) {
	if sliding {
		err := MakeVideoSlideshow(videoId, videoFilename)
		if err != nil {
			return "", "", err
		}
	} else {
		err := MakeVideo("collages/"+collageFilename, videoId, videoFilename)
		if err != nil {
			fmt.Println(err)
			return "", "", err
		}
	}

	videoWidth, videoHeight, err := GetVideoDimensions("collages/" + videoFilename)
	if err != nil {
		return "", "", err
	}
	return videoWidth, videoHeight, nil
}
