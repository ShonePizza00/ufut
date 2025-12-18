package sqliteRepoShowcase

import (
	"context"
	structsUFUT "ufut/lib"
)

func (r *SQLiteRepo) CreateItem(ctx context.Context, item *structsUFUT.ItemDataRSC) error {
	r.mutexCreateItem.Lock()
	defer r.mutexCreateItem.Unlock()
	var maxID int
	{
		q_res := r.DB.QueryRowContext(ctx,
			`SELECT MAX(itemID) FROM showcase_items WHERE category=?;`, item.Category)
		err := q_res.Scan(&maxID)
		if err != nil {
			return err
		}
	}
	item.ItemID = maxID + 1
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO showcase_items (itemID, sellerID, name, description, price, category)
		VALUES (?, ?, ?, ?, ?, ?);`,
		item.ItemID, item.SellerID, item.Name, item.Description, item.Price, item.Category)
	return err
}

func (r *SQLiteRepo) DeleteItem(ctx context.Context, item *structsUFUT.ItemDataRSC) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE showcase_items SET status='deleted' WHERE itemID=? AND category=?;`,
		item.ItemID, item.Category)
	return err
}
