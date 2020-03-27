package entity

// List struct (Model)
type List struct {
	ID   int64  `json:"ID"`
	UID  string `json:"UID"`
	Name string `json:"Name"`
}

// Task struct (Model)
type Task struct {
	ID     int    `json:"ID"`
	ListID int    `json:"ListID"`
	UID    int64  `json:"UID"`
	Name   string `json:"Name"`
}
