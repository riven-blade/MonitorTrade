import React, { useState } from 'react';
import './DataTable.css';

function DataTable({ data, loading, error }) {
  // 添加排序状态
  const [sortField, setSortField] = useState('priceDiff');
  const [sortDirection, setSortDirection] = useState('asc');

  // 格式化日期时间
  const formatDateTime = (timestamp) => {
    if (!timestamp) return '-';
    return new Date(timestamp).toLocaleString();
  };

  // 计算差价百分比
  const calculatePriceDiff = (current, monitor) => {
    if (!current || !monitor) return 0;
    return (Math.abs(current - monitor) / monitor * 100).toFixed(2);
  };

  // 处理排序
  const handleSort = (field) => {
    if (sortField === field) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
    } else {
      setSortField(field);
      setSortDirection('asc');
    }
  };

  // 获取排序图标
  const getSortIcon = (field) => {
    if (sortField !== field) return '⇵';
    return sortDirection === 'asc' ? '↑' : '↓';
  };

  // 渲染加载状态
  if (loading) {
    return <div className="table-status">加载中...</div>;
  }

  // 渲染错误状态
  if (error) {
    return <div className="table-error">加载失败: {error}</div>;
  }

  // 渲染空数据状态
  if (!data || data.length === 0) {
    return <div className="table-empty">没有数据</div>;
  }

  // 排序数据
  const sortedData = [...data].sort((a, b) => {
    const aData = a.pair_data;
    const bData = b.pair_data;
    const aMonitorData = a.pair_monitor_data;
    const bMonitorData = b.pair_monitor_data;

    let aValue, bValue;

    switch (sortField) {
      case 'pair':
        aValue = aMonitorData.pair;
        bValue = bMonitorData.pair;
        return sortDirection === 'asc'
          ? aValue.localeCompare(bValue)
          : bValue.localeCompare(aValue);

      case 'priceDiff':
        aValue = calculatePriceDiff(aData.close, aMonitorData.price);
        bValue = calculatePriceDiff(bData.close, bMonitorData.price);
        return sortDirection === 'asc'
          ? parseFloat(aValue) - parseFloat(bValue)
          : parseFloat(bValue) - parseFloat(aValue);

      default:
        aValue = aMonitorData.pair;
        bValue = bMonitorData.pair;
        return aValue.localeCompare(bValue);
    }
  });

  return (
    <div className="data-table-container">
      <table className="data-table">
        <thead>
        <tr>
          <th onClick={() => handleSort('pair')} className="sortable-header">
            交易对 {getSortIcon('pair')}
          </th>
          <th onClick={() => handleSort('priceDiff')} className="sortable-header">
            差价 {getSortIcon('priceDiff')}
          </th>
          <th>方向</th>
          <th>监控价格</th>
          <th>当前卖价</th>
          <th>当前买价</th>
          <th>监控时间</th>
          <th>过期时间(天)</th>
        </tr>
        </thead>
        <tbody>
        {sortedData.map((item, index) => {
          const pairData = item.pair_data;
          const monitorData = item.pair_monitor_data;
            const ttl = item.ttl;
            const priceDiffPercent = calculatePriceDiff(pairData.close, monitorData.price);

            return (
                <tr
                    key={`${monitorData.pair}-${index}`}
                    className={monitorData.direct === 'long' ? 'long' : 'short'}
                >
                  <td>{monitorData.pair}</td>
                  <td className={monitorData.direct === 'long' ? 'positive' : 'negative'}>
                    {priceDiffPercent}%
                  </td>
                  <td>{monitorData.direct === 'long' ? '多' : '空'}</td>
                  <td>{monitorData.price.toFixed(5)}</td>
                  <td>{pairData.ask_price.toFixed(5)}</td>
                  <td>{pairData.bid_price.toFixed(5)}</td>
                  <td>{formatDateTime(monitorData.timestamp)}</td>
                  <td>{(ttl / 3600 / 24).toFixed(1)}</td>
                </tr>
            );
        })}
        </tbody>
      </table>
    </div>
  );
}

export default DataTable;
