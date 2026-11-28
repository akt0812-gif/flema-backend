package main

import (
	"fmt"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// 許可するフロントのオリジン。本番では Cloud Run の環境変数で渡す
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:5173"
	}
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	fmt.Fprint(w, "Hello World")
}
func main() {
	// Cloud Run は環境変数 PORT でポート番号を渡すため、それを優先する。
	// ローカル実行時は 8080 を使う。
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)
}
