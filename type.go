package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	Token  string `json:"token"`
	Number int    `json:"number"`
}
type LogingRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type TransferRequest struct {
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}
type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	EncryptedPassword string    `json:"-"`
	Number            int64     `json:"number"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	encrytPW, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:                rand.Intn(10000),
		FirstName:         firstName,
		LastName:          lastName,
		EncryptedPassword: string(encrytPW),
		Number:            int64(rand.Intn(100000000)),
		CreatedAt:         time.Now().UTC(),
	}, err
}

func (a *Account) validatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(password)) == nil
}
