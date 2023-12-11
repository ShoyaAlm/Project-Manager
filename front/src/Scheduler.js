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

const Scheduler = () => {
  const [lists, setLists] = useState([]);
  const [todayMarker, setTodayMarker] = useState(moment().toDate());

  
  useEffect(() => {
  

    console.log('Scheduler component mounted');

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

  const getGroupsAndItems = () => {
    const groups = lists.map((list) => ({ id: list.id, title: list.name }));
    const items = lists.reduce((acc, list) => {
      const listItems = list.cards.map((card) => ({
        id: card.id,
        group: list.id,
        title: card.name,
        start_time: new Date(card.dates[0]),
        end_time: new Date(card.dates[1]),
        canMove: false,
        canResize: false,
      }));
      return acc.concat(listItems);
    }, []);

    return { groups, items };
  };

  const { groups, items } = getGroupsAndItems();

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