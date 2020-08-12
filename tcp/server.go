package tcp

import (
	"P2P-File-Sharing/common"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// BUFFERSIZE buffer size for file transfer
const BUFFERSIZE = 1024

// gets a string and adds 'filler' characters to the end of the string until it reaches to the wanted length
// this is needed because in TCP, the message's length must equal to the length of recievers buffer
func fillString(retunString string, toLength int, filler string) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + filler
			continue
		}
		break
	}
	return retunString
}

// responds to a client requesting a file transmit
func handleClient(clusterMap map[string]string, conn net.Conn, dir string, numServing *int, averageNumFiles *float64) {

	// at the end, close the connetion and decreament the number of clients being served
	defer func() {
		conn.Close()
		(*numServing)--
	}()

	// getting info = filename + the name of requesting node
	var bufferInfo [100]byte
	conn.Read(bufferInfo[:])

	info := strings.TrimRight(string(bufferInfo[:]), ":")
	fields := strings.Split(info, ":")

	// extracting client's node name and requested file name
	fileName := fields[0]
	nodeName := fields[1]

	// oepening the file
	file, err := os.Open(dir + fileName)
	common.CheckError(err)

	// getting the file's info
	fileInfo, err := file.Stat()
	common.CheckError(err)

	// finding the size of the file
	fileSize := fileInfo.Size()
	fileSizeStr := strconv.FormatInt(fileSize, 10)
	fileSizeStr = fillString(fileSizeStr, 64, ":")

	// sending file size to client
	_, err = conn.Write([]byte(fillString(fileSizeStr, 64, ":")))
	common.CheckError(err)

	// checking if the node is speed limited (free-riding node detection)
	nodeInfo := clusterMap[nodeName]
	numOfFiles, err := strconv.Atoi(strings.Split(nodeInfo, ";")[1])
	common.CheckError(err)

	// if the number of files shared by the client is less that the average number of files
	// shared each peer, then the user is trying to get a free ride! so must be speed limited
	isSpeedLimited := float64(numOfFiles) < *averageNumFiles

	// letting the client know if it is speed limited
	var limitStr string
	if limitStr = "0"; isSpeedLimited {
		limitStr = "1"
	}
	_, err = conn.Write([]byte(limitStr))
	common.CheckError(err)

	// sending the file
	var sendBuffer [BUFFERSIZE]byte

	// until we've sent all of the file
	for {
		// reading a chunk from file
		_, err = file.Read(sendBuffer[:])
		if err == io.EOF {
			break
		}

		// sending the chunk of file
		conn.Write(sendBuffer[:])

		// applying speed limit if the client is speed limited
		if isSpeedLimited {
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// Server responsble for sending files to peers
func Server(clusterMap map[string]string, myNode common.Node, dir string, numServing *int, averageNumFiles *float64) {

	// creating proper address
	service := myNode.LocalIP + ":" + myNode.TCPPort
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	common.CheckError(err)

	// listening...
	listener, err := net.ListenTCP("tcp4", tcpAddr)

	// forever
	for {

		// wait for incoming message
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		// adding 1 to number of clients being served
		(*numServing)++

		// responding to the client in a seperate goroutine
		go handleClient(clusterMap, conn, dir, numServing, averageNumFiles)
	}

}
