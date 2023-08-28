package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	// "github.com/codegangsta/gin"
	"project-manager/model"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

func GetAllCards(w http.ResponseWriter, r *http.Request) {

}

func GetACard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	cardID, err := strconv.Atoi(vars["cardID"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	fmt.Printf("list ID %v, card ID %v ", listID, cardID)
	cardRow := db.QueryRow(`
	    SELECT
	        c.id AS card_id, c.name AS card_name, c.description AS card_description, c.dates as card_dates
	    FROM cards c
	    WHERE c.list_id = $1 AND c.id = $2`, listID, cardID)

	card := &model.Card{}
	var cardDates pq.StringArray
	err = cardRow.Scan(&card.ID, &card.Name, &card.Description, &cardDates)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Card not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to fetch card data, %s", err), http.StatusInternalServerError)
		}
		return
	}

	card.Dates = []string(cardDates)

	memberRows, err := db.Query(`
	SELECT m.id AS member_id, m.name AS member_name
	FROM members m
	WHERE m.card_id = $1`, cardID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch members for card, %s", err), http.StatusInternalServerError)
		return
	}

	defer memberRows.Close()

	for memberRows.Next() {
		var (
			memberID   int
			memberName string
		)
		err := memberRows.Scan(&memberID, &memberName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning member rows, %s", err), http.StatusInternalServerError)
			return
		}
		card.Members = append(card.Members, &model.Member{ID: memberID, Name: memberName})
	}

	checklistRows, err := db.Query(`
	SELECT c.id AS checklist_id, c.name AS checklist_name,
	i.id AS item_id, i.name AS item_name, i.due_date AS item_duedate, i.assigned_to AS item_assignedto 
	FROM checklists c
	LEFT JOIN items i ON c.id = i.id
	WHERE c.card_id = $1`, cardID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch checklists for card, %s", err), http.StatusInternalServerError)
		return
	}

	defer checklistRows.Close()

	for checklistRows.Next() {
		var (
			checklistName, itemName, itemDueDate string
			checklistID, itemID                  int
			itemAssignedTo                       pq.StringArray
		)
		err := checklistRows.Scan(&checklistID, &checklistName, &itemID, &itemName, &itemDueDate, &itemAssignedTo)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning checklist rows, %s", err), http.StatusInternalServerError)
			return
		}

		// appending checklists to a card
		card.Checklists = append(card.Checklists, &model.Checklist{
			ID: checklistID, Name: checklistName})

		// appending items to that specific checklist
		checklistIndex := -1
		for idx, checklist := range checklists {
			if checklist.ID == checklistID {
				checklistIndex = idx
			}
		}

		if checklistIndex != -1 {
			newItem := &model.Item{
				ID:         itemID,
				Name:       itemName,
				DueDate:    itemDueDate,
				AssignedTo: []string(itemAssignedTo),
			}
			card.Checklists[checklistIndex].Items = append(card.Checklists[checklistIndex].Items, newItem)
		}
	}

	jsonData, err := json.Marshal(card)
	if err != nil {
		http.Error(w, "failed to marshal card data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

// func CreateCard(w http.ResponseWriter, r *http.Request) {

// 	var newCard model.Card

// 	err := json.NewDecoder(r.Body).Decode(&newCard)
// 	if err != nil {
// 		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
// 		return
// 	}

// 	newCard.ID = len(cards) + 1
// 	if newCard.Checklists == nil {
// 		newCard.Checklists = []*model.Checklist{}
// 	}
// 	cards = append(cards, &newCard)

// 	w.WriteHeader(http.StatusCreated)
// 	fmt.Fprintf(w, "Card created successfully")

// }

// func UpdateCard(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	cardID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid card ID", http.StatusBadRequest)
// 		return
// 	}

// 	var updatedCard model.Card
// 	err = json.NewDecoder(r.Body).Decode(&updatedCard)
// 	if err != nil {
// 		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
// 	}

// 	found := false
// 	for i, card := range cards {
// 		if card.ID == cardID {
// 			updatedCard.ID = card.ID
// 			updatedCard.Checklists = card.Checklists
// 			cards[i] = &updatedCard
// 			found = true
// 			break
// 		}
// 	}

// 	if !found {
// 		http.Error(w, "Card not found", http.StatusNotFound)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "Card updated successfully")
// }

// func DeleteCard(w http.ResponseWriter, r *http.Request) {
// 	// Parse the card ID from the request URL or request body
// 	vars := mux.Vars(r)
// 	cardID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid list ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Find and remove the card with the given ID from your 'cards' slice
// 	found := false
// 	for i, card := range cards {
// 		if card.ID == cardID {
// 			cards = append(cards[:i], cards[i+1:]...)
// 			found = true
// 			break
// 		}
// 	}

// 	if !found {
// 		http.Error(w, "Card not found", http.StatusNotFound)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "Card deleted successfully")
// }
