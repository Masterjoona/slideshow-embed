package main

type LogInfo struct {
	Order string `json:"order"`
}

type StatusInner struct {
	Status int `json:"status"`
}
type CommentConfig struct {
	EmojiRecommendList any `json:"emoji_recommend_list"`
}

type AigcInfo struct {
	AigcLabelType int `json:"aigc_label_type"`
}

type MatchedFriend struct {
	VideoItems any `json:"video_items"`
}

type LogPb struct {
	ImprID string `json:"impr_id"`
}

type LogInfoInner struct {
	ImprID   string `json:"impr_id"`
	PullType string `json:"pull_type"`
}

type BitrateImages struct {
	Name         string `json:"name"`
	BitrateImage Cover  `json:"bitrate_image"`
}

type OriginalClientText struct {
	MarkupText string           `json:"markup_text"`
	TextExtra  []TextExtraInner `json:"text_extra"`
}

type CommerceInfo struct {
	AdvPromotable      bool `json:"adv_promotable"`
	BrandedContentType int  `json:"branded_content_type"`
}
type Extra struct {
	Now          int64 `json:"now"`
	FatalItemIds any   `json:"fatal_item_ids"`
	APIDebugInfo any   `json:"api_debug_info"`
}
type ImagePostInfo struct {
	Images         []Image `json:"images"`
	ImagePostCover Image   `json:"image_post_cover"`
	PostExtra      string  `json:"post_extra"`
}

type TextExtraInner struct {
	Type        int    `json:"type"`
	HashtagName string `json:"hashtag_name"`
	SubType     int    `json:"sub_type"`
	TagID       string `json:"tag_id"`
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

type RiskInfos struct {
	Vote     bool   `json:"vote"`
	Warn     bool   `json:"warn"`
	RiskSink bool   `json:"risk_sink"`
	Type     int    `json:"type"`
	Content  string `json:"content"`
}

type Sharing struct {
	Code      int    `json:"code"`
	ShowType  int    `json:"show_type"`
	Extra     string `json:"extra"`
	Transcode int    `json:"transcode"`
	Mute      bool   `json:"mute"`
}

type VideoPlayAddr struct {
	URI       string   `json:"uri"`
	URLList   []string `json:"url_list"`
	Width     int      `json:"width"`
	Height    int      `json:"height"`
	URLKey    string   `json:"url_key"`
	URLPrefix any      `json:"url_prefix"`
}

type ContentDescExtra struct {
	Start       int    `json:"start"`
	End         int    `json:"end"`
	Type        int    `json:"type"`
	HashtagName string `json:"hashtag_name"`
	HashtagID   string `json:"hashtag_id"`
	IsCommerce  bool   `json:"is_commerce"`
	LineIdx     int    `json:"line_idx"`
}

type ShareInfo struct {
	ShareURL                   string `json:"share_url"`
	ShareDesc                  string `json:"share_desc"`
	ShareTitle                 string `json:"share_title"`
	ShareQrcodeURL             Cover  `json:"share_qrcode_url"`
	ShareTitleMyself           string `json:"share_title_myself"`
	ShareTitleOther            string `json:"share_title_other"`
	ShareDescInfo              string `json:"share_desc_info"`
	NowInvitationCardImageUrls any    `json:"now_invitation_card_image_urls"`
}

type Author struct {
	UID                        string        `json:"uid"`
	ShortID                    string        `json:"short_id"`
	Nickname                   string        `json:"nickname"`
	Signature                  string        `json:"signature"`
	AvatarThumb                Cover         `json:"avatar_thumb"`
	AvatarMedium               Cover         `json:"avatar_medium"`
	FollowStatus               int           `json:"follow_status"`
	IsBlock                    bool          `json:"is_block"`
	CustomVerify               string        `json:"custom_verify"`
	UniqueID                   string        `json:"unique_id"`
	RoomID                     int           `json:"room_id"`
	AuthorityStatus            int           `json:"authority_status"`
	VerifyInfo                 string        `json:"verify_info"`
	ShareInfo                  ShareInfo     `json:"share_info"`
	WithCommerceEntry          bool          `json:"with_commerce_entry"`
	VerificationType           int           `json:"verification_type"`
	EnterpriseVerifyReason     string        `json:"enterprise_verify_reason"`
	IsAdFake                   bool          `json:"is_ad_fake"`
	FollowersDetail            any           `json:"followers_detail"`
	Region                     string        `json:"region"`
	CommerceUserLevel          int           `json:"commerce_user_level"`
	PlatformSyncInfo           any           `json:"platform_sync_info"`
	IsDisciplineMember         bool          `json:"is_discipline_member"`
	Secret                     int           `json:"secret"`
	PreventDownload            bool          `json:"prevent_download"`
	Geofencing                 any           `json:"geofencing"`
	VideoIcon                  Cover         `json:"video_icon"`
	FollowerStatus             int           `json:"follower_status"`
	CommentSetting             int           `json:"comment_setting"`
	DuetSetting                int           `json:"duet_setting"`
	DownloadSetting            int           `json:"download_setting"`
	CoverURL                   []Cover       `json:"cover_url"`
	Language                   string        `json:"language"`
	ItemList                   any           `json:"item_list"`
	IsStar                     bool          `json:"is_star"`
	TypeLabel                  []any         `json:"type_label"`
	AdCoverURL                 any           `json:"ad_cover_url"`
	CommentFilterStatus        int           `json:"comment_filter_status"`
	Avatar168X168              Cover         `json:"avatar_168x168"`
	Avatar300X300              Cover         `json:"avatar_300x300"`
	RelativeUsers              any           `json:"relative_users"`
	ChaList                    any           `json:"cha_list"`
	SecUID                     string        `json:"sec_uid"`
	NeedPoints                 any           `json:"need_points"`
	HomepageBottomToast        any           `json:"homepage_bottom_toast"`
	CanSetGeofencing           any           `json:"can_set_geofencing"`
	WhiteCoverURL              any           `json:"white_cover_url"`
	UserTags                   any           `json:"user_tags"`
	BoldFields                 any           `json:"bold_fields"`
	SearchHighlight            any           `json:"search_highlight"`
	MutualRelationAvatars      any           `json:"mutual_relation_avatars"`
	Events                     any           `json:"events"`
	MatchedFriend              MatchedFriend `json:"matched_friend"`
	AdvanceFeatureItemOrder    any           `json:"advance_feature_item_order"`
	AdvancedFeatureInfo        any           `json:"advanced_feature_info"`
	UserProfileGuide           any           `json:"user_profile_guide"`
	ShieldEditFieldInfo        any           `json:"shield_edit_field_info"`
	CanMessageFollowStatusList any           `json:"can_message_follow_status_list"`
	AccountLabels              any           `json:"account_labels"`
}

type Music struct {
	ID                   int64   `json:"id"`
	IDStr                string  `json:"id_str"`
	Title                string  `json:"title"`
	Author               string  `json:"author"`
	Album                string  `json:"album"`
	CoverLarge           Cover   `json:"cover_large"`
	CoverMedium          Cover   `json:"cover_medium"`
	CoverThumb           Cover   `json:"cover_thumb"`
	PlayURL              Cover   `json:"play_url"`
	SourcePlatform       int     `json:"source_platform"`
	Duration             int     `json:"duration"`
	Extra                string  `json:"extra"`
	UserCount            int     `json:"user_count"`
	Position             any     `json:"position"`
	CollectStat          int     `json:"collect_stat"`
	Status               int     `json:"status"`
	OfflineDesc          string  `json:"offline_desc"`
	OwnerID              string  `json:"owner_id"`
	OwnerNickname        string  `json:"owner_nickname"`
	IsOriginal           bool    `json:"is_original"`
	Mid                  string  `json:"mid"`
	BindedChallengeID    int     `json:"binded_challenge_id"`
	AuthorDeleted        bool    `json:"author_deleted"`
	OwnerHandle          string  `json:"owner_handle"`
	AuthorPosition       any     `json:"author_position"`
	PreventDownload      bool    `json:"prevent_download"`
	ExternalSongInfo     []any   `json:"external_song_info"`
	SecUID               string  `json:"sec_uid"`
	AvatarThumb          Cover   `json:"avatar_thumb"`
	AvatarMedium         Cover   `json:"avatar_medium"`
	PreviewStartTime     float64 `json:"preview_start_time"`
	PreviewEndTime       float64 `json:"preview_end_time"`
	IsCommerceMusic      bool    `json:"is_commerce_music"`
	IsOriginalSound      bool    `json:"is_original_sound"`
	Artists              any     `json:"artists"`
	LyricShortPosition   any     `json:"lyric_short_position"`
	MuteShare            bool    `json:"mute_share"`
	TagList              any     `json:"tag_list"`
	IsAuthorArtist       bool    `json:"is_author_artist"`
	IsPgc                bool    `json:"is_pgc"`
	SearchHighlight      any     `json:"search_highlight"`
	MultiBitRatePlayInfo any     `json:"multi_bit_rate_play_info"`
	TtToDspSongInfos     any     `json:"tt_to_dsp_song_infos"`
	RecommendStatus      int     `json:"recommend_status"`
	UncertArtists        any     `json:"uncert_artists"`
}

type ChaList struct {
	Cid             string    `json:"cid"`
	ChaName         string    `json:"cha_name"`
	Desc            string    `json:"desc"`
	Schema          string    `json:"schema"`
	Author          Author    `json:"author"`
	UserCount       int       `json:"user_count"`
	ShareInfo       ShareInfo `json:"share_info"`
	ConnectMusic    []any     `json:"connect_music"`
	Type            int       `json:"type"`
	SubType         int       `json:"sub_type"`
	IsPgcshow       bool      `json:"is_pgcshow"`
	CollectStat     int       `json:"collect_stat"`
	IsChallenge     int       `json:"is_challenge"`
	ViewCount       int       `json:"view_count"`
	IsCommerce      bool      `json:"is_commerce"`
	HashtagProfile  string    `json:"hashtag_profile"`
	ChaAttrs        any       `json:"cha_attrs"`
	BannerList      any       `json:"banner_list"`
	ShowItems       any       `json:"show_items"`
	SearchHighlight any       `json:"search_highlight"`
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

type Status struct {
	AwemeID        string `json:"aweme_id"`
	IsDelete       bool   `json:"is_delete"`
	AllowShare     bool   `json:"allow_share"`
	AllowComment   bool   `json:"allow_comment"`
	IsPrivate      bool   `json:"is_private"`
	WithGoods      bool   `json:"with_goods"`
	PrivateStatus  int    `json:"private_status"`
	InReviewing    bool   `json:"in_reviewing"`
	Reviewed       int    `json:"reviewed"`
	SelfSee        bool   `json:"self_see"`
	IsProhibited   bool   `json:"is_prohibited"`
	DownloadStatus int    `json:"download_status"`
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

type ShareInfoInner struct {
	ShareURL                   string `json:"share_url"`
	ShareDesc                  string `json:"share_desc"`
	ShareTitle                 string `json:"share_title"`
	BoolPersist                int    `json:"bool_persist"`
	ShareTitleMyself           string `json:"share_title_myself"`
	ShareTitleOther            string `json:"share_title_other"`
	ShareLinkDesc              string `json:"share_link_desc"`
	ShareSignatureURL          string `json:"share_signature_url"`
	ShareSignatureDesc         string `json:"share_signature_desc"`
	ShareQuote                 string `json:"share_quote"`
	WhatsappDesc               string `json:"whatsapp_desc"`
	ShareDescInfo              string `json:"share_desc_info"`
	NowInvitationCardImageUrls any    `json:"now_invitation_card_image_urls"`
	ShareButtonDisplayMode     int    `json:"share_button_display_mode"`
}

type VideoControl struct {
	AllowDownload         bool `json:"allow_download"`
	ShareType             int  `json:"share_type"`
	ShowProgressBar       int  `json:"show_progress_bar"`
	DraftProgressBar      int  `json:"draft_progress_bar"`
	AllowDuet             bool `json:"allow_duet"`
	AllowReact            bool `json:"allow_react"`
	PreventDownloadType   int  `json:"prevent_download_type"`
	AllowDynamicWallpaper bool `json:"allow_dynamic_wallpaper"`
	TimerStatus           int  `json:"timer_status"`
	AllowMusic            bool `json:"allow_music"`
	AllowStitch           bool `json:"allow_stitch"`
}

type AwemeACL struct {
	DownloadGeneral   Sharing `json:"download_general"`
	DownloadMaskPanel Sharing `json:"download_mask_panel"`
	ShareListStatus   int     `json:"share_list_status"`
	ShareGeneral      Sharing `json:"share_general"`
	PlatformList      any     `json:"platform_list"`
	ShareActionList   any     `json:"share_action_list"`
	PressActionList   any     `json:"press_action_list"`
}

type InteractPermission struct {
	Duet                 int         `json:"duet"`
	Stitch               int         `json:"stitch"`
	DuetPrivacySetting   int         `json:"duet_privacy_setting"`
	StitchPrivacySetting int         `json:"stitch_privacy_setting"`
	Upvote               int         `json:"upvote"`
	AllowAddingToStory   int         `json:"allow_adding_to_story"`
	AllowCreateSticker   StatusInner `json:"allow_create_sticker"`
}

type AddedSoundMusicInfo struct {
	ID                   int64   `json:"id"`
	IDStr                string  `json:"id_str"`
	Title                string  `json:"title"`
	Author               string  `json:"author"`
	Album                string  `json:"album"`
	CoverLarge           Cover   `json:"cover_large"`
	CoverMedium          Cover   `json:"cover_medium"`
	CoverThumb           Cover   `json:"cover_thumb"`
	PlayURL              Cover   `json:"play_url"`
	SourcePlatform       int     `json:"source_platform"`
	Duration             int     `json:"duration"`
	Extra                string  `json:"extra"`
	UserCount            int     `json:"user_count"`
	Position             any     `json:"position"`
	CollectStat          int     `json:"collect_stat"`
	Status               int     `json:"status"`
	OfflineDesc          string  `json:"offline_desc"`
	OwnerID              string  `json:"owner_id"`
	OwnerNickname        string  `json:"owner_nickname"`
	IsOriginal           bool    `json:"is_original"`
	Mid                  string  `json:"mid"`
	BindedChallengeID    int     `json:"binded_challenge_id"`
	AuthorDeleted        bool    `json:"author_deleted"`
	OwnerHandle          string  `json:"owner_handle"`
	AuthorPosition       any     `json:"author_position"`
	PreventDownload      bool    `json:"prevent_download"`
	ExternalSongInfo     []any   `json:"external_song_info"`
	SecUID               string  `json:"sec_uid"`
	AvatarThumb          Cover   `json:"avatar_thumb"`
	AvatarMedium         Cover   `json:"avatar_medium"`
	PreviewStartTime     float64 `json:"preview_start_time"`
	PreviewEndTime       float64 `json:"preview_end_time"`
	IsCommerceMusic      bool    `json:"is_commerce_music"`
	IsOriginalSound      bool    `json:"is_original_sound"`
	Artists              any     `json:"artists"`
	LyricShortPosition   any     `json:"lyric_short_position"`
	MuteShare            bool    `json:"mute_share"`
	TagList              any     `json:"tag_list"`
	IsAuthorArtist       bool    `json:"is_author_artist"`
	IsPgc                bool    `json:"is_pgc"`
	SearchHighlight      any     `json:"search_highlight"`
	MultiBitRatePlayInfo any     `json:"multi_bit_rate_play_info"`
	TtToDspSongInfos     any     `json:"tt_to_dsp_song_infos"`
	RecommendStatus      int     `json:"recommend_status"`
	UncertArtists        any     `json:"uncert_artists"`
}

type Aweme struct {
	AwemeID                    string              `json:"aweme_id"`
	Desc                       string              `json:"desc"`
	CreateTime                 int                 `json:"create_time"`
	Author                     Author              `json:"author,omitempty"`
	Music                      Music               `json:"music,omitempty"`
	ChaList                    []ChaList           `json:"cha_list"`
	Video                      Video               `json:"video,omitempty"`
	ShareURL                   string              `json:"share_url"`
	UserDigged                 int                 `json:"user_digged"`
	Statistics                 Statistics          `json:"statistics"`
	Status                     Status              `json:"status"`
	Rate                       int                 `json:"rate"`
	TextExtra                  []TextExtra         `json:"text_extra"`
	IsTop                      int                 `json:"is_top"`
	LabelTop                   Cover               `json:"label_top"`
	ShareInfo                  ShareInfoInner      `json:"share_info,omitempty"`
	Distance                   string              `json:"distance"`
	VideoLabels                []any               `json:"video_labels"`
	IsVr                       bool                `json:"is_vr"`
	IsAds                      bool                `json:"is_ads"`
	AwemeType                  int                 `json:"aweme_type"`
	CmtSwt                     bool                `json:"cmt_swt"`
	ImageInfos                 any                 `json:"image_infos"`
	RiskInfos                  RiskInfos           `json:"risk_infos"`
	IsRelieve                  bool                `json:"is_relieve"`
	SortLabel                  string              `json:"sort_label"`
	Position                   any                 `json:"position"`
	UniqidPosition             any                 `json:"uniqid_position"`
	AuthorUserID               int64               `json:"author_user_id"`
	BodydanceScore             int                 `json:"bodydance_score"`
	Geofencing                 any                 `json:"geofencing"`
	IsHashTag                  int                 `json:"is_hash_tag"`
	IsPgcshow                  bool                `json:"is_pgcshow"`
	Region                     string              `json:"region"`
	VideoText                  []any               `json:"video_text"`
	CollectStat                int                 `json:"collect_stat"`
	LabelTopText               any                 `json:"label_top_text"`
	GroupID                    string              `json:"group_id"`
	PreventDownload            bool                `json:"prevent_download"`
	NicknamePosition           any                 `json:"nickname_position"`
	ChallengePosition          any                 `json:"challenge_position"`
	ItemCommentSettings        int                 `json:"item_comment_settings"`
	WithPromotionalMusic       bool                `json:"with_promotional_music"`
	LongVideo                  any                 `json:"long_video"`
	ItemDuet                   int                 `json:"item_duet"`
	ItemReact                  int                 `json:"item_react"`
	DescLanguage               string              `json:"desc_language"`
	InteractionStickers        any                 `json:"interaction_stickers"`
	MiscInfo                   string              `json:"misc_info"`
	OriginCommentIds           any                 `json:"origin_comment_ids"`
	CommerceConfigData         any                 `json:"commerce_config_data"`
	DistributeType             int                 `json:"distribute_type"`
	VideoControl               VideoControl        `json:"video_control"`
	HasVsEntry                 bool                `json:"has_vs_entry"`
	CommerceInfo               CommerceInfo        `json:"commerce_info"`
	NeedVsEntry                bool                `json:"need_vs_entry"`
	Anchors                    any                 `json:"anchors"`
	HybridLabel                any                 `json:"hybrid_label"`
	WithSurvey                 bool                `json:"with_survey"`
	GeofencingRegions          any                 `json:"geofencing_regions"`
	AwemeACL                   AwemeACL            `json:"aweme_acl"`
	CoverLabels                any                 `json:"cover_labels"`
	MaskInfos                  []any               `json:"mask_infos"`
	SearchHighlight            any                 `json:"search_highlight"`
	PlaylistBlocked            bool                `json:"playlist_blocked"`
	GreenScreenMaterials       any                 `json:"green_screen_materials"`
	InteractPermission         InteractPermission  `json:"interact_permission"`
	QuestionList               any                 `json:"question_list"`
	ContentDesc                string              `json:"content_desc"`
	ContentDescExtra           []ContentDescExtra  `json:"content_desc_extra"`
	ProductsInfo               any                 `json:"products_info"`
	FollowUpPublishFromID      int                 `json:"follow_up_publish_from_id"`
	DisableSearchTrendingBar   bool                `json:"disable_search_trending_bar"`
	ImagePostInfo              ImagePostInfo       `json:"image_post_info,omitempty"`
	MusicBeginTimeInMs         int                 `json:"music_begin_time_in_ms"`
	ItemDistributeSource       string              `json:"item_distribute_source"`
	ItemSourceCategory         int                 `json:"item_source_category"`
	BrandedContentAccounts     any                 `json:"branded_content_accounts"`
	IsDescriptionTranslatable  bool                `json:"is_description_translatable"`
	FollowUpItemIDGroups       string              `json:"follow_up_item_id_groups"`
	IsTextStickerTranslatable  bool                `json:"is_text_sticker_translatable"`
	TextStickerMajorLang       string              `json:"text_sticker_major_lang"`
	OriginalClientText         OriginalClientText  `json:"original_client_text"`
	MusicSelectedFrom          string              `json:"music_selected_from"`
	TtsVoiceIds                any                 `json:"tts_voice_ids"`
	ReferenceTtsVoiceIds       any                 `json:"reference_tts_voice_ids"`
	VoiceFilterIds             any                 `json:"voice_filter_ids"`
	ReferenceVoiceFilterIds    any                 `json:"reference_voice_filter_ids"`
	MusicTitleStyle            int                 `json:"music_title_style"`
	CommentConfig              CommentConfig       `json:"comment_config"`
	AddedSoundMusicInfo        AddedSoundMusicInfo `json:"added_sound_music_info,omitempty"`
	OriginVolume               string              `json:"origin_volume"`
	MusicVolume                string              `json:"music_volume"`
	SupportDanmaku             bool                `json:"support_danmaku"`
	HasDanmaku                 bool                `json:"has_danmaku"`
	MufCommentInfoV2           any                 `json:"muf_comment_info_v2"`
	BehindTheSongMusicIds      any                 `json:"behind_the_song_music_ids"`
	BehindTheSongVideoMusicIds any                 `json:"behind_the_song_video_music_ids"`
	ContentOriginalType        int                 `json:"content_original_type"`
	ShootTabName               string              `json:"shoot_tab_name,omitempty"`
	ContentType                string              `json:"content_type,omitempty"`
	ContentSizeType            int                 `json:"content_size_type,omitempty"`
	OperatorBoostInfo          any                 `json:"operator_boost_info"`
	LogInfo                    LogInfo             `json:"log_info"`
	MainArchCommon             string              `json:"main_arch_common"`
	AigcInfo                   AigcInfo            `json:"aigc_info"`
	Banners                    any                 `json:"banners"`
	PickedUsers                any                 `json:"picked_users"`
	CommentTopbarInfo          any                 `json:"comment_topbar_info"`
}

type TikTokAPIResponse struct {
	StatusCode    int          `json:"status_code"`
	MinCursor     int          `json:"min_cursor"`
	MaxCursor     int          `json:"max_cursor"`
	HasMore       int          `json:"has_more"`
	AwemeList     []Aweme      `json:"aweme_list"`
	HomeModel     int          `json:"home_model"`
	RefreshClear  int          `json:"refresh_clear"`
	Extra         Extra        `json:"extra"`
	LogPb         LogPb        `json:"log_pb"`
	PreloadAds    []any        `json:"preload_ads"`
	PreloadAwemes any          `json:"preload_awemes"`
	LogInfo       LogInfoInner `json:"log_info"`
}
