package config

import (
	"os"
)

var Domain string
var Public bool
var Port string
var UserAgent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"

func init() {
	Domain = os.Getenv("DOMAIN")
	Public = os.Getenv("PUBLIC") == "true"
	if Port = os.Getenv("PORT"); Port == "" {
		Port = "4232"
	}
	if Port[0] != ':' {
		Port = ":" + Port
	}
}
