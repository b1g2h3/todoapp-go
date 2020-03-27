package repository

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/b1g2h3/todoapp/entity"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// TodoRepository interface
type TodoRepository interface {
	GetLists(list *entity.List) ([]entity.List, error)
	AddList(list *entity.List) (*entity.List, error)
	GetTasks(ListID string) ([]entity.Task, error)
	GetTask(ListID, TaskID string) ([]entity.Task, error)
	AddTask(ListID string, task *entity.Task) (*entity.Task, error)
	UpdateTask(ListID, TaskID string, task *entity.Task) (*entity.Task, error)
	DestroyTask(Name string)
}

type repo struct{}

// Client
var (
	c     *firestore.Client
	lists []entity.List
	tasks []entity.Task
)

// initilaze firestore Client
func init() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("C:/Users/vlast/Desktop/credentionals.json")
	config := &firebase.Config{ProjectID: "todo-3840c"}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	c, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}
}

func (*repo) GetLists(list *entity.List) ([]entity.List, error) {
	ctx := context.Background()
	iter := c.Collection("lists").Where("UID", "==", list.UID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		list := entity.List{
			ID:   doc.Data()["ID"].(int64),
			UID:  doc.Data()["UID"].(string),
			Name: doc.Data()["Name"].(string),
		}
		lists = append(lists, list)
	}
	return lists, nil
}

func (*repo) AddList(list *entity.List) (*entity.List, error) {
	ctx := context.Background()
	list.ID = int64(rand.Intn(10000))
	_, _, err := c.Collection("lists").Add(ctx, list)
	if err != nil {
		log.Fatalf("Failed add list: %v", err)
	}
	return list, nil
}

func (*repo) GetTasks(ListID string) ([]entity.Task, error) {
	ctx := context.Background()
	fmt.Printf(ListID)
	iter := c.Collection("task").Where("ListID", "==", ListID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		task := entity.Task{
			ID:     doc.Data()["ID"].(string),
			ListID: doc.Data()["ListID"].(string),
			UID:    doc.Data()["UID"].(string),
			Name:   doc.Data()["Name"].(string),
		}
		tasks = append(tasks, task)

	}
	return tasks, nil
}

func (*repo) GetTask(ListID, TaskID string) ([]entity.Task, error) {
	ctx := context.Background()
	iter := c.Collection("task").Where("ListID", "==", ListID).Where("ID", "==", TaskID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		task := entity.Task{
			ID:     doc.Data()["ID"].(string),
			ListID: doc.Data()["ListID"].(string),
			UID:    doc.Data()["UID"].(string),
			Name:   doc.Data()["Name"].(string),
		}
		tasks = append(tasks, task)

	}
	return tasks, nil
}

func (*repo) AddTask(ListID string, task *entity.Task) (*entity.Task, error) {
	ctx := context.Background()
	// task.ID = strconv.FormatInt(0)
	task.ListID = ListID
	Name := task.ID + task.ListID
	_, err := c.Collection("task").Doc(Name).Set(ctx, task)
	if err != nil {
		log.Fatalf("Failed add task: %v", err)
	}
	return task, nil
}

func (*repo) UpdateTask(TaskID, ListID string, task *entity.Task) (*entity.Task, error) {
	ctx := context.Background()
	task.ID = TaskID
	task.ListID = ListID
	Name := task.ID + task.ListID
	_, err := c.Collection("task").Doc(Name).Set(ctx, task)
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
	return task, nil
}

func (*repo) DestroyTask(Name string) {
	ctx := context.Background()
	_, err := c.Collection("task").Doc(Name).Delete(ctx)
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
}

// NewTodoRepository init func for  repo
func NewTodoRepository() TodoRepository {
	return &repo{}
}
