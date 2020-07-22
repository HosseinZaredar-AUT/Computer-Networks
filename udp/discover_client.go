package udp

import (
	"fmt"
	"net"
	"time"
)

//DiscoverClient ...
func DiscoverClient(clusterList *[][2]string) {

	for {

		// for each node in cluster list
		for i, node := range *clusterList {
			fmt.Println(i, node)

			service := node[1]
			udpAddr, err := net.ResolveUDPAddr("udp4", service)
			checkError(err)

			conn, err := net.DialUDP("udp", nil, udpAddr)
			checkError(err)

			_, err = conn.Write([]byte("Hello UDP server!"))
			checkError(err)

			// var buffer [512]byte
			// conn.Read(buffer[:])

			fmt.Println("Sent some data")
		}

		// have some rest!
		time.Sleep(time.Second)
	}
}
