package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"
	"ufut/internal/orders_service"
	sqliteRepoOrders "ufut/internal/sqlite/orders_service"
	funcsUFUT "ufut/lib/funcs"

	_ "github.com/mattn/go-sqlite3"
	"github.com/segmentio/kafka-go"
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
	repo := sqliteRepoOrders.NewSQLiteRepo(db_)
	if err := repo.CreateTables(ctx); err != nil {
		log.Fatal(err)
	}
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{funcsUFUT.GetEnvDefault("KAFKA_ADDR", "localhost:9090")},
		Topic:   funcsUFUT.GetEnvDefault("KAFKA_ORDERS_TOPIC", "order_process"),
	})
	defer kafkaWriter.Close()
	service := orders_service.NewService(repo, kafkaWriter)
	handler := orders_service.NewHandler(service)
	orders_service.RegisterRoutes(srvMx, handler)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
