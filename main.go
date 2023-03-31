package main

import (
    "fmt"
    "log"
    "net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "<h1>Homepage</h1>")
}

func main() {
    http.HandleFunc("/", rootHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

