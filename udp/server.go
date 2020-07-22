package udp

import (
	"fmt"
	"net"
	"os"
	"time"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func handleClient(conn *net.UDPConn, clientAddr *net.UDPAddr, in [512]byte) {
	fmt.Printf("got '%s' from %s\n", string(in[:]), clientAddr)
	time.Sleep(3 * time.Second)
	conn.WriteToUDP([]byte("This is a UDP message!"), clientAddr)
}

//Server ...
func Server(clusterList *[][2]string, address string) {
	udpAddr, err := net.ResolveUDPAddr("udp4", address)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	fmt.Println("UDP server listining on", udpAddr)

	var buffer [512]byte

	for {
		_, clientAddr, err := conn.ReadFromUDP(buffer[:])
		if err != nil {
			continue
		}
		go handleClient(conn, clientAddr, buffer)
	}
}
