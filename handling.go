package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

var errorImages = []string{
	"https://media.discordapp.net/attachments/961445186280509451/980132677338423316/fuckmedaddyharderohyeailovecokcimsocissyfemboy.gif",
	"https://media.discordapp.net/attachments/901959319719936041/996927812927750264/chrome_2WOKI6Jm3v.gif",
	"https://cdn.discordapp.com/attachments/749030987530502197/980338691706880051/79587A35-FD36-41D3-8232-7A29B46D2543.gif",
}
var errorImagesIndex = 0

func renderTemplate(c *gin.Context, filename string, data gin.H) {
	tmpl, err := template.ParseFiles("templates/" + filename)
	if err != nil {
		handleError(c, err.Error(), errorImages[errorImagesIndexInt()])
		return
	}

	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		handleError(c, err.Error(), errorImages[errorImagesIndexInt()])
	}
}

func validateURL(url string) bool {
	if url == "" {
		return false
	}
	if !strings.Contains(url, "vm.tiktxk.com") && !strings.Contains(url, "vm.tiktok.com") {
		return false
	}
	return true
}

func errorImagesIndexInt() int {
	if errorImagesIndex == 2 {
		errorImagesIndex = 0
	} else {
		errorImagesIndex++
	}
	return errorImagesIndex
}

func handleError(c *gin.Context, errorMsg string, errorImageUrl string) {
	renderTemplate(c, "error.html", gin.H{
		"error":           errorMsg,
		"error_image_url": errorImageUrl,
	})
}

func handleDiscordEmbed(c *gin.Context, authorName string, caption string, details Counts, filename string) {
	detailsString := "❤️ " + details.Likes + " | 💬 " + details.Comments + " | 🔁 " + details.Shares + " | ⭐ " + details.Favorited + " | 👀 " + details.Views
	renderTemplate(c, "discord.html", gin.H{
		"authorName": authorName,
		"caption":    caption,
		"details":    detailsString,
		"imageUrl":   Domain + "/" + filename,
	})
}

func handleVideoDiscordEmbed(c *gin.Context, authorName string, caption string, details Counts, filename string, width string, height string) {
	detailsString := "❤️" + details.Likes + " | 💬 " + details.Comments + " | 🔁 " + details.Shares + " | ⭐ " + details.Favorited + " | 👀 " + details.Views
	authorName = strings.Split(authorName, "(@")[0]
	renderTemplate(c, "video.html", gin.H{
		"authorName": authorName,
		"details":    detailsString,
		"videoUrl":   Domain + "/" + filename,
		"caption":    caption,
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
	collageFiles, err := os.ReadDir("collages")
	if err != nil {
		handleError(c, err.Error(), errorImages[errorImagesIndexInt()])
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

	bytes, err := GetDirectorySize("collages")
	size := FormatSize(bytes)
	if err != nil {
		handleError(c, err.Error(), errorImages[errorImagesIndexInt()])
		return
	}
	renderTemplate(c, "index.html", gin.H{
		"FileLinks": filePaths,
		"count":     count,
		"size":      size,
	})
}

func HandleSoundCollageRequest(c *gin.Context) {
	tiktokURL := c.Query("v")

	randomErrorImage := errorImages[errorImagesIndexInt()]

	if !validateURL(tiktokURL) {
		handleError(c, "Invalid url", randomErrorImage)
		return
	}

	videoID, err := ExtractVideoID(tiktokURL)
	if err != nil {
		handleError(c, "Invalid url", randomErrorImage)
		return
	}

	filename := "video-" + videoID + ".mp4"
	authorName, caption, responseBody, err := GetVideoAuthorAndCaption(tiktokURL, videoID)
	if err != nil {
		handleError(c, "Couldn't get video author and caption. Is the video available?", randomErrorImage)
		return
	}

	details := GetVideoDetails(responseBody)

	if _, err := os.Stat("collages/" + filename); err == nil {
		videoWidth, videoHeight, err := GetVideoDimensions("collages/" + filename)
		if err != nil {
			handleError(c, "Couldn't get video dimensions", randomErrorImage)
			return
		}
		handleVideoDiscordEmbed(c, authorName, caption, details, filename, videoWidth, videoHeight)
		return
	}

	links, err := ExtractImageLinks(responseBody)
	if err != nil {
		handleError(c, "Couldn't get image links", randomErrorImage)
		return
	}

	err = DownloadImages(links, videoID)
	if err != nil {
		handleError(c, "Couldn't download images", randomErrorImage)
		return
	}

	audioLink, err := ExtractAudioLink(responseBody)
	if err != nil {
		handleError(c, "Couldn't get audio link", randomErrorImage)
		return
	}

	err = DownloadAudio(audioLink, "audio.mp3", videoID)
	if err != nil {
		handleError(c, "Couldn't download audio", randomErrorImage)
		return
	}

	err = MakeCollage(videoID, "collage-"+videoID+".png")
	if err != nil {
		handleError(c, "Couldn't make collage", randomErrorImage)
		return
	}

	err = MakeVideo("collages/collage-"+videoID+".png", videoID, filename)
	if err != nil {
		fmt.Println(err)
		handleError(c, "Couldn't make video", randomErrorImage)
		return
	}

	videoWidth, videoHeight, err := GetVideoDimensions("collages/" + filename)
	if err != nil {
		handleError(c, "Couldn't get video dimensions", randomErrorImage)
		return
	}
	handleVideoDiscordEmbed(c, authorName, caption, details, filename, videoWidth, videoHeight)
	os.RemoveAll(videoID)
}

func HandleRequest(c *gin.Context) {
	startTime := time.Now()
	tiktokURL := c.Query("v")
	debug := c.Query("d")

	randomErrorImage := errorImages[errorImagesIndexInt()]

	if !validateURL(tiktokURL) {
		handleError(c, "Invalid url", randomErrorImage)
		return
	}

	videoID, err := ExtractVideoID(tiktokURL)
	if err != nil {
		handleError(c, "Invalid url", randomErrorImage)
		return
	}

	filename := "collage-" + videoID + ".png"
	authorName, caption, responseBody, err := GetVideoAuthorAndCaption(tiktokURL, videoID)
	if err != nil {
		handleError(c, "Couldn't get video author and caption. Is the video available?", randomErrorImage)
		return
	}

	details := GetVideoDetails(responseBody)

	if _, err := os.Stat("collages/" + filename); err == nil {
		if debug == "true" {
			elapsed := time.Since(startTime)
			caption = caption + " | Took " + elapsed.String()
		}
		handleDiscordEmbed(c, authorName, caption, details, filename)
		return
	}

	links, err := ExtractImageLinks(responseBody)
	if err != nil {
		handleError(c, "Couldn't get image links", randomErrorImage)
		return
	}

	err = DownloadImages(links, videoID)
	if err != nil {
		handleError(c, "Couldn't download images", randomErrorImage)
		return
	}

	err = MakeCollage(videoID, filename)
	if err != nil {
		handleError(c, "Couldn't make collage", randomErrorImage)
		return
	}

	if debug == "true" || debug == "1" {
		elapsed := time.Since(startTime)
		filesizeBytes, err := GetFileSize("collages/" + filename)
		if err != nil {
			handleError(c, "Couldn't get filesize", randomErrorImage)
			return
		}
		filesize := FormatSize(filesizeBytes)
		caption = caption + " | Took " + elapsed.String() + " | " + filesize
	}

	handleDiscordEmbed(c, authorName, caption, details, filename)
	os.RemoveAll(videoID)
}

func HandleFancySlideshowRequest(c *gin.Context) {
	tiktokURL := c.Query("v")

	randomErrorImage := errorImages[errorImagesIndexInt()]

	if !validateURL(tiktokURL) {
		handleError(c, "Invalid url", randomErrorImage)
		return
	}

	videoID, err := ExtractVideoID(tiktokURL)
	if err != nil {
		handleError(c, "Invalid url", randomErrorImage)
		return
	}

	filename := "slide-" + videoID + ".mp4"
	authorName, caption, responseBody, err := GetVideoAuthorAndCaption(tiktokURL, videoID)
	if err != nil {
		handleError(c, "Couldn't get video author and caption. Is the video available?", randomErrorImage)
		return
	}

	details := GetVideoDetails(responseBody)

	if _, err := os.Stat("collages/" + filename); err == nil {
		videoWidth, videoHeight, err := GetVideoDimensions("collages/" + filename)
		if err != nil {
			handleError(c, "Couldn't get video dimensions", randomErrorImage)
			return
		}
		handleVideoDiscordEmbed(c, authorName, caption, details, filename, videoWidth, videoHeight)
		return
	}

	links, err := ExtractImageLinks(responseBody)
	if err != nil {
		handleError(c, "Couldn't get image links", randomErrorImage)
		return
	}

	err = DownloadImages(links, videoID)
	if err != nil {
		handleError(c, "Couldn't download images", randomErrorImage)
		return
	}

	audioLink, err := ExtractAudioLink(responseBody)
	if err != nil {
		handleError(c, "Couldn't get audio link", randomErrorImage)
		return
	}

	err = DownloadAudio(audioLink, "audio.mp3", videoID)
	if err != nil {
		handleError(c, "Couldn't download audio", randomErrorImage)
		return
	}

	err = MakeCollage(videoID, "collage-"+videoID+".png")
	if err != nil {
		handleError(c, "Couldn't make collage", randomErrorImage)
		return
	}

	err = MakeVideoSlideshow(videoID, filename)
	if err != nil {
		handleError(c, "Couldn't make slideshow video", randomErrorImage)
		return
	}

	videoWidth, videoHeight, err := GetVideoDimensions("collages/slide-" + videoID + ".mp4")
	if err != nil {
		handleError(c, "Couldn't get video dimensions", randomErrorImage)
		return
	}
	handleVideoDiscordEmbed(c, authorName, caption, details, "slide-"+videoID+".mp4", videoWidth, videoHeight)
	os.RemoveAll(videoID)
}

func HandleDirectCollage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleError(c, "No id provided", errorImages[errorImagesIndexInt()])
		return
	}

	filename := "collage-" + id
	if _, err := os.Stat("collages/" + filename); err != nil {
		handleError(c, "Collage not found", errorImages[errorImagesIndexInt()])
		return
	}

	c.File("collages/" + filename)

}

func HandleDirectVideo(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleError(c, "No id provided", errorImages[errorImagesIndexInt()])
		return
	}

	filename := "video-" + id
	if _, err := os.Stat("collages/" + filename); err != nil {
		handleError(c, "Collage not found", errorImages[errorImagesIndexInt()])
		return
	}

	c.File("collages/" + filename)
}

func HandleDirectFancyCollage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleError(c, "No id provided", errorImages[errorImagesIndexInt()])
		return
	}

	filename := "slide-" + id
	if _, err := os.Stat("collages/" + filename); err != nil {
		handleError(c, "Collage not found", errorImages[errorImagesIndexInt()])
		return
	}

	c.File("collages/" + filename)
}
