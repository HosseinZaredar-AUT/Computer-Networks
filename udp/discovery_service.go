package udp

import (
	"P2P-File-Sharing/common"
	"net"
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
func DiscoverService(clusterMap map[string]string, myNode common.Node) {
	for {

		// turn cluster map into an string
		flatList := flattenList(clusterMap)

		// for each node in cluster map
		for _, addr := range clusterMap {

			// no sending discovery message to myself
			if addr == (myNode.IP + ":" + myNode.UDPPPort) {
				continue
			}

			udpAddr, err := net.ResolveUDPAddr("udp4", addr)
			checkError(err)

			conn, err := net.DialUDP("udp", nil, udpAddr)
			checkError(err)

			_, err = conn.Write([]byte("dis:" + flatList))
			checkError(err)

			// fmt.Println("Sent the cluster list")
		}

		// have some rest!
		time.Sleep(4 * time.Second)

	}
}
