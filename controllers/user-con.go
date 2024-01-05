package controllers

import (
	"database/sql"
	_ "database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"

	// "time"

	// "log"
	"net/http"
	"project-manager/model"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	// "github.com/lib/pq"
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
			UserID                int
			userName, userEmail, userPassword string
			userBio sql.NullString
		)

		err := userRows.Scan(&UserID, &userName, &userEmail, &userPassword, &userBio)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning rows, %s", err), http.StatusInternalServerError)
			return
		}

		user := &model.User{
			ID:         UserID,
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
									WHERE c.owner_id = $1`, UserID)

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
	
	vars := mux.Vars(r)

	UserID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Fetch list details
	userRow := db.QueryRow("SELECT id, name, password, email, bio FROM users WHERE id = $1", UserID)

	var (
		userName, userPassword, userEmail string
		userBio sql.NullString
	)

	err = userRow.Scan(&UserID, &userName, &userPassword, &userEmail, &userBio)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to fetch user data, %s", err), http.StatusInternalServerError)
		}
		return
	}
	
	var userBioValue string
	if userBio.Valid {
		userBioValue = userBio.String
	}

	user := &model.User{
		ID:   UserID,
		Name: userName,
		Password: userPassword,
		Email: userEmail,
		Bio: userBioValue,
	}


	jsonData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to marshal user data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}


func GetUserByName(w http.ResponseWriter, r *http.Request) {
    // Extract the name parameter from the URL
	name := r.URL.Query().Get("name")

    fmt.Printf("name: %v\n", name)

	if name == "" {
        http.Error(w, "Name parameter is required", http.StatusBadRequest)
        return
    }

    // Fetch users based on the provided name
	rows, err := db.Query("SELECT id, name, email FROM users WHERE name LIKE '%' || $1 || '%'", name)


	if err != nil {
        http.Error(w, fmt.Sprintf("Failed to fetch users data, %s", err), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []*model.User

    for rows.Next() {
        var (
            UserID      int
            userName    string
            userEmail   string
            // userBio     string
        )

        if err := rows.Scan(&UserID, &userName, &userEmail); err != nil {
            http.Error(w, fmt.Sprintf("Failed to scan user data, %s", err), http.StatusInternalServerError)
            return
        }

        user := &model.User{
            ID:    UserID,
            Name:  userName,
            Email: userEmail,
            // Bio:   userBio,
        }

        users = append(users, user)
    }

    if err := rows.Err(); err != nil {
        http.Error(w, fmt.Sprintf("Error iterating over users data, %s", err), http.StatusInternalServerError)
        return
    }

    // Respond with the list of matching users
    jsonData, err := json.Marshal(users)
    if err != nil {
        http.Error(w, "Failed to marshal users data", http.StatusInternalServerError)
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
	UserID, err := strconv.Atoi(vars["UserID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE id = $1", UserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete the user, %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)


}



func UpdateUserInfo(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	UserID, err := strconv.Atoi(vars["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Decode the JSON payload from the request body
	var requestData struct {
		Name     string	        `json:"name"`
		Email    string         `json:"email"`
		Password string 		`json:"password"`
		Bio 	 string 		`json:"bio"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}


	// Hash the new password before updating the user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to hash password: %s", err), http.StatusInternalServerError)
		return
	}

	// Update the user's information in the database
	// Note: This assumes you have a PostgreSQL database connection in the "db" variable
	_, err = db.Exec("UPDATE users SET name=$1, email=$2, password=$3, bio=$4 WHERE id=$5",
		requestData.Name, requestData.Email, string(hashedPassword), requestData.Bio, UserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update user data: %s", err), http.StatusInternalServerError)
		return
	}

	// Respond with the updated user data
	updatedUser := &model.User{
		ID:       UserID,
		Name:     requestData.Name,
		Email:    requestData.Email,
		Bio:      requestData.Bio,
		Password: "", // Do not send the hashed password to the client
	}

	// Encode the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
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

