import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';
import SummaryCard from './components/SummaryCard';
import DataTable from './components/DataTable';
import FilterBar from './components/FilterBar';

function App() {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [prefix, setPrefix] = useState('*');
  const [stats, setStats] = useState({
    total: 0,
    long: 0,
    short: 0
  });
  // 检查访问权限
  const [authorized, setAuthorized] = useState(false);

  useEffect(() => {
    // 检查 URL 查询参数中是否包含正确的访问密钥
    const params = new URLSearchParams(window.location.search);
    const accessKey = params.get('access_key');
    // 在这里设置您的秘密访问密钥
    const secretKey = 'V8XoMgzeGY8we47AUlDYMxsvOwDy7SMemWVZV0zSmOWyO6CJmYM9EvlUS4LQpJZk'; // 建议使用更复杂的字符串

    if (accessKey === secretKey) {
      setAuthorized(true);
    }
  }, []);

  const fetchData = async (searchPrefix) => {
    setLoading(true);
    setError(null);

    try {
      const response = await axios.get(`/api/monitor?prefix=${searchPrefix}`);
      setData(response.data.data || []);

      // 计算统计数据
      const total = response.data.data ? response.data.data.length : 0;
      const longCount = response.data.data ? response.data.data.filter(item => item.pair_monitor_data.direct === 'long').length : 0;
      const shortCount = response.data.data ? response.data.data.filter(item => item.pair_monitor_data.direct === 'short').length : 0;

      setStats({
        total,
        long: longCount,
        short: shortCount
      });
    } catch (err) {
      setError(err.message || '获取数据失败');
      console.error('获取数据出错:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (authorized) {
      fetchData(prefix);
    }
  }, [authorized]);

  const handleFilterChange = (newPrefix) => {
    setPrefix(newPrefix);
    fetchData(newPrefix);
  };

  const handleRefresh = () => {
    fetchData(prefix);
  };

  // 如果未授权，显示 404 页面
  if (!authorized) {
    return (
      <div className="not-found">
        <h1>404</h1>
        <p>页面未找到</p>
      </div>
    );
  }

  return (
    <div className="app-container">
      <header className="app-header">
        <FilterBar
          prefix={prefix}
          onFilterChange={handleFilterChange}
          onRefresh={handleRefresh}
        />
      </header>

      <main className="app-main">
        <div className="summary-section">
          <SummaryCard title="监控总数" value={stats.total} />
          <SummaryCard title="多单数量" value={stats.long} type="long" />
          <SummaryCard title="空单数量" value={stats.short} type="short" />
        </div>

        <div className="data-section">
          <h2>交易对监控详情</h2>
          <DataTable
            data={data}
            loading={loading}
            error={error}
          />
        </div>
      </main>
    </div>
  );
}

export default App;
