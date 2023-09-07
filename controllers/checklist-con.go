package controllers

import (
	"database/sql"
	"project-manager/model"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	// "project-manager/model"
	"strconv"

	"github.com/lib/pq"

	// "github.com/codegangsta/gin"
	"github.com/gorilla/mux"
	// "project-manager/model"
)

func GetAllChecklists(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	cardID, err := strconv.Atoi(vars["cardID"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	checklistRows, err := db.Query("SELECT id, name FROM checklists WHERE card_id = $1", cardID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch checklists, %s", err), http.StatusInternalServerError)
		return
	}

	defer checklistRows.Close()

	var checklists []*model.Checklist

	for checklistRows.Next() {
		var (
			checklistID   int
			checklistName string
		)

		err := checklistRows.Scan(&checklistID, &checklistName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning rows, %s", err), http.StatusInternalServerError)
			return
		}

		checklist := &model.Checklist{
			ID:    checklistID,
			Name:  checklistName,
			Items: []*model.Item{},
		}

		// Start looking for items inside every checklist of every card
		itemRows, err := db.Query(`SELECT i.id AS item_id, i.name AS item_name, i.due_date AS item_due_date, i.assigned_to AS item_assigned_to
					FROM items i
					WHERE i.checklist_id = $1`, checklistID)

		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch items for checklists, %s", err), http.StatusInternalServerError)
			return

		}

		defer itemRows.Close()

		var items []*model.Item

		for itemRows.Next() {
			var (
				itemID                int
				itemName, itemDueDate string
				itemAssignedTo        pq.StringArray
			)

			err := itemRows.Scan(&itemID, &itemName, &itemDueDate, &itemAssignedTo)

			if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning itemRows, %s", err), http.StatusInternalServerError)
				return
			}

			item := &model.Item{
				ID:         itemID,
				Name:       itemName,
				DueDate:    itemDueDate,
				AssignedTo: itemAssignedTo,
			}

			items = append(items, item)

		}

		checklist.Items = items

		checklists = append(checklists, checklist)

	}

	jsonData, err := json.Marshal(checklists)
	if err != nil {
		http.Error(w, "Failed to marshal cbecklists data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

func GetAChecklist(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	cardID, err := strconv.Atoi(vars["cardID"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	checklistID, err := strconv.Atoi(vars["checklistID"])
	if err != nil {
		http.Error(w, "Invalid checklist ID", http.StatusBadRequest)
		return
	}

	// Fetch list details
	checklistRow := db.QueryRow("SELECT id, name FROM checklists WHERE id = $1 AND card_id = $2", checklistID, cardID)

	var (
		checklistName string
	)

	err = checklistRow.Scan(&checklistID, &checklistName)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Checklist not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to fetch checklist data, %s", err), http.StatusInternalServerError)
		}
		return
	}

	checklist := &model.Checklist{
		ID:    checklistID,
		Name:  checklistName,
		Items: []*model.Item{},
	}

	// Start looking for items inside every checklist of every card
	itemRows, err := db.Query(`SELECT i.id AS item_id, i.name AS item_name, i.due_date AS item_due_date, i.assigned_to AS item_assigned_to
				FROM items i
				WHERE i.checklist_id = $1`, checklistID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch items for checklists inside cards, %s", err), http.StatusInternalServerError)
		return

	}

	defer itemRows.Close()

	var items []*model.Item

	for itemRows.Next() {
		var (
			itemID                int
			itemName, itemDueDate string
			itemAssignedTo        pq.StringArray
		)

		err := itemRows.Scan(&itemID, &itemName, &itemDueDate, &itemAssignedTo)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning itemRows, %s", err), http.StatusInternalServerError)
			return
		}

		item := &model.Item{
			ID:         itemID,
			Name:       itemName,
			DueDate:    itemDueDate,
			AssignedTo: itemAssignedTo,
		}

		items = append(items, item)

	}

	checklist.Items = items

	jsonData, err := json.Marshal(checklist)
	if err != nil {
		http.Error(w, "Failed to marshal checklist data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

func CreateChecklist(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	cardID, err := strconv.Atoi(vars["cardID"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body, %s", err), http.StatusInternalServerError)
		return
	}

	var requestData struct {
		Name  string        `json:"name"`
		Items []*model.Item `json:"items"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON data, %s", err), http.StatusBadRequest)
		return
	}

	var newChecklistID, newItemID int

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create checklist, %s", err), http.StatusInternalServerError)
		return
	}

	emptyItem := &model.Item{
		ID:         newItemID,
		Name:       "آیتم 1",
		DueDate:    "2023-09-20T00:00:00Z",
		AssignedTo: []string{"شخص 1", "شخص 2"},
	}

	// Create a new card with non-null fields
	newChecklist := &model.Checklist{
		ID:    newChecklistID,
		Name:  requestData.Name,
		Items: []*model.Item{emptyItem},
	}

	err = db.QueryRow("INSERT INTO checklists (name, card_id) VALUES ($1, $2) RETURNING id",
		newChecklist.Name, cardID).Scan(&newChecklistID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert checklist, %s", err), http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("INSERT INTO items (name, due_date, assigned_to, checklist_id) VALUES ($1, $2, $3, $4) RETURNING id",
		emptyItem.Name, emptyItem.DueDate, pq.Array(emptyItem.AssignedTo), newChecklistID).Scan(&newItemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert item into checklist, %s", err), http.StatusInternalServerError)
		return
	}

	// Fetch the associated list
	cardRow := db.QueryRow("SELECT id, name, description, dates FROM cards WHERE id = $1", cardID)
	card := &model.Card{}
	var datesArray pq.StringArray
	err = cardRow.Scan(&card.ID, &card.Name, &card.Description, &datesArray)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch card data, %s", err), http.StatusInternalServerError)
		return
	}

	card.Dates = []string(datesArray)

	// Append the new card to the list's cards slice
	card.Checklists = append(card.Checklists, newChecklist)

	jsonData, err := json.Marshal(card)
	if err != nil {
		http.Error(w, "Failed to marshal list data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)

}

func UpdateChecklist(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["cardID"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	checklistID, err := strconv.Atoi(vars["checklistID"])
	if err != nil {
		http.Error(w, "Invalid checklist ID", http.StatusBadRequest)
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
		Name string `json:"name"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	// Update the list's name in the database
	_, err = db.Exec("UPDATE checklists SET name = $1 WHERE id = $2", requestData.Name, checklistID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update checklist, %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func DeleteChecklist(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	checklistID, err := strconv.Atoi(vars["checklistID"])
	if err != nil {
		http.Error(w, "Invalid checklist ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM items WHERE checklist_id = $1", checklistID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete items of checklist, %s", err), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM checklists WHERE id = $1", checklistID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete targeted checklist, %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}
