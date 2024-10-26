package net

import "container/list"

type Cache[T any] struct {
	capacity int
	items    map[string]*list.Element
	order    *list.List
}

type Item[T any] struct {
	key   string
	value T
}

type ImageWithIndex struct {
	Bytes []byte
	Index int
}

type ShortLinkInfo struct {
	VideoId      string
	UniqueUserId string
}

type SubtitlesResp struct {
	ServerTransTime      int    `json:"server_trans_time"`
	StatusCode           int    `json:"status_code"`
	StatusMsg            string `json:"status_msg"`
	TranslationCacheTime int    `json:"translation_cache_time"`
	Variant              string `json:"variant"`
	WebvttSubtitle       string `json:"webvtt_subtitle"`
}
