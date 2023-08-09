import lists from './lists'

import React, { useEffect, useContext } from 'react'
import {useState} from 'react'

import {Link, useParams} from 'react-router-dom'
import {Router} from 'react-router-dom'

// import { List } from './List'

const ListContext = React.createContext()
// has 2 components : provider - consumer

export const AllLists = () => {
    const [lsts, setLsts] = useState(lists)
    const removeList = (id) => {  // i want to pass this func into the SinglePerson, and for that i
        setLsts((lsts) => {            // create a context using createContext
            return lsts.filter((lsts) => lsts.id !== id)
        });
    };
    
    return (        //and in here, i'll pass it as a provider
        <ListContext.Provider value={{removeList, lsts}}>
            <List />
        </ListContext.Provider>
    )

}


const List = () => {
    const {removeList, lsts} = useContext(ListContext) // I'll be receiving it here
    return (
        <div>
            {lsts.map((lst) => {
                return (
                    <div key={lst.id} style={{display:'inline-block', width:'400px', 
                    height:'200px', marginLeft: '200px', textAlign:'center' }}>
                       <h2>{lst.id}</h2>
                       <Link to={`/lists/${lst.id}`}>
                       <h1>Name: {lst.name}</h1></Link>
                       <ShowCards list={lst}/>
                        <button onClick={() => removeList(lst.id)}>remove</button>
                    </div>
                )
            })}
        </div>
    )
}



const ShowCards = ({list}) => {

    const cards = list.cards
    console.log(cards);
    return (
        <>

            <h1 style={{textAlign:'center'}}>Cards</h1>
            <hr style={{marginLeft:'600px', marginRight:'600px'}}/> <br />
                <>
                    {cards.map((card) => {
                        return (
                            <div key={card.id} style={{display:'inline-block', width:'500px',
                             height:'200px', marginTop:'80px'}}>  
                                <h2>card id : {card.id}</h2>
                                <h2>name: {card.name}</h2>
                                <h3>status: {card.status}</h3> <br />
                                        
                                <Link to={`/lists/${list.id}/card/${card.id}`}> edit
                                </Link>
                                        
                                </div>
                                )
                            })}
                </>
        </>
    )
}



// const List = () => {

//     const {removeList, lsts} = useContext(ListContext)
    
//     const [listName, setListName] = useState('default')
//     const [listId, setListId] = useState(1)
//     const {id} = useParams();

//     const newList = lists.find((lst) => lst.id === parseInt(id))
//     useEffect(() => {
//         if(newList){
//             setListName(newList.name)
//             setListId(newList.id)
//         }
//     },[]);

//     return (
//         <div style={{textAlign:'center'}}>
//             <h1>{listName}</h1>
//             <h2>{listId}</h2> <br /> <br />

//             <ShowCards list={newList}/> <br />
//             <Link to="/">Back to Homepage</Link>
//         </div>
//     )
// }