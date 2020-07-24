package udp

import (
	"P2P-File-Sharing/common"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func requestNode(fileName string, nodeAddress string, ch chan [2]string) {
	udpAddr, err := net.ResolveUDPAddr("udp4", nodeAddress)
	common.CheckError(err)

	conn, err := net.DialUDP("udp", nil, udpAddr)
	common.CheckError(err)

	_, err = conn.Write([]byte("req:" + fileName))
	common.CheckError(err)

	// waiting for response

	var buffer [512]byte

	conn.SetReadDeadline(time.Now().Add(time.Second))
	l, _ := conn.Read(buffer[:])

	if l != 0 { // if we got any response

		response := strings.Split(string(buffer[:]), ",")

		// check if server is busy
		if string(buffer[0:4]) == "busy" {
			ch <- [2]string{"busy", response[1]}
		} else {
			sendTime, err := strconv.ParseInt(response[0], 10, 64)
			common.CheckError(err)

			delay := time.Now().UnixNano() - sendTime

			ch <- [2]string{strconv.FormatInt(delay, 10), response[1]}
		}
	}
}

// FileRequest ...
func FileRequest(fileName string, clusterMap map[string]string, myNode common.Node) string {

	// shared channel
	ch := make(chan [2]string, 10)

	for _, addr := range clusterMap {

		// no sending message to myself
		if addr == (myNode.IP + ":" + myNode.UDPPPort) {
			continue
		}

		go requestNode(fileName, addr, ch)

	}

	// wait for some time
	time.Sleep(2 * time.Second)

	// close the channel
	close(ch)

	// choose the best response
	allBusy := true
	bestNode := ""
	var bestDelay int64
	bestDelay = 9223372036854775807 // maximum possible value

	for n := range ch {
		fmt.Printf("Consumed %s\n", n)
		// check if node was busy
		if n[0] == "busy" {
			continue
		} else {
			allBusy = false
			delay, err := strconv.ParseInt(n[0], 10, 64)
			common.CheckError(err)

			if delay < bestDelay {
				bestDelay = delay
				bestNode = n[1]
			}
		}
	}

	if bestNode == "" && !allBusy { // if we didn't find the file
		return "not found"
	} else if bestNode == "" && allBusy { // if the file was found but the node(s) having that were busy
		return "busy"
	}

	// if the file was found and at least one node having that file isn't busy
	return bestNode
}
