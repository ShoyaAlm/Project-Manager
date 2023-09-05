package controllers

import (
	"database/sql"
	"project-manager/model"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	// "project-manager/model"
	"strconv"

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
	newMember := &model.Member{
		ID:   newMemberID,
		Name: requestData.Name,
	}

	err = db.QueryRow("INSERT INTO members (name, card_id) VALUES ($1, $2) RETURNING id",
		newMember.Name, cardID).Scan(&newMemberID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert member, %s", err), http.StatusInternalServerError)
		return
	}

	// Fetch the associated list
	cardRow := db.QueryRow("SELECT id, name, description, dates FROM cards WHERE id = $1", cardID)
	card := &model.Card{}
	err = cardRow.Scan(&card.ID, &card.Name, &card.Description, &card.Dates)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch card data, %s", err), http.StatusInternalServerError)
		return
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
