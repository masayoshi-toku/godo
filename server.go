package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type Todo struct {
	Id        int
	Content   string
	CreatedAt time.Time
}

var Db *gorm.DB

func init() {
	var err error
	Db, err = gorm.Open("postgres", "user=gwp dbname=gwp password=gwp sslmode=disable")
	if err != nil {
		panic(err)
	}
	Db.AutoMigrate(&Todo{})
}

func getAllToDo(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	Db.Order("Id asc").Find(&todos)

	t, _ := template.ParseFiles("index.html")
	t.Execute(w, todos)
}

func newToDo(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("form.html")
	t.Execute(w, nil)
}

func createToDo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	todo := Todo{Content: r.PostForm["content"][0]}
	Db.Create(&todo)

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("location", "http://127.0.0.1:8080/")
	w.WriteHeader(http.StatusSeeOther)
}

func editToDo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var todo Todo
	Db.Find(&todo, r.Form["id"][0])

	t, _ := template.ParseFiles("edit.html")
	t.Execute(w, &todo)
}

func updateToDo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var todo Todo
	Db.Find(&todo, r.PostForm["id"][0])
	Db.Model(&todo).Update("Content", r.PostForm["content"][0])

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("location", "http://127.0.0.1:8080/")
	w.WriteHeader(http.StatusSeeOther)
}

func deleteToDo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var todo Todo
	Db.Find(&todo, r.PostForm["id"][0])
	Db.Delete(&todo)

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("location", "http://127.0.0.1:8080/")
	w.WriteHeader(http.StatusSeeOther)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		switch r.URL.Path {
		case "/":
			getAllToDo(w, r)
		case "/new":
			newToDo(w, r)
		case "/edit":
			editToDo(w, r)
		}
	case "POST":
		switch r.URL.Path {
		case "/":
			createToDo(w, r)
		case "/update":
			updateToDo(w, r)
		case "/delete":
			deleteToDo(w, r)
		}
	}
}

func main() {
	server := http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/", handleRequest)
	server.ListenAndServe()
}
