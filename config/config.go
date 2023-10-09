package config

import (
	"os"
)

var Domain string
var Public bool

func init() {
	Domain = os.Getenv("DOMAIN")
	Public = os.Getenv("PUBLIC") == "true"
}
