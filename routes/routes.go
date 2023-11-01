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
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards", controllers.CreateCard).Methods("POST")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}", controllers.UpdateCard).Methods("PATCH")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}", controllers.DeleteCard).Methods("DELETE")
}

func SetChecklistRoutes(r *mux.Router) {
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}/checklists", controllers.GetAllChecklists).Methods("GET")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}/checklists/{checklistID:[0-9]+}", controllers.GetAChecklist).Methods("GET")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}/checklists", controllers.CreateChecklist).Methods("POST")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}/checklists/{checklistID:[0-9]+}", controllers.UpdateChecklist).Methods("PATCH")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}/checklists/{checklistID:[0-9]+}", controllers.DeleteChecklist).Methods("DELETE")
}

func SetItemRoutes(r *mux.Router) {
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}/checklists/{checklistID:[0-9]+}/items", controllers.GetAllItems).Methods("GET")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}/checklists/{checklistID:[0-9]+}/items/{itemID:[0-9]+}", controllers.GetAItem).Methods("GET")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}/checklists/{checklistID:[0-9]+}/items", controllers.CreateItem).Methods("POST")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}/checklists/{checklistID:[0-9]+}/items/{itemID:[0-9]+}", controllers.UpdateItem).Methods("PATCH")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}/checklists/{checklistID:[0-9]+}/items/{itemID:[0-9]+}", controllers.DeleteItem).Methods("DELETE")
}

func SetMemberRoutes(r *mux.Router) {
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}/members", controllers.GetAllMembers).Methods("GET")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}/members/{memberID:[0-9]+}", controllers.GetAMember).Methods("GET")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}/members", controllers.CreateMember).Methods("POST")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}/members/{memberID:[0-9]+}", controllers.UpdateMember).Methods("PATCH")
	r.HandleFunc("/api/lists/{id:[0-9]+}/cards/{cardID:[0-9]+}/members/{memberID:[0-9]+}", controllers.DeleteMember).Methods("DELETE")

}

func SetSignInUpRoutes(r *mux.Router) {
	r.HandleFunc("/api/users", controllers.GetAllUsers).Methods("GET")
	r.HandleFunc("/api/users/{userID:[0-9]+}", controllers.GetUser).Methods("GET")
	r.HandleFunc("/api/users/{userID:[0-9]+}", controllers.DeleteUser).Methods("DELETE")
	
	
	r.HandleFunc("/api/signup", controllers.SignUp).Methods("POST")
	
	r.HandleFunc("/api/login", controllers.Login).Methods("POST")
}


func SetNotifRoutes(r *mux.Router) {
	r.HandleFunc("/api/notifs/{userID:[0-9]+}", controllers.CreateNotif).Methods("POST")
	r.HandleFunc("/api/notifs", controllers.GetAllNotifs).Methods("GET")
	r.HandleFunc("/api/notifs/{notifID:[0-9]+}", controllers.DeleteNotif).Methods("DELETE")
	r.HandleFunc("/api/notifs/{userID:[0-9]+}", controllers.GetUserNotifs).Methods("GET")
	r.HandleFunc("/api/notifs/{userID:[0-9]+}", controllers.MarkAsReadNotifs).Methods("PATCH")
}