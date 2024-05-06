package main

import (
	"fmt"
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

func handleDiscordEmbed(c *gin.Context, tiktokData SimplifiedData, imageUrl string) {
	details := tiktokData.Details
	detailsString := "❤️ " + details.Likes + " | 💬 " + details.Comments + " | 🔁 " + details.Shares + " | ⭐ " + details.Favorites + " | 👀 " + details.Views
	renderTemplate(c, "discord.html", gin.H{
		"authorName": tiktokData.Author,
		"caption":    tiktokData.Caption,
		"details":    detailsString,
		"imageUrl":   imageUrl,
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

func HandleDirectFile(fileType string) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			HandleError(c, "No id provided")
			return
		}
		filename := fmt.Sprintf("collages/%s-%s", fileType, id)
		if _, err := os.Stat(filename); err != nil {
			HandleError(c, "File not found")
			return
		}
		c.File(filename)
	}
}

func preProcessTikTokRequest(c *gin.Context) (SimplifiedData, bool) {
	tiktokURL := c.Query("v")
	videoId, err := GetLongVideoId(tiktokURL)
	if err != nil {
		if err.Error() == "invalid URL" {
			HandleError(c, "link: "+tiktokURL+" is invalid")
		} else {
			HandleError(c, "Couldn't get tiktok")
		}
		return SimplifiedData{}, true
	}

	tiktokData, err := FetchTiktokData(videoId)
	if err != nil {
		HandleError(c, "Couldn't get tiktok")
		return SimplifiedData{}, true
	}
	if tiktokData.IsVideo {
		handleVideoDiscordEmbed(
			c,
			tiktokData,
			tiktokData.Video.Url,
			tiktokData.Video.Width,
			tiktokData.Video.Height,
		)
		return SimplifiedData{}, true
	}

	var filename string
	path := c.Request.URL.Path
	switch path {
	case PathCollage:
		filename = "collage-" + videoId + ".png"
	case PathCollageSound:
		filename = "video-" + videoId + ".mp4"
	case PathSlide:
		filename = "slide-" + videoId + ".mp4"
	}

	if _, err := os.Stat("collages/" + filename); err != nil {
		return tiktokData, false
	}

	if path == PathCollage {
		handleDiscordEmbed(c, tiktokData, Domain+filename)
		return SimplifiedData{}, true
	}

	width, height, err := GetVideoDimensions("collages/" + filename)
	if err != nil {
		HandleError(c, "Couldn't get video dimensions")
		return SimplifiedData{}, true
	}
	handleVideoDiscordEmbed(c, tiktokData, Domain+filename, width, height)
	return SimplifiedData{}, true
}

func processRequest(c *gin.Context, collageImages bool, downloadSound bool) (SimplifiedData, bool, error) {
	tiktokData, skip := preProcessTikTokRequest(c)
	if skip {
		return tiktokData, true, nil
	}
	tiktokData.DownloadImages()
	if downloadSound {
		err := tiktokData.DownloadSound()
		if err != nil {
			return SimplifiedData{}, false, err
		}
	}

	if !collageImages {
		return tiktokData, false, nil
	}

	err := tiktokData.MakeCollage()
	if err != nil {
		println(err.Error())
		HandleError(c, "Couldn't generate collage")
		return SimplifiedData{}, false, err
	}

	return tiktokData, false, nil
}

func HandleRequest(c *gin.Context) {
	tiktokData, skip, err := processRequest(c, true, false)
	if skip {
		return
	}
	if err != nil {
		HandleError(c, err.Error())
		return
	}
	handleDiscordEmbed(c, tiktokData, Domain+"collage-"+tiktokData.VideoID+".png")
	UpdateLocalStats()
}

func HandleSoundCollageRequest(c *gin.Context) {
	tiktokData, skip, err := processRequest(c, true, true)
	if skip {
		return
	}
	if err != nil {
		HandleError(c, err.Error())
		return
	}

	width, height, err := tiktokData.MakeCollageWithAudio()
	if err != nil {
		println(err.Error())
		HandleError(c, "Couldn't generate video")
		return
	}

	handleVideoDiscordEmbed(c, tiktokData, Domain+"video-"+tiktokData.VideoID+".mp4", width, height)
	UpdateLocalStats()
}

func HandleFancySlideshowRequest(c *gin.Context) {
	tiktokData, skip, err := processRequest(c, false, true)
	if skip {
		return
	}
	if err != nil {
		HandleError(c, err.Error())
		return
	}

	videoWidth, videoHeight, err := tiktokData.MakeVideoSlideshow()
	if err != nil {
		println(err.Error())
		HandleError(c, "Couldn't generate video")
		return
	}

	handleVideoDiscordEmbed(
		c,
		tiktokData,
		Domain+"slide-"+tiktokData.VideoID+".mp4",
		videoWidth,
		videoHeight,
	)

	os.RemoveAll(tiktokData.VideoID)
	UpdateLocalStats()
}
