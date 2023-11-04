/*
// REST API examople with some routing

// POST   /task/              :  create task and take in DB
// GET    /task/<taskid>      :  return task by id
// GET    /task/              :  return all tasks
// DELETE /task/<taskid>      :  delete task by id
// DELETE /task/              :  delete all tasks
// GET    /tag/<tagname>      :  return tasks by tag
// GET    /due/<yy>/<mm>/<dd> :  return tasks by date
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Task struct {
	//gorm.Model
	Id   uint           `json:"id"`
	Text string         `json:"text"`
	Tags []Tag          `json:"tags"` //`gorm:"serializer:json"` //Tags []string!!!
	Due  datatypes.Date `json:"due"`
}

type Tag struct {
	Id     uint
	Text   string
	TaskID uint
}

type App struct {
	DB *gorm.DB
}

func (a *App) Initialize(dbURI string) {
	db, err := gorm.Open(sqlite.Open(dbURI), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	a.DB = db

	//Migrate model in DB
	a.DB.AutoMigrate(&Task{}, &Tag{})
}

func (a *App) getAllTaskHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []Task

	// Select all tasks and convert to JSON.
	a.DB.Preload("Tags").Find(&tasks)
	tasksJSON, _ := json.Marshal(tasks)

	// Write to HTTP response.
	w.WriteHeader(200)
	w.Write([]byte(tasksJSON))
}

func (a *App) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	vars := mux.Vars(r)

	// Select the task with the given id, and convert to JSON.
	result := a.DB.Preload("Tags").First(&task, "id = ?", vars["id"])
	if result.RowsAffected == 0 {
		http.Error(w, "error: id not found in DataBase", http.StatusNotFound)
		return
	}
	taskJSON, _ := json.Marshal(task)

	// Write to HTTP response.
	w.WriteHeader(200)
	w.Write([]byte(taskJSON))
}

func (a *App) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling task create at %s\n", r.URL.Path)
	var newTask Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(400)
		return
	}

	rows := a.DB.Create(&newTask).RowsAffected
	log.Println("Added rows: ", rows)

	// create json for answer
	type ResponseId struct {
		Id uint `json:"id"`
	}

	taskJSON, err := json.Marshal(ResponseId{Id: newTask.Id})
	if err != nil {
		http.Error(w, "error: not create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(taskJSON))
}

func (a *App) getTaskHandlerByTag(w http.ResponseWriter, r *http.Request) {
	var task Task
	var tags []Tag
	var foundTasks []Task

	vars := mux.Vars(r)

	res := a.DB.Where("text = ?", vars["tag"]).Find(&tags)
	if res.RowsAffected == 0 {
		//w.WriteHeader(200)
		http.Error(w, "error: tag not found in DataBase", http.StatusNotFound)
		return
	}

	for _, tag := range tags {
		a.DB.Preload("Tags").Find(&task, tag.TaskID)
		foundTasks = append(foundTasks, task)
	}

	taskJSON, _ := json.Marshal(&foundTasks)

	// Write to HTTP response.
	w.WriteHeader(200)
	w.Write([]byte(taskJSON))

}

func (a *App) getTaskHandlerByDue(w http.ResponseWriter, r *http.Request) {
	var tasks []Task
	vars := mux.Vars(r)

	reqTimeStr := fmt.Sprintf("20%s-%s-%s", vars["yy"], vars["mm"], vars["dd"])
	needTime, _ := time.Parse("2006-01-02", reqTimeStr)

	err := a.DB.Model(&Task{}).Where("due = ?", datatypes.Date(needTime)).Find(&tasks)
	if err.RowsAffected == 0 {
		//w.WriteHeader(200)
		http.Error(w, "error: this date not found in DataBase", http.StatusNotFound)
		return
	}

	taskJSON, _ := json.Marshal(&tasks)

	// Write to HTTP response.
	w.WriteHeader(200)
	w.Write([]byte(taskJSON))
}

func (a *App) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Delete the task with the given id.
	err := a.DB.Where("id = ?", vars["id"]).Delete(Task{})
	if err.Error != nil {
		//w.WriteHeader(200)
		http.Error(w, "error: id not found in DataBase", http.StatusNotFound)
		return
	}

	// Write to HTTP response.
	w.WriteHeader(204)
}

func (a *App) deleteAllTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Delete all tasks.
	err := a.DB.Exec("DELETE FROM tasks")
	if err.Error != nil {
		//w.WriteHeader(200)
		http.Error(w, "error: not create task", http.StatusInternalServerError)
		return
	}

	// Write to HTTP response.
	w.WriteHeader(204)
}

func main() {
	a := &App{}
	a.Initialize("test.db")

	r := mux.NewRouter()

	r.HandleFunc("/task/", a.getAllTaskHandler).Methods("GET")
	r.HandleFunc("/task/{id}", a.getTaskHandler).Methods("GET")
	r.HandleFunc("/task/", a.createTaskHandler).Methods("POST")
	r.HandleFunc("/task/{id}", a.deleteTaskHandler).Methods("DELETE")
	r.HandleFunc("/task/", a.deleteAllTaskHandler).Methods("DELETE")
	r.HandleFunc("/due/{yy}/{mm}/{dd}", a.getTaskHandlerByDue).Methods("GET")
	r.HandleFunc("/tag/{tag}", a.getTaskHandlerByTag).Methods("GET")

	http.Handle("/", r)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
