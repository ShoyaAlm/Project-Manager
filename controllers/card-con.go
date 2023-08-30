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

	// Fetch card details
	cardRow := db.QueryRow("SELECT id, name, description FROM cards c WHERE c.list_id = $1", listID)
	card := &model.Card{}
	err = cardRow.Scan(&card.ID, &card.Name, &card.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Card not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to fetch card data, %s", err), http.StatusInternalServerError)
		}
		return
	}

	rows, err := db.Query(`
        SELECT c.id AS card_id, c.dates AS card_dates,
		m.id AS member_id, m.name AS member_name,
            cl.id AS checklist_id, cl.name AS checklist_name,
            i.id AS item_id, i.name AS item_name, i.due_date AS item_due_date, i.assigned_to AS item_assigned_to
        FROM cards c
        LEFT JOIN members m ON c.id = m.card_id
        LEFT JOIN checklists cl ON c.id = cl.card_id
        LEFT JOIN items i ON cl.id = i.checklist_id`)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch card, %s", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {

		var (
			cardID                                           int
			itemAssignedTo, cardDates                        pq.StringArray
			memberID, checklistID, itemID                    sql.NullInt64
			memberName, checklistName, itemName, itemDueDate sql.NullString
		)
		err = rows.Scan(&cardID, &cardDates, &memberID, &memberName, &checklistID, &checklistName, &itemID, &itemName, &itemDueDate, &itemAssignedTo)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, fmt.Sprintf("Card not found, %s", err), http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Failed to fetch card data, %s", err), http.StatusInternalServerError)
			}
			return
		}

		if memberName.Valid && memberID.Valid {
			card.Members = append(card.Members, &model.Member{ID: int(memberID.Int64), Name: memberName.String})
		}
		if checklistID.Valid && checklistName.Valid {
			checklist, ok := findChecklist(card.Checklists, int(checklistID.Int64))
			if !ok {
				checklist = &model.Checklist{
					ID:    int(checklistID.Int64),
					Name:  checklistName.String,
					Items: []*model.Item{},
				}

				if itemID.Valid {
					item := &model.Item{
						ID:         int(itemID.Int64),
						Name:       itemName.String,
						DueDate:    itemDueDate.String,
						AssignedTo: itemAssignedTo,
					}
					checklist.Items = append(checklist.Items, item)
				}

				card.Checklists = append(card.Checklists, checklist)
			}
		}
		card.Dates = cardDates
	}

	jsonData, err := json.Marshal(card)
	if err != nil {
		http.Error(w, "failed to marshal card data", http.StatusInternalServerError)
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

	// Fetch card details
	cardRow := db.QueryRow("SELECT id, name, description FROM cards c WHERE c.id = $1 AND c.list_id = $2", cardID, listID)
	card := &model.Card{}
	err = cardRow.Scan(&card.ID, &card.Name, &card.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Card not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to fetch card data, %s", err), http.StatusInternalServerError)
		}
		return
	}

	fmt.Printf("list ID %v, card ID %v ", listID, cardID)

	rows, err := db.Query(`
        SELECT c.dates AS card_dates,
		m.id AS member_id, m.name AS member_name,
            cl.id AS checklist_id, cl.name AS checklist_name,
            i.id AS item_id, i.name AS item_name, i.due_date AS item_due_date, i.assigned_to AS item_assigned_to
        FROM cards c
        LEFT JOIN members m ON $1 = m.card_id
        LEFT JOIN checklists cl ON $1 = cl.card_id
        LEFT JOIN items i ON cl.id = i.checklist_id`, cardID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch card, %s", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {

		var (
			itemAssignedTo, cardDates                        pq.StringArray
			memberID, checklistID, itemID                    int
			memberName, checklistName, itemName, itemDueDate string
		)
		err = rows.Scan(&cardDates, &memberID, &memberName, &checklistID, &checklistName, &itemID, &itemName, &itemDueDate, &itemAssignedTo)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, fmt.Sprintf("Card not found, %s", err), http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Failed to fetch card data, %s", err), http.StatusInternalServerError)
			}
			return
		}

		if memberName != "" {
			card.Members = append(card.Members, &model.Member{ID: memberID, Name: memberName})
		}
		if checklistID != 0 {
			checklist, ok := findChecklist(card.Checklists, checklistID)
			if !ok {
				checklist = &model.Checklist{
					ID:    int(checklistID),
					Name:  checklistName,
					Items: []*model.Item{},
				}

				if itemID != 0 {
					item := &model.Item{
						ID:         itemID,
						Name:       itemName,
						DueDate:    itemDueDate,
						AssignedTo: itemAssignedTo,
					}
					checklist.Items = append(checklist.Items, item)
				}

				card.Checklists = append(card.Checklists, checklist)
			}
		}
		card.Dates = cardDates
	}

	jsonData, err := json.Marshal(card)
	if err != nil {
		http.Error(w, "failed to marshal card data", http.StatusInternalServerError)
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
