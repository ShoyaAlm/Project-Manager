package routes

import (
	"project-manager/controllers"

	"github.com/gorilla/mux"
)

func SetListRoutes(r *mux.Router) {
	r.HandleFunc("/api/lists", controllers.GetAllLists).Methods("GET")
	r.HandleFunc("/api/lists/{id:[0-9]+}", controllers.GetAList).Methods("GET")
	r.HandleFunc("/api/lists", controllers.CreateList).Methods("POST")
	r.HandleFunc("/api/lists/{id:[0-9]+}", controllers.UpdateList).Methods("PATCH")
	r.HandleFunc("/api/lists/{id:[0-9]+}", controllers.DeleteList).Methods("DELETE")
}

func SetCardRoutes(r *mux.Router) {
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards", controllers.GetAllCards).Methods("GET")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{id:[0-9]+}", controllers.GetACard).Methods("GET")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards", controllers.CreateCard).Methods("POST")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{id:[0-9]+}", controllers.UpdateCard).Methods("PATCH")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{id:[0-9]+}", controllers.DeleteCard).Methods("DELETE")
}
