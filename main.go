package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
	"time"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./db.sql")
	checkErr(err)

	r := mux.NewRouter()
	r.HandleFunc("/pushups", GetTotalHandler).Methods("GET")
	r.HandleFunc("/pushups/{count}", AddPushupsHandler).Methods("POST")

	http.ListenAndServe(":8080", r)
}

func AddPushupsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	count := vars["count"]

	i, err := strconv.Atoi(count)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	AddPushups(i)

}

func GetTotalHandler(w http.ResponseWriter, r *http.Request) {
	total := GetTotal()

	fmt.Fprintf(w, "%d", total)

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
