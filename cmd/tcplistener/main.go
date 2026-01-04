package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jmaeagle99/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Error creating tcp listener: %v\n", err)
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %v", err)
		}

		go func() {
			defer connection.Close()

			request, err := request.RequestFromReader(connection)
			if err != nil {
				log.Fatalf("Error reading request: %v", err)
			}

			fmt.Println("Request line:")
			fmt.Printf("- Method: %s\n", request.RequestLine.Method)
			fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
			fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)
		}()
	}
}
