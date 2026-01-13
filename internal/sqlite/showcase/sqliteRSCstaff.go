package sqliteRepoShowcase

import (
	"context"
	structsUFUT "ufut/lib/structs"
)

func (r *SQLiteRepo) CreateItem(ctx context.Context, item *structsUFUT.ItemDataRSC) error {
	// r.mutexCreateItem.Lock()
	// defer r.mutexCreateItem.Unlock()
	// var maxID string
	// {
	// 	q_res := r.DB.QueryRowContext(ctx,
	// 		`SELECT MAX(itemID) FROM showcase_items WHERE category=?;`, item.Category)
	// 	err := q_res.Scan(&maxID)
	// 	if err != nil {
	// 		maxID = 0
	// 	}
	// }
	// item.ItemID = maxID + 1
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO showcase_items (itemID, sellerID, name, description, price, category, status, quantity)
		VALUES (?, ?, ?, ?, ?,?,?,?);`,
		item.ItemID, item.SellerID, item.Name, item.Description, item.Price, item.Category, item.Status, item.Quantity)
	return err
}

func (r *SQLiteRepo) DeleteItem(ctx context.Context, item *structsUFUT.ItemDataRSC) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE showcase_items SET status='deleted' WHERE itemID=? AND category=?;`,
		item.ItemID, item.Category)
	return err
}

/*
update quantity
*/
