package tcp

import (
	"P2P-File-Sharing/common"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

// GetFile ...
func GetFile(fileName string, serverName string, addr string, dir string, myNode common.Node) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	common.CheckError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	common.CheckError(err)

	// sending file name + myNode.name
	fileNameFilled := fillString(fileName+":"+myNode.Name, 100, ":")
	conn.Write([]byte(fileNameFilled))

	// getting file size
	var bufferFileSize [64]byte
	conn.Read(bufferFileSize[:])

	fileSize, err := strconv.ParseInt(strings.TrimRight(string(bufferFileSize[:]), ":"), 10, 64)
	common.CheckError(err)

	fmt.Printf("Getting '%s' (%s) from '%s (%s)'...\n", fileName, humanize.Bytes(uint64(fileSize)), serverName, addr)

	// checking if we are speed limited
	var bufferSpeedLimited [1]byte
	conn.Read(bufferSpeedLimited[:])

	if string(bufferSpeedLimited[:]) == "1" {
		fmt.Println("Your download speed is limited to 10kB/s (because you share less files than average!)")
	}

	// creating the file
	newFile, err := os.Create(dir + fileName)
	common.CheckError(err)

	defer newFile.Close()

	// getting the file in chunks
	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, conn, (fileSize - receivedBytes))
			conn.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, conn, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Println("File received!")
}
