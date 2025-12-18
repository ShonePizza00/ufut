package sqliteRepoAUTH

import (
	"context"
	"database/sql"
	"errors"
)

/*
Defines possible errors during user authentication and management
*/
var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrNoUser            = errors.New("user does not exists")
	ErrIncorrectPasswd   = errors.New("incorrect password")
	ErrIncorrectToken    = errors.New("incorrect token")
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
		CREATE TABLE IF NOT EXISTS users_auth (
		userID TEXT PRIMARY KEY,
		login TEXT NOT NULL UNIQUE,
		passwdHash TEXT NOT NULL,
		token TEXT NOT NULL
		);`)
		if err != nil {
			return err
		}
	}
	{
		_, err := r.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS staff_auth (
		userID TEXT PRIMARY KEY,
		login TEXT NOT NULL UNIQUE,
		passwdHash TEXT NOT NULL,
		token TEXT
		);`)
		if err != nil {
			return err
		}
	}
	return nil
}
