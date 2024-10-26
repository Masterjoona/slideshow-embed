package config

import (
	"os"
	"strings"
)

func addString(s string, r string, trailing bool) string {
	if (trailing && !strings.HasSuffix(s, r)) || (!trailing && !strings.HasPrefix(s, r)) {
		if trailing {
			return s + r
		}
		return r + s
	}
	return s
}

func checkEnvOrDefault(env string, def string) string {
	if val := os.Getenv(env); val != "" {
		return val
	}
	return def
}
