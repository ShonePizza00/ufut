package sqliteRepoInventory

import (
	"context"
	"errors"
)

var (
	ErrInvalidValue  = errors.New("invalid value")
	ErrItemNotFound  = errors.New("itemID not found")
	ErrItemNotEnough = errors.New("item's quantity is not enough")
)

func (r *SQLiteRepo) ReserveItems(ctx context.Context, itemsIDs []string, quantities []int) ([]bool, error) {
	if len(itemsIDs) == 0 {
		return nil, ErrItemNotFound
	}
	availabilities := make([]bool, len(itemsIDs))
	for i, item := range itemsIDs {
		_, err := r.DecreaseItemQuantity(ctx, item, quantities[i])
		if err == nil {
			availabilities[i] = true
		}
	}
	return availabilities, nil
}

func (r *SQLiteRepo) CancelItemReservation(ctx context.Context, itemsIDs []string, quantities []int) error {
	if len(itemsIDs) == 0 {
		return ErrItemNotFound
	}
	for i, item := range itemsIDs {
		r.IncreaseItemQuantity(ctx, item, quantities[i])
	}
	return nil
}

func (r *SQLiteRepo) IncreaseItemQuantity(ctx context.Context, itemID string, n int) (int, error) {
	if n < 1 {
		return 0, ErrInvalidValue
	}
	res := r.DB.QueryRowContext(ctx, `
	UPDATE itemsQuantities
	SET quantity = quantity + ?
	WHERE itemID = ?
	RETURNING quantity`, n, itemID)
	var q int
	if err := res.Scan(&q); err != nil {
		return 0, ErrItemNotEnough
	}
	return q, nil
}

func (r *SQLiteRepo) DecreaseItemQuantity(ctx context.Context, itemID string, n int) (int, error) {
	if n < 1 {
		return 0, ErrInvalidValue
	}
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var cur int
	{
		res := r.DB.QueryRowContext(ctx, `
		SELECT itemID
		FROM itemsQuantities
		WHERE itemID = ?
		FOR UPDATE`, itemID)
		if err := res.Scan(&cur); err != nil {
			return 0, ErrItemNotFound
		}
	}
	var q int
	{
		res := r.DB.QueryRowContext(ctx, `
		UPDATE itemsQuantities
		SET quantity = quantity - ?
		WHERE itemID = ?
		AND quantity >= ?
		RETURNING quantity`, n, itemID, n)
		if err := res.Scan(&q); err != nil {
			return 0, ErrItemNotEnough
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return q, nil
}
