package main

import (
	"fmt"

	"github.com/bl4ko/certstream-go"
)

func main() {
	certstreamServerURL := "wss://certstream.calidog.io"
	timeout := 15
	stream, errStream := certstream.EventStream(true, certstreamServerURL, timeout)
	for {
		select {
		case message := <-stream:
			fmt.Printf("Received stream: %+v\n\n", message)
		case err := <-errStream:
			fmt.Printf("Received error: %s\n\n", err)
		}
	}
}
