package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	master     bool
	slave      bool
	address    string
	slavesFile string
)

func init() {
	flag.BoolVar(&master, "m", false, "master time node")
	flag.BoolVar(&slave, "s", false, "slave time node")
	flag.StringVar(&address, "addr", "", "IP:Port to listen/request on")
	flag.StringVar(&slavesFile, "slaves", "", "Slaves json file")
	flag.Parse()
}

func validateFlags() bool {
	if master && slave {
		fmt.Println("Cannot be both master and slave node")
		return false
	} else if !master && !slave {
		fmt.Println("Program must be run with either -m (master) or -s (slave)")
		return false
	}

	if address == "" {
		fmt.Println("You must provide an ip:port with -addr")
		fmt.Println("eg. ./berkeley [-m or -s] -addr=123.123.123.123:1337")
		return false
	}

	if master && slavesFile == "" {
		fmt.Println("You must provide a json file containing a list of slave nodes with -slaves")
		fmt.Println("eg ./berkeley -m -addr=0.0.0.0:0 -slaves=slaves.json")
		return false
	}

	return true
}

func main() {
	if !validateFlags() {
		return
	}

	if master {
		fmt.Println("Running as master node")
		slaves := parseSlaves()
		runMaster(address, slaves)
	}

	if slave {
		fmt.Println("Running as slave node")
		runSlave(address)
	}
}

func parseSlaves() []string {
	file, fileErr := os.Open(slavesFile)
	if fileErr != nil {
		log.Fatal(fileErr)
	}

	var slaves []string
	decoder := json.NewDecoder(file)
	fileErr = decoder.Decode(&slaves)
	if fileErr != nil {
		log.Fatal(fileErr)
	}

	for index, slave := range slaves {
		fmt.Printf("Slave #%v: %v\n", index+1, slave)
	}

	return slaves
}
