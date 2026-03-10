package inventory_service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"
	structsUFUT "ufut/lib/structs"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

func NewRedisClient(ctx context.Context, cfg *structsUFUT.RedisConfig) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		Username:     cfg.User,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})
	if err := db.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return db, nil
}

type Service struct {
	Repo                     Repository
	kafkaOrdersReader        *kafka.Reader
	kafkaNotificationsWriter *kafka.Writer
	redisClient              *redis.Client
}

func NewService(
	repo Repository,
	kafkaReader1 *kafka.Reader,
	kafkaWriter1 *kafka.Writer,
	redisClient *redis.Client) *Service {
	return &Service{
		Repo:                     repo,
		kafkaOrdersReader:        kafkaReader1,
		kafkaNotificationsWriter: kafkaWriter1,
		redisClient:              redisClient,
	}
}

func (s *Service) handleOrdersMsg(ctx context.Context, msg kafka.Message) error {
	var list structsUFUT.ShoppingCartRMP
	if err := json.Unmarshal(msg.Value, &list); err != nil {
		return err
	}
	val, err := s.redisClient.Get(ctx, list.UserID).Result()
	if err != redis.Nil {
		return nil
	}
	action := strings.Split(val, ".")[1]
	var itemsAvailability []bool

	switch action {
	case "reserve":
		availability, err := s.Repo.ReserveItem(ctx, list.ItemsID, list.Quantities)
		if err != nil {
			return err
		}
		itemsAvailability = availability
	case "cancelReservation":
		err := s.Repo.CancelItemReservation(ctx, list.ItemsID)
		if err != nil {
			return err
		}
	}
	notificationMsg := structsUFUT.InventoryOrderNotification{
		Action:            action,
		ItemsAvailability: itemsAvailability,
		ItemsIDs:          list.ItemsID,
	}
	jsonData, err := json.Marshal(notificationMsg)
	if err != nil {
		log.Printf("%v%v\n", "failed marshal notification: ", err)
	} else {
		err = s.kafkaNotificationsWriter.WriteMessages(ctx, kafka.Message{
			Value: jsonData,
		})
		if err != nil {
			log.Println("failed send msg to kafka notification topic")
		}
	}
	if err := s.redisClient.Set(ctx, list.UserID, "", 10*time.Second); err != nil {
		log.Println("failed set trx: " + list.UserID)
	}
	return nil
}

func (s *Service) ServeKafka(ctx context.Context) error {
	for {
		msg, err := s.kafkaOrdersReader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return err
			}
			log.Printf("fetch error: %v\n", err)
			time.Sleep(time.Second)
			continue
		}

		if err := s.handleOrdersMsg(ctx, msg); err != nil {
			log.Printf("handle error: %v\n", err)
			continue
		}

		if err := s.kafkaOrdersReader.CommitMessages(ctx, msg); err != nil {
			log.Printf("commit error: %v\n", err)
		}
	}
}
