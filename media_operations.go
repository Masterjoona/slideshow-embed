package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func FetchAudio(c *gin.Context, tiktokURL string, tiktokData Data, randomErrorImage string) {
	audioLink, err := ExtractAudioLink(tiktokData.Body)
	if err != nil {
		handleError(c, "Couldn't get audio link", randomErrorImage)
		return
	}

	err = DownloadAudio(audioLink, "audio.mp3", tiktokData.VideoID)
	if err != nil {
		handleError(c, "Couldn't download audio", randomErrorImage)
		return
	}
}

func FetchImages(c *gin.Context, tiktokURL string, tiktokData Data, randomErrorImage string) {
	links, err := ExtractImageLinks(tiktokData.Body)
	if err != nil {
		handleError(c, "Couldn't get image links", randomErrorImage)
		return
	}

	err = DownloadImages(links, tiktokData.VideoID)
	if err != nil {
		handleError(c, "Couldn't download images", randomErrorImage)
		return
	}
}

func FetchTiktokData(c *gin.Context, tiktokURL string, errorImage string) (tiktokData Data) {
	if !validateURL(tiktokURL) {
		handleError(c, "Invalid url", errorImage)
		return
	}

	videoID, err := ExtractVideoID(tiktokURL)
	if err != nil {
		handleError(c, "Invalid url", errorImage)
		return
	}

	authorName, caption, responseBody, err := GetVideoAuthorAndCaption(tiktokURL, videoID)
	if err != nil {
		handleError(c, "Couldn't get video author and caption. Is the slideshow available?", errorImage)
		return
	}

	details := GetVideoDetails(responseBody)
	return Data{
		AuthorName: authorName,
		Caption:    caption,
		VideoID:    videoID,
		Body:       responseBody,
		Details:    details,
	}
}

func GenerateCollage(c *gin.Context, videoId string, collageFilename string, errorImage string) {
	err := MakeCollage(videoId, collageFilename)
	if err != nil {
		handleError(c, "Couldn't make collage", errorImage)
		return
	}
}

func GenerateVideo(c *gin.Context, videoId string, collageFilename string, videoFilename string, errorImage string, sliding bool) (videoWidth, videoHeight string) {
	if sliding {
		err := MakeVideoSlideshow(videoId, videoFilename)
		if err != nil {
			handleError(c, "Couldn't make slideshow video", errorImage)
			return "", ""
		}
	} else {
		err := MakeVideo("collages/"+collageFilename, videoId, videoFilename)
		if err != nil {
			fmt.Println(err)
			handleError(c, "Couldn't make video", errorImage)
			return "", ""
		}
	}

	videoWidth, videoHeight, err := GetVideoDimensions("collages/" + videoFilename)
	if err != nil {
		handleError(c, "Couldn't get video dimensions", errorImage)
		return
	}
	return videoWidth, videoHeight
}
