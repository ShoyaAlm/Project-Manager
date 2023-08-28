package main

import (
	"fmt"
	"log"
	"time"

	"net/http"

	"project-manager/routes"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	fmt.Printf("Server running at port 8080...\n")

	routes.SetListRoutes(r)
	routes.SetCardRoutes(r)
	http.Handle("/", r)

	maxRetries := 10
	retryInterval := time.Second * 5
	for retries := 0; retries < maxRetries; retries++ {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Printf("Failed to start server: %s\n", err)
			log.Printf("Retrying in %s...\n", retryInterval)
			time.Sleep(retryInterval)
		} else {
			break
		}
	}

	log.Fatal("Max retries reached, unable to start the server.")

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
