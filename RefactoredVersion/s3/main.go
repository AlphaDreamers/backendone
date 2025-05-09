package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	imagePath := "/Users/swanhtet1aungphyo/Downloads/IMG_4401.jpg"

	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal("Error closing file:", err)
		}
	}()

	all, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	// Encode the file content (image data) to Base64
	base64String := base64.StdEncoding.EncodeToString(all)

	// Print the Base64 encoded string
	fmt.Println(base64String)
}
