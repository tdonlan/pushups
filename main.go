package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
	"time"
)

var db *sql.DB

type PushupHistory struct {
	Id        int
	Count     int
	Timestamp time.Time
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./db.sql")
	checkErr(err)

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	r.HandleFunc("/pushups/total", GetTotalHandler).Methods("GET")
	r.HandleFunc("/pushups/graph", GetPushupsGraph).Methods("GET")
	r.HandleFunc("/pushups", GetPushupsHandler).Methods("GET")
	r.HandleFunc("/pushups/{count}", AddPushupsHandler).Methods("POST")

	log.Println("Pushups! Listening at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
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

	GetTotalHandler(w, r)

}

func GetTotalHandler(w http.ResponseWriter, r *http.Request) {
	total := GetTotal()

	fmt.Fprintf(w, "%d", total)

}

func GetPushupsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pushups := GetPushups()

	pJson, err := json.Marshal(pushups)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%s", pJson)
}

func GetPushupsGraph(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pushups := GetPushups()

	graphData := map[string]int{}

	for _, v := range *pushups {
		graphData[v.Timestamp.Format("2006-01-02 15:04:05")] = v.Count
	}

	gJson, err := json.Marshal(graphData)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%s", gJson)
}

//-----------------

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

func GetPushups() *[]PushupHistory {
	pushups := []PushupHistory{}

	sql := `select * from pushups`
	rows, err := db.Query(sql)
	checkErr(err)

	for rows.Next() {
		var p PushupHistory
		err = rows.Scan(&p.Id, &p.Count, &p.Timestamp)
		checkErr(err)

		pushups = append(pushups, p)
	}

	return &pushups
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
