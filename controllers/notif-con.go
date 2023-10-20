package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-manager/model" // Import your model package
	"time"
)

func CreateNotif(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request
	var requestData struct {
		Message string `json:"message"`
		UserID  int `json:"user_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON data: %s", err), http.StatusBadRequest)
		return
	}


	var newNotifID int

	// Create a new notification
	createdAt := time.Now()
	newNotif := model.Notification{
		Message:   requestData.Message,
		UserID:    requestData.UserID,
		CreatedAt: createdAt,
		Read:      false,
	}

	// Insert the new notification into your database (you'll need to implement this part)
	// For example, you can use your database's SQL query or ORM to insert the notification.

	
	err = db.QueryRow("INSERT INTO notifications (user_id, message, created_at, read) VALUES ($1, $2, $3, $4) RETURNING id",
	newNotif.UserID, newNotif.Message, newNotif.CreatedAt, newNotif.Read).Scan(newNotifID)

	// Respond with a success message or the newly created notification
	responseData := map[string]interface{}{
		"message": "Notification created successfully",
		"notification": newNotif,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseData)
}




func GetAllNotifs(w http.ResponseWriter, r *http.Request) {


	notifRows, err := db.Query("SELECT id, user_id, message, created_at, read FROM notifications")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch notifs, %s", err), http.StatusInternalServerError)
		return
	}

	defer notifRows.Close()

	var notifs []*model.Notification

	for notifRows.Next() {
		var (
			notifID, userID   int
			notifMessage string
			createdAt time.Time
			read bool
		)

		err := notifRows.Scan(&notifID, &userID, &notifMessage, &createdAt, &read)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning rows, %s", err), http.StatusInternalServerError)
			return
		}

		notif := &model.Notification{
			ID:   notifID,
			UserID: userID,
			Message: notifMessage,
			CreatedAt: createdAt,
			Read: read,
		}

		notifs = append(notifs, notif)

	}

	jsonData, err := json.Marshal(notifs)
	if err != nil {
		http.Error(w, "Failed to marshal notification data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}
