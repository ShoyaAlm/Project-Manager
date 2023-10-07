import {useEffect, useState} from 'react'
import {Link, useParams} from 'react-router-dom'
// import { useReducer } from 'react';
import './card.css'
export const Card = ({card, list}) => {

    const newCard = card
    const newList = list
    
    const {name, description, members, checklists} = newCard
    const [cardName, setCardName] = useState(name)
    const [cardDescription, setCardDescription] = useState(description)
    const [cardMembers, setCardMembers] = useState(members)
    
    // Define state for managing description editing
    const [isEditingDescription, setIsEditingDescription] = useState(false);
    
    // Define state to store the temporary edited description
    const [editedDescription, setEditedDescription] = useState(cardDescription);
    
    
    const [isAddingItem, setIsAddingItem] = useState(false)
    
    const [isAddingMember, setIsAddingMember] = useState(false)
    
    const [isNewMemberAddedOrRemoved, setIsNewMemberAddedOrRemoved] =  useState(false)
    
    
    
    const [isNewChecklistAddedOrRemoved, setIsNewChecklistAddedOrRemoved] = useState(false)
    const [isAddingChecklist, setIsAddingChecklist] = useState(false)
    const [cardChecklists, setCardChecklists] = useState(checklists)
    const [checklistCardID, setChecklistCardID] = useState('')
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
                console.log('newchecklistname: ', newChecklistName);
                setNewChecklistName('');

                const newChecklist = await response.json();
                setCardChecklists([...cardChecklists, newChecklist]);
                
                setIsNewChecklistAddedOrRemoved(true)
                
            } catch (error) {
                console.log("Error creating the checklist");
            }
            
        }
        
        const addNewChecklist = () => {
            if (newChecklistName.trim() !== '') {
                handleSaveChecklist(newChecklistName)
                setIsAddingChecklist(false)
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
    

    const AddMember = ({ card, list }) => {

        const [newMemberName, setNewMemberName] = useState('')

        useEffect(() => {
            
            const fetchMembers = async () => {
                try {
                    const response = await fetch(`http://localhost:8080/api/lists/${list.id}/cards/${card.id}/members`,{
                        method: 'GET',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                    });
                    if (!response.ok){
                        throw new Error('Error deleting the member')
                    }
            
                    const allMembers = await response.json();
                    setCardMembers(allMembers)
                    
                } catch (error) {
                    console.log('got error : ', error);
                }

            }


            if(isNewMemberAddedOrRemoved){
                fetchMembers();
                setIsNewMemberAddedOrRemoved(false)
            }

        },[isNewMemberAddedOrRemoved])
        
        
        const handleNewMember = async (newMemberName) => {
            try {
                const response = await fetch(`http://localhost:8080/api/lists/${list.id}/cards/${card.id}/members`, {
                    method: 'POST',
                    headers: {
                        'Content-type': 'application/json'
                    },
                    body: JSON.stringify({name: newMemberName}),
                });
                if(!response.ok){
                    throw new Error("Failed to create new member")
                }

                const newMember = await response.json()
                console.log('new member: ', newMember);
                // setCardMembers([...cardMembers, newMember])
                setIsNewMemberAddedOrRemoved(true)

            } catch (error) {
                console.log('Error creating the member');
            }
        }
        
        
        const addNewMember = () => {
            if(newMemberName !== ''){
                handleNewMember(newMemberName)
                setNewMemberName('')
                setIsAddingMember(false)
            }

        }

        return (
            <div className='add-member'>
                <button onClick={() => addNewMember()} style={{width:'50px', height:'40px'}}>ذخیره</button>
                <input
                    type="text"
                    value={newMemberName}
                    onChange={(e) => setNewMemberName(e.target.value)}
                    placeholder="Enter member name"
                    style={{width:'150px', height: '40px', marginLeft:'auto'}}/>
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
                const response = await fetch(`http://localhost:8080/api/lists/${newList.id}/cards/${newCard.id}`,{
                    method: 'DELETE',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                });
                if (!response.ok){
                    throw new Error('Error deleting the card')
                }
            } catch (error) {
                console.log(error);
            }

    }



    const AddItem = ({ checklist }) => {
        const [newItemName, setNewItemName] = useState('');
        const [checklistItems, setChecklistItems] = useState(checklist.items)
        
        const handleSaveItem = async () => {
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

                const newItem = await response.json();
                setChecklistItems([...checklistItems, newItem])
                                

            } catch (error) {
            console.log("Error creating the item");
                }

        }

        const addNewItem = () => {
            if (newItemName.trim() !== '') {
                handleSaveItem(newItemName)
                setNewItemName('');
                setIsAddingItem(false)
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


                        {isAddingMember === newCard.id && <AddMember card={card} list={list}/>}
                        
                        <button className='add-member-button' onClick={() => {
                            if(isAddingMember === ''){
                                setIsAddingMember(newCard.id)
                            } else {
                                setIsAddingMember('')
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
                        />
                        {isEditingDescription && (
                            <div className="description-buttons">
                                <button type="submit" onClick={() => {
                                    setCardDescription(editedDescription)
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
                                checklist.items.map((item, itemIndex) => (
                                    
                                        <div className="checklist-item" key={itemIndex}>
                                            <button className='remove-item-button' onClick={() => removeItem(checklist.id, item.id)}>حذف</button>
                                            <label htmlFor="item">{item.name}</label>
                                            <input type="checkbox" id="item"/>
                                        </div>
   
                                ))
                                
                                ) : (
                                    <span>بدون آیتم</span>
                                )}


                                {isAddingItem === checklist.id && <AddItem checklist={checklist}/> }


                                <button type='button' className='add-item-button' onClick={() => {
                                    if(isAddingItem === ''){
                                        setIsAddingItem(checklist.id)
                                    } else {
                                        setIsAddingItem('')
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
                            
                            {checklistCardID === newCard.id && <AddChecklist card={newCard} list={newList}/>}
                        
                            <button type='button' className='add-checklist-button' onClick={() => {
                                        if(checklistCardID === ''){
                                            setIsAddingChecklist(true)
                                            setChecklistCardID(card.id)
                                        } else {
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
                                         <button className='remove-member-button' onClick={() => removeMember(member.id)}>X</button>
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

                        <div className='dropdown'>
                        <button className='dropbtn'>تاریخ</button>
                            <div className='dropdown-content'>
                                <a href="#">{newCard.dates[0]} : شروع</a>
                                <a href="#">{newCard.dates[1]} : پایان</a>
                            </div>
                        </div>
 
                    </div>


                </div>

            </div>
        
        </div>  




);
    
}    

