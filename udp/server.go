package udp

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

//Server ...
func Server(clusterMap map[string]string, myAddress string) {

	udpAddr, err := net.ResolveUDPAddr("udp4", myAddress)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	fmt.Println("UDP server listining on", udpAddr)

	var buffer [512]byte

	for {
		_, _, err := conn.ReadFromUDP(buffer[:])
		if err != nil {
			continue
		}

		// updating cluster map
		message := string(buffer[:])
		nodes := strings.Split(message, ",")
		for _, node := range nodes {
			fields := strings.Fields(node)
			clusterMap[fields[0]] = fields[1]
		}

		fmt.Println("cluster map:", clusterMap)

	}
}
