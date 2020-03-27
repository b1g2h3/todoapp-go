package config

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
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
