package udp

import (
	"fmt"
	"net"
	"time"
)

// FileRequest ...
func FileRequest(fileName string, clusterMap map[string]string, myAddress string) {
	for _, addr := range clusterMap {

		// TODO check if you already got the file
		// TODO check if you a different file but with the same name

		// no sending message to myself
		if addr == myAddress {
			continue
		}

		udpAddr, err := net.ResolveUDPAddr("udp4", addr)
		checkError(err)

		conn, err := net.DialUDP("udp", nil, udpAddr)
		checkError(err)

		_, err = conn.Write([]byte("req:" + fileName))
		checkError(err)

		// waiting for response

		var buffer [512]byte

		conn.SetReadDeadline(time.Now().Add(time.Second))
		conn.Read(buffer[:])

		fmt.Println(string(buffer[:]))
	}
}