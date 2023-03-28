/*

Nama : Yusuf Valent Adyatomo
Kelas : Golang-03
ID Peserta : 1955617840-1275

*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"price"`
	Desc   string `json:"desc"`
}

type BookRequest struct {
	Title  string `json:"judul"`
	Author string `json:"author"`
	Desc   string `json:"deskripsi"`
}

var Books = []Book{
	{
		Id:     1,
		Title:  "Golang",
		Author: "Gopher",
		Desc:   "A book for Go",
	},
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})

	http.HandleFunc("/book", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			if r.URL.Query().Get("id") != "" {
				getBooksById(w, r)
				return
			}
			getBooks(w, r)
			return
		}

		// if r.Method == "GET" && r.URL.Query().Get("id") != "" {
		// 	getBooksById(w, r)
		// 	return
		// }

		if r.Method == "POST" {
			createBooks(w, r)
			return
		}

		if r.Method == "PUT" {
			updateBooks(w, r)
			return
		}

		if r.Method == "DELETE" {
			deleteBooks(w, r)
			return
		}

	})

	fmt.Println("Listening on PORT :8000")

	http.ListenAndServe(":8000", nil)
}

func deleteBooks(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("")

	idResult, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "id has to be an integer value")
		return
	}

	var bookIndex = 0

	for idx, val := range Books {
		if val.Id == idResult {
			bookIndex = idx
		}
	}

	copy(Books[bookIndex:], Books[bookIndex+1:])

	Books[len(Books)-1] = Book{}

	Books = Books[:len(Books)-1]

	fmt.Fprint(w, "Book has been deleted")
}

func updateBooks(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	idResult, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "id has to be an integer value")
		return
	}

	var updatedBook Book
	var bookIndex = 0

	for idx, val := range Books {
		if val.Id == idResult {
			updatedBook = val
			bookIndex = idx
		}
	}

	if updatedBook.Id == 0 {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Book with id %d does not exist", idResult)
		return
	}

	title := r.FormValue("title")
	author := r.FormValue("author")
	desc := r.FormValue("desc")

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "Bad Request")
		return
	}

	updatedBook.Title = title
	updatedBook.Author = author
	updatedBook.Desc = desc

	Books[bookIndex] = updatedBook

	fmt.Fprintf(w, "Book with id %d has been successfully updated", updatedBook.Id)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	bs, err := json.Marshal(Books)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Sumthing When Wrong"))
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bs))
}

func getBooksById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	idResult, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "id has to be an integer value")
		return
	}

	if idResult > len(Books) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Book with id %d does not exist", idResult)
		return
	}

	bs, err := json.Marshal(Books[idResult-1])

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Sumthing When Wrong"))
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bs))
}

func createBooks(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(400)
		fmt.Fprint(w, "invalid Content-Type")
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(422)
		fmt.Fprint(w, "Invalid request Body")
		return
	}

	var request BookRequest

	err = json.Unmarshal(body, &request)

	bookId := len(Books) + 1

	newBook := Book{
		Id:     bookId,
		Title:  request.Title,
		Author: request.Author,
		Desc:   request.Desc,
	}

	Books = append(Books, newBook)

	var response map[string]string = map[string]string{
		"message": "New Book has been successfully created",
	}

	bs, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "internal server error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(bs)
}
