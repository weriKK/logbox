package ingest

import (
	"bufio"
	"bytes"
	"log"
	"net"
)

func handleIngestClient(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 65536)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Error reading from connection: %v\n", err)
			return
		}
		if 0 < n {
			scanner := bufio.NewScanner(bytes.NewReader(buffer))
			for scanner.Scan() {
				log.Printf("Received: %s\n", scanner.Text())
			}

			// log.Printf("Received: %s\n", buffer[:n])
		}
	}
}
