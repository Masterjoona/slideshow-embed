package playerapi

type ImageItem struct {
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
			DisplayImage struct {
				Height  int      `json:"height"`
				URLList []string `json:"url_list"`
				Width   int      `json:"width"`
			} `json:"display_image"`
			OwnerWatermarkImage struct {
				Height  int      `json:"height"`
				URLList []string `json:"url_list"`
				Width   int      `json:"width"`
			} `json:"owner_watermark_image"`
			Thumbnail struct {
				Height  int      `json:"height"`
				URLList []string `json:"url_list"`
				Width   int      `json:"width"`
			} `json:"thumbnail"`
		} `json:"cover"`
		Images []struct {
			DisplayImage struct {
				Height  int      `json:"height"`
				URLList []string `json:"url_list"`
				Width   int      `json:"width"`
			} `json:"display_image"`
			OwnerWatermarkImage struct {
				Height  int      `json:"height"`
				URLList []string `json:"url_list"`
				Width   int      `json:"width"`
			} `json:"owner_watermark_image"`
			Thumbnail struct {
				Height  int      `json:"height"`
				URLList []string `json:"url_list"`
				Width   int      `json:"width"`
			} `json:"thumbnail"`
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
	OtherInfo struct {
	} `json:"other_info"`
	Region         string `json:"region"`
	StatisticsInfo struct {
		CommentCount int `json:"comment_count"`
		DiggCount    int `json:"digg_count"`
		ShareCount   int `json:"share_count"`
	} `json:"statistics_info"`
	VideoInfo struct {
		Meta struct {
			Duration int `json:"duration"`
			Height   int `json:"height"`
			Ratio    int `json:"ratio"`
			Width    int `json:"width"`
		} `json:"meta"`
		URI     string   `json:"uri"`
		URLList []string `json:"url_list"`
	} `json:"video_info"`
}

type PlayerAPIRespImage struct {
	Items []ImageItem `json:"items"`
}
