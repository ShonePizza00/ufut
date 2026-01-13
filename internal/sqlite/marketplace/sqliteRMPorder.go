package sqliteRepoMarketplace

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	structsUFUT "ufut/lib/structs"
)

var (
	ErrOrderAlreadyFinished = errors.New("order is already finished")
)

/*
userID - uuid (16 bytes)
Takes data from shopping cart, places order and clears user's shopping cart
*/
func (r *SQLiteRepo) PlaceOrder(ctx context.Context, userID string, availability []bool) error {
	q_row_res := r.DB.QueryRowContext(ctx,
		`SELECT MAX(orderID)
		FROM usersOrders
		WHERE userID=?`, userID)
	var maxID int
	q_row_res.Scan(&maxID)
	maxID++
	q_res, err := r.DB.QueryContext(ctx,
		`SELECT itemID, quantity
		FROM shopping_cart
		WHERE userID=?`, userID)
	if err != nil {
		return err
	}
	defer q_res.Close()
	{
		_, err := r.DB.ExecContext(ctx,
			`INSERT INTO usersOrders
			(orderID, userID, status, createdAt)
			VALUES (?,?,"CREATED", CURRENT_TIMESTAMP)`,
			maxID, userID)
		if err != nil {
			return err
		}
	}
	i := -1
	for q_res.Next() {
		i++
		if !availability[i] {
			continue
		}
		var (
			ItemID   string
			Quantity string
		)
		q_res.Scan(&ItemID, &Quantity)
		orderID := userID + strconv.Itoa(maxID)
		_, errEx := r.DB.ExecContext(ctx,
			`INSERT INTO orders
			(orderID, itemID, quantity)
			VALUES (?,?,?)`,
			orderID, ItemID, Quantity)
		if errEx != nil {
			return errEx
		}
	}
	_, errDel := r.DB.ExecContext(ctx,
		`DELETE FROM shopping_cart WHERE userID=?`, userID)
	if errDel != nil {
		return errDel
	}
	return nil
}

/*
req:

	UserID	- must be not null
	OrderID	- must be not null
	Status	- ignored

F checks whether the order is finished; if it isn’t, F cancels it
*/
func (r *SQLiteRepo) RemoveOrder(ctx context.Context, req *structsUFUT.OrderRequestRMP) error {
	q_res := r.DB.QueryRowContext(ctx,
		`SELECT status FROM usersOrders WHERE orderID=? AND userID=?`, req.OrderID, req.UserID)
	var stts string
	if err := q_res.Scan(&stts); err != nil {
		return err
	}
	if stts != "FINISHED" {
		_, err := r.DB.ExecContext(ctx,
			`UPDATE usersOrders
			SET status="CANCELLED"
			WHERE orderID=? AND userID=?`, req.OrderID, req.UserID)
		return err
	} else {
		return ErrOrderAlreadyFinished
	}
}

/*
req:

	UserID	- must be not null
	OrderID	- must be not null
	Status	- will be filled with current status of the order

Checks the status of the order
*/
func (r *SQLiteRepo) OrderStatus(ctx context.Context, req *structsUFUT.OrderRequestRMP) error {
	q_res := r.DB.QueryRowContext(ctx,
		`SELECT status FROM usersOrders WHERE orderID=? AND userID=?`, req.OrderID, req.UserID)
	var stts string
	if err := q_res.Scan(&stts); err != nil {
		return err
	}
	req.Status = stts
	return nil
}

/*
req:

	UserID	- must be not null
	OrderID - ignored
	Status	- optional; if provided, only orders with this status will be returned

Returns all orders of the user, otherwise filtered by status if provided
*/
func (r *SQLiteRepo) UserOrders(ctx context.Context, req *structsUFUT.OrderRequestRMP) (*structsUFUT.OrdersResponseRMP, error) {
	query := `SELECT orderID, status FROM usersOrders WHERE userID=?`
	var q_res *sql.Rows
	var err error
	if req.Status != "" {
		query += ` AND status=?`
		q_res, err = r.DB.QueryContext(ctx, query, req.UserID, req.Status)
		if err != nil {
			return nil, err
		}
	} else {
		q_res, err = r.DB.QueryContext(ctx, query, req.UserID)
		if err != nil {
			return nil, err
		}
	}
	ret := structsUFUT.OrdersResponseRMP{}
	for q_res.Next() {
		var orderID int
		var status string
		q_res.Scan(&orderID, &status)
		ret.OrderID = append(ret.OrderID, orderID)
		ret.Status = append(ret.Status, status)
	}
	return &ret, nil
}

func (r *SQLiteRepo) ItemsIDsByOrderID(ctx context.Context, req *structsUFUT.OrderRequestRMP) ([]string, error) {
	res, err := r.DB.QueryContext(ctx,
		`SELECT itemID FROM orders WHERE orderID=?`, strconv.Itoa(req.OrderID)+req.UserID)
	if err != nil {
		return nil, err
	}
	items := make([]string, 0, 10)
	for res.Next() {
		var itemID string
		res.Scan(&itemID)
		items = append(items, itemID)
	}
	return items, nil
}
