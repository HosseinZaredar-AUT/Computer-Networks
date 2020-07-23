package udp

import (
	"P2P-File-Sharing/common"
	"net"
	"strconv"
	"strings"
	"time"
)

func requestNode(fileName string, nodeAddress string, ch chan [2]string) {
	udpAddr, err := net.ResolveUDPAddr("udp4", nodeAddress)
	checkError(err)

	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)

	_, err = conn.Write([]byte("req:" + fileName))
	checkError(err)

	// waiting for response

	var buffer [512]byte

	conn.SetReadDeadline(time.Now().Add(time.Second))
	l, _ := conn.Read(buffer[:])

	if l != 0 { // if we got any response
		response := strings.Split(string(buffer[:]), ",")
		sendTime, err := strconv.ParseInt(response[0], 10, 64)
		checkError(err)

		delay := time.Now().UnixNano() - sendTime

		ch <- [2]string{strconv.FormatInt(delay, 10), response[1]}
	}
}

// FileRequest ...
func FileRequest(fileName string, clusterMap map[string]string, myNode common.Node) string {

	// TODO check if you already got the file
	// TODO check if you a different file but with the same name

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
	bestNode := ""
	var bestDelay int64
	bestDelay = 9223372036854775807 // maximum possible value

	for n := range ch {
		// fmt.Printf("Consumed %s\n", n)
		delay, err := strconv.ParseInt(n[0], 10, 64)
		checkError(err)

		if delay < bestDelay {
			bestDelay = delay
			bestNode = n[1]
		}
	}

	if bestNode == "" { // if we didn't find the file
		return "!"
	}

	// if the file was found
	return bestNode
}
