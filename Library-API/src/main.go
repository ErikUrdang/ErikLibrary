package main

// Library REST API
//
// Erik Urdang
//
// April 2018 for Redeam.com
//
// Sample URIs:
//
// - GET	Get all books			http://localhost:8000/books
// - GET	Get a book by ISBN 		http://localhost:8000/books/0451527127
// - PUT 	Add a book				http://localhost:8000/books/4/Stuff/urdang/erik/2018-04-02
// - DELETE	Delete a book by ISBN	http://localhost:8000/books/0451527127
// - POST	Update a book by ISBN	http://localhost:8000/books/4/Stuff and Nonsense/urdang/erik/2018-04-02


// Instructions
//
//Create a CRUD API to manage a list of Books, fulfilling the following requirements:
// 
//1. Books should have the following Attributes:
// 
//- Title
//- Author
//- Publisher
//- Publish Date
//- Rating (1-3)
//- Status (CheckedIn, CheckedOut)
// 
//2. Each endpoint should have test coverage of both successful and failed requests.
// 
//3. Use a data store of your choice
// 
//4. The project should be deployable via Kubernetes
// 
//5. Send the project along as a .git repository

import (
    "encoding/json"
    "log"
    "net/http"
    "mux"
    "time"
    "fmt"
    "strconv"
)


// mux found at: "github.com/gorilla/mux"

// Book is the base element in this system.

type Book struct {
    Title   	string 		`json:"title,omitempty"`
    Author		string   	`json:"author,omitempty"`
    Publisher	string   	`json:"publisher,omitempty"`
    PubDate		time.Time 	`json:"pubdate,omitempty"`
    Rating      int		  	`json:"rating,omitempty"`
    CheckedOut  bool	  	`json:"checkedOut,omitempty"`
    ISBN        string  	`json:"ISBN,omitempty"`
}


// An in-memory collection of all books

var books []Book

// Encode all of the books in books []

func GetBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting all books.\n")
    json.NewEncoder(w).Encode(books)
}

// Return one book matching the supplied ISBN if found. If not, just encode a string with notification.

func GetBook(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting one book....Hi!\n")
	//t := time.Now()
    params := mux.Vars(r)
    isbn := params["ISBN"]
	fmt.Printf("Getting book [%d].\n", isbn)
	index := FindBook (params["ISBN"], w, r)
	if index < 0 {
		json.NewEncoder(w).Encode("Couldn't find ISBN: [" + isbn + "].")
	} else {
		book := books [index]
		fmt.Printf ("Found it! (%s)\n", book.Title)
		json.NewEncoder(w).Encode(book)
		
	}
}

// Create a new book and add it to the array of books

func CreateBook(w http.ResponseWriter, r *http.Request) {

    params := mux.Vars(r)
    
	fmt.Printf("Creating one book - parameters = %s\n", params)
    
	ok, errStr := ValidateBookValues (r)
	if (! ok) {
	    json.NewEncoder(w).Encode("Error: " + errStr)		
	}

    var book Book
    _ = json.NewDecoder(r.Body).Decode(&book)
    
    // Fill in data structure from URI parameters
    
    books = append(books, book)
    AssignBookValues (len(books) - 1, r)
    
    json.NewEncoder(w).Encode(books)
}

// Find the book with the right ISBN or return -1 if not found.

func FindBook (isbn string, w http.ResponseWriter, r *http.Request) (int) {
    for index, book := range books {
        if book.ISBN == isbn {
            return index
        }
    }
	return -1
}

func ValidateBookValues (r *http.Request) (bool, string) {
    params := mux.Vars(r)

	fmt.Printf("Validating book values - parameters = %s\n", params)

	rating, err := strconv.Atoi(params["rating"])
	if (err != nil || rating < 1 || rating > 3) {
		return false, "Bad rating = " + (params["rating"])
	}
	_, err = strconv.ParseBool(params["checkedOut"])
	if (err != nil) {return false, "Bad boolean value"}
	
	// All tests passed:
	
	return true, "good"
}

// Update the values of the book at the index specified.

func AssignBookValues (index int, r *http.Request){
	// fmt.Printf ("Before assigning values: (%s)\n", books[index])
	
    params := mux.Vars(r)
    books[index].ISBN = params["ISBN"]
    books[index].Author = params["author"]
    books[index].PubDate, _ = time.Parse("2006-01-02", params["date"])
    books[index].Rating, _ = strconv.Atoi(params["rating"])
    books[index].Publisher = params["publisher"]
    books[index].Title = params["title"]
    books[index].CheckedOut, _ = strconv.ParseBool (params["checkedOut"])
	// fmt.Printf ("After assigning values: (%s in %s)\n", books[index], books)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {

    params := mux.Vars(r)
    
	fmt.Printf("Updating one book - parameters = %s\n", params)
	
	ok, errStr := ValidateBookValues (r)
	if (! ok) {
	    json.NewEncoder(w).Encode("Error: " + errStr)		
	}
	index := FindBook (params["ISBN"], w, r)
	if (index < 0) {
		fmt.Printf ("No matching book")
	    json.NewEncoder(w).Encode("No matching book")
		
	} else {
		fmt.Printf ("Found it! (%s)\n", books[index])
		AssignBookValues (index, r)
	    json.NewEncoder(w).Encode(books)
		
	}
}

// Replace item "index" with the value from the end of the array and
// then truncate the array.
//
// This assumes:
//
// - order doesn't matter
// - there are no duplicate entries (i.e., same ISBN)

func RemoveFromBooks (index int) {
	n := len(books)
	books [index] = books [n - 1]
	books = books[:n - 1]
}

// Delete an item
func DeleteBook(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    isbn := params ["ISBN"]
	fmt.Printf("Deleting one book (ISBN == %s)", isbn)
    for index, item := range books {
        if item.ISBN == isbn {
           // books = append(books[:index], books[index+1:]...)
			RemoveFromBooks (index)
            break
        }
        json.NewEncoder(w).Encode(books)
    }
}

// Create a couple of books for testing

func CreateSampleBooks() {
	
	ws, _ := time.Parse("2006-01-02", "1623-01-01")
	sk, _ := time.Parse("2006-01-02", "1986-09-05")
	
	books = append(books, 
		Book{
			ISBN: "0451527127", 
			Author: "William Shakespeare", 
			Title: "Tempest, The", 
			PubDate: ws, 
			Rating: 1, 
			CheckedOut: true})
    books = append(books, 
    	Book{
    		ISBN: "1444707868", 
    		Author: "Stephen King", 
    		Title: "It", 
    		PubDate: sk, 
			Rating: 2, 
    		CheckedOut: false})

}

// main function:
//
// - Create a router
// - Add a few sample books to the array
// - Specify which URIs are associated with which handler functions
// - Listen on port 8000 for incoming URIs

func main() {
    router := mux.NewRouter()
    
    CreateSampleBooks()
    
    router.HandleFunc("/books", GetBooks).Methods("GET")
    router.HandleFunc("/books/{ISBN}", GetBook).Methods("GET")
    router.HandleFunc("/books/{ISBN}/{title}/{author}/{date}/{publisher}/{rating}/{checkedOut}", CreateBook).Methods("POST")
    router.HandleFunc("/books/{ISBN}/{title}/{author}/{date}/{publisher}/{rating}/{checkedOut}", UpdateBook).Methods("PUT")
    router.HandleFunc("/books/{ISBN}", DeleteBook).Methods("DELETE")
    
    log.Fatal(http.ListenAndServe(":8000", router))
}