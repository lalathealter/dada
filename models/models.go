package models

import (
	"database/sql"
	"errors"
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
  IsUnique(string) bool
  HasValidLogin(User) bool
  GetInfo(User) UserData
  ChangeUsername(string, string) error
}

type UserModel struct {
	DB *sql.DB
}


func (um UserModel) IsUnique(name string) bool {
  row := um.DB.QueryRow(`
    SELECT username
    FROM users
    WHERE name_index = $1
    `, produceNameIndex(name))
  var res string
  row.Scan(&res)
  return res == ""
}

func (um UserModel) SaveNewUser(u User) error {
  res, err := um.DB.Query(`
    INSERT INTO USERS(username, name_index, password)
    VALUES($1, $2, crypt($3, gen_salt('md5')))
    `, u.Name, produceNameIndex(u.Name), u.Password)

  if err != nil {
    return err
  }
  defer res.Close()
  
  return res.Err()
}

func produceNameIndex(name string) string {
  return strings.ToLower(name)
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


type NewUsernamePackage struct {
  Name string `json:"new_username" binding:"required"`
}

var ErrFailedToChangeUsername = errors.New("Internal error; Failed to change the username")
func (um UserModel) ChangeUsername(prev, next string) error {
  res, err := um.DB.Exec(`
    UPDATE users
    SET username = $1,
      name_index = $2
    WHERE username = $3
    `, next, produceNameIndex(next), prev)

  if err != nil {
    return err
  }

  n, err := res.RowsAffected() 
  if err != nil {
    return err
  }
  if n < 0 {
    return ErrFailedToChangeUsername
  }

  return nil
}
