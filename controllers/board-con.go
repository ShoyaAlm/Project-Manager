package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"project-manager/model"
	"strconv"

	"github.com/gorilla/mux"
)


func CreateBoard(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to read request body, %s", err), http.StatusInternalServerError)
        return
    }

    var requestData struct {
        Name     string `json:"name"`
        UserID   int    `json:"user_id"`
        // Username string `json:"username"`
        // UserEmail string `json:"user_email"`
        // OwnerID  int    `json:"owner_id"`
    }

    err = json.Unmarshal(body, &requestData)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to parse JSON data, %s", err), http.StatusBadRequest)
        return
    }

    var newBoardID int

    // Insert a new board into the 'boards' table
    err = db.QueryRow("INSERT INTO boards (name) VALUES ($1) RETURNING id", requestData.Name).Scan(&newBoardID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to create board, %s", err), http.StatusInternalServerError)
        return
    }

    _, err = db.Exec("INSERT INTO lists (name, board_id, position) VALUES ($1, $2, $3)", "لیست پیش فرض", newBoardID, 1)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to create default list for the board, %s", err), http.StatusInternalServerError)
        return
    }


    _, err = db.Exec("INSERT INTO user_boards (user_id, board_id) VALUES ($1, $2)", requestData.UserID, newBoardID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to insert into the user_boards, %s", err), http.StatusInternalServerError)
        return
    }


	// Fetch the newly created board without including unnecessary attributes
    newBoard := &model.Board{ID: newBoardID, Name: requestData.Name}

    jsonData, err := json.Marshal(newBoard)
    if err != nil {
        http.Error(w, "Failed to marshal response data", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(jsonData)
}

func GetUserBoards(w http.ResponseWriter, r *http.Request) {

    vars := mux.Vars(r)
	UserID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
    
    rows, err := db.Query("SELECT b.id, b.name FROM user_boards ub JOIN boards b ON ub.board_id = b.id WHERE ub.user_id = $1", UserID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to retrieve user boards, %s", err), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var boards []model.Board
    for rows.Next() {
        var board model.Board
        if err := rows.Scan(&board.ID, &board.Name); err != nil {
            http.Error(w, fmt.Sprintf("Failed to scan board data, %s", err), http.StatusInternalServerError)
            return
        }
        boards = append(boards, board)
    }
    if err := rows.Err(); err != nil {
        http.Error(w, fmt.Sprintf("Error iterating over rows, %s", err), http.StatusInternalServerError)
        return
    }

    // Marshal the boards into JSON
    jsonData, err := json.Marshal(boards)
    if err != nil {
        http.Error(w, "Failed to marshal response data", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonData)
}


func AddUserToBoard(w http.ResponseWriter, r *http.Request) {

    vars := mux.Vars(r)
    BoardID, err := strconv.Atoi(vars["board_id"])
	if err != nil {
		http.Error(w, "Invalid Board ID", http.StatusBadRequest)
		return
	}

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to read request body, %s", err), http.StatusInternalServerError)
        return
    }



    var requestData struct {
        UserID        int   `json:"user_id"`
        // BoardID    int      `json:"board_id"`    
    }



    err = json.Unmarshal(body, &requestData)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to parse JSON data, %s", err), http.StatusBadRequest)
        return
    }



    // Retrieve the list of boards for the user
    rows, err := db.Query("SELECT board_id FROM user_boards WHERE user_id = $1", requestData.UserID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to retrieve user boards, %s", err), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var userBoards []int
    for rows.Next() {
        var userBoardID int
        if err := rows.Scan(&userBoardID); err != nil {
            http.Error(w, fmt.Sprintf("Failed to scan user board data, %s", err), http.StatusInternalServerError)
            return
        }
        userBoards = append(userBoards, userBoardID)
    }
    if err := rows.Err(); err != nil {
        http.Error(w, fmt.Sprintf("Error iterating over user board rows, %s", err), http.StatusInternalServerError)
        return
    }

    // If the board is not already in the user's list of boards, add it
    if !contains(userBoards, int(BoardID)) {
        _, err := db.Exec("INSERT INTO user_boards (user_id, board_id) VALUES ($1, $2)", requestData.UserID, BoardID)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to add board to user boards, %s", err), http.StatusInternalServerError)
            return
        }
    }

    w.WriteHeader(http.StatusOK)
}

// Helper function to check if an element is in a slice
func contains(slice []int, element int) bool {
    for _, e := range slice {
        if e == element {
            return true
        }
    }
    return false
}



func DeleteBoard(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    boardID, err := strconv.Atoi(vars["board_id"])
    if err != nil {
        http.Error(w, "Invalid board ID", http.StatusBadRequest)
        return
    }

    var listIDs []int
    rows, err := db.Query("SELECT id FROM lists WHERE board_id = $1", boardID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to fetch list IDs, %s", err), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var listID int
        if err := rows.Scan(&listID); err != nil {
            http.Error(w, fmt.Sprintf("Failed to scan list ID, %s", err), http.StatusInternalServerError)
            return
        }
        listIDs = append(listIDs, listID)
    }

    // loop for deleting the board's data
    for _, listID := range listIDs {
        var cardIDs []int
        cardRows, err := db.Query("SELECT id FROM cards WHERE list_id = $1", listID)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to fetch card IDs, %s", err), http.StatusInternalServerError)
            return
        }
        defer cardRows.Close()

        for cardRows.Next() {
            var cardID int
            if err := cardRows.Scan(&cardID); err != nil {
                http.Error(w, fmt.Sprintf("Failed to scan card ID, %s", err), http.StatusInternalServerError)
                return
            }
            cardIDs = append(cardIDs, cardID)
        }

        // loop for deleting the card's data
        for _, cardID := range cardIDs {


            _, err = db.Exec("DELETE FROM activities WHERE card_id = $1", cardID)
            if err != nil {
                http.Error(w, fmt.Sprintf("Failed to delete activities of card, %s", err), http.StatusInternalServerError)
                return
            }

            // deleting items & checklists
            var checklistIDs []int
            checklistRows, err := db.Query("SELECT id FROM checklists WHERE card_id = $1", cardID)
            if err != nil {
                http.Error(w, fmt.Sprintf("Failed to fetch checklist IDs, %s", err), http.StatusInternalServerError)
                return
            }
            defer checklistRows.Close()

            for checklistRows.Next() {
                var checklistID int
                if err := checklistRows.Scan(&checklistID); err != nil {
                    http.Error(w, fmt.Sprintf("Failed to scan checklist ID, %s", err), http.StatusInternalServerError)
                    return
                }
                checklistIDs = append(checklistIDs, checklistID)
            }

            for _, checklistID := range checklistIDs {

                // Fetch the list of item IDs associated with the checklist
                itemRows, err := db.Query("SELECT id FROM items WHERE checklist_id = $1", checklistID)
                if err != nil {
                    http.Error(w, fmt.Sprintf("Failed to fetch items of checklist, %s", err), http.StatusInternalServerError)
                    return
                }
        
                // Move defer outside the loop
                defer itemRows.Close()
        
                // Iterate through each item and delete associated item_members
                for itemRows.Next() {
                    var itemID int
                    if err := itemRows.Scan(&itemID); err != nil {
                        http.Error(w, fmt.Sprintf("Error scanning item ID, %s", err), http.StatusInternalServerError)
                        return
                    }
        
                    // Delete item_members associated with the current item
                    _, err := db.Exec("DELETE FROM item_members WHERE item_id = $1", itemID)
                    if err != nil {
                        http.Error(w, fmt.Sprintf("Failed to delete item_members for item %d, %s", itemID, err), http.StatusInternalServerError)
                        return
                    }
                }
        
                _, err = db.Exec("DELETE FROM items WHERE checklist_id = $1", checklistID)
                if err != nil {
                    http.Error(w, fmt.Sprintf("Failed to delete items of checklist, %s", err), http.StatusInternalServerError)
                    return
                }
            }
    
            _, err = db.Exec("DELETE FROM members WHERE card_id = $1", cardID)
            if err != nil {
                http.Error(w, fmt.Sprintf("Failed to delete members of list, %s", err), http.StatusInternalServerError)
                return
            }
            // delete checklists
            _, err = db.Exec("DELETE FROM checklists WHERE card_id = $1", cardID)
            if err != nil {
                http.Error(w, fmt.Sprintf("Failed to delete checklists of card, %s", err), http.StatusInternalServerError)
                return
            }
        }

        // delete the cards
        _, err = db.Exec("DELETE FROM cards WHERE list_id = $1", listID)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to delete card, %s", err), http.StatusInternalServerError)
            return
        }

        // delete the list
        _, err = db.Exec("DELETE FROM lists WHERE id = $1", listID)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to delete list, %s", err), http.StatusInternalServerError)
            return
        }
    }

    // Delete the board
    _, err = db.Exec("DELETE FROM boards WHERE id = $1", boardID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to delete board, %s", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusAccepted)

}



func GetAllBoards(w http.ResponseWriter, r *http.Request) {
    // Retrieve all boards from the 'boards' table
    rows, err := db.Query("SELECT id, name FROM boards")
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to retrieve all boards, %s", err), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var boards []model.Board
    for rows.Next() {
        var board model.Board
        if err := rows.Scan(&board.ID, &board.Name); err != nil {
            http.Error(w, fmt.Sprintf("Failed to scan board data, %s", err), http.StatusInternalServerError)
            return
        }
        boards = append(boards, board)
    }
    if err := rows.Err(); err != nil {
        http.Error(w, fmt.Sprintf("Error iterating over rows, %s", err), http.StatusInternalServerError)
        return
    }

    // Marshal the boards into JSON
    jsonData, err := json.Marshal(boards)
    if err != nil {
        http.Error(w, "Failed to marshal response data", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonData)
}
