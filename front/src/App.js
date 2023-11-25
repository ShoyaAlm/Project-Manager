// import projects from './projects.js'
import './css/index.css'

import React, { useState, useEffect, useContext } from 'react'
// import {useState} from 'react'


import { BrowserRouter as Router, Route, Switch} from 'react-router-dom'
import {Link, useParams} from 'react-router-dom'
import NotificationDropdown from './Notif.js'

import { AllLists } from './Board'
import { HandleSignupLogin } from './signup-login.js'
import { Card } from './Card'
import Profile from './Profile'
import jwt_decode from 'jwt-decode'

import Table from './Table'
import Timeline from './Timeline'

export const getJwtFromCookie = () => {
  // Get all cookies as a string
  const cookies = document.cookie;

  // Split the cookies into an array
  const cookieArray = cookies.split(';');

  // Find the cookie that contains the JWT (you may use a specific cookie name)
  for (let i = 0; i < cookieArray.length; i++) {
    const cookie = cookieArray[i].trim();
    if (cookie.startsWith('jwtToken=')) {
      // Extract the JWT value (remove 'jwtToken=')
      const jwt = cookie.substring('jwtToken='.length);
      return jwt;
    }
  }

  // Return null if JWT cookie is not found
  return null;
};

const App = () => {

    const [user, setUser] = useState(null)
    const [view, setView] = useState('board');

    useEffect(() => {
        // Check if a JWT is stored (in this example, we assume it's stored in a cookie)
        const jwt = getJwtFromCookie(); // You should replace this with the actual code to get the JWT
    
        if (jwt) {
          // Decode the JWT to get user data
          const decoded = jwt_decode(jwt);
          setUser(decoded); // Set the user data from the JWT
          console.log(user);
        } else {
          setUser(null); // No JWT, so no user data
        }
      }, []);


      
      

      const handleLogout = () => {
        // Clear the JWT token from cookies (you need to replace 'your_jwt_cookie_name' with your actual cookie name)
        document.cookie = 'jwtToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
      
        // Redirect the user to the "/signup" page
        window.location.href = '/signup'; // You can also use react-router's history for this.
      };

    return (
        <div>



            

            <Router>

           <Link to="/" style={{textDecoration:'none', fontFamily:'bardiya', color:'black'}}><h1 style={{textAlign:'center'}}>مدیریت پروژه</h1></Link> <br />

           {user ? (
            <div style={{textAlign:'right '}}>
                  <h3 style={{textAlign:'right', color:'black'}}>! خوش آمدی {user.name}</h3>
             <Link to="/profile" style={{ textDecoration: 'none' }}>
                    <h3 style={{ fontFamily: 'sans-serif', color: 'black' }}>پروفایل من</h3>
                </Link>


                {/* this is for the notifs */}
                <div>
                  <NotificationDropdown user={user}/>
                </div>


                  <Link to="#" style={{ textDecoration: 'none' }} onClick={handleLogout}>خروج</Link>
                    <hr />
                </div>
                ) : (
                <Link to="/signup" style={{ textDecoration: 'none' }}>
                    <h3 style={{ textAlign: 'right', fontFamily: 'sans-serif', color: 'black' }}>ورود/ثبت نام</h3>
                    <hr />
                </Link>
            )}


            <div style={{ textAlign: 'left', marginBottom: '10px', padding: '10px', borderBottom: '1px solid #ddd' }}>
              <span
                style={{
                  cursor: 'pointer',
                  marginRight: '20px',
                  color: view === 'board' ? '#007BFF' : '#000',
                  fontWeight: view === 'board' ? 'bold' : 'normal',
                }}
                onClick={() => setView('board')}
              >
                Board
              </span>
              |
              <span
                style={{
                  cursor: 'pointer',
                  marginLeft: '20px',
                  color: view === 'table' ? '#007BFF' : '#000',
                  fontWeight: view === 'table' ? 'bold' : 'normal',
                }}
                onClick={() => setView('table')}
              >
                Table
              </span>
              |
              <span
                style={{
                  cursor: 'pointer',
                  marginLeft: '20px',
                  color: view === 'timeline' ? '#007BFF' : '#000',
                  fontWeight: view === 'timeline' ? 'bold' : 'normal',
                }}
                onClick={() => setView('timeline')}
              >
                timeline
              </span>
            </div>





              <Switch>
              <Route exact path='/'>
                {view === 'board' ? <AllLists /> : (view === 'table' ? <Table /> : <Timeline />)}
              </Route>
                <Route path='/lists/:listId/cards/:cardId'>
                  <Card />
                </Route>
                <Route path='/signup'>
                  <HandleSignupLogin />
                </Route>
                <Route path='/profile'>
                  <Profile />
                </Route>
              </Switch>

        

            </Router>


        </div>
    )


}


export default App


