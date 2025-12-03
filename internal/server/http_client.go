package server

import (
	"net/http"
	"time"
)

func CreateHttpClient() *http.Client {
	client := http.Client{
		Timeout: time.Second * 50,
	}
	return &client
}
