package web

import (
	"encoding/json"
	"logbox/internal/db"
	"net/http"
)

func queryHandler(w http.ResponseWriter, r *http.Request) {

	// s(tart), c(ount), q(uery)

	q := r.URL.Query().Get("q")

	logs := db.Query(q)

	json.NewEncoder(w).Encode(logs)
}
