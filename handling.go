package main

import (
	"os"
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
	Private    bool
}

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

func handleDiscordEmbed(c *gin.Context, tiktokData Data, filename string) {
	details := tiktokData.Details
	detailsString := "❤️ " + details.Likes + " | 💬 " + details.Comments + " | 🔁 " + details.Shares + " | 👀 " + details.Views
	renderTemplate(c, "discord.html", gin.H{
		"authorName": tiktokData.AuthorName,
		"caption":    tiktokData.Caption,
		"details":    detailsString,
		"imageUrl":   Domain + "/" + filename,
	})
}

func handleVideoDiscordEmbed(
	c *gin.Context,
	tiktokData Data,
	filename string,
	width string,
	height string,
) {
	details := tiktokData.Details
	detailsString := "❤️" + details.Likes + " | 💬 " + details.Comments + " | 🔁 " + details.Shares + " | 👀 " + details.Views
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
		if video {
			videoWidth, videoHeight, err := GetVideoDimensions("collages/" + filename)
			if err != nil {
				HandleError(c, "Couldn't get video dimensions")
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

	tiktokData, err := FetchTiktokData(videoId)
	if err != nil {
		HandleError(c, "Couldn't fetch TikTok data")
		return
	}

	videoFilename := "video-" + tiktokData.VideoID + ".mp4"
	if handleExistingFile(c, videoFilename, true, tiktokData) {
		return
	}

	collageFilename := "collage-" + tiktokData.VideoID + ".png"
	collageFileExists, _ := os.Stat("collages/" + collageFilename)
	if collageFileExists != nil {
		err = FetchAudio(tiktokData.Body, videoId)
		if err != nil {
			HandleError(c, "Couldn't fetch audio")
			return
		}
		videoWidth, videoHeight, err := GenerateVideo(
			videoId,
			collageFilename,
			videoFilename,
			false,
		)
		if err != nil {
			HandleError(c, "Couldn't generate video")
			return
		}
		handleVideoDiscordEmbed(c, tiktokData, videoFilename, videoWidth, videoHeight)
		return
	}

	err = FetchImages(tiktokData.Body, videoId)
	if err != nil {
		HandleError(c, "Couldn't fetch images")
		return
	}

	err = FetchAudio(tiktokData.Body, videoId)
	if err != nil {
		HandleError(c, "Couldn't fetch audio")
		return
	}

	err = GenerateCollage(videoId, collageFilename)
	if err != nil {
		HandleError(c, "Couldn't generate collage")
		return
	}

	videoWidth, videoHeight, err := GenerateVideo(
		videoId,
		collageFilename,
		videoFilename,
		false,
	)
	if err != nil {
		HandleError(c, "Couldn't generate video")
		return
	}
	handleVideoDiscordEmbed(c, tiktokData, videoFilename, videoWidth, videoHeight)

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
	tiktokData, err := FetchTiktokData(videoId)
	if tiktokData.Private {
		// i dont understand i added println every where in the code and it still doesnt print
		// but it always returns "Couldn't fetch TikTok data" and i dont know why
		// test url https://vm.tiktok.com/ZGeDD2kT3/
		HandleError(c, "This TikTok is private")
		return
	}
	if err != nil {
		HandleError(c, "Couldn't fetch TikTok data")
		return
	}

	filename := "collage-" + tiktokData.VideoID + ".png"
	if handleExistingFile(c, filename, false, tiktokData) {
		return
	}
	err = FetchImages(tiktokData.Body, videoId)
	if err != nil {
		HandleError(c, "Couldn't fetch images")
		return
	}
	err = GenerateCollage(videoId, filename)
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
	tiktokData, err := FetchTiktokData(videoId)
	if err != nil {
		HandleError(c, "Couldn't fetch TikTok data")
		return
	}

	filename := "slide-" + videoId + ".mp4"
	if handleExistingFile(c, filename, true, tiktokData) {
		return
	}

	FetchImages(tiktokData.Body, videoId)
	FetchAudio(tiktokData.Body, videoId)

	videoWidth, videoHeight, err := GenerateVideo(videoId, "", filename, true)
	if err != nil {
		HandleError(c, "Couldn't generate video")
		return
	}

	handleVideoDiscordEmbed(c, tiktokData, filename, videoWidth, videoHeight)

	os.RemoveAll(tiktokData.VideoID)
	UpdateLocalStats()
}

/*
	func HandleSlideIndexRequest(c *gin.Context) {
		tiktokURLAndIndex := c.Query("v")
		tiktokURL, index, sound := SplitURLAndIndex(tiktokURLAndIndex)
		println(tiktokURL, index, sound)

		tiktokData := FetchTiktokData(c, tiktokURL)

		filename := "sIndex-" + tiktokData.VideoID + "-" + index + ".png"
		if handleExistingFile(c, filename, false, tiktokData) && !sound {
			return
		}

		videoFilename := "sIndex-" + tiktokData.VideoID + "-" + index + ".mp4"
		if handleExistingFile(c, videoFilename, true, tiktokData) {
			return
		}

		FetchImages(c, tiktokURL, tiktokData, index)
		if sound {
			FetchAudio(c, tiktokURL, tiktokData)
		}

		GenerateCollage(c, tiktokData.VideoID, filename)
		if sound {
			videoWidth, videoHeight := GenerateVideo(
				c,
				tiktokData.VideoID,
				filename,
				videoFilename,
				false,
			)
			handleVideoDiscordEmbed(c, tiktokData, videoFilename, videoWidth, videoHeight)
			os.RemoveAll(tiktokData.VideoID)
			return
		}
		handleDiscordEmbed(c, tiktokData, filename)
		os.RemoveAll(tiktokData.VideoID)
		UpdateLocalStats()
	}
*/

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
