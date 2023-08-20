package main

import (
	"fmt"
	"log"
	"net/http"
	"project-manager/controllers"

	"github.com/gorilla/mux"
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

func main() {

	r := mux.NewRouter()

	// fileServer := http.FileServer(http.Dir("./static"))
	// http.Handle("/", fileServer)
	// http.HandleFunc("/hello", helloHandler)

	fmt.Printf("Server running at port 8080...\n")
	r.HandleFunc("/api/lists", controllers.GetAllLists).Methods("GET")
	r.HandleFunc("/api/lists/{id:[0-9]+}", controllers.GetAList).Methods("GET")
	r.HandleFunc("/api/lists", controllers.CreateList).Methods("POST")
	r.HandleFunc("/api/lists/{id:[0-9]+}/delete", controllers.DeleteList).Methods("DELETE")
	r.HandleFunc("/api/lists/{id:[0-9]+}/update", controllers.UpdateList).Methods("PATCH")
	http.Handle("/", r)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
