package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lalathealter/dada/controllers"
	"github.com/lalathealter/dada/controllers/auth"
	"github.com/lalathealter/dada/models"
)
 
func MainHandler(w http.ResponseWriter, r *http.Request) {
  hname, err := os.Hostname()
  if err != nil {
    panic(err)
  }
  io.WriteString(w, hname + time.Now().Format("2006-01-02 15:04:05"))
}

func setupGin(db *sql.DB) *gin.Engine {

  wrapper := controllers.InitWrapper(db)

  g := gin.Default()
  g.Use(controllers.ObligateToUseJSON)
  g.Use(controllers.HandleErrors)
  g.POST("/register", wrapper.HandleRegistration)
  g.POST("/login", wrapper.HandleLogin)
  g.Use(auth.ValidateJWT)
  g.GET("/me", wrapper.HandleViewSelf)
  g.PATCH("/me", wrapper.HandleUsernameChange)

  return g
}

func main() {
  err := godotenv.Load()
  if err != nil {
    fmt.Println("Failed to load local .env file ", err)
  }

  host := os.Getenv("host")
  port := os.Getenv("port")
  hp := net.JoinHostPort(host, port)

  db, err := models.InitDB()
  if err != nil {
    log.Fatal(err)
  }
  g := setupGin(db)

  fmt.Println("Listening on port ", port)
  g.Run(hp)
}

