package fs

import "os"

func ReadFile(path string) []byte {
	bb, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bb
}
