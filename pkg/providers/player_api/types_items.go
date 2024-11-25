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

type SingleItemResp struct {
	Extra      Extra        `json:"extra"`
	Items      []SingleItem `json:"items"`
	LogPb      LogPb        `json:"log_pb"`
	Results    []Result     `json:"results"`
	StatusCode int64        `json:"status_code"`
	StatusMsg  string       `json:"status_msg"`
}

type Extra struct {
	FatalItemIDS []interface{} `json:"fatal_item_ids"`
	Logid        string        `json:"logid"`
	Now          int64         `json:"now"`
}

type SingleItem struct {
	AuthorInfo     AuthorInfo     `json:"author_info"`
	AwemeType      int64          `json:"aweme_type"`
	Desc           string         `json:"desc"`
	ID             float64        `json:"id"`
	IDStr          string         `json:"id_str"`
	ImagePostInfo  ImagePostInfo  `json:"image_post_info,omitempty"`
	MarkerInfo     MarkerInfo     `json:"marker_info"`
	MusicInfo      MusicInfo      `json:"music_info"`
	OtherInfo      OtherInfo      `json:"other_info"`
	Region         string         `json:"region"`
	StatisticsInfo StatisticsInfo `json:"statistics_info"`
	VideoInfo      VideoInfo      `json:"video_info"`
}

type AuthorInfo struct {
	AvatarURLList []string `json:"avatar_url_list"`
	Nickname      string   `json:"nickname"`
	SecretID      string   `json:"secret_id"`
	UniqueID      string   `json:"unique_id"`
}

type MarkerInfo struct {
	BrandedContentType int64 `json:"branded_content_type"`
	IsAds              bool  `json:"is_ads"`
}

type MusicInfo struct {
	Author string  `json:"author"`
	ID     float64 `json:"id"`
	IDStr  string  `json:"id_str"`
	Title  string  `json:"title"`
}

type OtherInfo struct {
}

type StatisticsInfo struct {
	CommentCount int64 `json:"comment_count"`
	DiggCount    int64 `json:"digg_count"`
	ShareCount   int64 `json:"share_count"`
}

type VideoInfo struct {
	ClaInfo ClaInfo  `json:"cla_info,omitempty"`
	Meta    Meta     `json:"meta"`
	URI     string   `json:"uri"`
	URLList []string `json:"url_list"`
}

type ClaInfo struct {
	CaptionInfos []CaptionInfo `json:"caption_infos"`
}

type CaptionInfo struct {
	CaptionFormat     string   `json:"caption_format"`
	CaptionLength     int64    `json:"caption_length"`
	ClaSubtitleID     float64  `json:"cla_subtitle_id"`
	ComplaintID       float64  `json:"complaint_id"`
	Expire            int64    `json:"expire"`
	IsAutoGenerated   bool     `json:"is_auto_generated"`
	IsOriginalCaption bool     `json:"is_original_caption"`
	Lang              string   `json:"lang"`
	LanguageCode      string   `json:"language_code"`
	LanguageID        int64    `json:"language_id"`
	SourceTag         string   `json:"source_tag"`
	SubID             int64    `json:"sub_id"`
	SubVersion        string   `json:"sub_version"`
	SubtitleType      int64    `json:"subtitle_type"`
	URL               string   `json:"url"`
	URLList           []string `json:"url_list"`
	Variant           string   `json:"variant"`
}

type Meta struct {
	Bitrate  int64 `json:"bitrate,omitempty"`
	Duration int64 `json:"duration"`
	Height   int64 `json:"height"`
	Ratio    int64 `json:"ratio"`
	Width    int64 `json:"width"`
}

type ImagePostInfo struct {
	Cover  Cover   `json:"cover"`
	Images []Cover `json:"images"`
}

type Cover struct {
	DisplayImage        DisplayImage `json:"display_image"`
	OwnerWatermarkImage DisplayImage `json:"owner_watermark_image"`
	Thumbnail           DisplayImage `json:"thumbnail"`
}

type DisplayImage struct {
	Height  int64    `json:"height"`
	URLList []string `json:"url_list"`
	Width   int64    `json:"width"`
}

type LogPb struct {
	ImprID string `json:"impr_id"`
}

type Result struct {
	Code  string  `json:"code"`
	ID    float64 `json:"id"`
	IDStr string  `json:"id_str"`
}
