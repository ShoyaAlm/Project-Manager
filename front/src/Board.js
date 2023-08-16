import lists from './lists'

import React, {useContext } from 'react'
import {useState} from 'react'

// import {desc} from 'icons/desc.png'

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

    const [newList, setNewList] = useState({})
    const [newListName, setNewListName] = useState('')
    
    
    const addNewList = () => {  // new code
    if (newListName.trim() !== '') {
    const newListToAdd = {
    id: Date.now(),
    name: newListName,
    cards: [
        {
            id: Date.now(),
            name: 'new card',
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
                                <button type="submit" onClick={() => handleSaveCard(lst)}>ذخیره</button>
                                <button type='submit' onClick={() => {
                                    setNewCardName('')
                                    setIsAddingCard(false)
                                }}>کنسل</button>
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
                        <button type="submit" onClick={() => addNewList()}>ذخیره</button>
                        <button type="submit" onClick={() => {
                            setNewListName('');
                            setIsAddingList(false);
                        }}>کنسل</button>
                    </div>
                )}
            </div>


        </div>
    );
};




// Inside the ShowCards component

const ShowCards = ({ list }) => {


    return (
        <div className="showcards-container">
            <hr />
            {list.cards.map((card) => (
                <div key={card.id} className="card-item">
                    <a href={`/lists/${list.id}/cards/${card.id}` } 
                    style={{textDecoration:'none', color:'black'}}>
                    <h4 style={{fontFamily:'sans-serif'}}>{card.name}</h4>
                    </a> 
                    <div className="icons-container" style={{display:'inline-block'}}>
                        
                    <h6><img src={require('./icons/members.png')} alt="members" style={{width:'16px', height:'25x'}} />{card.members && card.members.length}</h6>

                        <img src={require('./icons/desc.png')} alt="desc" style={{width:'16px', height:'25x'}} />
                        <h6><img src={require('./icons/tasks.png')} alt="tasks" style={{width:'16px', height:'25x'}} />1/2</h6>                

                    </div>
                    <br />
                </div>

            ))}

        </div>
    );
};
