// Package auth provides functions for password hashing and verification.
package auth

import (
	"crypto/sha256"
	"fmt"
)

// GenerateHashFromPassword returns a SHA256 hash string for the given password and salt.
func GenerateHashFromPassword(password, salt string) string {
	passwdWithSalt := fmt.Sprintf("%s%s", password, salt)
	hash := sha256.Sum256([]byte(passwdWithSalt))
	return fmt.Sprintf("%x", hash)
}

// CompareHashAndPassword compares the given password and salt against the given hashed password.
// It returns true if they match, false otherwise.
func CompareHashAndPassword(password, salt, hashedPassword string) bool {
	return GenerateHashFromPassword(password, salt) == hashedPassword
}
