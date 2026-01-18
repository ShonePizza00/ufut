package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	sqliteRepoAUTH "ufut/internal/sqlite/auth"
	funcsUFUT "ufut/lib/funcs"
	structsUFUT "ufut/lib/structs"

	"github.com/stretchr/testify/assert"

	"github.com/golang-jwt/jwt/v5"

	_ "github.com/mattn/go-sqlite3"
)

type cred struct {
	Login  string `json:"login"`
	Passwd string `json:"password"`
}
type wantS struct {
	errV bool
	RT   string
	code int
}
type testAuth struct {
	name  string
	value cred
	want  wantS
}

func CreateTestDataAuth() []testAuth {
	tests := []testAuth{
		{
			name: "abc+abc",
			value: cred{
				Login:  "abc",
				Passwd: "abc",
			},
			want: wantS{
				errV: false,
				RT:   "",
				code: 200,
			},
		},
		{
			name: "abc+incorrectPassword",
			value: cred{
				Login:  "abc",
				Passwd: "abcd",
			},
			want: wantS{
				errV: true,
				RT:   "",
				code: 401,
			},
		},
		{
			name: "empty",
			value: cred{
				Login:  "",
				Passwd: "",
			},
			want: wantS{
				errV: true,
				RT:   "",
				code: 400,
			},
		},
		{
			name: "empty+abc",
			value: cred{
				Login:  "",
				Passwd: "abc",
			},
			want: wantS{
				errV: true,
				RT:   "",
				code: 400,
			},
		},
		{
			name: "abc+empty",
			value: cred{
				Login:  "abc",
				Passwd: "",
			},
			want: wantS{
				errV: true,
				RT:   "",
				code: 400,
			},
		},
		{
			name: "a+a",
			value: cred{
				Login:  "a",
				Passwd: "a",
			},
			want: wantS{
				errV: false,
				RT:   "",
				code: 200,
			},
		},
		{
			name: "-+-",
			value: cred{
				Login:  "-",
				Passwd: "-",
			},
			want: wantS{
				errV: false,
				RT:   "",
				code: 200,
			},
		},
		{
			name: "0+0",
			value: cred{
				Login:  "0",
				Passwd: "0",
			},
			want: wantS{
				errV: false,
				RT:   "",
				code: 200,
			},
		},
	}
	return tests
}

func CreateAuthService(t *testing.T) (*Service, func() error) {
	dbFilePath := "usersAuth_test.db"
	db_Auth, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		t.Errorf("%v", err.Error())
	}
	repo_Auth := sqliteRepoAUTH.NewSQLiteRepo(db_Auth)
	if err := repo_Auth.CreateTables(t.Context()); err != nil {
		t.Fatalf("%v", err.Error())
	}
	return NewService(repo_Auth), func() error {
		err := db_Auth.Close()
		if err != nil {
			t.Fatalf("%v", err.Error())
			return err
		}
		os.Remove(dbFilePath)
		return nil
	}
}

func CheckJWT(t *testing.T, tk string) bool {
	token, err := jwt.ParseWithClaims(tk, &funcsUFUT.JWTCustomClaims{}, func(token *jwt.Token) (any, error) {
		return funcsUFUT.JWTkey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		return true
	}
	_, ok := token.Claims.(*funcsUFUT.JWTCustomClaims)
	if !ok {
		return true
	}
	return false
}

func TestHandler_AuthUser(t *testing.T) {
	srvc, cleanUp := CreateAuthService(t)
	defer cleanUp()
	h := NewHandler(srvc)
	tests := CreateTestDataAuth()
	for _, tt := range tests {
		body, _ := json.Marshal(tt.value)
		t.Run(tt.name+"_reg", func(t *testing.T) {
			r_reg := httptest.NewRequest(http.MethodPost, "/api/user/registerUser", bytes.NewReader(body))
			w_reg := httptest.NewRecorder()
			h.RegisterUser(w_reg, r_reg)
			res := w_reg.Result()
			var resJSON structsUFUT.TokenResponse
			json.NewDecoder(res.Body).Decode(&resJSON)
			res.Body.Close()
			tt.want.RT = resJSON.RT
		})
		t.Run(tt.name+"_auth", func(t *testing.T) {
			r_auth := httptest.NewRequest(http.MethodPost, "/api/user/authUser", bytes.NewReader(body))
			w_auth := httptest.NewRecorder()
			h.AuthUser(w_auth, r_auth)

			res := w_auth.Result()
			var resJSON structsUFUT.TokenResponse
			json.NewDecoder(res.Body).Decode(&resJSON)
			res.Body.Close()
			assert.Equal(t, tt.want.errV, CheckJWT(t, resJSON.JWT))
			assert.Equal(t, tt.want.code, res.StatusCode)
			// assert.Equal(t, tt.want.RT, resJSON.RT)
		})
	}
}

func TestHandler_AuthStaff(t *testing.T) {
	srvc, cleanUp := CreateAuthService(t)
	defer cleanUp()
	h := NewHandler(srvc)
	tests := CreateTestDataAuth()
	for _, tt := range tests {
		body, _ := json.Marshal(tt.value)
		t.Run(tt.name+"_reg", func(t *testing.T) {
			r_reg := httptest.NewRequest(http.MethodPost, "/api/staff/registerStaff", bytes.NewReader(body))
			w_reg := httptest.NewRecorder()
			h.RegisterUser(w_reg, r_reg)
			res := w_reg.Result()
			var resJSON structsUFUT.TokenResponse
			json.NewDecoder(res.Body).Decode(&resJSON)
			res.Body.Close()
			tt.want.RT = resJSON.RT
		})
		t.Run(tt.name+"_auth", func(t *testing.T) {
			r_auth := httptest.NewRequest(http.MethodPost, "/api/staff/authStaff", bytes.NewReader(body))
			w_auth := httptest.NewRecorder()
			h.AuthUser(w_auth, r_auth)

			res := w_auth.Result()
			var resJSON structsUFUT.TokenResponse
			json.NewDecoder(res.Body).Decode(&resJSON)
			res.Body.Close()
			assert.Equal(t, tt.want.errV, CheckJWT(t, resJSON.JWT))
			assert.Equal(t, tt.want.code, res.StatusCode)
			// assert.Equal(t, tt.want.RT, resJSON.RT)
		})
	}
}

func TestHandler_RegisterUser(t *testing.T) {
	srvc, cleanUp := CreateAuthService(t)
	defer cleanUp()
	h := NewHandler(srvc)
	tests := CreateTestDataAuth()
	for _, tt := range tests {
		body, _ := json.Marshal(tt.value)
		t.Run(tt.name+"_reg", func(t *testing.T) {
			r_reg := httptest.NewRequest(http.MethodPost, "/api/user/registerUser", bytes.NewReader(body))
			w_reg := httptest.NewRecorder()
			h.RegisterUser(w_reg, r_reg)
			res := w_reg.Result()
			var resJSON structsUFUT.TokenResponse
			json.NewDecoder(res.Body).Decode(&resJSON)
			res.Body.Close()
			tt.want.RT = resJSON.RT
			assert.Equal(t, tt.want.errV, CheckJWT(t, resJSON.JWT))
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func TestHandler_RegisterStaff(t *testing.T) {
	srvc, cleanUp := CreateAuthService(t)
	defer cleanUp()
	h := NewHandler(srvc)
	tests := CreateTestDataAuth()
	for _, tt := range tests {
		body, _ := json.Marshal(tt.value)
		t.Run(tt.name+"_reg", func(t *testing.T) {
			r_reg := httptest.NewRequest(http.MethodPost, "/api/staff/registerStaff", bytes.NewReader(body))
			w_reg := httptest.NewRecorder()
			h.RegisterUser(w_reg, r_reg)
			res := w_reg.Result()
			var resJSON structsUFUT.TokenResponse
			json.NewDecoder(res.Body).Decode(&resJSON)
			res.Body.Close()
			tt.want.RT = resJSON.RT
			assert.Equal(t, tt.want.errV, CheckJWT(t, resJSON.JWT))
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}
