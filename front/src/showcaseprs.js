import projects from './projects.js'
import React, { useEffect, useContext } from 'react'
import {useState} from 'react'

import {Link, useParams} from 'react-router-dom'
import {Router} from 'react-router-dom'

const ProjectContext = React.createContext()
// has 2 components : provider - consumer

export const AllProjects = () => {
    const [prs, setPrs] = useState(projects)
    const removeProject = (id) => {  // i want to pass this func into the SinglePerson, and for that i
        setPrs((prs) => {            // create a context using createContext
            return prs.filter((pr) => pr.id !== id)
        });
    };
    
    return (        //and in here, i'll pass it as a provider
        <ProjectContext.Provider value={{removeProject, prs}}>
            <SingleProject />
        </ProjectContext.Provider>
    )

}


const SingleProject = () => {
    const {removeProject, prs} = useContext(ProjectContext) // I'll be receiving it here
    return (
        <div>
            {prs.map((pr) => {
                return (
                    <div key={pr.id} style={{display:'inline-block', width:'500px', height:'200px', marginLeft: '200px' }}>
                       {/* <Router>
                       <Link to={`/project/:${pr.id}`}></Link>
                       </Router> */}
                       <h1>Name: {pr.name}</h1>
                        <h2>Level: {pr.level}</h2>
                        <button onClick={() => removeProject(pr.id)}>remove</button>
                    </div>
                )
            })}
        </div>
    )
}



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

