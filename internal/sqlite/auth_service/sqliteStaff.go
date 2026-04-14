package sqliteRepoAUTH

import (
	"context"
	"time"
	funcsUFUT "ufut/lib/funcs"
	structsUFUT "ufut/lib/structs"
)

/*
staff:

	Login		- uuid (16 bytes)
	PasswdHash	- uuid (16 bytes)

Creates new staff user in the database
*/
func (r *SQLiteRepo) CreateNewStaff(ctx context.Context, staff *structsUFUT.UserGeneral, tr *structsUFUT.TokenResponse) error {
	rtdie := time.Now().Unix() + funcsUFUT.RTtimeOffset
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO staff_auth 
		(userID, login, passwdHash, RTtoken, RTdie) 
		VALUES (?,?,?,?,?)`,
		staff.UserID, staff.Login, staff.PasswdHash, tr.RT, rtdie)
	return err
}

/*
data:

	Login			- uuid (16 bytes)
	PasswdHash		- current password hash, uuid (16 bytes)
	NewPasswdHash	- new password hash, uuid (16 bytes)

Updates staff user's password in the database
*/
func (r *SQLiteRepo) UpdateStaffPasswd(ctx context.Context, data *structsUFUT.UserUpdatePasswdHash, tr *structsUFUT.TokenResponse) (string, error) {
	q_res := r.DB.QueryRowContext(ctx, `SELECT passwdHash, userID FROM staff_auth WHERE login=?`, data.Login)
	var temp string
	var uid string
	err := q_res.Scan(&temp, &uid)
	if err != nil {
		return "", ErrNoUser
	}
	if temp != data.PasswdHash {
		return "", ErrIncorrectPasswd
	}
	rtdie := time.Now().Unix() + funcsUFUT.RTtimeOffset
	_, err = r.DB.ExecContext(ctx,
		`UPDATE staff_auth SET passwdHash=?, RTtoken=?, RTdie=? WHERE login=?`, data.NewPasswdHash, tr.RT, rtdie, data.Login)
	if err != nil {
		return "", err
	}
	return uid, nil
}

/*
user:

	Login		- uuid (16 bytes)
	PasswdHash	- sha256(password+hash_salt) (16 bytes)
	Token		- uuid+uuid (32 bytes)

Authenticates user and updates token in the database if successful
*/
func (r *SQLiteRepo) LoginStaff(ctx context.Context, staff *structsUFUT.UserGeneral, tr *structsUFUT.TokenResponse) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE staff_auth SET RTtoken=? WHERE login=? AND passwdHash=?`, tr.RT, staff.Login, staff.PasswdHash)
	if err != nil {
		return err
	}
	q := r.DB.QueryRowContext(ctx,
		`SELECT userID FROM staff_auth WHERE RTtoken=?`, tr.RT)
	err = q.Scan(&staff.UserID)
	return err
}

func (r *SQLiteRepo) UpdateJWTStaff(ctx context.Context, data *structsUFUT.JWTUpdate) (string, error) {
	var (
		userID    string
		timestamp int64
	)
	q := r.DB.QueryRowContext(ctx,
		`SELECT userID, RTdie FROM staff_auth WHERE RTtoken=?`, data.OldRT)
	if err := q.Scan(&userID, &timestamp); err != nil {
		return "", ErrNoUser
	}
	tNow := time.Now().Unix()
	if tNow > timestamp {
		return "", ErrTokenExpired
	}
	_, err := r.DB.ExecContext(ctx,
		`UPDATE staff_auth SET RTtoken=?, RTdie=? WHERE RTtoken=?`, data.NewRT, tNow+data.TimeOffset, data.OldRT)
	return userID, err
}
