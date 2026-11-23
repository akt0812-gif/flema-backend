package main

import (
	"fmt"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
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
