package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-manager/model"
	"strconv"

	"github.com/gorilla/mux"
)

var items = []*model.Item{
	{
		ID:         1,
		Name:       "item 1",
		DueDate:    "september 20th",
		AssignedTo: []*model.Member{},
	},
	{
		ID:         2,
		Name:       "item 2",
		DueDate:    "october 5th",
		AssignedTo: []*model.Member{},
	},
}

func GetAllItems(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Marshal the projects slice into JSON
	jsonData, err := json.Marshal(items)
	if err != nil {
		http.Error(w, "Failed to marshal items data", http.StatusInternalServerError)
		return
	}

	// Write the JSON data to the response
	w.Write(jsonData)

}

func GetAItem(w http.ResponseWriter, r *http.Request) {
	// Parse the list ID from the request URL
	vars := mux.Vars(r)
	itemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	// Find the Item with the given ID in your 'items' slice
	var foundItem *model.Item
	for _, item := range items {
		if item.ID == itemID {
			foundItem = item
			break
		}
	}

	if foundItem == nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	// Marshal the found list into JSON
	jsonData, err := json.Marshal(foundItem)
	if err != nil {
		http.Error(w, "Failed to marshal item data", http.StatusInternalServerError)
		return
	}

	// Write the JSON data to the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {

	var newItem model.Item

	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	newItem.ID = len(lists) + 1
	if newItem.AssignedTo == nil {
		newItem.AssignedTo = []*model.Member{}
	}
	items = append(items, &newItem)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Item created successfully")

}

func UpdateItem(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	itemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var updatedItem model.Item
	err = json.NewDecoder(r.Body).Decode(&updatedItem)
	if err != nil {
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
	}

	found := false
	for i, item := range items {
		if item.ID == itemID {
			updatedItem.ID = item.ID
			updatedItem.AssignedTo = item.AssignedTo
			items[i] = &updatedItem
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Item updated successfully")
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	// Parse the list ID from the request URL or request body
	vars := mux.Vars(r)
	itemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// Find and remove the list with the given ID from your 'lists' slice
	found := false
	for i, item := range items {
		if item.ID == itemID {
			items = append(items[:i], items[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Item deleted successfully")
}
