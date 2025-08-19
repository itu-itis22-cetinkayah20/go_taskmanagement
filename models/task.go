package models

type Task struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Details string `json:"details"`
}

var PublicTasks = []Task{
	{ID: 1, UserID: 0, Title: "Örnek Görev 1", Details: "Bu public bir görevdir."},
	{ID: 2, UserID: 0, Title: "Örnek Görev 2", Details: "Herkes görebilir."},
}

var Tasks = []Task{}
