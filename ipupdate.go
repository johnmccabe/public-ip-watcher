package main

import (
	"log"
	"sync"
	"time"
)

func update_ip_thread(db storage, ticker *time.Ticker, quit <-chan bool, wg *sync.WaitGroup) {
	update_ip(db) // run without waiting for first tick, so simple https://github.com/golang/go/issues/17601#issuecomment-307906597
	for {
		select {
		case <-ticker.C:
			update_ip(db)
		case <-quit:
			ticker.Stop()
			wg.Done()
		}
	}
}

func update_ip(db storage) {
	//
	// Query icanhazip.com for current ip address
	//
	cur_ip, err := getIP()
	if err != nil {
		log.Print(err)
		return
	}

	//
	// Get last detected ip address
	//
	last_ip, err := db.latest()
	if err != nil {
		log.Print(err)
		return // skip DB update
	}

	//
	// Store first ip address
	//
	if last_ip == nil {
		if err := db.insert(cur_ip); err != nil {
			log.Fatal(err)
		}
		log.Printf("initialised public ip: %s", cur_ip)
		return
	}

	//
	// Check if ip address updated
	//
	if cur_ip == last_ip.Addr {
		log.Printf("no change, old: %s (%s)", last_ip.Addr, last_ip.Created)
		return
	}

	//
	// Store updated public ip address
	//
	if err := db.insert(cur_ip); err != nil {
		log.Fatal(err)
	}
	log.Printf("updated public ip, old: %s (%s), new: %s", last_ip.Addr, last_ip.Created, cur_ip)
}
