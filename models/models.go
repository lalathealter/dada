package models

import (
	"database/sql"

  _ "github.com/lib/pq"
)

type User struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type UserModel struct {
	DB *sql.DB
}

func (um UserModel) IsUnique(u User) bool {
  row := um.DB.QueryRow(`
    SELECT username 
    FROM USERS
    WHERE username = $1
    `, u.Name)
  err := row.Err()
  return err == sql.ErrNoRows
}

func (um UserModel) SaveNewUser(u User) error {
  res, err := um.DB.Query(`
    INSERT INTO USERS(username, password)
    VALUES($1, crypt($2, gen_salt('md5')))
    `, u.Name, u.Password)

  if err != nil {
    return err
  }
  defer res.Close()
  
  return res.Err()
}

