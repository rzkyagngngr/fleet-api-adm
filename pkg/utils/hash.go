package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// HashSHA512 implements the hashing algorithm provided by the user (SHA-512 with uppercase hex)
func HashSHA512(password string) string {
	hasher := sha512.New()
	hasher.Write([]byte(password))
	hashBytes := hasher.Sum(nil)

	// In Java: sb.append(Integer.toString((bytes[i] & 0xff) + 0x100, 16).substring(1));
	// This is effectively hex.EncodeToString(hashBytes)
	hashHex := hex.EncodeToString(hashBytes)

	return strings.ToUpper(hashHex)
}

// HashPassword generates a bcrypt hash for the given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with a bcrypt hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
