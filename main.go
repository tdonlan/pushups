package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./db.sql")
	checkErr(err)

	AddPushups(5)
	AddPushups(10)

	total := GetTotal()

	fmt.Printf("Total pushups %d", total)
}

func AddPushups(count int) {
	sql := `insert into pushups (count, timestamp) values (?, datetime('now'));`
	stmt, err := db.Prepare(sql)
	checkErr(err)

	_, err = stmt.Exec(count)
	checkErr(err)
}

func GetTotal() int {
	total := 0
	sql := `select * from pushups`
	rows, err := db.Query(sql)
	checkErr(err)

	var id int
	var count int
	var datetime time.Time
	for rows.Next() {
		err = rows.Scan(&id, &count, &datetime)
		checkErr(err)
		total += count
	}

	return total
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
