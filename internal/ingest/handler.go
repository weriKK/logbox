package ingest

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"logbox/internal/common"
	"logbox/internal/db"
	"net"
)

func (is *ingestServer) handleIngestClient(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 65536)

	isNotificationNeeded := false

	for {
		bytesRead, err := conn.Read(buffer)
		log.Printf("CONNREAD ERR: %v\n", err)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from connection: %v\n", err)
			}
			break
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
				isNotificationNeeded = true
			}
		}
	}

	if isNotificationNeeded {
		is.clientNotifier.NotifyAll()
	}
}
