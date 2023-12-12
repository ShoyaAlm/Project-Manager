import React, {useContext } from 'react'
import {useState} from 'react'
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';
import { BrowserRouter as Router, Route, Switch, Link, useHistory, useParams } from 'react-router-dom';
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
    
    const { boardId } = useParams();

    
    useEffect(() => {
    const fetchData = async () => {
        try {
            const response = await fetch(`http://localhost:8080/api/boards/${boardId}/lists`, {
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

let user;

const List = () => {
    const { lsts, setLsts, setIsNewListAddedOrRemoved } = useContext(ListContext);
    
    const [isAddingCard, setIsAddingCard] = useState(false);    
    
    const [isAddingList, setIsAddingList] = useState(false)
    
    const [newListName, setNewListName] = useState('')

    const { boardId } = useParams();


    const findUser = () => {

        try {
            const jwt = getJwtFromCookie();
            if (jwt) {
                const decoded = jwt_decode(jwt);
                const user1 = decoded;
                // console.log(user1);
                return user1;
            }
        } catch (error) {
            console.log(error);
        }
    }
    
    user = findUser()
    
    

    const AddList = () => {

       

        const handleSaveList = async (newListName) => {
            
            try {
              const response = await fetch(`http://localhost:8080/api/boards/${boardId}/lists`, {
                method: 'POST',
                headers: {
                  'Content-Type': 'application/json',
                },
                body: JSON.stringify({ name: newListName, username: user.name, user_id: user.user_id, user_email: user.email, owner_id: user.user_id })
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
                // handleNotif()
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


     
    const handleDeleteList = async (lst) => {


        
        try {
            const response = await fetch(`http://localhost:8080/api/lists/${lst.id}`,{
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






        const handleSaveCard = async (newCardName) => {
            console.log(user);
            try{
                const response = await fetch(`http://localhost:8080/api/lists/${id}/cards`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ name: newCardName, username: user.name, user_id: user.id, user_email: user.email, owner_id: user.user_id })
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
                // handleNotif()
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
    

    
    const handleDragEnd = (result) => {
      if (!result.destination) {
        return;
      }

      const updatedLists = Array.from(lsts);
      const [movedList] = updatedLists.splice(result.source.index, 1);
      updatedLists.splice(result.destination.index, 0, movedList);

      setLsts(updatedLists);

      // Send a fetch request to update the list positions on the backend
      // Replace 'YOUR_BACKEND_API_ENDPOINT' with the actual endpoint
      fetch(`http://localhost:8080/api/boards/${boardId}/lists/update-lists-order`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          listOrder: updatedLists.map((list) => list.id),
        }),
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error('Failed to update list order');
          }
        })
        .catch((error) => {
          console.error('Error updating list order:', error);
        });
    };
      


    return (
      <DragDropContext onDragEnd={handleDragEnd}>
        <Droppable droppableId="list-container" direction="horizontal">
          {(provided) => (
            <div {...provided.droppableProps} ref={provided.innerRef} style={{ display: 'flex' }}>
              {lsts.map((lst, index) => (
                <Draggable key={lst.id} draggableId={lst.id.toString()} index={index}>
                  {(provided) => (
                    <div {...provided.draggableProps} {...provided.dragHandleProps} ref={provided.innerRef}>
                      <div className="list">
                        {/* Your existing list content */}
                        <h3>{lst.name}</h3>
                        <ShowCards list={lst} />
                        <input
                          type="text"
                          placeholder="add item"
                          onFocus={() => setIsAddingCard(lst.id)}
                          className={isAddingCard === lst.id ? 'add-card-active' : 'add-card'}
                          onChange={(e) => setNewCardName(e.target.value)}
                          style={{
                            margin: '10px',
                            padding: '10px',
                            width: '200px',
                            height: 'auto',
                            border: '2px solid #ccc',
                            borderRadius: '20px',
                          }}
                        />
  
                        {isAddingCard === lst.id && <AddCard id={lst.id} />}
  
                        <br />
                        <button onClick={() => handleDeleteList(lst)} className="remove-button">
                          پاک کردن
                        </button>
                      </div>
                    </div>
                  )}
                </Draggable>
              ))}
  
              {/* AddList component */}
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
  
              {provided.placeholder}
            </div>
          )}
        </Droppable>
      </DragDropContext>
    );


    

  };




// Inside the ShowCards component

const ShowCards = ({ list }) => {

    const [cardList, setCardList] = useState(list.cards)

    const history = useHistory();
    const [selectedCard, setSelectedCard] = useState(null);
    const [modalIsOpen, setModalIsOpen] = useState(false);

    const { boardId } = useParams();

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



    const onDragEnd = (result) => {
        // Check if the drag operation was completed successfully
        if (!result.destination) {
          return;
        }
      
        // Create a copy of the current list's cards for frontend update
        const updatedCards = [...cardList];
      
        // Reorder the cards in the frontend list
        const [movedCard] = updatedCards.splice(result.source.index, 1);
        updatedCards.splice(result.destination.index, 0, movedCard);
      
        // Update the state with the new order of cards on the frontend
        setCardList(updatedCards);
      
        // Prepare data to update the order on the backend
        const updatedOrder = updatedCards.map((card) => card.id);
      
        console.log(updatedOrder);
        // Make an API call to update the card order on the server
        // You need to replace 'your-api-endpoint' with your actual API endpoint
        fetch(`http://localhost:8080/api/boards/${boardId}/lists/${list.id}/update-cards-order`, {
          method: 'PUT', // Assuming you are using the PUT method to update the order
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            listId: list.id, // Assuming you pass the list ID to identify the list
            cardOrder: updatedOrder,
          }),
        })
          .then((response) => {
            if (!response.ok) {
              throw new Error('Failed to update card order on the server');
            }
            // Handle a successful response as needed
          })
          .catch((error) => {
            console.error('Error updating card order on the server:', error);
            // You can handle the error, show a message, or retry the operation
          });
      };
      


      return (
        <DragDropContext onDragEnd={onDragEnd}>
          <Droppable droppableId="cards">
            {(provided) => (
              <div ref={provided.innerRef}>
                {cardList && cardList.length > 0 ? (
                  cardList.map((card, index) => (
                    <Draggable
                      key={card.id.toString()}
                      draggableId={card.id.toString()}
                      index={index}
                    >
                      {(provided) => (
                        <div
                          ref={provided.innerRef}
                          {...provided.draggableProps}
                          {...provided.dragHandleProps}
                          className="card-item"
                        >
                          {/* ... your card content ... */}
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
                            {user ? (
                                user.name === selectedCard.owner.name ? (
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
                                )
                            ) : (
                                <>
                                <p>You need to sign up to view this content.</p>
                                <button onClick={() => {
                                    history.push('/signup')
                                }}>Sign Up</button>
                                <button onClick={closeModal}>Close Modal</button>
                                </>
                            )}
                            </div>
                        )}
                    </Modal>
                </div>
                        </div>
                      )}
                    </Draggable>
                  ))
                ) : (
                  <span></span>
                )}
                {provided.placeholder} {/* Add this line as a placeholder */}

              </div>
            )}
          </Droppable>
        </DragDropContext>
      );


};


Modal.setAppElement("#root")






// original list return section
// return (
//   <div className="list-container">
//       {lsts.map((lst) => (
//           <div key={lst.id} className="list">
//           <h3>{lst.name}</h3>
//           <ShowCards list={lst} />
//               <input
//                   type="text"
//                   placeholder='add item'
//                   onFocus={() => setIsAddingCard(lst.id)}
//                   className={isAddingCard === lst.id ? 'add-card-active' : 'add-card'}
//                   onChange={(e) => setNewCardName(e.target.value)}
//                   style={{margin: '10px', padding: '10px', 
//                   width: '200px', height: 'auto', 
//                   border: '2px solid #ccc', borderRadius: '20px'}}/>
              
              
//               {isAddingCard === lst.id && <AddCard id={lst.id}/>}
              

//               <br />
//               <button onClick={() => handleDeleteList(lst)} className="remove-button">
//                   پاک کردن
//               </button>
//           </div>
//       ))}

//       <div className="add-list-container">
//           <input
//               type="text"
//               placeholder="+ add a list"
//               className="add-list"
//               onFocus={() => setIsAddingList(true)}
//               onChange={(e) => setNewListName(e.target.value)}
//           />
//           {isAddingList === true &&  <AddList/>}
//       </div>


//   </div>
// );