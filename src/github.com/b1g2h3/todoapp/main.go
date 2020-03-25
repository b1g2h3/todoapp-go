package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Client
var client *firestore.Client

// initilaze firestore Client
func init() {
	// Use the application default credentials.
	ctx := context.Background()
	sa := option.WithCredentialsFile("C:/Users/vlast/Desktop/credentionals.json")

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}
}

// List struct (Model)
type List struct {
	ID   int    `json:"ID"`
	Name string `json:"Name"`
}

// Task struct (Model)
type Task struct {
	ID     int    `json:"ID"`
	Name   string `json:"Name"`
	ListID int    `json:"ListID"`
}

var lists []Task

// Get all Lists
func getLists(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	iter := client.Collection("lists").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(doc.Data())
	}

}

// Get Tasks
func getTasks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	ListID, err := strconv.Atoi(params["ListID"])
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	iter := client.Collection("task").Where("ListID", "==", ListID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		json.NewEncoder(w).Encode(doc.Data())

	}
}

// Get Task
func getTask(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	params := mux.Vars(r)
	ListID, err := strconv.Atoi(params["ListID"])
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	TaskID, err := strconv.Atoi(params["TaskID"])
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	iter := client.Collection("task").Where("ListID", "==", ListID).Where("ID", "==", TaskID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(doc.Data())

	}
}

// Add new List
func createList(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var list List
	_ = json.NewDecoder(r.Body).Decode(&list)
	list.ID = rand.Intn(100)
	json.NewEncoder(w).Encode(list)
	_, _, _ = client.Collection("lists").Add(ctx, list)
	w.Header().Set("Content-Type", "application/json")
}

// Add new Task
func createTask(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	ListID, err := strconv.Atoi(params["ListID"])
	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	task.ID = rand.Intn(10000)
	task.ListID = ListID
	Name := strconv.Itoa(task.ID) + strconv.Itoa(task.ListID)
	_, err = client.Collection("task").Doc(Name).Set(ctx, task)
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
	json.NewEncoder(w).Encode(task)
}

// Update Task
func updateTask(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	TaskID, err := strconv.Atoi(params["ListID"])
	ListID, err := strconv.Atoi(params["TaskID"])
	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	task.ID = TaskID
	task.ListID = ListID
	Name := strconv.Itoa(task.ListID) + strconv.Itoa(task.ID)
	_, err = client.Collection("task").Doc(Name).Set(ctx, task)
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
	json.NewEncoder(w).Encode(task)

}

// Delete Task
func deleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	TaskID, err := strconv.Atoi(params["ListID"])
	ListID, err := strconv.Atoi(params["TaskID"])
	Name := strconv.Itoa(ListID) + strconv.Itoa(TaskID)
	_, err = client.Collection("task").Doc(Name).Delete(ctx)
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}

}

// Handler
func handler() {
	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/lists", getLists).Methods("GET")
	r.HandleFunc("/lists", createList).Methods("POST")
	r.HandleFunc("/lists/{ListID}/tasks", getTasks).Methods("GET")
	r.HandleFunc("/lists/{ListID}/tasks", createTask).Methods("POST")
	r.HandleFunc("/lists/{ListID}/tasks/{TaskID}", getTask).Methods("GET")
	r.HandleFunc("/lists/{ListID}/tasks/{TaskID}", updateTask).Methods("PATCH")
	r.HandleFunc("/lists/{ListID}/tasks/{TaskID}", deleteTask).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":8000", r))
}

// Main function
func main() {
	handler()
}
