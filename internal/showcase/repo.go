package showcase

import (
	"context"
	structsUFUT "ufut/lib/structs"
)

type Repository interface {
	Categories(ctx context.Context) ([]string, error)
	ItemsByParams(ctx context.Context, req *structsUFUT.ItemsRequestRSC) (structsUFUT.ItemsResponseRSC, error)
	ItemByItemID(ctx context.Context, req *structsUFUT.ItemDataRSC) error
	ReserveItem(ctx context.Context, itemID []string) []bool
	CancelItemReservation(ctx context.Context, itemID []string) error

	CreateItem(ctx context.Context, item *structsUFUT.ItemDataRSC) error
	DeleteItem(ctx context.Context, item *structsUFUT.ItemDataRSC) error
}
