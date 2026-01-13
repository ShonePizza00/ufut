package funcsUFUT

import (
	"context"
	"crypto/rand"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	RTtimeOffset  = 604800
	JWTtimeOffset = 3600
)

var (
	JWTkey = []byte("123432")
)

type JWTCustomFields struct {
	GetterID string `json:"getterID"`
}

type JWTCustomClaims struct {
	JWTCustomFields
	jwt.RegisteredClaims
}

func GenerateJWT(fields JWTCustomFields) (string, error) {
	claims := &JWTCustomClaims{
		JWTCustomFields: fields,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(JWTtimeOffset * time.Second)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(JWTkey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func GenerateRT() (string, error) {
	rt := rand.Text()
	return rt, nil
}

func GetterIDFromContext(ctx context.Context) string {
	getterID, ok := ctx.Value("getterID").(string)
	if !ok {
		return ""
	}
	return getterID
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Invalid token1", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (any, error) {
			return JWTkey, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token2", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(*JWTCustomClaims)
		if !ok {
			http.Error(w, "Invalid token3", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "getterID", claims.GetterID)
		r = r.WithContext(ctx)
		next(w, r)
	}
}
