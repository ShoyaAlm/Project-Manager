package controllers

import (
	"database/sql"
	_ "database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"project-manager/model"
	"project-manager/secret"
	_ "strconv"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/gorilla/mux"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
)





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
		Bio 	 string `json:"bio"`
		Cards 	 []*model.Card `json:"cards"`
    }



	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON data, %s", err), http.StatusBadRequest)
		return
	}


	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to hash password, %s", err), http.StatusInternalServerError)
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
		Password: string(hashedPassword),
		Email: requestData.Email,
		Bio: "",
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




func Login(w http.ResponseWriter, r *http.Request) {
    // Parse the JSON request
    var requestData struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

	
    err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to parse JSON data: %s", err), http.StatusBadRequest)
        return
    }



	

    // Query the database to find a user with the given email
    user := model.User{}
    err = db.QueryRow("SELECT id, name, email, password FROM users WHERE email = $1", requestData.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
    if err != nil {
        if err == sql.ErrNoRows {
            // User not found
            http.Error(w, "User not found", http.StatusUnauthorized)
            return
        } else {
            log.Printf("Error querying the database: %v", err)
            http.Error(w, "Failed to query the database", http.StatusInternalServerError)
            return
        }
    }

	// Check if the password matches using bcrypt.CompareHashAndPassword
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestData.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			// Passwords do not match
			http.Error(w, "Incorrect password", http.StatusUnauthorized)
			return
		} else {
			log.Printf("Error comparing hashed passwords: %v", err)
			http.Error(w, "Failed to compare passwords", http.StatusInternalServerError)
			return
		}
	}


    // Create a new token
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["user_id"] = user.ID
	claims["name"] = user.Name
    claims["email"] = user.Email

    // Sign the token with the secret key
    tokenString, err := token.SignedString(secret.JwtSecret)
    if err != nil {
        http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
        return
    }

    // Authentication successful
    responseData := map[string]interface{}{
        "message": "Login successful",
        "user":    user,
        "token":   tokenString,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(responseData)
}



