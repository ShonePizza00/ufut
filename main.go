package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
	httpserver "ufut/internal/_httpserver"
	"ufut/internal/auth"
	"ufut/internal/marketplace"
	"ufut/internal/showcase"
	sqliteRepoAUTH "ufut/internal/sqlite/auth"
	sqliteRepoMarketplace "ufut/internal/sqlite/marketplace"
	sqliteRepoShowcase "ufut/internal/sqlite/showcase"

	_ "github.com/mattn/go-sqlite3"
)

var (
	addrString string = "127.0.0.1:80"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srvMx := http.NewServeMux()
	server := &http.Server{
		Addr:         addrString,
		Handler:      srvMx,
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 1000 * time.Millisecond,
	}
	db_Auth, err := sql.Open("sqlite3", "data/usersAuth.db")
	if err != nil {
		fmt.Println(err.Error())
	}
	db_MP, err := sql.Open("sqlite3", "data/marketplace.db")
	if err != nil {
		fmt.Println(err.Error())
	}
	db_SC, err := sql.Open("sqlite3", "data/showcase.db")
	if err != nil {
		fmt.Println(err.Error())
	}
	repo_Auth := sqliteRepoAUTH.NewSQLiteRepo(db_Auth)
	repo_MP := sqliteRepoMarketplace.NewSQLiteRepo(db_MP)
	repo_SC := sqliteRepoShowcase.NewSQLiteRepo(db_SC)
	if err := repo_Auth.CreateTables(ctx); err != nil {
		log.Fatalln(err)
	}
	if err := repo_MP.CreateTables(ctx); err != nil {
		log.Fatalln(err)
	}
	if err := repo_SC.CreateTables(ctx); err != nil {
		log.Fatalln(err)
	}
	service_Auth := auth.NewService(repo_Auth)
	service_MP := marketplace.NewService(repo_MP)
	service_SC := showcase.NewService(repo_SC)
	httpserver.AddRoutes(ctx, srvMx, &httpserver.Services{
		Service_Auth: service_Auth,
		Service_MP:   service_MP,
		Service_SC:   service_SC,
	})
	{
		err := server.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}
}
