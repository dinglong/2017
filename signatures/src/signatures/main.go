package main

import (
	"fmt"
	"log"
	"os"

	"signatures/client"
)

var (
	notaryEndpoint = "https://192.168.1.50:4443"
	fqRepo         = "192.168.1.50:5000/alpine"
)

func main() {
	if len(os.Args) > 1 {
		notaryEndpoint = os.Args[1]
	}
	if len(os.Args) > 2 {
		fqRepo = os.Args[2]
	}

	fmt.Printf("notary endpoint: %s\n", notaryEndpoint)
	fmt.Printf("fully-qualified repository: %s\n", fqRepo)

	ts, err := client.GetTargets(notaryEndpoint, fqRepo)
	if err != nil {
		log.Printf("get targets error, %v\n", err)
		return
	}

	for _, t := range ts {
		fmt.Printf("----------------------------------\n")
		fmt.Printf("name: %v\n", t.Name)
		fmt.Printf("role: %v\n", t.Role)
		fmt.Printf("size: %v\n", t.Length)
		for k, v := range t.Hashes {
			fmt.Printf("\tkey: %v\n", k)
			fmt.Printf("\tval: %v\n", v)
		}
	}
}
