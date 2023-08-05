package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project-manager/model"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello world")

}

var projects = []*model.Project{
	{
		ID:    1,
		Name:  "sample project 1",
		Tasks: []*model.Task{},
	},
	{
		ID:    2,
		Name:  "sample project 2",
		Tasks: []*model.Task{},
	},
}

func getAllProjects(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Marshal the projects slice into JSON
	jsonData, err := json.Marshal(projects)
	if err != nil {
		http.Error(w, "Failed to marshal projects data", http.StatusInternalServerError)
		return
	}

	// Write the JSON data to the response
	w.Write(jsonData)

}

func createProject(w http.ResponseWriter, r *http.Request) {

	var newProject model.Project

	err := json.NewDecoder(r.Body).Decode(&newProject)
	if err != nil {
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	newProject.ID = len(projects) + 1
	projects = append(projects, &newProject)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Project created successfully")

}

func main() {

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("Server running at port 8080...\n")
	http.HandleFunc("/api/projects", getAllProjects)
	http.HandleFunc("/api/projects/create-project", createProject)

	// fmt.Printf("new project : %v", newProject)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}

// newProject := model.Project{
// 	ID:    1,
// 	Name:  "Sample Project",
// 	Tasks: []*model.Task{},
// }

// Creating a new Task
// newTask := &Task{
//     ID:       1,
//     Name:     "Sample Task",
//     Category: Strategic, // You can directly use the constants to set the values
//     Priority: Priority2,
// }

// Assuming you have a task called newTask of type *Task
// newProject.Tasks = append(newProject.Tasks, newTask)

/*
	a project has : 1-name 2-startdate 3-duedate 4-tasks 5-members 6-
	each task has : 1-name 2-startdate 3-duedate 4-category 5-priority
*/
