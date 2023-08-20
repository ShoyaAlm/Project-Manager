package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-manager/model"
	"strconv"

	// "github.com/codegangsta/gin"
	"github.com/gorilla/mux"
	// "project-manager/model"
)

var lists = []*model.List{
	{
		ID:    1,
		Name:  "sample project 1",
		Cards: []*model.Card{},
	},
	{
		ID:    2,
		Name:  "sample project 2",
		Cards: []*model.Card{},
	},
}

func GetAllLists(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Marshal the projects slice into JSON
	jsonData, err := json.Marshal(lists)
	if err != nil {
		http.Error(w, "Failed to marshal lists data", http.StatusInternalServerError)
		return
	}

	// Write the JSON data to the response
	w.Write(jsonData)

}

func GetAList(w http.ResponseWriter, r *http.Request) {
	// Parse the list ID from the request URL
	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	// Find the list with the given ID in your 'lists' slice
	var foundList *model.List
	for _, list := range lists {
		if list.ID == listID {
			foundList = list
			break
		}
	}

	if foundList == nil {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	}

	// Marshal the found list into JSON
	jsonData, err := json.Marshal(foundList)
	if err != nil {
		http.Error(w, "Failed to marshal list data", http.StatusInternalServerError)
		return
	}

	// Write the JSON data to the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func CreateList(w http.ResponseWriter, r *http.Request) {

	var newList model.List

	err := json.NewDecoder(r.Body).Decode(&newList)
	if err != nil {
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	newList.ID = len(lists) + 1
	if newList.Cards == nil {
		newList.Cards = []*model.Card{}
	}
	lists = append(lists, &newList)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "List created successfully")

}

func UpdateList(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	var updatedList model.List
	err = json.NewDecoder(r.Body).Decode(&updatedList)
	if err != nil {
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
	}

	found := false
	for i, list := range lists {
		if list.ID == listID {
			updatedList.ID = list.ID
			updatedList.Cards = list.Cards
			lists[i] = &updatedList
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "List updated successfully")
}

func DeleteList(w http.ResponseWriter, r *http.Request) {
	// Parse the list ID from the request URL or request body
	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	// Find and remove the list with the given ID from your 'lists' slice
	found := false
	for i, list := range lists {
		if list.ID == listID {
			lists = append(lists[:i], lists[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "List deleted successfully")
}
