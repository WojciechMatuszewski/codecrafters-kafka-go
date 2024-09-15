package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

const (
	Host string = "0.0.0.0"
	Port string = "9092"
)

func main() {
	address := fmt.Sprintf("%s:%s", Host, Port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Failed to bind to port 9092", err)

		os.Exit(1)
	}

	_, err = l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
}
