package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type IPRecord struct {
	Addr    string
	Created string
}

type storage interface {
	latest() (*IPRecord, error)
	insert(string) error
	history() ([]IPRecord, error)
	close()
}

type sqliteDB struct {
	db *sql.DB
}

func newDB(path string) (storage, error) {
	db, err := sql.Open("sqlite3", path)

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS ipwatcher (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		ip TEXT NOT NULL,
		date_created DATATIME DEFAULT CURRENT_TIMESTAMP);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	s := sqliteDB{
		db: db,
	}
	return s, err
}

func (s sqliteDB) latest() (*IPRecord, error) {
	var last_ip, date_created string
	err := s.db.QueryRow("SELECT ip, date_created FROM ipwatcher ORDER BY id DESC LIMIT 0,1").Scan(&last_ip, &date_created)
	if err != nil {
		if err == sql.ErrNoRows { // http://go-database-sql.org/errors.html
			err = nil
		}
		return nil, err
	}

	result := &IPRecord{
		Addr:    last_ip,
		Created: date_created,
	}

	return result, nil
}

func (s sqliteDB) insert(ipAddr string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into ipwatcher(ip) values(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(ipAddr)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s sqliteDB) history() ([]IPRecord, error) {
	rows, err := s.db.Query("SELECT ip, date_created FROM ipwatcher ORDER BY id DESC LIMIT 0,50")
	if err != nil {
		if err == sql.ErrNoRows { // http://go-database-sql.org/errors.html
			err = nil
		}
		return nil, err
	}
	defer rows.Close()

	result := []IPRecord{}
	for rows.Next() {
		var last_ip, date_created string
		err := rows.Scan(&last_ip, &date_created)
		if err != nil {
			return nil, err
		}

		r := IPRecord{
			Addr:    last_ip,
			Created: date_created,
		}
		result = append(result, r)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s sqliteDB) close() {
	s.db.Close()
}
