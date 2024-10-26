package net

import "strings"

func validateURL(url string) bool {
	if url == "" {
		return false
	}
	if !strings.Contains(url, ".tiktxk.com") && !strings.Contains(url, ".tiktok.com") &&
		!strings.Contains(url, ".vxtiktok.com") {
		return false
	}
	return true
}
