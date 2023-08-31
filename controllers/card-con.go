package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

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

		card := &model.Card{
			ID:          cardID,
			Name:        cardName,
			Description: cardDescription,
			Dates:       cardDates,
			Members:     []*model.Member{},
			Checklists:  []*model.Checklist{},
		}

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

			checklists = append(checklists, checklist)

		}

		card.Checklists = checklists

		// Start looking for members inside every card
		memberRows, err := db.Query(`SELECT m.id AS member_id, m.name AS member_name
							FROM members m
						WHERE m.card_id = $1`, cardID)

		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch members for card, %s", err), http.StatusInternalServerError)
			return

		}

		defer memberRows.Close()

		var members []*model.Member

		for memberRows.Next() {

			var (
				memberID   int
				memberName string
			)

			err := memberRows.Scan(&memberID, &memberName)

			if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning memberRows, %s", err), http.StatusInternalServerError)
				return
			}

			member := &model.Member{
				ID:   memberID,
				Name: memberName,
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

	card := &model.Card{
		ID:          cardID,
		Name:        cardName,
		Description: cardDescription,
		Dates:       cardDates,
		Members:     []*model.Member{},
		Checklists:  []*model.Checklist{},
	}

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

		checklists = append(checklists, checklist)

	}

	card.Checklists = checklists

	// Start looking for members inside every card
	memberRows, err := db.Query(`SELECT m.id AS member_id, m.name AS member_name
			FROM members m
			WHERE m.card_id = $1`, cardID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch members for card, %s", err), http.StatusInternalServerError)
		return

	}

	defer memberRows.Close()

	var members []*model.Member

	for memberRows.Next() {

		var (
			memberID   int
			memberName string
		)

		err := memberRows.Scan(&memberID, &memberName)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning memberRows, %s", err), http.StatusInternalServerError)
			return
		}

		member := &model.Member{
			ID:   memberID,
			Name: memberName,
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

	var requestData struct {
		Name        string             `json:"name"`
		Description string             `json:"description"`
		Dates       []string           `json:"dates"`
		Checklists  []*model.Checklist `json:"checklists"`
		Members     []*model.Member    `json:"members"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON data, %s", err), http.StatusBadRequest)
		return
	}

	var newCardID, newChecklistID, newItemID, newMemberID int

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create card, %s", err), http.StatusInternalServerError)
		return
	}

	emptyItem := &model.Item{
		ID:         newItemID,
		Name:       "default item",
		DueDate:    "september 20th",
		AssignedTo: []string{"person1", "person2"},
	}

	emptyChecklist := &model.Checklist{
		ID:    newChecklistID,
		Name:  "default checklist",
		Items: []*model.Item{emptyItem},
	}

	emptyMember := &model.Member{
		ID:   newMemberID,
		Name: "member 1",
	}

	// Create a new card with non-null fields
	newCard := &model.Card{
		ID:          newCardID,
		Name:        requestData.Name,
		Description: "default description",
		Dates:       []string{"october 1st", "november 1st"}, // Initialize as empty slice
		Checklists:  []*model.Checklist{emptyChecklist},      // Initialize as empty slice
		Members:     []*model.Member{emptyMember},
	}

	err = db.QueryRow("INSERT INTO cards (name, description, dates, list_id) VALUES ($1, $2, $3, $4) RETURNING id",
		newCard.Name, newCard.Description, pq.Array(newCard.Dates), listID).Scan(&newCardID)

	err = db.QueryRow("INSERT INTO checklists (name, card_id) VALUES ($1, $2) RETURNING id",
		emptyChecklist.Name, newCardID).Scan(&newChecklistID)

	err = db.QueryRow("INSERT INTO items (name, duedate, assignedto, checklist_id) VALUES ($1, $2, $3, $4) RETURNING id",
		emptyItem.Name, emptyItem.DueDate, pq.Array(emptyItem.AssignedTo), newChecklistID).Scan(&newItemID)

	err = db.QueryRow("INSERT INTO members (name, card_id) VALUES ($1, $2) RETURNING id",
		emptyMember.Name, newCardID).Scan(&newMemberID)

	// Fetch the associated list
	listRow := db.QueryRow("SELECT id, name FROM lists WHERE id = $1", listID)
	list := &model.List{}
	err = listRow.Scan(&list.ID, &list.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch list data, %s", err), http.StatusInternalServerError)
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
		Name string `json:"name"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	// Update the list's name in the database
	_, err = db.Exec("UPDATE cards SET name = $1 WHERE id = $2", requestData.Name, cardID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update list, %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
