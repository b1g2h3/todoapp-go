package entity

// List struct (Model)
type List struct {
	ID   string `json:"ID"`
	UID  string `json:"UID"`
	Name string `json:"Name"`
}

// Task struct (Model)
type Task struct {
	ID     string `json:"ID"`
	ListID string `json:"ListID"`
	UID    string `json:"UID"`
	Name   string `json:"Name"`
}
