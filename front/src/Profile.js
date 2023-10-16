import React, { useState, useEffect } from 'react';
import './profile.css';

const Profile = () => {
  const [user, setUser] = useState(null);

  useEffect(() => {
    // Fetch the user's data from the backend
    fetch('http://localhost:8080/api/users/1') // Replace '1' with the actual user ID
      .then((response) => {
        if (response.ok) {
          return response.json();
        } else {
          throw new Error('Failed to fetch user data');
        }
      })
      .then((data) => {
        console.log('setting the user data : ', data);
        setUser(data);
      })
      .catch((error) => {
        console.error(error);
      });
  }, []);

  if (user === null) {
    return <div>Loading...</div>;
  }

  return (
    <div className="profile-container">
      <h2 className="profile-heading">User Profile</h2>
      <div className="profile-info">
        <p><strong>Name:</strong> {user.name}</p>
        <p><strong>Email:</strong> {user.email}</p>
        <p><strong>Bio:</strong> {user.bio || 'Not provided'}</p>
        {/* Display other user data as needed */}
      </div>
    </div>
  );
};

export default Profile;
