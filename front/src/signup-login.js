import React, { useState } from 'react';
import { BrowserRouter as Router, Switch, Route, Link, Redirect } from 'react-router-dom';
import './css/signup.css';
import Cookies from 'js-cookie'


export const HandleSignupLogin = () => {
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [isLogin, setIsLogin] = useState(false);
    const [users, setUsers] = useState([]);

    const handleSubmit = (e) => {
        console.log('e : ', e);
        e.preventDefault();
        if (name && email && password) {

            const signupData = {
                name : name,
                email : email,
                password : password
            }

            fetch('http://localhost:8080/api/signup', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(signupData),
            })

            .then((response) => {
                if (response.ok){
                    response.json().then((data) => {
                        console.log('Signup Successful: ', data);
                    })
                } else {
                    console.log('Signup Failed : ', response.statusText);
                }
            })
            .catch((error) => {
                console.log('Error : ', error);
            })



            // const person = { id: new Date().getTime().toString(), name, email, password };
            // setUsers(person);
            console.log('users: ', users);
        } else {
            console.log('Enter values for all inputs!');
        }
    };


    const handleLogin = (e) => {
        e.preventDefault();

        const loginData = {
            email: email,
            password: password,
        };

        console.log('Login data:', loginData);


        fetch('http://localhost:8080/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(loginData),
        })
            .then((response) => {
                if (response.ok) {
                    // Login successful
                    response.json().then((data) => {
                        console.log('Login successful:', data);

                        const token = data.token;

                        Cookies.set('jwtToken', token);

                        window.location.href = '/workspace';
                        // Redirect to the main page or perform other actions
                    });
                } else {
                    // Login failed
                    console.error('Login failed:', response.statusText);
                }
            })
            .catch((error) => {
                console.error('Error:', error);
            });
    };



    return (
        <Router>
            <div className="signupFrm">
                <div className="wrapper">
                    <Switch>
                        <Route path="/login">
                            <form
                                action=""
                                className="form"
                                onSubmit={(e) => handleLogin(e)}
                            >
                                <h1 className="title" style={{ textAlign: 'center' }}>
                                    ورود
                                </h1>
                                <div className="inputContainer">
                                    <input
                                        type="text"
                                        className="input"
                                        placeholder="a"
                                        onChange={(e) => setEmail(e.target.value)}
                                    />
                                    <label htmlFor="" className="label">
                                        ایمیل
                                    </label>
                                </div>
                                <div className="inputContainer">
                                    <input
                                        type="text"
                                        className="input"
                                        placeholder="a"
                                        onChange={(e) => setPassword(e.target.value)}
                                    />
                                    <label htmlFor="" className="label">
                                        رمز عبور
                                    </label>
                                </div>
                                <input
                                    type="submit"
                                    className="submitBtn"
                                    value="ورود"
                                />
                                <Link to="/signup" className="changeToLogin">
                                    ثبت نام
                                </Link>
                            </form>
                        </Route>
                        <Route path="/signup">
                            <form
                                action=""
                                className="form"
                                onSubmit={(e) => handleSubmit(e)}
                            >
                                <h1 className="title" style={{ textAlign: 'center' }}>
                                    ثبت نام
                                </h1>
                                <div className="inputContainer">
                                    <input
                                        type="text"
                                        className="input"
                                        placeholder="a"
                                        onChange={(e) => setName(e.target.value)}
                                    />
                                    <label htmlFor="" className="label">
                                        نام
                                    </label>
                                </div>
                                <div className="inputContainer">
                                    <input
                                        type="text"
                                        className="input"
                                        placeholder="a"
                                        onChange={(e) => setEmail(e.target.value)}
                                    />
                                    <label htmlFor="" className="label">
                                        ایمیل
                                    </label>
                                </div>
                                <div className="inputContainer">
                                    <input
                                        type="text"
                                        className="input"
                                        placeholder="a"
                                        onChange={(e) => setPassword(e.target.value)}
                                    />
                                    <label htmlFor="" className="label">
                                        رمز عبور
                                    </label>
                                </div>
                                <input
                                    type="submit"
                                    className="submitBtn"
                                    value="ثبت نام"
                                />
                                <Link to="/login" className="changeToLogin">
                                    ورود
                                </Link>
                            </form>
                        </Route>
                        <Redirect from="/" to="/signup" />
                    </Switch>
                </div>
            </div>
        </Router>
    );
};
