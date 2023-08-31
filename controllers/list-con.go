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

		cardRows, err := db.Query(`SELECT c.id AS card_id, c.name AS card_name, c.description AS cardDescription, c.dates AS card_dates
									FROM cards c
							   WHERE c.list_id = $1`, list.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch cards for list, %s", err), http.StatusInternalServerError)
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
				http.Error(w, fmt.Sprintf("Error scanning cardRows, %s", err), http.StatusInternalServerError)
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

	cardRows, err := db.Query(`SELECT c.id AS card_id, c.name AS card_name, c.description AS cardDescription, c.dates AS card_dates
		FROM cards c
		WHERE c.list_id = $1`, list.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch cards for list, %s", err), http.StatusInternalServerError)
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
			http.Error(w, fmt.Sprintf("Error scanning cardRows, %s", err), http.StatusInternalServerError)
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
		Name  string        `json:"name"`
		Cards []*model.Card `json:"cards"`
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

	emptyItem := &model.Item{
		ID:         0,
		Name:       "default item",
		DueDate:    "september 20th",
		AssignedTo: []string{"person1", "person2"},
	}

	emptyChecklist := &model.Checklist{
		ID:    0,
		Name:  "default checklist",
		Items: []*model.Item{emptyItem},
	}

	emptyMember := &model.Member{
		ID:   0,
		Name: "member 1",
	}

	emptyCard := &model.Card{
		ID:          0,
		Name:        "default card",
		Description: "default",
		Dates:       []string{"october 1st", "november 1st"},
		Checklists:  []*model.Checklist{emptyChecklist},
		Members:     []*model.Member{emptyMember},
	}

	responseData := &model.List{
		ID:    newListID,
		Name:  requestData.Name,
		Cards: []*model.Card{emptyCard}, // Initialize an empty cards attribute
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

// vars := mux.Vars(r)
// listID, err := strconv.Atoi(vars["id"])
// if err != nil {
// 	http.Error(w, "Invalid list ID", http.StatusBadRequest)
// 	return
// }

// // Fetch list details
// listRow := db.QueryRow("SELECT id, name FROM lists WHERE id = $1", listID)
// list := &model.List{}
// err = listRow.Scan(&list.ID, &list.Name)
// if err != nil {
// 	if err == sql.ErrNoRows {
// 		http.Error(w, "List not found", http.StatusNotFound)
// 	} else {
// 		http.Error(w, fmt.Sprintf("Failed to fetch list data, %s", err), http.StatusInternalServerError)
// 	}
// 	return
// }

// // Fetch related cards, their members, checklists, and items
// rows, err := db.Query(`
//     SELECT
//         c.id AS card_id, c.name AS card_name, c.description AS card_description, c.dates as card_dates,
//         m.id AS member_id, m.name AS member_name,
//         cl.id AS checklist_id, cl.name AS checklist_name,
//         i.id AS item_id, i.name AS item_name, i.due_date AS item_due_date, i.assigned_to AS item_assigned_to
//     FROM cards c
//     LEFT JOIN members m ON c.id = m.card_id
//     LEFT JOIN checklists cl ON c.id = cl.card_id
//     LEFT JOIN items i ON cl.id = i.checklist_id
//     WHERE c.list_id = $1`, listID)
// if err != nil {
// 	http.Error(w, fmt.Sprintf("Failed to fetch cards for list, %s", err), http.StatusInternalServerError)
// 	return
// }
// defer rows.Close()

// cardMap := make(map[int]*model.Card)

// for rows.Next() {
// 	var (
// 		cardID, checklistID, memberID, itemID int
// 		cardName, cardDescription, memberName string
// 		checklistName, itemName, itemDueDate  string
// 		itemAssignedTo, cardDates             pq.StringArray
// 	)
// 	err := rows.Scan(&cardID, &cardName, &cardDescription, &cardDates, &memberID, &memberName, &checklistID, &checklistName, &itemID, &itemName, &itemDueDate, &itemAssignedTo)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Error scanning rows, %s", err), http.StatusInternalServerError)
// 		return
// 	}

// 	card, ok := cardMap[cardID]
// 	if !ok {
// 		card = &model.Card{
// 			ID:          cardID,
// 			Name:        cardName,
// 			Description: cardDescription,
// 			Dates:       cardDates,
// 			Members:     []*model.Member{},
// 			Checklists:  []*model.Checklist{},
// 		}
// 		cardMap[cardID] = card
// 	}

// 	if memberName != "" {
// 		card.Members = append(card.Members, &model.Member{ID: memberID, Name: memberName})
// 	}

// 	if checklistID != 0 {
// 		checklist, ok := findChecklist(card.Checklists, checklistID)
// 		if !ok {
// 			checklist = &model.Checklist{
// 				ID:    checklistID,
// 				Name:  checklistName,
// 				Items: []*model.Item{},
// 			}
// 			card.Checklists = append(card.Checklists, checklist)
// 		}

// 		if itemID != 0 {
// 			item := &model.Item{
// 				ID:         itemID,
// 				Name:       itemName,
// 				DueDate:    itemDueDate,
// 				AssignedTo: itemAssignedTo,
// 			}
// 			checklist.Items = append(checklist.Items, item)
// 		}
// 	}
// }

// // Assign related cards to the list
// var cards []*model.Card
// for _, card := range cardMap {
// 	cards = append(cards, card)
// }
// list.Cards = cards

// jsonData, err := json.Marshal(list)
// if err != nil {
// 	http.Error(w, "Failed to marshal list data", http.StatusInternalServerError)
// 	return
// }

// w.Header().Set("Content-Type", "application/json")
// w.Write(jsonData)
