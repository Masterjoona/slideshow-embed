package tiktok_api

import (
	"encoding/json"
	"fmt"
	"io"
	"meow/pkg/types"
	"meow/pkg/util"
	"meow/pkg/vars"
	"net/http"
	"strconv"
	"strings"
)

const maxRetries = 5

func FetchTikok(videoId string) (types.TiktokInfo, error) {
	postAweme, err := FetchTiktokAPI(videoId)
	if err != nil {
		return types.TiktokInfo{}, err
	}

	videoUrl := postAweme.Video.PlayAddr.URLList[0]
	isVideo := !strings.Contains(videoUrl, "music")
	imageLinks := []string{}

	if !isVideo {
		imageLinks = postAweme.getImageLinks()
	}

	return types.TiktokInfo{
		Author:     postAweme.Author.Nickname + " (@" + postAweme.Author.UniqueID + ")",
		Caption:    postAweme.Desc,
		VideoID:    videoId,
		Details:    postAweme.getVideoDetails(),
		ImageLinks: imageLinks,
		SoundLink:  postAweme.Music.PlayURL.URI,
		Video: types.SimplifiedVideo{
			Url:    videoUrl,
			Width:  strconv.Itoa(postAweme.Video.Width),
			Height: strconv.Itoa(postAweme.Video.Height),
		},
	}, nil
}

func (a *Aweme) getVideoDetails() types.Counts {
	return types.Counts{
		Likes:     util.FormatLargeNumbers(strconv.Itoa(a.Statistics.DiggCount)),
		Comments:  util.FormatLargeNumbers(strconv.Itoa(a.Statistics.CommentCount)),
		Shares:    util.FormatLargeNumbers(strconv.Itoa(a.Statistics.ShareCount)),
		Views:     util.FormatLargeNumbers(strconv.Itoa(a.Statistics.PlayCount)),
		Favorites: util.FormatLargeNumbers(strconv.Itoa(a.Statistics.CollectCount)),
	}
}

func FetchTiktokAPI(awemeId string) (Aweme, error) {
	req, err := http.NewRequest("OPTIONS", "https://api22-normal-c-alisg.tiktokv.com/aweme/v1/feed/?aweme_id="+awemeId, nil)
	// yes, options is correct, it actually works troll
	if err != nil {
		return Aweme{}, err
	}

	req.Header.Set("user-agent", vars.UserAgent)

	for attempts := 1; attempts <= maxRetries; attempts++ {
		resp, err := vars.HttpClient.Do(req)
		if err != nil {
			return Aweme{}, err
		}
		defer resp.Body.Close()

		textByte, err := io.ReadAll(resp.Body)
		if err != nil {
			return Aweme{}, err
		}

		var response TikTokAPIResponse
		err = json.Unmarshal(textByte, &response)
		if err != nil {
			fmt.Printf("Error unmarshalling (attempt %d/%d): %v\n", attempts, maxRetries, err)
			if attempts == maxRetries {
				return Aweme{}, err
			}
			continue
		}
		return response.AwemeList[0], nil
	}

	return Aweme{}, fmt.Errorf("failed to unmarshal response after %d attempts", maxRetries)
}

func (t *Aweme) getImageLinks() []string {
	images := t.ImagePostInfo.Images
	imageLinks := make([]string, 0, len(images))
	for _, image := range images {
		imageLinks = append(imageLinks, image.DisplayImage.URLList[1])
	}
	return imageLinks
}
