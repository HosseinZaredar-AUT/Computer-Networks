package udp

import (
	"P2P-File-Sharing/common"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// a function responsble for handling incoming discovery messages
func handleDiscovery(message string, clusterMap map[string]string, cmMutex *sync.Mutex, averageNumFiles *float64) {

	// updating cluster map
	nodes := strings.Split(message, ",")
	cmMutex.Lock() // lock cluster map
	for _, node := range nodes {
		fields := strings.Fields(node)
		clusterMap[fields[0]] = fields[1]
	}
	cmMutex.Unlock() // unlock cluster map

	// updating average number of files being served by peers
	sum := 0
	numOfNodes := 0
	for _, info := range clusterMap {
		num, err := strconv.Atoi(strings.Split(info, ";")[1])
		common.CheckError(err)
		sum += num
		numOfNodes++
	}
	*averageNumFiles = float64(sum) / float64(numOfNodes)
}

// a function responsible for handing incoming file request messages
func handleFileRequest(fileName string, dir string, myNode common.Node, conn *net.UDPConn, clientAddr *net.UDPAddr, numServing *int, maxClients int) {

	// excluding "temp"
	if fileName == "temp" {
		return
	}

	// openning the user's directory
	f, err := os.Open(dir)
	common.CheckError(err)

	// reading the contents
	files, err := f.Readdir(-1)
	common.CheckError(err)

	// closing the directory
	err = f.Close()
	common.CheckError(err)

	// going through all entires of the directory in search of the requested file
	for _, file := range files {
		if !file.IsDir() && file.Name() == fileName {

			// if we are serving maximum number of clients
			if *numServing >= maxClients {
				// telling the client that we are busy
				conn.WriteToUDP([]byte("busy, "+myNode.Name), clientAddr)

			} else { // if the node is ready to serve
				// telling the clients that we're ready to send the file
				info := myNode.Name + " " + myNode.GlobalIP + ":" + myNode.TCPPort
				t := strconv.FormatInt(time.Now().UnixNano(), 10) // the time that we are sending this message
				conn.WriteToUDP([]byte(t+","+info), clientAddr)
			}

			break
		}
	}
}

//Server UDP server which is responsible for:
// 1. getting discovery messages and updating cluster map and average files shared by peers
// 2. getting file request and responding properly
func Server(clusterMap map[string]string, myNode common.Node, dir string, cmMutex *sync.Mutex, numServing *int, averageNumFiles *float64, maxClients int) {

	// creating proper address
	service := myNode.LocalIP + ":" + myNode.UDPPPort
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	common.CheckError(err)

	// listening...
	conn, err := net.ListenUDP("udp", udpAddr)
	common.CheckError(err)

	// forever
	for {

		// wait for incoming message
		var buffer [1024]byte
		_, clientAddr, err := conn.ReadFromUDP(buffer[:])
		if err != nil {
			continue
		}

		// trimming the buffer to get the message
		message := strings.TrimRight(string(buffer[:]), "\x00")

		if message[0:4] == "dis:" { // if it's discovery message
			handleDiscovery(message[4:], clusterMap, cmMutex, averageNumFiles)
		} else if message[0:4] == "req:" { // if it's file request
			go handleFileRequest(message[4:], dir, myNode, conn, clientAddr, numServing, maxClients)
		}

	}
}
