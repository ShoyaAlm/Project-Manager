package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-manager/model"
	"strconv"

	// "github.com/codegangsta/gin"
	"github.com/gorilla/mux"
	// "project-manager/model"
)

var cards = []*model.Card{
	{
		ID:          1,
		Name:        "ds",
		Description: "some Description",
		Dates:       [2]string{"30th august", "30th september"},
		Checklists:  []*model.Checklist{},
		Members:     []*model.Member{},
	},
	{
		ID:          2,
		Name:        "dsdsad",
		Description: "some Description",
		Dates:       [2]string{"30th august", "30th september"},
		Checklists:  []*model.Checklist{},
		Members:     []*model.Member{},
	},
}

func GetAllCards(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Marshal the projects slice into JSON
	jsonData, err := json.Marshal(cards)
	if err != nil {
		http.Error(w, "Failed to marshal cards data", http.StatusInternalServerError)
		return
	}

	// Write the JSON data to the response
	w.Write(jsonData)

}

func GetACard(w http.ResponseWriter, r *http.Request) {
	// Parse the list ID from the request URL
	vars := mux.Vars(r)
	cardID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	// Find the list with the given ID in your 'cards' slice
	var foundCard *model.Card
	for _, card := range cards {
		if card.ID == cardID {
			foundCard = card
			break
		}
	}

	if foundCard == nil {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	// Marshal the found list into JSON
	jsonData, err := json.Marshal(foundCard)
	if err != nil {
		http.Error(w, "Failed to marshal card data", http.StatusInternalServerError)
		return
	}

	// Write the JSON data to the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func CreateCard(w http.ResponseWriter, r *http.Request) {

	var newCard model.Card

	err := json.NewDecoder(r.Body).Decode(&newCard)
	if err != nil {
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	newCard.ID = len(cards) + 1
	if newCard.Checklists == nil {
		newCard.Checklists = []*model.Checklist{}
	}
	cards = append(cards, &newCard)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Card created successfully")

}

func UpdateCard(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	cardID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	var updatedCard model.Card
	err = json.NewDecoder(r.Body).Decode(&updatedCard)
	if err != nil {
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
	}

	found := false
	for i, card := range cards {
		if card.ID == cardID {
			updatedCard.ID = card.ID
			updatedCard.Checklists = card.Checklists
			cards[i] = &updatedCard
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Card updated successfully")
}

func DeleteCard(w http.ResponseWriter, r *http.Request) {
	// Parse the card ID from the request URL or request body
	vars := mux.Vars(r)
	cardID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid list ID", http.StatusBadRequest)
		return
	}

	// Find and remove the card with the given ID from your 'cards' slice
	found := false
	for i, card := range cards {
		if card.ID == cardID {
			cards = append(cards[:i], cards[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Card deleted successfully")
}
