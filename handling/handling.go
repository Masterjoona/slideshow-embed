package handling

import (
	"meow/collaging"
	"meow/config"
	"meow/extracting"
	"meow/files"
	"meow/httputil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

func renderTemplate(c *gin.Context, filename string, data gin.H) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
func handleError(c *gin.Context, errorMsg, errorImageUrl string) {
	renderTemplate(c, "error.html", gin.H{
		"error":           errorMsg,
		"error_image_url": errorImageUrl,
	})
}

func handleDiscordEmbed(c *gin.Context, authorName, caption, filename string) {
	renderTemplate(c, "discord.html", gin.H{
		"authorName": authorName,
		"caption":    caption,
		"imageUrl":   config.Domain + filename,
	})
}

var errorImages = []string{
	"https://media.discordapp.net/attachments/961445186280509451/980132677338423316/fuckmedaddyharderohyeailovecokcimsocissyfemboy.gif",
	"https://media.discordapp.net/attachments/901959319719936041/996927812927750264/chrome_2WOKI6Jm3v.gif",
	"https://cdn.discordapp.com/attachments/749030987530502197/980338691706880051/79587A35-FD36-41D3-8232-7A29B46D2543.gif",
}
var errorImagesIndex = 0

func isInvalidIntStr(str string, min, max int) bool {
	intValue, err := strconv.Atoi(str)
	return err != nil || intValue < min || intValue > max
}

func HandleIndex(c *gin.Context) {
	if !config.Public {
		renderTemplate(c, "index.html", gin.H{
			"FileLinks": nil,
			"count":     "0",
			"size":      "0",
		})
		return
	}
	collageFiles, err := os.ReadDir("collages")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	filePaths := make([]string, len(collageFiles))
	count := 0
	for index, file := range collageFiles {
		filePaths[index] = config.Domain + file.Name()
		count++
	}
	bytes, err := files.GetDirectorySize("collages")
	size := files.FormatSize(bytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	renderTemplate(c, "index.html", gin.H{
		"FileLinks": filePaths,
		"count":     count,
		"size":      size,
	})
}

func HandleTikTokRequest(c *gin.Context) {
	startTime := time.Now()
	tiktokURL := c.Query("v")
	width := c.Query("w")
	initHeight := c.Query("h")
	debug := c.Query("d")

	if width == "" || isInvalidIntStr(width, 1, 4096) {
		width = "1024"
	}

	if initHeight == "" || isInvalidIntStr(initHeight, 1, 1024) {
		initHeight = "320"
	}

	randomErrorImage := errorImages[errorImagesIndex]
	if errorImagesIndex == 2 {
		errorImagesIndex = 0
	} else {
		errorImagesIndex++
	}

	if tiktokURL == "" {
		handleError(c, "No url provided", randomErrorImage)
		return
	}
	if !strings.Contains(tiktokURL, "vm.tiktxk.com") && !strings.Contains(tiktokURL, "vm.tiktok.com") {
		handleError(c, "Invalid url", randomErrorImage)
		return
	}
	videoID, err := extracting.ExtractVideoID(tiktokURL)
	if err != nil {
		handleError(c, "Invalid url", randomErrorImage)
		return
	}
	filename := "collage-" + videoID + ".jpg"
	authorName, caption, responseBody, err := extracting.GetVideoAuthorAndCaption(tiktokURL, videoID)
	if err != nil {
		handleError(c, "Couldn't get video author and caption", randomErrorImage)
		return
	}

	if _, err := os.Stat("collages/" + filename); err == nil {
		if debug == "true" {
			elapsed := time.Since(startTime)
			caption = caption + " | Took " + elapsed.String()
		}
		handleDiscordEmbed(c, authorName, caption, filename)
		return
	}
	links, err := extracting.ExtractImageLinks(responseBody)
	if err != nil {
		handleError(c, "Couldn't get image links", randomErrorImage)
		return
	}
	err = httputil.DownloadImages(links, videoID)
	if err != nil {
		handleError(c, "Couldn't download images", randomErrorImage)
		return
	}
	err = collaging.MakeCollage(videoID, filename, width, initHeight)
	if err != nil {
		handleError(c, "Couldn't make collage", randomErrorImage)
		return
	}
	if debug == "true" {
		elapsed := time.Since(startTime)
		caption = caption + " | Took " + elapsed.String()
	}
	handleDiscordEmbed(c, authorName, caption, filename)
	os.RemoveAll(videoID)

}
func HandleCollage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleError(c, "No id provided", errorImages[errorImagesIndex])
		return
	}
	filename := "collage-" + id
	if _, err := os.Stat("collages/" + filename); err != nil {
		handleError(c, "Collage not found", errorImages[errorImagesIndex])
		return
	}
	c.File("collages/" + filename)

}
