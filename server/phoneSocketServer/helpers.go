package main

import (
	"fmt"
	"net"
	"os"
)

//this is a very hacky way to do it...
func findIpAddress() string {
	name, err := os.Hostname()
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
	}

	return addrs[0]
}
