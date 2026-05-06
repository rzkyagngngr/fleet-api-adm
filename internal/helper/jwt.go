package helper

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims holds the JWT payload for every authenticated session.
// Includes branch and terminal identity so downstream services
// can enforce multi-tenant isolation without an extra DB round-trip.
type Claims struct {
	UserID       uint64 `json:"user_id"`
	Email        string `json:"email"`
	EmployeeID   string `json:"employee_id"`
	FullName     string `json:"full_name"`
	BranchCode   string `json:"branch_code"`
	BranchName   string `json:"branch_name"`
	TerminalCode string `json:"terminal_code"`
	TerminalName string `json:"terminal_name"`
	CompanyCode  string `json:"company_code"`
	CompanyName  string `json:"company_name"`
	jwt.RegisteredClaims
}

type JWTUtil struct {
	secret      []byte
	expiryHours int
}

func NewJWTUtil(secret string, expiryHours int) *JWTUtil {
	return &JWTUtil{secret: []byte(secret), expiryHours: expiryHours}
}

// GenerateToken issues a signed JWT containing the full user identity context.
func (j *JWTUtil) GenerateToken(
	userID uint64,
	email, employeeID, fullName string,
	branchCode, branchName string,
	terminalCode, terminalName string,
	companyCode, companyName string,
) (string, error) {
	claims := Claims{
		UserID:       userID,
		Email:        email,
		EmployeeID:   employeeID,
		FullName:     fullName,
		BranchCode:   branchCode,
		BranchName:   branchName,
		TerminalCode: terminalCode,
		TerminalName: terminalName,
		CompanyCode:  companyCode,
		CompanyName:  companyName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ValidateToken parses and validates a JWT string, returning the embedded Claims.
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
