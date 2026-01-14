package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	funcsUFUT "ufut/lib/funcs"
	structsUFUT "ufut/lib/structs"

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
func (s *Service) AuthUser(ctx context.Context, req *structsUFUT.UpdatePasswdRequest) (*structsUFUT.TokenResponse, error) {
	passwdHash := generatePasswdHash(req.Passwd)
	user := &structsUFUT.UserGeneral{Login: req.Login, PasswdHash: passwdHash}
	newRT, err := funcsUFUT.GenerateRT()
	if err != nil {
		return nil, err
	}
	tr := structsUFUT.TokenResponse{RT: newRT}
	{
		err := s.Repo.LoginUser(ctx, user, &tr)
		if err != nil {
			return nil, err
		}
	}
	tr.JWT, err = funcsUFUT.GenerateJWT(funcsUFUT.JWTCustomFields{GetterID: user.UserID})
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

// authenticate staff and return token
func (s *Service) AuthStaff(ctx context.Context, req *structsUFUT.UpdatePasswdRequest) (*structsUFUT.TokenResponse, error) {
	passwdHash := generatePasswdHash(req.Passwd)
	staff := &structsUFUT.UserGeneral{Login: req.Login, PasswdHash: passwdHash}
	newRT, err := funcsUFUT.GenerateRT()
	if err != nil {
		return nil, err
	}
	tr := structsUFUT.TokenResponse{RT: newRT}
	{
		err := s.Repo.LoginStaff(ctx, staff, &tr)
		if err != nil {
			return nil, err
		}
	}
	tr.JWT, err = funcsUFUT.GenerateJWT(funcsUFUT.JWTCustomFields{GetterID: staff.UserID})
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

func (s *Service) RegisterUser(ctx context.Context, req *structsUFUT.UpdatePasswdRequest) (*structsUFUT.TokenResponse, error) {
	passwdHash := generatePasswdHash(req.Passwd)
	uid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	suid := uid.String()
	user := &structsUFUT.UserGeneral{Login: req.Login, PasswdHash: passwdHash, UserID: suid}
	newJWT, err := funcsUFUT.GenerateJWT(funcsUFUT.JWTCustomFields{GetterID: suid})
	if err != nil {
		return nil, err
	}
	newRT, err := funcsUFUT.GenerateRT()
	if err != nil {
		return nil, err
	}
	tr := structsUFUT.TokenResponse{JWT: newJWT, RT: newRT}
	err = s.Repo.CreateNewUser(ctx, user, &tr)
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

func (s *Service) RegisterStaff(ctx context.Context, req *structsUFUT.UpdatePasswdRequest) (*structsUFUT.TokenResponse, error) {
	passwdHash := generatePasswdHash(req.Passwd)
	uid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	suid := uid.String()
	staff := &structsUFUT.UserGeneral{Login: req.Login, PasswdHash: passwdHash, UserID: suid}
	newJWT, err := funcsUFUT.GenerateJWT(funcsUFUT.JWTCustomFields{GetterID: suid})
	if err != nil {
		return nil, err
	}
	newRT, err := funcsUFUT.GenerateRT()
	if err != nil {
		return nil, err
	}
	tr := structsUFUT.TokenResponse{JWT: newJWT, RT: newRT}
	err = s.Repo.CreateNewStaff(ctx, staff, &tr)
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

func (s *Service) UpdateUserPasswd(ctx context.Context, req *structsUFUT.UserUpdatePasswd) (*structsUFUT.TokenResponse, error) {
	newRT, err := funcsUFUT.GenerateRT()
	if err != nil {
		return nil, err
	}
	tr := structsUFUT.TokenResponse{RT: newRT}
	data := &structsUFUT.UserUpdatePasswdHash{
		Login:         req.Login,
		PasswdHash:    generatePasswdHash(req.Passwd),
		NewPasswdHash: generatePasswdHash(req.NewPasswd),
	}
	uid, err := s.Repo.UpdateUserPasswd(ctx, data, &tr)
	if err != nil {
		return nil, err
	}
	tr.JWT, err = funcsUFUT.GenerateJWT(funcsUFUT.JWTCustomFields{GetterID: uid})
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

func (s *Service) UpdateStaffPasswd(ctx context.Context, req *structsUFUT.UserUpdatePasswd) (*structsUFUT.TokenResponse, error) {
	newRT, err := funcsUFUT.GenerateRT()
	if err != nil {
		return nil, err
	}
	tr := structsUFUT.TokenResponse{RT: newRT}
	data := &structsUFUT.UserUpdatePasswdHash{
		Login:         req.Login,
		PasswdHash:    generatePasswdHash(req.Passwd),
		NewPasswdHash: generatePasswdHash(req.NewPasswd),
	}
	uid, err := s.Repo.UpdateStaffPasswd(ctx, data, &tr)
	if err != nil {
		return nil, err
	}
	tr.JWT, err = funcsUFUT.GenerateJWT(funcsUFUT.JWTCustomFields{GetterID: uid})
	if err != nil {
		return nil, err
	}
	return &tr, nil
}

func (s *Service) UpdateJWTUser(ctx context.Context, rt string) (*structsUFUT.TokenResponse, error) {
	data := structsUFUT.JWTUpdate{
		OldRT:      rt,
		TimeOffset: funcsUFUT.RTtimeOffset,
	}
	var err error
	data.NewRT, err = funcsUFUT.GenerateRT()
	if err != nil {
		return nil, err
	}
	userID, err := s.Repo.UpdateJWTUser(ctx, &data)
	if err != nil {
		return nil, err
	}
	newJWT, err := funcsUFUT.GenerateJWT(funcsUFUT.JWTCustomFields{GetterID: userID})
	if err != nil {
		return nil, err
	}
	return &structsUFUT.TokenResponse{JWT: newJWT, RT: data.NewRT}, nil
}

func (s *Service) UpdateJWTStaff(ctx context.Context, rt string) (*structsUFUT.TokenResponse, error) {
	data := structsUFUT.JWTUpdate{
		OldRT:      rt,
		TimeOffset: funcsUFUT.RTtimeOffset,
	}
	var err error
	data.NewRT, err = funcsUFUT.GenerateRT()
	if err != nil {
		return nil, err
	}
	userID, err := s.Repo.UpdateJWTStaff(ctx, &data)
	if err != nil {
		return nil, err
	}
	newJWT, err := funcsUFUT.GenerateJWT(funcsUFUT.JWTCustomFields{GetterID: userID})
	if err != nil {
		return nil, err
	}
	return &structsUFUT.TokenResponse{JWT: newJWT, RT: data.NewRT}, nil
}
