package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

var (
	oauth2Config = &oauth2.Config{
		ClientID:     "your_client_id",
		ClientSecret: "your_client_secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}
)

func main() {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		url := oauth2Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusFound)
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "未获取到code", http.StatusBadRequest)
			return
		}
		token, err := oauth2Config.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, "token获取失败", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "获取到的AccessToken: %s\n", token.AccessToken)
	})

	log.Println("OAuth2演示服务器启动于: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
