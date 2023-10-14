package controllers

import (
	"database/sql"
	_ "database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"project-manager/model"
	"strconv"

	"github.com/gorilla/mux"
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

func SignUp(w http.ResponseWriter, r *http.Request) {


	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body, %s", err), http.StatusInternalServerError)
		return
	}

	var requestData struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Password string `json:"password"`
		Bio string `json:"bio"`
		Cards []*model.Card `json:"cards"`
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
		Password: requestData.Password,
		Email: requestData.Email,
		Bio: requestData.Bio,
		Cards: []*model.Card{},
	}


	if requestData.Name == "" || requestData.Email == "" || requestData.Password == "" {
		http.Error(w, "Name, email, and password are required", http.StatusBadRequest)
        return
	}

	err = db.QueryRow("INSERT INTO users (name, password, email, bio) VALUES ($1, $2, $3, $4) RETURNING id",
			newUser.Name, newUser.Password, newUser.Email, newUser.Bio).Scan(&newUserID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to register user, %s", err), http.StatusInternalServerError)
		return
	}




	jsonData, err := json.Marshal(newUser)
	if err != nil {
		http.Error(w, "Failed to marshal list data", http.StatusInternalServerError)
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


func Login(w http.ResponseWriter, r *http.Request) {

}





