import {useState} from 'react'
import {Link, useParams} from 'react-router-dom'
import lists from './lists'
// import { useReducer } from 'react';
import './card.css'
export const Card = ({card, list}) => {

    const cardId = card.id
    const listId = list.id

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

    const AddItem = ({ checklist }) => {
        const [newItemName, setNewItemName] = useState('');
    
        const handleSaveItem = () => {
            if (newItemName.trim() !== '') {
                const updatedChecklist = {
                    ...checklist,
                    items: [
                        ...checklist.items,
                        {
                            id: Date.now(),
                            name: newItemName,
                            dueDate: '',
                            assignedTo: ['susan', 'jake'],
                        },
                    ],
                };
    
                const updatedCardChecklists = cardChecklists.map(c =>
                    c.id === checklist.id ? updatedChecklist : c
                );
    
                setCardChecklists(updatedCardChecklists);
                setNewItemName('');
            }
        };
    
        return (
            <div>
                <button onClick={handleSaveItem} style={{width:'auto', height:'auto'}}>ذخیره</button>
                <input
                    type="text"
                    value={newItemName}
                    onChange={(e) => setNewItemName(e.target.value)}
                    placeholder="Enter item name"
                    style={{width:'200px', height: '40px'}}/>
            </div>
        );
    };
    

    const removeChecklist = (id) => {
        setCardChecklists((cardChecklists) => {
            return cardChecklists.filter((cardChecklists) => cardChecklists.id !== id)
        })
    }

    return (
        
        <div>
           
            <div className='card-container'>
                
                <div className='card-details' style={{ textAlign: 'left', width: '800px', 
                    height: 'auto', marginLeft: '320px' }}>
                    <h2 className='card-title' style={{textAlign:'right'}}><img src={require('./icons/list.png')} alt="" style={{width:'50px', height:'50px', marginRight:'-20%' , marginBottom:'-30px' , position:'relative'}}/>
                    {cardName}</h2>
                    
                    <h3 className='list-name' style={{textAlign:'right', marginRight:'40px', marginTop:'10px'}}>در لیست <Link to={`/list/${newList.id}`} style={{
                        textDecoration:'none',
                        color:'blue'
                    }}>{newList.name}</Link></h3>



                    <div className='card-members' style={{marginRight:'30px'}}>
                        <img src={require('./icons/members.png')} alt="" style={{width:'24px', height:'24px', marginLeft:'780px', marginBottom:'-36px', position:'relative'}}/>
                        <h3 style={{textAlign:'right', marginRight:'6px'}}>اعضا</h3>
                        <h4 style={{textAlign:'right'}}>
                            {cardMembers.map((member, index) => (
                                <span key={index} style={{ marginRight: '8px', textAlign:'right' }}>{member} </span>
                            ))}
                        </h4>
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
                                    console.log('cancel');
                                    setEditedDescription(cardDescription);
                                    setIsEditingDescription(false)
                                    }}>لغو</button>
                            </div>
                        )}
                    </div>
                            





                    <div className='showcase-checklists' style={{marginRight:'auto'}}>
                        {cardChecklists.map((checklist, index) => (
                            <div className='checklist' key={index}>
                                <h2 className='checklist-title'><img src={require('./icons/checklist.png')} alt="" style={{width:'25px', height:'25px', marginBottom:'-5px', marginLeft:'-30px', marginRight:'10px'}}/>
                                {checklist.name}</h2>
                                {checklist.items.map((item, itemIndex) => (
                                    
                                        <div class="checklist-item" key={itemIndex}>
                                            <label for="item">{item.name}</label>
                                            <input type="checkbox" id="item"/>
                                        </div>
   
                                ))}


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
                        ))}
                    </div>
                







                    <div className='add-to-card' style={{width:'200px', height:'auto'}}>
                        
                        <div className='dropdown'>
                        <button className='dropbtn'>اعضا</button>
                            <div className='dropdown-content'>
                                {cardMembers.map((member) => {
                                    return <a href="#">{member}</a>
                                })}
                            </div>
                        </div>

                        <div className='dropdown'>
                        <button className='dropbtn'>چکلیست</button>
                            <div className='dropdown-content'>
                                {cardChecklists.map((checklist) => {
                                    return <a href="#">{checklist.name}</a>
                                })}
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

