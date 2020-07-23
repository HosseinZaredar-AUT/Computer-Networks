package udp

import (
	"P2P-File-Sharing/common"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func handleDiscovery(message string, clusterMap map[string]string, cmMutex *sync.Mutex) {
	// updating cluster map
	nodes := strings.Split(message, ",")
	cmMutex.Lock()
	for _, node := range nodes {
		fields := strings.Fields(node)
		clusterMap[fields[0]] = fields[1]
	}
	cmMutex.Unlock()
}

func handleFileRequest(fileName string, dir string, myNode common.Node, conn *net.UDPConn, clientAddr *net.UDPAddr) {

	f, err := os.Open(dir)
	checkError(err)

	files, err := f.Readdir(-1)
	checkError(err)

	err = f.Close()
	checkError(err)

	for _, file := range files {
		if !file.IsDir() && file.Name() == fileName[0:len(file.Name())] { //TODO: improve this
			// send message to client
			info := myNode.Name + " " + myNode.IP + ":" + myNode.TCPPort
			t := strconv.FormatInt(time.Now().UnixNano(), 10)
			conn.WriteToUDP([]byte(t+","+info), clientAddr)
			break
		}
	}
}

//Server ...
func Server(clusterMap map[string]string, myNode common.Node, dir string, cmMutex *sync.Mutex) {

	service := myNode.IP + ":" + myNode.UDPPPort
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	for {
		var buffer [512]byte
		_, clientAddr, err := conn.ReadFromUDP(buffer[:])
		if err != nil {
			continue
		}

		message := string(buffer[:])

		if message[0:4] == "dis:" { // if it's discovery message
			handleDiscovery(message[4:], clusterMap, cmMutex)
		} else if message[0:4] == "req:" { // if it's file request message
			go handleFileRequest(message[4:], dir, myNode, conn, clientAddr)
		}

	}
}
