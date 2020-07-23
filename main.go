package main

import (
	"P2P-File-Sharing/cli"
	"P2P-File-Sharing/common"
	"P2P-File-Sharing/tcp"
	"P2P-File-Sharing/udp"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/phayes/freeport"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// this function reads the cluster nodes file, updates clusterMap
// and returns the address of this machine
func readClusterNodes(clusterMap map[string]string, listPath string, myNode *common.Node) {
	f, err := os.Open(listPath)
	checkError(err)

	defer func() {
		err := f.Close()
		checkError(err)
	}()

	s := bufio.NewScanner(f)

	// getting my own node information from the first line
	s.Scan()
	fields := strings.Fields(s.Text())
	myNode.Name = fields[0]
	address := strings.Split(fields[1], ":")
	myNode.IP = address[0]
	myNode.UDPPPort = address[1]

	clusterMap[fields[0]] = fields[1]

	// getting other nodes
	for s.Scan() {
		fields := strings.Fields(s.Text())
		clusterMap[fields[0]] = fields[1]
	}

	err = s.Err()
	checkError(err)
}

func main() {

	// parse flags
	listPath := flag.String("l", "", "cluster list file path")
	dir := flag.String("d", "", "directory path")
	flag.Parse()

	// read list of cluster nodes from file
	clusterMap := make(map[string]string) // a map from name to IP address

	var myNode common.Node

	readClusterNodes(clusterMap, *listPath, &myNode)
	// fmt.Println("my node: ", myNode)
	// fmt.Println("initial cluster map:", clusterMap)

	// find a free TCP port
	port, err := freeport.GetFreePort()
	myNode.TCPPort = strconv.Itoa(port)
	checkError(err)

	// mutex for accessing cluster map
	var cmMutex sync.Mutex

	// run udp server
	go udp.Server(clusterMap, myNode, *dir, &cmMutex)

	// run discover client
	go udp.DiscoverService(clusterMap, myNode, &cmMutex)

	// run TCP server
	go tcp.Server(myNode, *dir)

	// run CLI in the main goroutine
	cli.RunCLI(clusterMap, myNode, *dir)

}
