package inventory_service

import "context"

type Repository interface {
	ReserveItem(ctx context.Context, itemsIDs []string, quantities []int) ([]bool, error)
	CancelItemReservation(ctx context.Context, itemsIDs []string) error
}
