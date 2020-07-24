package tcp

import (
	"P2P-File-Sharing/common"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

// BUFFERSIZE buffer size for file transfer
const BUFFERSIZE = 1024

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func handleClient(conn net.Conn, dir string) {
	defer conn.Close()

	// getting filename
	var bufferFileName [64]byte
	conn.Read(bufferFileName[:])

	// oepening the file
	file, err := os.Open(dir + strings.TrimRight(string(bufferFileName[:]), ":"))
	common.CheckError(err)

	fileInfo, err := file.Stat()
	common.CheckError(err)

	fileSize := fileInfo.Size()
	fileSizeStr := strconv.FormatInt(fileSize, 10)
	fileSizeStr = fillString(fileSizeStr, 64)

	// sending file size
	_, err = conn.Write([]byte(fileSizeStr))
	common.CheckError(err)

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
	common.CheckError(err)

	listener, err := net.ListenTCP("tcp4", tcpAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleClient(conn, dir)
	}

}
