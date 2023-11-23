package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-manager/model"
	"time"

	"github.com/lib/pq"
)



func TableInfo(w http.ResponseWriter, r *http.Request){
	

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


		cardRows, err := db.Query(`SELECT c.id AS card_id, c.name AS card_name, c.dates AS card_dates, c.label AS card_label
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
				cardID  		    	  int
				cardName				  string
				cardDates                 pq.StringArray
				cardLabel				  *string
			)
			err := cardRows.Scan(&cardID, &cardName, &cardDates, &cardLabel)
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

			card := &model.Card{
				ID:          cardID,
				Name:        cardName,
				Dates:       dates,
				Label: 		 cardLabel,
				Members:     []*model.User{},
			}

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