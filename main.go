// package main

// import "fmt"

// //Adding STRUCTURE

// type URL struct{
// 	ShortCode string
// 	LongURL string
// }

// // func add(a int, b int) int { //this func takes 2 ints and return an int
// // 	return a + b
// // }

// func greet(name string) string {
// 	return "Hello " + name
// }
// // func main(){

// // 	//Creating a map using the make in-built function
// // 	urlStore := make(map[string]string)

// // 	urlStore["abc123"] = "https://google.com"
// // 	urlStore["xyz789"] = "https://openai.com"

// // 	fmt.Println(urlStore["abc123"])

// // 	//intantiating a structure
// // 	u := URL{
// // 		ShortCode: "abc123",
// // 		LongURL: "https://google.com",
// // 	}
// // 	//adding a variable in the go file.
// // 	//There are two way - 1. var port(variable name) int = 2020 2. port := 2020
// // 	//the second method declare and assign using :=, second method can't be used for
// // 	//global initialization

// // 	fmt.Println(u.ShortCode)
// // 	fmt.Println(u.LongURL)

// // 	var message string = greet("URL Shortner")

// // 	port := 8080
// // 	fmt.Println("Hello, backend world \n", port);

// // 	fmt.Println(message)
// // }

package main

import (
	"math/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var urlStore = make(map[string]string)

func generateCode() string{
	//source of pool for random selection
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	//create a byte slice (array of byte) with length 6
	//byte is Go's type for a single character
	//this is to creat 6 empty memory slot for character 
	code := make([]byte, 6)


	// rand.Seed(time.Now().UnixNano()) this is an old version of manual seeding, in Go 1.20+ auto seeding is available
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	//this is to fill the code array with random int from range 0-61, which is length of letters
	for i := range code{
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

    longURL, ok := urlStore[code]
    if !ok {
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

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	code := generateCode()
	urlStore[code] = req.LongURL

	resp := ShortenResponse{
		ShortCode: code,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/shorten", shortenHandler)
    http.HandleFunc("/", redirectHandler)

	fmt.Println(("Server starting on port 8080..."))
	http.ListenAndServe(":8080", nil)

	// 	Test:
	// http://localhost:8080/
	// http://localhost:8080/shorten

}

// //why are we using struct for json?
// we use a struct:

// Type safety
// req.LongURL // guaranteed string
// No guessing, no parsing manually.

// Automatic mapping
// LongURL string `json:"long_url"`
// Go maps JSON → struct fields for you.

// Validation & clarity
// Missing fields → detectable
// Wrong types → errors
// Code is self-documenting

// Same struct used everywhere
// Request body
// Database rows (later)
// Redis values
// Responses


//yet to integrate the step 16 of the Phase II


