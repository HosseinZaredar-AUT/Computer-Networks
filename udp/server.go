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

const MAXSERVING = 1

func handleDiscovery(message string, clusterMap map[string]string, cmMutex *sync.Mutex) {
	// updating cluster map
	nodes := strings.Split(strings.TrimRight(message, "\x00"), ",")
	cmMutex.Lock()
	for _, node := range nodes {
		fields := strings.Fields(node)
		clusterMap[fields[0]] = fields[1]
	}
	cmMutex.Unlock()
}

func handleFileRequest(fileName string, dir string, myNode common.Node, conn *net.UDPConn, clientAddr *net.UDPAddr, numServing *int) {

	f, err := os.Open(dir)
	common.CheckError(err)

	files, err := f.Readdir(-1)
	common.CheckError(err)

	err = f.Close()
	common.CheckError(err)

	for _, file := range files {
		if !file.IsDir() && file.Name() == fileName[0:len(file.Name())] { //TODO: improve this

			// check if we are serving maximum number of clients
			if *numServing >= MAXSERVING {
				conn.WriteToUDP([]byte("busy, "+myNode.Name), clientAddr)
			} else { // if the node is ready to serve
				info := myNode.Name + " " + myNode.IP + ":" + myNode.TCPPort
				t := strconv.FormatInt(time.Now().UnixNano(), 10)
				conn.WriteToUDP([]byte(t+","+info), clientAddr)
			}

			break
		}
	}
}

//Server ...
func Server(clusterMap map[string]string, myNode common.Node, dir string, cmMutex *sync.Mutex, numServing *int) {

	service := myNode.IP + ":" + myNode.UDPPPort
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	common.CheckError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	common.CheckError(err)

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
			go handleFileRequest(message[4:], dir, myNode, conn, clientAddr, numServing)
		}

	}
}
