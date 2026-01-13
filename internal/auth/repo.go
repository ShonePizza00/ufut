package auth

import (
	"context"
	structsUFUT "ufut/lib/structs"
)

type Repository interface {
	CreateNewUser(ctx context.Context, user *structsUFUT.UserGeneral, tr *structsUFUT.TokenResponse) error
	UpdateUserPasswd(ctx context.Context, data *structsUFUT.UserUpdatePasswdHash, tr *structsUFUT.TokenResponse) (string, error)
	LoginUser(ctx context.Context, user *structsUFUT.UserGeneral, tr *structsUFUT.TokenResponse) error
	// VerifyTokenUser(ctx context.Context, user *structsUFUT.UserGeneral) (string, error)
	UpdateJWTUser(ctx context.Context, data *structsUFUT.JWTUpdate) (string, error)

	CreateNewStaff(ctx context.Context, staff *structsUFUT.UserGeneral, tr *structsUFUT.TokenResponse) error
	UpdateStaffPasswd(ctx context.Context, data *structsUFUT.UserUpdatePasswdHash, tr *structsUFUT.TokenResponse) (string, error)
	LoginStaff(ctx context.Context, staff *structsUFUT.UserGeneral, tr *structsUFUT.TokenResponse) error
	// VerifyTokenStaff(ctx context.Context, staff *structsUFUT.UserGeneral) (string, error)
	UpdateJWTStaff(ctx context.Context, data *structsUFUT.JWTUpdate) (string, error)
}
