package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

type HomePageDescriptor struct {
	PageTitle string
	Tasks     []Task
}

type Task struct {
	Index  int
	Title  string
	Status bool
}

var TodoTask = []Task{}

func indexHander(w http.ResponseWriter, r *http.Request) {
	templateFile, _ := template.ParseFiles("./template/index.html")

	TodoTask = []Task{
		{
			Index:  1,
			Title:  "Build Go Web Service",
			Status: true,
		},
		{
			Index:  2,
			Title:  "Automate Server Provisioning with Terraform",
			Status: false,
		},
		{
			Index:  3,
			Title:  "Build DSA into my Fabrics",
			Status: false,
		},
	}

	templateData := HomePageDescriptor{
		PageTitle: "Go Task : Home",
		Tasks:     TodoTask,
	}

	templateFile.Execute(w, templateData)
}

type TaskPayload struct {
	Title string `json:"title"`
}

func HttpErrorHandler(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := make(map[string]string)
	response["message"] = message

	json, _ := json.Marshal(response)
	w.Write(json)
}

func handleAddTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HttpErrorHandler(w, "Bad Request Method", http.StatusMethodNotAllowed)
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		HttpErrorHandler(w, "Unsupported Media Type", http.StatusNotAcceptable)
	}

	var unMarshalError *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var task = TaskPayload{}

	err := decoder.Decode(&task)
	if err != nil {
		if errors.As(err, &unMarshalError) {
			HttpErrorHandler(w, unMarshalError.Field, http.StatusBadRequest)
		} else {
			HttpErrorHandler(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		}
		return
	}

	json, _ := json.Marshal(task)
	w.Write(json)
}

func main() {
	fmt.Println("The Server has started")
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.HandleFunc("/", indexHander)
	http.HandleFunc("/tasks", handleAddTask)
	http.ListenAndServe(":8000", nil)
}
