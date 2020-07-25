package udp

import (
	"P2P-File-Sharing/common"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// opens the input directory and returns the number of files existing in it
func countNumOfFiles(dir string) int {

	// openning the directory
	f, err := os.Open(dir)
	common.CheckError(err)

	// reading the contents
	files, err := f.Readdir(-1)
	common.CheckError(err)

	// closing the directory
	err = f.Close()
	common.CheckError(err)

	count := 0

	// going through all entires of the directory in search of files
	for _, file := range files {
		if !file.IsDir() && file.Name() != "temp" {
			count++
		}
	}

	return count
}

// gets the cluster map as input and returns it as a comma-seperated string of values
func flattenList(clusterMap map[string]string) string {
	flatList := ""
	flag := false
	for key, value := range clusterMap {
		if flag {
			flatList = flatList + ","
		}
		flatList = flatList + key + " " + value
		flag = true
	}
	return flatList
}

//DiscoverService on certain intervals, sends discovery messages to all the nodes we are aware of
func DiscoverService(clusterMap map[string]string, myNode common.Node, cmMutex *sync.Mutex, dir string, discoveryInterval int) {

	// forever
	for {

		cmMutex.Lock() // lock cluster map

		// updating numOfFiles (= the number of files we're serving)
		numOfFiles := countNumOfFiles(dir)
		clusterMap[myNode.Name] = myNode.IP + ":" + myNode.UDPPPort + ";" + strconv.Itoa(numOfFiles)

		// getting a copy from cluster map
		clusterMapCopy := make(map[string]string)
		for key, value := range clusterMap {
			clusterMapCopy[key] = value
		}

		cmMutex.Unlock() // unlock cluster map

		// turn cluster map into an string
		flatList := flattenList(clusterMapCopy)

		// for each node in cluster map
		for _, info := range clusterMapCopy {

			// find the address of that node
			addr := strings.Split(info, ";")[0]

			// not sending discovery message to myself
			if addr == (myNode.IP + ":" + myNode.UDPPPort) {
				continue
			}

			// creating proper address
			udpAddr, err := net.ResolveUDPAddr("udp4", addr)
			common.CheckError(err)

			// connecting to node's UDP server
			conn, err := net.DialUDP("udp", nil, udpAddr)
			common.CheckError(err)

			// sending the discovery message
			_, err = conn.Write([]byte("dis:" + flatList))
			common.CheckError(err)

		}

		// have some rest!
		time.Sleep(time.Duration(discoveryInterval) * 1000000000)

	}
}
