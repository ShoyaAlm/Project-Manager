// Function to create and append a project element
function createProjectElement(project) {
    const projectDiv = document.createElement('div');
    projectDiv.classList.add('project-div');

    const projectName = document.createElement('h3');
    projectName.classList.add('project-name');
    projectName.textContent = project.name; // Access the 'name' property

    const projectId = document.createElement('h4');
    projectId.classList.add('project-id');
    projectId.textContent = `ID: ${project.id}`; // Access the 'id' property

    projectDiv.appendChild(projectName);
    projectDiv.appendChild(projectId);

    return projectDiv;
}

// Fetch projects from the API endpoint
fetch('/api/projects')
    .then(response => response.json())
    .then(data => {
        const projectsDiv = document.querySelector('.projects-div');
        data.forEach(project => {
            const projectElement = createProjectElement(project);
            projectsDiv.appendChild(projectElement);
        });
    })
    .catch(error => console.error('Error fetching projects:', error));