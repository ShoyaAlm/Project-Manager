import lists from './lists'

import React, {useContext } from 'react'
import {useState} from 'react'


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

// if (newCardName.trim() !== '') {
//     const updatedCards = lsts.cards.concat({ id: Date.now(), name: value });
//     // Update the cards array in the list object
//     const updatedList = { ...lsts, cards: updatedCards };
//     // Find the list index and update the entire list array
//     // const updatedLists = lsts.map((lst) => (lst.id === lsts.id ? updatedList : lst));
//     setLsts(prevLists => prevLists.map(lst => lst.id === updatedList.id ? updatedList : lst));
//     setNewCardName('');
// }
// setIsAddingCard(null);


// const handleCancelAddCard = () => {
//     setNewCardName('');
//     setIsAddingCard(null);
// };

const List = () => {
    const { removeList, lsts, setLsts } = useContext(ListContext);

    const [isAddingCard, setIsAddingCard] = useState(false);

    const [newCardName, setNewCardName] = useState('random')
    

    const handleSaveCard = (value) => {
        console.log('value: ', value);
    };


    const handleCancelCard = () => {
        console.log('handleCancel CardName: ', newCardName);
    }




    return (
        <div className="list-container">
            {lsts.map((lst) => (
                <div key={lst.id} className="list">
                <ShowCards list={lst} />
                    <input
                        type="text"
                        placeholder="+ Add an item"
                        onFocus={() => setIsAddingCard(lst.id)} // Pass the list ID to setIsAddingCard
                        onBlur={() => setIsAddingCard(null)}
                        className={isAddingCard === lst.id ? 'add-card-active' : 'add-card'}
                        // value={newCardName}        
                        // onChange={(e) => console.log(e.target.value)}
                        style={{margin: '10px', padding: '10px', 
                        width: '200px', height: 'auto', 
                        border: '2px solid #ccc', borderRadius: '20px'}}/>
                    
                    
                    {isAddingCard === lst.id && (
                            <div className="description-buttons">
                                <button onClick={() => handleSaveCard()}>Save</button>
                                <button onClick={() => console.log('cancel console log')}>Cancel</button>
                            </div>
                        )}
                    

                    <br />
                    <button onClick={() => removeList(lst.id)} className="remove-button">
                        Remove
                    </button>
                </div>
            ))}
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
                    <div className="icons-container" style={{display:'inline'}}>
                        
                        <img src="icons/members.png" alt="members" style={{width:'22px', height:'33x'}} />
                            <h6>{card.members.length}</h6>
                    

                        {/* <img src="https://thenounproject.com/api/private/icons/138377/edit/?backgroundShape=SQUARE&backgroundShapeColor=%23000000&backgroundShapeOpacity=0&exportSize=752&flipX=false&flipY=false&foregroundColor=%23000000&foregroundOpacity=1&imageFormat=png&rotation=0" alt="desc" style={{width:'22px', height:'33x'}}/> */}

                        {/* <img src="https://thenounproject.com/api/private/icons/5490508/edit/?backgroundShape=SQUARE&backgroundShapeColor=%23000000&backgroundShapeOpacity=0&exportSize=752&flipX=false&flipY=false&foregroundColor=%23000000&foregroundOpacity=1&imageFormat=png&rotation=0" alt="tasks" style={{width:'22px', height:'33x'}}/>                         */}
                            <h6>1/2</h6>
                    </div>
                    <br />
                </div>

            ))}
            
        </div>
    );
};
