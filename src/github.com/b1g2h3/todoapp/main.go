package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"text/template"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Client
var client *firestore.Client
var app *firebase.App

// initilaze firestore Client
func init() {
	// Use the application default credentials.
	ctx := context.Background()
	opt := option.WithCredentialsFile("C:/Users/vlast/Desktop/credentionals.json")
	config := &firebase.Config{ProjectID: "todo-3840c"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
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

type idToken struct {
	Value string `json:"Token"`
}

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     "25387032102-ql5jr8rkasa64nbg3ejf7ojjhv12g7ln.apps.googleusercontent.com",
		ClientSecret: "kxYYKu30m7tEY6mk0TCduLY9",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	tpl              *template.Template
	lists            []Task
	oauthStateString = "pseudo-random"
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func authTest(w http.ResponseWriter, r *http.Request) {

}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Content: %s\n", content)
}

func getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}

func login(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
func logout(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

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
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
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
	r.HandleFunc("/", index)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/callback", handleCallback)
	r.HandleFunc("/authTest", authTest)
	r.HandleFunc("/lists", getLists).Methods("GET")
	r.HandleFunc("/lists", createList).Methods("POST")
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
