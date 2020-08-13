package tcp

import (
	"P2P-File-Sharing/common"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/dustin/go-humanize"
)

// GetFile sends a transmit request to a peer
func GetFile(fileName string, serverName string, addr string, dir string, myNode common.Node, averageNumFiles *float64) {

	// creating proper address
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	common.CheckError(err)

	// getting connected to server
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	common.CheckError(err)

	// sending "filename:myNode.name"
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
		fmt.Println("Your download speed is limited (due to sharing less files than the average (=", *averageNumFiles, ")")
	}

	// creating an empty file named "temp"
	newFile, err := os.Create(dir + "temp")
	common.CheckError(err)

	// getting the file in chunks
	var receivedBytes int64

	bar := pb.Full.Start64(fileSize)
	bar.Set(pb.Bytes, true)

	// until we received all of the file
	for {

		// if the remaining chuck of file is smaller than buffer size
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, conn, (fileSize - receivedBytes))
			conn.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}

		// if the remaining chuck of file is equal to or more than buffer size
		io.CopyN(newFile, conn, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
		bar.Add(BUFFERSIZE)
	}

	bar.Finish()

	// closing the file
	newFile.Close()

	// renaming the file to the real name
	os.Rename(dir+"temp", dir+fileName)

	fmt.Println("File received!")
}
