import React from 'react';
import './SummaryCard.css';

function SummaryCard({ title, value, type }) {
  return (
    <div className="summary-card">
      <h3>{title}</h3>
      <p className={type ? type : ''}>{value}</p>
    </div>
  );
}

export default SummaryCard;
