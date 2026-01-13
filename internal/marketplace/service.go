package marketplace

import (
	"context"
	structsUFUT "ufut/lib/structs"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) PlaceOrder(ctx context.Context, userID string, callback func([]string) []bool) error {
	cart, err := s.ListCart(ctx, userID)
	if err != nil {
		return err
	}
	availability := callback(cart.ItemsID)
	return s.repo.PlaceOrder(ctx, userID, availability)
}

func (s *Service) RemoveOrder(ctx context.Context, req *structsUFUT.OrderRequestRMP, callback func([]string)) error {
	items, err := s.repo.ItemsIDsByOrderID(ctx, req)
	if err != nil {
		return err
	}
	callback(items)
	return s.repo.RemoveOrder(ctx, req)
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
