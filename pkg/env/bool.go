package env

import "os"

func Bool(key string) bool {
	val := os.Getenv(key)
	return val == "1" || val == "true"
}
