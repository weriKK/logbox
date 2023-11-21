package db

import (
	"database/sql"
	"fmt"
	"log"
	"logbox/internal/common"

	_ "modernc.org/sqlite"
)

const logfile string = "logs.db"

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite", logfile)
	if err != nil {
		log.Fatal(err)
	}
}

func Query(queryString string) *[]common.LogMessage {

	// TODO: normal query from website

	// TODO: TCP interface to stream messages to logbox
	// see oklog netcat streaming for inspiration

	// TODO Use Full Text search with highlights!
	// https://www.sqlitetutorial.net/sqlite-full-text-search/

	// TODO Filter from date

	// TODO json highlight/parse message on website

	// TODO stream/paginate as scrolling

	// TODO live updates as new logs enter the system

	q := "SELECT * FROM logs"

	if queryString != "" {
		q = fmt.Sprintf("SELECT * FROM logs WHERE message LIKE '%%%s%%'", queryString)
	}

	rows, err := db.Query(q)
	if err != nil {
		log.Println(err)
	}

	logs := make([]common.LogMessage, 0)

	for rows.Next() {
		l := common.LogMessage{}
		err = rows.Scan(&l.Id, &l.Message)
		if err != nil {
			log.Println(err)
		}

		logs = append(logs, l)
	}

	return &logs

}
