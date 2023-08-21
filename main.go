package main

import (
	"fmt"
	"log"
	"net/http"

	// "project-manager/controllers"
	"project-manager/routes"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	fmt.Printf("Server running at port 8080...\n")

	routes.SetListRoutes(r)
	routes.SetCardRoutes(r)
	routes.SetChecklistRoutes(r)
	routes.SetItemRoutes(r)

	http.Handle("/", r)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
