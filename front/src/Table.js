import React, { useState, useEffect } from 'react';

import './css/table.css'
const Table = () => {
  const [lists, setLists] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/table', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        });

        if (!response.ok) {
          throw new Error('Failed to get table info');
        }

        const tableData = await response.json();
        setLists(tableData);
      } catch (error) {
        console.error('Error getting table info:', error);
      }
    };

    fetchData();
  }, []);

  return (
    <table>
      <thead>
        <tr>
          <th>List</th>
          <th>Card Name</th>
          <th>Label</th>
          <th>Members</th>
          <th>Due Date</th>
        </tr>
      </thead>
      <tbody>
        {lists.map((list) =>
          list.cards.map((card) => (
            <tr key={card.id}>
              <td>{list.name}</td>
              <td>{card.name}</td>
              <td>{card.label || 'N/A'}</td>
              <td>{card.members.map((member) => member.name).join(', ')}</td>
              <td>{card.dueDate ? new Date(card.dueDate).toLocaleDateString() : 'N/A'}</td>
            </tr>
          ))
        )}
      </tbody>
    </table>
  );
};

export default Table;
