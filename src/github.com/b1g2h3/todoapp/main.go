package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

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
	sa := option.WithCredentialsJSON([]byte(`{
		"type": "service_account",
		"project_id": "todo-3840c",
		"private_key_id": "3f883b908a286fa7a4301a7ee8821eb5b79398d0",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDRNVU/EjckVNdx\n6RG9hhEG1se4UTK6gje/NAXvufcqmHtQBvkI4QxgamvKNMzsST2c5hGOZs9KPaZN\nqyJyqrxzhk2RYez17/HPlhvNOKPkaR0VJoyXA7g87HOMTInao03GNBuhiAkFaNin\n/67X6m0bycbDLTBGFxGFryMmT8L0zjP3TFHq7TI2lJO90cq45s2dM509HLu7DgXU\nO6PlDZkYNnnVqNTsxch98z3UinO7AlcQfKCNh4kzfx1K8GqSOId4nEKosFsviBUE\nRyC+PxXXWlwjHigQXcBDsdX9qqHYjGV4pBHrinWikEeeGWAd+NXenBIJTvGlgICJ\nWgmyLBfHAgMBAAECggEAXU1X02UgSoe7/gVf8BJWjaJEOCOeaCejDRb6fsuFO+39\nMNerQRZ9GpLbt7aMneScIdlJgyS+1fFgtcrY9iLHIQ6IkYoG0DhOs7HSfFgCX9+x\nJLmogcEa7bDWZ3/LC0NBcF/U4tl0jIER/vq803atanM2vdztZpTrL5/IIVH5NIvj\naI3puIWMUtpSQcc4BLmeOzxPVHN8B8jlnAPpSLN7GxJGdgkLzWO+tx1hiLUT79BC\nFP90ew5ORz5kMCl5wTAAMQa48fHzmPAO/m8NU1W8rtE115ygjvL4tb8nYWPnvf03\nAoLzQD3M81DRCdT95oQlDUmHSWeduAG9ihKtI7xegQKBgQDx8riYrOWLdo+C1P7L\n3Txwzrs+AO96jpF7BKrXHxScOYe0nUqZLla4/wkbEALOKGGSc3QurnFKaLkoAsmW\n/ZHtwHo8azQbpGFHzyKUGyrs9K1S03v3R/vNk2koQwIiBZ5koJV1kgQ3auCmEP2u\nVYIdr7yHjc82RQy2364jk2+SUQKBgQDdW9ZgBapGpVXG2EZZh2RoL3Wm3/C1xryt\nzXw5GIsMydMogaFhEdsbc+C0q3UAp9uBQ+z/4YPSoDMzsRENkzvoOTYDFMjSLKFj\nqSAaH0Ap05gmWQlL+TmyJeWx9q3igCs8TPOMYYDNXD88atOGYQ7PSyinXErfOIRq\ncIcViUOqlwKBgQCrtzmeWi98IMRP9b10kOshoQexRNayY9cKuVBK52soSYhv/qaA\nOywflhovU9i52l0NpNVTgEk1p0eqBvhuKj9UvyPCF8/ewnaskW0YMoPvsuQEgcZc\nxYEH8VRT1+L+pIA7KOGKlPxbHIaeNjblcRis2xnyFwp2mOEiNXSRGUW5UQKBgFn5\nECOradCZN0pBcibFv2wRjlKrx107UEmcshdLAInMJwXZ2sxnw5Ve/kCxSDdiAviB\nsX04Hqqn7ufd2r6Xz8vOJUQPWKkE9vxZK/EyLpRRqxA7NGoq/OaKPNifGYJs8iXq\naTvwDbhq/FEEYsHGBY0AUZ/lBZHBmSDiaCW6y0Q1AoGANG7uUrO+8oiO7sO/8258\nG07HSTIeASj9HKpHJToIdsVDvgtd1A9ECQrR3bO66RgYKzCdHl9ehDKg+CsZPpI/\nE+Ak6aKT1YlWegRY4MLfRlbit0oHrCmswUtPVLoAxi15KkUGMppNCt4ggYhM8o4b\neD/60Ow4SwXjJjJhiRu8PXc=\n-----END PRIVATE KEY-----\n",
		"client_email": "firebase-adminsdk-chomf@todo-3840c.iam.gserviceaccount.com",
		"client_id": "106237896442319976647",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-chomf%40todo-3840c.iam.gserviceaccount.com"
	  }
	  `))

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
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
	sa := option.WithCredentialsJSON([]byte(`{
		"type": "service_account",
		"project_id": "todo-3840c",
		"private_key_id": "3f883b908a286fa7a4301a7ee8821eb5b79398d0",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDRNVU/EjckVNdx\n6RG9hhEG1se4UTK6gje/NAXvufcqmHtQBvkI4QxgamvKNMzsST2c5hGOZs9KPaZN\nqyJyqrxzhk2RYez17/HPlhvNOKPkaR0VJoyXA7g87HOMTInao03GNBuhiAkFaNin\n/67X6m0bycbDLTBGFxGFryMmT8L0zjP3TFHq7TI2lJO90cq45s2dM509HLu7DgXU\nO6PlDZkYNnnVqNTsxch98z3UinO7AlcQfKCNh4kzfx1K8GqSOId4nEKosFsviBUE\nRyC+PxXXWlwjHigQXcBDsdX9qqHYjGV4pBHrinWikEeeGWAd+NXenBIJTvGlgICJ\nWgmyLBfHAgMBAAECggEAXU1X02UgSoe7/gVf8BJWjaJEOCOeaCejDRb6fsuFO+39\nMNerQRZ9GpLbt7aMneScIdlJgyS+1fFgtcrY9iLHIQ6IkYoG0DhOs7HSfFgCX9+x\nJLmogcEa7bDWZ3/LC0NBcF/U4tl0jIER/vq803atanM2vdztZpTrL5/IIVH5NIvj\naI3puIWMUtpSQcc4BLmeOzxPVHN8B8jlnAPpSLN7GxJGdgkLzWO+tx1hiLUT79BC\nFP90ew5ORz5kMCl5wTAAMQa48fHzmPAO/m8NU1W8rtE115ygjvL4tb8nYWPnvf03\nAoLzQD3M81DRCdT95oQlDUmHSWeduAG9ihKtI7xegQKBgQDx8riYrOWLdo+C1P7L\n3Txwzrs+AO96jpF7BKrXHxScOYe0nUqZLla4/wkbEALOKGGSc3QurnFKaLkoAsmW\n/ZHtwHo8azQbpGFHzyKUGyrs9K1S03v3R/vNk2koQwIiBZ5koJV1kgQ3auCmEP2u\nVYIdr7yHjc82RQy2364jk2+SUQKBgQDdW9ZgBapGpVXG2EZZh2RoL3Wm3/C1xryt\nzXw5GIsMydMogaFhEdsbc+C0q3UAp9uBQ+z/4YPSoDMzsRENkzvoOTYDFMjSLKFj\nqSAaH0Ap05gmWQlL+TmyJeWx9q3igCs8TPOMYYDNXD88atOGYQ7PSyinXErfOIRq\ncIcViUOqlwKBgQCrtzmeWi98IMRP9b10kOshoQexRNayY9cKuVBK52soSYhv/qaA\nOywflhovU9i52l0NpNVTgEk1p0eqBvhuKj9UvyPCF8/ewnaskW0YMoPvsuQEgcZc\nxYEH8VRT1+L+pIA7KOGKlPxbHIaeNjblcRis2xnyFwp2mOEiNXSRGUW5UQKBgFn5\nECOradCZN0pBcibFv2wRjlKrx107UEmcshdLAInMJwXZ2sxnw5Ve/kCxSDdiAviB\nsX04Hqqn7ufd2r6Xz8vOJUQPWKkE9vxZK/EyLpRRqxA7NGoq/OaKPNifGYJs8iXq\naTvwDbhq/FEEYsHGBY0AUZ/lBZHBmSDiaCW6y0Q1AoGANG7uUrO+8oiO7sO/8258\nG07HSTIeASj9HKpHJToIdsVDvgtd1A9ECQrR3bO66RgYKzCdHl9ehDKg+CsZPpI/\nE+Ak6aKT1YlWegRY4MLfRlbit0oHrCmswUtPVLoAxi15KkUGMppNCt4ggYhM8o4b\neD/60Ow4SwXjJjJhiRu8PXc=\n-----END PRIVATE KEY-----\n",
		"client_email": "firebase-adminsdk-chomf@todo-3840c.iam.gserviceaccount.com",
		"client_id": "106237896442319976647",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-chomf%40todo-3840c.iam.gserviceaccount.com"
	  }
	  `))

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	ListID, err := strconv.Atoi(params["ListID"])
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
	sa := option.WithCredentialsJSON([]byte(`{
		"type": "service_account",
		"project_id": "todo-3840c",
		"private_key_id": "3f883b908a286fa7a4301a7ee8821eb5b79398d0",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDRNVU/EjckVNdx\n6RG9hhEG1se4UTK6gje/NAXvufcqmHtQBvkI4QxgamvKNMzsST2c5hGOZs9KPaZN\nqyJyqrxzhk2RYez17/HPlhvNOKPkaR0VJoyXA7g87HOMTInao03GNBuhiAkFaNin\n/67X6m0bycbDLTBGFxGFryMmT8L0zjP3TFHq7TI2lJO90cq45s2dM509HLu7DgXU\nO6PlDZkYNnnVqNTsxch98z3UinO7AlcQfKCNh4kzfx1K8GqSOId4nEKosFsviBUE\nRyC+PxXXWlwjHigQXcBDsdX9qqHYjGV4pBHrinWikEeeGWAd+NXenBIJTvGlgICJ\nWgmyLBfHAgMBAAECggEAXU1X02UgSoe7/gVf8BJWjaJEOCOeaCejDRb6fsuFO+39\nMNerQRZ9GpLbt7aMneScIdlJgyS+1fFgtcrY9iLHIQ6IkYoG0DhOs7HSfFgCX9+x\nJLmogcEa7bDWZ3/LC0NBcF/U4tl0jIER/vq803atanM2vdztZpTrL5/IIVH5NIvj\naI3puIWMUtpSQcc4BLmeOzxPVHN8B8jlnAPpSLN7GxJGdgkLzWO+tx1hiLUT79BC\nFP90ew5ORz5kMCl5wTAAMQa48fHzmPAO/m8NU1W8rtE115ygjvL4tb8nYWPnvf03\nAoLzQD3M81DRCdT95oQlDUmHSWeduAG9ihKtI7xegQKBgQDx8riYrOWLdo+C1P7L\n3Txwzrs+AO96jpF7BKrXHxScOYe0nUqZLla4/wkbEALOKGGSc3QurnFKaLkoAsmW\n/ZHtwHo8azQbpGFHzyKUGyrs9K1S03v3R/vNk2koQwIiBZ5koJV1kgQ3auCmEP2u\nVYIdr7yHjc82RQy2364jk2+SUQKBgQDdW9ZgBapGpVXG2EZZh2RoL3Wm3/C1xryt\nzXw5GIsMydMogaFhEdsbc+C0q3UAp9uBQ+z/4YPSoDMzsRENkzvoOTYDFMjSLKFj\nqSAaH0Ap05gmWQlL+TmyJeWx9q3igCs8TPOMYYDNXD88atOGYQ7PSyinXErfOIRq\ncIcViUOqlwKBgQCrtzmeWi98IMRP9b10kOshoQexRNayY9cKuVBK52soSYhv/qaA\nOywflhovU9i52l0NpNVTgEk1p0eqBvhuKj9UvyPCF8/ewnaskW0YMoPvsuQEgcZc\nxYEH8VRT1+L+pIA7KOGKlPxbHIaeNjblcRis2xnyFwp2mOEiNXSRGUW5UQKBgFn5\nECOradCZN0pBcibFv2wRjlKrx107UEmcshdLAInMJwXZ2sxnw5Ve/kCxSDdiAviB\nsX04Hqqn7ufd2r6Xz8vOJUQPWKkE9vxZK/EyLpRRqxA7NGoq/OaKPNifGYJs8iXq\naTvwDbhq/FEEYsHGBY0AUZ/lBZHBmSDiaCW6y0Q1AoGANG7uUrO+8oiO7sO/8258\nG07HSTIeASj9HKpHJToIdsVDvgtd1A9ECQrR3bO66RgYKzCdHl9ehDKg+CsZPpI/\nE+Ak6aKT1YlWegRY4MLfRlbit0oHrCmswUtPVLoAxi15KkUGMppNCt4ggYhM8o4b\neD/60Ow4SwXjJjJhiRu8PXc=\n-----END PRIVATE KEY-----\n",
		"client_email": "firebase-adminsdk-chomf@todo-3840c.iam.gserviceaccount.com",
		"client_id": "106237896442319976647",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-chomf%40todo-3840c.iam.gserviceaccount.com"
	  }
	  `))

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	params := mux.Vars(r)
	ListID, err := strconv.Atoi(params["ListID"])
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
	sa := option.WithCredentialsJSON([]byte(`{
			"type": "service_account",
			"project_id": "todo-3840c",
			"private_key_id": "3f883b908a286fa7a4301a7ee8821eb5b79398d0",
			"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDRNVU/EjckVNdx\n6RG9hhEG1se4UTK6gje/NAXvufcqmHtQBvkI4QxgamvKNMzsST2c5hGOZs9KPaZN\nqyJyqrxzhk2RYez17/HPlhvNOKPkaR0VJoyXA7g87HOMTInao03GNBuhiAkFaNin\n/67X6m0bycbDLTBGFxGFryMmT8L0zjP3TFHq7TI2lJO90cq45s2dM509HLu7DgXU\nO6PlDZkYNnnVqNTsxch98z3UinO7AlcQfKCNh4kzfx1K8GqSOId4nEKosFsviBUE\nRyC+PxXXWlwjHigQXcBDsdX9qqHYjGV4pBHrinWikEeeGWAd+NXenBIJTvGlgICJ\nWgmyLBfHAgMBAAECggEAXU1X02UgSoe7/gVf8BJWjaJEOCOeaCejDRb6fsuFO+39\nMNerQRZ9GpLbt7aMneScIdlJgyS+1fFgtcrY9iLHIQ6IkYoG0DhOs7HSfFgCX9+x\nJLmogcEa7bDWZ3/LC0NBcF/U4tl0jIER/vq803atanM2vdztZpTrL5/IIVH5NIvj\naI3puIWMUtpSQcc4BLmeOzxPVHN8B8jlnAPpSLN7GxJGdgkLzWO+tx1hiLUT79BC\nFP90ew5ORz5kMCl5wTAAMQa48fHzmPAO/m8NU1W8rtE115ygjvL4tb8nYWPnvf03\nAoLzQD3M81DRCdT95oQlDUmHSWeduAG9ihKtI7xegQKBgQDx8riYrOWLdo+C1P7L\n3Txwzrs+AO96jpF7BKrXHxScOYe0nUqZLla4/wkbEALOKGGSc3QurnFKaLkoAsmW\n/ZHtwHo8azQbpGFHzyKUGyrs9K1S03v3R/vNk2koQwIiBZ5koJV1kgQ3auCmEP2u\nVYIdr7yHjc82RQy2364jk2+SUQKBgQDdW9ZgBapGpVXG2EZZh2RoL3Wm3/C1xryt\nzXw5GIsMydMogaFhEdsbc+C0q3UAp9uBQ+z/4YPSoDMzsRENkzvoOTYDFMjSLKFj\nqSAaH0Ap05gmWQlL+TmyJeWx9q3igCs8TPOMYYDNXD88atOGYQ7PSyinXErfOIRq\ncIcViUOqlwKBgQCrtzmeWi98IMRP9b10kOshoQexRNayY9cKuVBK52soSYhv/qaA\nOywflhovU9i52l0NpNVTgEk1p0eqBvhuKj9UvyPCF8/ewnaskW0YMoPvsuQEgcZc\nxYEH8VRT1+L+pIA7KOGKlPxbHIaeNjblcRis2xnyFwp2mOEiNXSRGUW5UQKBgFn5\nECOradCZN0pBcibFv2wRjlKrx107UEmcshdLAInMJwXZ2sxnw5Ve/kCxSDdiAviB\nsX04Hqqn7ufd2r6Xz8vOJUQPWKkE9vxZK/EyLpRRqxA7NGoq/OaKPNifGYJs8iXq\naTvwDbhq/FEEYsHGBY0AUZ/lBZHBmSDiaCW6y0Q1AoGANG7uUrO+8oiO7sO/8258\nG07HSTIeASj9HKpHJToIdsVDvgtd1A9ECQrR3bO66RgYKzCdHl9ehDKg+CsZPpI/\nE+Ak6aKT1YlWegRY4MLfRlbit0oHrCmswUtPVLoAxi15KkUGMppNCt4ggYhM8o4b\neD/60Ow4SwXjJjJhiRu8PXc=\n-----END PRIVATE KEY-----\n",
			"client_email": "firebase-adminsdk-chomf@todo-3840c.iam.gserviceaccount.com",
			"client_id": "106237896442319976647",
			"auth_uri": "https://accounts.google.com/o/oauth2/auth",
			"token_uri": "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-chomf%40todo-3840c.iam.gserviceaccount.com"
		  }
		  `))

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	var list List
	_ = json.NewDecoder(r.Body).Decode(&list)
	list.ID = rand.Intn(100)
	json.NewEncoder(w).Encode(list)
	_, _, err = client.Collection("lists").Add(ctx, list)
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
}

// Add new Task
func createTask(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	sa := option.WithCredentialsJSON([]byte(`{
			"type": "service_account",
			"project_id": "todo-3840c",
			"private_key_id": "3f883b908a286fa7a4301a7ee8821eb5b79398d0",
			"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDRNVU/EjckVNdx\n6RG9hhEG1se4UTK6gje/NAXvufcqmHtQBvkI4QxgamvKNMzsST2c5hGOZs9KPaZN\nqyJyqrxzhk2RYez17/HPlhvNOKPkaR0VJoyXA7g87HOMTInao03GNBuhiAkFaNin\n/67X6m0bycbDLTBGFxGFryMmT8L0zjP3TFHq7TI2lJO90cq45s2dM509HLu7DgXU\nO6PlDZkYNnnVqNTsxch98z3UinO7AlcQfKCNh4kzfx1K8GqSOId4nEKosFsviBUE\nRyC+PxXXWlwjHigQXcBDsdX9qqHYjGV4pBHrinWikEeeGWAd+NXenBIJTvGlgICJ\nWgmyLBfHAgMBAAECggEAXU1X02UgSoe7/gVf8BJWjaJEOCOeaCejDRb6fsuFO+39\nMNerQRZ9GpLbt7aMneScIdlJgyS+1fFgtcrY9iLHIQ6IkYoG0DhOs7HSfFgCX9+x\nJLmogcEa7bDWZ3/LC0NBcF/U4tl0jIER/vq803atanM2vdztZpTrL5/IIVH5NIvj\naI3puIWMUtpSQcc4BLmeOzxPVHN8B8jlnAPpSLN7GxJGdgkLzWO+tx1hiLUT79BC\nFP90ew5ORz5kMCl5wTAAMQa48fHzmPAO/m8NU1W8rtE115ygjvL4tb8nYWPnvf03\nAoLzQD3M81DRCdT95oQlDUmHSWeduAG9ihKtI7xegQKBgQDx8riYrOWLdo+C1P7L\n3Txwzrs+AO96jpF7BKrXHxScOYe0nUqZLla4/wkbEALOKGGSc3QurnFKaLkoAsmW\n/ZHtwHo8azQbpGFHzyKUGyrs9K1S03v3R/vNk2koQwIiBZ5koJV1kgQ3auCmEP2u\nVYIdr7yHjc82RQy2364jk2+SUQKBgQDdW9ZgBapGpVXG2EZZh2RoL3Wm3/C1xryt\nzXw5GIsMydMogaFhEdsbc+C0q3UAp9uBQ+z/4YPSoDMzsRENkzvoOTYDFMjSLKFj\nqSAaH0Ap05gmWQlL+TmyJeWx9q3igCs8TPOMYYDNXD88atOGYQ7PSyinXErfOIRq\ncIcViUOqlwKBgQCrtzmeWi98IMRP9b10kOshoQexRNayY9cKuVBK52soSYhv/qaA\nOywflhovU9i52l0NpNVTgEk1p0eqBvhuKj9UvyPCF8/ewnaskW0YMoPvsuQEgcZc\nxYEH8VRT1+L+pIA7KOGKlPxbHIaeNjblcRis2xnyFwp2mOEiNXSRGUW5UQKBgFn5\nECOradCZN0pBcibFv2wRjlKrx107UEmcshdLAInMJwXZ2sxnw5Ve/kCxSDdiAviB\nsX04Hqqn7ufd2r6Xz8vOJUQPWKkE9vxZK/EyLpRRqxA7NGoq/OaKPNifGYJs8iXq\naTvwDbhq/FEEYsHGBY0AUZ/lBZHBmSDiaCW6y0Q1AoGANG7uUrO+8oiO7sO/8258\nG07HSTIeASj9HKpHJToIdsVDvgtd1A9ECQrR3bO66RgYKzCdHl9ehDKg+CsZPpI/\nE+Ak6aKT1YlWegRY4MLfRlbit0oHrCmswUtPVLoAxi15KkUGMppNCt4ggYhM8o4b\neD/60Ow4SwXjJjJhiRu8PXc=\n-----END PRIVATE KEY-----\n",
			"client_email": "firebase-adminsdk-chomf@todo-3840c.iam.gserviceaccount.com",
			"client_id": "106237896442319976647",
			"auth_uri": "https://accounts.google.com/o/oauth2/auth",
			"token_uri": "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-chomf%40todo-3840c.iam.gserviceaccount.com"
		  }
		  `))

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
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
	sa := option.WithCredentialsJSON([]byte(`{
		"type": "service_account",
		"project_id": "todo-3840c",
		"private_key_id": "3f883b908a286fa7a4301a7ee8821eb5b79398d0",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDRNVU/EjckVNdx\n6RG9hhEG1se4UTK6gje/NAXvufcqmHtQBvkI4QxgamvKNMzsST2c5hGOZs9KPaZN\nqyJyqrxzhk2RYez17/HPlhvNOKPkaR0VJoyXA7g87HOMTInao03GNBuhiAkFaNin\n/67X6m0bycbDLTBGFxGFryMmT8L0zjP3TFHq7TI2lJO90cq45s2dM509HLu7DgXU\nO6PlDZkYNnnVqNTsxch98z3UinO7AlcQfKCNh4kzfx1K8GqSOId4nEKosFsviBUE\nRyC+PxXXWlwjHigQXcBDsdX9qqHYjGV4pBHrinWikEeeGWAd+NXenBIJTvGlgICJ\nWgmyLBfHAgMBAAECggEAXU1X02UgSoe7/gVf8BJWjaJEOCOeaCejDRb6fsuFO+39\nMNerQRZ9GpLbt7aMneScIdlJgyS+1fFgtcrY9iLHIQ6IkYoG0DhOs7HSfFgCX9+x\nJLmogcEa7bDWZ3/LC0NBcF/U4tl0jIER/vq803atanM2vdztZpTrL5/IIVH5NIvj\naI3puIWMUtpSQcc4BLmeOzxPVHN8B8jlnAPpSLN7GxJGdgkLzWO+tx1hiLUT79BC\nFP90ew5ORz5kMCl5wTAAMQa48fHzmPAO/m8NU1W8rtE115ygjvL4tb8nYWPnvf03\nAoLzQD3M81DRCdT95oQlDUmHSWeduAG9ihKtI7xegQKBgQDx8riYrOWLdo+C1P7L\n3Txwzrs+AO96jpF7BKrXHxScOYe0nUqZLla4/wkbEALOKGGSc3QurnFKaLkoAsmW\n/ZHtwHo8azQbpGFHzyKUGyrs9K1S03v3R/vNk2koQwIiBZ5koJV1kgQ3auCmEP2u\nVYIdr7yHjc82RQy2364jk2+SUQKBgQDdW9ZgBapGpVXG2EZZh2RoL3Wm3/C1xryt\nzXw5GIsMydMogaFhEdsbc+C0q3UAp9uBQ+z/4YPSoDMzsRENkzvoOTYDFMjSLKFj\nqSAaH0Ap05gmWQlL+TmyJeWx9q3igCs8TPOMYYDNXD88atOGYQ7PSyinXErfOIRq\ncIcViUOqlwKBgQCrtzmeWi98IMRP9b10kOshoQexRNayY9cKuVBK52soSYhv/qaA\nOywflhovU9i52l0NpNVTgEk1p0eqBvhuKj9UvyPCF8/ewnaskW0YMoPvsuQEgcZc\nxYEH8VRT1+L+pIA7KOGKlPxbHIaeNjblcRis2xnyFwp2mOEiNXSRGUW5UQKBgFn5\nECOradCZN0pBcibFv2wRjlKrx107UEmcshdLAInMJwXZ2sxnw5Ve/kCxSDdiAviB\nsX04Hqqn7ufd2r6Xz8vOJUQPWKkE9vxZK/EyLpRRqxA7NGoq/OaKPNifGYJs8iXq\naTvwDbhq/FEEYsHGBY0AUZ/lBZHBmSDiaCW6y0Q1AoGANG7uUrO+8oiO7sO/8258\nG07HSTIeASj9HKpHJToIdsVDvgtd1A9ECQrR3bO66RgYKzCdHl9ehDKg+CsZPpI/\nE+Ak6aKT1YlWegRY4MLfRlbit0oHrCmswUtPVLoAxi15KkUGMppNCt4ggYhM8o4b\neD/60Ow4SwXjJjJhiRu8PXc=\n-----END PRIVATE KEY-----\n",
		"client_email": "firebase-adminsdk-chomf@todo-3840c.iam.gserviceaccount.com",
		"client_id": "106237896442319976647",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-chomf%40todo-3840c.iam.gserviceaccount.com"
	  }
	  `))

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
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
	sa := option.WithCredentialsJSON([]byte(`{
		"type": "service_account",
		"project_id": "todo-3840c",
		"private_key_id": "3f883b908a286fa7a4301a7ee8821eb5b79398d0",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDRNVU/EjckVNdx\n6RG9hhEG1se4UTK6gje/NAXvufcqmHtQBvkI4QxgamvKNMzsST2c5hGOZs9KPaZN\nqyJyqrxzhk2RYez17/HPlhvNOKPkaR0VJoyXA7g87HOMTInao03GNBuhiAkFaNin\n/67X6m0bycbDLTBGFxGFryMmT8L0zjP3TFHq7TI2lJO90cq45s2dM509HLu7DgXU\nO6PlDZkYNnnVqNTsxch98z3UinO7AlcQfKCNh4kzfx1K8GqSOId4nEKosFsviBUE\nRyC+PxXXWlwjHigQXcBDsdX9qqHYjGV4pBHrinWikEeeGWAd+NXenBIJTvGlgICJ\nWgmyLBfHAgMBAAECggEAXU1X02UgSoe7/gVf8BJWjaJEOCOeaCejDRb6fsuFO+39\nMNerQRZ9GpLbt7aMneScIdlJgyS+1fFgtcrY9iLHIQ6IkYoG0DhOs7HSfFgCX9+x\nJLmogcEa7bDWZ3/LC0NBcF/U4tl0jIER/vq803atanM2vdztZpTrL5/IIVH5NIvj\naI3puIWMUtpSQcc4BLmeOzxPVHN8B8jlnAPpSLN7GxJGdgkLzWO+tx1hiLUT79BC\nFP90ew5ORz5kMCl5wTAAMQa48fHzmPAO/m8NU1W8rtE115ygjvL4tb8nYWPnvf03\nAoLzQD3M81DRCdT95oQlDUmHSWeduAG9ihKtI7xegQKBgQDx8riYrOWLdo+C1P7L\n3Txwzrs+AO96jpF7BKrXHxScOYe0nUqZLla4/wkbEALOKGGSc3QurnFKaLkoAsmW\n/ZHtwHo8azQbpGFHzyKUGyrs9K1S03v3R/vNk2koQwIiBZ5koJV1kgQ3auCmEP2u\nVYIdr7yHjc82RQy2364jk2+SUQKBgQDdW9ZgBapGpVXG2EZZh2RoL3Wm3/C1xryt\nzXw5GIsMydMogaFhEdsbc+C0q3UAp9uBQ+z/4YPSoDMzsRENkzvoOTYDFMjSLKFj\nqSAaH0Ap05gmWQlL+TmyJeWx9q3igCs8TPOMYYDNXD88atOGYQ7PSyinXErfOIRq\ncIcViUOqlwKBgQCrtzmeWi98IMRP9b10kOshoQexRNayY9cKuVBK52soSYhv/qaA\nOywflhovU9i52l0NpNVTgEk1p0eqBvhuKj9UvyPCF8/ewnaskW0YMoPvsuQEgcZc\nxYEH8VRT1+L+pIA7KOGKlPxbHIaeNjblcRis2xnyFwp2mOEiNXSRGUW5UQKBgFn5\nECOradCZN0pBcibFv2wRjlKrx107UEmcshdLAInMJwXZ2sxnw5Ve/kCxSDdiAviB\nsX04Hqqn7ufd2r6Xz8vOJUQPWKkE9vxZK/EyLpRRqxA7NGoq/OaKPNifGYJs8iXq\naTvwDbhq/FEEYsHGBY0AUZ/lBZHBmSDiaCW6y0Q1AoGANG7uUrO+8oiO7sO/8258\nG07HSTIeASj9HKpHJToIdsVDvgtd1A9ECQrR3bO66RgYKzCdHl9ehDKg+CsZPpI/\nE+Ak6aKT1YlWegRY4MLfRlbit0oHrCmswUtPVLoAxi15KkUGMppNCt4ggYhM8o4b\neD/60Ow4SwXjJjJhiRu8PXc=\n-----END PRIVATE KEY-----\n",
		"client_email": "firebase-adminsdk-chomf@todo-3840c.iam.gserviceaccount.com",
		"client_id": "106237896442319976647",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-chomf%40todo-3840c.iam.gserviceaccount.com"
	  }
	  `))

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

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
