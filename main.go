package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/pseudonator/portal-bookstore-demo/models"
	"github.com/pseudonator/portal-bookstore-demo/utils"
)

// Connection mongoDB with utils class
var collection = utils.ConnectDB()

type Remote struct {
	XFF string `json:"x-forwarded-for"`
}

/*
 * Operation: GET /remote
 * Description: Get x-forwarded-for header
 */
func getRemote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	remote := Remote{XFF: r.Header.Get("x-forwarded-for")}
	result, err := json.Marshal(remote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

/*
 * Operation: GET /api/books
 * Description: Get all books
 */
func getAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var books []models.Book

	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		utils.GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var book models.Book

		err := cur.Decode(&book)
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(books)
}

/*
 * Operation: GET /api/books/{id}
 * Description: Lookup a book
 */
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book models.Book
	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&book)
	if err != nil {
		utils.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(book)
}

/*
 * Operation: POST /api/books
 * Description: Add a book to store
 */
func addBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book models.Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	result, err := collection.InsertOne(context.TODO(), book)
	if err != nil {
		utils.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(result)
}

/*
 * Operation: PUT /api/books/{id}
 * Description: Update an existing book
 */
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book models.Book
	var params = mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": id}
	_ = json.NewDecoder(r.Body).Decode(&book)
	fmt.Println("Updated author", book.Author)
	author := book.Author
	update := bson.D{
		{"$set", bson.D{
			{"Name", book.Name},
			{"Author", book.Author},
		}},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&book)
	if err != nil {
		utils.GetError(err, w)
		return
	}
	book.ID = id
	book.Author = author
	json.NewEncoder(w).Encode(book)
}

/*
 * Operation: DELETE /api/books/{id}
 * Description: Remove a book from the store
 */
func removeBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	removeResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		utils.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(removeResult)
}

/*
 * Operation: DELETE /api/books
 * Description: Remove all books
 */
func removeAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	filter := bson.M{}
	removeResult, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		utils.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(removeResult)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/remote", getRemote).Methods("GET")
	r.HandleFunc("/api/books", getAllBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", addBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", removeBook).Methods("DELETE")
	r.HandleFunc("/api/books", removeAllBooks).Methods("DELETE")

	config := utils.GetConfiguration()

	log.Fatal(http.ListenAndServe(config.Port, r))
}
