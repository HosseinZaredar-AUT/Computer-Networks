package common

import (
	"fmt"
	"os"
	"runtime/debug"
)

// CheckError check the err input and in case of any error
// prints proper info and stops the whole program
func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		fmt.Println()
		debug.PrintStack()
		os.Exit(-1)
	}
}
