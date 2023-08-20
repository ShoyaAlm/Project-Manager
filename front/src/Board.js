import lists from './lists'

import React, {useContext } from 'react'
import {useState} from 'react'

import { BrowserRouter as Router, Route, Switch, Link, useHistory } from 'react-router-dom';
import Modal from 'react-modal'

import { Card } from './Card';

import './board.css'
const ListContext = React.createContext()

export const AllLists = () => {
    const [lsts, setLsts] = useState(lists)
    const removeList = (id) => {  // i want to pass this func into the SinglePerson, and for that i
        setLsts((lsts) => {            // create a context using createContext
            return lsts.filter((lsts) => lsts.id !== id)
        });
    };
    
    return (        //and in here, i'll pass it as a provider
        <ListContext.Provider value={{removeList, lsts, setLsts}}>
            <List />
        </ListContext.Provider>
    )

}


const List = () => {
    const { removeList, lsts, setLsts } = useContext(ListContext);
    
    const [isAddingCard, setIsAddingCard] = useState(false);
    
    const [newCardName, setNewCardName] = useState('random')
    
    
    const [isAddingList, setIsAddingList] = useState(false)
    
    const handleSaveCard = (lst) => {
        if (newCardName.trim() !== '') {
            const updatedList = lst.cards.push({ id: Date.now(), name:newCardName, 
                members:[], dates:[Date.now(), Date.now()],
            description:'',checklists:[]})
            console.log('updated List', updatedList)
            setNewCardName('')
            }
    };

    // const [newList, setNewList] = useState({})
    const [newListName, setNewListName] = useState('')
    
    
    const addNewList = () => {  // new code
    if (newListName.trim() !== '') {
    const newListToAdd = {
    id: Date.now(),
    name: newListName,
    cards: [
        {
            id: Date.now(),
            name: 'کارد جدید',
            members: ['alex', 'josh', 'lucas', 'peter'],
            dates: ['24th august', '21st september'],
            description: 'default description',
            checklists: [
                {
                    id: 1,
                    name: 'checklist 1',
                    items: [
                        {
                            id: 1,
                            name: 'item 1',
                            dueDate: '24th september',
                            assignedTo: ['josh', 'peter'],
                        },
                        {
                            id: 2,
                            name: 'item 2',
                            dueDate: '30th september',
                            assignedTo: ['alex', 'lucas'],
                        },
                            ],
                        },
                    ],
                },
            ],
        };

        setLsts(prevLsts => [...prevLsts, newListToAdd]);
        setNewListName('');
        setIsAddingList(false);
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
                                <button type="submit" onClick={() => handleSaveCard(lst)}>ذخیره</button>
                            </div>
                        )}
                    

                    <br />
                    <button onClick={() => removeList(lst.id)} className="remove-button">
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



      return (

            <div>
                {console.log('hello')}
            </div>
      )




    };
  
    const closeModal = () => {
      setSelectedCard(null);
      setModalIsOpen(false);
      history.push('/'); // Redirect to the main page when closing modal
    };


    return (
        <div className="showcards-container">
            <hr />
            {list.cards.map((card) => (
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
                    <Card card={card} list={list}></Card>
                    {/* Rest of the modal content */}
                    <button onClick={closeModal}>Close Modal</button>
                </div>
                )}
            </Modal>
                </div>

            ))}


        </div>
    );
};

Modal.setAppElement("#root")