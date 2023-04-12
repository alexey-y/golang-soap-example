package main

import (
	"fmt"
	"log"
	"soap-example/client"
)

func main() {
	result, err := client.CallAddService(12, 18)
	if err != nil {
		log.Fatalf("Error calling Add service: %v", err)
	}

	fmt.Printf("Result: %d\n", result)
}
