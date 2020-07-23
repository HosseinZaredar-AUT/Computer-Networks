package tcp

import (
	"P2P-File-Sharing/common"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

// BUFFERSIZE buffer size for file transfer
const BUFFERSIZE = 1024

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func handleClient(conn net.Conn, dir string) {
	defer conn.Close()

	var bufferFileName [64]byte
	conn.Read(bufferFileName[:])

	// oepening the file
	file, err := os.Open(dir + strings.TrimRight(string(bufferFileName[:]), "\x00"))
	checkError(err)

	fileInfo, err := file.Stat()
	checkError(err)

	fileSize := fileInfo.Size()

	// sending file size
	_, err = conn.Write([]byte(strconv.FormatInt(fileSize, 10)))
	checkError(err)

	// os.Exit(0)

	// sending the file
	var sendBuffer [BUFFERSIZE]byte
	for {
		_, err = file.Read(sendBuffer[:])
		if err == io.EOF {
			break
		}
		conn.Write(sendBuffer[:])
	}
}

// Server ...
func Server(myNode common.Node, dir string) {

	service := myNode.IP + ":" + myNode.TCPPort
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp4", tcpAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleClient(conn, dir)
	}

}
