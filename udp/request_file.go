package udp

import (
	"P2P-File-Sharing/common"
	"net"
	"strconv"
	"strings"
	"time"
)

// sends request to a single peer asking for the file
func requestNode(fileName string, nodeAddress string, ch chan [2]string, timeOut int) {

	// creating proper address
	udpAddr, err := net.ResolveUDPAddr("udp4", nodeAddress)
	common.CheckError(err)

	// connecting to node's UDP server
	conn, err := net.DialUDP("udp", nil, udpAddr)
	common.CheckError(err)

	// sending the request
	_, err = conn.Write([]byte("req:" + fileName))
	common.CheckError(err)

	// waiting for some time to see if we get any response

	conn.SetReadDeadline(time.Now().Add(time.Duration(timeOut * 500000000)))

	var buffer [150]byte
	l, _ := conn.Read(buffer[:])

	if l != 0 { // if we got any response

		response := strings.Split(strings.TrimRight(string(buffer[:]), "\x00"), ",")

		// if the server was busy
		if string(buffer[0:4]) == "busy" {
			// pushing proper info into channel
			ch <- [2]string{"busy", response[1]}

		} else { // if the server was

			// extracting the time the message was sent by the server
			sendTime, err := strconv.ParseInt(response[0], 10, 64)
			common.CheckError(err)

			// calculating the delay
			// delay = (now) - (the time the message was sent by the server)
			delay := time.Now().UnixNano() - sendTime

			// pushing proper info into channel
			ch <- [2]string{strconv.FormatInt(delay, 10), response[1]}
		}
	}
}

// FileRequest sends request for the input file to all peers in cluster map
func FileRequest(fileName string, clusterMap map[string]string, myNode common.Node, timeOut int) string {

	// a shared channel is used to contain response of each peer
	// in form of: {delay, server's IP:TCPPort}
	ch := make(chan [2]string, 10)

	// for each peer in cluster map
	for _, info := range clusterMap {

		// find the address of that node
		addr := strings.Split(info, ";")[0]

		// not sending message to myself
		if addr == (myNode.IP + ":" + myNode.UDPPPort) {
			continue
		}

		// send request to that node in a seperate goroutine
		go requestNode(fileName, addr, ch, timeOut)

	}

	// wait for some time (= timeout)
	time.Sleep(time.Duration(timeOut * 1000000000))

	// close the channel
	close(ch)

	// choose the best response (if there's any)

	allBusy := true // is used to track if all the peers having the file are busy
	exists := false // is used to check if the file was found

	bestNode := ""
	var bestDelay int64
	bestDelay = 9223372036854775807 // maximum possible value

	// pop all elemets of the channel
	for n := range ch {
		// fmt.Printf("Consumed %s\n", n)

		exists = true

		// check if the node was busy
		if n[0] == "busy" {
			continue
		} else { // if the node is not busy
			allBusy = false

			// getting the delay
			delay, err := strconv.ParseInt(n[0], 10, 64)
			common.CheckError(err)

			// if it has better delal
			if delay < bestDelay {
				bestDelay = delay
				bestNode = n[1]
			}
		}
	}

	if !exists { // if we didn't find the file at all
		return "not found"
	} else if allBusy { // if the file was found but the node(s) having the file was(were) busy
		return "busy"
	}

	// if the file was found and at least one node having that file isn't busy
	return bestNode
}
