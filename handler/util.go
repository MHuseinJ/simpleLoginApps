package handler

import (
	"crypto/sha1"
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

func createHash(password string) string {
	var salt = os.Getenv("SALT")
	var sha = sha1.New()
	sha.Write([]byte(password + salt))
	var encrypted = sha.Sum(nil)
	return fmt.Sprintf("%x", encrypted)
}

func createToken(username string) (string, error) {
	var secretKey = []byte(os.Getenv("SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) error {
	var secretKey = os.Getenv("SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
