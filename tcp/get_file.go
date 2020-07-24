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
func GetFile(fileName string, name string, addr string, dir string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	common.CheckError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	common.CheckError(err)

	// sending file name
	fileNameFilled := fillString(fileName, 64)
	conn.Write([]byte(fileNameFilled))

	// getting file size
	var bufferFileSize [64]byte
	conn.Read(bufferFileSize[:])

	fileSize, err := strconv.ParseInt(strings.TrimRight(string(bufferFileSize[:]), ":"), 10, 64)
	common.CheckError(err)

	fmt.Printf("Getting '%s' (%s) from '%s (%s)'...\n", fileName, humanize.Bytes(uint64(fileSize)), name, addr)

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
