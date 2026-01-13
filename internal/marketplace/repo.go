package marketplace

import (
	"context"
	structsUFUT "ufut/lib/structs"
)

type Repository interface {
	PlaceOrder(ctx context.Context, userID string, availability []bool) error
	RemoveOrder(ctx context.Context, req *structsUFUT.OrderRequestRMP) error
	OrderStatus(ctx context.Context, req *structsUFUT.OrderRequestRMP) error
	UserOrders(ctx context.Context, req *structsUFUT.OrderRequestRMP) (*structsUFUT.OrdersResponseRMP, error)
	ItemsIDsByOrderID(ctx context.Context, req *structsUFUT.OrderRequestRMP) ([]string, error)

	AddToCart(ctx context.Context, req *structsUFUT.ItemRequestRMP) error
	RemoveFromCart(ctx context.Context, req *structsUFUT.ItemRequestRMP) error
	IncreaseItemQuantity(ctx context.Context, req *structsUFUT.ItemRequestRMP) error
	DecreaseItemQuantity(ctx context.Context, req *structsUFUT.ItemRequestRMP) error
	ListCart(ctx context.Context, userID string) (*structsUFUT.ShoppingCartRMP, error)
	ClearCart(ctx context.Context, UserID string) error
}
