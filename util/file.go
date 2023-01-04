package util

import (
	"io/fs"
	"os"
)

func FileExists(output string) int {
	_, err := os.Stat(output)
	if err == nil {
		//fmt.Println("File exist")
		return 1
	}
	if os.IsNotExist(err) {
		//fmt.Println("File not exist")
		return -1
	}
	return 0
}

func CheckAndMkdir(output string) {
	if FileExists(output) == -1 {
		_ = os.MkdirAll(output, fs.ModePerm)
	}
}
