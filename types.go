package main

type SimplifiedVideo struct {
	Url    string
	Width  string
	Height string
}

type SimplifiedData struct {
	Author       string
	Caption      string
	VideoID      string
	Details      Counts
	ImageLinks   []string
	ImageBuffers [][]byte
	SoundLink    string
	SoundBuffer  []byte
	Video        SimplifiedVideo
	FileName     string
}

type Counts struct {
	Likes     string
	Comments  string
	Shares    string
	Views     string
	Favorites string
}

type ImageWithIndex struct {
	Bytes []byte
	Index int
}

type FileLink struct {
	Name string
	Path string
}

type Stats struct {
	FilePaths []FileLink
	FileCount string
	TotalSize string
}

type ShortLinkInfo struct {
	VideoId      string
	UniqueUserId string
}
