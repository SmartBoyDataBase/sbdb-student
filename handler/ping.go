package handler

import (
	"net/http"
)

func PingPongHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("pong"))
}
