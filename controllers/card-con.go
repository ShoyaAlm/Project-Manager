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
			cardID                                                                           int
			itemAssignedTo, cardDates                                                        pq.StringArray
			memberIDNullable, checklistIDNullable, itemIDNullable                            sql.NullInt64
			memberNameNullable, checklistNameNullable, itemNameNullable, itemDueDateNullable sql.NullString
		)
		err = rows.Scan(&cardID, &cardDates, &memberIDNullable, &memberNameNullable, &checklistIDNullable, &checklistNameNullable, &itemIDNullable, &itemNameNullable, &itemDueDateNullable, &itemAssignedTo)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, fmt.Sprintf("Card not found, %s", err), http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Failed to fetch card data, %s", err), http.StatusInternalServerError)
			}
			return
		}

		if memberNameNullable.Valid && memberIDNullable.Valid {
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
			itemAssignedTo, cardDates                                                        pq.StringArray
			memberIDNullable, checklistIDNullable, itemIDNullable                            sql.NullInt64
			memberNameNullable, checklistNameNullable, itemNameNullable, itemDueDateNullable sql.NullString
		)
		err = rows.Scan(&cardDates, &memberIDNullable, &memberNameNullable, &checklistIDNullable, &checklistNameNullable, &itemIDNullable, &itemNameNullable, &itemDueDateNullable, &itemAssignedTo)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, fmt.Sprintf("Card not found, %s", err), http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Failed to fetch card data, %s", err), http.StatusInternalServerError)
			}
			return
		}

		if memberNameNullable.Valid && memberIDNullable.Valid {
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
