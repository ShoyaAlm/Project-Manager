fetch('/api/projects')
            .then(response => response.json())
            .then(data => {
                const projectList = document.getElementById('projectList');
                data.forEach(project => {
                    const listItem = document.createElement('li');
                    listItem.textContent = project.Name;
                    projectList.appendChild(listItem);
                });
            })
            .catch(error => console.error('Error fetching projects:', error));
    