import React, { useState, useEffect } from "react";
import { Layout, Menu, Button, Avatar, Dropdown } from "antd";
import { Link, Outlet, useNavigate, useLocation } from "react-router-dom";
import {
  HomeOutlined,
  PictureOutlined,
  UserOutlined,
  LogoutOutlined,
} from "@ant-design/icons";
import {
  logout,
  getProfile,
  subscribeToUserState,
  notifyUserStateChange,
} from "../../services/userService";
import { User } from "../../types/user";

const { Header, Content, Footer } = Layout;

const MainLayout: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    console.log("MainLayout useEffect called");
    // 订阅用户状态变化
    const unsubscribe = subscribeToUserState((newUser) => {
      console.log("User state changed in MainLayout:", newUser);
      setUser(newUser);
    });

    // 如果有token，获取用户信息
    const token = localStorage.getItem("token");
    console.log("Token in MainLayout:", token);
    if (token) {
      getProfile()
        .then((userData) => {
          console.log("Profile fetched in MainLayout:", userData);
          notifyUserStateChange(userData);
        })
        .catch((error) => {
          console.error("Failed to fetch user profile:", error);
          localStorage.removeItem("token");
          notifyUserStateChange(null);
        });
    }

    return () => {
      console.log("MainLayout cleanup");
      unsubscribe();
    };
  }, []);

  const handleLogout = async () => {
    try {
      await logout();
      navigate("/login");
    } catch (error) {
      console.error("Logout failed:", error);
    }
  };

  const userMenu = (
    <Menu>
      <Menu.Item key="profile" icon={<UserOutlined />}>
        <Link to="/profile">个人资料</Link>
      </Menu.Item>
      <Menu.Item key="logout" icon={<LogoutOutlined />} onClick={handleLogout}>
        退出登录
      </Menu.Item>
    </Menu>
  );

  return (
    <Layout>
      <Header
        style={{ display: "flex", alignItems: "center", padding: "0 24px" }}
      >
        <div style={{ flex: 1 }}>
          <Menu
            theme="dark"
            mode="horizontal"
            selectedKeys={[location.pathname]}
          >
            <Menu.Item key="/" icon={<HomeOutlined />}>
              <Link to="/">首页</Link>
            </Menu.Item>
            <Menu.Item key="/images" icon={<PictureOutlined />}>
              <Link to="/images">容器镜像</Link>
            </Menu.Item>
          </Menu>
        </div>
        <div>
          {user ? (
            <Dropdown overlay={userMenu} placement="bottomRight">
              <div style={{ cursor: "pointer" }}>
                <Avatar src={user.avatar} icon={<UserOutlined />} />
                <span style={{ color: "#fff", marginLeft: 8 }}>
                  {user.nickname || user.username}
                </span>
              </div>
            </Dropdown>
          ) : (
            <Button type="primary" onClick={() => navigate("/login")}>
              登录
            </Button>
          )}
        </div>
      </Header>
      <Content style={{ padding: "24px 50px" }}>
        <Outlet />
      </Content>
      <Footer style={{ textAlign: "center" }}>
        Share AI Platform ©{new Date().getFullYear()} Created by Your Company
      </Footer>
    </Layout>
  );
};

export default MainLayout;
