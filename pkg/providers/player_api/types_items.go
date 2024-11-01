package playerapi

func (t *SingleItem) getImageLinks() []string {
	images := t.ImagePostInfo.Images
	if len(images) == 0 {
		return nil
	}
	imageLinks := make([]string, 0, len(images))
	for _, image := range images {
		imageLinks = append(imageLinks, image.DisplayImage.URLList[1])
	}
	return imageLinks
}

type ImageDetails struct {
	Height  int      `json:"height"`
	Width   int      `json:"width"`
	URLList []string `json:"url_list"`
}

type SingleItem struct {
	AuthorInfo struct {
		AvatarURLList []string `json:"avatar_url_list"`
		Nickname      string   `json:"nickname"`
		SecretID      string   `json:"secret_id"`
		UniqueID      string   `json:"unique_id"`
	} `json:"author_info"`
	AwemeType     int    `json:"aweme_type"`
	Desc          string `json:"desc"`
	ID            int64  `json:"id"`
	IDStr         string `json:"id_str"`
	ImagePostInfo struct {
		Cover struct {
			DisplayImage        ImageDetails `json:"display_image"`
			OwnerWatermarkImage ImageDetails `json:"owner_watermark_image"`
			Thumbnail           ImageDetails `json:"thumbnail"`
		} `json:"cover"`
		Images []struct {
			DisplayImage        ImageDetails `json:"display_image"`
			OwnerWatermarkImage ImageDetails `json:"owner_watermark_image"`
			Thumbnail           ImageDetails `json:"thumbnail"`
		} `json:"images"`
	} `json:"image_post_info"`
	MarkerInfo struct {
		BrandedContentType int  `json:"branded_content_type"`
		IsAds              bool `json:"is_ads"`
	} `json:"marker_info"`
	MusicInfo struct {
		Author string `json:"author"`
		ID     int64  `json:"id"`
		IDStr  string `json:"id_str"`
		Title  string `json:"title"`
	} `json:"music_info"`
	VideoInfo struct {
		Meta struct {
			Bitrate  int `json:"bitrate"`
			Duration int `json:"duration"`
			Height   int `json:"height"`
			Ratio    int `json:"ratio"`
			Width    int `json:"width"`
		} `json:"meta"`
		URI     string   `json:"uri"`
		URLList []string `json:"url_list"`
	} `json:"video_info"`
}

type SingleItemResp struct {
	Extra struct {
		FatalItemIds []any  `json:"fatal_item_ids"`
		Logid        string `json:"logid"`
		Now          int64  `json:"now"`
	} `json:"extra"`
	Items []SingleItem `json:"items"`
	LogPb struct {
		ImprID string `json:"impr_id"`
	} `json:"log_pb"`
	Results []struct {
		Code  string `json:"code"`
		ID    int64  `json:"id"`
		IDStr string `json:"id_str"`
	} `json:"results"`
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}
