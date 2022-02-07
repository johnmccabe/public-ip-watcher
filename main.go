package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Print("starting service")
	db, err := sql.Open("sqlite3", "./db/ipwatcher.db")
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
		log.Fatalf("%q: %s\n", err, sqlStmt)
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
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("failed to get IP address, response code: %d", resp.StatusCode)
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return
	}
	cur_ip := strings.TrimSpace(string(b))

	if !validIP(cur_ip) {
		maxLen := 10
		if len(cur_ip) < maxLen {
			maxLen = len(cur_ip)
		}
		log.Printf("not a valid ip: %s", sanitise(cur_ip[0:maxLen-1]))
		return
	}

	//
	// Get last detected ip address
	//

	var last_ip, date_created string
	err = db.QueryRow("SELECT ip, date_created FROM ipwatcher ORDER BY id DESC LIMIT 0,1").Scan(&last_ip, &date_created)
	if err != nil {
		if err == sql.ErrNoRows { // http://go-database-sql.org/errors.html
			log.Print("no stored ip addresses")
		} else {
			log.Print(err)
			return
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
		log.Print(err)
		return
	}
	stmt, err := tx.Prepare("insert into ipwatcher(ip) values(?)")
	if err != nil {
		log.Print(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(cur_ip)
	if err != nil {
		log.Print(err)
		return
	}
	tx.Commit()

	log.Printf("updated public ip, old: %s (%s), new: %s", last_ip, date_created, cur_ip)
}

func validIP(ip string) bool {
	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	return re.MatchString(ip)
}

func sanitise(s string) string {
	// limit valid characters
	reg, _ := regexp.Compile("[^a-zA-Z0-9<>!#='\"()]+")
	return reg.ReplaceAllString(s, " ")
}
