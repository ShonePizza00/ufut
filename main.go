package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
	sqliteRepoAUTH "ufut/internal/sqlite/auth"

	_ "github.com/mattn/go-sqlite3"
)

var (
	addrString string = "127.0.0.1:80"
)

func main() {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	srvMx := http.NewServeMux()
	server := &http.Server{
		Addr:         addrString,
		Handler:      srvMx,
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 1000 * time.Millisecond,
	}
	db, err := sql.Open("sqlite3", "data/usersAuth.db")
	authDB := sqliteRepoAUTH.NewSQLiteRepo(db)

	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
