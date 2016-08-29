package main

//  Steps:
//  1. Create UDP socket
//  2. Listen for requests from master node for current time
//  3. Return current time to master node
//  4. Listen for response from master node for time delta
//  5. Return to step 2.

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"
)

func runSlave(address string) {
	fmt.Println("Beginning Clock Synchronization...")
	fmt.Printf("Creating UDP socket for %s\n", address)

	serverAddr, err := net.ResolveUDPAddr("udp", address)
	checkError(true, err)

	sock, err := net.ListenUDP("udp", serverAddr)
	checkError(true, err)
	defer sock.Close()

	var delta int64

	for {
		fmt.Printf("Listening at %v\n", serverAddr)
		readBuf := make([]byte, 1024)
		n, addr, err := sock.ReadFromUDP(readBuf)
		if err != nil {
			fmt.Println(err)
			continue
		}

		now := time.Now().Unix()
		fmt.Printf("Before adjustment: %v\n", time.Unix(now+delta, 0))

		msg := strings.TrimSpace(string(readBuf[:n]))
		fmt.Println(msg)

		writeBuf := make([]byte, 1024)
		bytes := binary.PutVarint(writeBuf, now)
		if bytes <= 0 {
			fmt.Println("Error encoding time")
			continue
		}

		fmt.Println("Making request to ", addr)
		sock.WriteToUDP(writeBuf, addr)

		readBuf = make([]byte, 1024)
		n, addr, err = sock.ReadFromUDP(readBuf)
		if err != nil {
			fmt.Println(err)
			continue
		}

		delta, bytes = binary.Varint(readBuf[:n])
		if bytes <= 0 {
			continue
		}
		fmt.Printf("Received adjustment of %v\n", delta)
		fmt.Printf("After adjustment %v\n", time.Unix(now+delta, 0))
	}
}
