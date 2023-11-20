import {useEffect, useState} from 'react'
import {Link, useParams} from 'react-router-dom'
// import { useReducer } from 'react';
import './css/card.css'
import React from 'react'
import DatePicker from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';

import { getJwtFromCookie } from './App'
import jwt_decode from 'jwt-decode'



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

let user = findUser();



export const Card = ({card, list}) => {


    
    const newCard = card
    const newList = list
    
    const {name, description, members, checklists} = newCard
    const [cardName, setCardName] = useState(name)
    const [cardDescription, setCardDescription] = useState(description)
    
    // Define state for managing description editing
    const [isEditingDescription, setIsEditingDescription] = useState(false);
    
    // Define state to store the temporary edited description
    const [editedDescription, setEditedDescription] = useState(cardDescription);
    
    
    
    
    
    const [isNewChecklistAddedOrRemoved, setIsNewChecklistAddedOrRemoved] = useState(false)
    const [isAddingChecklist, setIsAddingChecklist] = useState(false)
    const [checklistCardID, setChecklistCardID] = useState('')
    const [cardChecklists, setCardChecklists] = useState(checklists)
    const [checklist, setChecklist] = useState([])

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
            <div>
                <button onClick={() => addNewChecklist()} style={{width:'auto', height:'auto', marginRight: '45px'}} type="button">ذخیره</button>
                <input
                    type="text"
                    value={newChecklistName}
                    onChange={(e) => setNewChecklistName(e.target.value)}
                    style={{width:'200px', height: '60px', marginRight: '45px', direction:'rtl'}}/>
            </div>
        );
        

        
    }


    const removeChecklist = async (id) => {
        try{
            const response = await fetch(`http://localhost:8080/api/lists/${list.id}/cards/${card.id}/checklists/${id}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                },
            });
            if (!response.ok){
                throw new Error("Failed to delete the checklist")
            }

            setIsNewChecklistAddedOrRemoved(true)
            // setCardChecklists(cardChecklists.filter((checklist) => checklist.id !== id))

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

        },[isNewItemAddedOrRemoved])  



        const handleSaveItem = async (newItemName) => {
            try{
                
                const response = await fetch(`http://localhost:8080/api/lists/${newList.id}/cards/${newCard.id}/checklists/${checklist.id}/items`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ name: newItemName }),
                });
                if (!response.ok){
                    throw new Error("Failed to create the item")
                }

                setIsNewItemAddedOrRemoved(true)

                const newItem = await response.json();

                setChecklistItems([...checklistItems, newItem])
                checklist.items = checklistItems
                
                setIsAddingItem(false);
                setNewItemName('');


            } catch (error) {
            console.log("Error creating the item");
                }

        }

        const addNewItem = () => {
            if (newItemName.trim() !== '') {
                handleSaveItem(newItemName)
            }
          }
    
        return (
            <div>
                <button onClick={() => addNewItem()} style={{width:'auto', height:'auto'}}>ذخیره</button>
                <input
                    type="text"
                    value={newItemName}
                    onChange={(e) => setNewItemName(e.target.value)}
                    placeholder="Enter item name"
                    style={{width:'200px', height: '40px', direction:"rtl"}}/>
            </div>
        );
    };
    


    const removeItem = async (checklistID, itemID) => {
        try{
            const response = await fetch(`http://localhost:8080/api/lists/${newList.id}/cards/${newCard.id}/checklists/${checklistID}/items/${itemID}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                },
            });
            if (!response.ok){
                throw new Error("Failed to delete the item")
            }

        } catch (error) {
        console.log("Error deleting the item");
        }
    }




    const [isAddingMember, setIsAddingMember] = useState(false)
    
    const [isNewMemberAddedOrRemoved, setIsNewMemberAddedOrRemoved] =  useState(false)
    
    const [cardMembers, setCardMembers] = useState(members)

    const [memberCardID, setMemberCardID] = useState('')

    

    const AddMember = ({ card, list }) => {

        const [newMemberName, setNewMemberName] = useState('')
        const [matchingUsers, setMatchingUsers] = useState([]);
        const [selectedUser, setSelectedUser] = useState(null);

        // useEffect(() => {

        //     const fetchMembers = async () => {
        //         try {
        //             const response = await fetch(`http://localhost:8080/api/lists/${list.id}/cards/${card.id}/members`,{
        //                 method: 'GET',
        //                 headers: {
        //                     'Content-Type': 'application/json'
        //                 },
        //             });
        //             if (!response.ok){
        //                 throw new Error('Error getting the member')
        //             }
            
        //             const allMembers = await response.json();
        //             setCardMembers(allMembers)
                    
        //         } catch (error) {
        //             console.log('got error : ', error);
        //         }
        
        //     }

        //     if(isNewMemberAddedOrRemoved){
        //         fetchMembers();
        //         setIsNewMemberAddedOrRemoved(false)
        //     }

        // },[isNewMemberAddedOrRemoved])

        console.log('newMemberName : ', newMemberName);
        console.log("selected user : ", selectedUser);

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
                    setMatchingUsers(users);
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
        
        const handleNewMember = async (newMemberName) => {

            // notification when a new member gets added to a card
       
            let user = null; // Define user variable here

            // First try-catch block: Get user info from JWT
            try {
                const jwt = getJwtFromCookie();
                if (jwt) {
                    const decoded = jwt_decode(jwt);
                    user = decoded; // Update user data from the JWT
                    console.log(user);
                }
            } catch (error) {
                console.log(error);
            }

            // Second try-catch block: Send a new notification
            // if(selectedUser){
            //     try {
            //             const notifResponse = await fetch(`http://localhost:8080/api/notifs`, {
            //                 method: 'POST',
            //                 headers: {
            //                     'Content-Type': 'application/json',
            //                 },
            //                 body: JSON.stringify({ message: `کاربر جدید بنام "${newMemberName}" به کارت "${card.name}" اضافه شد`, user_id: user.user_id }),
            //             });
    
            //             if (!notifResponse.ok) {
            //                 throw new Error('Error making a new notification');
            //             }
            //     } catch (error) {
            //         console.log(error);
            //     }
            // }



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

                const newMember = await response.json()
                console.log('new member: ', newMember);

                setNewMemberName('')
                // setIsAddingMember(false)

                // setCardMembers([...cardMembers, newMember])
                // setIsNewMemberAddedOrRemoved(true)

            } catch (error) {
                console.log('Error creating the member');
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
                <button onClick={() => {
                    if(memberCardID != ''){
                        setIsAddingMember(false)
                        addNewMember()
                        setMemberCardID('')
                    }
                    }} style={{ width: '50px', height: '40px', marginRight: '10px' }}>
                    ذخیره
                </button>
                <TextField
                    type="text"
                    value={newMemberName}
                    onChange={(e) => setNewMemberName(e.target.value)}
                    placeholder="Enter member name"
                    style={{ width: '150px', height: '40px', direction: 'rtl' }}
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
        // let user = null; // Define user variable here

        // // First try-catch block: Get user info from JWT
        // try {
        //     const jwt = getJwtFromCookie();
        //     if (jwt) {
        //         const decoded = jwt_decode(jwt);
        //         user = decoded; // Update user data from the JWT
        //         console.log(user);
        //     }
        // } catch (error) {
        //     console.log(error);
        // }

        // Second try-catch block: Send a new notification
        try {
            if (user) { // Check if user is available
                const notifResponse = await fetch(`http://localhost:8080/api/notifs`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ message: `کاربر "${user.name}" کارت با نام "${newCard.name}" را پاک کرد`, user_id: user.user_id }),
                });

                if (!notifResponse.ok) {
                    throw new Error('Error making a new notification');
                }
            }
        } catch (error) {
            console.log(error);
        }

        // Third try-catch block: Delete the card
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
    };




    const changeDescription = async (editedDescription) => {


        let user = null; // Define user variable here

        // First try-catch block: Get user info from JWT
        try {
            const jwt = getJwtFromCookie();
            if (jwt) {
                const decoded = jwt_decode(jwt);
                user = decoded; // Update user data from the JWT
                // console.log(user);
            }
        } catch (error) {
            console.log(error);
        }

        // Second try-catch block: Send a new notification
        try {
            // if (user) { // Check if user is available
                const notifResponse = await fetch(`http://localhost:8080/api/notifs`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ message: `کاربر ${user.name} توضیحات کارت "${newCard.name}" را تغییر داد`, user_id: user.user_id }),
                });

                if (!notifResponse.ok) {
                    throw new Error('Error making a new notification');
                }
            // }
        } catch (error) {
            console.log(error);
        }


        // now in here, you can send that newly created notif to every user that is a member of that card
        
        // write the code here...



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


    const changeDates = async () => {

        let editedDatesArray = [editedDates.start, editedDates.end]
        
        try {

            const requestBody = {
                dates: editedDatesArray
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

    
    const [isAssignItemModalOpen, setAssignItemModalOpen] = useState(false);
    

    const AssignItem = ({ listID, cardID, checklistID, itemID }) => {
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
                    const data = await response.json();
                    setSearchResults(data);
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
            const response = await fetch(`http://localhost:8080/api/lists/${listID}/cards/${cardID}/checklists/${checklistID}/items/${itemID}`, {
              method: 'PATCH',
              headers: {
                'Content-Type': 'application/json',
              },
              body: JSON.stringify(requestBody),
            });
            if (!response.ok) {
              console.error('Failed to update item assignedto array on the backend');
            }
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
              <div className="assignment-modal" style={{ marginLeft: '20px', marginRight: '20px', padding: '10px', position:'relative'}}>
                <button
                  onClick={() => {
                    // Add any additional conditions or actions before closing the modal
                    setAssignItemModalOpen(false);
                    setSearchResults([]);
                  }}
                  style={{ width: '50px', height: '40px', marginRight: '10px' }}
                >
                  Save
                </button>
                <TextField
                  type="text"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  placeholder="Search for a member"
                  style={{ width: '100px', height: '25px', direction: 'rtl' }}
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position="start">
                        <SearchIcon />
                      </InputAdornment>
                    ),
                  }}
                />
                
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
      

    // const AssignItem = ({ listID, cardID, checklistID, itemID }) => {
    //     const [searchQuery, setSearchQuery] = useState('');
    //     const [searchResults, setSearchResults] = useState([]);
      
        // useEffect(() => {
        //   const fetchSearchResults = async () => {
        //     if(searchQuery != ''){
        //         try {
        //           const response = await fetch(`http://localhost:8080/api/lists/${listID}/cards/${cardID}/members?name=${searchQuery}`, {
        //             method: 'GET',
        //             headers: {
        //               'Content-Type': 'application/json',
        //             },
        //           });
        //           const data = await response.json();
        //           setSearchResults(data);
        //         } catch (error) {
        //           console.error('Error fetching search results:', error);
        //         }
        //       };
        //     }
      
        //   if (isAssignItemModalOpen) {
        //     fetchSearchResults();
        //   }
        // }, [isAssignItemModalOpen, searchQuery, listID, cardID]);
      
    //     const handleMemberSelect = async (selectedMember) => {
    //       // assign that item to the member
    //       console.log('selectedMember : ', selectedMember);
    //       const requestBody = {
    //         assignedto: selectedMember,
    //       };
    //       try {
    //         const response = await fetch(`http://localhost:8080/api/lists/${listID}/cards/${cardID}/checklists/${checklistID}/items/${itemID}`, {
    //           method: 'PATCH',
    //           headers: {
    //             'Content-Type': 'application/json',
    //           },
    //           body: JSON.stringify(requestBody),
    //         });
    //         if (!response.ok) {
    //           console.error('Failed to update item assignedto array on the backend');
    //         }
    //       } catch (error) {
    //         console.log(error);
    //       }
      
    //       // Close the modal and clear search results
    //       setAssignItemModalOpen(false);
    //       setSearchResults([]);
    //     };
      
    //     return (
    //       <div>
    //         {isAssignItemModalOpen && (
    //           <div className="assignment-modal">
    //             <input type="text" value={searchQuery} onChange={(e) => setSearchQuery(e.target.value)} placeholder="Search" />
    //             <ul>
    //               {searchResults.map((member) => (
    //                 <li key={member.id} onClick={() => handleMemberSelect(member)}>
    //                   {member.name}
    //                 </li>
    //               ))}
    //             </ul>
    //           </div>
    //         )}
    //       </div>
    //     );
    //   };
      

    const handleCheckboxChange = async (checklist, checklistId, itemId, currentDoneValue) => {
        // Update the 'done' attribute on the front-end
        const updatedItems = checklist.items.map(item =>
            item.id === itemId ? { ...item, done: !currentDoneValue } : item
        );
    
        // Update the state or dispatch an action to update the items in your state management system
    
        setChecklist(prevChecklist => ({
            ...prevChecklist,
            items: updatedItems,
        }));
        // Call the changeItem function with the appropriate argument
        // changeItem('checkbox')   ;
    
        try {
            // Send a PATCH request to the backend
            const response = await fetch(`http://localhost:8080/api/lists/${newList.id}/cards/${newCard.id}/checklists/${checklistId}/items/${itemId}`, {
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
        } catch (error) {
            console.error('Error during PATCH request:', error);
        }
    };
    


    const [showOptions, setShowOptions] = useState(false);

    const handleToggleOptions = () => {
        setShowOptions(!showOptions);
    };

    return (
        
        <div>
           
            <div className='card-container'>
                
                <div className='card-details'>
                    <h2 className='card-title' style={{textAlign:'right'}}><img src={require('./icons/list.png')} alt="" style={{width:'50px', height:'50px', marginRight:'-6%' , position:'relative', float:'right'}}/>
                    {cardName} 
                    </h2>
                    <h3 className='list-name' style={{textAlign:'right', marginRight:'40px', marginTop:'10px'}}> لیست:  {newList.name}</h3>

                    <button className='remove-card-button' onClick={() => removeCard()}>حذف کارت</button>
                    

                    <div className='card-members' style={{marginRight:'30px'}}>
                        <img src={require('./icons/members.png')} alt="" style={{width:'24px', height:'24px', marginLeft:'800px', marginTop:'30px', position:'relative', float:'right'}}/>
                        <h3 style={{textAlign:'right', marginRight:'6px'}}>اعضا</h3>
                        <h4 style={{textAlign:'right'}}>
                        {cardMembers && cardMembers.length > 0 ? (
                            cardMembers.map((member, index) => (
                                <span key={index} style={{ marginRight: '8px', textAlign: 'right' }}>
                                 {member.name}
                                </span>
                            ))
                            ) : (
                            <span>بدون عضو</span>
                        )}

                        </h4>


                        {memberCardID === newCard.id && <AddMember card={card} list={list}/>}
                        
                        <button className='add-member-button' onClick={() => {
                            if(memberCardID === ''){
                                setMemberCardID(newCard.id)
                            } else {
                                setMemberCardID('')
                            }
                        }}>اضافه کردن عضو جدید</button>
                    
                    </div>





                    <div className="description-input" style={{marginRight:'30px'}}>
                    <img src={require('./icons/desc.png')} alt="" style={{width:'20px', height:'20px', marginRight:'-35px', marginTop:'30px', marginBottom:'-10%'}}/>
                    <h2 className='section-title' style={{textAlign:'right'}}>توضیحات</h2>
                        <input
                            type="text"
                            className={isEditingDescription ? 'card-description-active' : 'card-description'}
                            value={editedDescription}
                            onFocus={() => setIsEditingDescription(true)}
                            // onBlur={() => {setIsEditingDescription(false)}}
                            // onBlur={cancelDescription}
                            onChange={(e) => setEditedDescription(e.target.value)}
                        style={{direction:'rtl'}}/>
                        {isEditingDescription && (
                            <div className="description-buttons">
                                <button type="submit" onClick={() => {
                                    // setCardDescription(editedDescription)
                                    changeDescription(editedDescription)
                                    setIsEditingDescription(false);
                                }}>ذخیره</button>
                                <button type="submit" onClick={() => {
                                    setEditedDescription(cardDescription);
                                    setIsEditingDescription(false)
                                    }}>لغو</button>
                            </div>
                        )}
                    </div>
                            



                    <div className='showcase-checklists' style={{marginRight:'auto'}}>
                        
                        {cardChecklists && cardChecklists.length > 0 ? (
                        
                        cardChecklists.map((checklist, index) => (

                            <div className='checklist' key={index}>
                                
                                
                                <h2 className='checklist-title'><img src={require('./icons/checklist.png')} alt="" style={{width:'25px', height:'25px', marginBottom:'-5px', marginLeft:'-30px', marginRight:'10px'}}/>
                                {checklist.name}</h2>
                                {checklist.items && checklist.items.length > 0 ? (
                                checklist.items.map((item) => {
                                    // Parse the due date string into a Date object
                                    const dueDate = new Date(item.duedate);
                                
                                    // Extract the month and day from the Date object
                                    const month = dueDate.getMonth() + 1; // Months are zero-based, so add 1
                                    const day = dueDate.getDate();
                                
                                  

                                    return (
                                        <div className="checklist-item" key={item.id}>
                                          <div className="options-container">
                                            <div className="options-toggle" onClick={handleToggleOptions}>
                                              {/* Circular div for 3-dots icon */}
                                              <div className="circle">
                                                <span>...</span>
                                              </div>
                                            </div>
                                            {showOptions && (
                                              <div className="options-dropdown">
                                                <button className="option-button" onClick={() => removeItem(checklist.id, item.id)}>حذف</button>
                                                <button className="option-button" onClick={() => console.log('add date')}>تاریخ</button>
                                              </div>
                                            )}
                                          </div>
                                    
                                          <div className="clock-date-container">
                                            <div className="month-day" style={{ fontSize: '12px', marginLeft:'18px' }}>
                                              <img src={require('./icons/clock-date.png')} alt="" style={{ width: '14px', height: '14px' }} />
                                              {month}/{day}
                                            </div>
                                          </div>



                                          {isAssignItemModalOpen && (
                                            <AssignItem listID={newList.id} cardID={newCard.id} checklistID={checklist.id} itemID={item.id} />
                                            )}
                                            <div className='item-assigned-to'>
                                                <div className='assigned-to' onClick={() => setAssignItemModalOpen(true)} style={{marginLeft:'18px'}}>
                                                    <img
                                                    src={require('./icons/assignedto.png')}
                                                    alt=""
                                                    style={{ width: '14px', height: '14px', cursor: 'pointer' }}
                                                    />
                                                </div>
                                            </div>

                                    



                                          <label htmlFor="item">{item.name}</label>
                                    
                                          <input
                                            type="checkbox"
                                            id="item"
                                            checked={item.done}
                                            onChange={() => {
                                              handleCheckboxChange(checklist, checklist.id, item.id, item.done);
                                              item.done = !item.done;
                                            }}
                                          />
                                        </div>
                                      );





                                })
                                
                                
                                
                                ) : (
                                    <span>بدون آیتم</span>
                                )}


                                {itemChecklistID === checklist.id && isAddingItem && <AddItem checklist={checklist}/> }


                                <button type='button' className='add-item-button' onClick={() => {
                                    if(itemChecklistID === ''){
                                        setIsAddingItem(true)
                                        setItemChecklistID(checklist.id)
                                    } else {
                                        setIsAddingItem(false)
                                        setItemChecklistID('')
                                    } 
                                    
                                }}>اضافه کردن آیتم</button>
                                

                                <br />
                                <button type='submit' className='remove-checklist-button' onClick={() => removeChecklist(checklist.id)}>پاک کردن</button>
                            

                            </div>
                            
                        ))

                        
                            
                        ) : (
                            <span>بدون چکلیست</span>
                        )}
                        

                            {/* here we add the new checklists */}
                            <div className='add-checklist'>
                            
                            {checklistCardID === newCard.id && isAddingChecklist && <AddChecklist card={newCard} list={newList}/>}
                        
                            <button type='button' className='add-checklist-button' onClick={() => {
                                        if(checklistCardID === ''){
                                            setIsAddingChecklist(true)
                                            setChecklistCardID(card.id)
                                        } else {
                                            setIsAddingChecklist(false)
                                            setChecklistCardID('')
                                        } 
                                        
                                    }}>اضافه کردن چکلیست</button>
                                        

                            </div>
                        
                    </div>

                    <div className='add-to-card' style={{width:'200px', height:'auto'}}>
                        
                        <div className='dropdown'>
                        <button className='dropbtn'>اعضا</button>
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
                                <span>No members</span>
                            )}
                            </div>
                        </div>

                        <div className='dropdown'>
                        <button className='dropbtn'>چکلیست</button>
                            <div className='dropdown-content'>
                            {cardChecklists && cardChecklists.length > 0 ? (
                                cardChecklists.map((checklist, index) => {
                                    return (
                                    <div key={index}>
                                         <a href="#">{checklist.name}</a>
                                    </div>
                                    )
                                })

                            ) : (
                                <span>No Checklist</span>
                            )}
                            </div>
                        </div>

                        <div className='dropdown'>
                        <button className='dropbtn'>برچسب</button>
                            <div className='dropdown-content'>
                            </div>
                        </div>



                            <link type="text/css" rel="stylesheet" href="jalalidatepicker.min.css" />
                            <script type="text/javascript" src="jalalidatepicker.min.js"></script>

                        <div className='dropdown'>
                            <button className='dropbtn'>تاریخ</button>
                            <div className='dropdown-content'>
                                <a href="#">{newCard.dates[0]} : شروع</a>
                                <a href="#">{newCard.dates[1]} : پایان</a>

                                {/* Button to open the date picker for the start date */} 
                                {/* <button onClick={() => document.getElementById('start-date-picker').click()}>انتخاب تاریخ شروع</button> */}
                                <p>انتخاب تاریخ شروع</p>
                                <DatePicker
                                id="start-date-picker"
                                selected={editedDates.start}
                                onChange={(date) => handleDateChange(date, 'start')}
                                dateFormat="yyyy-MM-dd"
                                showYearDropdown
                                />

                                {/* Button to open the date picker for the end date */}
                                {/* <button onClick={() => document.getElementById('end-date-picker').click()}>انتخاب تاریخ پایان</button> */}
                                <p>انتخاب تاریخ پایان</p>
                                <DatePicker
                                id="end-date-picker"
                                selected={editedDates.end}
                                onChange={(date) => handleDateChange(date, 'end')}
                                dateFormat="yyyy-MM-dd"
                                showYearDropdown
                                />

                                {/* Button to save the edited dates to the backend */}
                                <button onClick={changeDates}>ثبت</button>
                            </div>
                        </div>

 
                    </div>


                </div>

            </div>
        
        </div>  




);
    
}    

