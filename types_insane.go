package main

type BitrateImages struct {
	Name         string `json:"name"`
	BitrateImage Cover  `json:"bitrate_image"`
}

type ImagePostInfo struct {
	Images         []Image `json:"images"`
	ImagePostCover Image   `json:"image_post_cover"`
	PostExtra      string  `json:"post_extra"`
}

type Cover struct {
	URI       string   `json:"uri"`
	URLList   []string `json:"url_list"`
	Width     int      `json:"width"`
	Height    int      `json:"height"`
	URLPrefix any      `json:"url_prefix"`
}

type Image struct {
	DisplayImage        Cover           `json:"display_image"`
	OwnerWatermarkImage Cover           `json:"owner_watermark_image"`
	UserWatermarkImage  Cover           `json:"user_watermark_image"`
	Thumbnail           Cover           `json:"thumbnail"`
	BitrateImages       []BitrateImages `json:"bitrate_images"`
}

type VideoPlayAddr struct {
	URI       string   `json:"uri"`
	URLList   []string `json:"url_list"`
	Width     int      `json:"width"`
	Height    int      `json:"height"`
	URLKey    string   `json:"url_key"`
	URLPrefix any      `json:"url_prefix"`
}
type Author struct {
	UID       string `json:"uid"`
	ShortID   string `json:"short_id"`
	Nickname  string `json:"nickname"`
	Signature string `json:"signature"`
	UniqueID  string `json:"unique_id"`
}

type Music struct {
	ID       int64  `json:"id"`
	IDStr    string `json:"id_str"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	PlayURL  Cover  `json:"play_url"`
	Duration int    `json:"duration"`
}

type Video struct {
	PlayAddr      VideoPlayAddr `json:"play_addr"`
	Cover         Cover         `json:"cover"`
	Height        int           `json:"height"`
	Width         int           `json:"width"`
	DynamicCover  Cover         `json:"dynamic_cover"`
	OriginCover   Cover         `json:"origin_cover"`
	Ratio         string        `json:"ratio"`
	DownloadAddr  Cover         `json:"download_addr"`
	HasWatermark  bool          `json:"has_watermark"`
	BitRate       []any         `json:"bit_rate"`
	Duration      int           `json:"duration"`
	CdnURLExpired int           `json:"cdn_url_expired"`
	NeedSetToken  bool          `json:"need_set_token"`
	Tags          any           `json:"tags"`
	BigThumbs     any           `json:"big_thumbs"`
	IsBytevc1     int           `json:"is_bytevc1"`
	Meta          string        `json:"meta"`
	BitRateAudio  []any         `json:"bit_rate_audio"`
}

type Statistics struct {
	AwemeID            string `json:"aweme_id"`
	CommentCount       int    `json:"comment_count"`
	DiggCount          int    `json:"digg_count"`
	DownloadCount      int    `json:"download_count"`
	PlayCount          int    `json:"play_count"`
	ShareCount         int    `json:"share_count"`
	ForwardCount       int    `json:"forward_count"`
	LoseCount          int    `json:"lose_count"`
	LoseCommentCount   int    `json:"lose_comment_count"`
	WhatsappShareCount int    `json:"whatsapp_share_count"`
	CollectCount       int    `json:"collect_count"`
}

type TextExtra struct {
	Start       int    `json:"start"`
	End         int    `json:"end"`
	UserID      string `json:"user_id"`
	Type        int    `json:"type"`
	HashtagName string `json:"hashtag_name"`
	HashtagID   string `json:"hashtag_id"`
	IsCommerce  bool   `json:"is_commerce"`
	SecUID      string `json:"sec_uid"`
}

type Aweme struct {
	AwemeID       string        `json:"aweme_id"`
	Desc          string        `json:"desc"`
	CreateTime    int           `json:"create_time"`
	Author        Author        `json:"author,omitempty"`
	Music         Music         `json:"music,omitempty"`
	Video         Video         `json:"video,omitempty"`
	Statistics    Statistics    `json:"statistics"`
	ImagePostInfo ImagePostInfo `json:"image_post_info,omitempty"`
}

type TikTokAPIResponse struct {
	AwemeList []Aweme `json:"aweme_list"`
}

type HtmlVideo struct {
	ID           string `json:"id"`
	Height       int    `json:"height"`
	Width        int    `json:"width"`
	Duration     int    `json:"duration"`
	PlayAddr     string `json:"playAddr"`
	DownloadAddr string `json:"downloadAddr"`
}

type HtmlAuthor struct {
	ID       string `json:"id"`
	ShortID  string `json:"shortId"`
	UniqueID string `json:"uniqueId"`
	Nickname string `json:"nickname"`
}

type PreciseDuration struct {
	PreciseDuration         float64 `json:"preciseDuration"`
	PreciseShootDuration    float64 `json:"preciseShootDuration"`
	PreciseAuditionDuration float64 `json:"preciseAuditionDuration"`
	PreciseVideoDuration    float64 `json:"preciseVideoDuration"`
}
type HtmlMusic struct {
	ID              string          `json:"id"`
	Title           string          `json:"title"`
	PlayURL         string          `json:"playUrl"`
	AuthorName      string          `json:"authorName"`
	Duration        int             `json:"duration"`
	PreciseDuration PreciseDuration `json:"preciseDuration"`
}

type HtmlImage struct {
	ImageURL struct {
		URLList []string `json:"urlList"`
	} `json:"imageURL"`
	ImageWidth  int `json:"imageWidth"`
	ImageHeight int `json:"imageHeight"`
}

type HtmlImagePost struct {
	Images     []HtmlImage `json:"images"`
	Cover      HtmlImage   `json:"cover"`
	ShareCover HtmlImage   `json:"shareCover"`
	Title      string      `json:"title"`
}

type TiktokHTMLScript struct {
	ID           string     `json:"id"`
	Desc         string     `json:"desc"`
	CreateTime   string     `json:"createTime"`
	ScheduleTime int        `json:"scheduleTime"`
	Video        HtmlVideo  `json:"video"`
	Author       HtmlAuthor `json:"author"`
	Music        HtmlMusic  `json:"music"`
	Stats        struct {
		DiggCount    int    `json:"diggCount"`
		ShareCount   int    `json:"shareCount"`
		CommentCount int    `json:"commentCount"`
		PlayCount    int    `json:"playCount"`
		CollectCount string `json:"collectCount"`
	} `json:"stats"`
	ImagePost HtmlImagePost `json:"imagePost"`
	Contents  []struct {
		Desc      string      `json:"desc"`
		TextExtra []TextExtra `json:"textExtra"`
	} `json:"contents"`
}

type TikWmResp struct {
	Code          int     `json:"code"`
	Msg           string  `json:"msg"`
	ProcessedTime float64 `json:"processed_time"`
	Data          struct {
		ID        string   `json:"id"`
		Region    string   `json:"region"`
		Title     string   `json:"title"`
		Cover     string   `json:"cover"`
		Duration  int      `json:"duration"`
		Play      string   `json:"play"`
		Wmplay    string   `json:"wmplay"`
		Hdplay    string   `json:"hdplay"`
		Size      int      `json:"size"`
		WmSize    int      `json:"wm_size"`
		HdSize    int      `json:"hd_size"`
		Images    []string `json:"images"`
		Music     string   `json:"music"`
		MusicInfo struct {
			ID       string `json:"id"`
			Title    string `json:"title"`
			Play     string `json:"play"`
			Author   string `json:"author"`
			Original bool   `json:"original"`
			Duration int    `json:"duration"`
			Album    string `json:"album"`
		} `json:"music_info"`
		PlayCount     int         `json:"play_count"`
		DiggCount     int         `json:"digg_count"`
		CommentCount  int         `json:"comment_count"`
		ShareCount    int         `json:"share_count"`
		DownloadCount int         `json:"download_count"`
		CollectCount  int         `json:"collect_count"`
		CreateTime    int         `json:"create_time"`
		Anchors       interface{} `json:"anchors"`
		AnchorsExtras string      `json:"anchors_extras"`
		IsAd          bool        `json:"is_ad"`
		CommerceInfo  struct {
			AdvPromotable          bool `json:"adv_promotable"`
			AuctionAdInvited       bool `json:"auction_ad_invited"`
			BrandedContentType     int  `json:"branded_content_type"`
			WithCommentFilterWords bool `json:"with_comment_filter_words"`
		} `json:"commerce_info"`
		CommercialVideoInfo string `json:"commercial_video_info"`
		ItemCommentSettings int    `json:"item_comment_settings"`
		Author              struct {
			ID       string `json:"id"`
			UniqueID string `json:"unique_id"`
			Nickname string `json:"nickname"`
			Avatar   string `json:"avatar"`
		} `json:"author"`
	} `json:"data"`
}
