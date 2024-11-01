package ttsave

import (
	"fmt"
	"io"
	provider_util "meow/pkg/providers/util"
	"meow/pkg/types"
	"meow/pkg/util"
	"meow/pkg/vars"
	"net/http"
	"regexp"
	"strings"
)

var AudioSrcRe = regexp.MustCompile(`<a href="(.*)" onclick="bdl\(this, event\)" type="audio`)
var VideoSrcLinkRe = regexp.MustCompile(`<a href="(.*)" onclick="bdl\(this, event\)" rel`)

func fetchTTSave(tiktokUrl string) (*string, error) {
	var data = strings.NewReader(fmt.Sprintf(`{"language_id":"1","query":"%s"}`, tiktokUrl))
	req, err := http.NewRequest("POST", "https://ttsave.app/download", data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	resp, err := vars.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	text := string(bodyText)
	return &text, nil
}

func getData(body *string) (string, string, types.Counts, error) {
	authorRe := regexp.MustCompile(`mb-2">(.*)</a>`)
	match := authorRe.FindStringSubmatch(*body)
	if len(match) == 0 {
		return "", "", types.Counts{}, fmt.Errorf("could not find author")
	}
	author := match[1]

	nicknameRe := regexp.MustCompile(`text-center">(.*)</h2>`)
	match = nicknameRe.FindStringSubmatch(*body)
	author = match[1] + " (" + author + ")"

	captionRe := regexp.MustCompile(`oneliner">(.*)<\/p>`)
	match = captionRe.FindStringSubmatch(*body)
	caption := match[1]

	countRe := regexp.MustCompile(`text-gray-500">(.*)</span>`)
	matches := countRe.FindAllStringSubmatch(*body, -1)

	return author, caption, types.Counts{
		Views:     matches[0][1],
		Likes:     matches[1][1],
		Comments:  matches[2][1],
		Favorites: matches[3][1],
		Shares:    matches[4][1],
	}, nil
}

func getMediaLinks(body *string) (string, string) {
	return VideoSrcLinkRe.FindStringSubmatch(*body)[1], AudioSrcRe.FindStringSubmatch(*body)[1]
}

func getSlideLinks(body *string) []string {
	re := regexp.MustCompile(`<img src="([^"]*)">`)
	matches := re.FindAllStringSubmatch(*body, -1)
	var slideLinks []string
	for _, match := range matches {
		slideLinks = append(slideLinks, match[1])
	}
	return slideLinks
}

func FetchTiktok(videoId string) (types.TiktokInfo, error) {
	url := "https://www.tiktok.com/@placeholder/video/" + videoId

	data, err := fetchTTSave(url)
	if err != nil {
		return types.TiktokInfo{}, err
	}

	slideLinks := getSlideLinks(data)
	author, caption, stats, err := getData(data)
	if err != nil {
		return types.TiktokInfo{}, err
	}

	videoSrc, audioSrc := getMediaLinks(data)

	videoUrl := util.Ternary(len(slideLinks) == 0, "", videoSrc)
	videoInfo, err := provider_util.GetDimensionsOrNil(videoUrl, videoUrl != "")
	if err != nil {
		return types.TiktokInfo{}, err
	}

	return types.TiktokInfo{
		Author:     author,
		Caption:    caption,
		Details:    stats,
		VideoID:    videoId,
		SoundLink:  audioSrc,
		ImageLinks: slideLinks,
		Video:      videoInfo,
	}, nil
}
