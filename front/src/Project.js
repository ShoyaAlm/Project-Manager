import projects from './projects.js'

import {useState, useEffect} from 'react'
import {Link, useParams} from 'react-router-dom'


export const Project = () => {
    
    const [projectName, setProjectName] = useState('defautl')
    const [projectLevel, setProjectLevel] = useState(1)
    const {id} = useParams();
    const newProject = projects.find((pr) => pr.id === parseInt(id))
    console.log(newProject);
    useEffect(() => {
        setProjectName(newProject.name)
        setProjectLevel(newProject.level)
    },[]);

    return (
        <div style={{textAlign:'center'}}>
            <h1>{projectName}</h1>
            <h2>{projectLevel}</h2>
            <Link to="/">Back to Homepage</Link>
        </div>
    )
}