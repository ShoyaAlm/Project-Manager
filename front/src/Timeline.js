import React from 'react';
import './css/timeline.css'; // You can style the timeline as needed
import { useState, useEffect } from 'react';


 



  const Timeline = () => {
    const [dateFormat, setDateFormat] = useState('weeks');
  


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

    
    const handleFormatChange = (event) => {
      setDateFormat(event.target.value);
    };
  
    return (
      <div className="timeline">
        <div className="format-selector">
          <label htmlFor="date-format">Select Date Format:</label>
          <select id="date-format" value={dateFormat} onChange={handleFormatChange}>
            <option value="weeks">Weeks</option>
            <option value="months">Months</option>
          </select>
        </div>
        <table className="timeline-table">
          <thead>
            <tr>
              <th>List Name</th>
              <th>Date Range</th>
            </tr>
          </thead>
          <tbody>
            {lists.map((list) => (
              <tr key={list.id} className="list-timeline">
                <td className="list-name">{list.name}</td>
                <td className="dates-timeline">
                  <div className="horizontal-scroll">
                    {list.cards.map((card) => (
                      <div
                        key={card.id}
                        className="card-timeline"
                        style={{
                          width: calculateWidth(card.dates[0], card.dates[1], dateFormat),
                        }}
                      >
                        <div className="date-range">
                          {formatDateRange(card.dates[0], card.dates[1], dateFormat)}
                        </div>
                        <div className="card-title">{card.name}</div>
                      </div>
                    ))}
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    );
  };
  
  const formatDate = (date) => {
    const options = { year: 'numeric', month: 'long', day: 'numeric' };
    return new Date(date).toLocaleDateString(undefined, options);
  };
  
  const formatDateRange = (startDate, endDate, format) => {
    const startDateTime = new Date(startDate).getTime();
    const endDateTime = new Date(endDate).getTime();
    const totalDays = Math.ceil((endDateTime - startDateTime) / (1000 * 60 * 60 * 24));
  
    if (format === 'weeks') {
      const daysOfWeek = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
      return Array.from({ length: totalDays }, (_, index) => {
        const day = new Date(startDateTime + index * 24 * 60 * 60 * 1000).getDay();
        return daysOfWeek[day];
      }).join(' ');
    } else if (format === 'months') {
      return Array.from({ length: totalDays }, (_, index) => {
        const day = new Date(startDateTime + index * 24 * 60 * 60 * 1000).getDate();
        return day;
      }).join(' ');
    }
  
    return '';
  };
  
  const calculateWidth = (startDate, endDate, format) => {
    const startDateTime = new Date(startDate).getTime();
    const endDateTime = new Date(endDate).getTime();
    const totalDays = Math.ceil((endDateTime - startDateTime) / (1000 * 60 * 60 * 24));
  
    if (format === 'weeks') {
      return totalDays * 30; // Adjust the factor based on your design for weeks
    } else if (format === 'months') {
      return totalDays * 15; // Adjust the factor based on your design for months
    }
  
    return totalDays * 30; // Default factor for other formats
  };
  
  export default Timeline;
  