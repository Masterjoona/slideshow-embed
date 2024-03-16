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
	isVideo := !strings.Contains(postAweme.Video.PlayAddr.URLList[0], "music")
	imageLinks := []string{}
	if !isVideo {
		imageLinks = GetImageLinks(postAweme)
	}
	return SimplifiedData{
		Author:     GetAuthor(postAweme),
		Caption:    postAweme.Desc + "\n\n" + postAweme.Music.Title + " - " + postAweme.Music.Author,
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
