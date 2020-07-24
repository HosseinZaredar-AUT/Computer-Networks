package common

import (
	"fmt"
	"os"
	"runtime/debug"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		fmt.Println()
		debug.PrintStack()
		os.Exit(-1)
	}
}

func CountNumOfFiles(dir string) int {
	f, err := os.Open(dir)
	CheckError(err)

	files, err := f.Readdir(-1)
	CheckError(err)

	err = f.Close()
	CheckError(err)

	count := 0

	for _, file := range files {
		if !file.IsDir() {
			count++
		}
	}

	return count
}
