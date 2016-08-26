package main

//  Steps:
//  1. Create UDP socket
//  2. Listen for requests from master node for current time
//  3. Return current time to master node
//  4. Listen for response from master node for time delta
//  5. Return to step 2.

import (
	"fmt"
)

func runSlave(address string) {
	fmt.Println("Beginning Clock Synchronization...")

	fmt.Printf("Creating UDP socket for %s\n", address)
}
