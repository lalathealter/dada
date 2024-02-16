package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

)
 
func MainHandler(w http.ResponseWriter, r *http.Request) {
  hname, err := os.Hostname()
  if err != nil {
    panic(err)
  }
  io.WriteString(w, hname + time.Now().Format("2006-01-02 15:04:05"))
}
 
func main() {
  http.HandleFunc("/", MainHandler)

  port := "8080"
  host := ""
  hp := net.JoinHostPort(host, port)

  fmt.Println("Listening on port ", port)
  http.ListenAndServe(hp, nil)
}

