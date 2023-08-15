import {useState} from 'react'
import {Link, useParams} from 'react-router-dom'
import lists from './lists'
// import { useReducer } from 'react';
import './card.css'
export const Card = () => {

    const {cardId, listId} = useParams()

    const newList = lists.find((project) => project.id === parseInt(listId))
    
    const newCard = newList.cards.find((card) => card.id === parseInt(cardId))
    
    
    const {name, description, members, checklists} = newCard
    const [cardName, setCardName] = useState(name)
    const [cardDescription, setCardDescription] = useState(description)
    const [cardMembers, setCardMembers] = useState(members)
    const [cardChecklists, setCardChecklists] = useState(checklists)


    // Define state for managing description editing
    const [isEditingDescription, setIsEditingDescription] = useState(false);

    // Define state to store the temporary edited description
    const [editedDescription, setEditedDescription] = useState(cardDescription);


    const [isAddingItem, setIsAddingItem] = useState(false)
    const [newItemName, setNewItemName] = useState('')



    const AddItem = (checklist) => {
        

        const handleSaveItem = () => {
            if(newItemName.trim() !== ''){
                const updatedChecklist = {
                    ...checklist,
                    items: [
                        ...checklist.items,
                        {
                            id: Date.now(),
                            name:newItemName,
                            dueDate:'24th september',
                            assignedTo:['josh', 'peter']
                        }
                    ]
                }
        
                const updatedCardChecklists = [...cardChecklists]
        
                updatedCardChecklists[checklist.id] = updatedChecklist;
        
                setCardChecklists(updatedCardChecklists);
    
                setNewItemName('');
            }

        };
        
        return (
            <div>
                <input
                    type="text"
                    value={newItemName}
                    onChange={(e) => setNewItemName(e.target.value)}
                    placeholder="Enter item name"
                />
                <button onClick={() => handleSaveItem()}>Save</button>
            </div>
        );
    }

    const removeChecklist = (id) => {
        setCardChecklists((cardChecklists) => {
            return cardChecklists.filter((cardChecklists) => cardChecklists.id !== id)
        })
    }

    return (
        
        <div>
           
            <div className='card-container'>
                
                <div className='card-details' style={{ textAlign: 'left', width: '800px', 
                    height: '800px', marginLeft: '320px' }}>   
                    <h2 className='card-title'>{cardName}</h2>
                    
                    <h3 className='list-name'>in list @<Link to={`/list/${newList.id}`} style={{
                        textDecoration:'none',
                        color:'blue'
                    }}>{newList.name}</Link></h3>
                    
                    <div className='card-members'>
                        <h3>Members</h3>
                        <h4>
                            {cardMembers.map((member, index) => (
                                <span key={index} style={{ marginRight: '5px' }}>{member} </span>
                            ))}
                        </h4>
                    </div>
                    <h2 className='section-title'>Description</h2>
                    <div className="description-input">
                        <input
                            type="text"
                            className={isEditingDescription ? 'card-description-active' : 'card-description'}
                            value={editedDescription}
                            onFocus={() => setIsEditingDescription(true)}
                            // onBlur={() => {setIsEditingDescription(false)}}
                            // onBlur={cancelDescription}
                            onChange={(e) => setEditedDescription(e.target.value)}
                        />
                        {isEditingDescription && (
                            <div className="description-buttons">
                                <button type="submit" onClick={() => {
                                    setCardDescription(editedDescription)
                                    setIsEditingDescription(false);
                                }}>Save</button>
                                <button type="submit" onClick={() => {
                                    console.log('cancel');
                                    setEditedDescription(cardDescription);
                                    setIsEditingDescription(false)
                                    }}>Cancel</button>
                            </div>
                        )}
                    </div>
                            
                    
                    <div className='showcase-checklists'>
                        {cardChecklists.map((checklist, index) => (
                            <div className='checklist' key={index}>
                                <h2 className='checklist-title'>{checklist.name}</h2>
                                {checklist.items.map((item, itemIndex) => (
                                    
                                        <div class="checklist-item" key={itemIndex}>
                                            <input type="checkbox" id="item"/>
                                            <label for="item">{item.name}</label>
                                        </div>


                                    
                                ))}

                                {isAddingItem === checklist.id && <AddItem checklist={checklist}/> }



                                <button type='button' className='add-item-button' onClick={() => setIsAddingItem(checklist.id)}>Add an item</button>
                                








                                <br />
                                <button type='submit' className='remove-checklist-button' onClick={() => removeChecklist(checklist.id)}>remove</button>
                            </div>
                        ))}
                    </div>
                

                    <div className='add-to-card'>
                        
                        <div className='dropdown'>
                        <button className='dropbtn'>Members</button>
                            <div className='dropdown-content'>
                                {cardMembers.map((member) => {
                                    return <a href="#">{member}</a>
                                })}
                            </div>
                        </div>

                        <div className='dropdown'>
                        <button className='dropbtn'>Checklist</button>
                            <div className='dropdown-content'>
                                {cardChecklists.map((checklist) => {
                                    return <a href="#">{checklist.name}</a>
                                })}
                            </div>
                        </div>

                        <div className='dropdown'>
                        <button className='dropbtn'>Lables</button>
                            <div className='dropdown-content'>
                            </div>
                        </div>

                        <div className='dropdown'>
                        <button className='dropbtn'>Dates</button>
                            <div className='dropdown-content'>
                                <a href="#">starting: {newCard.dates[0]}</a>
                                <a href="#">ending: {newCard.dates[1]}</a>
                            </div>
                        </div>
 
                    </div>


                </div>

            </div>
        
        </div>  




);
    
}    

