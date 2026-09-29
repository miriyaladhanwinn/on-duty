package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/miriyaladhanwinn/on-duty/internal/db"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
	JwtSecretKey    = []byte("amrita-cse-onduty-super-secure-production-jwt-key-2026")
)

type Claims struct {
	UserID     string      `json:"user_id"`
	Email      string      `json:"email"`
	Name       string      `json:"name"`
	RollNo     string      `json:"roll_no"`
	Department string      `json:"department"`
	Role       db.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(user *db.User, duration time.Duration) (string, error) {
	claims := Claims{
		UserID:     user.ID,
		Email:      user.Email,
		Name:       user.Name,
		RollNo:     user.RollNo,
		Department: user.Department,
		Role:       user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "amrita-onduty-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecretKey)
}

func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return JwtSecretKey, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
