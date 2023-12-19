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

func (ws *webServer) queryHandler(w http.ResponseWriter, r *http.Request) {

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

func (ws *webServer) websocketHandler(wsock *websocket.Conn) {
	defer wsock.Close()

	log.Printf("%s [%s] \"%s %s\" connected.\n", wsock.Request().RemoteAddr, time.Now().Format(time.RFC3339), wsock.Request().Method, wsock.Request().RequestURI)

	var browserState BrowserStateMessage
	if err := websocket.JSON.Receive(wsock, &browserState); err != nil {
		log.Println("Error receiving JSON from WebSocket:", err)
		return
	}

	ingestUpdateNotification := make(chan struct{}, 1)
	ws.clientNotif.Register(ingestUpdateNotification)

	log.Printf("%+v\n", browserState)

	pingTicker := time.NewTicker(15 * time.Second)

	keep_looping := true
	for keep_looping {

		select {
		case <-pingTicker.C:
			{
				if err := websocket.Message.Send(wsock, "PING"); err != nil {
					log.Println("Failed to send PING to websocket:", err)
					keep_looping = false
					break
				}

				var expectedPong string
				if err := websocket.Message.Receive(wsock, &expectedPong); err != nil {
					log.Println("Failed to receive PONG on websocket:", err)
					keep_looping = false
					break
				}

				if expectedPong != "PONG" {
					log.Printf("Invalid PING response on websocket: %q\n", expectedPong)
					keep_looping = false
					break
				}
			}

		case <-ingestUpdateNotification:
			{
				log.Println("Websocket - Ingest update notification received")
				q := db.NewSelectMessageQuery().WithStartingId(browserState.LastLogMessageId).WithPattern(browserState.QueryString).Build()
				logs := db.Query(q)

				if 0 < len(*logs) {
					err := json.NewEncoder(wsock).Encode(logs)
					if err != nil {
						log.Printf("Error writing to websocket of %s: %v\n", wsock.RemoteAddr().String(), err)
						keep_looping = false
						break
					}

					browserState.LastLogMessageId = lastMessageId(browserState.LastLogMessageId, logs)
				}
			}
		}

		time.Sleep(1 * time.Second)
	}

	log.Printf("%s [%s] \"%s %s\" connection closed.\n", wsock.Request().RemoteAddr, time.Now().Format(time.RFC3339), wsock.Request().Method, wsock.Request().RequestURI)
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
