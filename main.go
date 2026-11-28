package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *sql.DB

func main() {
	// 接続先は環境変数で切り替える(ローカル: tcp(localhost:3306)、Cloud Run: unix(/cloudsql/接続名))
	dsn := fmt.Sprintf("%s:%s@%s/%s?parseTime=true",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PWD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_DATABASE"),
	)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("fail: sql.Open, %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("fail: db.Ping, %v", err)
	}
	defer db.Close()

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/users", usersHandler)

	// Cloud Run は環境変数 PORT でポート番号を渡すため、それを優先する。
	// ローカル実行時は 8080 を使う。
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// 別のオリジン(フロント)からの fetch を許可する
func setCORS(w http.ResponseWriter) {
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:5173"
	}
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	fmt.Fprint(w, "Hello World")
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT id, name, age FROM user")
	if err != nil {
		log.Printf("fail: db.Query, %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			log.Printf("fail: rows.Scan, %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		log.Printf("fail: rows.Err, %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
