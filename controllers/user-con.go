package controllers

import (
	"database/sql"
	_ "database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"

	// "log"
	"net/http"
	"project-manager/model"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)


func GetAllUsers(w http.ResponseWriter, r *http.Request){

	
	userRows, err := db.Query("SELECT id, name, email, password, bio FROM users")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch users, %s", err), http.StatusInternalServerError)
		return
	}

	defer userRows.Close()

	var users []*model.User

	for userRows.Next() {
		var (
			userID                int
			userName, userEmail, userPassword string
			userBio sql.NullString
		)

		err := userRows.Scan(&userID, &userName, &userEmail, &userPassword, &userBio)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning rows, %s", err), http.StatusInternalServerError)
			return
		}

		user := &model.User{
			ID:         userID,
			Name:       userName,
			Email:		userEmail,
			Password: 	userPassword,
			// Bio: 		userBio,
		}


		if userBio.Valid { // Check if the bio column is not NULL
            user.Bio = userBio.String
        } else {
            user.Bio = "" // Set to an empty string or handle it as needed
        }


		// SAVE THIS PART FOR WHEN YOU ARE READY TO CHANGE THE CARDS AND GIVE THEM AN ACTUAL OWNER

		cardsRows, err := db.Query(`SELECT c.id AS card_id, c.name AS card_name, c.description AS card_description
									FROM cards c
									WHERE c.owner_id = $1`, userID)

		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch the cards inside user, %s", err), http.StatusInternalServerError)
					return						
		}
								
		defer cardsRows.Close()
										
		
		var cards []*model.Card

		for cardsRows.Next() {
			var (
				cardID int
				cardName, cardDescription string
			)

			err := cardsRows.Scan(&cardID, &cardName, &cardDescription)

			if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning cardsRows, %s", err), http.StatusInternalServerError)
				return
			}

			card := &model.Card{
				ID: cardID,
				Name: cardName,
				Description: cardDescription,
			}

			cards = append(cards, card)

		}

		user.Cards = cards

		users = append(users, user)

	}

	jsonData, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to marshal users data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)


}





func GetUser(w http.ResponseWriter, r *http.Request){
	
	// vars := mux.Vars(r)

	// userID, err := strconv.Atoi(vars["userID"])
	// if err != nil {
	// 	http.Error(w, "Invalid user ID", http.StatusBadRequest)
	// 	return
	// }

	userID := 1

	// Fetch list details
	userRow := db.QueryRow("SELECT id, name, email, bio FROM users WHERE id = $1", userID)

	var (
		userName, userEmail string
		userBio sql.NullString
	)

	err := userRow.Scan(&userID, &userName, &userEmail, &userBio)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to fetch user data, %s", err), http.StatusInternalServerError)
		}
		return
	}

	user := &model.User{
		ID:   userID,
		Name: userName,
		Email: userEmail,
		// Bio: userBio,
	}

	if userBio.Valid { // Check if the bio column is not NULL
		user.Bio = userBio.String
	} else {
		user.Bio = "" // Set to an empty string or handle it as needed
	}



	cardsRows, err := db.Query(`SELECT c.id AS card_id, c.name AS card_name, 
								c.description AS card_description, c.dates AS card_dates
								FROM cards c
								WHERE c.owner_id = $1`, userID)

		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch the cards inside user, %s", err), http.StatusInternalServerError)
					return						
		}
								
		defer cardsRows.Close()
										
		
		var cards []*model.Card

		for cardsRows.Next() {
			var (
				cardID int
				cardName, cardDescription string
				cardDates                 pq.StringArray
			)

			err := cardsRows.Scan(&cardID, &cardName, &cardDescription, &cardDates)

			if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning cardsRows, %s", err), http.StatusInternalServerError)
				return
			}

			card := &model.Card{
				ID: cardID,
				Name: cardName,
				Description: cardDescription,
				Dates:       cardDates,
				Members:     []*model.Member{},
				Checklists:  []*model.Checklist{},
			}

			
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

		user.Cards = cards






	jsonData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to marshal user data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}


func CreateUser(w http.ResponseWriter, r *http.Request){

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body, %s", err), http.StatusInternalServerError)
		return
	}

	var requestData struct {
		Name     string	        `json:"name"`
		Email    string         `json:"email"`
		Password string 		`json:"password"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON data, %s", err), http.StatusBadRequest)
		return
	}

	var newUserID int

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user, %s", err), http.StatusInternalServerError)
		return
	}


	newUser := &model.User{
		ID: newUserID,
		Name: requestData.Name,
		Email: requestData.Email,
		Password: requestData.Password,
		Bio: "",
	}

	err = db.QueryRow("INSERT INTO users (name, email, password, bio) VALUES ($1, $2, $3, $4) RETURNING id",
	newUser.Name, newUser.Email, newUser.Password, newUser.Bio).Scan(&newUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert user, %s", err), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(newUser)
	if err != nil {
		http.Error(w, "Failed to marshal user data", http.StatusInternalServerError)
		return
	}	

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}


func DeleteUser(w http.ResponseWriter, r *http.Request){

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete the user, %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)


}


// func JwtHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost {
// 		// Parse the JSON request
// 		var request model.JwtRequest
// 		decoder := json.NewDecoder(r.Body)
// 		if err := decoder.Decode(&request); err != nil {
// 			http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
// 			return
// 		}

// 		// Access the JWT in request.Jwt
// 		jwt := request.Jwt

// 		// Use the JWT as needed (e.g., for authentication)
// 		// ...

// 		// Send a response
// 		responseData := map[string]interface{}{
// 			"message": "Received JWT",
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		json.NewEncoder(w).Encode(responseData)
// 	} else {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 	}
// }

