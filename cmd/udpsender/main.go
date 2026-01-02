package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	address, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatalf("Error resolving address: %v\n", err)
	}

	connection, err := net.DialUDP("udp", nil, address)
	if err != nil {
		log.Fatalf("Error dialing address: %v\n", err)
	}
	defer connection.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading from stdin: %v\n", err)
		}
		_, err = connection.Write([]byte(line))
		if err != nil {
			log.Fatalf("Error write to address: %v\n", err)
		}
	}
}
