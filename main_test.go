package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/lalathealter/dada/models"
	"github.com/stretchr/testify/assert"
)

type PostRegistrationCase struct {
  ExpStatus int
  User models.User
}



func init() {
  err := godotenv.Load(".test.env")
  if err != nil {
    log.Fatal("failed to initiate a database for testing - could not find an .env file with proper settings")
  }
}

var postregcases = []PostRegistrationCase{
  {http.StatusCreated, models.User{Name:"name", Password:"pass1"}},
  {http.StatusCreated, models.User{Name:"nametwo", Password:"pass2"}},
  {http.StatusBadRequest, models.User{Name:"name3", Password:"pass3"}},
  {http.StatusBadRequest, models.User{Name:"3123", Password:"pass22"}},
  {http.StatusBadRequest, models.User{Name:"333qweq", Password:"pass44"}},
  {http.StatusBadRequest, models.User{Name:"name w i", Password:"ppp"}},
  {http.StatusBadRequest, models.User{Name:"name_weqwe", Password:"pppp"}},
  {http.StatusBadRequest, models.User{Name:"", Password:""}},
  {http.StatusBadRequest, models.User{Name:"3", Password:""}},
  {http.StatusBadRequest, models.User{Name:"", Password:"password"}},
  {http.StatusBadRequest, models.User{Name:"validname", Password:""}},
  {http.StatusBadRequest, models.User{Name:"name", Password:"oops,thatwasnotunique"}},
  {http.StatusBadRequest, models.User{Name:"____", Password:"eqweq"}},
  {http.StatusCreated, models.User{Name:"dan", Password:"dana"}},
  {http.StatusCreated, models.User{Name:"daAAAA", Password:"validana"}},
  {http.StatusBadRequest, models.User{Name:"da", Password:"tooshort"}},
  {http.StatusBadRequest, models.User{Name:"d", Password:"toooooshort"}},
  {http.StatusBadRequest, models.User{Name:"TOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOLOOOOOOOOOOOOOOOOOOOOOOOONG", Password:"toolong"}},
  {http.StatusBadRequest, models.User{Name:"THISNAMEISALMOSTVALIDBUTTOOLONGEH", Password:"oops"}},
  {http.StatusCreated, models.User{Name:"ThisNameIsNicelyFitActuallyHeeey", Password:"yay"}},
  {http.StatusCreated, models.User{Name:"ThisShouldBeOkay", Password:"dhqwihei"}},
  {http.StatusBadRequest, models.User{Name:"####", Password:"meeeh"}},
  {http.StatusCreated, models.User{Name:"namethree", Password:"pass2"}},
}

func TestHandleRegistration(t *testing.T) {
  db, err := models.InitDB()
  if err != nil {
    t.Fatalf("Failed to properly init a database")
  }

  g := setupGin(db)

  for _, prcase := range postregcases {
    w := httptest.NewRecorder()
    var body bytes.Buffer
    err = json.NewEncoder(&body).Encode(prcase.User)
    if err != nil {
      t.Fatalf("Failed to encode a test case into JSON")
    }  
    req, _ := http.NewRequest("POST", "/register", &body)
    g.ServeHTTP(w, req)
    assert.Equal(t, prcase.ExpStatus, w.Code, prcase)
  }


  t.Cleanup(func() {
    res, err := db.Exec("DELETE FROM users")

    if err != nil  {
      t.Error(CleanupError{err})
    }
    n, err := res.RowsAffected()
    if err != nil {
      t.Error(CleanupError{err})
    }

    if n <= 0 {
      t.Error(CleanupError{ErrNoRowsDeleted})
    }
  })
}

var ErrNoRowsDeleted = errors.New("no rows were deleted")
type CleanupError struct {
  Err error
}
func (err CleanupError) Error() string {
  return "//FAILED CLEANUP: " + err.Error()
}
