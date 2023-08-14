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

    // Function to save the edited description
    const saveDescription = (e) => {
        e.preventDefault()
        console.log('save');
        setIsEditingDescription(false);
        // setEditedDescription(cardDescription);
    };
    
    // Function to cancel editing and revert to original description
    const cancelDescription = () => {
        console.log('cancel');
        setEditedDescription(cardDescription);
        setIsEditingDescription(false);
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
                                <button type="submit" onSubmit={() => {
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
                                <button type='submit' className='remove-checklist-button' onClick={() => removeChecklist(checklist.id)}>remove</button>
                            </div>
                        ))}
                    </div>
                

                    <div className='add-to-card'>
                        <h5 className='members'>Members</h5> <h5 className='labels'>Labels</h5>
                        <h5 className='atc-checklist'>Checklist</h5> <h5 className='dates'>Dates</h5>
                        <h5 className='attachment'>Attachment</h5> 
                    </div>


                </div>

            </div>
        
        </div>  




);
    
}    





    // const reducer = (state, action) => {
    
    //     console.log(state);
    //     if (action.type === 'edit_card'){
    //         const newcard = [...state.card, action.payload]
    //         return {
    //             ...state,
    //             card: newcard,
    //             isModalOpen: true,
    //             modalContent: 'card added'
    //         }
    //     }
    //     if(action.type === 'no_value'){
    //         return {
    //             ...state,
    //             isModalOpen: true,
    //             modalContent: 'please enter value'
    //         }
    //     }    
    
    //     throw new Error('No matching action type')
    // };
    
    
    // const defaultState = {
    //     card: [],
    //     name: '',
    //     description:'',
    //     isModalOpen: true,
    //     modalContent: 'Edit Card'
    // };
    
    
    
    // const Modal = ({modalContent, closeModal}) => {
    //     useEffect(() => {
    //         setTimeout(() => {
    //             closeModal();
    //         }, 3000);
    //     });
    //     return (
    //         <div>
    //             <p>{modalContent}</p>
    //         </div>
    //         )
    // }
    
    
    
    
    // const index = () => {
        
    //     const [name, setName] = useState('')
    //     const [description, setDescription] = useState('')
    //     const [state, dispatch] = useReducer(reducer, defaultState)
    
    //     const handleSubmit = (e) => {
    //         e.preventDefault();
    //         if (name && description) {
    //             const newCard = { id: new Date().getTime().toString(), name, description}
    //             dispatch({type: 'edit_card', payload: newCard })
    //             setName(name)
    //             setDescription(description)
    //         } else {
    //             dispatch({type: 'no_value'})
    //         }
    //     };
    
    //     return (
    //         <>
    //         {state.isModalOpen && (
    //             <Modal closeModal={closeModal} modalContent={state.modalContent} />
    //           )}
    //           <form onSubmit={handleSubmit} className='form'>
    //             <div>
    //               <input
    //                 type='text'
    //                 value={name}
    //                 onChange={(e) => setName(e.target.value)}
    //               />
    //             </div>
    //             <button type='submit'>add</button>
    //           </form>
    //             {state.users.map((person) => {
    //                 return (
    //                     <div key={person.id} className="card">
    //                         <h4>{person.name}</h4>
    //                     </div>
    //                 )
    //             })}
    //         </>
    //     )
    
    // }

