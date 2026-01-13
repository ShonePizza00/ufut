package sqliteRepoShowcase

import (
	"context"
	"errors"
	structsUFUT "ufut/lib/structs"
)

var (
	ErrSoldOut = errors.New("item is sold out")
)

/*
req: None
resp: all categories at showcase
*/
func (r *SQLiteRepo) Categories(ctx context.Context) ([]string, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT categoryName FROM showcase_categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

/*
req:

	Category: always (specifies which category to search in)
	Price: TODO
	StartIndex: optional, 0 if not provided (offset from begging)
	Count: optional, 10 if not provided (number of items in response)
	OrderBy: "asc" or "desc". optional, "desc" if not provided. (specifies order)

resp:

	ItemsID: array of <string>ItemID
*/
func (r *SQLiteRepo) ItemsByParams(ctx context.Context, req *structsUFUT.ItemsRequestRSC) (structsUFUT.ItemsResponseRSC, error) {
	var resp structsUFUT.ItemsResponseRSC
	query := `SELECT itemID FROM showcase_items WHERE category=? LIMIT ? OFFSET ? ORDER BY `
	if req.OrderBy == "asc" {
		query += "price ASC"
	} else {
		query += "price DESC"
	}
	rows, err := r.DB.QueryContext(ctx, query, req.Category, req.Count, req.StartIndex)
	if err != nil {
		return resp, err
	}
	defer rows.Close()
	for rows.Next() {
		var itemID string
		if err := rows.Scan(&itemID); err != nil {
			return resp, err
		}
		resp.ItemsIDs = append(resp.ItemsIDs, itemID)
	}
	return resp, nil
}

/*
req:

	ItemID:			always (identifies the exact item)
	SellerID:		ignored
	Name:			ignored
	Description:		ignored
	Price:			ignored
	Category:		ignored
	Status:			ignored
	Quantity:		ignored

resp:

	ItemID:			not changed
	SellerID:		item's "SellerID"
	Name:			item's "Name"
	Description:		item's "Description"
	Price:			item's "Price"
	Category:		item's "Category"
	Status:			item's "Status"
	Quantity:		item's "Quantity"
*/
func (r *SQLiteRepo) ItemByItemID(ctx context.Context, req *structsUFUT.ItemDataRSC) error {
	q_res := r.DB.QueryRowContext(ctx,
		`SELECT * FROM showcase_items WHERE itemID=? AND category=?`, req.ItemID, req.Category)
	return q_res.Scan(
		&req.ItemID,
		&req.SellerID,
		&req.Name,
		&req.Description,
		&req.Price,
		&req.Category,
		&req.Status,
		&req.Quantity)
}

func (r *SQLiteRepo) ReserveItem(ctx context.Context, itemsID []string) []bool {
	r.mutexReserveItem.Lock()
	defer r.mutexReserveItem.Unlock()
	successful := make([]bool, len(itemsID))
	for i, itemID := range itemsID {
		q_res := r.DB.QueryRowContext(ctx,
			`SELECT quantity, status FROM showcase_items WHERE itemID=?`, itemID)
		var quantity int64
		var status string
		if err := q_res.Scan(&quantity, &status); err != nil {
			continue
		}
		if quantity <= 0 || status != "available" {
			continue
		}
		_, err := r.DB.ExecContext(ctx,
			`UPDATE showcase_items SET quantity=? WHERE itemID=?`, quantity-1, itemID)
		if err != nil {
			continue
		}
		successful[i] = true
	}
	return successful
}

func (r *SQLiteRepo) CancelItemReservation(ctx context.Context, itemsID []string) error {
	r.mutexReserveItem.Lock()
	defer r.mutexReserveItem.Unlock()
	for _, itemID := range itemsID {
		r.DB.ExecContext(ctx,
			`UPDATE showcase_items SET quantity = quantity+1 WHERE itemID=?`, itemID)
	}
	return nil
}
