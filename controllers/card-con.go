package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	// "github.com/codegangsta/gin"
	"project-manager/model"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

func GetAllCards(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	cardRows, err := db.Query("SELECT id, name, description, dates FROM cards WHERE list_id = $1 ORDER BY position ASC", listID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch cards, %s", err), http.StatusInternalServerError)
		return
	}

	defer cardRows.Close()

	var cards []*model.Card

	for cardRows.Next() {
		var (
			cardID                    int
			cardName, cardDescription string
			cardDates                 pq.StringArray
		)

		err := cardRows.Scan(&cardID, &cardName, &cardDescription, &cardDates)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning rows, %s", err), http.StatusInternalServerError)
			return
		}

		var dates []time.Time
		
		for _, dateString := range cardDates {
			date, err := time.Parse("2006-01-02", dateString)
			if err != nil {
				// Handle the error, e.g., log it or return an error response
			}
			dates = append(dates, date)
		}

		card := &model.Card{
			ID:          cardID,
			Name:        cardName,
			Description: cardDescription,
			Dates:       dates,
			Members:     []*model.User{},
			Checklists:  []*model.Checklist{},
		}


		owner := &model.User{}

		ownerRow := db.QueryRow("SELECT u.id, u.name, u.email, u.bio FROM users u JOIN user_cards uc ON u.id = uc.user_id WHERE uc.card_id = $1", cardID)
		err = ownerRow.Scan(&owner.ID, &owner.Name, &owner.Email, &owner.Bio)
		if err != nil {
			owner = nil
		}

		card.Owner = owner



		// Start checking for checklists inside every card
		checklistRows, err := db.Query(`SELECT cl.id AS checklist_id, cl.name AS checklist_name
								FROM checklists cl
							   WHERE cl.card_id = $1`, cardID)

		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch checklists for card, %s", err), http.StatusInternalServerError)
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
				http.Error(w, fmt.Sprintf("Error scanning checklistRows, %s", err), http.StatusInternalServerError)
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
				http.Error(w, fmt.Sprintf("Failed to fetch items for checklists inside cards, %s", err), http.StatusInternalServerError)
				return

			}

			defer itemRows.Close()

			var items []*model.Item

			for itemRows.Next() {
				var (
					itemID                		int
					itemName				 	string
					itemStartDate, itemDueDate 	time.Time
					itemDone					bool
					// itemAssignedTo        		[]*model.Member
				)

				err := itemRows.Scan(&itemID, &itemName, &itemStartDate, &itemDueDate, &itemDone)

				if err != nil {
					http.Error(w, fmt.Sprintf("Error scanning itemRows, %s", err), http.StatusInternalServerError)
					return
				}

				item := &model.Item{
					ID:         itemID,
					Name:       itemName,
					DueDate:    itemDueDate,
					// AssignedTo: itemAssignedTo,
				}

				items = append(items, item)

			}

			checklist.Items = items

			checklists = append(checklists, checklist)

		}

		card.Checklists = checklists

		// Start looking for members inside every card
		memberRows, err := db.Query(`SELECT m.id AS member_id, m.name AS member_name, m.email AS member_email
							FROM members m
						WHERE m.card_id = $1`, cardID)

		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch members for card, %s", err), http.StatusInternalServerError)
			return

		}

		defer memberRows.Close()

		var members []*model.User

		for memberRows.Next() {

			var (
				memberID   int
				memberName, memberEmail string
			)

			err := memberRows.Scan(&memberID, &memberName, &memberEmail)

			if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning memberRows, %s", err), http.StatusInternalServerError)
				return
			}

			member := &model.User{
				ID:       memberID,     
				Name:     memberName,
				Email:    memberEmail,
			}

			members = append(members, member)

		}

		card.Members = members

		cards = append(cards, card)

	}

	jsonData, err := json.Marshal(cards)
	if err != nil {
		http.Error(w, "Failed to marshal list data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

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

	// Fetch list details
	cardRow := db.QueryRow("SELECT id, name, description, dates, label FROM cards WHERE id = $1 AND list_id = $2", cardID, listID)

	var (
		cardName, cardDescription string
		cardDates                 pq.StringArray
		cardLabel 				  sql.NullString
	)

	err = cardRow.Scan(&cardID, &cardName, &cardDescription, &cardDates, &cardLabel)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Card not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to fetch card data, %s", err), http.StatusInternalServerError)
		}
		return
	}

	var dates []time.Time

	for _, dateString := range cardDates {
		date, err := time.Parse("2006-01-02", dateString)
		if err != nil {
			// Handle the error, e.g., log it or return an error response
		}
		dates = append(dates, date)
	}



	var label *string

			if cardLabel.Valid {
				// Use cardLabel.String when the label is present
				actualLabel := cardLabel.String
				label = &actualLabel
			} else {
				// Handle the case where the label is NULL
				label = nil // or any other default value for a missing label
			}

	card := &model.Card{
		ID:          cardID,
		Name:        cardName,
		Description: cardDescription,
		Dates:       dates,
		Members:     []*model.User{},
		Checklists:  []*model.Checklist{},
		Label: 		 label,
	}


	owner := &model.User{}

	ownerRow := db.QueryRow("SELECT u.id, u.name, u.email, u.bio FROM users u JOIN user_cards uc ON u.id = uc.user_id WHERE uc.card_id = $1", cardID)
	err = ownerRow.Scan(&owner.ID, &owner.Name, &owner.Email, &owner.Bio)
	if err != nil {
		owner = nil
	}

	card.Owner = owner

	// Start checking for checklists inside every card
	checklistRows, err := db.Query(`SELECT cl.id AS checklist_id, cl.name AS checklist_name
			FROM checklists cl
			WHERE cl.card_id = $1`, cardID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch checklists for card, %s", err), http.StatusInternalServerError)
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
			http.Error(w, fmt.Sprintf("Error scanning checklistRows, %s", err), http.StatusInternalServerError)
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
			http.Error(w, fmt.Sprintf("Failed to fetch items for checklists inside cards, %s", err), http.StatusInternalServerError)
			return

		}

		defer itemRows.Close()

		var items []*model.Item

		for itemRows.Next() {
			var (
				itemID                		int
				itemName			  		string
				itemStartDate, itemDueDate 	time.Time
				itemDone					bool
				// itemAssignedTo        		[]*model.Member
			)

			err := itemRows.Scan(&itemID, &itemName, &itemStartDate, &itemDueDate, &itemDone)

			if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning itemRows, %s", err), http.StatusInternalServerError)
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

		checklist.Items = items

		checklists = append(checklists, checklist)

	}

	card.Checklists = checklists

	// Start looking for members inside every card
	memberRows, err := db.Query(`SELECT m.id AS member_id, m.name AS member_name, m.email AS member_email
			FROM members m
			WHERE m.card_id = $1`, cardID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch members for card, %s", err), http.StatusInternalServerError)
		return

	}

	defer memberRows.Close()

	var members []*model.User

	for memberRows.Next() {

		var (
			memberID   int
			memberName, memberEmail string	
		)

		err := memberRows.Scan(&memberID, &memberName, &memberEmail)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning memberRows, %s", err), http.StatusInternalServerError)
			return
		}

		member := &model.User{
			ID:   memberID,
			Name: memberName,
			Email: memberEmail,
		}

		members = append(members, member)

	}

	card.Members = members

	jsonData, err := json.Marshal(card)
	if err != nil {
		http.Error(w, "Failed to marshal card data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

func DeleteCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	cardID, err := strconv.Atoi(vars["cardID"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	

	_, err = db.Exec("DELETE FROM activities WHERE card_id = $1", cardID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete activities of card, %s", err), http.StatusInternalServerError)
		return
	}

	// Fetch the IDs of checklists associated with the card
	var checklistIDs []int
	rows, err := db.Query("SELECT id FROM checklists WHERE card_id = $1", cardID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch checklist IDs, %s", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var checklistID int
		if err := rows.Scan(&checklistID); err != nil {
			http.Error(w, fmt.Sprintf("Failed to scan checklist ID, %s", err), http.StatusInternalServerError)
			return
		}
		checklistIDs = append(checklistIDs, checklistID)
	}

	// Delete the items associated with the checklists
	for _, checklistID := range checklistIDs {

		// Fetch the list of item IDs associated with the checklist
		itemRows, err := db.Query("SELECT id FROM items WHERE checklist_id = $1", checklistID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch items of checklist, %s", err), http.StatusInternalServerError)
			return
		}

		// Move defer outside the loop
		defer itemRows.Close()

		// Iterate through each item and delete associated item_members
		for itemRows.Next() {
			var itemID int
			if err := itemRows.Scan(&itemID); err != nil {
				http.Error(w, fmt.Sprintf("Error scanning item ID, %s", err), http.StatusInternalServerError)
				return
			}

			// Delete item_members associated with the current item
			_, err := db.Exec("DELETE FROM item_members WHERE item_id = $1", itemID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to delete item_members for item %d, %s", itemID, err), http.StatusInternalServerError)
				return
			}
		}

		_, err = db.Exec("DELETE FROM items WHERE checklist_id = $1", checklistID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete items of checklist, %s", err), http.StatusInternalServerError)
			return
		}
	}

	_, err = db.Exec("DELETE FROM members WHERE card_id = $1", cardID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete members of card, %s", err), http.StatusInternalServerError)
		return
	}
	
	_, err = db.Exec("DELETE FROM checklists WHERE card_id = $1", cardID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete checklists of card, %s", err), http.StatusInternalServerError)
		return
	}

	// Delete the card itself
	_, err = db.Exec("DELETE FROM cards WHERE id = $1", cardID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete card, %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func CreateCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body, %s", err), http.StatusInternalServerError)
		return
	}

	type Date time.Time

	var requestData struct {
		Name        string             `json:"name"`
		UserID      int                `json:"user_id"`
		Username    string             `json:"username"`
		UserEmail   string             `json:"user_email"`
		Description string             `json:"description"`
		Dates       []Date             `json:"dates"`
		Members     []*model.User    `json:"members"`
		OwnerID     int                `json:"owner_id"`
		Owner       *model.User        `json:"owner"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON data, %s", err), http.StatusBadRequest)
		return
	}

	log.Printf("userID: %v", requestData.OwnerID)

	var newCardID, newMemberID int

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create card, %s", err), http.StatusInternalServerError)
		return
	}

	// Get the current time
	currentDate := time.Now()

	// Fetch the maximum position value from the cards table for the given list
	var maxCardPosition int
	err = db.QueryRow("SELECT COALESCE(MAX(position), 0) FROM cards WHERE list_id = $1", listID).Scan(&maxCardPosition)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch maximum card position, %s", err), http.StatusInternalServerError)
		return
	}

	// Increment the position value for the new card
	newCardPosition := maxCardPosition + 1

	// Calculate one month later
	oneMonthLater := currentDate.AddDate(0, 1, 0)

	dates := []time.Time{currentDate, oneMonthLater}

	// Create a new card with non-null fields
	newCard := &model.Card{
		ID:          newCardID,
		Name:        requestData.Name,
		Description: requestData.Description,
		Dates:       dates,
		Members:     requestData.Members,
		OwnerID:     requestData.OwnerID,
		Owner:       requestData.Owner,
	}

	owner := &model.User{
		ID:    requestData.OwnerID,
		Name:  requestData.Username,
		Email: requestData.UserEmail,
	}

	// Query the database to fetch owner details based on UserID
	ownerRow := db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", requestData.OwnerID)
	err = ownerRow.Scan(&owner.ID, &owner.Name, &owner.Email)
	if err != nil {
		log.Printf("Error fetching owner details: %v", err)
		http.Error(w, "Failed to fetch owner details", http.StatusInternalServerError)
		return
	}

	fmt.Printf("owner: %v", owner)

	newCard.Owner = owner

	err = db.QueryRow("INSERT INTO cards (name, description, dates, owner_id, list_id, position) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		newCard.Name, newCard.Description, pq.Array(dates), newCard.OwnerID, listID, newCardPosition).Scan(&newCardID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert card, %s", err), http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("INSERT INTO members (card_id, name, email) VALUES ($1, $2, $3) RETURNING id",
		newCardID, owner.Name, owner.Email).Scan(&newMemberID)
	if err != nil {
		log.Printf("Failed to insert members: %v", err)
		http.Error(w, fmt.Sprintf("Failed to insert members, %s", err), http.StatusInternalServerError)
		return
	}

	// Fetch the associated list
	listRow := db.QueryRow("SELECT id, name FROM lists WHERE id = $1", listID)
	list := &model.List{}
	err = listRow.Scan(&list.ID, &list.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch list data, %s", err), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO user_cards (user_id, card_id) VALUES ($1, $2)", 1, newCardID)
	if err != nil {
		log.Printf("Failed to insert user_card: %v", err)
		http.Error(w, "Failed to insert user_card", http.StatusInternalServerError)
		return
	}

	// Append the new card to the list's cards slice
	list.Cards = append(list.Cards, newCard)

	jsonData, err := json.Marshal(list)
	if err != nil {
		http.Error(w, "Failed to marshal list data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}


func UpdateCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	cardID, err := strconv.Atoi(vars["cardID"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
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
		Name        *string       `json:"name"`
		Description *string       `json:"description"`
		Dates       []time.Time   `json:"dates"`
		Label       *string       `json:"label"` // Added label field
	}
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	fmt.Printf("requestData : %v", requestData)

	var query string
	var args []interface{}

	// Check for label in requestData
	if requestData.Label != nil {
		query = "UPDATE cards SET label = $1 WHERE id = $2"
		args = []interface{}{*requestData.Label, cardID}
	} else {
		// Handle other update scenarios (name, description, dates) as before
		if requestData.Name != nil && requestData.Description != nil && len(requestData.Dates) > 0 {
			query = "UPDATE cards SET name = $1, description = $2, dates = $3 WHERE id = $4"
			args = []interface{}{*requestData.Name, *requestData.Description, pq.Array(requestData.Dates), cardID}
		} else if requestData.Name != nil && requestData.Description != nil {
			query = "UPDATE cards SET name = $1, description = $2 WHERE id = $3"
			args = []interface{}{*requestData.Name, *requestData.Description, cardID}
		} else if len(requestData.Dates) > 0 {
			query = "UPDATE cards SET dates = $1 WHERE id = $2"
			args = []interface{}{pq.Array(requestData.Dates), cardID}
		} else if requestData.Name != nil {
			query = "UPDATE cards SET name = $1 WHERE id = $2"
			args = []interface{}{*requestData.Name, cardID}
		} else if requestData.Description != nil {
			query = "UPDATE cards SET description = $1 WHERE id = $2"
			args = []interface{}{*requestData.Description, cardID}
		} else {
			http.Error(w, "No update data provided", http.StatusBadRequest)
			return
		}
	}

	// Update the card in the database
	_, err = db.Exec(query, args...)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update card, %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

