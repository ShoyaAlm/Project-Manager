// import projects from './projects.js'
import './index.css'

import React, { useEffect, useContext } from 'react'
// import {useState} from 'react'


import { BrowserRouter as Router, Route, Switch} from 'react-router-dom'
import {Link, useParams} from 'react-router-dom'


import { AllLists } from './Board'
import { HandleSigninLogin } from './signup-login.js'
import { List } from './List'

const App = () => {

    return (
        <div>

            <Router>

           <Link to="/"><h1 style={{textAlign:'center'}}>Welcome to Board</h1></Link> <br />
            <Link to="/signup"><h3 style={{textAlign: 'right'}}>Signin/Login</h3></Link> <hr /> <br />
            
                <Switch>
                    <Route exact path='/'> <AllLists/> </Route>
                    
                    <Route exact path='/signup'> <HandleSigninLogin/> </Route>
                    
                    <Route path='/lists/:id' > <List/> </Route>
                    
                    {/* <Route path='/projects/:projectId/task/:taskId'> <EditTask/> </Route> */}
                    

                </Switch>

        

            </Router>


        </div>
    )


}



// ############ TEMPLATE ###########
// #################################
// #################################
// #################################

// we have a board, and in that board, we have some lists.
// those lists have names(like todo, doing & done)
// those lists consist of some actions(or cards)
// each card has members, category, startdate, due date, checklists
// each of those checklists have a name, and a list of actual todo's
//projects -> lists, tasks -> cards 

// showcase lists in the main page





// const Error = () => {
//     return (
//         <>
//             <h1>Error Page</h1>
//         </>
//     )
// }


// const AllProjects = () => {
//     const [prs, setPrs] = useState(projects)
//     return (
//         <div>
//             <h1>Projects</h1>
//             {prs.map((pr) => {
//                 return (
//                     <div key={pr.id}>
//                         <h2>Name: {pr.name}</h2>
//                         <h2>Level: {pr.level}</h2>
//                         <Link to={`/project/:${pr.id}`}>Learn more</Link>
//                     </div>
//                 )
//             })}
//         </div>
//     )
// }

// const Project = () => {
    
//     const [projectName, setProjectName] = useState('adsla')
//     const [projectLevel, setProjectLevel] = useState('')
//     const {id} = useParams();
//     console.log(id);
//     const newProject = projects.find((project) => project.id === parseInt(id))
//     console.log('lol : ',newProject);
//     // useEffect(() => {
//     //     setProjectName(newProject.name)
//     //     setProjectLevel(newProject.level)
//     // },[]);

//     return (
//         <div>
//             <h1>{projectName}</h1>
//             <h2>{projectLevel}</h2>
//             <Link to="/">Back to Homepage</Link>
//         </div>
//     )
// }


// const About = () => {
//     return (
//         <>
//             <h1>About Page</h1>
//         </>

//     )
// }


// const Navbar = () => {
//     return (
//         <nav>
//             <ul>
//                 <li>
//                     <Link to="/">Home</Link>
//                 </li>
//                 <li>
//                     <Link to="/about">About</Link>
//                 </li>
//             </ul>
//         </nav>
//     )
// }


// const App = () => {

//     return (
//         <div>
//     <Router>
//     <Navbar/>
//             <Switch>
//                 <Route exact path='/'> <AllProjects/> </Route>
                
//                 <Route path='/project/:id' > <Project/> </Route>
                
//                 <Route exact path='/about'> <About/> </Route>
                
//                 <Route path='*'> <Error/> </Route>
//             </Switch>
//     </Router>
//     </div>)

// }



export default App








// ########## Routing ###########
// ############################################
// ############################################
// ############################################


// import { BrowserRouter as Router, Route, Switch} from 'react-router-dom'
// import {Link, useParams} from 'react-router-dom'

// const Home = () => {
//     return (
//         <>
//             <h1>Home Page</h1>
//         </>
//     )
// }

// const AllProjects = () => {
//     const [prs, setPrs] = useState(projects)
//     return (
//         <div>
//             <h1>Projects</h1>
//             {prs.map((pr) => {
//                 return (
//                     <div key={pr.id}>
//                         <h2>Name: {pr.name}</h2>
//                         <h2>Level: {pr.level}</h2>
//                         <Link to={`/project/:${pr.id}`}>Learn more</Link>
//                     </div>
//                 )
//             })}
//         </div>
//     )
// }

// const Project = () => {
    
//     const [projectName, setProjectName] = useState('adsla')
//     const [projectLevel, setProjectLevel] = useState('')
//     const {id} = useParams();
//     console.log(id);
//     const newProject = projects.find((project) => project.id === parseInt(id))
//     console.log('lol : ',newProject);
//     // useEffect(() => {
//     //     setProjectName(newProject.name)
//     //     setProjectLevel(newProject.level)
//     // },[]);

//     return (
//         <div>
//             <h1>{projectName}</h1>
//             <h2>{projectLevel}</h2>
//             <Link to="/">Back to Homepage</Link>
//         </div>
//     )
// }


// const About = () => {
//     return (
//         <>
//             <h1>About Page</h1>
//         </>

//     )
// }


// const Navbar = () => {
//     return (
//         <nav>
//             <ul>
//                 <li>
//                     <Link to="/">Home</Link>
//                 </li>
//                 <li>
//                     <Link to="/about">About</Link>
//                 </li>
//             </ul>
//         </nav>
//     )
// }

// const Error = () => {
//     return (
//         <>
//             <h1>Error Page</h1>
//         </>
//     )
// }

// const App = () => {

//     return (
//         <div>
//     <Router>
//     <Navbar/>
//             <Switch>
//                 <Route exact path='/'> <AllProjects/> </Route>
                
//                 <Route path='/project/:id' > <Project/> </Route>
                
//                 <Route exact path='/about'> <About/> </Route>
                
//                 <Route path='*'> <Error/> </Route>
//             </Switch>
//     </Router>
//     </div>)

// }






// ########## Custom Hooks ###########
// ############################################
// ############################################
// ############################################



// const url = 'https://api.github.com/users'

// const useFetch = (url) => {  // everytime you tryin' to make a custom hook,
//                              //   stick the work 'use' at the beginning    
//         const [loading, setLoading] = useState(true)
//         const [products, setProducts] = useState([])
    
//         const getProducts = async () => {
//             const response = await fetch(url);
//             const products = await response.json();
//             setProducts(products)
//             setLoading(false)
//         };
    
//         useEffect(() => {
//             getProducts();
//         }, []);

//         return {loading, products}
// }

// const App = () => {

//     const {loading, products} = useFetch(url)

//     console.log(products);
    
//     return (
//         <div>
//             <h1>{loading ? 'loading...' : 'projects'}</h1>
//         </div>
//     )

// }



// ########## useContext ###################### // u can use this 'probably' in showcasing projects, or tasks
// ############################################
// ############################################
// ############################################

// essentially the usecase of this one is when you try to pass a function down to another
// component that is 2 lvl's or more deeper than where the function actually currently is.





// ########## useReducer ###########  USE THIS FOR ADDING TASKS
// ############################################
// ############################################
// ############################################




// const reducer = (state, action) => {
//     console.log(state);
//     if (action.type === 'add_item'){
//         const newusers = [...state.users, action.payload]
//         return {
//             ...state,
//             users: newusers,
//             isModalOpen: true,
//             modalContent: 'item added'
//         }
//     }
//     if(action.type === 'no_value'){
//         return {
//             ...state,
//             isModalOpen: true,
//             modalContent: 'please enter value'
//         }
//     }

//     if(action.type === 'close_modal'){
//         return {
//             ...state,
//             isModalOpen: false
//         }
//     }

//     if(action.type === 'remove_item'){
//         const newusers = state.users.filter((person) => person.id !== action.payload)
//        return {
//             ...state,
//             users: newusers,
//             isModalOpen: true,
//             modalContent: 'item removed'
//         }
//     }

//     throw new Error('No matching action type')
// };


// const defaultState = {
//     users: [],
//     isModalOpen: true,
//     modalContent: 'hello world'
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




// const App = () => {
    
//     const [name, setName] = useState('')
//     const [state, dispatch] = useReducer(reducer, defaultState)

//     const handleSubmit = (e) => {
//         e.preventDefault();
//         if (name) {
//             const newItem = { id: new Date().getTime().toString(), name}
//             dispatch({type: 'add_item', payload: newItem })
//             setName('')
//         } else {
//             dispatch({type: 'no_value'})
//         }
//     };

//     const closeModal = () => {
//         dispatch({ type: 'close_modal' })
//     }

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
//                     <div key={person.id} className="item">
//                         <h4>{person.name}</h4>
//                         <button onClick={() => 
//                             dispatch({ type: 'remove_item', payload: person.id})
//                         }>remove</button>
//                     </div>
//                 )
//             })}
//         </>
//     )

// }






// ########## for signup/login page ###########
// ############################################
// ############################################
// ############################################

// const [firstName, setFirstName] = useState('')
// const [email, setEmail] = useState('')
// const [age, setAge] = useState('')
// const [users, setusers] = useState([])


// const handleSubmit = (e) => {
//     e.preventDefault()
//     if (firstName && email && age) {
//         const person = {id: new Date().getTime().toString(), firstName, email, age}
//         setusers((users) => {
//             return [...users, person]
//         })
//         console.log('person: ',person);
//     }
//     console.log('users: ',users);
// }

// return (
//     <>
//         <div>
//             <form onSubmit={handleSubmit}>

//                 <label htmlFor="firstName">Name: </label>
//                 <input type="text" name='firstName' 
//                 value={firstName} onChange={(e) => setFirstName(e.target.value)}/> <br/>

//                 <label htmlFor="email">Email: </label>
//                 <input type="email" name='email' 
//                 value={email} onChange={(e) => setEmail(e.target.value)}/> <br/>

//                 <label htmlFor="age">Age: </label>
//                 <input type="age" name='age' 
//                 value={age} onChange={(e) => setAge(e.target.value)}/> <br/>

//             <button type="submit">Submit</button>
//             </form>


//         </div>
//     </>
// )




// ########## MAIN ############
// ############################
// ############################
// ############################

// const [prs, setPrs] = useState(projects)
// return (
//     <>
//         <h1 className='header'>Project Manager</h1><br />
//         <h3 className='signin-login-header'>signin/login</h3>
//         <hr />
//         <br />


//         {prs.map((pr) => {
//             return (
//                 <section className='projectlist'>
//                 <a href="./project.html">
//             <h1>{pr.name}</h1>
//             <h3>level {pr.level}</h3>
//             <h3>priority {pr.priority}</h3>
//                 </a>
//                 </section>

//             )
//         })}

//         <div className='new-projects-showcase'>
//             <h3>Newly created Projects</h3>
//             <hr />
//             <h1>Project 4</h1>
//             <h1>Project 3</h1>
//         </div>
//     </>
// )


// ########### useEffect ############
// ##############################
// ##############################
// useEffect takes two arguments: 1- a callback function 2- how many times it is rendered
// useEffect(() => {

// }, []) // if empty, it only renders once, if not, it will render whenever the variable inside it will change




// ########### useState ############
// ##############################
// ##############################

// const App = () => {
//     const [prs, setPrs] = useState(projects)
//     const removeItem = (id) => {
//         const newProjects = prs.filter((pr) => pr.id !== id)
//         setPrs(newProjects)
//     }
//     return (
//         <React.Fragment>
//             {prs.map((pr) => {
//                 return (
//                     <div key={pr.id}>
//                         <h1>Name: {pr.name}</h1>
//                         <h1>Priority: {pr.priority}</h1>
//                         <h1>Level: {pr.level}</h1>
//                         <button onClick={() => removeItem(pr.id)}>Delete</button>
//                     </div>
//                 )
//             })}
//         </React.Fragment>
//     )
// }








// ########### PROPS ############
// ##############################
// ##############################
// function View() {
//     return (
//     <div>
//     <h1>Project Manager</h1>
//        <hr />
//        <section className='projectlist'>
//            <Project name={projects.project1.name} priority={projects.project1.priority} level={projects.project1.level}/>
//            <Project name={projects.project2.name} priority={projects.project2.priority} level={projects.project2.level} />
//        </section>
       
       
//        </div>
//     );
//    }
   
//    const Project = (props) => {
   
//        const {name, priority, level} = props
//        return (
//            <article className='project'>
//                <h1>Project Name</h1>
//                <h4>Project id</h4>
//                <h4>{name}</h4>
//                <h4>{priority}</h4>
//                <h4>{level}</h4>
//            </article>
//        )
   
//    }