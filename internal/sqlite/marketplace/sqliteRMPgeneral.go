package sqliteRepoMarketplace

import (
	"context"
	"database/sql"
	structsUFUT "ufut/lib/structs"
)

type SQLiteRepo struct {
	DB *sql.DB
}

/*
Creates new SQLiteRepo instance
DB_instance - instance of sql.DB connected to SQLite database
Returns pointer to SQLiteRepo
*/
func NewSQLiteRepo(DB_instance *sql.DB) *SQLiteRepo {
	return &SQLiteRepo{DB: DB_instance}
}

/*
Closes the database connection
*/
func (r *SQLiteRepo) Close() error {
	return r.DB.Close()
}

/*
req:

	UserID  - uuid (16 bytes)
	ItemID  - uuid (16 bytes)
	Quantity - ignored

Returns quantity of the item in user's shopping cart
*/
func (r *SQLiteRepo) quantityOfItemInCart(ctx context.Context, req *structsUFUT.ItemRequestRMP) (int, error) {
	q_res := r.DB.QueryRowContext(ctx,
		`SELECT quantity FROM shopping_cart WHERE userID=? AND itemID=?`, req.UserID, req.ItemID)
	var t int
	err := q_res.Scan(&t)
	return t, err
}

/*
Creates necessary tables if they do not exist
*/
func (r *SQLiteRepo) CreateTables(ctx context.Context) error {
	{
		_, err := r.DB.ExecContext(ctx,
			`CREATE TABLE IF NOT EXISTS orders (
			orderID TEXT NOT NULL,
			itemID TEXT NOT NULL,
			quantity INTEGER,
			PRIMARY KEY(orderID, itemID)
			);`)
		if err != nil {
			return err
		}
	}
	{
		_, err := r.DB.ExecContext(ctx,
			`CREATE TABLE IF NOT EXISTS shopping_cart (
			userID TEXT NOT NULL,
			itemID TEXT NOT NULL,
			quantity INTEGER DEFAULT 1,
			PRIMARY KEY(userID, itemID)
			);`)
		if err != nil {
			return err
		}
	}
	{
		_, err := r.DB.ExecContext(ctx,
			`CREATE TABLE IF NOT EXISTS usersOrders (
			orderID INTEGER NOT NULL,
			userID TEXT NOT NULL,
			status TEXT NOT NULL,
			createdAt DATETIME NOT NULL,
			PRIMARY KEY(userID, orderID)
			);`)
		if err != nil {
			return err
		}
	}
	return nil
}
