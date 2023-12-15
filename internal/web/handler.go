package web

import (
	"encoding/json"
	"log"
	"logbox/internal/db"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

func queryHandler(w http.ResponseWriter, r *http.Request) {

	// s(tart), c(ount), q(uery)

	q := r.URL.Query().Get("q")

	logs := db.Query(q)

	json.NewEncoder(w).Encode(logs)

	log.Printf("%s [%s] \"%s %s\" %s\n", r.RemoteAddr, time.Now().Format(time.RFC3339), r.Method, r.RequestURI, "200 OK")
}

func wsHandler(ws *websocket.Conn) {

	log.Printf("%s [%s] \"%s %s\" connected.\n", ws.Request().RemoteAddr, time.Now().Format(time.RFC3339), ws.Request().Method, ws.Request().RequestURI)

	for {
		data := "Data from server: " + time.Now().Format(time.RFC3339)

		n, err := ws.Write([]byte(data))
		if n != len(data) || err != nil {
			log.Printf("Error writing to websocket of %s: %v\n", ws.RemoteAddr().String(), err)
			break
		}

		time.Sleep(2 * time.Second)
	}

	log.Printf("%s [%s] \"%s %s\" connection closed.\n", ws.Request().RemoteAddr, time.Now().Format(time.RFC3339), ws.Request().Method, ws.Request().RequestURI)
}
