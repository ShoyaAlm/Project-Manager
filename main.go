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

	http.Handle("/", r)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}

// routes.SetCardRoutes(r)
// routes.SetChecklistRoutes(r)
// routes.SetItemRoutes(r)

// func main() {
// 	r := mux.NewRouter()

// 	fmt.Printf("Server running at port 8080...\n")

// 	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	err = client.Connect(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	routes.SetListRoutes(r, client)

// 	http.Handle("/", r)

// 	if err := http.ListenAndServe(":8080", nil); err != nil {
// 		log.Fatal(err)
// 	}
//}
