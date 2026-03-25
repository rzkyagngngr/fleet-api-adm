package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"strings"
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
