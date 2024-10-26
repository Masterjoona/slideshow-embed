package types

import (
	"meow/pkg/media"
	"meow/pkg/net"
	"meow/pkg/util"
)

func (t *TiktokInfo) DecodeStrings() {
	t.Author = util.UrlDecodeString(t.Author)
	t.Caption = util.UrlDecodeString(t.Caption)
	t.Video.Url = util.UrlDecodeString(t.Video.Url)
	t.SoundLink = util.UrlDecodeString(t.SoundLink)
	t.ImageLinks = util.UrlDecodeStrings(t.ImageLinks)

}

func (t *TiktokInfo) MakeCollage() error {
	return media.MakeCollage(t.VideoID, &t.ImageBuffers)
}

func (t *TiktokInfo) MakeCollageWithAudio(filetype string) (string, string, error) {
	return media.MakeCollageWithAudio(t.VideoID, t.SoundBuffer, filetype)
}

func (t *TiktokInfo) MakeVideoSlideshow() (string, string, error) {
	return media.MakeVideoSlideshow(t.VideoID, t.FileName, t.SoundBuffer, &t.ImageBuffers)
}

func (t *TiktokInfo) MakeVideoSubtitles(lang string) (string, string, error) {
	return media.MakeVideoSubtitles(t.VideoID, t.FileName, lang)
}

func (t *TiktokInfo) DownloadVideo() error {
	buffer, err := net.DownloadMedia(t.Video.Url)
	if err != nil {
		return err
	}
	t.Video.Buffer = buffer
	return nil
}
