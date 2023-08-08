import {useState} from 'react'

export const HandleSigninLogin = () => {

    const [firstName, setFirstName] = useState('')
    const [email, setEmail] = useState('')
    const [age, setAge] = useState('')
    const [users, setusers] = useState([])
    
    
    const handleSubmit = (e) => {
        e.preventDefault()
        if (firstName && email && age) {
            const person = {id: new Date().getTime().toString(), firstName, email, age}
            setusers((users) => {
                return [...users, person]
            })
            console.log('person: ',person);
        }
        console.log('users: ',users);
    }
    
    return (
        <>
            <div>
                <form onSubmit={handleSubmit}>
    
                    <label htmlFor="firstName">Name: </label>
                    <input type="text" name='firstName' 
                    value={firstName} onChange={(e) => setFirstName(e.target.value)}/> <br/>
    
                    <label htmlFor="email">Email: </label>
                    <input type="email" name='email' 
                    value={email} onChange={(e) => setEmail(e.target.value)}/> <br/>
    
                    <label htmlFor="age">Age: </label>
                    <input type="age" name='age' 
                    value={age} onChange={(e) => setAge(e.target.value)}/> <br/>
    
                <button type="submit">Submit</button>
                </form>
    
    
            </div>
        </>
    )

}
