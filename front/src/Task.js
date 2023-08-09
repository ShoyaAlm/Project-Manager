import {useState} from 'react'
import {Link, useParams} from 'react-router-dom'
import projects from './cards'

export const EditTask = () => {

    const {projectId, taskId} = useParams()
    console.log('projectId: ', projectId);
    console.log('taskId: ', taskId);
    
    const [pr, setPr] = useState([])
    const newProject = projects.find((project) => project.id === parseInt(projectId))
    setPr(newProject)
    console.log('project: ', pr);
    
    const [task, setTask] = useState([])
    const newTask = newProject.tasks.find((task) => task.id === parseInt(taskId))
    setTask(newTask)
    console.log('task: ', task);

    
    const [name, setName] = useState(task.name)
    const [category, setCategory] = useState(task.category)
    const [status, setStatus] = useState(task.status)
    const [priority, setPriority] = useState(task.priority)
    
    return (
        <div style={{display:'inline-grid'}}>
            <form>
    
                <label htmlFor="name" > Name: <input type="text" value={name}/> </label>
                <label htmlFor="category"> Category: <input type="text" value={category}/> </label>
                <label htmlFor="status"> Status: <input type="text" value={status}/> </label>
                <label htmlFor="priority"> Priority: <input type="number" value={priority}/> </label>
                
                <button type='submit'>submit</button>
                
    
            </form>
        </div>
    )
    
}

