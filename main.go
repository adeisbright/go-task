package main

import (
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

func main() {
	fmt.Println("The Server has started")
	http.HandleFunc("/", indexHander)
	http.ListenAndServe(":8000", nil)
}
