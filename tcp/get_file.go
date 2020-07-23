package tcp

import (
	"fmt"
	"net"
)

// GetFile ...
func GetFile(addr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	var buffer [512]byte
	conn.Read(buffer[:])

	fmt.Println(string(buffer[:]))
}
