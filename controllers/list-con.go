package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"project-manager/model"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "host=localhost user=postgres password=2220819 dbname=postgres sslmode=disable")
	fmt.Printf("Connecting to the server...")
	if err != nil {
		fmt.Printf("error : %s", err)
		log.Fatal(err)
	}
}

func GetAllLists(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT id, name FROM lists")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch lists, %s", err), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var lists []*model.List

	for rows.Next() {
		var (
			listID   int
			listName string
		)

		err := rows.Scan(&listID, &listName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning rows, %s", err), http.StatusInternalServerError)
			return
		}

		list := &model.List{
			ID:   listID,
			Name: listName,
		}

		lists = append(lists, list)
	}

	for _, list := range lists {
		rows, err := db.Query(`SELECT c.id AS card_id, c.name AS card_name, c.description AS cardDescription, c.dates AS card_dates,
								m.id AS member_id, m.name AS member_name,
								cl.id AS checklist_id, cl.name AS checklist_name,
								i.id AS item_id, i.name AS item_name, i.due_date AS item_due_date, i.assigned_to AS item_assigned_to
								FROM cards c
								LEFT JOIN members m ON m.card_id = c.id
								LEFT JOIN checklists cl ON cl.card_id = c.id
								LEFT JOIN items i ON cl.id = i.checklist_id
								WHERE list_id = $1`, list.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch cards for list, %s", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var cards []*model.Card

		for rows.Next() {
			var (
				cardID                                                                           int
				cardName, cardDescription                                                        string
				cardDates, itemAssignedTo                                                        pq.StringArray
				memberIDNullable, checklistIDNullable, itemIDNullable                            sql.NullInt64
				memberNameNullable, checklistNameNullable, itemNameNullable, itemDueDateNullable sql.NullString
			)

			err := rows.Scan(&cardID, &cardName, &cardDescription, &cardDates,
				&memberIDNullable, &memberNameNullable,
				&checklistIDNullable, &checklistNameNullable,
				&itemIDNullable, &itemNameNullable, &itemDueDateNullable, &itemAssignedTo)
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

			if memberIDNullable.Valid && memberNameNullable.Valid {
				card.Members = append(card.Members, &model.Member{ID: int(memberIDNullable.Int64), Name: memberNameNullable.String})
			}

			if checklistIDNullable.Valid && checklistNameNullable.Valid {
				checklist, ok := findChecklist(card.Checklists, int(checklistIDNullable.Int64))
				if !ok {
					checklist = &model.Checklist{
						ID:    int(checklistIDNullable.Int64),
						Name:  checklistNameNullable.String,
						Items: []*model.Item{},
					}

					if itemIDNullable.Valid {
						item := &model.Item{
							ID:         int(itemIDNullable.Int64),
							Name:       itemNameNullable.String,
							DueDate:    itemDueDateNullable.String,
							AssignedTo: itemAssignedTo,
						}
						checklist.Items = append(checklist.Items, item)
					}

					card.Checklists = append(card.Checklists, checklist)
				}
			}

			cards = append(cards, card)
		}

		list.Cards = cards
	}

	jsonData, err := json.Marshal(lists)
	if err != nil {
		http.Error(w, "Failed to marshal list data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

// #############################################
// #############################################
// #############################################
// #############################################
// #############################################
// #############################################

func GetAList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	// Fetch list details
	listRow := db.QueryRow("SELECT id, name FROM lists WHERE id = $1", listID)
	list := &model.List{}
	err = listRow.Scan(&list.ID, &list.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "List not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to fetch list data, %s", err), http.StatusInternalServerError)
		}
		return
	}

	// Fetch related cards, their members, checklists, and items
	rows, err := db.Query(`
        SELECT
            c.id AS card_id, c.name AS card_name, c.description AS card_description, c.dates as card_dates,
           	m.id AS member_id, m.name AS member_name,
            cl.id AS checklist_id, cl.name AS checklist_name,
            i.id AS item_id, i.name AS item_name, i.due_date AS item_due_date, i.assigned_to AS item_assigned_to
        FROM cards c
        LEFT JOIN members m ON c.id = m.card_id
        LEFT JOIN checklists cl ON c.id = cl.card_id
        LEFT JOIN items i ON cl.id = i.checklist_id
        WHERE c.list_id = $1`, listID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch cards for list, %s", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cards []*model.Card
	cardMap := make(map[int]*model.Card)

	for rows.Next() {
		var (
			cardID, checklistID, memberID, itemID int
			cardName, cardDescription, memberName string
			checklistName, itemName, itemDueDate  string
			itemAssignedTo, cardDates             pq.StringArray
		)
		err := rows.Scan(&cardID, &cardName, &cardDescription, &cardDates, &memberID, &memberName, &checklistID, &checklistName, &itemID, &itemName, &itemDueDate, &itemAssignedTo)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning rows, %s", err), http.StatusInternalServerError)
			return
		}

		card, ok := cardMap[cardID]
		if !ok {
			card = &model.Card{
				ID:          cardID,
				Name:        cardName,
				Description: cardDescription,
				Dates:       cardDates,
				Members:     []*model.Member{},
				Checklists:  []*model.Checklist{},
			}
			cardMap[cardID] = card
			cards = append(cards, card)
		}

		if memberName != "" {
			card.Members = append(card.Members, &model.Member{ID: memberID, Name: memberName})
		}

		if checklistID != 0 {
			checklist, ok := findChecklist(card.Checklists, checklistID)
			if !ok {
				checklist = &model.Checklist{
					ID:    checklistID,
					Name:  checklistName,
					Items: []*model.Item{},
				}
				card.Checklists = append(card.Checklists, checklist)
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
		}
	}

	// Assign related cards to the list
	list.Cards = cards

	jsonData, err := json.Marshal(list)
	if err != nil {
		http.Error(w, "Failed to marshal list data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func findChecklist(checklists []*model.Checklist, id int) (*model.Checklist, bool) {
	for _, c := range checklists {
		if c.ID == id {
			return c, true
		}
	}
	return nil, false
}

func CreateList(w http.ResponseWriter, r *http.Request) {

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

	var newListID int
	err = db.QueryRow("INSERT INTO lists (name) VALUES ($1) RETURNING id", requestData.Name).Scan(&newListID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create list, %s", err), http.StatusInternalServerError)
		return
	}

	responseData := &model.List{
		ID:    newListID,
		Name:  requestData.Name,
		Cards: []*model.Card{}, // Initialize an empty cards attribute
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Failed to marshal response data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)

}

func DeleteAList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	// Delete the list and related data
	_, err = db.Exec("DELETE FROM lists WHERE id = $1", listID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete list, %s", err), http.StatusInternalServerError)
		return
	}

	// You may also want to delete related cards, members, checklists, and items
	// Here, I'm assuming you have foreign key constraints that automatically handle this
	// If not, you should handle the deletion of related data accordingly.

	w.WriteHeader(http.StatusAccepted)
}

func UpdateAList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
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
	_, err = db.Exec("UPDATE lists SET name = $1 WHERE id = $2", requestData.Name, listID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update list, %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
