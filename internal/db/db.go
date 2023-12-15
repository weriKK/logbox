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

type SelectMessageQuery struct {
	table      string
	field      string
	pattern    string
	startingId int
}

func (q *SelectMessageQuery) defaultValues() {
	q.table = "logs"
	q.field = "message"
	q.pattern = ""
	q.startingId = 0
}

func NewSelectMessageQuery() *SelectMessageQuery {
	q := new(SelectMessageQuery)
	q.defaultValues()
	return q
}

func (q *SelectMessageQuery) WithPattern(p string) *SelectMessageQuery {
	q.pattern = p
	return q
}

func (q *SelectMessageQuery) WithStartingId(id int) *SelectMessageQuery {
	q.startingId = id
	return q
}

func (q *SelectMessageQuery) Build() string {

	conditions := make([]string, 0)

	if q.pattern != "" {
		cond := fmt.Sprintf("%s LIKE '%%%s%%'", q.field, q.pattern)
		conditions = append(conditions, cond)
	}

	if q.startingId > 0 {
		cond := fmt.Sprintf("id > %d", q.startingId)
		conditions = append(conditions, cond)
	}

	query := fmt.Sprintf("SELECT * FROM %s", q.table)
	if 0 < len(conditions) {
		query += " WHERE "
		query += strings.Join(conditions, " AND ")
	}

	log.Println(query)

	return query
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
		q = queryString
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
