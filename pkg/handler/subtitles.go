package handler

import (
	"encoding/json"
	"errors"
	"io"
	"meow/pkg/config"
	"meow/pkg/files"
	"meow/pkg/net"
	playerapi "meow/pkg/providers/player_api"
	"meow/pkg/vars"
	"net/http"
	"os"
)

// i dont like this being here
// but i must so that there is import cycle

type SubtitlesResp struct {
	ServerTransTime      int    `json:"server_trans_time"`
	StatusCode           int    `json:"status_code"`
	StatusMsg            string `json:"status_msg"`
	TranslationCacheTime int    `json:"translation_cache_time"`
	Variant              string `json:"variant"`
	WebvttSubtitle       string `json:"webvtt_subtitle"`
}

func DownloadVideoAndSubtitles(videoId, videoUrl, fileName, lang string) error {
	subtitles, err := fetchSubtitles(videoId, lang)
	if err != nil {
		return err
	}

	videoTmpDir := config.TmpCollageDir + videoId
	if err := files.CreateDirectory(videoTmpDir); err != nil {
		return err
	}

	if err := os.WriteFile(videoTmpDir+"/subtitles.vtt", []byte(subtitles), 0644); err != nil {
		return err
	}

	videoBytes, err := net.DownloadMedia(videoUrl)
	if err != nil {
		return err
	}

	return os.WriteFile(videoTmpDir+"/video.mp4", videoBytes, 0644)
}

func fetchSubtitles(videoId, lang string) (string, error) {
	url := vars.SubtitlesHost + "subtitle_id=02981317794434464&target_language=" + lang + "&item_id=" + videoId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", vars.UserAgent)

	resp, err := vars.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var subtitlesResp SubtitlesResp
	if err := json.Unmarshal(bodyText, &subtitlesResp); err != nil {
		return "", err
	}

	if subtitlesResp.WebvttSubtitle != "" {
		return subtitlesResp.WebvttSubtitle, nil
	}

	data, err := playerapi.FetchSingleAPI(videoId)
	if err != nil {
		return "", err
	}

	if data.VideoInfo.ClaInfo.CaptionInfos[0].LanguageCode == lang {
		req, err = http.NewRequest("GET", data.VideoInfo.ClaInfo.CaptionInfos[0].URL, nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("User-Agent", vars.UserAgent)

		resp, err = vars.HttpClient.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		subtitlesBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		return string(subtitlesBytes), nil
	}

	return "", errors.New("no subtitles found")
}
