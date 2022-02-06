package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Print("starting service")
	db, err := sql.Open("sqlite3", "./ipwatcher.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Print("ensure db initialised")

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS ipwatcher (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		ip TEXT NOT NULL,
		date_created DATATIME DEFAULT CURRENT_TIMESTAMP);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	ticker := time.NewTicker(1 * time.Hour)
	quit := make(chan bool)

	update_ip(db) // run without waiting for first tick, so simple https://github.com/golang/go/issues/17601#issuecomment-307906597
	for {
		select {
		case <-ticker.C:
			update_ip(db)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func update_ip(db *sql.DB) {

	//
	// Query icanhazip.com for current ip address
	//

	resp, err := http.Get("https://icanhazip.com")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("failed to get IP address, response code: %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	cur_ip := strings.TrimSpace(string(b))

	//
	// Get last detected ip address
	//

	var last_ip, date_created string
	err = db.QueryRow("SELECT ip, date_created FROM ipwatcher ORDER BY id DESC LIMIT 0,1").Scan(&last_ip, &date_created)
	if err != nil {
		if err == sql.ErrNoRows { // http://go-database-sql.org/errors.html
			log.Print("no stored ip addresses")
		} else {
			log.Fatal(err)
		}
	}

	//
	// Check if ip address updated
	//

	if cur_ip == last_ip {
		log.Printf("no change, old: %s (%s)", last_ip, date_created)
		return
	}

	//
	// Store updated public ip address
	//

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into ipwatcher(ip) values(?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(cur_ip)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()

	log.Printf("updated public ip, old: %s (%s), new: %s", last_ip, date_created, cur_ip)
}
