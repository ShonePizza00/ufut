package auth

import (
	"context"
	structsUFUT "ufut/lib"
)

type Repository interface {
	CreateNewUser(ctx context.Context, user *structsUFUT.UserGeneral) error
	UpdateUserPasswd(ctx context.Context, data *structsUFUT.UserUpdatePasswdHash) error
	LoginUser(ctx context.Context, user *structsUFUT.UserGeneral) error
	VerifyTokenUser(ctx context.Context, user *structsUFUT.UserGeneral) (string, error)

	CreateNewStaff(ctx context.Context, staff *structsUFUT.UserGeneral) error
	UpdateStaffPasswd(ctx context.Context, data *structsUFUT.UserUpdatePasswdHash) error
	LoginStaff(ctx context.Context, staff *structsUFUT.UserGeneral) error
	VerifyTokenStaff(ctx context.Context, staff *structsUFUT.UserGeneral) (string, error)
}
