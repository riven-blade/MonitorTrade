import React, { useState } from 'react';
import './FilterBar.css';

function FilterBar({ prefix, onFilterChange, onRefresh }) {
  const [inputValue, setInputValue] = useState(prefix);

  const handleChange = (e) => {
    setInputValue(e.target.value);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onFilterChange(inputValue || '*');
  };

  return (
    <div className="filter-bar">
      <form onSubmit={handleSubmit}>
        <div className="filter-controls">
          <label htmlFor="prefix-filter">交易对过滤: </label>
          <input
            type="text"
            id="prefix-filter"
            placeholder="输入交易对前缀..."
            value={inputValue}
            onChange={handleChange}
          />
          <button type="submit" className="filter-btn">过滤</button>
          <button
            type="button"
            className="refresh-btn"
            onClick={onRefresh}
          >
            刷新数据
          </button>
        </div>
      </form>
    </div>
  );
}

export default FilterBar;
