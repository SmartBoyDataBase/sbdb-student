package handler

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8000/ping", nil)
	w := httptest.NewRecorder()
	PingPongHandler(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != "pong" {
		t.Error("PingPongTest Failed")
	}
}
