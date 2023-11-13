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

	cardRows, err := db.Query("SELECT id, name, description, dates FROM cards WHERE list_id = $1", listID)
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
	cardRow := db.QueryRow("SELECT id, name, description, dates FROM cards WHERE id = $1 AND list_id = $2", cardID, listID)

	var (
		cardName, cardDescription string
		cardDates                 pq.StringArray
	)

	err = cardRow.Scan(&cardID, &cardName, &cardDescription, &cardDates)
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

	_, err = db.Exec("DELETE FROM members WHERE card_id = $1", cardID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete members of card, %s", err), http.StatusInternalServerError)
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
		_, err := db.Exec("DELETE FROM items WHERE checklist_id = $1", checklistID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete items of checklist, %s", err), http.StatusInternalServerError)
			return
		}
	}

	_, err = db.Exec("DELETE FROM checklists WHERE card_id = $1", cardID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete checklists of card, %s", err), http.StatusInternalServerError)
		return
	}


	// Delete the list and related data
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
		UserID 		int			   	   `json:"user_id"`
		Username 	string			   `json:"username"`
		UserEmail 	string			   `json:"user_email"`
		Description string             `json:"description"`
		Dates       []Date             `json:"dates"`
		Checklists  []*model.Checklist `json:"checklists"`
		Members     []*model.Member    `json:"members"`
		OwnerID 	int 	    	   `json:"owner_id"`
		Owner 		*model.User 	   `json:"owner"`
	}

	
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON data, %s", err), http.StatusBadRequest)
		return
	}



	log.Printf("userID : %v", requestData.OwnerID)

	var newCardID, newChecklistID, newItemID, newMemberID int

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create card, %s", err), http.StatusInternalServerError)
		return
	}


	// Get the current time
	currentDate := time.Now()
	oneWeekLater := currentDate.AddDate(0, 0, 7)

	emptyItem := &model.Item{
		ID:         newItemID,
		Name:       "آیتم 1",
		StartDate:  currentDate,
		DueDate:    oneWeekLater,
		AssignedTo: []*model.Member{},
		Done: 		false,	
	}

	emptyChecklist := &model.Checklist{
		ID:    newChecklistID,
		Name:  "چکلیست جدید",
		Items: []*model.Item{emptyItem},
	}

	emptyMember := &model.User{
		ID:   newMemberID,
		Name: requestData.Username,
		Email: requestData.UserEmail,
	}



	// Calculate one month later
	oneMonthLater := currentDate.AddDate(0, 1, 0)

	dates := []time.Time{currentDate, oneMonthLater}
	// Create a new card with non-null fields
	newCard := &model.Card{
		ID:          newCardID,
		Name:        requestData.Name,
		Description: "توضیحات",
		Dates:       dates,
		Checklists:  []*model.Checklist{emptyChecklist},      // Initialize as empty slice
		Members:     []*model.User{emptyMember},
		OwnerID: 	 requestData.OwnerID,
		Owner:       &model.User{},
	}


	owner := &model.User{
		ID: requestData.OwnerID,
		Name: requestData.Username,
		Email: requestData.UserEmail,
	}



	// Query the database to fetch owner details based on UserID
	ownerRow := db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", requestData.OwnerID)
	err = ownerRow.Scan(&owner.ID, &owner.Name, &owner.Email,)
	if err != nil {
		log.Printf("Error fetching owner details: %v", err)
		http.Error(w, "Failed to fetch owner details", http.StatusInternalServerError)
		return
	}
	// ownerRow := db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", requestData.UserID)
	// err = ownerRow.Scan(&owner.ID, &owner.Name, &owner.Email)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		// Handle the case where no matching owner was found
	// 		log.Printf("Owner not found for UserID: %v", requestData.UserID)
	// 		// You can set a default owner or handle the situation as needed
	// 		// For example:
	// 		// owner = &model.User{ID: 0, Name: "Default Owner", Email: "default@example.com"}
	// 	} else {
	// 		log.Printf("Error fetching owner details: %v", err)
	// 		http.Error(w, "Failed to fetch owner details", http.StatusInternalServerError)
	// 		return
	// 	}
	// }
	
	fmt.Printf("owner : %v", owner)

	newCard.Owner = owner


	err = db.QueryRow("INSERT INTO cards (name, description, dates, owner_id, list_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		newCard.Name, newCard.Description, pq.Array(dates), newCard.OwnerID, listID).Scan(&newCardID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert card, %s", err), http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("INSERT INTO checklists (name, card_id) VALUES ($1, $2) RETURNING id",
		emptyChecklist.Name, newCardID).Scan(&newChecklistID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert checklists, %s", err), http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("INSERT INTO items (name, start_date, due_date, done, checklist_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		emptyItem.Name, emptyItem.StartDate, emptyItem.DueDate, emptyItem.Done, newChecklistID).Scan(&newItemID)
	if err != nil {
		log.Printf("Failed to insert items: %v", err)
		http.Error(w, fmt.Sprintf("Failed to insert items, %s", err), http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("INSERT INTO members (card_id, name, email) VALUES ($1, $2, $3) RETURNING id",
	newCardID, emptyMember.Name, emptyMember.Email).Scan(&newMemberID)
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
		Name 		*string 		`json:"name"`
		Description *string 		`json:"description"`
		Dates		[]time.Time 	`json:"dates"` 
	}
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
	
	fmt.Printf("requestData : %v", requestData)

	var query string
	var args []interface{}

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


	// Update the list's name in the database
	_, err = db.Exec(query, args...)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update card, %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}


// func UpdateCard(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	_, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid list ID", http.StatusBadRequest)
// 		return
// 	}

// 	cardID, err := strconv.Atoi(vars["cardID"])
// 	if err != nil {
// 		http.Error(w, "Invalid card ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Read the request body
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
// 		return
// 	}

// 	// Parse the JSON request body
// 	var requestData struct {
// 		Name *string `json:"name"`
// 		Description *string `json:"description"`
// 	}

// 	err = json.Unmarshal(body, &requestData)
// 	if err != nil {
// 		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
// 		return
// 	}


// 	var query string
// 	var args []interface{}

// 	if requestData.Name != nil  && requestData.Description != nil {
// 		query = "UPDATE cards SET name = $1, description = $2 WHERE id = $3"
// 		args = []interface{}{*requestData.Name, *requestData.Description, cardID}
// 	} else if requestData.Name != nil {
// 		query = "UPDATE cards SET name = $1 WHERE id = $2"
// 		args = []interface{}{*requestData.Name, cardID}
// 	} else if requestData.Description != nil {
// 		query = "UPDATE cards SET description = $1 WHERE id = $2"
// 		args = []interface{}{*requestData.Description, cardID}
// 	} else {
// 		http.Error(w, "No update data provided", http.StatusBadRequest)
// 		return
// 	}


// 	// Update the list's name in the database
// 	_, err = db.Exec(query, args...)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to update card, %s", err), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// }
