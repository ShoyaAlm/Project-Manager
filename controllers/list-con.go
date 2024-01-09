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
	"time"

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


	vars := mux.Vars(r)
	boardID, err := strconv.Atoi(vars["board_id"])
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
		return
	}

	
	// listRows, err := db.Query("SELECT id, name FROM lists ORDER BY position WHERE board_id = 1")

	listRows, err := db.Query("SELECT id, name FROM lists WHERE board_id = $1 ORDER BY position", boardID)

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

		cardRows, err := db.Query(`SELECT c.id AS card_id, c.name AS card_name, c.description AS card_description,
								   c.dates AS card_dates, c.position AS card_position, c.label AS card_label
									FROM cards c
							   WHERE c.list_id = $1
							   ORDER BY c.position;`,
							    list.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch cards for list, %s", err), http.StatusInternalServerError)
			return
		}
		defer cardRows.Close()

		var cards []*model.Card

		for cardRows.Next() {

			var (
				cardID  		    	  int
				cardName				  string
				cardDescription 		  sql.NullString
				cardDates                 pq.StringArray
				cardPosition  			  sql.NullInt64
				cardLabel				  sql.NullString
				// ownerID					  int
			)
			err := cardRows.Scan(&cardID, &cardName, &cardDescription, &cardDates, &cardPosition, &cardLabel)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning cardRows, %s", err), http.StatusInternalServerError)
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


			var position int

			// Check if cardPosition is valid (not NULL)
			if cardPosition.Valid {
				position = int(cardPosition.Int64)
			} else {
				position = 0 // or any other default value
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
				Description: getStringOrNil(cardDescription),
				Dates:       dates,
				Members:     []*model.User{},
				Checklists:  []*model.Checklist{},
				Position: 	 position,
				Label:		 label,	
				// OwnerID: 	 ownerID,
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
						itemName 					string
						itemStartDate, itemDueDate 	time.Time
						itemDone 					bool
						itemAssignedTo		        []*model.Member
					)

					err := itemRows.Scan(&itemID, &itemName, &itemStartDate, &itemDueDate, &itemDone)

					if err != nil {
						http.Error(w, fmt.Sprintf("Error scanning itemRows, %s", err), http.StatusInternalServerError)
						return
					}

					item := &model.Item{
						ID:         itemID,
						Name:       itemName,
						StartDate:  itemStartDate,
						DueDate:    itemDueDate,
						Done: 		itemDone,	
						AssignedTo: itemAssignedTo,
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
	boardID, err := strconv.Atoi(vars["board_id"])
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
		return
	}

	// vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	// Fetch list details
	listRow := db.QueryRow("SELECT id, name FROM lists WHERE id = $1 and board_id = $2", listID, boardID)
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

	cardRows, err := db.Query(`SELECT c.id AS card_id, c.name AS card_name, c.description AS cardDescription, c.dates AS card_dates, c.owner_id AS owner_id
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
			ownerID 				  int
		)
		err := cardRows.Scan(&cardID, &cardName, &cardDescription, &cardDates, &ownerID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning cardRows, %s", err), http.StatusInternalServerError)
			return
		}

		var dates []time.Time

		for _, dateString := range cardDates {
			date, err := time.Parse("2006-01-02", dateString)
			if err != nil {
				// Handle the error
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
			OwnerID: 	 ownerID,
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
					itemName 					string
					itemStartDate, itemDueDate 	time.Time
					itemDone 					bool
					itemAssignedTo		        []*model.Member
				)

				err := itemRows.Scan(&itemID, &itemName, &itemStartDate, &itemDueDate, &itemDone)

				if err != nil {
					http.Error(w, fmt.Sprintf("Error scanning itemRows, %s", err), http.StatusInternalServerError)
					return
				}

				item := &model.Item{
					ID:         itemID,
					Name:       itemName,
					StartDate:  itemStartDate,
					DueDate:    itemDueDate,
					Done: 		itemDone,	
					AssignedTo: itemAssignedTo,
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

func getStringOrNil(nullString sql.NullString) string {
    if nullString.Valid {
        return nullString.String
    }
    return ""
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

	vars := mux.Vars(r)
	boardID, err := strconv.Atoi(vars["board_id"])
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
		return
	}



	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body, %s", err), http.StatusInternalServerError)
		return
	}

	var requestData struct {
		Name     string `json:"name"`
		UserID   int    `json:"user_id"`
		Username string `json:"username"`
		UserEmail string `json:"user_email"`
		OwnerID  int    `json:"owner_id"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON data, %s", err), http.StatusBadRequest)
		return
	}

	var newListID int

	var maxPosition int
	err = db.QueryRow("SELECT COALESCE(MAX(position), 0) FROM lists").Scan(&maxPosition)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch maximum position, %s", err), http.StatusInternalServerError)
		return
	}

	newPosition := maxPosition + 1

	err = db.QueryRow("INSERT INTO lists (name, board_id, position) VALUES ($1, $2, $3) RETURNING id", requestData.Name, boardID, newPosition).Scan(&newListID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create list, %s", err), http.StatusInternalServerError)
		return
	}

	newList := &model.List{ID: newListID, Name: requestData.Name}

	jsonData, err := json.Marshal(newList)
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

	for _, cardID := range cardIDs {


		_, err = db.Exec("DELETE FROM activities WHERE card_id = $1", cardID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete activities of card, %s", err), http.StatusInternalServerError)
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

			// Fetch the list of item IDs associated with the checklist
			itemRows, err := db.Query("SELECT id FROM items WHERE checklist_id = $1", checklistID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to fetch items of checklist, %s", err), http.StatusInternalServerError)
				return
			}
	
			defer itemRows.Close()
	
			for itemRows.Next() {
				var itemID int
				if err := itemRows.Scan(&itemID); err != nil {
					http.Error(w, fmt.Sprintf("Error scanning item ID, %s", err), http.StatusInternalServerError)
					return
				}
	
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
			http.Error(w, fmt.Sprintf("Failed to delete members of list, %s", err), http.StatusInternalServerError)
			return
		}
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




func UpdateListOrder(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusInternalServerError)
        return
    }

    // Parse the JSON request body
    var requestData struct {
		ListOrder  []int    `json:"listOrder"`
	}
	
    err = json.Unmarshal(body, &requestData)
    if err != nil {
        http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
        return
    }

    if len(requestData.ListOrder) > 0 {

        tx, err := db.Begin()
        if err != nil {
            http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
            return
        }
        defer tx.Rollback()

        for i, cardID := range requestData.ListOrder {
            _, err := tx.Exec("UPDATE lists SET position = $1 WHERE id = $2", i, cardID)
            if err != nil {
                http.Error(w, "Failed to update card order", http.StatusInternalServerError)
                fmt.Printf("error : %v", err)
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



func UpdateCardOrder(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusInternalServerError)
        return
    }

    // Parse the JSON request body
    var requestData struct {
		ListID 	   int	`json:"listId"`
		CardOrder  []int    `json:"cardOrder"`
	}
	
    err = json.Unmarshal(body, &requestData)
    if err != nil {
        http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		fmt.Printf("error : %s", err)
        return
    }

    if len(requestData.CardOrder) > 0 {

        tx, err := db.Begin()
        if err != nil {
            http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
            return
        }
        defer tx.Rollback()

        for i, cardID := range requestData.CardOrder {
            _, err := tx.Exec("UPDATE cards SET position = $1 WHERE id = $2", i, cardID)
            if err != nil {
                http.Error(w, "Failed to update card order", http.StatusInternalServerError)
                fmt.Printf("error : %v\n", err)
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


func MoveCardToList(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusInternalServerError)
        return
    }

    // Parse the JSON request body
    var requestData struct {
        SourceListID       int    `json:"sourceListId"`
        DestinationListID  int    `json:"destinationListId"`
        CardID             int    `json:"cardId"`
        CardName           string `json:"cardName"`
        Position           int    `json:"newPosition"`
    }

    err = json.Unmarshal(body, &requestData)
    if err != nil {
        http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
        fmt.Printf("error : %s\n", err)
        return
    }

    if requestData.SourceListID != requestData.DestinationListID {
        tx, err := db.Begin()
        if err != nil {
            http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
            return
        }
        defer tx.Rollback()

        _, err = tx.Exec("UPDATE cards SET list_id = $1, position = $2 WHERE id = $3",
            requestData.DestinationListID, requestData.Position, requestData.CardID)
        if err != nil {
            http.Error(w, "Failed to move card", http.StatusInternalServerError)
            fmt.Printf("error : %s\n", err)
            return
        }

        err = tx.Commit()
        if err != nil {
            http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
            return
        }
    }

    w.WriteHeader(http.StatusOK)
}

