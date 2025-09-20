package models

import (
	"time"
)

type User struct {
	Id        string     `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Email     string     `db:"email" json:"email"`
	Password  string     `db:"password" json:"-"`
	ApiKey    *string    `db:"api_key" json:"-"`
	CreatedAt *time.Time `db:"created_at" json:"-"`
	UpdatedAt *time.Time `db:"updated_at" json:"-"`
}
