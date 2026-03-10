package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"
	auth_service "ufut/internal/auth_service"
	sqliteRepoAUTH "ufut/internal/sqlite/auth_service"
	funcsUFUT "ufut/lib/funcs"

	_ "github.com/mattn/go-sqlite3"
)

var (
	_PORT string = funcsUFUT.GetEnvDefault("PORT", "8080")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srvMx := http.NewServeMux()
	server := &http.Server{
		Addr:         ":" + _PORT,
		Handler:      srvMx,
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 1000 * time.Millisecond,
	}
	db_, err := sql.Open(
		funcsUFUT.GetEnvDefault("SQLDATABASE", "sqlite3"),
		funcsUFUT.GetEnvDefault("SQLCONNECT", "data.db"))
	if err != nil {
		log.Fatal(err)
	}
	defer db_.Close()
	repo := sqliteRepoAUTH.NewSQLiteRepo(db_)
	if err := repo.CreateTables(ctx); err != nil {
		log.Fatal(err)
	}
	service := auth_service.NewService(repo)
	handler := auth_service.NewHandler(service)
	auth_service.RegisterRoutes(srvMx, handler)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
