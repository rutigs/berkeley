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
		usage()
		return false
	} else if !master && !slave {
		fmt.Println("Program must be run with either -m (master) or -s (slave)")
		usage()
		return false
	}

	if address == "" {
		fmt.Println("You must provide an ip:port with -addr")
		fmt.Println("eg. ./berkeley -s -addr=123.123.123.123:1337")
		fmt.Println("eg. ./berkeley -m -addr=:1337")
		usage()
		return false
	}

	if master && slavesFile == "" {
		fmt.Println("You must provide a json file containing a list of slave nodes with -slaves")
		fmt.Println("eg ./berkeley -m -addr=0.0.0.0:0 -slaves=slaves.json")
		usage()
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

func usage() {
	fmt.Println(`
Usage: ./berkeley (-m or -s) -addr=0.0.0.0:0 [-slaves=slavesJsonFile.json]
  -m      Run program as master node that will compute the synchronization algorithm
  -s      Run program as slave node that will listen for requests from the master node 
          for its current time, and receives a synchronization value
  -addr   IP:Port string for the program to run under eg. "-addr=127.0.0.1:1337"
  -slaves Name of json file containing the list of slaves nodes addresses
          Must be used with -m
	`)
}

func checkError(isFatal bool, err error) {
	if err != nil {
		if isFatal {
			log.Fatal(err)
		}
		fmt.Println(err)
	}
}
