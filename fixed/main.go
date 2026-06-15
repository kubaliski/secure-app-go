package main

import (
	"database/sql"
	"html"
	"log"
	"net/http"
	"os/exec"
	"regexp"

	_ "modernc.org/sqlite"
)

var db *sql.DB
var validHost = regexp.MustCompile(`^[a-zA-Z0-9.-]+$`)

func initDB() {
	var err error
	db, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("CREATE TABLE users (id INTEGER, username TEXT, password TEXT)")
	db.Exec("INSERT INTO users VALUES (1, 'admin', 'supersecret')")
	db.Exec("INSERT INTO users VALUES (2, 'user', 'password123')")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	row := db.QueryRow("SELECT id FROM users WHERE username = ? AND password = ?", username, password)
	var id int
	if err := row.Scan(&id); err != nil {
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}
	w.Write([]byte("Welcome!"))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		w.Write([]byte("<h1>Search</h1><form><input name=q><button>Search</button></form>"))
		return
	}
	w.Write([]byte("<h1>Search results for: " + html.EscapeString(q) + "</h1>"))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	if host == "" {
		w.Write([]byte("<h1>Ping</h1><form><input name=host value=localhost><button>Ping</button></form>"))
		return
	}
	if !validHost.MatchString(host) {
		http.Error(w, "Invalid host", http.StatusBadRequest)
		return
	}
	out, err := exec.Command("ping", "-c1", host).Output()
	if err != nil {
		http.Error(w, "Ping failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("<pre>" + string(out) + "</pre>"))
}

func main() {
	initDB()
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/ping", pingHandler)
	log.Println("Fixed app on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
