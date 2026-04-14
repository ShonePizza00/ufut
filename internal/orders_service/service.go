package orders_service

import (
	"context"
	"encoding/json"
	structsUFUT "ufut/lib/structs"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Service struct {
	repo        Repository
	kafkaWriter *kafka.Writer
}

func NewService(repo Repository, kafkaWriter *kafka.Writer) *Service {
	return &Service{
		repo:        repo,
		kafkaWriter: kafkaWriter}
}

func (s *Service) PlaceOrder(ctx context.Context, userID string) error {
	cart, err := s.ListCart(ctx, userID)
	if err != nil {
		return err
	}
	trx, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	cart.UserID = trx.String() + ".reserve"
	jsonData, err := json.Marshal(cart)
	if err != nil {
		return err
	}
	err = s.kafkaWriter.WriteMessages(ctx, kafka.Message{
		Value: jsonData,
	})
	if err != nil {
		return err
	}
	return nil
	// availability := callback(cart.ItemsID)
	// return s.repo.PlaceOrder(ctx, userID, availability)
}

func (s *Service) RemoveOrder(ctx context.Context, req *structsUFUT.OrderRequestRMP) error {
	items, err := s.repo.ItemsIDsByOrderID(ctx, req)
	if err != nil {
		return err
	}
	trx, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	items.UserID = trx.String() + ".cancelReservation"
	jsonData, err := json.Marshal(items)
	if err != nil {
		return err
	}
	err = s.kafkaWriter.WriteMessages(ctx, kafka.Message{
		Value: jsonData,
	})
	return nil
	// return s.repo.RemoveOrder(ctx, req)
}

func (s *Service) OrderStatus(ctx context.Context, req *structsUFUT.OrderRequestRMP) error {
	return s.repo.OrderStatus(ctx, req)
}

func (s *Service) UserOrders(ctx context.Context, req *structsUFUT.OrderRequestRMP) (*structsUFUT.OrdersResponseRMP, error) {
	return s.repo.UserOrders(ctx, req)
}

func (s *Service) AddToCart(ctx context.Context, req *structsUFUT.ItemRequestRMP) error {
	return s.repo.AddToCart(ctx, req)
}

func (s *Service) RemoveFromCart(ctx context.Context, req *structsUFUT.ItemRequestRMP) error {
	return s.repo.RemoveFromCart(ctx, req)
}

func (s *Service) IncreaseItemQuantity(ctx context.Context, req *structsUFUT.ItemRequestRMP) error {
	return s.repo.IncreaseItemQuantity(ctx, req)
}

func (s *Service) DecreaseItemQuantity(ctx context.Context, req *structsUFUT.ItemRequestRMP) error {
	return s.repo.DecreaseItemQuantity(ctx, req)
}

func (s *Service) ListCart(ctx context.Context, userID string) (*structsUFUT.ShoppingCartRMP, error) {
	return s.repo.ListCart(ctx, userID)
}

func (s *Service) ClearCart(ctx context.Context, userID string) error {
	return s.repo.ClearCart(ctx, userID)
}
