package main

import (
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
)

type Data struct {
	AuthorName string
	Caption    string
	VideoID    string
	Body       string
	Details    Counts
}

func renderTemplate(c *gin.Context, filename string, data gin.H) {
	tmpl, err := template.ParseFiles("templates/" + filename)
	if err != nil {
		handleError(c, err.Error(), errorImage())
		return
	}

	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		handleError(c, err.Error(), errorImage())
	}
}

func handleError(c *gin.Context, errorMsg string, errorImageUrl string) {
	renderTemplate(c, "error.html", gin.H{
		"error":           errorMsg,
		"error_image_url": errorImageUrl,
	})
}

func handleDiscordEmbed(c *gin.Context, tiktokData Data, filename string) {
	details := tiktokData.Details
	detailsString := "❤️ " + details.Likes + " | 💬 " + details.Comments + " | 🔁 " + details.Shares + " | ⭐ " + details.Favorited + " | 👀 " + details.Views
	renderTemplate(c, "discord.html", gin.H{
		"authorName": tiktokData.AuthorName,
		"caption":    tiktokData.Caption,
		"details":    detailsString,
		"imageUrl":   Domain + "/" + filename,
	})
}

func handleVideoDiscordEmbed(c *gin.Context, tiktokData Data, filename string, width string, height string) {
	details := tiktokData.Details
	detailsString := "❤️" + details.Likes + " | 💬 " + details.Comments + " | 🔁 " + details.Shares + " | ⭐ " + details.Favorited + " | 👀 " + details.Views
	authorName := strings.Split(tiktokData.AuthorName, "(@")[0]
	renderTemplate(c, "video.html", gin.H{
		"authorName": authorName,
		"details":    detailsString,
		"videoUrl":   Domain + "/" + filename,
		"caption":    tiktokData.Caption,
		"width":      width,
		"height":     height,
	})
}

func handleExistingFile(c *gin.Context, filename string, video bool, tiktokData Data) bool {
	if _, err := os.Stat("collages/" + filename); err == nil {
		println("File exists")
		if video {
			println("File is video")
			videoWidth, videoHeight, err := GetVideoDimensions("collages/" + filename)
			if err != nil {
				handleError(c, "Couldn't get video dimensions", errorImage())
				return true
			}
			handleVideoDiscordEmbed(c, tiktokData, filename, videoWidth, videoHeight)
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
	collageFiles, err := os.ReadDir("collages")
	if err != nil {
		handleError(c, err.Error(), errorImage())
		return
	}
	filePaths := make([]string, len(collageFiles))
	count := 0
	sort.Slice(collageFiles, func(i, j int) bool {
		fileI, err1 := collageFiles[i].Info()
		fileJ, err2 := collageFiles[j].Info()
		if err1 != nil || err2 != nil {
			return collageFiles[i].Name() < collageFiles[j].Name()
		}
		return fileI.ModTime().After(fileJ.ModTime())
	})

	for index, file := range collageFiles {
		filePaths[index] = Domain + "/" + file.Name()
		count++
	}
	countString := strconv.Itoa(count)
	if LimitPublicAmount > 0 && len(filePaths) > LimitPublicAmount {
		filePaths = filePaths[:LimitPublicAmount]
		countString += " (Only showing " + strconv.Itoa(len(filePaths)) + ")"
	}

	bytes, err := GetDirectorySize("collages")
	size := FormatSize(bytes)
	if err != nil {
		handleError(c, err.Error(), errorImage())
		return
	}
	renderTemplate(c, "index.html", gin.H{
		"FileLinks": filePaths,
		"count":     countString,
		"size":      size,
	})
}

func HandleSoundCollageRequest(c *gin.Context) {
	tiktokURL := c.Query("v")

	randomErrorImage := errorImage()
	tiktokData := FetchTiktokData(c, tiktokURL, randomErrorImage)

	filename := "video-" + tiktokData.VideoID + ".mp4"
	if handleExistingFile(c, filename, true, tiktokData) {
		return
	}

	collageFilename := "collage-" + tiktokData.VideoID + ".png"
	FetchImages(c, tiktokURL, tiktokData, randomErrorImage)
	FetchAudio(c, tiktokURL, tiktokData, randomErrorImage)
	GenerateCollage(c, tiktokData.VideoID, collageFilename, randomErrorImage)

	videoWidth, videoHeight := GenerateVideo(c, tiktokData.VideoID, collageFilename, filename, randomErrorImage, false)
	handleVideoDiscordEmbed(c, tiktokData, filename, videoWidth, videoHeight)

	os.RemoveAll(tiktokData.VideoID)
}

func HandleRequest(c *gin.Context) {
	tiktokURL := c.Query("v")

	randomErrorImage := errorImage()
	tiktokData := FetchTiktokData(c, tiktokURL, randomErrorImage)

	filename := "collage-" + tiktokData.VideoID + ".png"
	if handleExistingFile(c, filename, false, tiktokData) {
		return
	}

	FetchImages(c, tiktokURL, tiktokData, randomErrorImage)
	println("Images fetched")
	GenerateCollage(c, tiktokData.VideoID, filename, randomErrorImage)
	println("Collage generated")

	handleDiscordEmbed(c, tiktokData, filename)
	//os.RemoveAll(tiktokData.VideoID)
}

func HandleFancySlideshowRequest(c *gin.Context) {
	tiktokURL := c.Query("v")

	randomErrorImage := errorImage()
	tiktokData := FetchTiktokData(c, tiktokURL, randomErrorImage)

	filename := "slide-" + tiktokData.VideoID + ".mp4"
	if handleExistingFile(c, filename, true, tiktokData) {
		return
	}

	FetchImages(c, tiktokURL, tiktokData, randomErrorImage)
	FetchAudio(c, tiktokURL, tiktokData, randomErrorImage)

	videoWidth, videoHeight := GenerateVideo(c, tiktokData.VideoID, "", filename, randomErrorImage, true)
	handleVideoDiscordEmbed(c, tiktokData, filename, videoWidth, videoHeight)

	//os.RemoveAll(tiktokData.VideoID)
}

func HandleDirectFile(c *gin.Context) {
	id := c.Param("id")
	mediaType := strings.Split(c.Request.URL.Path, "-")[0][1:]
	if id == "" || mediaType == "" {
		handleError(c, "No id provided", errorImage())
		return
	}
	filename := mediaType + "-" + id
	if _, err := os.Stat("collages/" + filename); err != nil {
		handleError(c, "File not found", errorImage())
		return
	}

	c.File("collages/" + filename)
}
