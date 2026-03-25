package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID     uint64 `json:"user_id"`
	Email      string `json:"email"`
	EmployeeID string `json:"employee_id"`
	FullName   string `json:"full_name"`
	BranchCode *int64 `json:"branch_code"`
	jwt.RegisteredClaims
}

type JWTUtil struct {
	secret      []byte
	expiryHours int
}

func NewJWTUtil(secret string, expiryHours int) *JWTUtil {
	return &JWTUtil{
		secret:      []byte(secret),
		expiryHours: expiryHours,
	}
}

func (j *JWTUtil) GenerateToken(userID uint64, email string, employeeID string, fullName string, branchCode *int64) (string, error) {
	claims := Claims{
		UserID:     userID,
		Email:      email,
		EmployeeID: employeeID,
		FullName:   fullName,
		BranchCode: branchCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTUtil) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
