package models

import (
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
)

type User struct {
	Name     string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required"`
}

type UserCollectionI interface {
	SaveNewUser(User) error
  IsUnique(User) bool
  HasValidLogin(User) bool
}

type UserModel struct {
	DB *sql.DB
}


func (um UserModel) IsUnique(u User) bool {
  row := um.DB.QueryRow(`
    SELECT username
    FROM users
    WHERE name_index = $1
    `, produceNameIndex(u))
  var res string
  row.Scan(&res)
  return res == ""
}

func (um UserModel) SaveNewUser(u User) error {
  res, err := um.DB.Query(`
    INSERT INTO USERS(username, name_index, password)
    VALUES($1, $2, crypt($3, gen_salt('md5')))
    `, u.Name, produceNameIndex(u), u.Password)

  if err != nil {
    return err
  }
  defer res.Close()
  
  return res.Err()
}

func produceNameIndex(u User) string {
  return strings.ToLower(u.Name)
}

func (um UserModel) HasValidLogin(u User) (bool) {
  row := um.DB.QueryRow(`
      SELECT password
      FROM users
      WHERE username = $1 
      AND password = crypt($2, password)
    `, u.Name, u.Password)
  
  var pass string
  row.Scan(&pass)
  return pass != ""
}

