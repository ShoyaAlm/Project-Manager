import React, { useState, useEffect } from 'react';

import './css/workspace.css'

import { getJwtFromCookie } from './App';
import jwt_decode from 'jwt-decode'

import { Link, useHistory } from 'react-router-dom';


const Workspace = () => {
  const [boards, setBoards] = useState([]);
  const [otherBoards, setOtherBoards] = useState([]);
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
    const fetchData = async () => {
      try {
        // Fetch user boards
        const userBoardsResponse = await fetch(`http://localhost:8080/api/boards/user/${getUserId()}`);
        if (userBoardsResponse.ok) {
          const userData = await userBoardsResponse.json();
          setBoards(userData);
  
          // Fetch all other boards
          const allBoardsResponse = await fetch('http://localhost:8080/api/boards', {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
            },
          });
          if (allBoardsResponse.ok) {
            const allBoardsData = await allBoardsResponse.json();
            
            // Filter out user's boards using the userData directly
            if (userData && userData.length > 0) {
              // User has boards, filter out user's boards using the userData directly
              const otherUserBoards = allBoardsData.filter(
                (board) => !userData.some((userBoard) => userBoard.id === board.id)
              );
              setOtherBoards(otherUserBoards);
            } else {
              // User has no boards, set otherBoards to allBoardsData
              setOtherBoards(allBoardsData);
            }
          } else {
            console.error('Failed to fetch all boards');
          }
        } else {
          console.error('Failed to fetch user boards');
        }
      } catch (error) {
        console.error('Error during fetch:', error);
      }
    };
  
    fetchData();
  }, []);

  
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

  const handleDeleteBoard = async (boardId) => {
    try {
      const response = await fetch(`http://localhost:8080/api/boards/${boardId}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        }
      });
  
      if (response.ok) {
        // Update the state by removing the deleted board
        setBoards(boards.filter(board => board.id !== boardId));
      } else {
        console.error('Failed to delete board');
      }
    } catch (error) {
      console.error('Error during fetch:', error);
    }
  }
  

  return (
    <div className="workspace-container">
      <h1 style={{ textAlign: 'center', fontFamily: 'vazirmatn', color: '#333', marginBottom: '20px', padding: '10px', backgroundColor: '#f8f8f8', borderRadius: '10px', boxShadow: '0 0 10px rgba(0, 0, 0, 0.1)', marginRight: '630px', marginLeft: '630px' }}>صفحه بورد</h1>
  
      <div className="boards-container">
        {/* User's boards */}
        <div className="user-boards-container">
        <h2>بوردهای شما</h2>
        {boards && boards.length > 0 ? (
          <div className="user-boards">
          {boards.map((board, index) => (
            <div key={index} className="board" onClick={() => handleBoardClick(board.id)}>
              <p>{board.name}</p>
              <button onClick={() => handleDeleteBoard(board.id)} style={{fontFamily:'vazirmatn', fontSize:'10px'}}>حذف بورد</button>
            </div>
          ))}
        </div>
        ) : (
          <p className="no-boards-message">شما بوردی ندارید!</p>
        )}
      </div>
  
        {/* All other boards */}
        <div className="other-boards-container">
          <h2>بوردهای دیگر</h2>
          {otherBoards && otherBoards.length > 0 ? (
            <div className="other-boards">
              {otherBoards.map((board, index) => (
                <div key={index} className="board" onClick={() => handleBoardClick(board.id)}>
                  <p>{board.name}</p>
                </div>
              ))}
            </div>
          ) : (
            <p className="no-boards-message">شما به بوردها دسترسی ندارید.</p>
          )}
        </div>
      </div>
  
      {/* Create Board Section */}
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
          <p>ساختن بورد</p>
        )}
      </div>
  
      <div className="create-board-container">
        {isCreatingBoard ? (
          <button onClick={handleCreateBoard}>ساختن بورد</button>
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

