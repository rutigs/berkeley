package main

//  Steps:
//  1. Create UDP socket
//  2. Create requests to each slave for its time
//  3. After max reached or all slaves responded compute average
//  4. Respond to each slave the delta of the average and their reported time
//  5. Return to step 2.

import (
	"fmt"
)

func runMaster(address string, slavesList []string) {
	fmt.Println("Beginning Clock Synchronization...")
	fmt.Println("Creating UDP socket to request slave nodes times")
}
