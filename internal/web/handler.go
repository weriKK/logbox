package web

import (
	"encoding/json"
	"log"
	"logbox/internal/common"
	"logbox/internal/db"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

func queryHandler(w http.ResponseWriter, r *http.Request) {

	// s(tart), c(ount), q(uery)

	q := r.URL.Query().Get("q")

	logs := db.Query(db.NewSelectMessageQuery().WithPattern(q).Build())

	json.NewEncoder(w).Encode(logs)

	log.Printf("%s [%s] \"%s %s\" %s\n", r.RemoteAddr, time.Now().Format(time.RFC3339), r.Method, r.RequestURI, "200 OK")

}

type BrowserStateMessage struct {
	LastLogMessageId int    `json:"lastLogMessageId"`
	QueryString      string `json:"queryString"`
	Timestamp        int    `json:"timestamp"`
}

func wsHandler(ws *websocket.Conn) {
	defer ws.Close()

	log.Printf("%s [%s] \"%s %s\" connected.\n", ws.Request().RemoteAddr, time.Now().Format(time.RFC3339), ws.Request().Method, ws.Request().RequestURI)

	var browserState BrowserStateMessage
	if err := websocket.JSON.Receive(ws, &browserState); err != nil {
		log.Println("Error receiving JSON from WebSocket:", err)
		return
	}

	log.Printf("%+v\n", browserState)

	for {

		q := db.NewSelectMessageQuery().WithStartingId(browserState.LastLogMessageId).WithPattern(browserState.QueryString).Build()
		logs := db.Query(q)

		if 0 < len(*logs) {
			err := json.NewEncoder(ws).Encode(logs)
			if err != nil {
				log.Printf("Error writing to websocket of %s: %v\n", ws.RemoteAddr().String(), err)
				break
			}

			browserState.LastLogMessageId = lastMessageId(browserState.LastLogMessageId, logs)
		}

		time.Sleep(10 * time.Second)
	}

	log.Printf("%s [%s] \"%s %s\" connection closed.\n", ws.Request().RemoteAddr, time.Now().Format(time.RFC3339), ws.Request().Method, ws.Request().RequestURI)
}

func lastMessageId(currentId int, logs *[]common.LogMessage) int {

	var lastId = currentId

	for _, m := range *logs {
		if lastId < m.Id {
			lastId = m.Id
		}
	}

	return lastId
}
