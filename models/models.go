package models

import "database/sql"

type User struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type UserModel struct {
	DB *sql.DB
}

func (um UserModel) IsUnique(name string) bool {

  return false
}


func (um UserModel) SaveNewUser(u User) error {

  return nil
}
