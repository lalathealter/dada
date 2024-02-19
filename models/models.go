package models

import (
	"database/sql"
	"strings"
	"time"

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
  GetInfo(User) UserData
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

type UserData struct {
  Username string `json:"username"`
  Id int `json:"id"`
  RegDate time.Time `json:"registration_date"`
  PassChange time.Time `json:"password_changed_at"`
}

func (um UserModel) GetInfo(u User) UserData{
  row := um.DB.QueryRow(`
    SELECT username, id, registration_date, password_changed_at
    FROM users 
    WHERE username = $1
    `, u.Name)

  var ud UserData
  row.Scan(&ud.Username, &ud.Id, &ud.RegDate, &ud.PassChange)
  return ud
}
