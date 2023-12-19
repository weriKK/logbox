package main

import (
	"logbox/internal/common"
	"logbox/internal/ingest"
	"logbox/internal/web"
)

func main() {

	cn := common.NewClientNotifier()

	webServer := web.NewWebServer(cn)
	ingestServer := ingest.NewIngestServer(cn)

	go webServer.Run()
	go ingestServer.Run()

	select {}
}
