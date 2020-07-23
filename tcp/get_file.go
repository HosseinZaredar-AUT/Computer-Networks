package tcp

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

// GetFile ...
func GetFile(fileName string, name string, addr string, dir string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	// sending file name
	conn.Write([]byte(fileName))

	// getting file size
	var bufferFileSize [10]byte
	conn.Read(bufferFileSize[:])

	fmt.Println(bufferFileSize[:])
	fileSize, err := strconv.ParseInt(strings.TrimRight(string(bufferFileSize[:]), "\x00"), 10, 64)
	checkError(err)

	fmt.Printf("Getting '%s' (%d Bytes) from '%s (%s)'...\n", fileName, fileSize, name, addr)

	// os.Exit(0)

	// creating the file
	newFile, err := os.Create(dir + fileName)
	checkError(err)

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
