package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	structsUFUT "ufut/lib"

	"github.com/google/uuid"
)

type Service struct {
	Repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{Repo: repo}
}

func generatePasswdHash(passwd string) string {
	passwdHashEncoder := sha256.New()
	passwdHashEncoder.Write([]byte(passwd))
	passwdHashEncoder.Write([]byte(structsUFUT.PASSWD_HASH_SALT))
	passwdHash := hex.EncodeToString(passwdHashEncoder.Sum(nil))
	return passwdHash
}

// authenticate user and return token
func (s *Service) AuthUser(ctx context.Context, req *structsUFUT.UpdatePasswdRequest) (string, error) {
	passwdHash := generatePasswdHash(req.Passwd)
	user := &structsUFUT.UserGeneral{Login: req.Login, PasswdHash: passwdHash}
	err := s.Repo.LoginUser(ctx, user)
	if err != nil {
		return "", err
	}
	return user.Token, nil
}

// authenticate staff and return token
func (s *Service) AuthStaff(ctx context.Context, req *structsUFUT.UpdatePasswdRequest) (string, error) {
	passwdHash := generatePasswdHash(req.Passwd)
	staff := &structsUFUT.UserGeneral{Login: req.Login, PasswdHash: passwdHash}
	err := s.Repo.LoginStaff(ctx, staff)
	if err != nil {
		return "", err
	}
	return staff.Token, nil
}

func (s *Service) RegisterUser(ctx context.Context, req *structsUFUT.UpdatePasswdRequest) (string, error) {
	passwdHash := generatePasswdHash(req.Passwd)
	token := uuid.New().String() + uuid.New().String()
	user := &structsUFUT.UserGeneral{Login: req.Login, PasswdHash: passwdHash, Token: token}
	err := s.Repo.CreateNewUser(ctx, user)
	if err != nil {
		return "", err
	}
	return user.Token, nil
}

func (s *Service) RegisterStaff(ctx context.Context, req *structsUFUT.UpdatePasswdRequest) (string, error) {
	passwdHash := generatePasswdHash(req.Passwd)
	token := uuid.New().String() + uuid.New().String()
	staff := &structsUFUT.UserGeneral{Login: req.Login, PasswdHash: passwdHash, Token: token}
	err := s.Repo.CreateNewStaff(ctx, staff)
	if err != nil {
		return "", err
	}
	return staff.Token, nil
}

func (s *Service) UpdateUserPasswd(ctx context.Context, req *structsUFUT.UserUpdatePasswd) (string, error) {
	nTk := uuid.New().String() + uuid.New().String()
	data := &structsUFUT.UserUpdatePasswdHash{
		Login:         req.Login,
		PasswdHash:    generatePasswdHash(req.Passwd),
		NewPasswdHash: generatePasswdHash(req.NewPasswd),
		NewToken:      nTk,
	}
	return nTk, s.Repo.UpdateUserPasswd(ctx, data)
}

func (s *Service) UpdateStaffPasswd(ctx context.Context, req *structsUFUT.UserUpdatePasswd) (string, error) {
	nTk := uuid.New().String() + uuid.New().String()
	data := &structsUFUT.UserUpdatePasswdHash{
		Login:         req.Login,
		PasswdHash:    generatePasswdHash(req.Passwd),
		NewPasswdHash: generatePasswdHash(req.NewPasswd),
		NewToken:      nTk,
	}
	return nTk, s.Repo.UpdateStaffPasswd(ctx, data)
}

func (s *Service) VerifyTokenUser(ctx context.Context, token string) (string, error) {
	user := &structsUFUT.UserGeneral{Token: token}
	return s.Repo.VerifyTokenUser(ctx, user)
}

func (s *Service) VerifyTokenStaff(ctx context.Context, token string) (string, error) {
	staff := &structsUFUT.UserGeneral{Token: token}
	return s.Repo.VerifyTokenStaff(ctx, staff)
}
