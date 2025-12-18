package sqliteRepoAUTH

import (
	"context"
	structsUFUT "ufut/lib"

	"github.com/google/uuid"
)

/*
staff:

	Login		- uuid (16 bytes)
	PasswdHash	- uuid (16 bytes)

Creates new staff user in the database
*/
func (r *SQLiteRepo) CreateNewStaff(ctx context.Context, staff *structsUFUT.UserGeneral) error {
	id := uuid.New()
	for {
		q_res := r.DB.QueryRowContext(ctx, `SELECT userID FROM staff_auth WHERE userID=?`, id)
		var t string
		if err := q_res.Scan(&t); err != nil {
			break
		}
		id = uuid.New()
	}
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO staff_auth 
		(userID, login, passwdHash, token) 
		VALUES (?,?,?,?)`,
		id, staff.Login, staff.PasswdHash, staff.Token)
	return err
}

/*
data:

	Login			- uuid (16 bytes)
	PasswdHash		- current password hash, uuid (16 bytes)
	NewPasswdHash	- new password hash, uuid (16 bytes)

Updates staff user's password in the database
*/
func (r *SQLiteRepo) UpdateStaffPasswd(ctx context.Context, data *structsUFUT.UserUpdatePasswdHash) error {
	q_res := r.DB.QueryRowContext(ctx, `SELECT passwdHash FROM staff_auth WHERE login=?`)
	var temp string
	err := q_res.Scan(&temp)
	if err != nil {
		return ErrNoUser
	}
	if temp != data.PasswdHash {
		return ErrIncorrectPasswd
	}
	_, err = r.DB.ExecContext(ctx, `
		UPDATE staff_auth SET passwdHash=?, token=? WHERE login=?`, data.NewPasswdHash, data.NewToken, data.Login)
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
func (r *SQLiteRepo) LoginStaff(ctx context.Context, staff *structsUFUT.UserGeneral) error {
	q_res := r.DB.QueryRowContext(ctx,
		`SELECT token FROM staff_auth WHERE login=? AND passwdHash=?`, staff.Login, staff.PasswdHash)
	var token string
	err := q_res.Scan(&token)
	if err != nil {
		return ErrNoUser
	}
	// r.updateTokenStaff(ctx, staff)
	return nil
}

/*
user:

	Login		- uuid (16 bytes)
	PasswdHash	- sha256(password+hash_salt) (16 bytes)
	Token		- uuid+uuid (32 bytes)

Updates user's token in the database
*/
func (r *SQLiteRepo) updateTokenStaff(ctx context.Context, staff *structsUFUT.UserGeneral) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE staff_auth SET token=? WHERE login=? AND passwdHash=?`, staff.Token, staff.Login, staff.PasswdHash)
	return err
}

/*
user:

	Login		- uuid (16 bytes)
	Token		- uuid+uuid (32 bytes)

Verifies user's token in the database
*/
func (r *SQLiteRepo) VerifyTokenStaff(ctx context.Context, staff *structsUFUT.UserGeneral) (string, error) {
	q_res := r.DB.QueryRowContext(ctx, `SELECT userID FROM staff_auth WHERE token=?`, staff.Token)
	var userID string
	err := q_res.Scan(&userID)
	if err != nil {
		return "", ErrNoUser
	}
	return userID, nil
}
