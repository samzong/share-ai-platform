import React from 'react';
import { Layout, Menu, Button, Space } from 'antd';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import {
  HomeOutlined,
  AppstoreOutlined,
  UserOutlined,
  LoginOutlined,
  UserAddOutlined,
} from '@ant-design/icons';
import { getToken } from '../../services/auth';

const { Header, Content, Footer } = Layout;

const MainLayout: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const token = getToken();

  const menuItems = [
    {
      key: '/',
      icon: <HomeOutlined />,
      label: '首页',
    },
    {
      key: 'images',
      icon: <AppstoreOutlined />,
      label: '镜像',
    },
    {
      key: 'algorithms',
      icon: <AppstoreOutlined />,
      label: '算法',
      disabled: true,
    },
    {
      key: 'models',
      icon: <AppstoreOutlined />,
      label: '模型',
      disabled: true,
    },
    {
      key: 'datasets',
      icon: <AppstoreOutlined />,
      label: '数据集',
      disabled: true,
    },
  ];

  return (
    <Layout>
      <Header style={{ 
        position: 'fixed', 
        zIndex: 1, 
        width: '100%', 
        padding: '0 24px',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center'
      }}>
        <div style={{ display: 'flex', alignItems: 'center' }}>
          <div style={{ 
            width: 120, 
            height: 31, 
            margin: '16px 24px 16px 0', 
            background: 'rgba(255, 255, 255, 0.2)' 
          }} />
          <Menu
            theme="dark"
            mode="horizontal"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={({ key }) => navigate(key)}
            style={{ lineHeight: '64px', border: 'none' }}
          />
        </div>
        
        <Space>
          {!token ? (
            <>
              <Button 
                type="link" 
                icon={<LoginOutlined />} 
                onClick={() => navigate('/login')}
                style={{ color: '#fff' }}
              >
                登录
              </Button>
              <Button 
                type="primary" 
                icon={<UserAddOutlined />} 
                onClick={() => navigate('/register')}
              >
                注册
              </Button>
            </>
          ) : (
            <Button 
              type="link" 
              icon={<UserOutlined />} 
              onClick={() => navigate('/profile')}
              style={{ color: '#fff' }}
            >
              个人中心
            </Button>
          )}
        </Space>
      </Header>
      <Content style={{ padding: '0 50px', marginTop: 64 }}>
        <div style={{ padding: 24, minHeight: 380 }}>
          <Outlet />
        </div>
      </Content>
      <Footer style={{ textAlign: 'center' }}>
        Share AI Platform ©{new Date().getFullYear()} Created by Share AI Team
      </Footer>
    </Layout>
  );
};

export default MainLayout; 