package sqliteRepoAUTH

import (
	"context"
	structsUFUT "ufut/lib"

	"github.com/google/uuid"
)

/*
user:

	Login		- uuid (16 bytes)
	PasswdHash	- uuid (16 bytes)

Creates new user in the database
*/
func (r *SQLiteRepo) CreateNewUser(ctx context.Context, user *structsUFUT.UserGeneral) error {
	id := uuid.New()
	for {
		q_res := r.DB.QueryRowContext(ctx, `SELECT userID FROM users_auth WHERE userID=?`, id)
		var t string
		if err := q_res.Scan(&t); err == nil {
			break
		}
		id = uuid.New()
	}
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO users_auth
		(userID, login, passwdHash, token)
		VALUES (?,?,?,?)`,
		id, user.Login, user.PasswdHash, user.Token)
	return err
}

/*
data:

	Login			- uuid (16 bytes)
	PasswdHash		- current password hash, uuid (16 bytes)
	NewPasswdHash	- new password hash, uuid (16 bytes)

Updates user's password in the database
*/
func (r *SQLiteRepo) UpdateUserPasswd(ctx context.Context, data *structsUFUT.UserUpdatePasswdHash) error {
	q_res := r.DB.QueryRowContext(ctx, `SELECT passwdHash FROM users_auth WHERE login=?`)
	var temp string
	err := q_res.Scan(&temp)
	if err != nil {
		return ErrNoUser
	}
	if temp != data.PasswdHash {
		return ErrIncorrectPasswd
	}
	_, err = r.DB.ExecContext(ctx,
		`UPDATE users_auth SET passwdHash=?, token=? WHERE login=?`, data.NewPasswdHash, data.NewToken, data.Login)
	if err != nil {
		return err
	}
	return nil
}

/*
user:

	Login		- uuid (16 bytes)
	PasswdHash	- sha256(password+hash_salt) (16 bytes)
	Token		- uuid+uuid (32 bytes)

Authenticates user and updates token in the database if successful
*/
func (r *SQLiteRepo) LoginUser(ctx context.Context, user *structsUFUT.UserGeneral) error {
	q_res := r.DB.QueryRowContext(ctx,
		`SELECT token FROM users_auth WHERE login=? AND passwdHash=?`, user.Login, user.PasswdHash)
	err := q_res.Scan(&user.Token)
	if err != nil {
		return ErrNoUser
	}
	// r.updateTokenUser(ctx, user)
	return nil
}

/*
user:

	Login		- uuid (16 bytes)
	PasswdHash	- sha256(password+hash_salt) (16 bytes)
	Token		- uuid+uuid (32 bytes)

Updates user's token in the database
*/
func (r *SQLiteRepo) updateTokenUser(ctx context.Context, user *structsUFUT.UserGeneral) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE users_auth SET token=? WHERE login=? AND passwdHash=?`, user.Token, user.Login, user.PasswdHash)
	return err
}

/*
user:

	Login		- uuid (16 bytes)
	Token		- uuid+uuid (32 bytes)

Verifies user's token in the database
*/
func (r *SQLiteRepo) VerifyTokenUser(ctx context.Context, user *structsUFUT.UserGeneral) (string, error) {
	q_res := r.DB.QueryRowContext(ctx, `SELECT userID FROM users_auth WHERE token=?`, user.Token)
	var userID string
	err := q_res.Scan(&userID)
	if err != nil {
		return "", ErrNoUser
	}
	return userID, nil
}
