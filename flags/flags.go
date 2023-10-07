package flags

import "flag"

var Domain *string
var Public *bool

func init() {
	Domain = flag.String("domain", "", "Specify domain where the collages are available")
	Public = flag.Bool("public", false, "Specify whether the collages are publicly available")
	flag.Parse()
}
