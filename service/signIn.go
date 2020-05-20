package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func SignIn(username string, password string) (uint64, error) {
	body, _ := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: username,
		Password: password,
	})
	resp, err := http.Post(os.Getenv("SIGN_IN_URL"),
		"application/json",
		bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("sign in failed with errcode %d", resp.StatusCode)
	}
	var result struct {
		Id uint64 `json:"id"`
	}
	body, _ = ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &result)
	return result.Id, nil
}
