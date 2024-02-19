package controllers

import (
	"database/sql"

	"github.com/lalathealter/dada/models"
)


func InitWrapper(db *sql.DB) *Wrapper {
	wr := Wrapper{
		models.UserModel{DB: db},
	}

	return &wr
}

type Wrapper struct {
	users UserCollectionI
}

type UserCollectionI interface {
	SaveNewUser(models.User) error
  IsUnique(models.User) bool
}

