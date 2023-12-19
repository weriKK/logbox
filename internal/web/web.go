package web

import (
	"log"
	"logbox/internal/common"
	"net"
	"net/http"
	"syscall"

	"golang.org/x/net/websocket"
)

type webServer struct {
	clientNotif *common.ClientNotifier
}

func NewWebServer(cn *common.ClientNotifier) *webServer {
	return &webServer{
		clientNotif: cn,
	}
}

func (ws *webServer) Run() {

	listener, err := net.Listen("tcp", "localhost:8080")
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

	fs := http.FileServer(http.Dir("./webui"))
	http.Handle("/", fs)

	http.Handle("/events", websocket.Handler(ws.websocketHandler))

	http.HandleFunc("/query", ws.queryHandler)

	// go func() {
	if err := http.Serve(listener, nil); err != nil {
		log.Fatal(err)
	}
	// }()

	// select {}
}
