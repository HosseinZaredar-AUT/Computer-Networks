package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func readClusterList(clusterList *[][2]string, listPath string) {
	f, err := os.Open(listPath)
	checkError(err)

	defer func() {
		f.Close()
		checkError(err)
	}()

	s := bufio.NewScanner(f)
	for s.Scan() {
		fields := strings.Fields(s.Text())
		*clusterList = append(*clusterList, [2]string{fields[0], fields[1]})
	}

	err = s.Err()
	checkError(err)
}

func main() {
	listPath := flag.String("l", "", "cluster list file path")
	// dirPath := flag.String("d", "", "directory path")
	flag.Parse()

	clusterList := make([][2]string, 0, 10)
	readClusterList(&clusterList, *listPath)
	fmt.Println(clusterList)
}
