package main

//  Steps:
//  1. Create UDP socket
//  2. Create requests to each slave for its time
//  3. After max reached or all slaves responded compute average
//  4. Respond to each slave the delta of the average and their reported time
//  5. Return to step 2.

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func runMaster(address string, slavesList []string) {
	fmt.Println("Beginning Clock Synchronization...")
	fmt.Println("Creating UDP socket to request slave nodes times")

	serverAddr, err := net.ResolveUDPAddr("udp", address)
	checkError(true, err)

	sock, err := net.ListenUDP("udp", serverAddr)
	checkError(true, err)
	defer sock.Close()

	slaveNodes := make([]*slaveNode, len(slavesList))

	for {
		for i, slaveAddr := range slavesList {
			tempNodeAddr, err := net.ResolveUDPAddr("udp", slaveAddr)
			if err != nil {
				fmt.Println(err)
				continue
			}
			slaveNodes[i] = &slaveNode{addr: tempNodeAddr}
		}

		// Poll all slave nodes for their time
		masterDelta := time.Now().Unix()
		fmt.Printf("Before adjustment: %v\n", time.Unix(masterDelta, 0))
		slaveResponses := pollSlaves(slaveNodes, sock)

		// Compute the algorithm for nodes that have responded
		masterDelta = berkeleyTime(masterDelta, slaveResponses, slaveNodes)
		fmt.Printf("After adjustment: %v\n", time.Unix(masterDelta, 0))

		// Send the nodes their new time deltas
		tellTheSlaves(slaveNodes, sock)
		// Sleep
		time.Sleep(5 * time.Second)
	}
}

type slaveNode struct {
	addr      *net.UDPAddr
	timeTicks int64
	delta     int64
}

type slaveResponse struct {
	addr      *net.UDPAddr
	timeTicks int64
}

func pollSlaves(slaves []*slaveNode, masterSock *net.UDPConn) chan slaveResponse {
	res := make(chan slaveResponse, len(slaves))
	for _, slave := range slaves {
		masterSock.WriteToUDP([]byte("Master requesting slave time"), slave.addr)
		go func(address *net.UDPAddr) {
			buf := make([]byte, 1024)
			fmt.Println("Making request to ", address)
			n, addr, err := masterSock.ReadFromUDP(buf)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Err is number of bytes read, == 0: buf too small, < 0: overflow
			ticks, bytes := binary.Varint(buf[:n])
			if bytes <= 0 {
				return
			}
			fmt.Printf("Received %v from %v\n", time.Unix(ticks, 0), addr)
			res <- slaveResponse{addr: address, timeTicks: ticks}
		}(slave.addr)
	}

	return res
}

func berkeleyTime(now int64, slaveResponses chan slaveResponse, slaves []*slaveNode) int64 {
	totalTime := now

	timeouts := make(chan bool, len(slaves))
	for _ = range slaves {
		go func() {
			time.Sleep(1 * time.Second)
			timeouts <- true
		}()
	}

	var responses []slaveResponse
	for _ = range slaves {
		select {
		case <-timeouts:
			break
		case res := <-slaveResponses:
			responses = append(responses, res)
			break
		}
	}

	numResponses := 1
	fmt.Printf("%v - Master\n", time.Unix(totalTime, 0))
	for i, nodeRes := range responses {
		if nodeRes.timeTicks != 0 {
			numResponses++

			fmt.Printf("%v - Slave #%v\n", time.Unix(nodeRes.timeTicks, 0), i)

			totalTime += nodeRes.timeTicks
			for _, node := range slaves {
				if node.addr.String() == nodeRes.addr.String() {
					node.timeTicks = nodeRes.timeTicks
				}
			}
		}
	}

	// TODO throw out outliers

	// 1 response is only master
	fmt.Printf("%v node(s) responded.\n", numResponses-1)
	var masterDelta int64
	if numResponses > 1 {
		masterDelta = totalTime / int64(numResponses)
	} else {
		masterDelta = now
	}

	for _, node := range slaves {
		node.delta = masterDelta - node.timeTicks
	}

	return masterDelta
}

func tellTheSlaves(slaves []*slaveNode, sock *net.UDPConn) {
	for _, node := range slaves {
		writeBuf := make([]byte, 1024)
		bytes := binary.PutVarint(writeBuf, node.delta)
		if bytes <= 0 {
			fmt.Println("Error encoding the node's delta")
			continue
		}

		fmt.Println("Making request to ", node.addr)
		sock.WriteToUDP(writeBuf, node.addr)
	}
}
