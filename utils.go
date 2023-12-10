package main

import "strings"

var errorImages = []string{
	"https://media.discordapp.net/attachments/961445186280509451/980132677338423316/fuckmedaddyharderohyeailovecokcimsocissyfemboy.gif",
	"https://media.discordapp.net/attachments/901959319719936041/996927812927750264/chrome_2WOKI6Jm3v.gif",
	"https://cdn.discordapp.com/attachments/749030987530502197/980338691706880051/79587A35-FD36-41D3-8232-7A29B46D2543.gif",
}
var errorImagesIndex = 0

func errorImage() string {
	if errorImagesIndex == 2 {
		errorImagesIndex = 0
	} else {
		errorImagesIndex++
	}
	return errorImages[errorImagesIndex]
}

func validateURL(url string) bool {
	if url == "" {
		return false
	}
	if !strings.Contains(url, "vm.tiktxk.com") && !strings.Contains(url, "vm.tiktok.com") {
		return false
	}
	return true
}
