package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/AustinfHunter/blog/server/data"
	"github.com/golang-jwt/jwt/v5"
)

type authResponse struct {
	SignInError string `json:"signInError,omitempty"`
	Redirect    string `json:"redirect,omitempty"`
	Message     string `json:"message,omitempty"`
}

type claims struct {
	Admin bool `json:"admin"`
	jwt.RegisteredClaims
}

func getHashedPassword(u *data.User) string {
	s := os.Getenv("SECRET_KEY")
	fmt.Println(s)
	p := u.Email + u.Password + s
	h := hmac.New(sha256.New, []byte(p))
	return hex.EncodeToString(h.Sum(nil))
}

func comparePasswords(u *data.User, p2 string) bool {
	p1 := getHashedPassword(u)
	return p1 == p2
}

func getNewToken(u *data.User) *jwt.Token {
	claims := &claims{
		u.Admin,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    os.Getenv("HOST_NAME"),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

// ParseToken parses a token string and returns a pointer to the decoded token. Uses the custom claims defined in the jwtClaims struct.
// The token string must use the HMAC signing method.
func parseToken(tStr string) *jwt.Token {
	pt, err := jwt.ParseWithClaims(tStr, &claims{}, func(tk *jwt.Token) (interface{}, error) {
		if _, ok := tk.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tk.Header["alg"])
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	return pt
}
