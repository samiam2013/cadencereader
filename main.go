package main

import (
    "fmt"
    "net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
    for name, headers := range req.Header {
        for _, h := range headers {
            fmt.Fprintf(w, "%v: %v\n", name, h)
        }
    }
}

func main() {
    fmt.Println("binary started")
    http.HandleFunc("/hello", hello)
    http.HandleFunc("/headers", headers)
    fmt.Println("server starting")
    err := http.ListenAndServe(":80", nil)
    fmt.Println("Error on exit", err)
}
