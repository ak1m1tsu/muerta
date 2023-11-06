package entity

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  Password
	CreatedAt time.Time
	IsDeleted bool
}

type Users []User

type Password struct {
	Hash      []byte
	plainText string
}

func (p *Password) Generate(plainText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.plainText = plainText
	p.Hash = hash

	return nil
}

func (p *Password) Matches(plainText string) bool {
	return bcrypt.CompareHashAndPassword(p.Hash, []byte(plainText)) == nil
}
