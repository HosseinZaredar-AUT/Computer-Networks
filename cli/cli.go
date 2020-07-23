package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// RunCLI ...
func RunCLI(clusterMap map[string]string) {
	state := 0
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to Netwolf P2P File Sharing System!")

	for {
		switch state {
		case 0: // main menu
			fmt.Println("")
			fmt.Println("1. See cluster list.")
			fmt.Println("2. Get a file.")
			fmt.Printf("Please choose a command: ")

			command, _ := reader.ReadString('\n')
			command = strings.TrimRight(command, "\n")
			if command == "1" {
				state = 1
			} else if command == "2" {
				state = 2
			}

		case 1: // list of nodes
			fmt.Println("Cluster List:")
			fmt.Println(clusterMap)
			state = 0
		}

		fmt.Println()
	}
}
