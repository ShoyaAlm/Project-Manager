package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-manager/model" // Import your model package
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func CreateNotif(w http.ResponseWriter, r *http.Request) {

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

	createdAt := time.Now()
	newNotif := model.Notification{
		Message:   requestData.Message,
		UserID:    requestData.UserID,
		CreatedAt: createdAt,
		Read:      false,
	}


	
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


func GetUserNotifs(w http.ResponseWriter, r *http.Request) {


	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}


	notifRows, err := db.Query("SELECT id, user_id, message, created_at, read FROM notifications WHERE user_id = $1", userID)
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


func MarkAsReadNotifs(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["userID"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    _, err = db.Exec("UPDATE notifications SET read = true WHERE user_id = $1", userID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to mark notifications as read, %s", err), http.StatusInternalServerError)
        return
    }

    // Return a success response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Notifications marked as read successfully"})
}




func DeleteNotif(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	notifID, err := strconv.Atoi(vars["notifID"])
	if err != nil {
		http.Error(w, "Invalid notif ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM notifications WHERE id = $1", notifID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete the notification, %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}