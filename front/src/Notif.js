import React, { useState, useEffect } from 'react';

function Notifications() {
  const [notifications, setNotifications] = useState([]);
  const [unreadCount, setUnreadCount] = useState(0);

  useEffect(() => {
    // Fetch notifications from the API endpoint
    fetch('http://localhost:8080/api/notifications')
      .then((response) => response.json())
      .then((data) => {
        setNotifications(data.notifications);
        setUnreadCount(data.unreadCount);
      })
      .catch((error) => console.error('Error fetching notifications: ', error));
  }, []);

  return (
    <div>
      <h3>Notifications</h3>
      {unreadCount > 0 && (
        <p>You have {unreadCount} unread notifications.</p>
      )}
      <ul>
        {notifications.map((notification) => (
          <li key={notification.notification_id}>
            {notification.message}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default Notifications;
