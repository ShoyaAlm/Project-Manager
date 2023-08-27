package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"project-manager/model"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
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
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, name FROM lists")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch lists from database, %s", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lists []*model.List
	for rows.Next() {
		list := &model.List{}
		err := rows.Scan(&list.ID, &list.Name)
		if err != nil {
			http.Error(w, "Error scanning rows", http.StatusInternalServerError)
			return
		}
		lists = append(lists, list)
	}

	jsonData, err := json.Marshal(lists)
	if err != nil {
		http.Error(w, "Failed to marshal lists data", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

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
            c.id AS card_id, c.name AS card_name, c.description AS card_description,
            m.name AS member_name,
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
			cardID, checklistID, itemID           int
			cardName, cardDescription, memberName string
			checklistName, itemName               string
			itemDueDate                           sql.NullTime
			itemAssignedTo                        pq.StringArray
		)
		err := rows.Scan(&cardID, &cardName, &cardDescription, &memberName, &checklistID, &checklistName, &itemID, &itemName, &itemDueDate, pq.Array(&itemAssignedTo))
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
				Members:     []*model.Member{},
				Checklists:  []*model.Checklist{},
			}
			cardMap[cardID] = card
			cards = append(cards, card)
		}

		if memberName != "" {
			card.Members = append(card.Members, &model.Member{Name: memberName})
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
					DueDate:    itemDueDate.Time,
					AssignedTo: []string(itemAssignedTo),
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

// func GetAList(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	listID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid list ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Fetch list details
// 	listRow := db.QueryRow("SELECT id, name FROM lists WHERE id = $1", listID)
// 	list := &model.List{}
// 	err = listRow.Scan(&list.ID, &list.Name)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			http.Error(w, "List not found", http.StatusNotFound)
// 		} else {
// 			http.Error(w, fmt.Sprintf("Failed to fetch list data, %s", err), http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	// Fetch related cards and their members
// 	rows, err := db.Query(`
//         SELECT
//             c.id, c.name, c.description,
//             m.name AS member_name, c.dates
//         FROM cards c
//         LEFT JOIN members m ON c.id = m.card_id
//         WHERE c.list_id = $1`, listID)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to fetch cards for list, %s", err), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var cards []*model.Card
// 	cardMap := make(map[int]*model.Card)

// 	for rows.Next() {
// 		var cardID int
// 		var cardName, cardDescription, memberName string
// 		var cardDate pq.StringArray
// 		err := rows.Scan(&cardID, &cardName, &cardDescription, &memberName, &cardDate)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Error scanning card rows, %s", err), http.StatusInternalServerError)
// 			return
// 		}

// 		card, ok := cardMap[cardID]
// 		if !ok {
// 			card = &model.Card{
// 				ID:          cardID,
// 				Name:        cardName,
// 				Description: cardDescription,
// 				Members:     []*model.Member{},
// 				Dates:       cardDate,
// 			}
// 			cardMap[cardID] = card
// 			cards = append(cards, card)
// 		}

// 		if memberName != "" {
// 			card.Members = append(card.Members, &model.Member{Name: memberName})
// 		}
// 	}

// 	// Assign related cards to the list
// 	list.Cards = cards

// 	jsonData, err := json.Marshal(list)
// 	if err != nil {
// 		http.Error(w, "Failed to marshal list data", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(jsonData)

// }

func CreateList(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var newList model.List
	err := json.NewDecoder(r.Body).Decode(&newList)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO lists (name) VALUES ($1) RETURNING id", newList.Name)
	if err != nil {
		http.Error(w, "Failed to insert new list", http.StatusInternalServerError)
		return
	}

	var newID int
	err = db.QueryRow("SELECT LASTVAL()").Scan(&newID)
	if err != nil {
		http.Error(w, "Failed to get new list ID", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Result : %s", result)
	response := map[string]interface{}{
		"message": "List created successfully",
		"id":      newID,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response data", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)

}

// func UpdateList(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	listID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid list ID", http.StatusBadRequest)
// 		return
// 	}

// 	var updatedList model.List
// 	err = json.NewDecoder(r.Body).Decode(&updatedList)
// 	if err != nil {
// 		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
// 	}

// 	found := false
// 	for i, list := range lists {
// 		if list.ID == listID {
// 			updatedList.ID = list.ID
// 			updatedList.Cards = list.Cards
// 			lists[i] = &updatedList
// 			found = true
// 			break
// 		}
// 	}

// 	if !found {
// 		http.Error(w, "List not found", http.StatusNotFound)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "List updated successfully")
// }

// func DeleteList(w http.ResponseWriter, r *http.Request) {
// 	// Parse the list ID from the request URL or request body
// 	vars := mux.Vars(r)
// 	listID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid list ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Find and remove the list with the given ID from your 'lists' slice
// 	found := false
// 	for i, list := range lists {
// 		if list.ID == listID {
// 			lists = append(lists[:i], lists[i+1:]...)
// 			found = true
// 			break
// 		}
// 	}

// 	if !found {
// 		http.Error(w, "List not found", http.StatusNotFound)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "List deleted successfully")
// }
