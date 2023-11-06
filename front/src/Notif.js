import React, { useState, useEffect } from 'react';

import { getJwtFromCookie } from './App';
import jwt_decode from 'jwt-decode'

import './css/notif.css'

const NotificationDropdown = ({user}) => {
  const [notifications, setNotifications] = useState([]);
  const [isOpen, setIsOpen] = useState(false);
  const [unreadCount, setUnreadCount] = useState(0);
  const [notificationsRead, setNotificationsRead] = useState(false);


  const fetchNotifications = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/notifs/${user.user_id}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      if (response.ok) {
        const data = await response.json();
  
        // Sort notifications by CreatedAt in ascending order (oldest first)
        data.sort((a, b) => new Date(a.CreatedAt) - new Date(b.CreatedAt));
  
        // Get the last 4 notifications (most recent) and reverse the order
        const last4Notifications = data.slice(-4).reverse();
  
        setNotifications(last4Notifications);
  
        // Calculate the count of unread notifications
        const newUnreadCount = last4Notifications.filter((notification) => !notification.read).length;
        setUnreadCount(newUnreadCount);
      }
    } catch (error) {
      console.error('Error fetching notifications:', error);
    }
  };
  

  const markNotificationsAsRead = async () => {
    try {
      // Send a request to mark unread notifications as read
      const unreadIds = notifications.filter((notification) => !notification.read).map((notification) => notification.id);

      if (unreadIds.length > 0) {
        const response = await fetch(`http://localhost:8080/api/notifs/${user.user_id}`, {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
          },
          // body: JSON.stringify({ ids: unreadIds }),
        });

        if (response.ok) {
          setUnreadCount(0); // Reset the unread count
          setNotificationsRead(true); // Notifications have been marked as read
        }
      }
    } catch (error) {
      console.error('Error marking notifications as read:', error);
    }
  };

  useEffect(() => {
    if (isOpen) {
      fetchNotifications();
    }
  }, [isOpen]);


  return (
    <div className="notification-dropdown" style={{ marginLeft: 'auto', marginRight: 'auto' }}>
      <button onClick={() => setIsOpen(!isOpen)}>
        {/* <img src="./icons/notif.png" alt="notif" style={{ width: '30px', height: '30px' }} /> */}
        اعلان
        {unreadCount > 0 && <span className="unread-count">{unreadCount}</span>}
      </button>
      {isOpen && (
        <div className="notification-list">
          
          {notifications && notifications.length > 0 ? (
            <div>
              {notifications.map((notification) => (
                <div key={notification.id} className="notification-item">
                    {notification.read ? (
                        // Render the message as a paragraph
                        <p>{notification.message}</p>
                    ) : (
                        // Render the message as a paragraph with "new" text
                        
                        <p> <span className="new-notification">جدید : </span> {notification.message} </p>
                    )}
                </div>
              ))}

              {user && unreadCount > 0 && !notificationsRead && (
                <button onClick={markNotificationsAsRead} className="mark-read-button">
                  Mark as Read
                </button>
              )}
            </div>
          ) : (
            <h3>هیچ اعلانی وجود ندارد</h3>
          )}



        </div>
      )}
    </div>
  );
};

// notifications && notifications.length > 0 ?
export default NotificationDropdown;
