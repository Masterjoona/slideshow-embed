package playerapi

import (
	"encoding/json"
	"fmt"
	"io"
	"meow/pkg/types"
	"meow/pkg/util"
	"meow/pkg/vars"
	"net/http"
	"strconv"
)

func FetchTikok(videoId string) (types.TiktokInfo, error) {
	_, _, err := fetch(videoId)
	if err != nil {
		return types.TiktokInfo{}, err
	}

	return types.TiktokInfo{}, nil
}

func (a *GenericItem) getVideoDetails() types.Counts {
	return types.Counts{
		Likes:     util.FormatLargeNumbers(strconv.Itoa(a.Stats.DiggCount)),
		Comments:  util.FormatLargeNumbers(strconv.Itoa(a.Stats.CommentCount)),
		Shares:    util.FormatLargeNumbers(strconv.Itoa(a.Stats.ShareCount)),
		Views:     util.FormatLargeNumbers(strconv.Itoa(a.Stats.PlayCount)),
		Favorites: util.FormatLargeNumbers(strconv.Itoa(a.Stats.CollectCount)),
	}
}

func fetch(SingleItemId string) (GenericItem, ImageItem, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.tiktok.com/api/related/item_list/?itemID="+SingleItemId, nil)
	// tiktok also calls https://www.tiktok.com/player/api/v1/items?item_ids=
	// but it misses some data like view and favorite count
	if err != nil {
		fmt.Println("Error creating request:", err)
		return GenericItem{}, ImageItem{}, err
	}

	req.Header.Set("user-agent", vars.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return GenericItem{}, ImageItem{}, err
	}
	defer resp.Body.Close()

	textByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return GenericItem{}, ImageItem{}, err
	}

	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(textByte), &raw); err != nil {
		fmt.Println("Error unmarshaling:", err)
		return GenericItem{}, ImageItem{}, err
	}

	if items, ok := raw["items"]; ok {
		length := len(items.([]interface{}))
		if length == 0 {
			return GenericItem{}, ImageItem{}, fmt.Errorf("No items found")
		}

		if length == 1 {
			var imageResp = PlayerAPIRespImage{}
			if err := json.Unmarshal([]byte(textByte), &imageResp); err != nil {
				fmt.Println("Error unmarshaling video item:", err)
				return GenericItem{}, ImageItem{}, err
			}
			return GenericItem{}, imageResp.Items[0], nil
		}

		var videoResp = GenericResp{}
		if err := json.Unmarshal([]byte(textByte), &videoResp); err != nil {
			fmt.Println("Error unmarshaling video item:", err)
			return GenericItem{}, ImageItem{}, err
		}

		return videoResp.ItemList[0], ImageItem{}, nil
	}

	return GenericItem{}, ImageItem{}, fmt.Errorf("No items found")
}

func (t *ImageItem) getImageLinks() []string {
	images := t.ImagePostInfo.Images
	imageLinks := make([]string, 0, len(images))
	for _, image := range images {
		imageLinks = append(imageLinks, image.DisplayImage.URLList[1])
	}
	return imageLinks
}
