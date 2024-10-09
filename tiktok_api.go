package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func (a *Aweme) getVideoDetails() Counts {
	return Counts{
		Likes:     FormatLargeNumbers(strconv.Itoa(a.Statistics.DiggCount)),
		Comments:  FormatLargeNumbers(strconv.Itoa(a.Statistics.CommentCount)),
		Shares:    FormatLargeNumbers(strconv.Itoa(a.Statistics.ShareCount)),
		Views:     FormatLargeNumbers(strconv.Itoa(a.Statistics.PlayCount)),
		Favorites: FormatLargeNumbers(strconv.Itoa(a.Statistics.CollectCount)),
	}
}

func FetchTiktokDataTiktokAPI(videoId string) (SimplifiedData, error) {
	postAweme, err := fetch(videoId)
	if err != nil {
		return SimplifiedData{}, err
	}
	videoUrl := postAweme.Video.PlayAddr.URLList[0]
	isVideo := !strings.Contains(videoUrl, "music")
	imageLinks := []string{}
	if !isVideo {
		imageLinks = postAweme.getImageLinks()
	}
	return SimplifiedData{
		Author:     postAweme.Author.Nickname + " (@" + postAweme.Author.UniqueID + ")",
		Caption:    postAweme.Desc,
		VideoID:    videoId,
		Details:    postAweme.getVideoDetails(),
		ImageLinks: imageLinks,
		SoundLink:  postAweme.Music.PlayURL.URI,
		Video: SimplifiedVideo{
			Url:    videoUrl,
			Width:  strconv.Itoa(postAweme.Video.Width),
			Height: strconv.Itoa(postAweme.Video.Height),
		},
	}, nil
}

func fetch(awemeId string) (Aweme, error) {
	client := &http.Client{}
	req, err := http.NewRequest("OPTIONS", "https://api22-normal-c-alisg.tiktokv.com/aweme/v1/feed/?aweme_id="+awemeId, nil)
	// yes, options is correct, it actually returns data (most of the time)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return Aweme{}, err
	}

	req.Header.Set("user-agent", UserAgent)

	for attempts := 1; attempts <= MaxRetriesForTiktokAPI; attempts++ {
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return Aweme{}, err
		}
		defer resp.Body.Close()

		textByte, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return Aweme{}, err
		}

		var response TikTokAPIResponse
		err = json.Unmarshal(textByte, &response)
		if err != nil {
			fmt.Printf("Error unmarshalling (attempt %d/%d): %v\n", attempts, MaxRetriesForTiktokAPI, err)
			if attempts == MaxRetriesForTiktokAPI {
				return Aweme{}, err
			}
			continue
		}
		return response.AwemeList[0], nil
	}

	return Aweme{}, fmt.Errorf("failed to unmarshal response after %d attempts", MaxRetriesForTiktokAPI)
}

func (t *Aweme) getImageLinks() []string {
	images := t.ImagePostInfo.Images
	imageLinks := make([]string, 0, len(images))
	for _, image := range images {
		imageLinks = append(imageLinks, image.DisplayImage.URLList[1])
	}
	return imageLinks
}
