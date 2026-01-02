package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
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
			log.Fatalf("Error accepting connection: %v\n", err)
		}

		fmt.Println("Server accepted connection.")

		for line := range getLinesChannel(connection) {
			fmt.Println(line)
		}

		fmt.Println("Server connection closed.")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesChannel := make(chan string)

	go func() {
		defer close(linesChannel)

		readBuffer := make([]byte, 8)
		currentLineBuffer := make([]byte, 0, 100)

		for {
			n, err := f.Read(readBuffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatalf("Error reading file: %v\n", err)
			}

			parts := bytes.Split(readBuffer[:n], []byte("\n"))

			for partIndex := 0; partIndex < len(parts)-1; partIndex++ {
				currentLineBuffer = append(currentLineBuffer, parts[partIndex]...)
				linesChannel <- string(currentLineBuffer)
				currentLineBuffer = currentLineBuffer[:0]
			}

			currentLineBuffer = append(currentLineBuffer, parts[len(parts)-1]...)
		}

		linesChannel <- string(currentLineBuffer)
	}()

	return linesChannel
}
