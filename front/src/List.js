import lists from './lists.js'

import {useState, useEffect} from 'react'
import {Link, useParams} from 'react-router-dom'
import {Router, Route} from 'react-router-dom'
import './list.css'

export const List = () => {
    
    const [listName, setListName] = useState('default')
    const [listLevel, setListLevel] = useState(1)
    const {id} = useParams();

    const newList = lists.find((list) => list.id === parseInt(id))
    useEffect(() => {
        if(newList){
            setListName(newList.name)
            setListLevel(newList.level)
        }
    },[]);

    return (
        <div style={{textAlign:'center'}}>
            <h1>{listName}</h1>
            <h2>{listLevel}</h2> <br /> <br />

            <ShowCards list={newList}/> <br />
            <Link to="/">Back to Homepage</Link>
        </div>
    )
}


const ShowCards = ({list}) => {

    return (
        <>
        <h1 style={{ textAlign: 'center' }}>Cards</h1>
        <hr style={{ marginLeft: 'auto', marginRight: 'auto', width: '50%' }} /> <br />
        <div style={{ display: 'flex', flexWrap: 'wrap', justifyContent: 'center' }}>
            {list.cards.map((card) => (
                <div key={card.id} className="card">
                    <h2>Card {card.id}</h2>
                    <h3>Name: {card.name}</h3>
                    <Link to={`/lists/${list.id}/cards/${card.id}`} className="edit-link">
                        Edit
                    </Link>
                </div>
            ))}
        </div>
    </>
    
    )
}


// const CreateList = () => {
//     const [name, setName] = useState('')
//     const [card, setcard] = useState([])
//     const [level, setLevel] = useState()


//     const handleSubmit = (e) => {
//         if(name && card && level) {
//             const createdList = {id: new Date().getTime().toString(), name, card, level}
//         }
//     }


//     return (
//         <div style={{textAlign:'center'}}>

//             <form >
                
//                 <label htmlFor="name">Name: <input type="text"/></label>

//                 <label htmlFor="card">card: <input type="text"/></label>

//                 <label htmlFor="level">Level: <input type="text"/></label>

//             </form>

//         </div>
//     )
// }