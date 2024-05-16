package websocket

import (
	"fmt"
	"net/http"
	"time"
)

type Authentication interface {
	Authenticate(w http.ResponseWriter, r *http.Request) bool
	UserId(r *http.Request) string
}

type webSocketAuthentication struct{}

func (a *webSocketAuthentication) Authenticate(w http.ResponseWriter, r *http.Request) bool {
	return true
}

func (a *webSocketAuthentication) UserId(r *http.Request) string {
	query := r.URL.Query()
	if query != nil && query["userId"] != nil {
		return fmt.Sprintf("%v", query["userId"])
	}

	return fmt.Sprintf("%v", time.Now().UnixMilli())
}
