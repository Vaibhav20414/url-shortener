package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

var urlStore = make(map[string]string)

func generateCode() string {
	//source of pool for random selection
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	//create a byte slice (array of byte) with length 6
	//byte is Go's type for a single character
	//this is to creat 6 empty memory slot for character
	code := make([]byte, 6)

	// rand.Seed(time.Now().UnixNano()) this is an old version of manual seeding, in Go 1.20+ auto seeding is available
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	//this is to fill the code array with random int from range 0-61, which is length of letters
	for i := range code {
		code[i] = letters[rng.Intn(len(letters))]
	}

	return string(code)
}

type ShortenRequest struct {
	LongURL string `json:"long_url"`
}

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Server is running")
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:] // remove leading "/"

	var longURL string
	err := db.QueryRow("SELECT long_url FROM  urls WHERE short_code = $1", code).Scan(&longURL)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest

	er := json.NewDecoder(r.Body).Decode(&req)
	if er != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var (
		code string
		err  error
	)

	for i := 0; i < 5; i++ {
		code = generateCode()

		_, err = db.Exec("INSERT INTO urls (short_code, long_url) VALUES ($1, $2) ", code, req.LongURL)

		if err == nil {
			break
		}
	}

	if err != nil {
		http.Error(w, "Failed to store URL", http.StatusInternalServerError)
		return
	}

	resp := ShortenResponse{
		ShortCode: code,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	var err error

	connStr := "postgres://postgres:postgres@localhost:5432/url_shortener?sslmode=disable"

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to PostgreSQL")

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)

	fmt.Println(("Server starting on port 8080..."))
	http.ListenAndServe(":8080", nil)

	// 	Test:
	// http://localhost:8080/
	// http://localhost:8080/shorten

}
