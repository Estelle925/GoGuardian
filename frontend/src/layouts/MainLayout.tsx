import { useState, useEffect } from 'react';
import { Layout, Menu, Button, theme, Avatar, Dropdown, Space } from 'antd';
import { MenuFoldOutlined, MenuUnfoldOutlined, UserOutlined, TeamOutlined, KeyOutlined, MenuOutlined, ControlOutlined, LogoutOutlined } from '@ant-design/icons';
import { useNavigate, useLocation, Routes, Route, Navigate } from 'react-router-dom';
import UserList from '../pages/users/UserList';
import RoleList from '../pages/roles/RoleList';
import PermissionList from '../pages/permissions/PermissionList';
import MenuList from '../pages/menus/MenuList';
import ButtonList from '../pages/buttons/ButtonList';
import Login from '../pages/login/Login';

const { Header, Sider, Content } = Layout;

const MainLayout = () => {
  const [collapsed, setCollapsed] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const { token: { colorBgContainer, borderRadiusLG } } = theme.useToken();

  // 检查用户是否已登录
  const isAuthenticated = !!localStorage.getItem('token');
  const username = localStorage.getItem('username') || '用户';

  // 如果用户未登录且不在登录页面，重定向到登录页面
  useEffect(() => {
    if (!isAuthenticated && location.pathname !== '/login') {
      navigate('/login');
    }
  }, [isAuthenticated, location.pathname, navigate]);

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    navigate('/login');
  };

  // 如果在登录页面，不显示布局
  if (location.pathname === '/login') {
    return (
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="*" element={<Navigate to="/login" />} />
      </Routes>
    );
  }

  const menuItems = [
    {
      key: '/users',
      icon: <UserOutlined />,
      label: '用户管理',
    },
    {
      key: '/roles',
      icon: <TeamOutlined />,
      label: '角色管理',
    },
    {
      key: '/permissions',
      icon: <KeyOutlined />,
      label: '权限管理',
    },
    {
      key: '/menus',
      icon: <MenuOutlined />,
      label: '菜单管理',
    },
    {
      key: '/buttons',
      icon: <ControlOutlined />,
      label: '按钮管理',
    },
  ];

  const handleMenuClick = (key: string) => {
    navigate(key);
  };

  const userMenuItems = [
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ];

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider trigger={null} collapsible collapsed={collapsed}>
        <div style={{ 
          height: 32, 
          margin: 16, 
          background: 'rgba(255, 255, 255, 0.2)', 
          display: 'flex', 
          alignItems: 'center', 
          justifyContent: 'center',
          color: '#fff',
          fontSize: '16px',
          fontWeight: 'bold'
        }}>
          {!collapsed && '系统管理平台'}
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => handleMenuClick(key)}
        />
      </Sider>
      <Layout>
        <Header style={{ 
          padding: '0 16px', 
          background: colorBgContainer, 
          display: 'flex', 
          alignItems: 'center', 
          justifyContent: 'space-between'
        }}>
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <Button
              type="text"
              icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
              onClick={() => setCollapsed(!collapsed)}
              style={{
                fontSize: '16px',
                width: 64,
                height: 64,
              }}
            />
            <h1 style={{
              margin: 0,
              fontSize: '22px',
              fontWeight: '600',
              background: 'linear-gradient(45deg, #1677ff, #4096ff)',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
              letterSpacing: '1px'
            }}>系统管理平台</h1>
          </div>
          <div>
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <Space style={{
                cursor: 'pointer',
                padding: '4px 8px',
                borderRadius: '4px',
                transition: 'all 0.3s',
                ':hover': {
                  backgroundColor: 'rgba(0, 0, 0, 0.025)'
                }
              }}>
                <Avatar 
                  style={{
                    backgroundColor: '#1677ff',
                    backgroundImage: 'linear-gradient(45deg, #1677ff, #4096ff)',
                    transition: 'all 0.3s'
                  }}
                  icon={<UserOutlined style={{ color: '#fff' }} />}
                />
                <span style={{
                  fontSize: '14px',
                  fontWeight: '500',
                  color: 'rgba(0, 0, 0, 0.88)'
                }}>{username}</span>
              </Space>
            </Dropdown>
          </div>
        </Header>
        <Content
          style={{
            margin: '24px 16px',
            padding: 24,
            background: colorBgContainer,
            borderRadius: borderRadiusLG,
          }}
        >
          <Routes>
            <Route path="/users" element={<UserList />} />
            <Route path="/roles" element={<RoleList />} />
            <Route path="/permissions" element={<PermissionList />} />
            <Route path="/menus" element={<MenuList />} />
            <Route path="/buttons" element={<ButtonList />} />
            <Route path="/" element={<Navigate to="/users" />} />
          </Routes>
        </Content>
      </Layout>
    </Layout>
  );
};

export default MainLayout;