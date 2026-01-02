package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const filePath = "./messages.txt"

func main() {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file %s: %v\n", filePath, err)
	}
	defer file.Close()

	readBuffer := make([]byte, 8)
	currentLineBuffer := make([]byte, 0, 100)

	for {
		n, err := file.Read(readBuffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("Error reading file: %v\n", err)
		}

		parts := bytes.Split(readBuffer[:n], []byte("\n"))

		for partIndex := 0; partIndex < len(parts)-1; partIndex++ {
			currentLineBuffer = append(currentLineBuffer, parts[partIndex]...)
			fmt.Printf("read: %s\n", string(currentLineBuffer))
			currentLineBuffer = currentLineBuffer[:0]
		}

		currentLineBuffer = append(currentLineBuffer, parts[len(parts)-1]...)
	}

	fmt.Printf("read: %s\n", string(currentLineBuffer))
}
