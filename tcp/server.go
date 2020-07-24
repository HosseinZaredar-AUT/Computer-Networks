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

func handleClient(clusterMap map[string]string, conn net.Conn, dir string, numServing *int, averageNumFiles *float64) {
	defer func() {
		conn.Close()
		(*numServing)--
	}()

	// getting info = filename + the name of requesting node
	var bufferInfo [100]byte
	conn.Read(bufferInfo[:])

	info := strings.TrimRight(string(bufferInfo[:]), ":")
	fields := strings.Split(info, ":")

	fileName := fields[0]
	nodeName := fields[1]

	// oepening the file
	file, err := os.Open(dir + fileName)
	common.CheckError(err)

	fileInfo, err := file.Stat()
	common.CheckError(err)

	fileSize := fileInfo.Size()
	fileSizeStr := strconv.FormatInt(fileSize, 10)
	fileSizeStr = fillString(fileSizeStr, 64, ":")

	// sending file size
	_, err = conn.Write([]byte(fillString(fileSizeStr, 64, ":")))
	common.CheckError(err)

	// checking if the node is speed limited (free-riding node detection)
	nodeInfo := clusterMap[nodeName]
	numOfFiles, err := strconv.Atoi(strings.Split(nodeInfo, ";")[1])
	common.CheckError(err)

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
	for {
		_, err = file.Read(sendBuffer[:])
		if err == io.EOF {
			break
		}
		conn.Write(sendBuffer[:])

		// applying speed limit
		if isSpeedLimited {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// Server ...
func Server(clusterMap map[string]string, myNode common.Node, dir string, numServing *int, averageNumFiles *float64) {

	service := myNode.IP + ":" + myNode.TCPPort
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	common.CheckError(err)

	listener, err := net.ListenTCP("tcp4", tcpAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		// adding 1 to number of clients being served
		(*numServing)++
		go handleClient(clusterMap, conn, dir, numServing, averageNumFiles)
	}

}
