package config

import "os"

func getEnv(k, defaultVal string) string {
	v := os.Getenv(k)
	if v != "" {
		return v
	}
	return defaultVal
}
