import React, { useState, useEffect } from 'react';
import './css/profile.css';


const Profile = ({ userId }) => {
  const [user, setUser] = useState(null);
  const [isEditing, setIsEditing] = useState(false);
  const [editedName, setEditedName] = useState('');
  const [editedEmail, setEditedEmail] = useState('');
  const [editedBio, setEditedBio] = useState('');
  const [edtiedPassword, setEdtiedPassword] = useState('')

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
        console.error('data: ', data);
        setUser(data);
        setEditedName(data.name);
        setEditedEmail(data.email);
        // setEdtiedPassword(data.password)
        setEditedBio(data.bio || '');
      })
      .catch((error) => {
        console.error(error);
      });
  }, [userId]);

  const handleEditClick = () => {
    setIsEditing(true);
  };

  const handleSaveClick = async () => {

  // const hashedPassword = await bcrypt.hash(edtiedPassword, bcrypt.genSaltSync(bcrypt.DefaultCost));

  const updatedUser = {
    ...user,
    name: editedName,
    email: editedEmail,
    password: edtiedPassword,
    bio: editedBio,
  };

  fetch(`http://localhost:8080/api/users/${userId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(updatedUser),
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      } else {
        throw new Error('Failed to update user data');
      }
    })
    .then((data) => {
      setUser(data);
      setIsEditing(false);
    })
    .catch((error) => {
      console.error(error);
    });
  };

  const handleCancelClick = () => {
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
              dir='ltr'
              onChange={(e) => setEditedEmail(e.target.value)}
              className='user-info-input'
            />
          ) : (
            <span>{editedEmail}</span>
          )}
        </div>
        <div>
          {isEditing ? (
            <>
            <strong>پسورد جدید :</strong>
            <input
              type="text"
              value={edtiedPassword}
              dir='ltr'
              onChange={(e) => setEdtiedPassword(e.target.value)}
              className='user-info-input'
            />
            </>
          ) : (
            <></>
          )}
        </div>
        <div>
          <strong>بیوی کاربر :</strong>
          {isEditing ? (
            <input
              type="text"
              value={editedBio}
              onChange={(e) => setEditedBio(e.target.value)}
              className='user-info-input'
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
