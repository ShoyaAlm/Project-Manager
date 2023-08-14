import {useEffect, useState} from 'react'


import './signup.css'
export const HandleSigninLogin = () => {

    const [name, setName] = useState('')
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [users, setUsers] = useState([])
    
    
    const handleSubmit = (e) => {
        e.preventDefault()

        if (name && email && password) {
            const person = {id: new Date().getTime().toString(), name, email, password}
            setUsers(person)
            console.log('users: ',users);
        } else {
            console.log('enter value for all inputs!');
        }
        console.log('users: ',users);
    }
    
    return (
        <>

        <div class="signupFrm">
            <div class="wrapper">
            <form action="" class="form" onSubmit={(e) => handleSubmit(e)}>
            <h1 class="title">Sign up</h1>

            <div class="inputContainer">
                <input type="text" class="input" placeholder="a" 
                onChange={(e) => setEmail(e.target.value)}/>
                <label for="" class="label">Email </label>
            </div>

            <div class="inputContainer">
                <input type="text" class="input" placeholder="a" 
                onChange={(e) => setName(e.target.value)}/>
                <label for="" class="label">Name </label>
            </div>

            <div class="inputContainer">
                <input type="text" class="input" placeholder="a" 
                onChange={(e) => setPassword(e.target.value)}/>
                <label for="" class="label">Password </label>
            </div>


            <input type="submit" class="submitBtn" value="Sign up"/>
            </form>
            </div>
        </div>


        
        </>
    )

}


