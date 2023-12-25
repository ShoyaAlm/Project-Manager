import React, { useState, useEffect } from 'react';
import Timeline, {
  TimelineMarkers,
  TodayMarker,
  CustomMarker,
  CursorMarker,
  DateHeader,
  SidebarHeader,
  TimelineHeaders,
} from 'react-calendar-timeline';
import 'react-calendar-timeline/lib/Timeline.css';
import moment from 'moment';
import './css/scheduler.css';
import { useParams } from 'react-router-dom/cjs/react-router-dom.min';

const Scheduler = () => {
  const [lists, setLists] = useState([]);
  const [todayMarker, setTodayMarker] = useState(moment().toDate());
  const { boardId } = useParams();

  


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
  
        // Assuming timelineUnit is an array within tableData
        const [timelineUnit] = tableData || [];
  
        if (!timelineUnit || !Array.isArray(timelineUnit)) {
          console.error('Invalid table info structure:', tableData);
          console.error('timelineUnit:', timelineUnit);
          throw new Error('Invalid table info structure');
        }
  
        setLists([timelineUnit]); // Wrap timelineUnit in an array
      } catch (error) {
        console.error('Error getting table info:', error);
      }
    };
  
    fetchData();
  }, [boardId]);
  



  // Wait for lists to be populated before rendering the Timeline
  if (lists.length === 0) {
    return <div>Loading...</div>; // or render a loading indicator
  }


  const getGroupsAndItems = () => {
    const groups = lists.map((list) => ({ id: list.id, title: list.name }));
    const items = lists.reduce((acc, list) => {
      // Check if list.cards exists and is an array, otherwise set it to an empty array
      const cardsArray = Array.isArray(list.cards) ? list.cards : [];
  
      const listItems = cardsArray.map((card) => {
        if (card && card.dates && Array.isArray(card.dates) && card.dates.length >= 2) {
          return {
            id: card.id,
            group: list.id,
            title: card.name,
            start_time: new Date(card.dates[0]),
            end_time: new Date(card.dates[1]),
            canMove: false,
            canResize: false,
          };
        } else {
          console.warn('Card dates are empty:', card);
          return null;
        }
      });
  
      return acc.concat(listItems.filter((item) => item !== null));
    }, []);
  
    console.log('Items:', items);
    console.log('Groups:', groups);
  
    return { groups, items };
  };
  
  
  
  const { groups, items } = getGroupsAndItems();

  console.log(items);
  console.log(groups);
  
  return (
    <div className="scheduler-wrapper">
      <Timeline
        groups={groups}
        items={items}
        defaultTimeStart={moment().add(-1, 'months')}
        defaultTimeEnd={moment().add(1, 'months')}
        canMove={false}
        canResize={false}
      />
      <TimelineMarkers
        items={items}
        groups={groups}
        showCursor={false}
        markers={[{ date: todayMarker, className: 'today-marker' }]}
        onTodayMarkerClick={(date) => setTodayMarker(date)}
      >
        <TodayMarker updateMarker={() => {}} />
        <CustomMarker date={moment().add(2, 'weeks').valueOf()} updateMarker={() => {}} />
        <CursorMarker />
      </TimelineMarkers>
      <TimelineHeaders>
        <SidebarHeader>
          {({ getRootProps }) => <div {...getRootProps()}>Your Sidebar Content Here</div>}
        </SidebarHeader>
        <DateHeader />
      </TimelineHeaders>
    </div>
  );
};

export default Scheduler;

  {/* {lists.map((list) => (
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
            ))} */}