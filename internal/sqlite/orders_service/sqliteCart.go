package sqliteRepoMarketplace

import (
	"context"
	structsUFUT "ufut/lib/structs"
)

/*
req:

	UserID  - uuid (16 bytes)
	ItemID  - uuid (16 bytes)
	Quantity - int, quantity to add

Adds item to user's shopping cart, or increases quantity if already present
*/
func (r *SQLiteRepo) AddToCart(ctx context.Context, req *structsUFUT.ItemRequestRMP) error {
	t, err := r.quantityOfItemInCart(ctx, req)
	if err != nil {
		_, err := r.DB.ExecContext(ctx,
			`INSERT INTO shopping_cart
			(userID, itemID, quantity)
			VALUES (?,?,?)`,
			req.UserID, req.ItemID, req.Quantity)
		if err != nil {
			return err
		}
	} else {
		t += req.Quantity
		_, err := r.DB.ExecContext(ctx,
			`UPDATE shopping_cart
			SET quantity=?
			WHERE userID=? AND itemID=?`,
			t, req.UserID, req.ItemID)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
req:

	UserID  - uuid (16 bytes)
	ItemID  - uuid (16 bytes)
	Quantity - ignored

Removes item from user's shopping cart
*/
func (r *SQLiteRepo) RemoveFromCart(ctx context.Context, req *structsUFUT.ItemRequestRMP) error {
	{
		_, err := r.quantityOfItemInCart(ctx, req)
		if err != nil {
			return err
		}
	}
	_, err := r.DB.ExecContext(ctx,
		`DELETE FROM shopping_cart
		WHERE userID=? AND itemID=?`,
		req.UserID, req.ItemID)
	if err != nil {
		return err
	}
	return nil
}

/*
req:

	UserID  - uuid (16 bytes)
	ItemID  - uuid (16 bytes)
	Quantity - int, quantity to increase

Increases quantity of the item in user's shopping cart
*/
func (r *SQLiteRepo) IncreaseItemQuantity(ctx context.Context, req *structsUFUT.ItemRequestRMP) error {
	t, err := r.quantityOfItemInCart(ctx, req)
	if err != nil {
		return err
	}
	t += req.Quantity
	_, err = r.DB.ExecContext(ctx,
		`UPDATE shopping_cart
		SET quantity=?
		WHERE userID=? AND itemID=?`,
		t, req.UserID, req.ItemID)
	if err != nil {
		return err
	}
	return nil
}

/*
req:

	UserID  - uuid (16 bytes)
	ItemID  - uuid (16 bytes)
	Quantity - int, quantity to decrease

Decreases quantity of the item in user's shopping cart
*/
func (r *SQLiteRepo) DecreaseItemQuantity(ctx context.Context, req *structsUFUT.ItemRequestRMP) error {
	t, err := r.quantityOfItemInCart(ctx, req)
	if err != nil {
		return err
	}
	t -= req.Quantity
	if t == 0 {
		_, err := r.DB.ExecContext(ctx,
			`DELETE FROM shopping_cart
			WHERE userID=? AND itemID=?`,
			req.UserID, req.ItemID)
		if err != nil {
			return err
		}
		return nil
	}
	{
		_, err := r.DB.ExecContext(ctx,
			`UPDATE shopping_cart
			SET quantity=?
			WHERE userID=? AND itemID=?`,
			t, req.UserID, req.ItemID)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
req:

	UserID  - uuid (16 bytes)

Returns contents of the user's shopping cart
*/
func (r *SQLiteRepo) ListCart(ctx context.Context, userID string) (*structsUFUT.ShoppingCartRMP, error) {
	q_res, err := r.DB.QueryContext(ctx,
		`SELECT itemID, quantity FROM shopping_cart WHERE userID=?`, userID)
	if err != nil {
		return nil, err
	}
	defer q_res.Close()
	var cart structsUFUT.ShoppingCartRMP
	cart.UserID = userID
	for q_res.Next() {
		var itemID string
		var quantity int
		err := q_res.Scan(&itemID, &quantity)
		if err != nil {
			return nil, err
		}
		cart.ItemsID = append(cart.ItemsID, itemID)
		cart.Quantities = append(cart.Quantities, quantity)
	}
	return &cart, nil
}

/*
req:

	UserID  - uuid (16 bytes)

Clears the user's shopping cart
*/
func (r *SQLiteRepo) ClearCart(ctx context.Context, UserID string) error {
	_, err := r.DB.ExecContext(ctx,
		`DELETE FROM shopping_cart WHERE userID=?`, UserID)
	return err
}
