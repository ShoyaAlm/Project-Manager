package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"project-manager/model"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

func GetAllItems(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	checklistID, err := strconv.Atoi(vars["checklistID"])
	if err != nil {
		http.Error(w, "Invalid checklist ID", http.StatusBadRequest)
		return
	}

	itemRows, err := db.Query("SELECT id, name, due_date, assigned_to FROM items WHERE checklist_id = $1", checklistID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch items, %s", err), http.StatusInternalServerError)
		return
	}

	defer itemRows.Close()

	var items []*model.Item

	for itemRows.Next() {
		var (
			itemID                      int
			itemName 			  		string
			itemStartDate, itemDueDate  time.Time
			itemDone					bool
			// itemAssignedTo        		pq.StringArray
		)

		err := itemRows.Scan(&itemID, &itemName, &itemStartDate, &itemDueDate, &itemDone)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning rows, %s", err), http.StatusInternalServerError)
			return
		}

		item := &model.Item{
			ID:         itemID,
			Name:       itemName,
			StartDate: 	itemStartDate,
			DueDate:    itemDueDate,
			Done: 		itemDone,
			// AssignedTo: itemAssignedTo,
		}

		items = append(items, item)	

	}

	jsonData, err := json.Marshal(items)
	if err != nil {
		http.Error(w, "Failed to marshal cbecklists data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

func GetAItem(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	checklistID, err := strconv.Atoi(vars["checklistID"])
	if err != nil {
		http.Error(w, "Invalid checklist ID", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.Atoi(vars["itemID"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// Fetch list details
	itemRow := db.QueryRow("SELECT id, name, start_date due_date, done FROM items WHERE id = $1 AND checklist_id = $2", itemID, checklistID)

	var (
		itemName					string
		itemStartDate, itemDueDate 	time.Time 
		itemDone 					bool
		// itemAssignedTo        		[]*model.Member
	)

	err = itemRow.Scan(&itemID, &itemName, &itemStartDate, &itemDueDate, &itemDone)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "item not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to fetch item data, %s", err), http.StatusInternalServerError)
		}
		return
	}

	item := &model.Item{
		ID:         itemID,
		Name:       itemName,
		DueDate:    itemDueDate,
		StartDate: 	itemStartDate,
		Done: 		itemDone,
		// AssignedTo: itemAssignedTo,
	}

	jsonData, err := json.Marshal(item)
	if err != nil {
		http.Error(w, "Failed to marshal checklist data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

func CreateItem(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	checklistID, err := strconv.Atoi(vars["checklistID"])
	if err != nil {
		http.Error(w, "Invalid checklist ID", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body, %s", err), http.StatusInternalServerError)
		return
	}

	var requestData struct {
		Name       string         `json:"name"`
		DueDate    string         `json:"duedate"`
		AssignedTo pq.StringArray `json:"assignedto"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON data, %s", err), http.StatusBadRequest)
		return
	}

	var newItemID int

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create item, %s", err), http.StatusInternalServerError)
		return
	}

	currentDate := time.Now()
	oneWeekLater := currentDate.AddDate(0, 0, 7)


	// Create a new card with non-null fields
	newItem := &model.Item{
		ID:         newItemID,
		Name:       requestData.Name,
		StartDate: 	currentDate,
		DueDate:    oneWeekLater,
		Done: 		false,
		AssignedTo: []*model.Member{},
	}

	err = db.QueryRow("INSERT INTO items (name, start_date, due_date, done, checklist_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		newItem.Name, newItem.StartDate, newItem.DueDate, newItem.Done, checklistID).Scan(&newItemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert items, %s", err), http.StatusInternalServerError)
		return
	}

	// Fetch the associated list
	checklistRow := db.QueryRow("SELECT id, name FROM checklists WHERE id = $1", checklistID)
	checklist := &model.Checklist{}
	err = checklistRow.Scan(&checklist.ID, &checklist.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch checklist data, %s", err), http.StatusInternalServerError)
		return
	}

	// Append the new card to the list's cards slice
	checklist.Items = append(checklist.Items, newItem)

	jsonData, err := json.Marshal(checklist)
	if err != nil {
		http.Error(w, "Failed to marshal list data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)

}

func UpdateItem(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["checklistID"])
	if err != nil {
		http.Error(w, "Invalid checklist ID", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.Atoi(vars["itemID"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Parse the JSON request body
	var requestData struct {
		Name       string 		`json:"name"`
		StartDate  time.Time 	`json:"startdate"`
		DueDate    time.Time 	`json:"duedate"`
		Done 	   bool			`json:"done"`
		AssignedTo string 		`json:"assignedto"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	// Update the list's name in the database
	_, err = db.Exec("UPDATE items SET name = $1 WHERE id = $2", requestData.Name, itemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update item, %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func DeleteItem(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	itemID, err := strconv.Atoi(vars["itemID"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM items WHERE id = $1", itemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete the item, %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}
