package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"
	"ufut/internal/catalog_service"
	sqliteRepoCatalog "ufut/internal/sqlite/catalog_service"
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
	repo := sqliteRepoCatalog.NewSQLiteRepo(db_)
	if err := repo.CreateTables(ctx); err != nil {
		log.Fatal(err)
	}
	service := catalog_service.NewService(repo)
	handler := catalog_service.NewHandler(service)
	catalog_service.RegisterRoutes(srvMx, handler)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
