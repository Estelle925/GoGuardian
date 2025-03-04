import { BrowserRouter as Router } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import Layout from './layouts/MainLayout';

const App = () => {
  return (
    <ConfigProvider locale={zhCN}>
      <Router>
        <Layout />
      </Router>
    </ConfigProvider>
  );
};

export default App;