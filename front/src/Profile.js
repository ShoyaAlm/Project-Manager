import React, { useState, useEffect } from 'react';
import './css/profile.css';


const Profile = ({ userId }) => {
  const [user, setUser] = useState(null);
  const [isEditing, setIsEditing] = useState(false);
  const [editedName, setEditedName] = useState('');
  const [editedEmail, setEditedEmail] = useState('');
  const [editedBio, setEditedBio] = useState('');

  useEffect(() => {
    fetch(`http://localhost:8080/api/users/${userId}`)
      .then((response) => {
        if (response.ok) {
          return response.json();
        } else {
          throw new Error('Failed to fetch user data');
        }
      })
      .then((data) => {
        setUser(data);
        setEditedName(data.name);
        setEditedEmail(data.email);
        setEditedBio(data.bio || '');
      })
      .catch((error) => {
        console.error(error);
      });
  }, [userId]);

  const handleEditClick = () => {
    setIsEditing(true);
  };

  const handleSaveClick = () => {
    // Perform save logic (e.g., update data on the backend)
    // After saving, set isEditing to false
    setIsEditing(false);
  };

  const handleCancelClick = () => {
    // Reset edited values to original values
    setEditedName(user.name);
    setEditedEmail(user.email);
    setEditedBio(user.bio || '');
    setIsEditing(false);
  };

  return (
    <div className="profile-container">
      <h1 className="profile-heading">حساب کاربری</h1>
      <div className="profile-info">
        <div>
          <strong>نام کاربری :</strong>
          {isEditing ? (
            <input
              type="text"
              value={editedName}
              onChange={(e) => setEditedName(e.target.value)}
              className='user-info-input'
            />
          ) : (
            <span>{editedName}</span>
          )}
        </div>
        <div>
          <strong>آدرس ایمیل :</strong>
          {isEditing ? (
            <input
              type="text"
              value={editedEmail}
              onChange={(e) => setEditedEmail(e.target.value)}
              className='user-info-input'
            />
          ) : (
            <span>{editedEmail}</span>
          )}
        </div>
        <div>
          <strong>بیوی کاربر :</strong>
          {isEditing ? (
            <textarea
              value={editedBio}
              onChange={(e) => setEditedBio(e.target.value)}
            />
          ) : (
            <span>{editedBio}</span>
          )}
        </div>
        {isEditing ? (
          <div className="edit-buttons">
            <button onClick={handleCancelClick}>لغو</button>
            <button onClick={handleSaveClick}>ثبت</button>
          </div>
        ) : (
          <div className="edit-buttons">
            <button onClick={handleEditClick}>ویرایش</button>
          </div>
        )}
      </div>
    </div>
  );
};



export default Profile;
