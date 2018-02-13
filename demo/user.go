package demo

import (
	"math/big"
)

// This example demonstrates
//   * indexes
//   * unexported fields

//go:generate sqlgen -type demo.User -o user_sql.go -gofmt -v -prefix Db -all -setters=all user.go role.go

type User struct {
	Uid          int64    `sql:"pk: true, auto: true"`
	Name         string   `sql:"unique: user_login"`
	EmailAddress string   `sql:"nk: true"`
	AddressId    *int64   `sql:"fk: addresses.id, onupdate: restrict, ondelete: restrict"`
	Avatar       *string
	Role         *Role    `sql:"type: text, size: 20"`
	Active       bool
	Admin        bool
	Fave         *big.Int `sql:"encode: json"`
	LastUpdated  int64

	// oauth token and secret
	token  string
	secret string

	// randomly generated hash used to sign user
	// session and application tokens.
	hash string `sql:"-"`
}

// These hooks are just used here for testing, but you could put whatever you like in them.
func (u *User) PreInsert() error {
	u.hash = "PreInsert"
	return nil
}

func (u *User) PreUpdate() error {
	u.hash = "PreUpdate"
	return nil
}

func (u *User) PostGet() error {
	u.hash = "PostGet"
	return nil
}
