package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/b1g2h3/todoapp/entity"
	"github.com/b1g2h3/todoapp/repository"
	"github.com/gorilla/mux"
)

var (
	list *entity.List
	tpl  *template.Template
)

var (
	repo repository.TodoRepository = repository.NewTodoRepository()
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

// Get all Lists
func getLists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var list *entity.List
	_ = json.NewDecoder(r.Body).Decode(&list)
	lists, err := repo.GetLists(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error getting the lists"}`))
		return
	}
	result, err := json.Marshal(lists)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error marshalling data"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// Add new List
func createList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var list *entity.List
	_ = json.NewDecoder(r.Body).Decode(&list)
	newList, err := repo.AddList(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error adding list"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newList)
}

// Get Tasks
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	ListID := params["ListID"]
	tasks, err := repo.GetTasks(ListID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error getting the lists"}`))
		return
	}
	result, err := json.Marshal(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error marshalling data"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// Get Task
func getTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ListID := params["ListID"]
	TaskID := params["TaskID"]
	task, err := repo.GetTask(ListID, TaskID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error marshalling data"}`))
	}
	result, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error marshalling data"}`))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// Add new Task
func createTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ListID := params["ListID"]
	var Task *entity.Task
	_ = json.NewDecoder(r.Body).Decode(&Task)
	newTask, err := repo.AddTask(ListID, Task)
	result, err := json.Marshal(newTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error marshalling data"}`))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// Update Task
func updateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ListID := params["ListID"]
	TaskID := params["TaskID"]
	var task *entity.Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	newTask, err := repo.UpdateTask(TaskID, ListID, task)
	result, err := json.Marshal(newTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error marshalling data"}`))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// Delete Task
func deleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	TaskID, err := strconv.Atoi(params["ListID"])
	ListID, err := strconv.Atoi(params["TaskID"])
	Name := strconv.Itoa(ListID) + strconv.Itoa(TaskID)
	repo.DestroyTask(Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error marshalling data"}`))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"Success": "Task was deleted"}`))
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
	log.Fatal(http.ListenAndServe(":8090", r))
}

// Main function
func main() {
	handler()
}
