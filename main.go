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

// this function reads the cluster nodes file, fills clusterMap and user's own node info
func readClusterNodes(clusterMap map[string]string, listPath string, myNode *common.Node, dir string) {

	// openning the file
	f, err := os.Open(listPath)
	common.CheckError(err)

	// close the file at the end
	defer func() {
		err := f.Close()
		common.CheckError(err)
	}()

	// scanner to read the file
	s := bufio.NewScanner(f)

	// getting my own node information from the first line
	s.Scan()
	fields := strings.Fields(s.Text())
	clusterMap[fields[0]] = fields[1] + ";" + "0"

	// updating myNode struct
	myNode.Name = fields[0]
	address := strings.Split(fields[1], ":")
	myNode.IP = address[0]
	myNode.UDPPPort = address[1]

	// getting other nodes
	for s.Scan() {
		fields := strings.Fields(s.Text())
		clusterMap[fields[0]] = fields[1] + ";" + "0"
	}

	err = s.Err()
	common.CheckError(err)
}

func main() {

	// parse flags
	listPath := flag.String("l", "", "cluster list file path")
	dir := flag.String("d", "", "directory path")
	interval := flag.String("i", "1", "discovery message interval (seconds)")
	tOut := flag.String("t", "2", "timeout for file requests (seconds)")
	max := flag.String("m", "5", "maximum simultaneous TCP clients")

	flag.Parse()

	// adding '/' ('\') at the end of dir, if it doesn't have that
	if !strings.HasSuffix(*dir, string(os.PathSeparator)) {
		*dir = *dir + string(os.PathSeparator)
	}

	// converting 'interval', 'tOut' and 'max' to int
	discoveryInterval, err := strconv.Atoi(*interval)
	common.CheckError(err)
	timeOut, err := strconv.Atoi(*tOut)
	common.CheckError(err)
	MaxClients, err := strconv.Atoi(*max)
	common.CheckError(err)

	// read list of cluster nodes from file

	// a map from name to (IP:Port;numOfFilesShared)
	// for exmaple: "N1" -> "127.0.0.1:1500;7"
	clusterMap := make(map[string]string)

	// user's own node information
	var myNode common.Node

	readClusterNodes(clusterMap, *listPath, &myNode, *dir)

	// find a free TCP port
	port, err := freeport.GetFreePort()
	myNode.TCPPort = strconv.Itoa(port)
	common.CheckError(err)

	// mutex for accessing cluster map
	// (in udp.server.go for updating it, in discovery_service.go for reading it)
	var cmMutex sync.Mutex

	// number of clients being served, which is the number of clients that are getting a file from TCP server
	// (TCP and UDP servers will have access to this variable)
	numServing := 0

	// average number of files being served by peers
	// (which gets updated by dicovery messages)
	averageNumFiles := 0.0

	// run udp server (responsibe for getting discovery messages and file requests)
	go udp.Server(clusterMap, myNode, *dir, &cmMutex, &numServing, &averageNumFiles, MaxClients)

	// run discover service (responsible for sending discovery messages)
	go udp.DiscoverService(clusterMap, myNode, &cmMutex, *dir, discoveryInterval)

	// run TCP server (responsible for getting a file name and transfering the file)
	go tcp.Server(clusterMap, myNode, *dir, &numServing, &averageNumFiles)

	// run CLI in the main goroutine
	cli.RunCLI(clusterMap, myNode, *dir, &averageNumFiles, timeOut)
}
