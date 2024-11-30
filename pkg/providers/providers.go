package providers

import (
	"fmt"
	"meow/pkg/config"
	"meow/pkg/net"
	playerapi "meow/pkg/providers/player_api"
	"meow/pkg/providers/tiktok_api"
	"meow/pkg/providers/tikwm"
	"meow/pkg/providers/ttsave"
	"meow/pkg/types"
)

var RecentTiktokReqs = net.NewCache[string, types.TiktokInfo]()

var fetchers = map[string]func(string) (types.TiktokInfo, error){
	"tikwm":     tikwm.FetchTiktok,
	"tiktok":    tiktok_api.FetchTikok,
	"playerapi": playerapi.FetchTikok,
	"ttsave":    ttsave.FetchTiktok,
}

var Fetchers = map[string]func(string) (types.TiktokInfo, error){}
var defaultFetcher func(string) (types.TiktokInfo, error)

func MakeMap() {
	if fetcher := config.TiktokProvider; fetcher != "" {
		defaultFetcher = fetchers[fetcher]
	} else {
		defaultFetcher = fetchers["tikwm"]
	}

	if !config.FallbackProvider {
		Fetchers = map[string]func(string) (types.TiktokInfo, error){
			config.TiktokProvider: defaultFetcher,
		}
		return
	} else {
		Fetchers = fetchers
	}

	if defaultFetcher != nil {
		delete(Fetchers, config.TiktokProvider)
	}
}

func fetchAndCache(videoId string, fetcher func(string) (types.TiktokInfo, error)) (types.TiktokInfo, error) {
	data, err := fetcher(videoId)
	if err != nil {
		return types.TiktokInfo{}, err
	}

	data.DecodeStrings()
	if data.Video.Url != "" {
		data.FileName = fmt.Sprintf("%s.mp4", videoId)
		if err := data.DownloadVideo(); err != nil {
			fmt.Println("Failed to download video", err)
			return types.TiktokInfo{}, err
		}
	}

	RecentTiktokReqs.Set(videoId, data)
	return data, nil
}

func FetchTiktok(videoId string) (types.TiktokInfo, error) {
	if data, ok := RecentTiktokReqs.Get(videoId); ok {
		return data, nil
	}

	defaultResponse, err := fetchAndCache(videoId, defaultFetcher)
	if err != nil {
		for _, fetcher := range Fetchers {
			response, err := fetchAndCache(videoId, fetcher)
			if err == nil {
				return response, nil
			}
		}
	}
	return defaultResponse, nil
}
