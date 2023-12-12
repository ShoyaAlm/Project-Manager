import React, { useState, useEffect } from 'react';

import './css/table.css';
import { useParams } from 'react-router-dom/cjs/react-router-dom.min';

const Table = () => {
  const [lists, setLists] = useState([]);

  const { boardId } = useParams();

  console.log(boardId);
  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch(`http://localhost:8080/api/boards/${boardId}/table`, {
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
          <th>نام لیست</th>
          <th>نام کارت</th>
          <th>برچسب</th>
          <th>اعضا</th>
          <th>موعد مقرر</th>
        </tr>
      </thead>
      <tbody>
        {lists.map((list) => (
          <React.Fragment key={list.name}>
            {list.cards && list.cards.length !== 0 ? (
              list.cards.map((card, index) => (
                <tr key={card.id}>
                  {index === 0 ? <td rowSpan={list.cards.length}>{list.name}</td> : null}
                  <td>{card.name}</td>
                  <td>
                    {card.label && (
                      <div
                        style={{
                          width: '20px',
                          height: '20px',
                          backgroundColor: card.label.toLowerCase(),
                          display: 'inline-block',
                          marginRight: '5px',
                        }}
                      ></div>
                    )}
                  </td>
                  <td>{card.members.map((member) => member.name).join(', ')}</td>
                  <td>{card.dueDate ? new Date(card.dueDate).toLocaleDateString() : 'ندارد'}</td>
                </tr>
              ))
            ) : (
              <tr>
                <td>{list.name}</td>
                <td colSpan={4}>هیچ کارتی در این لیست وجود ندارد</td>
              </tr>
            )}
          </React.Fragment>
        ))}
      </tbody>
    </table>
  );
};

export default Table;
