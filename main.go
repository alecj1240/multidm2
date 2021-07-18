package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
  http.Handle("/", http.FileServer(http.Dir("./static")))
  http.HandleFunc("/messages/webhook", receiveMessage)
  http.HandleFunc("/user/redirect", appAuthorized)

  fmt.Printf("Starting server at port 8080\n")

  if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal(err)
  }
}

