package types

import "fmt"

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

func (c *Counts) ToString() string {
	return fmt.Sprintf("❤️ %s | 💬 %s | 🔁 %s | ⭐ %s | 👀 %s", c.Likes, c.Comments, c.Shares, c.Favorites, c.Views)
}
