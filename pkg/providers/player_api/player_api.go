package playerapi

import (
	"encoding/json"
	"fmt"
	"io"
	provider_util "meow/pkg/providers/util"
	"meow/pkg/types"
	"meow/pkg/util"
	"meow/pkg/vars"
	"net/http"
)

func FetchTikok(videoId string) (types.TiktokInfo, error) {
	data, err := fetchList(videoId)
	if err != nil {
		return types.TiktokInfo{}, err
	}

	singleItemData, err := fetchSingle(videoId)
	if err != nil {
		return types.TiktokInfo{}, err
	}

	isVideo := singleItemData.VideoInfo.Meta.Duration != 0
	author := data.Author.Nickname + " (@" + data.Author.UniqueID + ")"
	videoInfo, err := provider_util.GetDimensionsOrNil(singleItemData.VideoInfo.URLList[0], isVideo)
	if err != nil {
		return types.TiktokInfo{}, err
	}

	return types.TiktokInfo{
		Author:     author,
		Caption:    data.Desc,
		Details:    data.getItemDetails(),
		VideoID:    videoId,
		ImageLinks: singleItemData.getImageLinks(),
		SoundLink:  data.Music.PlayURL,
		Video:      videoInfo,
	}, nil
}

func (i *GenericItem) getItemDetails() types.Counts {
	return types.Counts{
		Likes:     util.FormatLargeNumbers(i.StatsV2.DiggCount),
		Comments:  util.FormatLargeNumbers(i.StatsV2.CommentCount),
		Shares:    util.FormatLargeNumbers(i.StatsV2.ShareCount),
		Views:     util.FormatLargeNumbers(i.StatsV2.PlayCount),
		Favorites: util.FormatLargeNumbers(i.StatsV2.CollectCount),
	}
}

func fetchSingle(videoId string) (SingleItem, error) {
	req, err := http.NewRequest("GET", "https://www.tiktok.com/player/api/v1/items?item_ids="+videoId, nil)
	if err != nil {
		return SingleItem{}, err
	}

	req.Header.Set("user-agent", vars.UserAgent)

	resp, err := vars.HttpClient.Do(req)
	if err != nil {
		return SingleItem{}, err
	}
	defer resp.Body.Close()

	textByte, err := io.ReadAll(resp.Body)

	if err != nil {
		return SingleItem{}, err
	}

	var data SingleItemResp
	err = json.Unmarshal(textByte, &data)
	if err != nil {
		return SingleItem{}, err
	}

	if len(data.Items) == 0 {
		return SingleItem{}, fmt.Errorf("no items found")
	}

	return data.Items[0], nil
}

func fetchList(videoId string) (GenericItem, error) {
	req, err := http.NewRequest("GET", "https://www.tiktok.com/api/related/item_list/?aid=1284&count=14&itemID="+videoId, nil)
	if err != nil {
		return GenericItem{}, err
	}

	req.Header.Set("user-agent", vars.UserAgent)

	resp, err := vars.HttpClient.Do(req)
	if err != nil {
		return GenericItem{}, err
	}
	defer resp.Body.Close()

	textByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return GenericItem{}, err
	}

	var data ItemListResp
	err = json.Unmarshal(textByte, &data)
	if err != nil {
		return GenericItem{}, err
	}

	return data.ItemList[0], nil
}
