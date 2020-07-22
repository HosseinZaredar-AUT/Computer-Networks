package main

import (
	"P2P-File-Sharing/udp"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// this function reads the cluster list file, updates clusterList slice
// and returns the address of the machine
func readClusterList(clusterList *[][2]string, listPath string) string {
	f, err := os.Open(listPath)
	checkError(err)

	defer func() {
		f.Close()
		checkError(err)
	}()

	s := bufio.NewScanner(f)

	// getting clients own address from the first line
	s.Scan()
	address := strings.Fields(s.Text())[1]

	// getting other nodes
	for s.Scan() {
		fields := strings.Fields(s.Text())
		*clusterList = append(*clusterList, [2]string{fields[0], fields[1]})
	}

	err = s.Err()
	checkError(err)

	return address
}

func main() {

	// parse flags
	listPath := flag.String("l", "", "cluster list file path")
	// dirPath := flag.String("d", "", "directory path")
	flag.Parse()

	// read cluster list from file
	clusterList := make([][2]string, 0, 10)
	address := readClusterList(&clusterList, *listPath)
	fmt.Println("initial cluster list:", clusterList)

	// run udp server
	go udp.Server(&clusterList, address)

	// waiting for goroutines
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
