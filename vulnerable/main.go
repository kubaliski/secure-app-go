package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	_ "modernc.org/sqlite"
)

var db *sql.DB

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
	query := fmt.Sprintf("SELECT id FROM users WHERE username = '%s' AND password = '%s'", username, password)
	row := db.QueryRow(query)
	var id int
	if err := row.Scan(&id); err != nil {
		fmt.Fprintf(w, "Login failed for user: %s", username)
		return
	}
	fmt.Fprintf(w, "Welcome user %d!", id)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		fmt.Fprint(w, "<h1>Search</h1><form><input name=q><button>Search</button></form>")
		return
	}
	fmt.Fprintf(w, "<h1>Search results for: %s</h1>", q)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	if host == "" {
		fmt.Fprint(w, "<h1>Ping</h1><form><input name=host value=localhost><button>Ping</button></form>")
		return
	}
	out, err := exec.Command("bash", "-c", "ping -c1 "+host).Output()
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "<pre>%s</pre>", out)
}

func main() {
	initDB()
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/ping", pingHandler)
	log.Println("Vulnerable app on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
