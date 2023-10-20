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
      <h1 className="profile-heading">حساب کاربری</h1>
      <div className="profile-info">
        <p> <strong>نام کاربری :</strong> {user.name}</p>
        <p> {user.email} <strong>: آدرس ایمیل</strong> </p>
        <p> <strong>بیوی کاربر : </strong> {user.bio || 'ارائه نشده'}</p>
        <br />
        <div className='user-cards'>
          <h2 className='user-cards-title-h2'>کارت های کاربر</h2>
          
        {user.cards.map((card) => {

        return (
          <div key={card.id} className='user-card'>
            <div className='user-card-details'>
              <h3> نام کارت : </h3>
              <p>{card.name}</p>
            </div>
            <div className='user-card-details'>
              <h3>توضیحات : </h3>
              <p>{card.description}</p>
            </div>
            <div className='user-card-details'>
              <h3>تاریخ : </h3>
              <p>{card.dates.join(', ')}</p>
            </div>
            <div className='user-card-checklists'>
              <h3>چکلیست های کارت</h3>
              <div className='user-card-checklist'>
                {card.checklists && card.checklists.length > 0 ? (
                  card.checklists.map((checklist) => (
                    <h4 key={checklist.id}>{checklist.name}</h4>
                  ))
                ) : (
                  <p>بدون چکلیست</p>
                )}
              </div>
            </div>
            <hr />
          </div>
        );

        })}



        </div>


      </div>
    </div>
  );
};

export default Profile;
