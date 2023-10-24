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

	listRows, err := db.Query("SELECT id, name FROM lists")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch lists, %s", err), http.StatusInternalServerError)
		return
	}

	defer listRows.Close()

	var lists []*model.List

	for listRows.Next() {
		var (
			listID   int
			listName string
		)

		err := listRows.Scan(&listID, &listName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning rows, %s", err), http.StatusInternalServerError)
			return
		}

		list := &model.List{
			ID:   listID,
			Name: listName,
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

		lists = append(lists, list)
	}

	jsonData, err := json.Marshal(lists)
	if err != nil {
		http.Error(w, "Failed to marshal list data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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
		Name  		string	           `json:"name"`
		UserID 		int			   	   `json:"user_id"`
		Username 	string			   `json:"username"`
		UserEmail 	string			   `json:"user_email"`
		Cards []*model.Card 		   `json:"cards"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON data, %s", err), http.StatusBadRequest)
		return
	}

	var newListID, newCardID, newChecklistID, newItemID, newMemberID int
	// err = db.QueryRow("INSERT INTO lists (name) VALUES ($1) RETURNING id", requestData.Name).Scan(&newListID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create list, %s", err), http.StatusInternalServerError)
		return
	}

	emptyItem := &model.Item{
		ID:         newItemID,
		Name:       "آیتم 1",
		DueDate:    "2023-09-20T00:00:00Z",
		AssignedTo: []string{"شخص 1", "شخص 2"},
	}

	emptyChecklist := &model.Checklist{
		ID:    newChecklistID,
		Name:  "چکلیست جدید",
		Items: []*model.Item{emptyItem},
	}

	emptyMember := &model.Member{
		ID:   newMemberID,
		Name: requestData.Username,
		Email: requestData.UserEmail,
	}

	emptyCard := &model.Card{
		ID:          newCardID,
		Name:        "کارت جدید",
		Description: "توضیحات",
		Dates:       []string{"1 شهریور", "1 مهر"},
		Checklists:  []*model.Checklist{emptyChecklist},
		Members:     []*model.Member{emptyMember},
		Owner:       &model.User{ID: requestData.UserID},
	}

	newList := &model.List{
		ID:    newListID,
		Name:  requestData.Name,
		Cards: []*model.Card{emptyCard}, // Initialize an empty cards attribute
	}

	err = db.QueryRow("INSERT INTO lists (name) VALUES ($1) RETURNING id", newList.Name).Scan(&newListID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert card, %s", err), http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("INSERT INTO cards (name, description, dates, list_id) VALUES ($1, $2, $3, $4) RETURNING id",
		emptyCard.Name, emptyCard.Description, pq.Array(emptyCard.Dates), newListID).Scan(&newCardID)
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

	err = db.QueryRow("INSERT INTO items (name, due_date, assigned_to, checklist_id) VALUES ($1, $2, $3, $4) RETURNING id",
		emptyItem.Name, emptyItem.DueDate, pq.Array(emptyItem.AssignedTo), newChecklistID).Scan(&newItemID)
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

	// if err := tx.Commit(); err != nil {
	// 	http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
	// 	return
	// }

	// Fetch the associated list
	listRow := db.QueryRow("SELECT id, name FROM lists WHERE id = $1", newListID)
	list := []*model.List{}
	err = listRow.Scan(&newList.ID, &newList.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch list data, %s", err), http.StatusInternalServerError)
		return
	}

	// Append the new card to the list's cards slice
	list = append(list, newList)

	jsonData, err := json.Marshal(list)
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

	var cardIDs []int
	rows, err := db.Query("SELECT id FROM cards WHERE list_id = $1", listID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch card IDs, %s", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cardID int
		if err := rows.Scan(&cardID); err != nil {
			http.Error(w, fmt.Sprintf("Failed to scan card ID, %s", err), http.StatusInternalServerError)
			return
		}
		cardIDs = append(cardIDs, cardID)
	}

	// loop for deleting the list's data
	for _, cardID := range cardIDs {

		_, err := db.Exec("DELETE FROM members WHERE card_id = $1", cardID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete members of list, %s", err), http.StatusInternalServerError)
			return
		}

		// deleting items & checklists
		var checklistIDs []int
		checklistRows, err := db.Query("SELECT id FROM checklists WHERE card_id = $1", cardID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch checklist IDs, %s", err), http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		for checklistRows.Next() {
			var checklistID int
			if err := checklistRows.Scan(&checklistID); err != nil {
				http.Error(w, fmt.Sprintf("Failed to scan checklist ID, %s", err), http.StatusInternalServerError)
				return
			}
			checklistIDs = append(checklistIDs, checklistID)
		}

		for _, checklistID := range checklistIDs {
			_, err := db.Exec("DELETE FROM items WHERE checklist_id = $1", checklistID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to delete items of checklist, %s", err), http.StatusInternalServerError)
				return
			}
		}

		//delete checklists
		_, err = db.Exec("DELETE FROM checklists WHERE card_id = $1", cardID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete checklists of a list, %s", err), http.StatusInternalServerError)
			return
		}

	}

	_, err = db.Exec("DELETE FROM cards WHERE list_id = $1", listID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete card, %s", err), http.StatusInternalServerError)
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
