package playerapi

type GenericItem struct {
	Author struct {
		Nickname string `json:"nickname"`
		UniqueID string `json:"uniqueId"`
	} `json:"author"`
	Desc  string `json:"desc"`
	ID    string `json:"id"`
	Music struct {
		PlayURL string `json:"playUrl"`
	} `json:"music"`
	Stats struct {
		CollectCount int `json:"collectCount"`
		CommentCount int `json:"commentCount"`
		DiggCount    int `json:"diggCount"`
		PlayCount    int `json:"playCount"`
		ShareCount   int `json:"shareCount"`
	} `json:"stats"`
	Video struct {
		VQScore     string `json:"VQScore"`
		Bitrate     int    `json:"bitrate"`
		BitrateInfo []struct {
			Bitrate   int    `json:"Bitrate"`
			CodecType string `json:"CodecType"`
			GearName  string `json:"GearName"`
			MVMAF     string `json:"MVMAF"`
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
			CaptionInfos []struct {
				CaptionFormat     string   `json:"captionFormat"`
				ClaSubtitleID     string   `json:"claSubtitleID"`
				Expire            string   `json:"expire"`
				IsAutoGen         bool     `json:"isAutoGen"`
				IsOriginalCaption bool     `json:"isOriginalCaption"`
				Language          string   `json:"language"`
				LanguageCode      string   `json:"languageCode"`
				LanguageID        string   `json:"languageID"`
				SubID             string   `json:"subID"`
				SubtitleType      string   `json:"subtitleType"`
				URL               string   `json:"url"`
				URLList           []string `json:"urlList"`
				Variant           string   `json:"variant"`
			} `json:"captionInfos"`
			CaptionsType         int  `json:"captionsType"`
			EnableAutoCaption    bool `json:"enableAutoCaption"`
			HasOriginalAudio     bool `json:"hasOriginalAudio"`
			OriginalLanguageInfo struct {
				Language     string `json:"language"`
				LanguageCode string `json:"languageCode"`
				LanguageID   string `json:"languageID"`
			} `json:"originalLanguageInfo"`
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
		SubtitleInfos []struct {
			Format           string `json:"Format"`
			LanguageCodeName string `json:"LanguageCodeName"`
			LanguageID       string `json:"LanguageID"`
			Size             int    `json:"Size"`
			Source           string `json:"Source"`
			URL              string `json:"Url"`
			URLExpire        int    `json:"UrlExpire"`
			Version          string `json:"Version"`
		} `json:"subtitleInfos"`
		VideoQuality string `json:"videoQuality"`
		VolumeInfo   struct {
			Loudness int     `json:"Loudness"`
			Peak     float64 `json:"Peak"`
		} `json:"volumeInfo"`
		Width     int `json:"width"`
		ZoomCover struct {
			Num240 string `json:"240"`
			Num480 string `json:"480"`
			Num720 string `json:"720"`
			Num960 string `json:"960"`
		} `json:"zoomCover"`
	} `json:"video"`
}

type GenericResp struct {
	ItemList []GenericItem `json:"itemList"`
}
