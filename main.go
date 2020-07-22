package main

import (
	"P2P-File-Sharing/udp"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// this function reads the cluster nodes file, updates clusterMap
// and returns the address of this machine
func readclusterMap(clusterMap map[string]string, listPath string) string {
	f, err := os.Open(listPath)
	checkError(err)

	defer func() {
		err := f.Close()
		checkError(err)
	}()

	s := bufio.NewScanner(f)

	// getting my own address from the first line
	s.Scan()
	fields := strings.Fields(s.Text())
	myAddress := fields[1]
	clusterMap[fields[0]] = fields[1]

	// getting other nodes
	for s.Scan() {
		fields := strings.Fields(s.Text())
		clusterMap[fields[0]] = fields[1]
	}

	err = s.Err()
	checkError(err)

	return myAddress
}

func main() {

	// parse flags
	listPath := flag.String("l", "", "cluster list file path")
	dir := flag.String("d", "", "directory path")
	flag.Parse()

	// read list of cluster nodes from file
	clusterMap := make(map[string]string) // a map from name to IP address
	myAddress := readclusterMap(clusterMap, *listPath)
	fmt.Println("initial cluster map:", clusterMap)

	// run udp server
	go udp.Server(clusterMap, myAddress, *dir)

	// run discover client
	go udp.DiscoverService(clusterMap, myAddress)

	go func() {
		for {
			fmt.Println("sent file request!")
			udp.FileRequest("a.txt", clusterMap, myAddress)
			time.Sleep(2 * time.Second)
		}
	}()

	// waiting for goroutines
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
