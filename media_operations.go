package main

import (
	"fmt"
)

func FetchImages(content string, videoId string) error {
	links, err := GetImageLinks(content)
	if err != nil {
		return err
	}

	err = DownloadImages(links, videoId)
	if err != nil {
		return err
	}
	return nil
}

func FetchAudio(content string, videoId string) error {
	audioUrl, err := GetAudioLink(content)
	if err != nil {
		return err
	}

	err = DownloadAudio(audioUrl, videoId)
	if err != nil {
		return err
	}
	return nil
}

func FetchTiktokData(videoID string) (Data, error) {
	var content string
	retryCount := 0
	if ProxiTokInstance != "" && ProxiTokInstance != "/" {
		retryCount = 3
	}
	var err error
	for retryCount < 4 {
		content, err = FetchProxiTokVideo(videoID)
		if content == "private" {
			return Data{Private: true}, err
		}
		if err != nil {
			print("Instance failed, retrying... ")
			retryCount++
			continue
		}
		break
	}
	styledAuthor, caption, err := GetAuthorAndCaption(content)
	if err != nil {
		return Data{}, err
	}

	details := GetVideoDetails(content)

	return Data{
		AuthorName: styledAuthor,
		Caption:    caption,
		VideoID:    videoID,
		Details:    details,
		Body:       content,
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
