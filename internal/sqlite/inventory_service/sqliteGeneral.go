package sqliteRepoInventory

import (
	"context"
	"database/sql"
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
Creates necessary tables if they do not exist
*/
func (r *SQLiteRepo) CreateTables(ctx context.Context) error {
	{
		_, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS itemsQuantities (
		itemID TEXT PRIMARY KEY,
		quantity INT
		);`)
		if err != nil {
			return err
		}
	}
	return nil
}
