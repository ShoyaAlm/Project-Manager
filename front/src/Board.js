import React, {useContext } from 'react'
import {useState} from 'react'
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';
import { BrowserRouter as Router, Route, Switch, Link, useHistory, useParams } from 'react-router-dom';
import Modal from 'react-modal'




import { getJwtFromCookie } from './App';
import jwt_decode from 'jwt-decode'


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
            'Content-Type': 'application/json',
          },
        });

        if (!response.ok) {
          throw new Error('Failed to get lists');
        }

        const Lists = await response.json();
        setLsts(Lists);
        console.error(Lists);
      } catch (error) {
        console.error('Error getting all lists:', error);
      }
    };

    fetchData();

    if (isNewListAddedOrRemoved) {
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

    const [isNewCardAddedOrRemoved, setIsNewCardAddedOrRemoved] = useState(false);


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
                handleSaveList(newListName)
                setNewListName('');
                setIsAddingList(false)
            }
          }

         return (
            <div className="add-list-buttons">
              <button type="submit" onClick={() => addNewList()} style={{fontFamily:'shabnam', fontSize:'12px'}}>ذخیره</button>
                <button type="submit" onClick={() => {
                    setIsAddingList(false)
                }} style={{fontFamily:'shabnam', fontSize:'12px'}}>لغو</button>
            </div>
        )

    }


     
    const handleDeleteList = async (lst) => {
        
        try {
            const response = await fetch(`http://localhost:8080/api/boards/${boardId}/lists/${lst.id}`,{
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
        try {
            const response = await fetch(`http://localhost:8080/api/lists/${id}/cards`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              name: newCardName,
              username: user.name,
              user_id: user.id,
              user_email: user.email,
              owner_id: user.user_id,
            }),
          });

          if (!response.ok) {
            throw new Error('Failed to create a new card');
          }
      
          const newCard = await response.json();
      
          const updatedLists = lsts.map((list) => {
            console.error('list: ', list);
          
            // Check if list.cards is null or undefined, then initialize it as an empty array
            const updatedList = {
              ...list,
              cards: list.cards ? [...list.cards, ...newCard.cards] : [newCard],
            };
          
            if (list.id === id) {
              return updatedList;
            }
          
            return list;
          });
          
          setIsNewCardAddedOrRemoved(true); // Move this line after setLsts
          
          setLsts(updatedLists);

          window.location.reload();
        } catch (error) {
          console.error('Error creating a new card:', error);
        }

      };
        
        
        const addNewCard = () => {
            if (newCardName.trim() !== '') {
                handleSaveCard(newCardName)
                setNewCardName('');
                setIsAddingCard(false);
            }
        }


        return (
            <div className="add-card-buttons">
              <button type="submit" onClick={() => addNewCard()} style={{fontFamily:'shabnam', fontSize:'12px'}}>ذخیره</button>
                <button type="submit" onClick={() => {
                    setIsAddingCard(false)
                }} style={{fontFamily:'shabnam', fontSize:'12px'}}>لغو</button>
            </div>
        )

    }


    

      

    const handleDragEnd = (result) => {
      // Check if the drag operation was completed successfully
      if (!result.destination) {
        return;
      }
      const updatedLists = Array.from(lsts);

      if (result.type === 'CARD') {
      // Handle card drag and drop
      const sourceListId = result.source.droppableId;
      const sourceListIdNumber = parseInt(sourceListId.split('-')[1], 10);

      const movedCardId = parseInt(result.draggableId.split('-')[0], 10); // Extract cardId from the draggableId

      const destinationListId = result.destination.droppableId;
      const destinationListIdNumber = parseInt(destinationListId.split('-')[1], 10);

      // Clone the current state for modification
      const updatedLists = [...lsts];

      // Find the source list and destination list using their IDs
      const sourceListIndex = updatedLists.findIndex((list) => list.id === sourceListIdNumber);
      const destinationListIndex = updatedLists.findIndex((list) => list.id === destinationListIdNumber);

      if (sourceListIndex === -1 || destinationListIndex === -1) {
        console.error('Source or destination list not found. Lists:', updatedLists);
        return;
      }

      const sourceList = updatedLists[sourceListIndex];
      
      const destinationList = updatedLists[destinationListIndex];
      destinationList.cards = destinationList.cards || [];

      if (!destinationList || !destinationList.cards) {
        console.error('Destination list or its cards array is null or undefined');
        console.error('Destination list:', destinationList);
        return;
      }


      // Find the moved card in the source list
      const movedCardIndex = sourceList.cards.findIndex((card) => card && card.id === movedCardId);

      if (movedCardIndex === -1) {
        console.error('Moved card not found in source list. Lists:', updatedLists);
        console.error('Moved card id:', movedCardId);
        return;
      }

      // Remove the moved card from the source list and add it to the destination list
      const movedCard = sourceList.cards.splice(movedCardIndex, 1)[0];
      destinationList.cards.splice(result.destination.index, 0, movedCard);

      // Update the positions of cards in the source list
      sourceList.cards.forEach((card, index) => {
        card.position = index;
      });

      // Update the positions of cards in the destination list
      destinationList.cards.forEach((card, index) => {
        card.position = index;
      });

      // Update the state with the modified lists
      setLsts(updatedLists);

      // Now, make an API call to update the card order and positions on the server
      fetch(`http://localhost:8080/api/boards/${boardId}/lists/${sourceListIdNumber}/update-cards-order`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          listId: sourceListIdNumber,
          cardOrder: sourceList.cards.map((card) => card.id),
        }),
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error('Failed to update card order on the server');
          }
        })
        .catch((error) => {
          console.error('Error updating card order on the server:', error);
        })
        .finally(() => {
          // Now, make another API call to update the list of the moved card in the destination list
          fetch(`http://localhost:8080/api/boards/${boardId}/lists/${destinationListIdNumber}/card-to-list-order`, {
            method: 'PUT',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              sourceListId: sourceListIdNumber,
              destinationListId: destinationListIdNumber,
              cardId: movedCardId,
              cardName: movedCard.name,  // Include the name of the card
              newPosition: result.destination.index,  // Include the new position of the card
            }),
          })
            .then((response) => {
              if (!response.ok) {
                throw new Error('Failed to update card order on the server');
              }
            })
            .catch((error) => {
              console.error('Error updating card order on the server:', error);
            });
        });
      } else {
        // Handle list drag and drop
        const movedList = updatedLists.splice(result.source.index, 1)[0];
        updatedLists.splice(result.destination.index, 0, movedList);

        // Update the positions of lists
        updatedLists.forEach((list, index) => {
          list.position = index;
        });

        // Update the frontend immediately
        setLsts(updatedLists);

        // Now, send a fetch request to update the list positions on the backend
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
      }
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
                        <ShowCards list={lst} isNewCardAddedOrRemoved={isNewCardAddedOrRemoved} setIsNewCardAddedOrRemoved={setIsNewCardAddedOrRemoved}/>
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
                            direction:'rtl'
                          }}
                        />
  
                        {isAddingCard === lst.id && <AddCard id={lst.id} />}
  
                        <br />
                        <button onClick={() => handleDeleteList(lst)} className="remove-button">
                          حذف لیست
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
  
              {/* {provided.placeholder} */}
            </div>
          )}
        </Droppable>
      </DragDropContext>
    );


  };




// Inside the ShowCards component

const ShowCards = ({ list, isNewCardAddedOrRemoved, setIsNewCardAddedOrRemoved }) => {

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

    const formatCardDate = (dateString) => {
      const options = { month: 'long', day: 'numeric' };
      const formattedDate = new Date(dateString).toLocaleDateString('fa-IR', options);
      return formattedDate;
    };

    
    useEffect(() => {
      if (isNewCardAddedOrRemoved) {
        setCardList(list.cards);
        setIsNewCardAddedOrRemoved(false); // Reset the variable after updating the cards
      }
    }, [isNewCardAddedOrRemoved, list.cards, setIsNewCardAddedOrRemoved]);
    

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
        // <DragDropContext onDragEnd={onDragEnd}>
          <Droppable droppableId={`cards-${list.id}`} type="CARD">
            {(provided) => (
              <div ref={provided.innerRef}>
                {cardList && cardList.length > 0 ? (
                  cardList.map((card, index) => (
                    <Draggable
                      key={card.id.toString()}
                      draggableId={`${card.id}-${list.id}`}
                      index={index}
                    >
                      {(provided) => (
                        <div
                          ref={provided.innerRef}
                          {...provided.draggableProps}
                          {...provided.dragHandleProps}
                          className="card-item"
                        >
                          {card.label && (
                          <div
                            style={{
                              height: '2px',
                              width: '100px',
                              backgroundColor: card.label,
                              marginLeft: '10px',
                              marginLeft:'40px'
                            }}
                          ></div>
                        )}
                          {/* ... your card content ... */}
                          <div key={card.id} className="card-item">
                <div
                        style={{ textDecoration: 'none', color: 'black', cursor: 'pointer' }}
                        onClick={() => openModal(card)}>
                        <h4 style={{ fontFamily: 'Shabnam-Medium', fontWeight: 'normal' }}>{card.name}</h4>
                </div>
                     
                    <div className="icons-container" style={{display:'inline-flex' , direction: 'rtl'}}>
                        
                        <h6><img src={require('./icons/members.png')} alt="members" style={{width:'15px', height:'24x'}} />{card.members && card.members.length}</h6>
                        <img src={require('./icons/desc.png')} alt="desc" style={{width:'15px', marginRight:'20px', marginLeft:'20px'}} />
                        {card.checklists && card.checklists.length > 0 && (
                          <h6>
                            <img src={require('./icons/tasks.png')} alt="tasks" style={{width:'15px', height:'24x'}} />
                            {card.checklists.length}
                          </h6>
                        )}

                        {card.dates && card.dates.length > 0 && (
                            <div style={{fontSize:'13px', marginRight:'10px', marginTop:'3px'}}>
                              <h6>
                                {formatCardDate(card.dates[0])} - {formatCardDate(card.dates[1])} 
                              </h6>
                            </div>
                          )}
                        
                    </div>
                    <br />
                    <Modal isOpen={modalIsOpen} onRequestClose={closeModal} className="modal">
                    {selectedCard && (
                      <div className="modal-content">
                        {user ? (
                          selectedCard.members.some((member) => member.name === user.name) ? (
                            <>
                              {/* Display modal content with selectedCard data */}
                              <Card card={selectedCard} list={list} userIsMember={true}></Card>
                              {/* Rest of the modal content */}
                              <div className="center-button-container">
                                <button onClick={() => { window.location.href = `/board/${boardId}/lists` }}
                                style={{fontFamily:'shabnam'}}>
                                  بستن صفحه
                                  </button>
                              </div>
                            </>
                          ) : (
                            <>
                              <Card card={selectedCard} list={list} userIsMember={false}></Card>
                              <button  className="center-button-container" onClick={() => {
                                window.location.href = `/board/${boardId}/lists`
                              }} style={{fontFamily:'shabnam'}}>بستن صفحه</button>
                            </>
                          )
                        ) : (<div className="authentication-prompt">
                              <p>شما برای دیدن محتوای کارت باید وارد شوید</p>
                              <div className="buttons-container">
                                <button
                                  className="signup-button"
                                  onClick={() => {
                                    history.push('/signup');
                                  }}
                                  style={{fontFamily:'vazirmatn'}}
                                >
                                  ثبت نام
                                </button>
                                <button
                                  className="close-modal-button"
                                  onClick={() => {
                                    window.location.href = `/board/${boardId}/lists`;
                                  }}
                                >
                                  Close Modal
                                </button>
                              </div>
                            </div>
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
        // </DragDropContext>
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