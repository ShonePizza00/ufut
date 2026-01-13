package showcase

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

func (s *Service) Categories(ctx context.Context) ([]string, error) {
	return s.repo.Categories(ctx)
}

func (s *Service) ItemsByParams(ctx context.Context, req *structsUFUT.ItemsRequestRSC) (structsUFUT.ItemsResponseRSC, error) {
	return s.repo.ItemsByParams(ctx, req)
}

func (s *Service) ItemByItemID(ctx context.Context, req *structsUFUT.ItemDataRSC) error {
	return s.repo.ItemByItemID(ctx, req)
}

func (s *Service) CreateItem(ctx context.Context, item *structsUFUT.ItemDataRSC) error {
	return s.repo.CreateItem(ctx, item)
}

func (s *Service) DeleteItem(ctx context.Context, item *structsUFUT.ItemDataRSC) error {
	return s.repo.DeleteItem(ctx, item)
}

func (s *Service) ReserveItem(ctx context.Context, itemID []string) []bool {
	return s.repo.ReserveItem(ctx, itemID)
}

func (s *Service) CancelItemReservation(ctx context.Context, itemID []string) error {
	return s.repo.CancelItemReservation(ctx, itemID)
}
