package handler

import (
	"fmt"
	"meow/pkg/config"
	"meow/pkg/files"
	"meow/pkg/net"
	"meow/pkg/providers"
	"meow/pkg/types"
	"meow/pkg/util"
	"meow/pkg/vars"
	"strconv"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
)

func renderTemplate(c *gin.Context, filename string, data gin.H) {
	tmpl, err := template.ParseFiles("templates/" + filename)
	if err != nil {
		msg := fmt.Sprintf("Couldn't parse template: %s", filename)
		HandleError(c, msg, err)
		return
	}

	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		msg := fmt.Sprintf("Couldn't execute template: %s", filename)
		HandleError(c, msg, err)
	}
}

func HandleError(c *gin.Context, errorMsg string, err error) {
	if err != nil {
		println(err.Error())
	}
	renderTemplate(c, "error.html", gin.H{
		"error":           errorMsg,
		"error_image_url": ErrorImage(),
	})
}
func handleDiscordEmbed(c *gin.Context, t types.TiktokInfo, imageUrl string) {
	renderTemplate(c, "discord.html", gin.H{
		"authorName": t.Author,
		"caption":    t.Caption,
		"details":    t.Details.ToString(),
		"imageUrl":   imageUrl,
	})
}

func handleDiscordVideoEmbed(
	c *gin.Context,
	t types.TiktokInfo,
	videoUrl string,
) {
	renderTemplate(c, "video.html", gin.H{
		"authorName": strings.Split(t.Author, "(@")[0],
		"details":    t.Details.ToString(),
		"caption":    t.Caption,
		"videoUrl":   videoUrl,
		"width":      t.Video.Width,
		"height":     t.Video.Height,
	})
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

	renderTemplate(c, "index.html", gin.H{
		"FileLinks": config.LocalStats.FilePaths,
		"count":     config.LocalStats.FileCount,
		"size":      config.LocalStats.TotalSize,
	})
}

func HandleDirectFile(fileType string) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			HandleError(c, "No id provided", nil)
			return
		}
		filename := fmt.Sprintf("collages/%s-%s", fileType, id)
		if _, err := files.GetFileSize(filename); err != nil {
			HandleError(c, "File not found", nil)
			return
		}
		c.File(filename)
	}
}

func HandleDownloader(c *gin.Context) {
	t, errored := getTiktokData(c, "", false)
	if errored {
		return
	}

	if t.Video.Width != "" {
		handleDiscordVideoEmbed(
			c,
			t,
			vars.PathVideoProxy+"/"+t.VideoID,
		)
		return
	}

	renderTemplate(c, "images.html", gin.H{
		"authorName": t.Author,
		"caption":    t.Caption,
		"details":    t.Details.ToString(),
		"imageLinks": t.ImageLinks,
		"imageUrl":   t.ImageLinks[0],
		"soundUrl":   t.SoundLink,
	})
}

func getTiktokData(c *gin.Context, filePrefix string, isVideo bool) (types.TiktokInfo, bool) {
	tiktokURL := c.Query("v")

	uniqueUserId, videoId, err := net.GetLongVideoId(tiktokURL)
	if err != nil {
		if err.Error() == "invalid URL" {
			HandleError(c, "link: "+tiktokURL+" is invalid", nil)
		} else {
			HandleError(c, "Couldn't get tiktok", err)
		}
		return types.TiktokInfo{}, true
	}

	t, err := providers.FetchTiktok(videoId)

	if err != nil {
		HandleError(c, "Couldn't get the tiktok", err)
		return types.TiktokInfo{}, true
	}

	if !strings.Contains(t.Author, uniqueUserId) {
		t.Caption += "\n\ntiktok returned a different user, is the post available?"
	}

	hasFile := filePrefix != ""

	if hasFile {
		if t.Video.Width != "" && filePrefix[1] != 'u' {
			if strings.HasSuffix(c.Request.URL.String(), ".mp4") {
				c.Redirect(302, vars.PathVideoProxy+"/"+videoId)
				return types.TiktokInfo{}, true
			}
			handleDiscordVideoEmbed(
				c,
				t,
				vars.PathVideoProxy+"/"+videoId,
			)
			return types.TiktokInfo{}, true
		}

		fileExt := util.Ternary(filePrefix[0] == 'c', ".png", ".mp4")
		fileName := fmt.Sprintf("%s-%s%s", filePrefix, videoId, fileExt)
		t.FileName = fileName
		filePath := "collages/" + fileName

		if _, err := files.GetFileSize(filePath); err == nil {
			if isVideo {
				if isAwemeBeingRendered(videoId) {
					HandleError(c, "This video is being rendered, please request again in some time!", nil)
					return types.TiktokInfo{}, true
				}

				width, height, err := files.GetVideoDimensions(filePath)

				if err != nil {
					HandleError(c, "Couldn't get video dimensions", nil)
					return types.TiktokInfo{}, true
				}

				t.Video = types.SimplifiedVideo{
					Url:    config.Domain + fileName,
					Width:  width,
					Height: height,
				}

				handleDiscordVideoEmbed(c, t, config.Domain+fileName)
			} else {
				handleDiscordEmbed(c, t, config.Domain+fileName)
			}
			return types.TiktokInfo{}, true
		}

		t.ImageBuffers = net.DownloadImages(t.VideoID, t.ImageLinks)
		if isVideo {
			t.SoundBuffer, err = net.DownloadSound(t.SoundLink)
			if err != nil {
				HandleError(c, "Couldn't download sound", err)
				return types.TiktokInfo{}, true
			}
		}
	}

	return t, false
}

func HandleJsonRequest(c *gin.Context) {
	t, errored := getTiktokData(c, "", false)
	if !errored {
		t.Video.Buffer = nil
		c.JSON(200, t)
	}
}

func HandleVideoProxy(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		HandleError(c, "No id provided", nil)
		return
	}

	_, err := strconv.Atoi(idStr)
	if err != nil {
		HandleError(c, "Invalid id", nil)
		return
	}

	t, ok := providers.RecentTiktokReqs.Get(idStr)
	if !ok {
		var err error
		t, err = providers.FetchTiktok(idStr)
		if err != nil {
			HandleError(c, "Couldn't fetch TikTok data", err)
			return
		}
	}

	c.Header("Content-Type", "video/mp4")
	c.Header("Content-Disposition", "inline; filename="+t.FileName)
	c.Header("Content-Length", strconv.Itoa(len(t.Video.Buffer)))
	c.Header("Accept-Ranges", "bytes")
	c.Data(200, "video/mp4", t.Video.Buffer)
}

func HandleRequest(c *gin.Context) {
	t, errored := getTiktokData(c, "collage", false)
	if errored {
		return
	}

	err := t.MakeCollage()
	if err != nil {
		HandleError(c, "Couldn't generate collage", err)
		return
	}

	handleDiscordEmbed(c, t, config.Domain+t.FileName)
	config.UpdateLocalStats()
}

func HandleSoundCollageRequest(c *gin.Context) {
	t, errored := getTiktokData(c, "video", true)
	if errored {
		return
	}

	collageFilePath := "collages/collage-" + t.VideoID + ".png"
	if _, err := files.GetFileSize(collageFilePath); err != nil {
		err := t.MakeCollage()
		if err != nil {
			HandleError(c, "Couldn't generate collage", err)
			return
		}
	}

	width, height, err := t.MakeCollageWithAudio("video")
	if err != nil {
		HandleError(c, "Couldn't generate video", err)
		return
	}

	t.Video = types.SimplifiedVideo{
		Url:    config.Domain + t.FileName,
		Width:  width,
		Height: height,
	}

	handleDiscordVideoEmbed(c, t, config.Domain+t.FileName)
	config.UpdateLocalStats()
}

func HandleFancySlideshowRequest(c *gin.Context) {
	t, errored := getTiktokData(c, "slide", true)
	if errored {
		return
	}

	if len(t.ImageLinks) == 1 {
		c.Redirect(302, vars.PathCollageSound+"?v="+c.Query("v"))
		return
	}

	addAwemeToRendering(t.VideoID)

	go func() {
		_, _, err := t.MakeVideoSlideshow()
		if err != nil {
			println(err.Error())
			return
		}

		config.UpdateLocalStats()
		removeAwemeFromRendering(t.VideoID)
	}()

	HandleError(c, "This slideshow was sent to be rendered, please request again in some time!", nil)
}

func HandleSubtitleVideo(c *gin.Context) {
	subLang := c.Query("lang")
	if subLang == "" {
		HandleError(c, "No language provided", nil)
		return
	}

	t, errored := getTiktokData(c, "subs-"+subLang, true)
	if errored {
		return
	}

	err := net.DownloadVideoAndSubtitles(t.VideoID, t.FileName, subLang)
	if err != nil {
		// This is due to tiktok goofery
		errorMsg := "Couldn't download video with subtitles. Only translations are available. e.g if the non-translated subtitles are in English, you can only get translations in other languages."
		HandleError(c, errorMsg, err)
		return
	}

	addAwemeToRendering(t.VideoID)
	go func() {
		_, _, err := t.MakeVideoSubtitles(subLang)
		if err != nil {
			println(err.Error())
			return
		}
		config.UpdateLocalStats()
		removeAwemeFromRendering(t.VideoID)
	}()

	HandleError(c, "This video was sent to be subtitled, please request again in some time!", nil)
}
