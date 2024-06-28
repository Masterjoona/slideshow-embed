package main

var FetchTiktokData func(string) (SimplifiedData, error)

func SetTiktokTiktokProvider(provider string) {
	if provider == "tikwm" {
		FetchTiktokData = FetchTiktokDataTikWm
		return
	} else if provider == "ttsave" {
		FetchTiktokData = FetchTiktokDataTTSave
		return
	}
	FetchTiktokData = FetchTiktokDataTiktokAPI
}
