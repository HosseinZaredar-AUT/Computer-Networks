package cli

import (
	"P2P-File-Sharing/common"
	"P2P-File-Sharing/tcp"
	"P2P-File-Sharing/udp"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// RunCLI a command-line user inteface
func RunCLI(clusterMap map[string]string, myNode common.Node, dir string, averageNumFiles *float64, timeOut int) {
	state := 0
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to Netwolf P2P File Sharing System!")

	for {
		switch state {
		case 0: // main menu
			fmt.Println("")
			fmt.Println("1. Get a file")
			fmt.Println("2. See the list of peers")
			fmt.Println("3. See your status")
			fmt.Println("4. See your files")
			fmt.Printf("Please choose a command: ")

			command, _ := reader.ReadString('\n')
			command = strings.TrimRight(command, "\n")
			if command == "1" {
				state = 1
			} else if command == "2" {
				state = 2
			} else if command == "3" {
				state = 3
			} else if command == "4" {
				state = 4
			}

		case 1: // get file
			fmt.Printf("Please enter file name: ")
			fileName, _ := reader.ReadString('\n')
			fileName = strings.TrimRight(fileName, "\n")

			// requesting for the file
			res := udp.FileRequest(fileName, clusterMap, myNode, timeOut)

			if res == "not found" {
				fmt.Println("Not found!")
			} else if res == "busy" {
				fmt.Println("The file was found, but the node(s) are busy at the moment.")
				fmt.Println("Please try again later.")
			} else {
				fields := strings.Fields(res)

				// getting the file
				tcp.GetFile(fileName, fields[0], fields[1], dir, myNode, averageNumFiles)
			}

			state = 0

		case 2: // list of nodes
			fmt.Println("Cluster List:")
			fmt.Println(clusterMap)
			state = 0

		case 3: // status
			fmt.Println("Status:")
			fmt.Printf("Name: %s\n", myNode.Name)
			fmt.Printf("Global IP: %s\n", myNode.GlobalIP)
			fmt.Printf("Local IP: %s\n", myNode.LocalIP)
			fmt.Printf("UDP server running on port '%s'\n", myNode.UDPPPort)
			fmt.Printf("TCP server running on port '%s'\n", myNode.TCPPort)
			state = 0

		case 4: // list of files

			// openning the directory
			f, err := os.Open(dir)
			common.CheckError(err)

			// reading the contents
			files, err := f.Readdir(-1)
			common.CheckError(err)

			// closing the directory
			err = f.Close()
			common.CheckError(err)

			fmt.Println("Your files:")

			// going through all entires of the directory in search of files
			for _, file := range files {
				if !file.IsDir() {
					fmt.Println(file.Name())
				}
			}

			state = 0
		}

		fmt.Println()
	}
}
