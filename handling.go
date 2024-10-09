package main

import (
	"fmt"
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

func generateDetailsString(details Counts) string {
	return fmt.Sprintf("❤️ %s | 💬 %s | 🔁 %s | ⭐ %s | 👀 %s",
		details.Likes, details.Comments, details.Shares, details.Favorites, details.Views)
}

func handleDiscordEmbed(c *gin.Context, tiktokData SimplifiedData, imageUrl string) {
	detailsString := generateDetailsString(tiktokData.Details)
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
	detailsString := generateDetailsString(tiktokData.Details)
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
		if _, err := GetFileSize(filename); err != nil {
			HandleError(c, "File not found")
			return
		}
		c.File(filename)
	}
}

func HandleDownloader(c *gin.Context) {
	tiktokData, errored := getTiktokData(c, "", false)
	if errored {
		return
	}

	detailsString := generateDetailsString(tiktokData.Details)

	if tiktokData.Video.Width != "" {
		handleVideoDiscordEmbed(c, tiktokData, tiktokData.Video.Url, tiktokData.Video.Width, tiktokData.Video.Height)
		return
	}

	renderTemplate(c, "images.html", gin.H{
		"authorName": tiktokData.Author,
		"caption":    tiktokData.Caption,
		"details":    detailsString,
		"imageLinks": tiktokData.ImageLinks,
		"imageUrl":   tiktokData.ImageLinks[0],
		"soundUrl":   tiktokData.SoundLink,
	})
}

func getTiktokData(c *gin.Context, filePrefix string, isVideo bool) (SimplifiedData, bool) {
	tiktokURL := c.Query("v")
	uniqueUserId, videoId, err := GetLongVideoId(tiktokURL)
	if err != nil {
		if err.Error() == "invalid URL" {
			HandleError(c, "link: "+tiktokURL+" is invalid")
		} else {
			HandleError(c, "Couldn't get tiktok")
		}
		return SimplifiedData{}, true
	}

	var tiktokData SimplifiedData

	if cachedData, ok := RecentTiktokReqs.Get(videoId); ok {
		tiktokData = cachedData
	} else {
		tiktokData, err = FetchTiktokData(videoId)
		if err != nil {
			println(err.Error())
			HandleError(c, "Couldn't get the tiktok")
			return SimplifiedData{}, true
		}
		tiktokData.DecodeStrings()
		RecentTiktokReqs.Put(videoId, tiktokData)
	}

	if !strings.Contains(tiktokData.Author, uniqueUserId) {
		tiktokData.Caption += "\n\ntiktok returned a different user, is the post available?"
	}

	hasFile := filePrefix != ""

	if hasFile {
		if tiktokData.Video.Width != "" && filePrefix[1] != 'u' {
			handleVideoDiscordEmbed(c, tiktokData, tiktokData.Video.Url, tiktokData.Video.Width, tiktokData.Video.Height)
			return SimplifiedData{}, true
		}

		var fileExt string

		if filePrefix[0] == 'c' {
			fileExt = ".png"
		} else {
			fileExt = ".mp4"
		}
		fileName := fmt.Sprintf("%s-%s%s", filePrefix, videoId, fileExt)
		tiktokData.FileName = fileName
		filePath := "collages/" + fileName
		if _, err := GetFileSize(filePath); err == nil {
			if isVideo {
				if IsAwemeBeingRendered(videoId) {
					HandleError(c, "This video is being rendered, please request again in some time!")
					return SimplifiedData{}, true
				}
				width, height, err := GetVideoDimensions(filePath)
				if err != nil {
					println(err.Error())
					HandleError(c, "Couldn't get video dimensions")
					return SimplifiedData{}, true
				}
				handleVideoDiscordEmbed(c, tiktokData, Domain+fileName, width, height)
			} else {
				handleDiscordEmbed(c, tiktokData, Domain+fileName)
			}
			return SimplifiedData{}, true
		}
		tiktokData.DownloadImages()
		if isVideo {
			tiktokData.DownloadSound()
		}
	}

	return tiktokData, false
}

func HandleJsonRequest(c *gin.Context) {
	tiktokData, errored := getTiktokData(c, "", false)
	if !errored {
		c.JSON(200, tiktokData)
	}
}

// not really a proxy but whatever
func HandleVideoProxy(c *gin.Context) {
	tiktokData, errored := getTiktokData(c, "", false)
	if errored {
		return
	}

	if tiktokData.Video.Width == "" {
		HandleError(c, "This is not a video tiktok")
		return
	}
	c.Redirect(302, tiktokData.Video.Url)
}

func HandleRequest(c *gin.Context) {
	tiktokData, errored := getTiktokData(c, "collage", false)
	if errored {
		return
	}
	tiktokData.MakeCollage()

	handleDiscordEmbed(c, tiktokData, Domain+tiktokData.FileName)
	UpdateLocalStats()
}

func HandleSoundCollageRequest(c *gin.Context) {
	tiktokData, errored := getTiktokData(c, "video", true)
	if errored {
		return
	}

	collageFilePath := "collages/collage-" + tiktokData.VideoID + ".png"
	if _, err := GetFileSize(collageFilePath); err != nil {
		tiktokData.MakeCollage()
	}

	width, height, err := tiktokData.MakeCollageWithAudio("video")
	if err != nil {
		println(err.Error())
		HandleError(c, "Couldn't generate video")
		return
	}

	handleVideoDiscordEmbed(c, tiktokData, Domain+tiktokData.FileName, width, height)
	UpdateLocalStats()
}

func HandleFancySlideshowRequest(c *gin.Context) {
	tiktokData, errored := getTiktokData(c, "slide", true)
	if errored {
		return
	}

	if len(tiktokData.ImageLinks) == 1 {
		c.Redirect(302, PathCollageSound+"?v="+c.Query("v"))
		return
	}

	AddAwemeToRendering(tiktokData.VideoID)

	go func() {
		_, _, err := tiktokData.MakeVideoSlideshow()
		if err != nil {
			println(err.Error())
			return
		}

		UpdateLocalStats()
		RemoveAwemeFromRendering(tiktokData.VideoID)
	}()

	HandleError(c, "This slideshow was sent to be rendered, please request again in some time!")
}

func HandleSubtitleVideo(c *gin.Context) {
	subLang := c.Query("lang")
	if subLang == "" {
		HandleError(c, "No language provided")
		return
	}

	tiktokData, errored := getTiktokData(c, "subs-"+subLang, true)
	if errored {
		return
	}

	err := tiktokData.DownloadVideoAndSubtitles(subLang)
	if err != nil {
		// This is due to tiktok goofery
		errorMsg := "Couldn't download video with subtitles. Only translations are available. e.g if the non-translated subtitles are in English, you can only get translations in other languages."
		HandleError(c, errorMsg)
		return
	}

	AddAwemeToRendering(tiktokData.VideoID)
	go func() {
		_, _, err := tiktokData.MakeVideoSubtitles(subLang)
		if err != nil {
			println(err.Error())
			return
		}
		UpdateLocalStats()
		RemoveAwemeFromRendering(tiktokData.VideoID)
	}()

	HandleError(c, "This video was sent to be subtitled, please request again in some time!")
}
