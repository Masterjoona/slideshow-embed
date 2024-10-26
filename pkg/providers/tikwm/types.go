package tikwm

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
