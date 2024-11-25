package tikwm

import (
	"encoding/json"
	"io"
	provider_util "meow/pkg/providers/util"
	"meow/pkg/types"
	"meow/pkg/util"
	"meow/pkg/vars"
	"net/http"
	"strconv"
	"strings"
)

func FetchTikwm(videoId string) (TikWmResp, error) {
	var data = strings.NewReader("url=https://www.tiktok.com/@placeholder/video/" + videoId)
	req, err := http.NewRequest("POST", "https://tikwm.com/api/", data)
	if err != nil {
		return TikWmResp{}, err
	}

	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := vars.HttpClient.Do(req)
	if err != nil {
		return TikWmResp{}, err
	}

	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return TikWmResp{}, err
	}

	var tikWmResp TikWmResp
	err = json.Unmarshal(bodyText, &tikWmResp)
	if err != nil {
		return TikWmResp{}, err
	}

	return tikWmResp, nil
}

func FetchTiktok(videoId string) (types.TiktokInfo, error) {
	tikWmApiResp, err := FetchTikwm(videoId)
	data := tikWmApiResp.Data

	if err != nil {
		return types.TiktokInfo{}, err
	}

	stats := types.Counts{
		Likes:     util.FormatLargeNumbers(strconv.Itoa(data.DiggCount)),
		Comments:  util.FormatLargeNumbers(strconv.Itoa(data.CommentCount)),
		Shares:    util.FormatLargeNumbers(strconv.Itoa(data.ShareCount)),
		Views:     util.FormatLargeNumbers(strconv.Itoa(data.PlayCount)),
		Favorites: util.FormatLargeNumbers(strconv.Itoa(data.CollectCount)),
	}

	author := data.Author.Nickname + " (@" + data.Author.UniqueID + ")"
	caption := data.Title

	videoUrl := util.Ternary(data.Duration != 0, data.Play, "")
	videoInfo, err := provider_util.GetDimensionsOrNil(videoUrl, data.Duration != 0)
	if err != nil {
		return types.TiktokInfo{}, err
	}

	return types.TiktokInfo{
		Author:     author,
		Caption:    caption,
		Details:    stats,
		VideoID:    videoId,
		SoundLink:  data.MusicInfo.Play,
		ImageLinks: util.Ternary(data.Duration == 0, data.Images, nil),
		Video:      videoInfo,
	}, nil
}
