package main

import (
	"P2P-File-Sharing/cli"
	"P2P-File-Sharing/common"
	"P2P-File-Sharing/tcp"
	"P2P-File-Sharing/udp"
	"bufio"
	"flag"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/phayes/freeport"
)

// this function reads the cluster nodes file, updates clusterMap
// and returns the address of this machine
func readClusterNodes(clusterMap map[string]string, listPath string, myNode *common.Node) {
	f, err := os.Open(listPath)
	common.CheckError(err)

	defer func() {
		err := f.Close()
		common.CheckError(err)
	}()

	s := bufio.NewScanner(f)

	// getting my own node information from the first line
	s.Scan()
	fields := strings.Fields(s.Text())
	clusterMap[fields[0]] = fields[1]
	myNode.Name = fields[0]
	address := strings.Split(fields[1], ":")
	myNode.IP = address[0]
	myNode.UDPPPort = address[1]

	// getting other nodes
	for s.Scan() {
		fields := strings.Fields(s.Text())
		clusterMap[fields[0]] = fields[1]
	}

	err = s.Err()
	common.CheckError(err)
}

func main() {

	// parse flags
	listPath := flag.String("l", "", "cluster list file path")
	dir := flag.String("d", "", "directory path")
	flag.Parse()

	// adding '/' ('\') at the end of dir, if it doesn't have that
	if !strings.HasSuffix(*dir, string(os.PathSeparator)) {
		*dir = *dir + string(os.PathSeparator)
	}

	// read list of cluster nodes from file
	clusterMap := make(map[string]string) // a map from name to IP:Port
	var myNode common.Node
	readClusterNodes(clusterMap, *listPath, &myNode)

	// find a free TCP port
	port, err := freeport.GetFreePort()
	myNode.TCPPort = strconv.Itoa(port)
	common.CheckError(err)

	// mutex for accessing cluster map
	var cmMutex sync.Mutex

	// number of clients being served, which is the number of clients that are getting a file from TCP server
	// (TCP and UDP servers will have acess to this variable)
	numServing := 0

	// run udp server (responsibe fot getting discovery messages and file requests)
	go udp.Server(clusterMap, myNode, *dir, &cmMutex, &numServing)

	// run discover service (responsible for sending discovery messages)
	go udp.DiscoverService(clusterMap, myNode, &cmMutex)

	// run TCP server (responsible for getting file name and transfering the file)
	go tcp.Server(myNode, *dir, &numServing)

	// run CLI in the main goroutine
	cli.RunCLI(clusterMap, myNode, *dir)

}
