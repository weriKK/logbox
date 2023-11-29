package db

import (
	"database/sql"
	"fmt"
	"log"
	"logbox/internal/common"
	"strings"

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

	// unixepoch() was added in sqlite3 3.38.., ubuntu 22.04 only installs 3.37...
	// const create string = `CREATE TABLE IF NOT EXISTS logs (
	// 	id INTEGER NOT NULL PRIMARY KEY,
	// 	timestamp INTEGER DEFAULT (unixepoch('now','subsec')),
	// 	message TEXT
	// );`

	const create string = `CREATE TABLE IF NOT EXISTS logs (
		id INTEGER NOT NULL PRIMARY KEY,
		timestamp INTEGER DEFAULT (cast(strftime('%s', 'now') as int)),
		message TEXT
	);`

	if _, err := db.Exec(create); err != nil {
		log.Fatal(err)
	}
}

func Store(msg common.LogMessage) {

	m := strings.TrimSpace(msg.Message)
	log.Printf("%x", []byte(m))
	if msg.Message == "" {
		return
	}

	// TODO: timestamp
	insert := "INSERT INTO logs (message) VALUES(?)"
	if _, err := db.Exec(insert, msg.Message); err != nil {
		log.Println(err)
	}
}

func Query(queryString string) *[]common.LogMessage {

	// TODO: stream messages to logbox
	// What's the best/simplest way to get logs into logbox from files, folders, stdouts
	// see oklog netcat streaming for inspiration

	// TODO: Show last entry after query

	// TODO Use Full Text search with highlights!
	// https://www.sqlitetutorial.net/sqlite-full-text-search/

	// TODO Filter from date

	// TODO json highlight/parse message on website

	// TODO stream/paginate as scrolling

	// TODO live updates as new logs enter the system

	// TODO: show last entry during live updates

	// TODO: unshittyfy the code :)

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
		err = rows.Scan(&l.Id, &l.Timestamp, &l.Message)
		if err != nil {
			log.Println(err)
		}

		logs = append(logs, l)
	}

	return &logs

}
