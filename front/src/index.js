import React from 'react'
import ReactDOM from 'react-dom'
import './index.css'

function Greet() {
 return (
 <div>
 <h1>Project Manager</h1>
    <hr />
    <section>
        <Project/>
        <Project/>
        <Project/>
    </section>
    
    
    </div>
 );
}

const Project = () => {

    return (
        <article>
            <h2>Project Name</h2>
            <h4>Project id</h4>
        </article>
    )

}


ReactDOM.render(<Greet/>, document.getElementById("root"))