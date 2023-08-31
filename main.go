package main

import (
	"fmt"
	"log"

	"net/http"

	"project-manager/routes"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	fmt.Printf("Server running at port 8080...\n")

	routes.SetListRoutes(r)
	routes.SetCardRoutes(r)
	routes.SetChecklistRoutes(r)
	http.Handle("/", r)

	// err := http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	log.Printf("Failed to start server: %s\n", err)
	// }

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
