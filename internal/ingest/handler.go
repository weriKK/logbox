package ingest

import (
	"bufio"
	"bytes"
	"log"
	"logbox/internal/common"
	"logbox/internal/db"
	"net"
)

func handleIngestClient(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 65536)

	for {
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Error reading from connection: %v\n", err)
			return
		}

		if 0 < bytesRead {

			data := buffer[0:bytesRead]
			scanner := bufio.NewScanner(bytes.NewReader(data))

			for scanner.Scan() {
				msg := scanner.Text()
				log.Printf("Received: %s\n", msg)

				logMessage := common.LogMessage{
					Message: msg,
				}

				db.Store(logMessage)
			}

		}
	}
}
