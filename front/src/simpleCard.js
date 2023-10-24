import React, { useState } from 'react';

// import './css/simplecard.css'


const SimpleCard = ({ card, list }) => {
  

  return (
    <div>
      
        <div className="modal">
          <div className="modal-content">
            <h2>{card.name}</h2>
            <h3> لیست : {list.name}</h3>
            <p>ادمین : {card.owner.name}</p>


            <h4 className='members'>اعضا</h4>
                    {card.members && card.members.length > 0 ? (
                                card.members.map((member, index) => {
                                    return (
                                    <div key={index}>
                                         <span >{member.name}</span>
                                    </div>
                                    )
                                })

                            ) : (
                                <span>این کارت عضوی ندارد</span>
                            )}

          </div>
        </div>
      
    </div>
  );
};

export default SimpleCard;
