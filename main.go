package main

import (
	"logbox/internal/ingest"
	"logbox/internal/web"
)

func main() {
	go web.RunWebServer()
	go ingest.RunIngestServer()

	select {}
}
