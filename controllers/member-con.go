package controllers

import (
	"database/sql"
	"project-manager/model"
	"strings"
	"time"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	// "project-manager/model"
	"strconv"

	"github.com/lib/pq"
	_ "github.com/lib/pq"

	// "github.com/codegangsta/gin"
	"github.com/gorilla/mux"
	// "project-manager/model"
)

func GetAllMembers(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	cardID, err := strconv.Atoi(vars["cardID"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	memberRows, err := db.Query("SELECT id, name FROM members WHERE card_id = $1", cardID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch members, %s", err), http.StatusInternalServerError)
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
			http.Error(w, fmt.Sprintf("Error scanning rows, %s", err), http.StatusInternalServerError)
			return
		}

		member := &model.Member{
			ID:   memberID,
			Name: memberName,
		}

		members = append(members, member)

	}

	jsonData, err := json.Marshal(members)
	if err != nil {
		http.Error(w, "Failed to marshal cbecklists data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

func GetAMember(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	cardID, err := strconv.Atoi(vars["cardID"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	memberID, err := strconv.Atoi(vars["memberID"])
	if err != nil {
		http.Error(w, "Invalid member ID", http.StatusBadRequest)
		return
	}

	// Fetch list details
	memberRow := db.QueryRow("SELECT id, name FROM members WHERE id = $1 AND card_id = $2", memberID, cardID)

	var (
		memberName string
	)

	err = memberRow.Scan(&memberID, &memberName)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Member not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to fetch member data, %s", err), http.StatusInternalServerError)
		}
		return
	}

	member := &model.Member{
		ID:   memberID,
		Name: memberName,
	}

	jsonData, err := json.Marshal(member)
	if err != nil {
		http.Error(w, "Failed to marshal checklist data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

func CreateMember(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	cardID, err := strconv.Atoi(vars["cardID"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body, %s", err), http.StatusInternalServerError)
		return
	}

	var requestData struct {
		Name string `json:"name"`
		Email string `json:"email"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON data, %s", err), http.StatusBadRequest)
		return
	}

	var newMemberID int

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create member, %s", err), http.StatusInternalServerError)
		return
	}

	// Create a new card with non-null fields
	newMember := &model.User{
		ID:   newMemberID,
		Name: requestData.Name,
		Email: requestData.Email,
	}

	err = db.QueryRow("INSERT INTO members (name, email, card_id) VALUES ($1, $2, $3) RETURNING id",
		newMember.Name,newMember.Email, cardID).Scan(&newMemberID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert member, %s", err), http.StatusInternalServerError)
		return
	}

	// Fetch the associated list
	var datesArray pq.StringArray
	cardRow := db.QueryRow("SELECT id, name, description, dates FROM cards WHERE id = $1", cardID)
	card := &model.Card{}
	err = cardRow.Scan(&card.ID, &card.Name, &card.Description, &datesArray)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch card data, %s", err), http.StatusInternalServerError)
		return
	}

	var dates []time.Time

	for _, dateString := range datesArray {
		date, err := time.Parse("2006-01-02", dateString)
		if err != nil {
			// Handle the error, e.g., log it or return an error response
		}
		dates = append(dates, date)
	}


	card.Dates = dates

	// Fetch the associated checklists for the card
	checklistsRows, err := db.Query("SELECT id, name FROM checklists WHERE card_id = $1", cardID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch checklists for card, %s", err), http.StatusInternalServerError)
		return
	}
	defer checklistsRows.Close()

	for checklistsRows.Next() {
		checklist := &model.Checklist{}
		err := checklistsRows.Scan(&checklist.ID, &checklist.Name)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to scan checklist data, %s", err), http.StatusInternalServerError)
			return
		}

		// Fetch checklist items for each checklist
		itemsRows, err := db.Query("SELECT id, name, start_date due_date, done FROM items WHERE checklist_id = $1", checklist.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch checklist items for checklist, %s", err), http.StatusInternalServerError)
			return
		}
		defer itemsRows.Close()

		for itemsRows.Next() {
			item := &model.Item{}
			// var assignedTo pq.StringArray
			err := itemsRows.Scan(&item.ID, &item.Name, &item.DueDate, &item.StartDate, &item.Done)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to scan checklist item data, %s", err), http.StatusInternalServerError)
				return
			}

			// item.AssignedTo = []string(assignedTo)

			// Append checklist item to checklist
			checklist.Items = append(checklist.Items, item)
		}

		// Append checklist to card's checklists
		card.Checklists = append(card.Checklists, checklist)
	}

	// Append the new card to the list's cards slice
	card.Members = append(card.Members, newMember)

	jsonData, err := json.Marshal(card)
	if err != nil {
		http.Error(w, "Failed to marshal card data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)

}





func GetMemberByName(w http.ResponseWriter, r *http.Request) {

	
	vars := mux.Vars(r)

	cardID, err := strconv.Atoi(vars["cardID"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}
    // Extract the name parameter from the URL
	name := r.URL.Query().Get("name")

    fmt.Printf("name: %v\n", name)

	if name == "" {
        http.Error(w, "Name parameter is required", http.StatusBadRequest)
        return
    }

    // Fetch members based on the provided name
	rows, err := db.Query("SELECT id, name, email FROM members WHERE name LIKE $1 || '%' AND card_id = $2", name, cardID)

	if err != nil {
        http.Error(w, fmt.Sprintf("Failed to fetch members data, %s", err), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var members []*model.Member

    for rows.Next() {
        var (
            memberID      int
            memberName    string
            memberEmail   string
            // memberBio     string
        )

        if err := rows.Scan(&memberID, &memberName, &memberEmail); err != nil {
            http.Error(w, fmt.Sprintf("Failed to scan member data, %s", err), http.StatusInternalServerError)
            return
        }

        member := &model.Member{
            ID:    memberID,
            Name:  memberName,
            Email: memberEmail,
            // Bio:   memberBio,
        }

        members = append(members, member)
    }

    if err := rows.Err(); err != nil {
        http.Error(w, fmt.Sprintf("Error iterating over members data, %s", err), http.StatusInternalServerError)
        return
    }

    // Respond with the list of matching members
    jsonData, err := json.Marshal(members)
    if err != nil {
        http.Error(w, "Failed to marshal members data", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)
}



func UpdateMember(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["cardID"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	memberID, err := strconv.Atoi(vars["memberID"])
	if err != nil {
		http.Error(w, "Invalid member ID", http.StatusBadRequest)
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
	_, err = db.Exec("UPDATE members SET name = $1 WHERE id = $2", requestData.Name, memberID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update member, %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func DeleteMember(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	memberID, err := strconv.Atoi(vars["memberID"])
	if err != nil {
		http.Error(w, "Invalid member ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM members WHERE id = $1", memberID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete the member, %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}



func SearchMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	cardID, err := strconv.Atoi(vars["cardID"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	// Fetch all members of the specific card
	rows, err := db.Query("SELECT id, name FROM members WHERE card_id = $1", cardID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch members data, %s", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var members []*model.Member

	for rows.Next() {
		var (
			memberID   int
			memberName string
		)

		if err := rows.Scan(&memberID, &memberName); err != nil {
			http.Error(w, fmt.Sprintf("Failed to scan member data, %s", err), http.StatusInternalServerError)
			return
		}

		member := &model.Member{
			ID:   memberID,
			Name: memberName,
		}

		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error iterating over members data, %s", err), http.StatusInternalServerError)
		return
	}

	// Extract the search query from the URL
	searchQuery := r.URL.Query().Get("name")

	// Filter members based on the search query
	var searchResults []*model.Member

	for _, member := range members {
		if strings.Contains(strings.ToLower(member.Name), strings.ToLower(searchQuery)) {
			searchResults = append(searchResults, member)
		}
	}

	// Respond with the list of matching members
	jsonData, err := json.Marshal(searchResults)
	if err != nil {
		http.Error(w, "Failed to marshal members data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
