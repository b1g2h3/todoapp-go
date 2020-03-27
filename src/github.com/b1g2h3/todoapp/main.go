package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"text/template"
	"./repository"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
)


var (
	repo rep.TodoRepository = rep.NewDatabaseHandler()
	list entity.List
	tpl   *template.Template
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

// Get all Lists
func getLists(w http.ResponseWriter, r *http.Request) {
	var list entity.List
	_ = json.NewDecoder(r.Body).Decode(&list)
	lists, err : =repo.GetLists(list)
	result, err := json.Marshal(lists)
	if err != nil {
		w.WriteHeader(http.StatusIntervalServerError)
		w.Write([](`{"error": "Error marshalling data"}`))
	}
	r.WriteHeader(http.StatusOK)
	r.Write(result)
}

// Add new List
func createList(w http.ResponseWriter, r *http.Request) {
	r.WriteHeader(http.StatusOK)
	var list entity.List
	_ = json.NewDecoder(r.Body).Decode(&list)
	list, err : =repo.addList(list)
	result, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusIntervalServerError)
		w.Write([](`{"error": "Error marshalling data"}`))
	}
	r.WriteHeader(http.StatusOK)
	r.Write(result)
}

// Get Tasks
func getTasks(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ListID, err := strconv.Atoi(params["ListID"])
	tasks, err : =repo.getTasks(ListID)
	result, err := json.Marshal(tasks)
	if err != nil {
		w.WriteHeader(http.StatusIntervalServerError)
		w.Write([](`{"error": "Error marshalling data"}`))
	}
	r.WriteHeader(http.StatusOK)
	r.Write(result)
}

// Get Task
func getTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ListID, err := strconv.Atoi(params["ListID"])
	TaskID, err := strconv.Atoi(params["TaskID"])
	task, err : =repo.getTask(ListListID, TaskTaskID)
	result, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusIntervalServerError)
		w.Write([](`{"error": "Error marshalling data"}`))
	}
	r.WriteHeader(http.StatusOK)
	r.Write(result)
}

// Add new Task
func createTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ListID, err := strconv.Atoi(params["ListID"])
	var task entity.Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	newTask, err : =repo.addTask(ListID, task)
	result, err := json.Marshal(newTask)
	if err != nil {
		w.WriteHeader(http.StatusIntervalServerError)
		w.Write([](`{"error": "Error marshalling data"}`))
	}
	r.WriteHeader(http.StatusOK)
	r.Write(result)
}

// Update Task
func updateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	TaskID, err := strconv.Atoi(params["ListID"])
	ListID, err := strconv.Atoi(params["TaskID"])
	var task entity.Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	newTask, err : =repo.updateTask(TaskListID, LisListID, task)
	result, err := json.Marshal(newTask)
	if err != nil {
		w.WriteHeader(http.StatusIntervalServerError)
		w.Write([](`{"error": "Error marshalling data"}`))
	}
	r.WriteHeader(http.StatusOK)
	r.Write(result)
}

// Delete Task
func deleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	TaskID, err := strconv.Atoi(params["ListID"])
	ListID, err := strconv.Atoi(params["TaskID"])
	Name := strconv.Itoa(ListID) + strconv.Itoa(TaskID)
	repo.destroyTask(Name)
	if err != nil {
		w.WriteHeader(http.StatusIntervalServerError)
		w.Write([](`{"error": "Error Delete data"}`))
	}
	r.WriteHeader(http.StatusOK)
	r.Write(result)
}

// Handler
func handler() {
	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/", index)
	r.HandleFunc("/lists", getLists).Methods("POST")
	r.HandleFunc("/lists", createList).Methods("PUT")
	r.HandleFunc("/lists/{ListID}/tasks", getTasks).Methods("GET")
	r.HandleFunc("/lists/{ListID}/tasks", createTask).Methods("POST")
	r.HandleFunc("/lists/{ListID}/tasks/{TaskID}", getTask).Methods("GET")
	r.HandleFunc("/lists/{ListID}/tasks/{TaskID}", updateTask).Methods("PATCH")
	r.HandleFunc("/lists/{ListID}/tasks/{TaskID}", deleteTask).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Main function
func main() {
	handler()
}
