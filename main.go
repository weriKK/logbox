package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"logbox/internal/db"
	"net"
	"net/http"
	"syscall"
)

func queryHandler(w http.ResponseWriter, r *http.Request) {

	// s(tart), c(ount), q(uery)

	q := r.URL.Query().Get("q")

	logs := db.Query(q)

	json.NewEncoder(w).Encode(logs)
}

func main() {
	go runWebServer()
	go runIngestServer()

	select {}
}

func runIngestServer() {

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

		go handleIngestClient(conn)
	}
}

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

func runWebServer() {

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

	// go func() {
	if err := http.Serve(listener, nil); err != nil {
		log.Fatal(err)
	}
	// }()

	// select {}
}
