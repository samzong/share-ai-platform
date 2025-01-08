import React from 'react';
import { Layout, Menu } from 'antd';
import { Outlet, useNavigate } from 'react-router-dom';
import {
  HomeOutlined,
  AppstoreOutlined,
  UserOutlined,
} from '@ant-design/icons';

const { Header, Content, Footer } = Layout;

const MainLayout: React.FC = () => {
  const navigate = useNavigate();

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
      <Header style={{ position: 'fixed', zIndex: 1, width: '100%', padding: 0 }}>
        <div style={{ float: 'left', width: 120, height: 31, margin: '16px 24px 16px 0', background: 'rgba(255, 255, 255, 0.2)' }} />
        <Menu
          theme="dark"
          mode="horizontal"
          defaultSelectedKeys={['/']}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
          style={{ lineHeight: '64px' }}
        />
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