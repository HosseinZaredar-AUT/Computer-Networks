package common

import (
	"fmt"
	"os"
	"runtime/debug"
)

// CheckError ...
func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		fmt.Println()
		debug.PrintStack()
		os.Exit(-1)
	}
}
