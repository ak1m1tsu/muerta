package auth

import (
	"crypto/sha256"
	"fmt"
)

func GenerateHashFromPassword(password, salt string) string {
	passwdWithSalt := fmt.Sprintf("%s%s", password, salt)
	res := sha256.Sum256([]byte(passwdWithSalt))
	return fmt.Sprintf("%x", res)
}

func CompareHashAndPassword(password, salt, hashedPassword string) bool {
	passwdWithSalt := fmt.Sprintf("%s%s", password, salt)
	hash := sha256.Sum256([]byte(passwdWithSalt))
	return fmt.Sprintf("%x", hash) == hashedPassword
}
