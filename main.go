package main

import (
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	configureDefaultLogger()
	logMultiline(title, " ")

	//
	// Initialise DB
	//
	db, err := newDB("./db/ipwatcher.db")
	if err != nil {
		log.Fatal(err)
	}

	//
	// Start public ip change polling
	//
	quit := make(chan bool)
	wg := new(sync.WaitGroup)
	wg.Add(1)

	ticker := time.NewTicker(1 * time.Hour)
	go update_ip_thread(db, ticker, quit, wg)

	//
	// Start webapi
	//
	r, err := newServer(db)
	if err != nil {
		log.Fatal(err)
	}
	r.Run(":8080") // listen and serve on port 8080

	wg.Wait()
}
