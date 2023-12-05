package controllers

import (
	"database/sql"
	"project-manager/model"
	"time"

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
		itemRows, err := db.Query(`SELECT i.id AS item_id, i.name AS item_name, i.start_date AS item_start_date, i.due_date AS item_due_date, i.done AS item_done
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
				itemID                		int
				itemName 			 		string
				itemStartDate, itemDueDate 	time.Time
				itemDone 			 		bool
				itemAssignedTo        		[]*model.Member
			)

			err := itemRows.Scan(&itemID, &itemName, &itemStartDate, &itemDueDate, &itemDone)

			if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning itemRows, %s", err), http.StatusInternalServerError)
				return
			}

			item := &model.Item{
				ID:         	itemID,
				Name:       	itemName,
				StartDate:    	itemStartDate,
				DueDate:    	itemDueDate,
				Done: 			itemDone,	
				AssignedTo: 	itemAssignedTo,
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
			itemID                		int
			itemName 			 		string
			itemStartDate, itemDueDate 	time.Time
			itemDone 			 		bool
			itemAssignedTo        		[]*model.Member
		)

		err := itemRows.Scan(&itemID, &itemName, &itemStartDate, &itemDueDate, &itemDone)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning itemRows, %s", err), http.StatusInternalServerError)
			return
		}

		item := &model.Item{
			ID:         	itemID,
			Name:       	itemName,
			StartDate:    	itemStartDate,
			DueDate:    	itemDueDate,
			Done: 			itemDone,	
			AssignedTo: 	itemAssignedTo,
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
		Name string `json:"name"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON data, %s", err), http.StatusBadRequest)
		return
	}

	var newChecklistID int

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create checklist, %s", err), http.StatusInternalServerError)
		return
	}

	// Fetch the maximum position value from the checklists table for the given card
	var maxChecklistPosition int
	err = db.QueryRow("SELECT COALESCE(MAX(position), 0) FROM checklists WHERE card_id = $1", cardID).Scan(&maxChecklistPosition)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch maximum checklist position, %s", err), http.StatusInternalServerError)
		return
	}

	// Increment the position value for the new checklist
	newChecklistPosition := maxChecklistPosition + 1

	// Create a new checklist with non-null fields
	newChecklist := &model.Checklist{
		ID:       newChecklistID,
		Name:     requestData.Name,
		Items:    []*model.Item{}, // No items for now
		// Position: newChecklistPosition,
	}

	err = db.QueryRow("INSERT INTO checklists (name, card_id, position) VALUES ($1, $2, $3) RETURNING id",
		newChecklist.Name, cardID, newChecklistPosition).Scan(&newChecklistID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert checklist, %s", err), http.StatusInternalServerError)
		return
	}

	// Fetch the associated card
	cardRow := db.QueryRow("SELECT id, name, description, dates FROM cards WHERE id = $1", cardID)
	card := &model.Card{}
	var datesArray pq.StringArray
	err = cardRow.Scan(&card.ID, &card.Name, &card.Description, &datesArray)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch card data, %s", err), http.StatusInternalServerError)
		return
	}

	var dates []time.Time

	for _, dateString := range datesArray {
		date, err := time.Parse("2006-01-02", dateString)
		if err != nil {
			// Handle the error, e.g., log it or return an error response
		}
		dates = append(dates, date)
	}

	card.Dates = dates

	// Append the new checklist to the card's checklists slice
	card.Checklists = append(card.Checklists, newChecklist)

	jsonData, err := json.Marshal(card)
	if err != nil {
		http.Error(w, "Failed to marshal card data", http.StatusInternalServerError)
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



func UpdateChecklistOrder(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusInternalServerError)
        return
    }

    type UpdateChecklistOrderRequest struct {
        ChecklistOrder []int `json:"checklistOrder"`
    }

    var requestData UpdateChecklistOrderRequest
    err = json.Unmarshal(body, &requestData)
    if err != nil {
        http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
        return
    }

    if len(requestData.ChecklistOrder) > 0 {
        tx, err := db.Begin()
        if err != nil {
            http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
            return
        }
        defer tx.Rollback()

        for i, checklistID := range requestData.ChecklistOrder {
            _, err := tx.Exec("UPDATE checklists SET position = $1 WHERE id = $2", i, checklistID)
            if err != nil {
                fmt.Printf("error: %v", err)
                http.Error(w, "Failed to update checklist order", http.StatusInternalServerError)
                return
            }
        }

        err = tx.Commit()
        if err != nil {
            http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
            return
        }
    }

    w.WriteHeader(http.StatusOK)
}

func UpdateItemOrder(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
	  http.Error(w, "Failed to read request body", http.StatusInternalServerError)
	  return
	}


	type UpdateItemOrderRequest struct {
		ChecklistID int   `json:"checklistId"`
		ItemOrder   []int `json:"itemOrder"`
	  }
  
	var requestData UpdateItemOrderRequest
	err = json.Unmarshal(body, &requestData)
	if err != nil {
	  http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
	  return
	}
  
	if len(requestData.ItemOrder) > 0 {
	  tx, err := db.Begin()
	  if err != nil {
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	  }
	  defer tx.Rollback()
  
	  for i, itemID := range requestData.ItemOrder {
		_, err := tx.Exec("UPDATE items SET position = $1 WHERE id = $2", i, itemID)
		if err != nil {
			fmt.Printf("error: %v", err)
			http.Error(w, "Failed to update item order", http.StatusInternalServerError)
		//   fmt.Printf("error: %v", err)
		  return
		}
	  }
  
	  err = tx.Commit()
	  if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	  }
	}
  
	w.WriteHeader(http.StatusOK)
  }
  