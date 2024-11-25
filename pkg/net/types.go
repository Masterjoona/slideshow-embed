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
