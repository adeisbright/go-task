package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
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

var TodoTask = []Task{
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

var templateData = HomePageDescriptor{
	PageTitle: "Go Task : Home",
	Tasks:     TodoTask,
}

func indexHander(w http.ResponseWriter, r *http.Request) {
	templateFile, _ := template.ParseFiles("./template/index.html")
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
	if r.Method == http.MethodGet {
		//How to create a union of many types
		// How to create and assign a struct immediately

		json, _ := json.Marshal(TodoTask)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(json)
		return
	} else if r.Method == http.MethodDelete {
		fmt.Fprint(w, "This deletes the data")
	} else if r.Method == http.MethodPut {
		fmt.Fprint(w, "This updates the data")
	} else {
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
		var title = "Salomon" //json["title"]

		TodoTask = append(TodoTask, Task{
			Index:  4,
			Title:  "I add this manually",
			Status: false,
		})

		fmt.Println("The title is ", title)

		w.Write(json)
	}
}

//Todo :
// Add an Item to a List of Struct
// Remove an Item from a list of struct
// Find an Item in a list of struct
// Edit an Item in a list of struct
// Get a property from a JSON data or payload
// Handle diffeent kinds of HTTP request
// Middleware and Routing
// Add Authentication and Authorization
// Add Database Connection for storing data
// Setup a Proper project structure

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	json, _ := json.Marshal(TodoTask)
	w.Write(json)
}

func removeTasksHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	index := vars["taskId"]
	fmt.Fprint(w, "This will remove the item whose index is", index)
}

func updateTasksHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "This will update the item")
}

func main() {

	fmt.Println("The Server has started")
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	router := mux.NewRouter()
	router.HandleFunc("/", indexHander)
	router.HandleFunc("/tasks", getTasksHandler).Methods("GET")
	router.HandleFunc("/tasks", handleAddTask).Methods("POST")
	router.HandleFunc("/tasks/{taskId}", updateTasksHandler).Methods("PUT")
	router.HandleFunc("/tasks/{taskId}", removeTasksHandler).Methods("DELETE")

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}
