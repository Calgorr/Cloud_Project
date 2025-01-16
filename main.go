package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"order_system/internal/http"
)

func main() {
	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	fmt.Println("Current directory:", dir)

	// List files in the current directory
	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	fmt.Println("Files in the current directory:")
	for _, file := range files {
		fmt.Println(file.Name())
	}

	// Read and print the contents of config.yaml
	configFilePath := "/app/config.yaml"
	file, err := os.Open(configFilePath)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", configFilePath, err)
		return
	}
	defer file.Close()

	fmt.Printf("Contents of %s:\n", configFilePath)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading %s: %v\n", configFilePath, err)
		return
	}

	// Start the HTTP server
	s := http.NewServer()
	err = s.Start("8080")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}
}
