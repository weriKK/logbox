package main

import (
	"encoding/json"
	"log"
	"logbox/internal/db"
	"net"
	"net/http"
	"syscall"
)

func queryHandler(w http.ResponseWriter, r *http.Request) {

	// s(tart), c(ount), q(uery)

	logs := db.Query("chicken")

	json.NewEncoder(w).Encode(logs)
}

func main() {
	runServer()
}

func runServer() {

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
	http.HandleFunc("/query", queryHandler)

	go func() {
		if err := http.Serve(listener, nil); err != nil {
			log.Fatal(err)
		}
	}()

	select {}
}
