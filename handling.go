package main

import (
	"os"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
)

func renderTemplate(c *gin.Context, filename string, data gin.H) {
	tmpl, err := template.ParseFiles("templates/" + filename)
	if err != nil {
		HandleError(c, err.Error())
		return
	}

	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		HandleError(c, err.Error())
	}
}

func HandleError(c *gin.Context, errorMsg string) {
	renderTemplate(c, "error.html", gin.H{
		"error":           errorMsg,
		"error_image_url": ErrorImage(),
	})
}

func handleDiscordEmbed(c *gin.Context, tiktokData SimplifiedData, filename string) {
	details := tiktokData.Details
	detailsString := "❤️ " + details.Likes + " | 💬 " + details.Comments + " | 🔁 " + details.Shares + " | ⭐ " + details.Favorites + " | 👀 " + details.Views
	renderTemplate(c, "discord.html", gin.H{
		"authorName": tiktokData.Author,
		"caption":    tiktokData.Caption,
		"details":    detailsString,
		"imageUrl":   Domain + filename,
	})
}

func handleVideoDiscordEmbed(
	c *gin.Context,
	tiktokData SimplifiedData,
	url string,
	width string,
	height string,
) {
	details := tiktokData.Details
	detailsString := "❤️" + details.Likes + " | 💬 " + details.Comments + " | 🔁 " + details.Shares + " | ⭐ " + details.Favorites + " | 👀 " + details.Views
	renderTemplate(c, "video.html", gin.H{
		"authorName": strings.Split(tiktokData.Author, "(@")[0],
		"details":    detailsString,
		"videoUrl":   url,
		"caption":    tiktokData.Caption,
		"width":      width,
		"height":     height,
	})
}

func handleExistingFile(
	c *gin.Context,
	filename string,
	video bool,
	tiktokData SimplifiedData,
) bool {
	if _, err := os.Stat("collages/" + filename); err == nil {
		if video {
			videoWidth, videoHeight, err := GetVideoDimensions("collages/" + filename)
			if err != nil {
				HandleError(c, "Couldn't get video dimensions")
				return true
			}
			handleVideoDiscordEmbed(c, tiktokData, Domain+filename, videoWidth, videoHeight)
			return true
		}
		handleDiscordEmbed(c, tiktokData, filename)
		return true
	}
	return false
}

func HandleIndex(c *gin.Context) {
	if !Public {
		renderTemplate(c, "index.html", gin.H{
			"FileLinks": nil,
			"count":     "0",
			"size":      "0",
		})
		return
	}

	renderTemplate(c, "index.html", gin.H{
		"FileLinks": LocalStats.FilePaths,
		"count":     LocalStats.FileCount,
		"size":      LocalStats.TotalSize,
	})
}

func HandleSoundCollageRequest(c *gin.Context) {
	tiktokURL := c.Query("v")

	videoId, err := GetLongVideoId(tiktokURL)
	if err != nil {
		HandleError(c, "Couldn't fetch slideshow or your URL is invalid")
		return
	}

	collageFilename := "collages/collage-" + videoId + ".png"
	videoFilename := "video-" + videoId + ".mp4"
	tiktokData, err := FetchTiktokData(videoId)

	if err != nil {
		HandleError(c, "Couldn't fetch TikTok data")
		return
	}

	if handleExistingFile(c, videoFilename, true, tiktokData) {
		return
	}

	collageFileExists, _ := os.Stat(collageFilename)
	if collageFileExists != nil {
		CreateDirectory(videoId)
		audioBuffer, err := DownloadAudio(tiktokData.SoundUrl)
		if err != nil {
			HandleError(c, "Couldn't fetch audio")
			return
		}
		videoWidth, videoHeight, err := MakeCollageWithAudio(
			collageFilename,
			audioBuffer,
			videoId,
		)
		if err != nil {
			HandleError(c, "Couldn't generate video")
			return
		}
		handleVideoDiscordEmbed(c, tiktokData, Domain+videoFilename, videoWidth, videoHeight)
		return
	}

	imageBuffers, err := DownloadImages(tiktokData.ImageLinks)
	if err != nil {
		HandleError(c, "Couldn't fetch images")
		return
	}

	audioBuffer, err := DownloadAudio(tiktokData.SoundUrl)
	if err != nil {
		HandleError(c, "Couldn't fetch audio")
		return
	}

	err = MakeCollage(imageBuffers, videoId)
	if err != nil {
		HandleError(c, "Couldn't generate collage")
		return
	}

	videoWidth, videoHeight, err := MakeCollageWithAudio(
		collageFilename,
		audioBuffer,
		videoId,
	)

	if err != nil {
		HandleError(c, "Couldn't generate video")
		return
	}
	handleVideoDiscordEmbed(c, tiktokData, Domain+videoFilename, videoWidth, videoHeight)

	os.RemoveAll(videoId)
	UpdateLocalStats()
}

func HandleRequest(c *gin.Context) {
	tiktokURL := c.Query("v")

	videoId, err := GetLongVideoId(tiktokURL)
	if err != nil {
		HandleError(c, "Couldn't fetch slideshow")
		return
	}
	filename := "collage-" + videoId + ".png"
	tiktokData, err := FetchTiktokData(videoId)
	if err != nil {
		HandleError(c, "Couldn't fetch TikTok data")
		return
	}

	if tiktokData.IsVideo {
		handleVideoDiscordEmbed(
			c,
			tiktokData,
			tiktokData.Video.Url,
			tiktokData.Video.Width,
			tiktokData.Video.Height,
		)
		return
	}

	if handleExistingFile(c, filename, false, tiktokData) {
		return
	}
	images, err := DownloadImages(tiktokData.ImageLinks)
	if err != nil {
		HandleError(c, "Couldn't fetch images")
		return
	}
	err = MakeCollage(images, videoId)
	if err != nil {
		HandleError(c, "Couldn't generate collage")
		return
	}

	handleDiscordEmbed(c, tiktokData, filename)
	os.RemoveAll(videoId)
	UpdateLocalStats()
}

func HandleFancySlideshowRequest(c *gin.Context) {
	tiktokURL := c.Query("v")

	videoId, err := GetLongVideoId(tiktokURL)
	if err != nil {
		HandleError(c, "Couldn't fetch slideshow")
		return
	}
	filename := "slide-" + videoId + ".mp4"
	tiktokData, err := FetchTiktokData(videoId)
	if err != nil {
		HandleError(c, "Couldn't fetch TikTok data")
		return
	}

	if handleExistingFile(c, filename, true, tiktokData) {
		return
	}

	imageBuffers, err := DownloadImages(tiktokData.ImageLinks)
	if err != nil {
		HandleError(c, "Couldn't fetch images")
		return
	}
	audioBuffer, err := DownloadAudio(tiktokData.SoundUrl)
	if err != nil {
		HandleError(c, "Couldn't fetch audio")
		return
	}

	videoWidth, videoHeight, err := MakeVideoSlideshow(imageBuffers, audioBuffer, videoId)
	if err != nil {
		HandleError(c, "Couldn't generate video")
		return
	}

	handleVideoDiscordEmbed(c, tiktokData, Domain+filename, videoWidth, videoHeight)

	os.RemoveAll(tiktokData.VideoID)
	UpdateLocalStats()
}

func HandleDirectFile(c *gin.Context) {
	id := c.Param("id")
	mediaType := strings.Split(c.Request.URL.Path, "-")[0][1:]
	if id == "" || mediaType == "" {
		HandleError(c, "No id provided")
		return
	}
	filename := mediaType + "-" + id
	if _, err := os.Stat("collages/" + filename); err != nil {
		HandleError(c, "File not found")
		return
	}
	c.File("collages/" + filename)
}
