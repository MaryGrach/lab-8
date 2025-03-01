package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "ps1"
	password = "1103"
	dbname   = "lr8"
)

type Handlers struct {
	dbProvider DatabaseProvider
}

type DatabaseProvider struct {
	db *sql.DB
}

func (h *Handlers) GetQuery(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("The parameter is not entered"))
		return
	}

	test, err := h.dbProvider.SelectQuery(name)
	if !test {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("The note has not been added to DB"))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello," + name + "!"))
}

func (h *Handlers) PostQuery(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("The parameter is not entered"))
		return
	}

	test, err := h.dbProvider.SelectQuery(name)
	if test && err == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("The note has already been added to DB"))
		return
	}

	err = h.dbProvider.InsertQuery(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Note added"))
}

func (dp *DatabaseProvider) SelectQuery(msg string) (bool, error) {
	var rec string

	row := dp.db.QueryRow("SELECT name_query FROM query WHERE name_query = ($1)", msg)
	err := row.Scan(&rec)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (dp *DatabaseProvider) InsertQuery(msg string) error {
	_, err := dp.db.Exec("INSERT INTO query (name_query) VALUES ($1)", msg)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	address := flag.String("address", "127.0.0.1:8083", "server startup address")
	flag.Parse()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dp := DatabaseProvider{db: db}
	h := Handlers{dbProvider: dp}

	http.HandleFunc("/get", h.GetQuery)
	http.HandleFunc("/post", h.PostQuery)

	err = http.ListenAndServe(*address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
