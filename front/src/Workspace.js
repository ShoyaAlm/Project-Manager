import React, { useState, useEffect } from 'react';

import './css/workspace.css'

import { getJwtFromCookie } from './App';
import jwt_decode from 'jwt-decode'

import { Link, useHistory } from 'react-router-dom';


const Workspace = () => {
  const [boards, setBoards] = useState([]);
  const [newBoardName, setNewBoardName] = useState('');
  const [isCreatingBoard, setIsCreatingBoard] = useState(false);

  // const history = useHistory();


  const handleBoardClick = (boardId) => {
    // Redirect to the Board component with the specific boardId

    window.location.href = `/board/${boardId}/lists`;

  };

  const getUserId = () => {
    
    const findUser = () => {

      try {
          const jwt = getJwtFromCookie();
          if (jwt) {
              const decoded = jwt_decode(jwt);
              const user1 = decoded;
              console.log(user1);
              // console.log(user1);
              return user1;
          }
      } catch (error) {
          console.log(error);
      }
  }
    // Implement your logic to get the user ID
    // This could be from the decoded JWT, from a cookie, or another method
    // For simplicity, assuming there's a global function `findUser` that returns the user object
    const user = findUser();
    return user.user_id ? user.user_id : null;
  };


  useEffect(() => {
    // Fetch user boards when the component mounts
    const fetchUserBoards = async () => {
      try {
        const response = await fetch(`http://localhost:8080/api/boards/user/${getUserId()}`);
        if (response.ok) {
          const data = await response.json();
          setBoards(data);
        } else {
          console.error('Failed to fetch user boards');
        }
      } catch (error) {
        console.error('Error during fetch:', error);
      }
    };

    fetchUserBoards();
  }, []); // Empty dependency array means this effect runs once when the component mounts

  const handleCreateBoard = async () => {
    if (newBoardName.trim() !== '') {
      // Make a POST request to create a new board
      try {
        const response = await fetch('http://localhost:8080/api/boards', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            name: newBoardName,
            user_id: getUserId(),
            // You may need to include other user-related information in the request body
          }),
        });

        if (response.ok) {
          // Update the state with the newly created board
          const newBoard = await response.json();
          setBoards([...boards, newBoard]);
          setNewBoardName('');
        } else {
          console.error('Failed to create board');
        }
      } catch (error) {
        console.error('Error during fetch:', error);
      }
    }

    setIsCreatingBoard(false);
  };

  const handleCreateBoardClick = () => {
    setIsCreatingBoard(true);
  };


  return (
    <div className="workspace-container">
      <h1>Boards</h1>

      <div className="boards-container">
        {boards.map((board, index) => (
          <div key={index} className="board" onClick={() => handleBoardClick(board.id)}>
            <p>{board.name}</p>
          </div>
        ))}

        <div className={`board create-board ${isCreatingBoard ? 'editing' : ''}`} onClick={handleCreateBoardClick}>
          {isCreatingBoard ? (
            <input
              type="text"
              placeholder="Enter board name"
              value={newBoardName}
              onChange={(e) => setNewBoardName(e.target.value)}
              onBlur={handleCreateBoard}
              autoFocus
            />
          ) : (
            <p>Create Board</p>
          )}
        </div>
      </div>

      <div className="create-board-container">
        {isCreatingBoard ? (
          <button onClick={handleCreateBoard}>Create Board</button>
        ) : (
          <input
            type="text"
            placeholder="Enter board name"
            value={newBoardName}
            onChange={(e) => setNewBoardName(e.target.value)}
          />
        )}
      </div>
    </div>
  );

};

export default Workspace;

