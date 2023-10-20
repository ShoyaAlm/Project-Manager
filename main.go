package main

import (
	"fmt"
	"log"

	"net/http"
	"project-manager/routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	fmt.Printf("Server running at port 8080...\n")

	routes.SetListRoutes(r)
	routes.SetCardRoutes(r)
	routes.SetChecklistRoutes(r)
	routes.SetItemRoutes(r)
	routes.SetMemberRoutes(r)
	routes.SetSignInUpRoutes(r)
	routes.SetNotifRoutes(r)

	
	http.Handle("/", cors(r))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
