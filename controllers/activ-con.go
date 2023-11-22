package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-manager/model"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func CreateActivity(w http.ResponseWriter, r *http.Request) {


	// Get the card ID from the URL parameter
	cardID, err := strconv.Atoi(mux.Vars(r)["cardID"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	var requestData struct {
		Message string `json:"message"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON data: %s", err), http.StatusBadRequest)
		return
	}

	var newActivityID int

	createdAt := time.Now()
	newActivity := model.Activity{
		Message:   requestData.Message,
		CardID:    cardID,
		CreatedAt: createdAt,
	}

	err = db.QueryRow("INSERT INTO activities (message, card_id, created_at) VALUES ($1, $2, $3) RETURNING id",
		newActivity.Message, newActivity.CardID, newActivity.CreatedAt).Scan(&newActivityID)

	// Respond with a success message or the newly created activity
	responseData := map[string]interface{}{
		"message":  "Activity created successfully",
		"activity": newActivity,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseData)
}


func GetAllActivities(w http.ResponseWriter, r *http.Request) {
	// Get the card ID from the URL parameter
	cardID, err := strconv.Atoi(mux.Vars(r)["cardID"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	// Query activities for the specific card ID
	activityRows, err := db.Query("SELECT id, message, card_id, created_at FROM activities WHERE card_id = $1", cardID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch activities, %s", err), http.StatusInternalServerError)
		return
	}

	defer activityRows.Close()

	var activities []*model.Activity

	for activityRows.Next() {
		var (
			activityID    int
			activityMessage string
			cardID         int
			createdAt      time.Time
		)

		err := activityRows.Scan(&activityID, &activityMessage, &cardID, &createdAt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning rows, %s", err), http.StatusInternalServerError)
			return
		}

		activity := &model.Activity{
			ID:        activityID,
			Message:   activityMessage,
			CardID:    cardID,
			CreatedAt: createdAt,
		}

		activities = append(activities, activity)
	}

	jsonData, err := json.Marshal(activities)
	if err != nil {
		http.Error(w, "Failed to marshal activity data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
