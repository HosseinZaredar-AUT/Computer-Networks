package cli

import (
	"P2P-File-Sharing/common"
	"P2P-File-Sharing/udp"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// RunCLI ...
func RunCLI(clusterMap map[string]string, myNode common.Node) {
	state := 0
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to Netwolf P2P File Sharing System!")

	for {
		switch state {
		case 0: // main menu
			fmt.Println("")
			fmt.Println("1. See cluster list")
			fmt.Println("2. Get a file")
			fmt.Println("3. Status")
			fmt.Printf("Please choose a command: ")

			command, _ := reader.ReadString('\n')
			command = strings.TrimRight(command, "\n")
			if command == "1" {
				state = 1
			} else if command == "2" {
				state = 2
			} else if command == "3" {
				state = 3
			}

		case 1: // list of nodes
			fmt.Println("Cluster List:")
			fmt.Println(clusterMap)
			state = 0

		case 2: // get file
			fmt.Printf("Please enter file name: ")
			fileName, _ := reader.ReadString('\n')
			res := udp.FileRequest(fileName, clusterMap, myNode)

			if res == "!" {
				fmt.Println("Not found!")
			} else {
				fmt.Println(res)
			}

			state = 0

		case 3: // status
			fmt.Println("Status:")
			fmt.Printf("Name (IP): %s (%s)\n", myNode.Name, myNode.IP)
			fmt.Printf("UDP server running on port '%s'\n", myNode.UDPPPort)
			fmt.Printf("TCP server running on port '%s'\n", myNode.TCPPort)
			state = 0

		}

		fmt.Println()
	}
}
