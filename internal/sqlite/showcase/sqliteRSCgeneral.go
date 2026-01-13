package sqliteRepoShowcase

import (
	"context"
	"database/sql"
	"sync"
)

type SQLiteRepo struct {
	DB               *sql.DB
	mutexCreateItem  sync.Mutex
	mutexReserveItem sync.Mutex
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
Creates necessary tables if they do not exist
*/
func (r *SQLiteRepo) CreateTables(ctx context.Context) error {
	{
		_, err := r.DB.ExecContext(ctx,
			`CREATE TABLE IF NOT EXISTS showcase_items (
			itemID TEXT PRIMARY KEY,
			sellerID TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			price INTEGER NOT NULL,
			category TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'available',
			quantity INTEGER NOT NULL DEFAULT 0,
			);`)
		if err != nil {
			return err
		}
	}
	{
		_, err := r.DB.ExecContext(ctx,
			`CREATE TABLE IF NOT EXISTS showcase_categories (
			categoryName TEXT PRIMARY KEY
			);`)
		if err != nil {
			return err
		}
	}
	{
		_, err := r.DB.ExecContext(ctx,
			`CREATE TABLE IF NOT EXISTS showcase_reviews (
			reviewID TEXT NOT NULL,
			itemID TEXT NOT NULL,
			userID TEXT NOT NULL,
			rating INTEGER NOT NULL,
			comment TEXT,
			PRIMARY KEY(reviewID, itemID)
			);`)
		if err != nil {
			return err
		}
	}
	{
		_, err := r.DB.ExecContext(ctx,
			`INSERT INTO showcase_categories (categoryName)
			VALUES ('books'), ('electronics'), ('clothing'), ('home'), ('toys'), ('18+')
			ON CONFLICT(categoryName) DO NOTHING;`)
		if err != nil {
			return err
		}
	}
	return nil
}
