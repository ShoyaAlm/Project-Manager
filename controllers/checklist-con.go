package controllers

import (
	"project-manager/model"

	"encoding/json"
	"fmt"
	"net/http"

	"strconv"

	"github.com/gorilla/mux"
)

var checklists = []*model.Checklist{
	{
		ID:    1,
		Name:  "checklist 1",
		Items: []*model.Item{},
	},
	{
		ID:    2,
		Name:  "checklist 2",
		Items: []*model.Item{},
	},
}

func GetAllChecklists(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Marshal the projects slice into JSON
	jsonData, err := json.Marshal(checklists)
	if err != nil {
		http.Error(w, "Failed to marshal lists data", http.StatusInternalServerError)
		return
	}

	// Write the JSON data to the response
	w.Write(jsonData)

}

func GetAChecklist(w http.ResponseWriter, r *http.Request) {
	// Parse the list ID from the request URL
	vars := mux.Vars(r)
	checklistID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	// Find the list with the given ID in your 'lists' slice
	var foundChecklist *model.Checklist
	for _, checklist := range checklists {
		if checklist.ID == checklistID {
			foundChecklist = checklist
			break
		}
	}

	if foundChecklist == nil {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	}

	// Marshal the found list into JSON
	jsonData, err := json.Marshal(foundChecklist)
	if err != nil {
		http.Error(w, "Failed to marshal list data", http.StatusInternalServerError)
		return
	}

	// Write the JSON data to the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func CreateChecklist(w http.ResponseWriter, r *http.Request) {

	var newChecklist model.Checklist

	err := json.NewDecoder(r.Body).Decode(&newChecklist)
	if err != nil {
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	newChecklist.ID = len(checklists) + 1
	if newChecklist.Items == nil {
		newChecklist.Items = []*model.Item{}
	}
	checklists = append(checklists, &newChecklist)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "List created successfully")

}

func UpdateChecklist(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	checklistID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	var updatedChecklist model.Checklist
	err = json.NewDecoder(r.Body).Decode(&updatedChecklist)
	if err != nil {
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
	}

	found := false
	for i, checklist := range checklists {
		if checklist.ID == checklistID {
			updatedChecklist.ID = checklist.ID
			updatedChecklist.Items = checklist.Items
			checklists[i] = &updatedChecklist
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

func DeleteChecklist(w http.ResponseWriter, r *http.Request) {
	// Parse the list ID from the request URL or request body
	vars := mux.Vars(r)
	checklistID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	// Find and remove the list with the given ID from your 'lists' slice
	found := false
	for i, checklist := range checklists {
		if checklist.ID == checklistID {
			checklists = append(checklists[:i], checklists[i+1:]...)
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
