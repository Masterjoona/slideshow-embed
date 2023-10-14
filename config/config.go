package config

import (
	"os"
)

var Domain string
var Public bool
var Port string

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
