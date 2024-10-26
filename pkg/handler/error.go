package handler

var ErrorImages = []string{
	"https://media.discordapp.net/attachments/961445186280509451/980132677338423316/fuckmedaddyharderohyeailovecokcimsocissyfemboy.gif",
	"https://media.discordapp.net/attachments/901959319719936041/996927812927750264/chrome_2WOKI6Jm3v.gif",
	"https://cdn.discordapp.com/attachments/749030987530502197/980338691706880051/79587A35-FD36-41D3-8232-7A29B46D2543.gif",
	"https://media.discordapp.net/attachments/880335303984943154/1237439972290859140/looksinside.gif",
}

var ErrorImagesIndex = 0

func ErrorImage() string {
	if ErrorImagesIndex == 3 {
		ErrorImagesIndex = 0
	} else {
		ErrorImagesIndex++
	}
	return ErrorImages[ErrorImagesIndex]
}
