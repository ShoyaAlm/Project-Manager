// import projects from './projects.js'
import './css/index.css'
import './css/app.css'
import React, { useState, useEffect, useContext } from 'react'
// import {useState} from 'react'


import { BrowserRouter as Router, Route, Switch, Redirect} from 'react-router-dom'
import {Link, useParams} from 'react-router-dom'


import { AllLists } from './Board'
import { HandleSignupLogin } from './signup-login.js'
import { Card } from './Card'
import Profile from './Profile'
import jwt_decode from 'jwt-decode'
import Workspace from './Workspace'

import Table from './Table'
import Scheduler from './Scheduler'

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


      const styles = {
        container: {
          textAlign: 'left',
          marginBottom: '10px',
          padding: '10px',
          borderBottom: '1px solid #ddd',
          backgroundColor: '#f4f4f4',
          border: '1px solid #ccc',
          borderRadius: '5px',
          display: 'flex',
          justifyContent: 'flex-end',  // Align items to the right
        },
        link: (currentView) => ({
          cursor: 'pointer',
          marginRight: '20px',  // Adjusted margin to separate links
          color: view === currentView ? '#007BFF' : '#000',
          fontWeight: view === currentView ? 'bold' : 'normal',
        }),
      };
      
      

    return (
        <div>



            

            <Router>

            <Link to="/" style={{ textDecoration: 'none', fontFamily: 'bardiya', color: 'black', display: 'block', textAlign: 'center' }}>
              <h1 style={{ margin: '20px 0', fontSize: '2.5rem', fontWeight: 'bold', fontFamily:'shabnam', color: 'white' }}>مدیریت پروژه</h1>
            </Link>
           {user ? (
            
            <div className="header-container">
              <h3 className="welcome-message" style={{fontSize:'16px'}}>خوش آمدی {user.name}!</h3>
              <div className="link-container">
                <Link to="/profile" className="nav-link profile-link">
                  <h3 style={{fontFamily: 'shabnam'}}>پروفایل من</h3>
                </Link>
                <Link to="/workspace" className="nav-link workspace-link">
                  <h3>Workspace</h3>
                </Link>
              </div>
              <Link to="#" className="logout-link" onClick={handleLogout}>
                خروج
              </Link>
              <hr className="divider" />
            </div>
                        ) : (
                          <Link to="/signup" style={{ textDecoration: 'none', display: 'flex', justifyContent: 'flex-end' }}>
                            <div className="auth-button">
                              <h3 style={{fontFamily: 'shabnam' }}>ورود/ثبت نام</h3>
                              <hr />
                            </div>
                          </Link>
            )}


            <div style={styles.container}>
              {/* <span style={styles.link('timeline')} onClick={() => setView('timeline')}>
                تایملاین
              </span> */}
              <span style={styles.link('table')} onClick={() => setView('table')}>
                جدول کارت ها
              </span>
              <span style={styles.link('board')} onClick={() => setView('board')}>
                صفحه بورد
              </span>
            </div>





              <Switch>
              <Route path='/board/:boardId/lists'>
                {view === 'board' && <AllLists /> }
                {view === 'table' && <Table /> }
                {view === 'timeline' && <Scheduler /> }
              </Route>
              <Route path='/workspace'>
                  <Workspace />
                </Route>
                <Route path='/boards/:boardId/lists/:listId/cards/:cardId'>
                  <Card />
                </Route>
                <Route path='/signup'>
                  <HandleSignupLogin />
                </Route>
                <Route path='/profile'>
                  <Profile />
                </Route>

                <Redirect from="/" to="/workspace" />

              </Switch>

        

            </Router>


        </div>
    )


}


export default App


