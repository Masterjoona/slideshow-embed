package providers

import (
	"fmt"
	"meow/pkg/net"
	"meow/pkg/providers/tiktok_api"
	"meow/pkg/providers/tikwm"
	"meow/pkg/providers/ttsave"
	"meow/pkg/types"
)

var RecentTiktokReqs = net.NewCache[types.TiktokInfo](20)

var Fetchers = map[string]func(string) (types.TiktokInfo, error){
	"tikwm":  tikwm.FetchTiktok,
	"tiktok": tiktok_api.FetchTikok,
	"ttsave": ttsave.FetchTiktok,
}

func FetchTiktokData(videoId string) (types.TiktokInfo, error) {
	if data, ok := RecentTiktokReqs.Get(videoId); ok {
		return data, nil
	}
	for _, fetcher := range Fetchers {
		if data, err := fetcher(videoId); err == nil {
			data.DecodeStrings()
			if data.Video.Url != "" {
				data.FileName = fmt.Sprintf("%s.mp4", videoId)
				err := data.DownloadVideo()
				if err != nil {
					fmt.Println("Failed to download video")
				}
			}
			RecentTiktokReqs.Put(videoId, data)
			return data, nil
		}
	}
	return types.TiktokInfo{}, fmt.Errorf("failed to fetch data for video %s", videoId)
}
