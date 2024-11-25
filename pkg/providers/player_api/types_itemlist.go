package playerapi

type ItemListResp struct {
	Cursor                 string        `json:"cursor"`
	Extra                  Extra         `json:"extra"`
	ItemList               []GenericItem `json:"itemList"`
	LogPb                  LogPb         `json:"log_pb"`
	StatusCode             int64         `json:"statusCode"`
	ItemListRespStatusCode int64         `json:"status_code"`
	StatusMsg              string        `json:"status_msg"`
}

type GenericItem struct {
	AIGCDescription            string      `json:"AIGCDescription"`
	CategoryType               int64       `json:"CategoryType"`
	Author                     Author      `json:"author"`
	AuthorStats                AuthorStats `json:"authorStats"`
	BackendSourceEventTracking string      `json:"backendSourceEventTracking"`
	Challenges                 []Challenge `json:"challenges"`
	Collected                  bool        `json:"collected"`
	Contents                   []Content   `json:"contents"`
	CreateTime                 int64       `json:"createTime"`
	Desc                       string      `json:"desc"`
	Digged                     bool        `json:"digged"`
	DiversificationID          int64       `json:"diversificationId"`
	DuetDisplay                int64       `json:"duetDisplay"`
	DuetEnabled                bool        `json:"duetEnabled"`
	ForFriend                  bool        `json:"forFriend"`
	ID                         string      `json:"id"`
	ItemCommentStatus          int64       `json:"itemCommentStatus"`
	ItemControl                ItemControl `json:"item_control"`
	Music                      Music       `json:"music"`
	OfficalItem                bool        `json:"officalItem"`
	OriginalItem               bool        `json:"originalItem"`
	PrivateItem                bool        `json:"privateItem"`
	Secret                     bool        `json:"secret"`
	ShareEnabled               bool        `json:"shareEnabled"`
	Stats                      Stats       `json:"stats"`
	StatsV2                    StatsV2     `json:"statsV2"`
	StitchDisplay              int64       `json:"stitchDisplay"`
	StitchEnabled              bool        `json:"stitchEnabled"`
	TextExtra                  []TextExtra `json:"textExtra"`
	Video                      Video       `json:"video"`
	PlaylistID                 *string     `json:"playlistId,omitempty"`
}

type Author struct {
	AvatarLarger    string `json:"avatarLarger"`
	AvatarMedium    string `json:"avatarMedium"`
	AvatarThumb     string `json:"avatarThumb"`
	CommentSetting  int64  `json:"commentSetting"`
	DownloadSetting int64  `json:"downloadSetting"`
	DuetSetting     int64  `json:"duetSetting"`
	Ftc             bool   `json:"ftc"`
	ID              string `json:"id"`
	IsADVirtual     bool   `json:"isADVirtual"`
	IsEmbedBanned   bool   `json:"isEmbedBanned"`
	Nickname        string `json:"nickname"`
	OpenFavorite    bool   `json:"openFavorite"`
	PrivateAccount  bool   `json:"privateAccount"`
	Relation        int64  `json:"relation"`
	SECUid          string `json:"secUid"`
	Secret          bool   `json:"secret"`
	Signature       string `json:"signature"`
	StitchSetting   int64  `json:"stitchSetting"`
	UniqueID        string `json:"uniqueId"`
	Verified        bool   `json:"verified"`
}

type AuthorStats struct {
	DiggCount      int64 `json:"diggCount"`
	FollowerCount  int64 `json:"followerCount"`
	FollowingCount int64 `json:"followingCount"`
	FriendCount    int64 `json:"friendCount"`
	Heart          int64 `json:"heart"`
	HeartCount     int64 `json:"heartCount"`
	VideoCount     int64 `json:"videoCount"`
}

type Challenge struct {
	CoverLarger   string `json:"coverLarger"`
	CoverMedium   string `json:"coverMedium"`
	CoverThumb    string `json:"coverThumb"`
	Desc          string `json:"desc"`
	ID            string `json:"id"`
	ProfileLarger string `json:"profileLarger"`
	ProfileMedium string `json:"profileMedium"`
	ProfileThumb  string `json:"profileThumb"`
	Title         string `json:"title"`
}

type Content struct {
	Desc      string      `json:"desc"`
	TextExtra []TextExtra `json:"textExtra"`
}

type TextExtra struct {
	AwemeID     string `json:"awemeId"`
	End         int64  `json:"end"`
	HashtagName string `json:"hashtagName"`
	IsCommerce  bool   `json:"isCommerce"`
	Start       int64  `json:"start"`
	SubType     int64  `json:"subType"`
	Type        int64  `json:"type"`
}

type ItemControl struct {
	CanRepost bool `json:"can_repost"`
}

type Music struct {
	AuthorName  string  `json:"authorName"`
	CoverLarge  string  `json:"coverLarge"`
	CoverMedium string  `json:"coverMedium"`
	CoverThumb  string  `json:"coverThumb"`
	Duration    int64   `json:"duration"`
	ID          string  `json:"id"`
	Original    bool    `json:"original"`
	PlayURL     string  `json:"playUrl"`
	Title       string  `json:"title"`
	Album       *string `json:"album,omitempty"`
}

type Stats struct {
	CollectCount int64 `json:"collectCount"`
	CommentCount int64 `json:"commentCount"`
	DiggCount    int64 `json:"diggCount"`
	PlayCount    int64 `json:"playCount"`
	ShareCount   int64 `json:"shareCount"`
}

type StatsV2 struct {
	CollectCount string `json:"collectCount"`
	CommentCount string `json:"commentCount"`
	DiggCount    string `json:"diggCount"`
	PlayCount    string `json:"playCount"`
	RepostCount  string `json:"repostCount"`
	ShareCount   string `json:"shareCount"`
}

type Video struct {
	VQScore       string            `json:"VQScore"`
	Bitrate       int64             `json:"bitrate"`
	BitrateInfo   []BitrateInfo     `json:"bitrateInfo"`
	ClaInfoList   ClaInfoList       `json:"ClaInfoList"`
	CodecType     string            `json:"codecType"`
	Cover         string            `json:"cover"`
	Definition    string            `json:"definition"`
	DownloadAddr  string            `json:"downloadAddr"`
	Duration      int64             `json:"duration"`
	DynamicCover  string            `json:"dynamicCover"`
	EncodeUserTag string            `json:"encodeUserTag"`
	EncodedType   string            `json:"string"`
	Format        string            `json:"format"`
	Height        int64             `json:"height"`
	ID            string            `json:"id"`
	OriginCover   string            `json:"originCover"`
	PlayAddr      string            `json:"playAddr"`
	Ratio         string            `json:"ratio"`
	SubtitleInfos []SubtitleInfo    `json:"subtitleInfos,omitempty"`
	VideoQuality  string            `json:"videoQuality"`
	VolumeInfo    VolumeInfo        `json:"volumeInfo"`
	Width         int64             `json:"width"`
	ZoomCover     map[string]string `json:"zoomCover"`
}

type BitrateInfo struct {
	Bitrate     int64    `json:"Bitrate"`
	CodecType   string   `json:"CodecType"`
	GearName    string   `json:"GearName"`
	Mvmaf       string   `json:"MVMAF"`
	PlayAddr    PlayAddr `json:"PlayAddr"`
	QualityType int64    `json:"QualityType"`
}

type PlayAddr struct {
	DataSize int64    `json:"DataSize"`
	FileCS   string   `json:"FileCs"`
	FileHash string   `json:"FileHash"`
	Height   int64    `json:"Height"`
	URI      string   `json:"Uri"`
	URLKey   string   `json:"UrlKey"`
	URLList  []string `json:"UrlList"`
	Width    int64    `json:"Width"`
}

type ClaInfoList struct {
	CaptionInfoLists     []CaptionInfoList     `json:"CaptionInfoLists,omitempty"`
	CaptionsType         *int64                `json:"captionsType,omitempty"`
	EnableAutoCaption    bool                  `json:"enableAutoCaption"`
	HasOriginalAudio     bool                  `json:"hasOriginalAudio"`
	OriginalLanguageInfo *OriginalLanguageInfo `json:"originalLanguageInfo,omitempty"`
	NoCaptionReason      *int64                `json:"noCaptionReason,omitempty"`
}

type CaptionInfoList struct {
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
}

type OriginalLanguageInfo struct {
	Language     string `json:"language"`
	LanguageCode string `json:"languageCode"`
	LanguageID   string `json:"languageID"`
}

type SubtitleInfo struct {
	Format           string `json:"Format"`
	LanguageCodeName string `json:"LanguageCodeName"`
	LanguageID       string `json:"LanguageID"`
	Size             int64  `json:"Size"`
	Source           string `json:"Source"`
	URL              string `json:"Url"`
	URLExpire        int64  `json:"UrlExpire"`
	Version          string `json:"Version"`
}

type VolumeInfo struct {
	Loudness float64 `json:"Loudness"`
	Peak     float64 `json:"Peak"`
}
