package routes

import (
	"project-manager/controllers"

	"github.com/gorilla/mux"
)

func SetListRoutes(r *mux.Router) {
	r.HandleFunc("/api/lists", controllers.GetAllLists).Methods("GET")
	r.HandleFunc("/api/lists/{id:[0-9]+}", controllers.GetAList).Methods("GET")
	r.HandleFunc("/api/lists", controllers.CreateList).Methods("POST")
	r.HandleFunc("/api/lists/{id:[0-9]+}", controllers.UpdateAList).Methods("PATCH")
	r.HandleFunc("/api/lists/{id:[0-9]+}", controllers.DeleteAList).Methods("DELETE")
}

func SetCardRoutes(r *mux.Router) {
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards", controllers.GetAllCards).Methods("GET")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}", controllers.GetACard).Methods("GET")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards", controllers.CreateACard).Methods("POST")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}", controllers.UpdateCard).Methods("PATCH")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}", controllers.DeleteCard).Methods("DELETE")
}

// func SetChecklistRoutes(r *mux.Router) {
// 	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{id:[0-9]+}/checklists", controllers.GetAllChecklists).Methods("GET")
// 	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{id:[0-9]+}/checklists/{id:[0-9]+}", controllers.GetAChecklist).Methods("GET")
// 	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{id:[0-9]+}/checklists", controllers.CreateChecklist).Methods("POST")
// 	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{id:[0-9]+}/checklists/{id:[0-9]+}", controllers.UpdateChecklist).Methods("PATCH")
// 	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{id:[0-9]+}/checklists/{id:[0-9]+}", controllers.DeleteChecklist).Methods("DELETE")
// }

// func SetItemRoutes(r *mux.Router) {
// 	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{id:[0-9]+}/checklists/{id:[0-9]+}/items", controllers.GetAllItems).Methods("GET")
// 	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{id:[0-9]+}/checklists/{id:[0-9]+}/items/{id:[0-9]+}", controllers.GetAItem).Methods("GET")
// 	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{id:[0-9]+}/checklists/{id:[0-9]+}/items", controllers.CreateItem).Methods("POST")
// 	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{id:[0-9]+}/checklists/{id:[0-9]+}/items/{id:[0-9]+}", controllers.UpdateItem).Methods("PATCH")
// 	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{id:[0-9]+}/checklists/{id:[0-9]+}/items/{id:[0-9]+}", controllers.DeleteItem).Methods("DELETE")
// }
