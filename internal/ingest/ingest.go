package ingest

import (
	"log"
	"logbox/internal/common"
	"net"
	"syscall"
)

type ingestServer struct {
	clientNotifier *common.ClientNotifier
}

func NewIngestServer(cn *common.ClientNotifier) *ingestServer {
	return &ingestServer{
		clientNotifier: cn,
	}
}

func (is *ingestServer) Run() {

	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		log.Fatalf("Error creating TCP listener: %v\n", err)
	}

	file, err := listener.(*net.TCPListener).File()
	if err != nil {
		log.Fatalf("Error retrieving listener file descriptor: %v\n", err)
	}

	err = syscall.SetsockoptInt(int(file.Fd()), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		log.Fatalf("Error setting SO_REUSEADDR: %v\n", err)
	}

	defer listener.Close()

	log.Printf("Server listening on %s\n", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %v\n", err)
		}

		go is.handleIngestClient(conn)
	}
}
