package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"
	"ufut/internal/inventory_service"
	sqliteRepoInventory "ufut/internal/sqlite/inventory_service"
	funcsUFUT "ufut/lib/funcs"
	structsUFUT "ufut/lib/structs"

	// "ufut/internal/inventory_service"
	// sqliteRepoInventory "ufut/internal/sqlite/inventory_service"

	_ "github.com/mattn/go-sqlite3"
	"github.com/segmentio/kafka-go"
)

var (
	_PORT         string                  = funcsUFUT.GetEnvDefault("PORT", "8080")
	_REDIS_CONFIG structsUFUT.RedisConfig = structsUFUT.RedisConfig{
		Addr:        funcsUFUT.GetEnvDefault("REDIS_ADDR", "9091"),
		Password:    funcsUFUT.GetEnvDefault("REDIS_PASSWORD", ""),
		User:        funcsUFUT.GetEnvDefault("REDIS_USERNAME", ""),
		DB:          0,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}
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
		log.Panicln(err)
	}
	defer db_.Close()
	repo := sqliteRepoInventory.NewSQLiteRepo(db_)
	if err := repo.CreateTables(ctx); err != nil {
		log.Panicln(err)
	}
	kafkaOrdersReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{funcsUFUT.GetEnvDefault("KAFKA_ADDR", "localhost:9090")},
		Topic:   funcsUFUT.GetEnvDefault("KAFKA_ORDERS_TOPIC", "order_process"),
	})
	kafkaNotificationsWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{funcsUFUT.GetEnvDefault("KAFKA_ADDR", "localhost:9090")},
		Topic:   funcsUFUT.GetEnvDefault("KAFKA_NOTIFICATIONS_TOPIC", "notifications"),
	})
	redisClient, err := inventory_service.NewRedisClient(ctx, &_REDIS_CONFIG)
	if err != nil {
		log.Panicln(err)
	}
	service := inventory_service.NewService(repo, kafkaOrdersReader, kafkaNotificationsWriter, redisClient)
	handler := inventory_service.NewHandler(service)
	inventory_service.RegisterRoutes(srvMx, handler)
	go service.ServeKafka(ctx)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
