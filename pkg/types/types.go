package types

type SimplifiedVideo struct {
	Buffer []byte
	Url    string
	Width  string
	Height string
}

type TiktokInfo struct {
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
