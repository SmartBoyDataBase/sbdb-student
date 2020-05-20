package main

import (
	"log"
	"net/http"
	"sbdb-student/handler"
)

func main() {
	http.HandleFunc("/ping", handler.PingPongHandler)
	http.HandleFunc("/student", handler.Handler)
	http.HandleFunc("/students", handler.AllHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
