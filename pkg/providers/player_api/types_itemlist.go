package playerapi

type GenericItem struct {
	AIGCDescription     string `json:"AIGCDescription"`
	ExploreCategoryType int    `json:"ExploreCategoryType"`
	Author              struct {
		AvatarLarger    string `json:"avatarLarger"`
		AvatarMedium    string `json:"avatarMedium"`
		AvatarThumb     string `json:"avatarThumb"`
		CommentSetting  int    `json:"commentSetting"`
		DownloadSetting int    `json:"downloadSetting"`
		DuetSetting     int    `json:"duetSetting"`
		Ftc             bool   `json:"ftc"`
		ID              string `json:"id"`
		IsADVirtual     bool   `json:"isADVirtual"`
		IsEmbedBanned   bool   `json:"isEmbedBanned"`
		Nickname        string `json:"nickname"`
		OpenFavorite    bool   `json:"openFavorite"`
		PrivateAccount  bool   `json:"privateAccount"`
		Relation        int    `json:"relation"`
		SecUID          string `json:"secUid"`
		Secret          bool   `json:"secret"`
		Signature       string `json:"signature"`
		StitchSetting   int    `json:"stitchSetting"`
		UniqueID        string `json:"uniqueId"`
		Verified        bool   `json:"verified"`
	} `json:"author"`
	CreateTime        int    `json:"createTime"`
	Desc              string `json:"desc"`
	Digged            bool   `json:"digged"`
	DiversificationID int    `json:"diversificationId"`
	DuetDisplay       int    `json:"duetDisplay"`
	DuetEnabled       bool   `json:"duetEnabled,omitempty"`
	ForFriend         bool   `json:"forFriend"`
	ID                string `json:"id"`
	ItemCommentStatus int    `json:"itemCommentStatus"`
	ItemControl       struct {
		CanRepost bool `json:"can_repost"`
	} `json:"item_control"`
	Music struct {
		Album       string `json:"album"`
		AuthorName  string `json:"authorName"`
		CoverLarge  string `json:"coverLarge"`
		CoverMedium string `json:"coverMedium"`
		CoverThumb  string `json:"coverThumb"`
		Duration    int    `json:"duration"`
		ID          string `json:"id"`
		Original    bool   `json:"original"`
		PlayURL     string `json:"playUrl"`
		Title       string `json:"title"`
	} `json:"music,omitempty"`
	StatsV2 struct {
		CollectCount string `json:"collectCount"`
		CommentCount string `json:"commentCount"`
		DiggCount    string `json:"diggCount"`
		PlayCount    string `json:"playCount"`
		RepostCount  string `json:"repostCount"`
		ShareCount   string `json:"shareCount"`
	} `json:"statsV2"`
	Video struct {
		VQScore     string `json:"VQScore"`
		Bitrate     int    `json:"bitrate"`
		BitrateInfo []struct {
			Bitrate   int    `json:"Bitrate"`
			CodecType string `json:"CodecType"`
			GearName  string `json:"GearName"`
			Mvmaf     string `json:"MVMAF"`
			PlayAddr  struct {
				DataSize int      `json:"DataSize"`
				FileCs   string   `json:"FileCs"`
				FileHash string   `json:"FileHash"`
				Height   int      `json:"Height"`
				URI      string   `json:"Uri"`
				URLKey   string   `json:"UrlKey"`
				URLList  []string `json:"UrlList"`
				Width    int      `json:"Width"`
			} `json:"PlayAddr"`
			QualityType int `json:"QualityType"`
		} `json:"bitrateInfo"`
		ClaInfo struct {
			EnableAutoCaption bool `json:"enableAutoCaption"`
			HasOriginalAudio  bool `json:"hasOriginalAudio"`
			NoCaptionReason   int  `json:"noCaptionReason"`
		} `json:"claInfo"`
		CodecType     string `json:"codecType"`
		Cover         string `json:"cover"`
		Definition    string `json:"definition"`
		DownloadAddr  string `json:"downloadAddr"`
		Duration      int    `json:"duration"`
		DynamicCover  string `json:"dynamicCover"`
		EncodeUserTag string `json:"encodeUserTag"`
		EncodedType   string `json:"encodedType"`
		Format        string `json:"format"`
		Height        int    `json:"height"`
		ID            string `json:"id"`
		OriginCover   string `json:"originCover"`
		PlayAddr      string `json:"playAddr"`
		Ratio         string `json:"ratio"`
		VideoQuality  string `json:"videoQuality"`
		VolumeInfo    struct {
			Loudness float64 `json:"Loudness"`
			Peak     float64 `json:"Peak"`
		} `json:"volumeInfo"`
		Width     int `json:"width"`
		ZoomCover struct {
			Num240 string `json:"240"`
			Num480 string `json:"480"`
			Num720 string `json:"720"`
			Num960 string `json:"960"`
		} `json:"zoomCover"`
	} `json:"video,omitempty"`
	PlaylistID string `json:"playlistId,omitempty"`
}

type ItemListResp struct {
	ItemList []GenericItem `json:"itemList"`
}
