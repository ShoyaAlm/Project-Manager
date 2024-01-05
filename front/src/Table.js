import React, { useState, useEffect } from 'react';

import './css/table.css';
import { useParams } from 'react-router-dom/cjs/react-router-dom.min';



const Table = () => {
  const [lists, setLists] = useState([]);

  const { boardId } = useParams();

  const formatDate = (dateString) => {
    const options = { year: 'numeric', month: 'long', day: 'numeric' };
    const dateObject = new Date(dateString);
    
    const year = dateObject.toLocaleDateString('fa-IR', { year: 'numeric' });
    const month = dateObject.toLocaleDateString('fa-IR', { month: 'long' });
    const day = dateObject.toLocaleDateString('fa-IR', { day: 'numeric' });
  
    return `${day} ${month} ${year}`;
  };

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
        console.error('lists: ', lists);
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
          <th style={{fontFamily:'shabnam', fontSize:'15px'}}>نام لیست</th>
          <th>نام کارت</th>
          <th>برچسب</th>
          <th>اعضا</th>
          <th>موعد مقرر</th>
        </tr>
      </thead>
      <tbody>
        {lists.map((list) => (
          <React.Fragment key={list.name} >
            {list.cards && list.cards.length !== 0 ? (
              list.cards.map((card, index) => (
                <tr key={card.id}>
                  {index === 0 ? <td rowSpan={list.cards.length} style={{fontFamily:'vazirmatn'}}>{list.name}</td> : null}
                  <td style={{fontFamily:'vazirmatn'}}>{card.name}</td>
                  <td>
                    {card.label && (
                      <div
                        style={{
                          height: '20px',
                              width: '70px',
                              backgroundColor: card.label,
                              marginLeft:'auto',
                              marginRight:'auto'
                        }}
                      ></div>
                    )}
                  </td>
                  <td>{card.members.map((member) => member.name).join(', ')}</td>
                  <td style={{direction:'rtl'}}>{card.dates[1] ? formatDate(card.dates[1]) : 'ندارد'}</td>
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
