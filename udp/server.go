package udp

import (
	"P2P-File-Sharing/common"
	"fmt"
	"net"
	"os"
	"strings"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func handleDiscovery(message string, clusterMap map[string]string) {
	// updating cluster map
	nodes := strings.Split(message, ",")
	for _, node := range nodes {
		fields := strings.Fields(node)
		clusterMap[fields[0]] = fields[1]
	}

	// fmt.Println("cluster map:", clusterMap)
}

func handleFileRequest(fileName string, dir string, myNode common.Node, conn *net.UDPConn, clientAddr *net.UDPAddr) {

	fmt.Println("got a file request!")
	f, err := os.Open(dir)
	checkError(err)

	files, err := f.Readdir(-1)
	checkError(err)

	err = f.Close()
	checkError(err)

	for _, file := range files {
		if !file.IsDir() && file.Name() == fileName[0:len(file.Name())] { //TODO: improve this
			// send message to client
			conn.WriteToUDP([]byte(myNode.Address+": I have '"+fileName+"'"), clientAddr)
			break
		}
	}
}

//Server ...
func Server(clusterMap map[string]string, myNode common.Node, dir string) {

	udpAddr, err := net.ResolveUDPAddr("udp4", myNode.Address)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	// fmt.Println("UDP server listining on", udpAddr)

	for {
		var buffer [512]byte
		_, clientAddr, err := conn.ReadFromUDP(buffer[:])
		if err != nil {
			continue
		}

		message := string(buffer[:])

		if message[0:4] == "dis:" { // if it's discovery message
			handleDiscovery(message[4:], clusterMap)
		} else if message[0:4] == "req:" { // if it's file request message
			go handleFileRequest(message[4:], dir, myNode, conn, clientAddr)
		}

	}
}
