package main

import (
	"bytes"
	"encoding/base64"
	"html/template"
	"image/jpeg"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type Todo struct {
	Id        int
	Content   string
	Image     []byte
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
	r.ParseMultipartForm(5242880) // 5MB
	file, _, _ := r.FormFile("image")
	defer file.Close()

	img, _ := jpeg.Decode(file)
	buffer := new(bytes.Buffer)
	jpeg.Encode(buffer, img, nil)
	imageBytes := buffer.Bytes()

	todo := Todo{Content: r.PostFormValue("content"), Image: imageBytes}
	Db.Create(&todo)

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("location", "http://127.0.0.1:8080/")
	w.WriteHeader(http.StatusSeeOther)
}

func showToDo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	Db.Find(&todo, r.FormValue("id"))

	imageStr := base64.StdEncoding.EncodeToString(todo.Image)

	todoData := map[string]interface{}{
		"todo":     todo,
		"imageStr": imageStr,
	}

	t, _ := template.ParseFiles("show.html")
	t.Execute(w, todoData)
}

func editToDo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	Db.Find(&todo, r.FormValue("id"))

	imageStr := base64.StdEncoding.EncodeToString(todo.Image)

	todoData := map[string]interface{}{
		"todo":     todo,
		"imageStr": imageStr,
	}

	t, _ := template.ParseFiles("edit.html")
	t.Execute(w, todoData)
}

func updateToDo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	Db.Find(&todo, r.PostFormValue("id"))

	r.ParseMultipartForm(5242880) // 5MB
	file, _, _ := r.FormFile("image")
	defer file.Close()

	img, _ := jpeg.Decode(file)
	buffer := new(bytes.Buffer)
	jpeg.Encode(buffer, img, nil)
	imageBytes := buffer.Bytes()

	Db.Model(&todo).Updates(Todo{Content: r.PostFormValue("content"), Image: imageBytes})

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("location", "http://127.0.0.1:8080/")
	w.WriteHeader(http.StatusSeeOther)
}

func deleteToDo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	Db.Find(&todo, r.PostFormValue("id"))
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
		case "/show":
			showToDo(w, r)
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
