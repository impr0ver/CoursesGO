/*
// Пример REST сервера с несколькими маршрутами(используем только стандартную библиотеку)

// POST   /task/              :  создаёт задачу и возвращает её ID
// GET    /task/<taskid>      :  возвращает одну задачу по её ID
// GET    /task/              :  возвращает все задачи
// DELETE /task/<taskid>      :  удаляет задачу по ID
// DELETE /task/              :  удаляет все задачи
// GET    /tag/<tagname>      :  возвращает список задач с заданным тегом
// GET    /due/<yy>/<mm>/<dd> :  возвращает список задач, запланированных на указанную дату

Структура проекта
https://github.com/golang-standards/project-layout/blob/master/README_ru.md
*/
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm" //"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Task struct {
	//gorm.Model
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `gorm:"serializer:json"` //Tags []string!!!
	Due  time.Time `json:"due"`
}

type App struct {
	DB *gorm.DB
}

func (a *App) Initialize(dbURI string) {
	db, err := gorm.Open(sqlite.Open(dbURI), &gorm.Config{})
	//db, err := gorm.Open(dbDriver, dbURI)
	if err != nil {
		panic("failed to connect database")
	}
	a.DB = db

	// Мигрируем в базу
	a.DB.AutoMigrate(&Task{})
}

func (a *App) getAllTaskHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []Task

	// Select all tasks and convert to JSON.
	a.DB.Find(&tasks)
	tasksJSON, _ := json.Marshal(tasks)

	// Write to HTTP response.
	w.WriteHeader(200)
	w.Write([]byte(tasksJSON))
}

func (a *App) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	vars := mux.Vars(r)

	// Select the task with the given id, and convert to JSON.
	result := a.DB.First(&task, "id = ?", vars["id"])
	if result.RowsAffected == 0 {
		http.Error(w, "error: id not found in DataBase", http.StatusNotFound)
		//w.WriteHeader(404)
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
		w.WriteHeader(400)
		return
	}

	rows := a.DB.Create(&newTask).RowsAffected
	log.Println("Added rows: ", rows)

	// Создаем json для ответа
	type ResponseId struct {
		Id int `json:"id"`
	}

	taskJSON, err := json.Marshal(ResponseId{Id: newTask.Id})
	if err != nil {
		http.Error(w, "error: not create task", http.StatusInternalServerError)
		//w.WriteHeader(500) // код ошибки 500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(taskJSON))
}

func (a *App) getTaskHandlerByTag(w http.ResponseWriter, r *http.Request) {
	var task Task
	vars := mux.Vars(r)

	// Select the task with the given id, and convert to JSON.
	err := a.DB.Where(&task, "tags = ?", vars["tag"])
	if err.Error != nil {
		http.Error(w, "error: tag not found in DataBase", http.StatusNotFound)
		//w.WriteHeader(404)
		return
	}
	taskJSON, _ := json.Marshal(task)

	// Write to HTTP response.
	w.WriteHeader(200)
	w.Write([]byte(taskJSON))
}

func (a *App) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Delete the task with the given id.
	err := a.DB.Where("id = ?", vars["id"]).Delete(Task{})
	if err.Error != nil {
		http.Error(w, "error: id not found in DataBase", http.StatusNotFound)
		//w.WriteHeader(404)
		return
	}

	// Write to HTTP response.
	w.WriteHeader(204)
}

func (a *App) deleteAllTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Delete all tasks.
	err := a.DB.Exec("DELETE FROM tasks")
	if err.Error != nil {
		http.Error(w, "error: not create task", http.StatusInternalServerError)
		//w.WriteHeader(500)
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

	r.HandleFunc("/task/tag/{tag}", a.getTaskHandlerByTag).Methods("GET")

	http.Handle("/", r)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
