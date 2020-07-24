package udp

import (
	"P2P-File-Sharing/common"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

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

//DiscoverService ...
func DiscoverService(clusterMap map[string]string, myNode common.Node, cmMutex *sync.Mutex, dir string) {
	for {

		cmMutex.Lock()

		// updating numOfFiles
		numOfFiles := common.CountNumOfFiles(dir)
		clusterMap[myNode.Name] = myNode.IP + ":" + myNode.UDPPPort + ";" + strconv.Itoa(numOfFiles)

		// getting a copy from cluster map
		clusterMapCopy := make(map[string]string)
		for key, value := range clusterMap {
			clusterMapCopy[key] = value
		}
		cmMutex.Unlock()

		// turn cluster map into an string
		flatList := flattenList(clusterMapCopy)

		// for each node in cluster map
		for _, info := range clusterMapCopy {

			addr := strings.Split(info, ";")[0]
			// no sending discovery message to myself
			if addr == (myNode.IP + ":" + myNode.UDPPPort) {
				continue
			}

			udpAddr, err := net.ResolveUDPAddr("udp4", addr)
			common.CheckError(err)

			conn, err := net.DialUDP("udp", nil, udpAddr)
			common.CheckError(err)

			_, err = conn.Write([]byte("dis:" + flatList))
			common.CheckError(err)

			// fmt.Println("Sent the cluster list")
		}

		// have some rest!
		time.Sleep(4 * time.Second)

	}
}
