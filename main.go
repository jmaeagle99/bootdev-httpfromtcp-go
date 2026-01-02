package main

import (
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

	buffer := make([]byte, 8)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("Error reading file: %v\n", err)
		}
		fmt.Printf("read: %s\n", string(buffer[:n]))
	}
}
