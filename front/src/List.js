import lists from './lists.js'

import {useState, useEffect} from 'react'
import {Link, useParams} from 'react-router-dom'
import {Router, Route} from 'react-router-dom'


export const List = () => {
    
    const [listName, setListName] = useState('default')
    const [listLevel, setListLevel] = useState(1)
    const {id} = useParams();

    const newList = lists.find((lst) => lst.id === parseInt(id))
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
            <h1 style={{textAlign:'center'}}>Cards</h1>
            <hr style={{marginLeft:'600px', marginRight:'600px'}}/> <br />
                <>
                    {list.cards.map((card) => {
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


const CreateList = () => {
    const [name, setName] = useState('')
    const [card, setcard] = useState([])
    const [level, setLevel] = useState()


    const handleSubmit = (e) => {
        if(name && card && level) {
            const createdList = {id: new Date().getTime().toString(), name, card, level}
        }
    }


    return (
        <div style={{textAlign:'center'}}>

            <form >
                
                <label htmlFor="name">Name: <input type="text"/></label>

                <label htmlFor="card">card: <input type="text"/></label>

                <label htmlFor="level">Level: <input type="text"/></label>

            </form>

        </div>
    )
}