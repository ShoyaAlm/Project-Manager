import {useEffect, useState, useRef} from 'react'
import {Link, useParams} from 'react-router-dom'
// import { useReducer } from 'react';
import './css/card.css'
import React from 'react'
import DatePicker from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';

import { getJwtFromCookie } from './App'
import jwt_decode from 'jwt-decode'

import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';


import SearchIcon from '@mui/icons-material/Search';
import InputAdornment from '@mui/material/InputAdornment';
import TextField from '@mui/material/TextField';


const findUser = () => {

    try {
        const jwt = getJwtFromCookie();
        if (jwt) {
            const decoded = jwt_decode(jwt);
            const user1 = decoded;
            // console.log(user);
            return user1;
        }
    } catch (error) {
        console.log(error);
    }
}

// let user = findUser();



export const Card = ({card, list, userIsMember}) => {


    
    const newCard = card
    const newList = list
    const isUserMember = userIsMember

    const user = findUser()
    
    const {name, description, members, checklists} = newCard
    const [cardName, setCardName] = useState(name)
    const [cardDescription, setCardDescription] = useState(description)
    
    // Define state for managing description editing
    const [isEditingDescription, setIsEditingDescription] = useState(false);
    
    // Define state to store the temporary edited description
    const [editedDescription, setEditedDescription] = useState(cardDescription);
    
    

    // const [isNewActivityAdded, setIsNewActivityAdded] = useState(false)
      
    const createActivity = async (message) => {

                try{
                    const response = await fetch(`http://localhost:8080/api/lists/${list.id}/cards/${card.id}/activity`, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({
                            message: message
                        })
                        
                    });
                    if (!response.ok){
                        throw new Error("Failed to add activity")
                    }
        
                    
                } catch (error) {
                    console.log("Error adding the activity");
                }


    }
    
    
    
    const [isNewChecklistAddedOrRemoved, setIsNewChecklistAddedOrRemoved] = useState(false)
    const [isAddingChecklist, setIsAddingChecklist] = useState(false)
    const [checklistCardID, setChecklistCardID] = useState('')
    const [cardChecklists, setCardChecklists] = useState(checklists)
    const [checklist, setChecklist] = useState([])
    const [showAddChecklist, setShowAddChecklist] = useState(false); // Track whether to show the checklist input

    const AddChecklist = ({ card, list }) => {
        const [newChecklistName, setNewChecklistName] = useState('');

        useEffect(() => {
            
            const fetchChecklists = async () => {
                try {
                    const response = await fetch(`http://localhost:8080/api/lists/${list.id}/cards/${card.id}/checklists`,{
                        method: 'GET',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                    });
                    if (!response.ok){
                        throw new Error('Error getting the checklists')
                    }
            
                    const allChecklists = await response.json();
                    setCardChecklists(allChecklists)

                
                } catch (error) {
                    console.log('got error : ', error);
                }

            }
            

            if(isNewChecklistAddedOrRemoved){
                fetchChecklists();
                setIsNewChecklistAddedOrRemoved(false)
            }

        },[isNewChecklistAddedOrRemoved])
        

        const handleSaveChecklist = async (newChecklistName) => {
            try{

                const response = await fetch(`http://localhost:8080/api/lists/${list.id}/cards/${card.id}/checklists`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ name: newChecklistName }),
                });
                if (!response.ok){
                    throw new Error("Failed to create the checklist")
                }

                setIsNewChecklistAddedOrRemoved(true)
                
                const newChecklist = await response.json();
                
                setCardChecklists([...cardChecklists, newChecklist]);
                

                createActivity(`${user.name} ${newChecklistName} را به کارت اضافه کرد`)

                
                setIsAddingChecklist(false)
                setNewChecklistName('');

                
            } catch (error) {
                console.log("Error creating the checklist");
            }
            
        }
        
        const addNewChecklist = () => {
            if (newChecklistName.trim() !== '') {
                handleSaveChecklist(newChecklistName)
            }
        }
        


        return (
            <div className="add-checklist-container">
              <h3 className="add-checklist-title">نام چکلیست جدید</h3>
              <div className="add-checklist-input-container">
                <input
                  className="add-checklist-input"
                  type="text"
                  value={newChecklistName}
                  onChange={(e) => setNewChecklistName(e.target.value)}
                  placeholder="نام چکلیست جدید"
                />
                <button className="add-checklist-button" onClick={() => addNewChecklist()} type="button"
                style={{fontFamily:'shabnam', fontSize:'12px'}}>
                  ذخیره
                </button>
              </div>
            </div>
          );
          
          
          

        
    }


    const removeChecklist = async (checklist) => {
        try{
            const response = await fetch(`http://localhost:8080/api/lists/${list.id}/cards/${card.id}/checklists/${checklist.id}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                },
            });
            if (!response.ok){
                throw new Error("Failed to delete the checklist")
            }
            setIsNewChecklistAddedOrRemoved(true)
            setCardChecklists(cardChecklists.filter((cl) => cl.id !== checklist.id))


            createActivity(`${user.name} ${checklist.name} را از کارت حذف کرد`)


        } catch (error) {
        console.log("Error deleting the checklist");
        }
    }
    



    const [isNewItemAddedOrRemoved, setIsNewItemAddedOrRemoved] = useState(false)
    const [isAddingItem, setIsAddingItem] = useState(false)
    const [itemChecklistID, setItemChecklistID] = useState('')
  
    const AddItem = ({ checklist }) => {
        
        const [checklistItems, setChecklistItems] = useState(checklist.items)
        const [newItemName, setNewItemName] = useState('');
        // setChecklistItems(checklist.items)

        useEffect(() => {
            
            const fetchChecklistItems = async () => {
                try {
                    const response = await fetch(`http://localhost:8080/api/lists/${list.id}/cards/${card.id}/checklists/${checklist.id}/items`,{
                        method: 'GET',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                    });
                    if (!response.ok){
                        throw new Error('Error getting the checklists')
                    }
            
                    const allChecklistItems = await response.json();
                    setChecklistItems(allChecklistItems)
                    
                } catch (error) {
                    console.log('got error : ', error);
                }

            }
            
            

            if(isNewItemAddedOrRemoved){
                fetchChecklistItems();
                setIsNewItemAddedOrRemoved(false)
            }

        },[isNewItemAddedOrRemoved, checklist.id, newList.id, newCard.id])  



        const handleSaveItem = async (newItemName) => {
            try {
                const response = await fetch(`http://localhost:8080/api/lists/${newList.id}/cards/${newCard.id}/checklists/${checklist.id}/items`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ name: newItemName }),
                });
                if (!response.ok) {
                    throw new Error("Failed to create the item");
                }
        
                const newItem = await response.json();
        
                // Avoid direct mutation, create a new array with the updated item
                const updatedItems = [...checklist.items, newItem];
        
                // Update the state with the new items
                setChecklistItems(updatedItems);
        
                // Assuming checklist is part of a larger state (not shown in your code)
                // Find the index of the current checklist in the array
                const checklistIndex = cardChecklists.findIndex((c) => c.id === checklist.id);
        
                if (checklistIndex !== -1) {
                    // Create a new array with the updated checklist
                    const updatedChecklists = [
                        ...cardChecklists.slice(0, checklistIndex),
                        { ...cardChecklists[checklistIndex], items: updatedItems },
                        ...cardChecklists.slice(checklistIndex + 1),
                    ];
        
                    // Update the state with the new checklists
                    setCardChecklists(updatedChecklists);
        
                    // Call createActivity after the state has been updated
                    createActivity(`${user.name} ${newItemName} را به ${checklist.name} اضافه کرد`);
                }
        
                setIsNewItemAddedOrRemoved(true);
                setIsAddingItem(false);
                setNewItemName('');
        
            } catch (error) {
                console.log("Error creating the item", error);
            }
        }

        const addNewItem = () => {
            if (newItemName.trim() !== '') {
                handleSaveItem(newItemName)
            }
          }
    
          return (
            <div className="add-item-container">
              <input
                type="text"
                value={newItemName}
                onChange={(e) => setNewItemName(e.target.value)}
                placeholder="نام آیتم را وارد کنید"
                className="add-item-input"
              />
              <button onClick={() => addNewItem()} className="add-item-save">
                ذخیره
              </button>
            </div>
          );
          
    };
    


    const removeItem = async (checklist, item) => {
        // const [checklistItems, setChecklistItems] = useState(checklist.items)

        try{
            const response = await fetch(`http://localhost:8080/api/lists/${newList.id}/cards/${newCard.id}/checklists/${checklist.id}/items/${item.id}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                },
            });
            if (!response.ok){
                throw new Error("Failed to delete the item")
            }


            // Remove the item from the local state
            const updatedItems = checklist.items.filter((checklistItem) => checklistItem.id !== item.id);
            // setChecklistItems(updatedItems);

            // Assuming checklist is part of a larger state (not shown in your code)
            // Find the index of the current checklist in the array
            const checklistIndex = cardChecklists.findIndex((c) => c.id === checklist.id);

            if (checklistIndex !== -1) {
                // Create a new array with the updated checklist
                const updatedChecklists = [
                    ...cardChecklists.slice(0, checklistIndex),
                    { ...cardChecklists[checklistIndex], items: updatedItems },
                    ...cardChecklists.slice(checklistIndex + 1),
                ];

                // Update the state with the new checklists
                setCardChecklists(updatedChecklists);
            }

            createActivity(`${user.name} ${item.name} را از ${checklist.name} حذف کرد`)


        } catch (error) {
        console.log("Error deleting the item");
        }
    }




    const [isAddingMember, setIsAddingMember] = useState(false)
    
    const [isNewMemberAddedOrRemoved, setIsNewMemberAddedOrRemoved] =  useState(false)
    
    const [cardMembers, setCardMembers] = useState(members)

    const [memberCardID, setMemberCardID] = useState('')

    

    const AddMember = ({ card, list }) => {
        const {boardId} = useParams()

        const boardIdAsInt = parseInt(boardId, 10);

        const [newMemberName, setNewMemberName] = useState('')
        const [matchingUsers, setMatchingUsers] = useState([]);
        const [selectedUser, setSelectedUser] = useState(null);

        useEffect(() => {
            
            // Fetch matching users based on the entered name
            const fetchMatchingUsers = async () => {
                try {
                    const response = await fetch(`http://localhost:8080/api/users?name=${encodeURIComponent(newMemberName)}`, {
                        method: 'GET',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                    });
    
                    if (!response.ok) {
                        throw new Error('Error getting matching users');
                    }
    
                    const users = await response.json();
                    setMatchingUsers(() => {
                        const similarUsers = users.filter(user => {
                            // You can use any string matching algorithm here.
                            // For example, here we're using a simple case-insensitive substring match.
                            const lowerCaseName = user.name.toLowerCase();
                            const lowerCaseNewMemberName = newMemberName.toLowerCase();
                            return lowerCaseName.includes(lowerCaseNewMemberName);
                        });
                    
                        // Update state with the filtered users (up to 4)
                        return similarUsers.slice(0, 4);
                    });
                } catch (error) {
                    console.log('Error fetching matching users:', error);
                }
            };
    
            if (newMemberName.trim() !== '') {
                fetchMatchingUsers();
            } else {
                setMatchingUsers([]);
            }
        }, [newMemberName]);
        
        useEffect(() => {
            // Fetch card members when component mounts or when isNewMemberAddedOrRemoved changes
            const fetchCardMembers = async () => {
                try {
                    const response = await fetch(`http://localhost:8080/api/lists/${list.id}/cards/${card.id}/members`, {
                        method: 'GET',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                    });
    
                    if (!response.ok) {
                        throw new Error('Error getting card members');
                    }
    
                    const cardMembersData = await response.json();
                    console.error('cardMembersData: ', cardMembersData);
                    setCardMembers(cardMembersData);
                } catch (error) {
                    console.log('Error fetching card members', error);
                }
            };
            
            if(isNewMemberAddedOrRemoved){
                fetchCardMembers();
                setIsNewMemberAddedOrRemoved(false)
            }
        }, [list.id, card.id, isNewMemberAddedOrRemoved]);

        
        const handleNewMember = async (newMemberName) => {


            console.log('newMemberName in handleNewMember ', newMemberName);
            try {
                const response = await fetch(`http://localhost:8080/api/lists/${list.id}/cards/${card.id}/members`, {
                    method: 'POST',
                    headers: {
                        'Content-type': 'application/json'
                    },
                    body: JSON.stringify({name: selectedUser.name, email: selectedUser.email}),
                });
                if(!response.ok){
                    throw new Error("Failed to create new member")
                }



                const responseData = await response.json();
                console.log('new member response data: ', responseData);

                // Check if the new member details are nested within a 'members' property
                const newMember = responseData.members ? responseData.members[0] : responseData;

                console.log('new member: ', newMember);

                // Update the cardMembers state with the new member using the callback form
                setCardMembers((prevMembers) => [...prevMembers, newMember]);

                setIsNewMemberAddedOrRemoved(true);

                setNewMemberName('');
        
                // Update the isNewMemberAddedOrRemoved state to trigger a re-render
                // setIsAddingMember(false)
                
                // setCardMembers([...cardMembers, newMember])
                // setIsNewMemberAddedOrRemoved(true)

                // createActivity(`${user.name} ${newMemberName} را به کارت اضافه کرد`)

            } catch (error) {
                console.log('Error creating the member');
            }


            // add the newly added user to the board, because they have a card in it

            const requestBody = {
                user_id: selectedUser.id,
                // board_id: boardIdAsInt,
              };


            try {
                
                const response = await fetch(`http://localhost:8080/api/boards/${boardIdAsInt}`, {
                    method: 'PUT',
                    headers: {
                        'Content-type': 'application/json'
                    },
                    body: JSON.stringify({requestBody}),
                });
                if(!response.ok){
                    throw new Error("Failed to add new member to the board")
                }
    
            } catch (error) {
                console.log('Error adding new member to the board');
            }


        }
        
        
        const addNewMember = () => {
            if(newMemberName !== ''){
                handleNewMember(newMemberName)
            }

        }



        const [isListVisible, setListVisibility] = useState(true);

        const handleUserClick = (user) => {
            setNewMemberName(user.name);
            setSelectedUser(user);
            setListVisibility(false); // Hide the list when a user is selected
        };

        return (
            <div className="add-member" style={{ marginLeft: '620px', marginRight: 'auto', padding: '10px', position: 'relative' }}>
                <button
                    onClick={() => {
                        if (memberCardID !== '') {
                        setIsAddingMember(false);
                        addNewMember();
                        setMemberCardID('');
                        }
                    }}
                    className="custom-button" // Add this class name for styling
                    >
                    ذخیره
                </button>

                <TextField
                    type="text"
                    value={newMemberName}
                    onChange={(e) => setNewMemberName(e.target.value)}
                    placeholder="نام را وارد کنید"
                    className="custom-textfield" // Add this class name for styling
                    InputProps={{
                        startAdornment: (
                        <InputAdornment position="start">
                            <SearchIcon />
                        </InputAdornment>
                        ),
                    }}
                    />
                {isListVisible && (
                    <ul
                        style={{
                            listStyleType: 'none',
                            padding: 0,
                            position: 'absolute',
                            top: '100%', // Position the list below the input
                            left: 0,
                            width: '100%', // Make the list the same width as the input
                            boxShadow: '0px 4px 4px rgba(0, 0, 0, 0.1)', // Add a shadow for a visual effect
                            backgroundColor: 'white', // Add a background color
                            zIndex: 1, // Ensure the list is above other elements
                            display: 'flex', // Ensure proper rendering of list items
                            flexDirection: 'column', // Align items vertically
                            textAlign: 'right',
                        }}
                    >
                        {matchingUsers.map((user) => (
                            <li
                                key={user.id}
                                onClick={() => handleUserClick(user)}
                                className="listItem"
                                style={{
                                    cursor: 'pointer',
                                    padding: '8px',
                                    borderBottom: '1px solid #ddd',
                                    borderRadius: '4px',
                                    transition: 'background-color 0.3s',
                                    color: '#333', // Set text color to something visible
                                }}
                            >
                                {user.name}
                            </li>
                        ))}
                    </ul>
                )}
            </div>
        );    

        

    }
    


    const removeMember = async (id) => {

        try {
            const response = await fetch(`http://localhost:8080/api/lists/${newList.id}/cards/${newCard.id}/members/${id}`,{
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                },
            });
            if (!response.ok){
                throw new Error('Error deleting the member')
            }

            setIsNewMemberAddedOrRemoved(true)
            // const listID = newList.id
            // const cardID = newCard.id
            // fetchMembers(listID, cardID)
            // setCardMembers(setCardMembers.filter((member) => member.id !== id))


        } catch (error) {
            console.log(error);
        }

}


    const removeCard = async () => {
        
        try {
            const response = await fetch(`http://localhost:8080/api/lists/${newList.id}/cards/${newCard.id}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                },
            });
            if (!response.ok) {
                throw new Error('Error deleting the card');
            }
        } catch (error) {
            console.log(error);
        }

        window.location.reload();
    };




    const changeDescription = async (editedDescription) => {

        try {

            const requestBody = {
                description: editedDescription // Include the updated description
            };

            const response = await fetch(`http://localhost:8080/api/lists/${newList.id}/cards/${newCard.id}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(requestBody)
            });
            if(!response.ok){
                throw new Error("Failed to change description")
            }

            createActivity(`${user.name} بخش توضیحات را تغییر داد`)

            const updatedCard = await response.json();
            console.log('Updated card:', updatedCard);

        } catch (error) {
            console.log('Error changing the description');
        }
    }

    const [editedDates, setEditedDates] = useState({
        start: new Date(newCard.dates[0]), // Convert to Date object
        end: new Date(newCard.dates[1]) // Convert to Date object
      });
    

      // Define a function to handle changes in the edited dates
      const handleDateChange = (date, fieldName) => {
        setEditedDates({
          ...editedDates,
          [fieldName]: date
        });
      };


    const changeDatesOfCard = async () => {

        let editedDatesArray = [editedDates.start, editedDates.end]
        
        try {

            const requestBody = {
                startdate: editedDatesArray
            }
            
            const backendEndpoint = await fetch(`http://localhost:8080/api/lists/${newList.id}/cards/${newCard.id}`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestBody)
            })
            if(!backendEndpoint.ok){
                throw new Error("Failed to change dates")
            }

            const updatedDates = await backendEndpoint.json();
            console.log('Updated dates:', updatedDates);


        } catch (error) {
            // Handle any errors that occur during the API call
            console.error('Error saving dates:', error);
        }

        
            
    }


    // EDITING ITEM's Dates
    // #####################
    // #####################
    // #####################
    // #####################
    // #####################
    // #####################
    const [isEditingDates, setIsEditingDates] = useState(false);
    const [editedItemDates, setEditedItemDates] = useState({
        start: null,
        end: null,
      });
    

    // const [startDate, setStartDate] = useState(new Date());
    // const [endDate, setEndDate] = useState(new Date);
    const changeDatesOfItem = async (checklistId, itemId, startDate, endDate) => {
        try {
        const requestBody = {
            startdate: startDate,
            enddate: endDate
        };

        const backendEndpoint = await fetch(`http://localhost:8080/api/checklists/${checklistId}/items/${itemId} `, {
            method: 'PATCH',
            headers: {
            'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestBody),
        });

        if (!backendEndpoint.ok) {
            throw new Error('Failed to change dates');
        }

        const updatedDates = await backendEndpoint.json();
        console.log('Updated dates:', updatedDates);
        } catch (error) {
        // Handle any errors that occur during the API call
        console.error('Error saving dates:', error);
        }
  };

  const handleItemDateChange = (date, type) => {
    setEditedItemDates((prevDates) => ({ ...prevDates, [type]: date }));
  };

  const saveEditedDates = (checklistId, itemId) => {
    // Assuming you have access to the item's ID and want to update the dates for that item
    const startDate = editedItemDates.start;
    const endDate = editedItemDates.end;

    // Call the function to change dates
    changeDatesOfItem(checklistId, itemId, startDate, endDate);

    // Reset the editing state after saving
    setIsEditingDates(false);
  };     

  
    // #####################
    // #####################
    // #####################
    // #####################
    // #####################
    // #####################


    
    const [isAssignItemModalOpen, setAssignItemModalOpen] = useState(false);
    
    const [selectedItemId, setSelectedItemId] = useState(null);

    const AssignItem = ({ listID, cardID, checklistID, item }) => {
        // const [isAssignItemModalOpen, setAssignItemModalOpen] = useState(false);
        const [searchQuery, setSearchQuery] = useState('');
        const [searchResults, setSearchResults] = useState([]);
      
        useEffect(() => {
            const fetchSearchResults = async () => {
              if(searchQuery != ''){
                  try {
                    const response = await fetch(`http://localhost:8080/api/lists/${listID}/cards/${cardID}/members?name=${searchQuery}`, {
                      method: 'GET',
                      headers: {
                        'Content-Type': 'application/json',
                      },
                    });
                    const members = await response.json();
                    setSearchResults(() => {
                        const similarUsers = members.filter(member => {
                            // You can use any string matching algorithm here.
                            // For example, here we're using a simple case-insensitive substring match.
                            const lowerCaseName = member.name.toLowerCase();
                            const lowerCaseSearchQuery = searchQuery.toLowerCase();
                            return lowerCaseName.includes(lowerCaseSearchQuery);
                        });
                    
                        // Update state with the filtered users (up to 4)
                        return similarUsers.slice(0, 4);
                    });
                  } catch (error) {
                    console.error('Error fetching search results:', error);
                  }
                };
              }
        
            if (isAssignItemModalOpen) {
              fetchSearchResults();
            }
          }, [isAssignItemModalOpen, searchQuery, listID, cardID]);
        

          
        const handleMemberSelect = async (selectedMember) => {
          // assign that item to the member
          const requestBody = {
            assignedto: selectedMember,
          };
          try {
            const response = await fetch(`http://localhost:8080/api/lists/${listID}/cards/${cardID}/checklists/${checklistID}/items/${item.id}`, {
              method: 'PATCH',
              headers: {
                'Content-Type': 'application/json',
              },
              body: JSON.stringify(requestBody),
            });
            if (!response.ok) {
              console.error('Failed to update item assignedto array on the backend');
            }
            createActivity(`${item.name} به ${selectedMember.name} اختصاص داده شد`)

        
        } catch (error) {
            console.log(error);
          }
      
          // Close the modal and clear search results
          setAssignItemModalOpen(false);
          setSearchResults([]);
        };
      
        return (
          <div>
            {isAssignItemModalOpen && (
              <div className="assignment-modal" style={{ marginLeft: '10px', marginRight: 'auto', position:'relative'}}>

                <div className="button-and-search-container">
                <button
                    onClick={() => {
                    // Add any additional conditions or actions before closing the modal
                    setAssignItemModalOpen(false);
                    setSearchResults([]);
                    }}
                    className="custom-button"
                >
                    ذخیره
                </button>
                <TextField
                    type="text"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    placeholder="نام کاربر..."
                    className="custom-textfield"
                    InputProps={{
                    startAdornment: (
                        <InputAdornment position="start">
                        <SearchIcon />
                        </InputAdornment>
                    ),
                    }}
                />
                </div>


                
                {searchResults.length > 0 && (
                    <div
                    style={{
                        listStyleType: 'none',
                        padding: 0,
                        position: 'absolute',
                        top: '100%', // Position the list below the input
                        left: 0,
                        width: '100%', // Make the list the same width as the input
                        boxShadow: '0px 4px 4px rgba(0, 0, 0, 0.1)', // Add a shadow for a visual effect
                        backgroundColor: 'white', // Add a background color
                        zIndex: 1, // Ensure the list is above other elements
                        display: 'flex', // Ensure proper rendering of list items
                        flexDirection: 'column', // Align items vertically
                        textAlign: 'right',
                    }}
                    >
                    {searchResults.map((member) => (
                        <div
                        key={member.id}
                        onClick={() => handleMemberSelect(member)}
                        className="listItem"
                        style={{
                            cursor: 'pointer',
                            padding: '8px',
                            borderBottom: '1px solid #ddd',
                            borderRadius: '4px',
                            transition: 'background-color 0.3s',
                            color: '#333', // Set text color to something visible
                        }}
                        >
                        {member.name}
                        </div>
                    ))}
                    </div>
                )}


              </div>
            )}
          </div>
        );
      };
      



      const [tooltipContent, setTooltipContent] = useState(null);

      const [hoveredItemId, setHoveredItemId] = useState(null);

      const ItemMembers = ({ itemId }) => {
        const [assignedMembers, setAssignedMembers] = useState([]);
      
        useEffect(() => {
          // Fetch assigned members for the item with the given itemId
          const fetchData = async () => {


            try {
                const response = await fetch(`http://localhost:8080/api/items/${itemId}`, {
                  method: 'GET',
                  headers: {
                    'Content-Type': 'application/json',
                  },
                //   body: JSON.stringify(requestBody),
                });
                if (!response.ok) {
                  console.error('Failed to update item assignedto array on the backend');
                }
              
                
              const data = await response.json();
              setAssignedMembers(data);
    
            
            } catch (error) {
                console.log(error);
              }

          };
      
          fetchData();
        }, [itemId]);
      
        return (
            <div className="item-members">
              {assignedMembers && assignedMembers.length > 0 ? (
                <ul>
                  <span className="tooltip-header">تخصیص داده شده به</span>
                  {assignedMembers.map((member) => (
                    <li key={member.id}>{member.name}</li>
                  ))}
                </ul>
              ) : (
                <span className="tooltip-no-assignment">تاکنون به کسی اختصاص نیافته است</span>
              )}
            </div>
          );
          
          
      };


      const formatCardDate = (dateString) => {
        const options = { year: 'numeric', month: 'long', day: 'numeric' };
        const formattedDate = new Date(dateString).toLocaleDateString('fa-IR', options);
        return formattedDate;
      };


      const formatItemDate = (dateString) => {
        console.log('dateString for item date: ', dateString);
        const options = { year: 'numeric', month: 'long', day: 'numeric' };
        const dateObject = new Date(dateString);
        
        const year = dateObject.toLocaleDateString('fa-IR', { year: 'numeric' });
        const month = dateObject.toLocaleDateString('fa-IR', { month: 'long' });
        const day = dateObject.toLocaleDateString('fa-IR', { day: 'numeric' });
      
        return `${day} ${month} ${year}`;
      };

    const handleCheckboxChange = async (checklist, checklistId, item, currentDoneValue) => {
        // Update the 'done' attribute on the front-end
    const updatedItems = checklist.items.map(currentItem =>
        currentItem.id === item.id ? { ...currentItem, done: !currentDoneValue } : currentItem
    );

    // Update the state or dispatch an action to update the items in your state management system
    setChecklist(prevChecklist => ({
        ...prevChecklist,
        items: updatedItems,
    }));
    
        try {
            // Send a PATCH request to the backend
            const response = await fetch(`http://localhost:8080/api/lists/${newList.id}/cards/${newCard.id}/checklists/${checklistId}/items/${item.id}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    done: !currentDoneValue,
                }),
            });
    
            if (!response.ok) {
                // Handle error
                console.error('Failed to update item on the backend');
            }

            if(!currentDoneValue == true){
                createActivity(`${item.name} مربوط به ${checklist.name} انجام شد`)
            } else {
                createActivity(`${item.name} مربوط به ${checklist.name} به حالت ناتمام بازگشت`)
            }

        } catch (error) {
            console.error('Error during PATCH request:', error);
        }
    };
    


    const [showOptions, setShowOptions] = useState(false);

    const handleToggleOptions = () => {
        setShowOptions(!showOptions);
    };



    const [isLabelOpen, setIsLabelOpen] = useState(false);
    const LabelDropdown = ({ newCard, isLabelOpen, newList, newCardId }) => {
        const labels = ['green', 'yellow', 'orange', 'red', 'purple', 'blue'];
      
        const url = `http://localhost:8080/api/lists/${newList.id}/cards/${newCard.id}`;
        const handleLabelClick = async (label) => {
          try {
            const response = await fetch(url, {
              method: 'PATCH',
              headers: {
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({ label }),
            });
      
            if (response.ok) {
              console.log('adding label has ok response');
            } else {
              console.error('Failed to update label');
            }
          } catch (error) {
            console.error('Error updating label:', error);
          }
        };
      
        return (
          <div>
            {isLabelOpen && (
              <div className="label-dropdown">
                {newCard.label ? (
                  <div className="current-label">
                    <p>: برچسب فعلی</p>
                    <div className={`label-color ${newCard.label}`}></div>
                    <p className="choose-another">انتخاب برچسب دیگر</p>
                  </div>
                ) : (
                  <div className="no-label">
                    <p className="choose-label">انتخاب برچسب</p>
                  </div>
                )}
      
                <div className="label-colors">
                  {labels.map((label) => (
                    <div
                      key={label}
                      className={`label-color ${label}`}
                      onClick={() => handleLabelClick(label)}
                    ></div>
                  ))}
                </div>
              </div>
            )}
          </div>
        );
      };
      



    // const [activityInput, setActivityInput] = useState('');

    const CardActivity = ({ list, card, user }) => {
        const [activityInput, setActivityInput] = useState('');
        const [cardActivities, setCardActivities] = useState([]);
        const inputRef = useRef(null);

        const message = ` ${user.name} : ${activityInput}`
        const submitActivity = async () => {
            try {
              const response = await fetch(`http://localhost:8080/api/lists/${list.id}/cards/${card.id}/activity`, {
                method: 'POST',
                headers: {
                  'Content-Type': 'application/json'
                },
                body: JSON.stringify({ message: message }),
              });
          
              if (!response.ok) {
                throw new Error("Failed to create the activity");
              }
          
              // Clear the input field after successful submission
                setActivityInput('');

                // Fetch card activities after successful submission
                const updatedActivitiesResponse = await fetch(`http://localhost:8080/api/lists/${list.id}/cards/${card.id}/activity`);
                const updatedActivities = await updatedActivitiesResponse.json();
                
                setCardActivities(updatedActivities);
                
            } catch (error) {
              console.error('Error submitting activity:', error);
            }
          };
          


      
          useEffect(() => {
            // Fetch card activities when the component mounts
            fetch(`http://localhost:8080/api/lists/${list.id}/cards/${card.id}/activity`)
              .then(response => response.json())
              .then(data => {
                setCardActivities(data);
              })
              .catch(error => {
                console.error('Error fetching activities:', error);
              });
          }, [list.id, card.id]);
          
          useEffect(() => {
            // Focus the input element when the component mounts or when cardActivities change
            if (inputRef.current) {
              inputRef.current.focus();
            }
          }, [cardActivities]);

          return (
            <>
              <div className="activity-input-container">
                {isUserMember && (
                  <>
                    <input
                      ref={inputRef}
                      type="text"
                      placeholder="Comment"
                      value={activityInput}
                      onChange={(e) => setActivityInput(e.target.value)}
                      style={{ direction: 'rtl' }}
                      className="activity-input"
                    />
                    <button onClick={submitActivity} className="submit-button" style={{ fontFamily: 'vazirmatn', fontSize: '10px' }}>
                      نظر بدهید
                    </button>
                  </>
                )}
              </div>
              {cardActivities && cardActivities.length > 0 && (
                <div className="existing-activities">
                  <ul className="activity-list">
                    {cardActivities.map((activity, index) => (
                      <li key={index} className="activity-message" style={{ fontFamily: 'vazirmatn', fontSize: '15px' }}>
                        {activity.message}
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </>
          );
          
      
      };





    const onChecklistDragEnd = (result, listID, cardID, cardChecklists) => {
    if (!result.destination) {
        return;
    }

    const updatedChecklists = [...cardChecklists];
    const [movedChecklist] = updatedChecklists.splice(result.source.index, 1);
    updatedChecklists.splice(result.destination.index, 0, movedChecklist);

    // Update the position property based on the new order
    const updatedChecklistsWithPosition = updatedChecklists.map((checklist, index) => ({
        ...checklist,
        position: index,
    }));

    // Assuming you have a function to update the card checklists order
    // You might need to update the state for the entire card, including checklists
    setCardChecklists(updatedChecklistsWithPosition);

    // Prepare data to update the order on the backend
    const updatedOrder = updatedChecklistsWithPosition.map((checklist) => checklist.id);

    // Make an API call to update the checklist order on the server
    fetch(`http://localhost:8080/api/lists/${listID}/cards/${cardID}/update-checklists-order`, {
        method: 'PUT', // Assuming you are using the PUT method to update the order
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            checklistOrder: updatedOrder,
        }),
    })
        .then((response) => {
            if (!response.ok) {
                throw new Error('Failed to update checklist order on the server');
            }
            // Handle a successful response as needed
        })
        .catch((error) => {
            console.error('Error updating checklist order on the server:', error);
            // You can handle the error, show a message, or retry the operation
        });
        };


    // drag and drop items
    const onChecklistItemsDragEnd = (result, listID, cardID, draggedChecklist) => {
        const { source, destination } = result;
        
        // Ensure that the dragged checklist is an object and has the necessary properties
        if (!draggedChecklist || !draggedChecklist.id || !draggedChecklist.items) {
          console.error('Invalid checklist data');
          console.log('result: ', result);
          console.log('listID: ', listID);
          console.log('cardID: ', cardID);
          console.log('draggedChecklist: ', draggedChecklist);
          return;
        }
      
        const updatedItems = [...draggedChecklist.items];
        const [movedItem] = updatedItems.splice(source.index, 1);
        updatedItems.splice(destination.index, 0, movedItem);
      
        // Update the position property based on the new order
        const updatedItemsWithPosition = updatedItems.map((item, index) => ({
          ...item,
          position: index,
        }));
      
        // Update the state for the entire checklist in setCardChecklists
        setCardChecklists((prevChecklists) =>
          prevChecklists.map((prevChecklist) =>
            prevChecklist.id === draggedChecklist.id
              ? { ...prevChecklist, items: updatedItemsWithPosition }
              : prevChecklist
          )
        );
      
        // Prepare data to update the order on the backend
        const updatedOrder = updatedItemsWithPosition.map((item) => item.id);
      
        // Make an API call to update the item order on the server
        fetch(`http://localhost:8080/api/lists/${listID}/cards/${cardID}/checklists/${draggedChecklist.id}/update-items-order`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            checklistId: draggedChecklist.id,
            itemOrder: updatedOrder,
          }),
        })
          .then((response) => {
            if (!response.ok) {
              throw new Error('Failed to update item order on the server');
            }
          })
          .catch((error) => {
            console.error('Error updating item order on the server:', error);
          });
      };
      




      const onItemDrop = (result, sourceChecklist, destinationChecklist, listID, cardID) => {
        if (!result.destination) {
          return;
        }
      
        const { source, destination } = result;
      
        const updatedSourceChecklistItems = [...sourceChecklist.items];
        const [movedItem] = updatedSourceChecklistItems.splice(source.index, 1);
      
        const updatedDestinationChecklistItems = [...destinationChecklist.items];
        updatedDestinationChecklistItems.splice(destination.index, 0, movedItem);
      
        // Update position for the source checklist
        const updatedSourceItemsWithPosition = updatedSourceChecklistItems.map((item, index) => ({
          ...item,
          position: index,
        }));
      
        // Update position for the destination checklist
        const updatedDestinationItemsWithPosition = updatedDestinationChecklistItems.map((item, index) => ({
          ...item,
          position: index,
        }));
      
        // Update state to reflect the changes
        setCardChecklists((prevChecklists) =>
          prevChecklists.map((prevChecklist) => {
            if (prevChecklist.id === sourceChecklist.id) {
              return { ...prevChecklist, items: updatedSourceItemsWithPosition };
            } else if (prevChecklist.id === destinationChecklist.id) {
              return { ...prevChecklist, items: updatedDestinationItemsWithPosition };
            } else {
              return prevChecklist;
            }
          })
        );


        // Prepare data to update the order on the backend
        const updatedSourceOrder = updatedSourceItemsWithPosition.map((item) => item.id);
        const updatedDestinationOrder = updatedDestinationItemsWithPosition.map((item) => item.id);

        // Make API calls to update item order on the server
        // You might need to adjust the API endpoint and payload structure based on your backend implementation
        // ...

        // Handle other backend updates if needed
        // ...
        };



      


      const [isCardJoined, setIsCardJoined] = useState(false);
      const handleJoinCard = async () => {
        try {
            const response = await fetch(`http://localhost:8080/api/lists/${list.id}/cards/${card.id}/members`, {
                method: 'POST',
                headers: {
                    'Content-type': 'application/json'
                },
                body: JSON.stringify({name: user.name, email: user.email}),
            });
            if(!response.ok){
                throw new Error("Failed to create new member")
            }

            createActivity(`${user.name} !عضو کارت شد`)
            
            const newMember = await response.json()
            console.log('new member: ', newMember);

            // setNewMemberName('')
            // setIsAddingMember(false)

            // setCardMembers([...cardMembers, newMember])
            // setIsNewMemberAddedOrRemoved(true)

        } catch (error) {
            console.log('Error creating the member');
        }
          setIsCardJoined(true)
      };
      

      const onDragEnd = (result, listID, cardID, draggedChecklist) => {
        
        if (!result.destination) {
            return;
          }
        
          // Handle checklist drag and drop
          if (result.type === 'checklist') {
            const updatedChecklists = [...cardChecklists];
            const [movedChecklist] = updatedChecklists.splice(result.source.index, 1);
            updatedChecklists.splice(result.destination.index, 0, movedChecklist);
        
            // Update the position property based on the new order
            const updatedChecklistsWithPosition = updatedChecklists.map((checklist, index) => ({
              ...checklist,
              position: index,
            }));
        
            // Assuming you have a function to update the card checklists order
            // You might need to update the state for the entire card, including checklists
            setCardChecklists(updatedChecklistsWithPosition);
        
            // Prepare data to update the order on the backend
            const updatedOrder = updatedChecklistsWithPosition.map((checklist) => checklist.id);
        
            // Make an API call to update the checklist order on the server
            fetch(`http://localhost:8080/api/lists/${listID}/cards/${cardID}/update-checklists-order`, {
              method: 'PUT', // Assuming you are using the PUT method to update the order
              headers: {
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({
                checklistOrder: updatedOrder,
              }),
            })
              .then((response) => {
                if (!response.ok) {
                  throw new Error('Failed to update checklist order on the server');
                }
                // Handle a successful response as needed
              })
              .catch((error) => {
                console.error('Error updating checklist order on the server:', error);
                // You can handle the error, show a message, or retry the operation
              });
          }




        // Handle item drag and drop
    else if (result.type === 'checklist-item') {
        const updatedChecklists = [...cardChecklists];
        const sourceChecklist = updatedChecklists.find((checklist) => checklist.id === result.source.droppableId);
        const destinationChecklist = updatedChecklists.find((checklist) => checklist.id === result.destination.droppableId);

        // Ensure that the dragged checklist is an object and has the necessary properties
        if (!sourceChecklist || !sourceChecklist.id || !sourceChecklist.items) {
        console.error('Invalid checklist data');
        console.log('result: ', result);
        console.log('listID: ', listID);
        console.log('cardID: ', cardID);
        console.log('sourceChecklist: ', sourceChecklist);
        return;
        }

        const updatedItems = [...sourceChecklist.items];
        const [movedItem] = updatedItems.splice(result.source.index, 1);
        destinationChecklist.items.splice(result.destination.index, 0, movedItem);

        // Update the position property based on the new order
        const updatedItemsWithPosition = destinationChecklist.items.map((item, index) => ({
        ...item,
        position: index,
        }));

        // Update the state for the destination checklist in setCardChecklists
        setCardChecklists((prevChecklists) =>
        prevChecklists.map((prevChecklist) =>
            prevChecklist.id === destinationChecklist.id
            ? { ...prevChecklist, items: updatedItemsWithPosition }
            : prevChecklist
        )
        );

        // Prepare data to update the order on the backend
        const updatedOrder = updatedItemsWithPosition.map((item) => item.id);

        // Make an API call to update the item order on the server
        fetch(`http://localhost:8080/api/lists/${listID}/cards/${cardID}/checklists/${destinationChecklist.id}/update-items-order`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            checklistId: destinationChecklist.id,
            itemOrder: updatedOrder,
        }),
        })
        .then((response) => {
            if (!response.ok) {
            throw new Error('Failed to update item order on the server');
            }
        })
        .catch((error) => {
            console.error('Error updating item order on the server:', error);
        });
        }
      };

      const onInterChecklistItemDragEnd = (result, listID, cardID, cardChecklists) => {
    if (!result.destination) {
        return;
    }

    const { source, destination, draggableId, type } = result;

    console.log('cardChecklists:', cardChecklists);

    if (type === 'item-to-checklist') {
        // Extract numeric part from checklist IDs
        const sourceChecklistId = source.droppableId.match(/\d+/)[0];
        const destinationChecklistId = destination.droppableId.match(/\d+/)[0];

        // Find the source and destination checklists directly
        const sourceChecklist = cardChecklists.find((checklist) => checklist.id === +sourceChecklistId);
        const destinationChecklist = cardChecklists.find((checklist) => checklist.id === +destinationChecklistId);

        // Log source and destination checklists for debugging
        // console.log('Source Checklist:', sourceChecklist);
        // console.log('Destination Checklist:', destinationChecklist);

            if (!sourceChecklist || !destinationChecklist) {
    
                console.error('SourceChecklistID:', sourceChecklistId);
                console.error('DestinationChecklistID:', destinationChecklistId);

                console.error('Invalid checklist data');
    
                console.error('sourceChecklist: ', sourceChecklist);
                console.error('destinationChecklist: ', destinationChecklist);
                
                return;
            }
    
            // Find the moved item in the source checklist
            const movedItem = sourceChecklist.items.find((item) => item.id === +draggableId);
    
            if (!movedItem) {
                console.error('sourceChecklist.items: ', sourceChecklist.items);
                console.error('Invalid item data');
                return;
            }
    
            // Remove the item from the source checklist
            const updatedSourceItems = sourceChecklist.items.filter((item) => item.id !== +draggableId);

            // If destination checklist has no items, add a temporary item
            if (!destinationChecklist.items || destinationChecklist.items.length === 0) {
                destinationChecklist.items = [{ id: 'tempItem', name: 'Temporary Item' }];
            }

            // Rest of the code for moving items...

            // Remove the temporary item after the actual item has been added
            if (destinationChecklist.items.length === 1 && destinationChecklist.items[0].id === 'tempItem') {
                destinationChecklist.items = [];
            }

            // Add the item to the destination checklist
            const updatedDestinationItems = [...destinationChecklist.items];
            updatedDestinationItems.splice(destination.index, 0, movedItem);
    
            // Update the state with the new checklist items
            setCardChecklists((prevChecklists) =>
                prevChecklists.map((prevChecklist) =>
                    prevChecklist.id === sourceChecklist.id
                        ? { ...prevChecklist, items: updatedSourceItems }
                        : prevChecklist.id === destinationChecklist.id
                        ? { ...prevChecklist, items: updatedDestinationItems }
                        : prevChecklist
                )
            );
    
            // Prepare data to update the order on the backend
            const updatedOrder = updatedDestinationItems.map((item) => item.id);
    
            // Make an API call to update the item order on the server
            fetch(
                `http://localhost:8080/api/lists/${listID}/cards/${cardID}/checklists/${destinationChecklistId}/item-to-checklist-order`,
                {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        sourceChecklistID: parseInt(sourceChecklistId, 10), // Convert to integer
                        destinationChecklistID: parseInt(destinationChecklistId, 10),
                        checklistID: parseInt(destinationChecklistId, 10),
                        itemID: movedItem.id,
                        itemName: movedItem.name,
                        position: destination.index,
                    }),
                }
            )
                .then((response) => {
                    if (!response.ok) {
                        throw new Error('Failed to update item order on the server');
                    }
                })
                .catch((error) => {
                    console.error('Error updating item order on the server:', error);
                });
        }
    };
    
      

    return (
        
        <div className='card-container-parent'>
           
            <div className='card-container'>

                <div className='card-details'>
                    
                    
                    {!isUserMember && (
                        <div className={`join-card-button ${isCardJoined ? 'joined' : ''}`}>
                        <button onClick={handleJoinCard} style={{fontFamily:'vazirmatn'}}>
                            {isCardJoined ? '!عضو شدی' : 'عضو بشو'}
                        </button>
                        </div>
                    )}            

                    <h2 className='card-title' style={{textAlign:'right'}}><img src={require('./icons/list.png')} alt="" style={{width:'50px', height:'50px', marginRight:'-6%' , position:'relative', float:'right'}}/>
                    {cardName} 
                    </h2>
                    <h3 className='list-name' style={{textAlign:'right', marginRight:'40px', marginTop:'10px'}}> لیست:  {newList.name}</h3>


                    {isUserMember && (
                        <>
                            <button className='remove-card-button' onClick={() => removeCard()}>حذف کارت</button>
                        </>
                    )}

                    

                    <div className='card-members' style={{ marginRight: '30px', display: 'flex', flexDirection: 'column' }}>
                    <div>
                        <img
                        src={require('./icons/members.png')}
                        alt=""
                        style={{
                            width: '28px',
                            height: '24px',
                            marginLeft: '800px',
                            marginTop: '30px',
                            position: 'relative',
                            float: 'right',
                        }}
                        />
                        <h3 style={{ textAlign: 'right', marginRight: '6px' }}>اعضا</h3>
                        <div className="member-circle-container">
                        {cardMembers && cardMembers.length > 0 ? (
                            cardMembers.map((member, index) => (
                            <div className="member-circle" key={index}>
                                <div className="member-color">{member.name.charAt(0)}</div>
                                <span className="member-name">{member.name}</span>
                            </div>
                            ))
                        ) : (
                            <span>بدون عضو</span>
                        )}
                        </div>
                    </div>

                    {userIsMember && (
                        <>
                        {/* ... other member-specific options ... */}
                        {memberCardID === card.id && <AddMember card={card} list={list} />}
                        <button className='add-member-button' onClick={() => setMemberCardID(memberCardID === card.id ? '' : card.id)}>
                            اضافه کردن عضو جدید
                        </button>
                        </>
                    )}
                    </div>






                    <div className="description-input" style={{ marginRight: '30px' }}>
                        <img src={require('./icons/desc.png')} alt="" style={{ width: '20px', height: '20px', marginRight: '-35px', marginTop: '30px', marginBottom: '-10%' }} />
                        <h2 className='section-title' style={{ textAlign: 'right' }}>توضیحات</h2>

                        {/* Display description inside an input field for both non-members and members */}
                        <input
                        type="text"
                        className={isEditingDescription ? 'card-description-active' : 'card-description'}
                        value={editedDescription}
                        onFocus={() => setIsEditingDescription(true)}
                        readOnly={!userIsMember} // Make the input read-only for non-members
                        onChange={(e) => setEditedDescription(e.target.value)}
                        style={{ direction: 'rtl' }}
                        />

                        {isEditingDescription && userIsMember && (
                        <div className="description-buttons">
                            <button type="submit" onClick={() => {
                            changeDescription(editedDescription);
                            setIsEditingDescription(false);
                            }}>ذخیره</button>
                            <button type="submit" onClick={() => {
                            setEditedDescription(card.description);
                            setIsEditingDescription(false);
                            }}>لغو</button>
                        </div>
                        )}
                    </div>



                    <div className='showcase-checklists' style={{ marginRight: 'auto', maxWidth: '1200px' }}>
                        <div className='activity'>
                            <h3 style={{textAlign:'center', fontFamily:'vazirmatn'}}>فعالیت</h3>
                            <CardActivity list={newList} card={newCard} user={user}/>
                        </div>

                        {cardChecklists && cardChecklists.length > 0 ? (
                            <DragDropContext onDragEnd={(result) => onChecklistDragEnd(result, newList.id, newCard.id, cardChecklists)}>
                                <Droppable droppableId="checklist-container" type="checklist">
                                    {(provided) => (
                                        <div ref={provided.innerRef} {...provided.droppableProps} className="checklists" style={{ display: 'flex', flexDirection: 'column', flex: '1', marginLeft:'100px' }}>
                                            {cardChecklists.map((checklist, index) => (
                                                <Draggable key={checklist.id.toString()} draggableId={`checklist-draggable-${checklist.id}`} index={index}>
                                                    {(provided) => (
                                                    <div ref={provided.innerRef} {...provided.draggableProps} {...provided.dragHandleProps}>
                                                        <div className='checklist'>
                                                            <h2 className='checklist-title'>
                                                                <img src={require('./icons/checklist.png')} alt="" style={{ width: '25px', height: '25px', marginBottom: '-5px', marginLeft: '-30px', marginRight: '10px' }} />
                                                                {checklist.name}
                                                            </h2>
                                                                {/* ... (existing code for items) */}

                                        {/* Add a separate DragDropContext for onChecklistItemsDragEnd */}
                                        <DragDropContext onDragEnd={(result) => onChecklistItemsDragEnd(result, newList.id, newCard.id, checklist)}>
                                                            
                                            {checklist.items && checklist.items.length > 0 ? (
                                            <Droppable droppableId={`checklist-${checklist.id}`} type={`checklist-${checklist.id}`}>
                                            {(provided) => (
                                            <div ref={provided.innerRef} {...provided.droppableProps}>
                                            {checklist.items.map((item, itemIndex) => (
    <Draggable key={item.id.toString()} draggableId={item.id.toString()} index={itemIndex}>
        {(provided) => (
            <div
                ref={provided.innerRef}
                {...provided.draggableProps}
                {...provided.dragHandleProps}
                className="checklist-item"
            >
                {userIsMember && (
                    <>
                        {/* Options container */}
                        <div className="remove-item-container">
                            
                                <div className="remove-item">
                                <button
                                    className="remove-item-button"
                                    onClick={() => removeItem(checklist, item)}
                                >
                                    حذف
                                </button>
                                </div>
                            
                        </div>





                        {/* Clock and date container */}
                <div className="clock-date-container">
                    <div className="month-day" style={{ fontSize: '12px', marginLeft: '18px' }}>
                        <div className='dropdown'>
                        <img
                            src={require('./icons/clock-date.png')}
                            alt=""
                            style={{ width: '14px', height: '14px', cursor: 'pointer' }}
                            onClick={() => setIsEditingDates(!isEditingDates)}
                        />
                        {isEditingDates && (
            <div className="dropdown-content">
                            {/* <a href="#" className="start-date" style={{fontFamily:'vazirmatn', fontSize:'12px'}}>  شروع : {formatCardDate(item.startDate)}</a>
                            <a href="#" className="due-date" style={{fontFamily:'vazirmatn', fontSize:'12px'}}>  پایان : {formatCardDate(item.dueDate)}</a> */}

            <span>
                {item.startDate && `شروع: ${formatItemDate(item.startDate)}`}
                {item.dueDate && ` | پایان: ${formatItemDate(item.dueDate)}`}
            </span>
              <p className="date-picker-text">انتخاب تاریخ شروع</p>
              <DatePicker
                selected={editedItemDates.start}
                onChange={(date) => handleItemDateChange(date, 'start')}
                dateFormat="yyyy-MM-dd"
                showYearDropdown
                className="date-picker-input"
              />

              <p className="date-picker-text">انتخاب تاریخ پایان</p>
              <DatePicker
                selected={editedItemDates.end}
                onChange={(date) => handleItemDateChange(date, 'end')}
                dateFormat="yyyy-MM-dd"
                showYearDropdown
                className="date-picker-input"
              />

              <button onClick={() => saveEditedDates(checklist.id, item.id)} className='submit-new-date-item'>ثبت</button>
            </div>
          )}
          {/* {!isEditingDates && (
            <span>
              {item.startDate && `شروع: ${formatCardDate(item.startDate)}`}
              {item.dueDate && ` | پایان: ${formatCardDate(item.dueDate)}`}
            </span>
          )} */}
                        </div>
                    </div>

                </div>








                        {/* Assign item section */}
                        {isAssignItemModalOpen === item.id && (
                            <AssignItem listID={newList.id} cardID={newCard.id} checklistID={checklist.id} item={item} />
                        )}

                        <div className='item-assigned-to'>
                            <div className='assigned-to' onClick={() => setAssignItemModalOpen(item.id)} style={{ marginLeft: '18px' }}>
                                <img
                                    src={require('./icons/assignedto.png')}
                                    alt=""
                                    style={{ width: '14px', height: '14px', cursor: 'pointer' }}
                                />
                            </div>
                            {/* <ItemMembers itemId={item.id} /> */}
                        </div>
                    </>
                )}

                        
                <label htmlFor={`item-${item.id}`} style={{ fontSize: '14px', position: 'relative' }}
                onMouseEnter={() => {
                    setHoveredItemId(item.id);
                    const timeoutId = setTimeout(() => {
                    setTooltipContent(<ItemMembers itemId={item.id} />);
                    }, 1000); // Adjust the delay time as needed

                    return () => clearTimeout(timeoutId);
                }}
                onMouseLeave={() => {
                    setHoveredItemId(null);
                    setTooltipContent(null);
                }}
                >
                {item.name}
                {hoveredItemId === item.id && tooltipContent && (
                    <div className="custom-tooltip" style={{ position: 'absolute', top: '100%', left: 0, zIndex: 1000 }}>
                    {tooltipContent}
                    </div>
                )}
                </label>

                        {isUserMember && (
                        <>
                            <input
                            type="checkbox"
                            id={`item-${item.id}`}
                            checked={item.done}
                            onChange={() => {
                                handleCheckboxChange(checklist, checklist.id, item, item.done);
                            }}
                            />
                        </>
                        )}

            </div>
        )}
    </Draggable>
))}
                                        {provided.placeholder}
                                        </div>
                                    )}
                                    </Droppable>
                                ) : (
                                    <span className="no-items-message">بدون آیتم</span>
                                )}
                                </DragDropContext>

                                {/* ... (your existing code) */}

                                {itemChecklistID === checklist.id && isAddingItem && <AddItem checklist={checklist}/> }


                                {isUserMember && (

                                    <>
                                    <button
                                        type="button"
                                        className="add-item-button"
                                        onClick={() => {
                                            if (itemChecklistID === '') {
                                            setIsAddingItem(true);
                                            setItemChecklistID(checklist.id);
                                            } else {
                                            setIsAddingItem(false);
                                            setItemChecklistID('');
                                            }
                                        }}
                                        >
                                        اضافه کردن آیتم
                                    </button>
                                    </>

                                )}
                                

                                <br />

                                {isUserMember && (
                                    <>
                                    <button type='submit' className='remove-checklist-button'
                                    style={{fontFamily:'vazirmatn', fontSize:'10px'}}
                                    onClick={() => removeChecklist(checklist)}>پاک کردن</button>
                                    </>
                                )}
                            
                                                                
                                                            </div>
                                                        </div>
                                                    )}
                                                </Draggable>
                                            ))}
                                            {provided.placeholder}
                                        </div>
                                    )}
                                </Droppable>
                            </DragDropContext>
                        ) : (
                            <div style={{ backgroundColor: '#eafbea', padding: '10px',
                             border: '1px solid #4caf50', borderRadius: '5px',
                              display: 'inline-block', width:'290px', textAlign:'center',
                              marginRight:'30px',
                              height:'50px' }}>
                            <span style={{ color: '#4caf50', fontSize: '18px' }}>
                                بدون چکلیست
                            </span>
                        </div>
                        )}
                    </div>


                    







                    <div className='add-to-card' style={{width:'200px', height:'auto', fontFamily:'vazirmatn'}}>
                        
                        <div className='dropdown'>
                        <button className='dropbtn' style={{fontFamily:'shabnam'}}>اعضا</button>
                            <div className='dropdown-content'>
                            {cardMembers && cardMembers.length > 0 ? (
                                cardMembers.map((member, index) => {
                                    return (
                                    <div key={index}>
                                         <span >{member.name}</span>
                                         {/* <button className='remove-member-button' onClick={() => removeMember(member.id)}>X</button> */}
                                    </div>
                                    )
                                })

                            ) : (
                                <span style={{fontFamily:'shabnam'}}>بدون عضو</span>
                            )}
                            </div>
                        </div>


                        <div className='dropdown'>
                            <button
                                className='dropbtn'
                                onClick={() => setShowAddChecklist(true)} // Show the checklist input when the button is clicked
                                style={{fontFamily:'shabnam'}}
                            >
                                چکلیست
                            </button>
                            <div className='dropdown-content'>
                                {showAddChecklist && (
                                    <AddChecklist
                                        card={newCard}
                                        list={newList}
                                        onChecklistAdded={() => {
                                            setIsNewChecklistAddedOrRemoved(true);
                                            setShowAddChecklist(false); // Hide the checklist input after adding
                                        }}
                                    />
                                )}
                            </div>
                        </div>


                        {/* <div className='dropdown'>
                            <button className='dropbtn'>چکلیست</button>
                            <div className='dropdown-content'>
                                <AddChecklist card={newCard} list={newList} />

                            </div>
                        </div> */}





                        <div className='dropdown'>
                        <button className='dropbtn' onClick={() => setIsLabelOpen(!isLabelOpen)}
                         style={{fontFamily:'shabnam'}}
                         >برچسب
                         </button>
                            <div className='dropdown-content'>
                                {isLabelOpen && <LabelDropdown newCard={newCard} isLabelOpen={isLabelOpen} newList={newList} newCardId={newCard.id} />}
                            </div>
                        </div>



                            <link type="text/css" rel="stylesheet" href="jalalidatepicker.min.css" />
                            <script type="text/javascript" src="jalalidatepicker.min.js"></script>

                        <div className='dropdown'>
                            <button className='dropbtn' style={{fontFamily:'shabnam'}}>تاریخ</button>
                            <div className='dropdown-content'>
                            <a href="#" className="start-date" style={{fontFamily:'vazirmatn', fontSize:'15px'}}>  شروع : {formatCardDate(newCard.dates[0])}</a>
                            <a href="#" className="due-date" style={{fontFamily:'vazirmatn', fontSize:'15px'}}>  پایان : {formatCardDate(newCard.dates[1])}</a>

                                {isUserMember && (

                                    <>
                                        <p className="date-picker-text">انتخاب تاریخ شروع</p>
                                        <DatePicker
                                        id="start-date-picker"
                                        selected={editedDates.start}
                                        onChange={(date) => handleDateChange(date, 'start')}
                                        dateFormat="yyyy-MM-dd"
                                        showYearDropdown
                                        className="date-picker-input"
                                        />

                                        {/* Date picker for the end date */}
                                        <p className="date-picker-text">انتخاب تاریخ پایان</p>
                                        <DatePicker
                                        id="end-date-picker"
                                        selected={editedDates.end}
                                        onChange={(date) => handleDateChange(date, 'end')}
                                        dateFormat="yyyy-MM-dd"
                                        showYearDropdown
                                        className="date-picker-input"
                                        />

                                        {/* Button to save the edited dates to the backend */}
                                        <button onClick={changeDatesOfCard} className='changeCardDateButton'>ثبت</button>

                                    </>

                                )}

                            </div>
                        </div>

 
                    </div>


                </div>

            </div>
        
        </div>  




);
    
}    




// items can be dragged and dropped from one checklist into another
/*
<div className='showcase-checklists' style={{ marginRight: 'auto', maxWidth: '1200px' }}>
                        <div className='activity'>
                            <h3 style={{textAlign:'center'}}>فعالیت</h3>
                            <CardActivity list={newList} card={newCard} user={user}/>
                        </div>

                        {cardChecklists && cardChecklists.length > 0 ? (
                        <DragDropContext
                            onDragEnd={(result) => onInterChecklistItemDragEnd(result, newList.id, newCard.id, cardChecklists)}
                        >
                            <Droppable droppableId="checklist-container" type="checklist">
                                {(provided) => (
                                    <div ref={provided.innerRef} {...provided.droppableProps} className="checklists" style={{ display: 'flex', flexDirection: 'column', flex: '1', marginLeft: '100px' }}>
                                        {cardChecklists.map((checklist, checklistIndex) => (
                                            <Draggable key={checklist.id.toString()} draggableId={`checklist-draggable-${checklist.id}`} index={checklistIndex}>
                                                {(provided) => (
                                                    <div ref={provided.innerRef} {...provided.draggableProps} {...provided.dragHandleProps}>
                                                        <div className='checklist'>
                                                            <h2 className='checklist-title'>
                                                                <img src={require('./icons/checklist.png')} alt="" style={{ width: '25px', height: '25px', marginBottom: '-5px', marginLeft: '-30px', marginRight: '10px' }} />
                                                                {checklist.name}
                                                            </h2>

                                                            {/* Add Droppable for items within each checklist */
                                                            /*
                                                            <Droppable droppableId={`checklist-items-${checklist.id}`} type="item-to-checklist">
                                                                {(provided) => (
                                                                    <div ref={provided.innerRef} {...provided.droppableProps} className="checklist-items">
                                                                        {checklist.items && checklist.items.length > 0 ? (
                                                                            checklist.items.map((item, itemIndex) => (
                                                                                <Draggable key={item.id.toString()} draggableId={item.id.toString()} index={itemIndex}>
                                                                                    {(provided) => (
                                                                                        <div
                                                                                            ref={provided.innerRef}
                                                                                            {...provided.draggableProps}
                                                                                            {...provided.dragHandleProps}
                                                                                            className="checklist-item"
                                                                                        >
                                                                            {userIsMember && (
                                                                                <>
                                                                                {/* Options container */
                                                                                /*
                                                                                <div className="options-container">
                                                                                    <div className="options-toggle" onClick={handleToggleOptions}>
                                                                                    <div className="circle">
                                                                                        <span>...</span>
                                                                                    </div>
                                                                                    </div>
                                                                                    {showOptions && (
                                                                                    <div className="options-dropdown">
                                                                                        <button className="option-button" onClick={() => removeItem(checklist, item)}>حذف</button>
                                                                                        <button className="option-button" onClick={() => console.log('add date')}>تاریخ</button>
                                                                                    </div>
                                                                                    )}
                                                                                </div>

                                                                                {/* Clock and date container */
                                                                                /*
                                                                                <div className="clock-date-container">
                                                                                    <div className="month-day" style={{ fontSize: '12px', marginLeft: '18px' }}>
                                                                                    <img src={require('./icons/clock-date.png')} alt="" style={{ width: '14px', height: '14px' }} />
                                                                                    {item.dueDate} {/* Assuming item has a dueDate property */
                                                                                    /*
                                                                                    </div>
                                                                                </div>

                                                                                {/* Assign item section */
                                                                                /*
                                                                                {isAssignItemModalOpen && (
                                                                                    <AssignItem listID={newList.id} cardID={newCard.id} checklistID={checklist.id} itemID={item.id} />
                                                                                )}
                                                                                <div className='item-assigned-to'>
                                                                                    <div className='assigned-to' onClick={() => setAssignItemModalOpen(true)} style={{ marginLeft: '18px' }}>
                                                                                    <img
                                                                                        src={require('./icons/assignedto.png')}
                                                                                        alt=""
                                                                                        style={{ width: '14px', height: '14px', cursor: 'pointer' }}
                                                                                    />
                                                                                    </div>
                                                                                </div>
                                                                                </>
                                                                            )}

                                                                


                                                                            <label htmlFor="item">{item.name}</label>
                                                                            {isUserMember && (
                                                                                <>
                                                                                <input
                                                                                    type="checkbox"
                                                                                    id="item"
                                                                                    checked={item.done}
                                                                                    onChange={() => {
                                                                                    handleCheckboxChange(checklist, checklist.id, item, item.done);
                                                                                    // item.done = !item.done; // Don't mutate state directly
                                                                                    }}
                                                                                />
                                                                                </>
                                                                            )}
                                                                    </div>
                                                                )}
                                                            </Draggable>
                                                        ))
                                                    ) : (
                                                        <Droppable
                                                            droppableId={`checklist-items-${checklist.id}`}
                                                            type="item-to-checklist"
                                                            direction="vertical"
                                                            // Add onDragEnd for handling dragging items within a checklist
                                                            onDragEnd={(result) => onChecklistItemsDragEnd(result, newList.id, newCard.id, checklist)}
                                                        >
                                                            {(provided) => (
                                                                <div ref={provided.innerRef} {...provided.droppableProps} className="checklist-items">
                                                                    {/* Placeholder or message for an empty checklist */
                                                                    /*
                                                                    <div>No items</div>
                                                                    {provided.placeholder}
                                                                </div>
                                                            )}
                                                        </Droppable>
                                                    )}
                                                </div>
                                            )}
                                        </Droppable>
                                        {/* Add a separate DragDropContext for onChecklistItemsDragEnd */

                                {/* ... (your existing code) */}
/*
                                {itemChecklistID === checklist.id && isAddingItem && <AddItem checklist={checklist}/> }


                                {isUserMember && (

                                    <>
                                    <button type='button' className='add-item-button' onClick={() => {
                                        if(itemChecklistID === ''){
                                            setIsAddingItem(true)
                                            setItemChecklistID(checklist.id)
                                        } else {
                                            setIsAddingItem(false)
                                            setItemChecklistID('')
                                        } 
                                        
                                    }}>اضافه کردن آیتم</button>
                                    
                                    </>

                                )}

                                */
                                /*

                                <br />

                                {isUserMember && (
                                    <>
                                    <button type='submit' className='remove-checklist-button' onClick={() => removeChecklist(checklist)}>پاک کردن</button>
                                    </>
                                )}
                            
                                                                
                                    </div>
                                        </div>
                                    )}
                                                </Draggable>
                                            ))}
                                            {provided.placeholder}
                                        </div>
                                    )}
                                </Droppable>
                            </DragDropContext>
                        ) : (
                            <div className="checklist-message">
                            <span style={{ color: 'green', fontSize: '18px' }}>
                                بدون چکلیست
                            </span>
                            </div>
                        )}
                    </div>

                    */








/* make the options for items hidden unless we hover over them

<div
                                                ref={provided.innerRef}
                                                {...provided.draggableProps}
                                                {...provided.dragHandleProps}
                                                className="checklist-item"
                                                onMouseEnter={handleToggleOptions}
                                                onMouseLeave={() => handleToggleOptions(false)}
                                              >
                                                {userIsMember && (
                                                  <>
                                                    {/* Options container */
                                                    /*
                                                    <div className="options-container">
                                                      <div className="options-toggle">
                                                        <div className="circle">
                                                          <span>...</span>
                                                        </div>
                                                      </div>
                                                      {showOptions && (
                                                        <div className="options-dropdown">
                                                          <button className="option-button" onClick={() => removeItem(checklist, item)}>
                                                            حذف
                                                          </button>
                                                          <button className="option-button" onClick={() => console.log('add date')}>
                                                            تاریخ
                                                          </button>
                                                        </div>
                                                      )}
                                                    </div>
                                              
                                                    {/* Clock and date container */
                                                    /*
                                                    <div className="clock-date-container">
                                                      <div className="month-day" style={{ fontSize: '12px', marginLeft: '18px' }}>
                                                        <img src={require('./icons/clock-date.png')} alt="" style={{ width: '14px', height: '14px' }} />
                                                        {item.dueDate} {/* Assuming item has a dueDate property */
                                                        /*
                                                      </div>
                                                    </div>
                                              
                                                    {/* Assign item section */
                                                    /*
                                                    {isAssignItemModalOpen && (
                                                      <AssignItem listID={newList.id} cardID={newCard.id} checklistID={checklist.id} itemID={item.id} />
                                                    )}
                                                    <div className="item-assigned-to">
                                                      <div className="assigned-to" onClick={() => setAssignItemModalOpen(true)} style={{ marginLeft: '18px' }}>
                                                        <img
                                                          src={require('./icons/assignedto.png')}
                                                          alt=""
                                                          style={{ width: '14px', height: '14px', cursor: 'pointer' }}
                                                        />
                                                      </div>
                                                    </div>
                                                  </>
                                                )}
                                              
                                                <label htmlFor="item">{item.name}</label>
                                                {isUserMember && (
                                                  <>
                                                    <input
                                                      type="checkbox"
                                                      id="item"
                                                      checked={item.done}
                                                      onChange={() => {
                                                        handleCheckboxChange(checklist, checklist.id, item, item.done);
                                                      }}
                                                    />
                                                  </>
                                                )}
                                              </div>


*/