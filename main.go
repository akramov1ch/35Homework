package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book
var mutex sync.Mutex
const fileName = "books.json"

func main() {
	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/books/", bookHandler)
	loadBooks()
	log.Println("Server ishga tushdi 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadBooks() {
	mutex.Lock()
	defer mutex.Unlock()
	file, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			books = []Book{}
			return
		}
		log.Fatal(err)
	}
	err = json.Unmarshal(file, &books)
	if err != nil {
		log.Fatal(err)
	}
}

func saveBooks() {
	mutex.Lock()
	defer mutex.Unlock()
	file, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(fileName, file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
func booksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getBooks(w, r)
	case "POST":
		createBook(w, r)
	default:
		http.Error(w, "xato metod ishlatildi", http.StatusMethodNotAllowed)
	}
}




func bookHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/books/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "kitobning id si xato", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "GET":
		getBook(w, r, id)
	case "PUT":
		updateBook(w, r, id)
	case "DELETE":
		deleteBook(w, r, id)
	default:
		http.Error(w, "xato metod ishlatildi", http.StatusMethodNotAllowed)
	}
}
func getBooks(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}




func getBook(w http.ResponseWriter, r *http.Request, id int) {
	mutex.Lock()
	defer mutex.Unlock()
	for _, book := range books {
		if book.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.NotFound(w, r)
}
func createBook(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	book.ID = len(books) + 1
	books = append(books, book)
	saveBooks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}
func updateBook(w http.ResponseWriter, r *http.Request, id int) {
	mutex.Lock()
	defer mutex.Unlock()
	for i, book := range books {
		if book.ID == id {
			var updatedBook Book
			err := json.NewDecoder(r.Body).Decode(&updatedBook)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			updatedBook.ID = id // Ensure the ID remains unchanged
			books[i] = updatedBook
			saveBooks()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedBook)
			return
		}
	}
	http.NotFound(w, r)
}


func deleteBook(w http.ResponseWriter, r *http.Request, id int) {
	mutex.Lock()
	defer mutex.Unlock()
	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			saveBooks()
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}
