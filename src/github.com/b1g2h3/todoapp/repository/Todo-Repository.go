package rep

import (
	"../entity"
)
type NewDatabaseHandler interface {
	getLists(list *entity.List) ([]entity.List, error)
	addList(list *entity.List) (*entity.List, error)
	getTasks(ListID int, task *entity.Task) ([]entity.Task, error)
	getTask(ListID, TaskID int,task *entity.Task) (*entity.Task, error)
	addTask(ListID int,task *entity.Task) (*entity.Task, error)
	updateTask(task *entity.Task) (*entity.Task, error)
	destroyTask(Name String) 
}

type repo struct{}

//NewListsRepository
func NewDatabaseHandler() {
	return &repo{}
}


// Client
var (
	client *firestore.Client
	app *firebase.App
)
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

var lists            []List

func (*repo) getLists(list *List) ([]List, error){
	ctx := context.Background()
	iter := client.Collection("lists").Where("UID", "==", list.UID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		list := List{
			ID:  doc.Data()["ID"].(string),
			UID:  doc.Data()["UID"].(string),
			Name: doc.Data()["Name"].(string),
		}
		lists = append(lists, list)
	}
	return lists, nil
}

func (*repo) addList(list *entity.List) (entity.List, error){
	ctx := context.Background()
	list.ID = rand.Intn(10000)
	_, err = client.Collection("task").Doc(Name).Set(ctx, list)
	return list, nil
}

func (*repo) getTasks(ListID int) ([]entity.Task, error){
	ctx := context.Background()
	iter := client.Collection("task").Where("ListID", "==", ListID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		task := entity.Task{
			ID:  doc.Data()["ID"].(string),
			UID:  doc.Data()["UID"].(string),
			Name: doc.Data()["Name"].(string),
			ListID: doc.Data()["ListID"].(string),
		}
		tasks = append(tasks, task)

	}
	tasks, nil
}

func (*repo) getTask(ListID, TaskID int) ([]entity.Task, error){
	ctx := context.Background()
	iter := client.Collection("task").Where("ListID", "==", ListID).Where("ID", "==", TaskID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		task := entity.Task{
			ID:  doc.Data()["ID"].(string),
			UID:  doc.Data()["UID"].(string),
			Name: doc.Data()["Name"].(string),
			ListID: doc.Data()["ListID"].(string),
		}
		tasks = append(tasks, task)

	}
	tasks, nil
}

func (*repo) addTask(ListID int, task *entity.List) (*entity.Task, error){
	ctx := context.Background()
	task.ID = rand.Intn(10000)
	task.ListID = ListID
	Name := strconv.Itoa(task.ID) + strconv.Itoa(task.ListID)
	_, err = client.Collection("task").Doc(Name).Set(ctx, task)
	return task, nil
}

func (*repo) updateTask(TaskID,ListID int,task *entity.List) (*entity.Task, error){
	ctx := context.Background()
	task.ID = TaskID
	task.ListID = ListID
	Name := strconv.Itoa(task.ListID) + strconv.Itoa(task.ID)
	_, err = client.Collection("task").Doc(Name).Set(ctx, task)
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
	json.NewEncoder(w).Encode(task)
	return task, nil
}

func (*repo) destroyTask(Name string) {
	_, err = client.Collection("task").Doc(Name).Delete(ctx)
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
}