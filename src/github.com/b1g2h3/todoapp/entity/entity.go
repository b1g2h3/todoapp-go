package entity

// List struct (Model)
type List struct {
	ID   int    `json:"ID"`
	UID  string `json:"UID"`
	Name string `json:"Name"`
}

// Task struct (Model)
type Task struct {
	ID     int    `json:"ID"`
	UID    string `json:"UID"`
	Name   string `json:"Name"`
	ListID int    `json:"ListID"`
}
