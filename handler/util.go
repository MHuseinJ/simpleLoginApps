package handler

import (
	"crypto/sha1"
	"fmt"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/golang-jwt/jwt"
	"os"
	"regexp"
	"time"
)

func CreateHash(password string) string {
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

func validateRegisterRequest(request generated.RegisterRequest) []string {
	var errors []string
	patternPhone := "^\\+62\\d{9,12}$"
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(request.Password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(request.Password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(request.Password)
	isIndonesia := regexp.MustCompile(patternPhone).MatchString(request.Phone)
	if !hasUpper {
		errors = append(errors, "Password must containing at least 1 capital characters")
	}
	if !hasSpecial {
		errors = append(errors, "Password must containing at least 1 special (non alpha-numeric) characters")
	}
	if !hasDigit {
		errors = append(errors, "Password must containing at least 1 number")
	}
	if !isIndonesia {
		errors = append(errors, "Phone numbers must start with the Indonesia country code “+62”")
	}

	if len(request.Fullname) < 3 || len(request.Fullname) > 60 {
		errors = append(errors, "Full name must be at minimum 3 characters and maximum 60 characters")
	}
	if len(request.Phone) < 12 || len(request.Phone) > 15 {
		errors = append(errors, "Phone numbers must be at minimum 10 characters and maximum 13 characters (+62 count as 1 char)")
	}
	if len(request.Password) < 6 || len(request.Password) > 64 {
		errors = append(errors, "Passwords must be minimum 6 characters and maximum 64 characters")
	}
	return errors
}

func verifyToken(tokenString string) error {
	var secretKey = []byte(os.Getenv("SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	fmt.Println(token.Raw)
	return nil
}
