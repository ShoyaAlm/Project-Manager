package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"project-manager/model"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "host=localhost user=postgres password=2220819 dbname=project_manager sslmode=disable")
	// db, err = sql.Open("postgres", "host=localhost user=postgres password=2220819 dbname=projectmanager sslmode=disable")
	fmt.Printf("Connecting to the server...")
	if err != nil {
		fmt.Printf("error : %s", err)
		log.Fatal(err)
	}
}
func GetAllLists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, name FROM lists")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch lists from database, %s", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lists []*model.List
	for rows.Next() {
		list := &model.List{}
		err := rows.Scan(&list.ID, &list.Name)
		if err != nil {
			http.Error(w, "Error scanning rows", http.StatusInternalServerError)
			return
		}
		lists = append(lists, list)
	}

	jsonData, err := json.Marshal(lists)
	if err != nil {
		http.Error(w, "Failed to marshal lists data", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

// w.Header().Set("Content-Type", "application/json")

// // Marshal the projects slice into JSON
// jsonData, err := json.Marshal(lists)
// if err != nil {
// 	http.Error(w, "Failed to marshal lists data", http.StatusInternalServerError)
// 	return
// }

// // Write the JSON data to the response
// w.Write(jsonData)

func GetAList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT id, name FROM lists WHERE id = $1", listID)
	list := &model.List{}
	err = row.Scan(&list.ID, &list.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "List not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to fetch list data, %s", err), http.StatusInternalServerError)
		}
		return
	}

	jsonData, err := json.Marshal(list)
	if err != nil {
		http.Error(w, "Failed to marshal list data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	// Parse the list ID from the request URL
	// vars := mux.Vars(r)
	// listID, err := strconv.Atoi(vars["id"])
	// if err != nil {
	// 	http.Error(w, "Invalid list ID", http.StatusBadRequest)
	// 	return
	// }

	// // Find the list with the given ID in your 'lists' slice
	// var foundList *model.List
	// for _, list := range lists {
	// 	if list.ID == listID {
	// 		foundList = list
	// 		break
	// 	}
	// }

	// if foundList == nil {
	// 	http.Error(w, "List not found", http.StatusNotFound)
	// 	return
	// }

	// // Marshal the found list into JSON
	// jsonData, err := json.Marshal(foundList)
	// if err != nil {
	// 	http.Error(w, "Failed to marshal list data", http.StatusInternalServerError)
	// 	return
	// }

	// // Write the JSON data to the response
	// w.Header().Set("Content-Type", "application/json")
	// w.Write(jsonData)

}

// var lists = []*model.List{
// 	{
// 		ID:    1,
// 		Name:  "sample project 1",
// 		Cards: []*model.Card{},
// 	},
// 	{
// 		ID:    2,
// 		Name:  "sample project 2",
// 		Cards: []*model.Card{},
// 	},
// }

// func CreateList(w http.ResponseWriter, r *http.Request) {

// 	var newList model.List

// 	err := json.NewDecoder(r.Body).Decode(&newList)
// 	if err != nil {
// 		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
// 		return
// 	}

// 	newList.ID = len(lists) + 1
// 	if newList.Cards == nil {
// 		newList.Cards = []*model.Card{}
// 	}
// 	lists = append(lists, &newList)

// 	w.WriteHeader(http.StatusCreated)
// 	fmt.Fprintf(w, "List created successfully")

// }

// func UpdateList(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	listID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid list ID", http.StatusBadRequest)
// 		return
// 	}

// 	var updatedList model.List
// 	err = json.NewDecoder(r.Body).Decode(&updatedList)
// 	if err != nil {
// 		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
// 	}

// 	found := false
// 	for i, list := range lists {
// 		if list.ID == listID {
// 			updatedList.ID = list.ID
// 			updatedList.Cards = list.Cards
// 			lists[i] = &updatedList
// 			found = true
// 			break
// 		}
// 	}

// 	if !found {
// 		http.Error(w, "List not found", http.StatusNotFound)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "List updated successfully")
// }

// func DeleteList(w http.ResponseWriter, r *http.Request) {
// 	// Parse the list ID from the request URL or request body
// 	vars := mux.Vars(r)
// 	listID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid list ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Find and remove the list with the given ID from your 'lists' slice
// 	found := false
// 	for i, list := range lists {
// 		if list.ID == listID {
// 			lists = append(lists[:i], lists[i+1:]...)
// 			found = true
// 			break
// 		}
// 	}

// 	if !found {
// 		http.Error(w, "List not found", http.StatusNotFound)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "List deleted successfully")
// }
