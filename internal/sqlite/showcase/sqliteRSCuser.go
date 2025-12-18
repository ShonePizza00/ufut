package sqliteRepoShowcase

import (
	"context"
	"strconv"
	structsUFUT "ufut/lib"
)

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
		var itemID int
		if err := rows.Scan(&itemID); err != nil {
			return resp, err
		}
		resp.ItemsIDs = append(resp.ItemsIDs, req.Category+":"+strconv.Itoa(itemID))
	}
	return resp, nil
}
