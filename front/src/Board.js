import React, {useContext } from 'react'
import {useState} from 'react'

import { BrowserRouter as Router, Route, Switch, Link, useHistory } from 'react-router-dom';
import Modal from 'react-modal'


import { getJwtFromCookie } from './App';
import jwt_decode from 'jwt-decode'

import SimpleCard from './simpleCard';

import { Card } from './Card';

import './css/board.css'
import { useEffect } from 'react';
const ListContext = React.createContext()

export const AllLists = () => {
    const [lsts, setLsts] = useState([]);
    const [isNewListAddedOrRemoved, setIsNewListAddedOrRemoved] = useState(false);
    
    
    useEffect(() => {
    const fetchData = async () => {
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

            const Lists = await response.json();
            setLsts(Lists);

        } catch (error) {
            console.error('Error getting all lists:', error);
        }
    }

    fetchData();
    
        if (isNewListAddedOrRemoved){
            fetchData();
            setIsNewListAddedOrRemoved(false);
        }
    }, [isNewListAddedOrRemoved]);


    return (
        <ListContext.Provider value={{ lsts, setLsts, setIsNewListAddedOrRemoved }}>
            <List />
        </ListContext.Provider>
    );
};

const List = () => {
    const { lsts, setLsts, setIsNewListAddedOrRemoved } = useContext(ListContext);
    
    const [isAddingCard, setIsAddingCard] = useState(false);    
    
    const [isAddingList, setIsAddingList] = useState(false)
    
    const [newListName, setNewListName] = useState('')


    const findUser = () => {
        let user = null; // Define user variable here

        // First try-catch block: Get user info from JWT
        try {
            const jwt = getJwtFromCookie();
            if (jwt) {
                const decoded = jwt_decode(jwt);
                user = decoded; // Update user data from the JWT
                // console.log(user);
                return user;
            }
        } catch (error) {
            console.log(error);
        }
    }

    const AddList = () => {


        let user = findUser(); // Define user variable here

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

        const handleNotif = async () => {
            try{
                const response = await fetch(`http://localhost:8080/api/notifs/${user.user_id}`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ message: `شما لیست "${newListName}" را ساختید `, 
                    user_id: user.user_id, })
                });
                // console.log(user);
                if (!response.ok) {
                    throw new Error('Failed to create a new card');
                  } 
                  
                }
                catch (error) {
                console.log('error : ' , error);
            }
        }



        const handleSaveList = async (newListName) => {
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
          
              setIsNewListAddedOrRemoved(true)

            } catch (error) {
              console.error('Error creating a new list:', error);
            }

            
        };
        
        
        const addNewList = () => {
            if (newListName.trim() !== '') {
                handleNotif()
                handleSaveList(newListName)
                setNewListName('');
                setIsAddingList(false)
            }
          }

         return (
            <div className="add-list-buttons">
                <button type="submit" onClick={() => {
                    setIsAddingList(false)
                }}>لغو</button>
                <button type="submit" onClick={() => addNewList()}>ذخیره</button>
            </div>
        )

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

            setIsNewListAddedOrRemoved(true)

        } catch (error) {
            console.log("Error deleting the list");
        }
    };



    const [newCardName, setNewCardName] = useState('')


    const AddCard = ({id}) => {



        let user = findUser(); // Define user variable here


        const handleNotif = async () => {
            try{
                const response = await fetch(`http://localhost:8080/api/notifs/${user.user_id}`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ message: `شما کارت "${newCardName}" را ساختید `, 
                    user_id: user.user_id })
                });
                if (!response.ok) {
                    throw new Error('Failed to create a new card');
                  } 
                  
                }
                catch (error) {
                console.log('error : ' , error);
            }
        }




        const handleSaveCard = async (newCardName) => {
            try{
                const response = await fetch(`http://localhost:8080/api/lists/${id}/cards`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ name: newCardName, username: user.name, user_id: user.id })
                });
                if (!response.ok) {
                    throw new Error('Failed to create a new card');
                  }
              
                  const newCard = await response.json();
                  const updatedLists = lsts.map((list) => {
                    if (list.id === id) {
                      return {
                        ...list,
                        cards: [...list.cards, newCard],
                      };
                    }
                    return list;
                  });
                  
                  setLsts(updatedLists);
                  


            } catch (error) {
                  console.error('Error creating a new card:', error);
                  // Handle the error as needed
                }      
        };
        
        
        const addNewCard = () => {
            if (newCardName.trim() !== '') {
                handleNotif()
                handleSaveCard(newCardName)
                setIsNewListAddedOrRemoved(true);
                setNewCardName('');
                setIsAddingCard(false);
            }
        }


        return (
            <div className="add-card-buttons">
                <button type="submit" onClick={() => {
                    setIsAddingCard(false)
                }}>لغو</button>
                <button type="submit" onClick={() => addNewCard()}>ذخیره</button>
            </div>
        )

    }
    
   


    return (
        <div className="list-container">
            {lsts.map((lst) => (
                <div key={lst.id} className="list">
                <h3>{lst.name}</h3>
                <ShowCards list={lst} />
                    <input
                        type="text"
                        placeholder='add item'
                        onFocus={() => setIsAddingCard(lst.id)}
                        className={isAddingCard === lst.id ? 'add-card-active' : 'add-card'}
                        onChange={(e) => setNewCardName(e.target.value)}
                        style={{margin: '10px', padding: '10px', 
                        width: '200px', height: 'auto', 
                        border: '2px solid #ccc', borderRadius: '20px'}}/>
                    
                    
                    {isAddingCard === lst.id && <AddCard id={lst.id}/>}
                    

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
                {isAddingList === true &&  <AddList/>}
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
                { user != null && user.name === selectedCard.owner.name ? (
                    <>
                    {/* Display modal content with selectedCard data */}
                    <Card card={selectedCard} list={list}></Card>
                    {/* Rest of the modal content */}
                    <button onClick={closeModal}>Close Modal</button>
                    </>
                 ) : ( 
                    <>
                        <SimpleCard card={selectedCard} list={list}></SimpleCard>

                        <button onClick={closeModal}>Close Modal</button>
                    </> 
                )}
                </div>
                )}
            </Modal>
                </div>

            ))

            ) : (
                <span>بدون کارت</span>
            )}


        </div>
    );
};

Modal.setAppElement("#root")