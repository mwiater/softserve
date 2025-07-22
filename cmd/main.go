// cmd/main.go
package main

import (
	"log"

	cmd "github.com/mwiater/softserve/cmd/softserve"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
