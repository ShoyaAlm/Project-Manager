import lists from './lists'

import React, {useContext } from 'react'
import {useState} from 'react'

import { BrowserRouter as Router, Route, Switch, Link, useHistory } from 'react-router-dom';
import Modal from 'react-modal'

import { Card } from './Card';

import './board.css'
import { useEffect } from 'react';
const ListContext = React.createContext()

export const AllLists = () => {
    const [lsts, setLsts] = useState([]);

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await fetch('http://localhost:8080/api/lists', {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });

                if (!response.ok) {
                    throw new Error('Failed to get lists');
                }

                const allLists = await response.json();
                // Handle the response as needed, e.g., update your component state
                // allLists contains the data returned from your Go backend
                setLsts(allLists);
            } catch (error) {
                console.error('Error getting all lists:', error);
                // Handle the error as needed
            }
        }

        // Call the async fetchData function
        fetchData();
    }, []); // Empty dependency array to run the effect only once

    const removeList = (id) => {
        setLsts((currentLists) => {
            return currentLists.filter((list) => list.id !== id);
        });
    };

    return (
        // and in here, I'll pass it as a provider
        <ListContext.Provider value={{ removeList, lsts, setLsts }}>
            <List />
        </ListContext.Provider>
    );
};

const List = () => {
    const { removeList, lsts, setLsts } = useContext(ListContext);
    
    const [isAddingCard, setIsAddingCard] = useState(false);
    
    const [newCardName, setNewCardName] = useState('random')
    
    
    const [isAddingList, setIsAddingList] = useState(false)
    
    

    // const [newList, setNewList] = useState({})
    const [newListName, setNewListName] = useState('')
    
    
    const createListOnServer = async (newListName) => {
        try {
          const response = await fetch('http://localhost:8080/api/lists', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name: newListName }), // Send the new list name in the request body
          });
          
          if (!response.ok) {
            throw new Error('Failed to create a new list');
          }
      
          const newList = await response.json();
          // Handle the response as needed, e.g., update your component state
          // newList contains the data returned from your Go backend
        } catch (error) {
          console.error('Error creating a new list:', error);
          // Handle the error as needed
        }
      };

      
      const addNewList = () => {
        if (newListName.trim() !== '') {
            createListOnServer(newListName)
            setNewListName('');
            setIsAddingList(false)
        }
      }



    const createCardOnServer = async (id) => {
        
        try{
            const response = await fetch(`http://localhost:8080/api/lists/${id}/cards`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ name: newCardName })
            });
            if (!response.ok) {
                throw new Error('Failed to create a new card');
              }
          
              const newCard = await response.json();
              // Handle the response as needed, e.g., update your component state
              // newList contains the data returned from your Go backend

        } catch (error) {
              console.error('Error creating a new card:', error);
              // Handle the error as needed
            }      
    };


    const addNewCard = (id) => {
        if (newCardName.trim() !== '') {
            createCardOnServer(id)
            setNewCardName('');
            setIsAddingCard(false);
        }
    }


    const handleDeleteList = async (id) => {
        try {
            const response = await fetch(`http://localhost:8080/api/lists/${id}`,{
                method: "DELETE",
                headers: {
                    "Content-Type": "application/json"
                }
                
            });

            if (!response.ok){
                throw new Error("Failed to delete the list")
            }
        } catch (error) {
            console.log("Error deleting the list");
        }
    };


    return (
        <div className="list-container">
            {lsts.map((lst) => (
                <div key={lst.id} className="list">
                <h3>{lst.name}</h3>
                <ShowCards list={lst} />
                    <input
                        type="text"
                        placeholder='add item'
                        onFocus={() => setIsAddingCard(lst.id)} // Pass the list ID to setIsAddingCard
                        // onBlur={() => setIsAddingCard(null)}
                        className={isAddingCard === lst.id ? 'add-card-active' : 'add-card'}
                        // value={newCardName}        
                        onChange={(e) => setNewCardName(e.target.value)}
                        style={{margin: '10px', padding: '10px', 
                        width: '200px', height: 'auto', 
                        border: '2px solid #ccc', borderRadius: '20px'}}/>
                    
                    
                    {isAddingCard === lst.id && (
                            <div className="add-item-buttons">
                                <button type='submit' onClick={() => {
                                    setNewCardName('')
                                    setIsAddingCard(false)
                                }}>لغو</button>
                                <button type="submit" onClick={() => addNewCard(lst.id)}>ذخیره</button>
                            </div>
                        )}
                    

                    <br />
                    <button onClick={() => handleDeleteList(lst.id)} className="remove-button">
                        پاک کردن
                    </button>
                </div>
            ))}





            <div className="add-list-container">
                <input
                    type="text"
                    placeholder="+ add a list"
                    className="add-list"
                    onFocus={() => setIsAddingList(true)}
                    onChange={(e) => setNewListName(e.target.value)}
                />
                {isAddingList === true && (
                    <div className="add-list-buttons">
                        <button type="submit" onClick={() => {
                            setNewListName('');
                            setIsAddingList(false);
                        }}>لغو</button>
                        <button type="submit" onClick={() => addNewList()}>ذخیره</button>
                    </div>
                )}
            </div>


        </div>
    );
};




// Inside the ShowCards component

const ShowCards = ({ list }) => {


    const history = useHistory();
    const [selectedCard, setSelectedCard] = useState(null);
    const [modalIsOpen, setModalIsOpen] = useState(false);


    const openModal = (card) => {
      setSelectedCard(card);
      setModalIsOpen(true);
    };
  
    const closeModal = () => {
      setSelectedCard(null);
      setModalIsOpen(false);
      history.push('/'); // Redirect to the main page when closing modal
    };


    return (
        <div className="showcards-container">
            <hr />
            {list.cards && list.cards.length > 0 ? (
            list.cards.map((card) => (
                <div key={card.id} className="card-item">
                <div
                        style={{ textDecoration: 'none', color: 'black', cursor: 'pointer' }}
                        onClick={() => openModal(card)}>
                        <h4 style={{ fontFamily: 'sans-serif' }}>{card.name}</h4>
                </div>
                     
                    <div className="icons-container" style={{display:'inline-flex' , direction: 'rtl'}}>
                        
                        <h6><img src={require('./icons/members.png')} alt="members" style={{width:'15px', height:'24x'}} />{card.members && card.members.length}</h6>
                        <img src={require('./icons/desc.png')} alt="desc" style={{width:'15px', height:'24x', marginRight:'20px', marginLeft:'20px'}} />
                        <h6><img src={require('./icons/tasks.png')} alt="tasks" style={{width:'15px', height:'24x'}} />1/2</h6>                

                    </div>
                    <br />
            <Modal isOpen={modalIsOpen} onRequestClose={closeModal} className="modal">
                {selectedCard && (
                <div className="modal-content">
                
                    {/* Display modal content with selectedCard data */}
                    <Card card={selectedCard} list={list}></Card>
                    {/* Rest of the modal content */}
                    <button onClick={closeModal}>Close Modal</button>
                </div>
                )}
            </Modal>
                </div>

            ))

            ) : (
                <span>no cards</span>
            )}


        </div>
    );
};

Modal.setAppElement("#root")