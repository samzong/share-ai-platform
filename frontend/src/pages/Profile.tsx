import React, { useState, useEffect } from 'react';
import { Card, Form, Input, Button, Upload, message, Avatar, Row, Col, Spin } from 'antd';
import { UploadOutlined, UserOutlined, LoadingOutlined } from '@ant-design/icons';
import { getProfile, updateProfile, subscribeToUserState } from '../services/userService';
import { User } from '../types/user';
import { RcFile } from 'antd/lib/upload/interface';
import { useNavigate } from 'react-router-dom';

const Profile: React.FC = () => {
  const [form] = Form.useForm();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(false);
  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const [avatarLoading, setAvatarLoading] = useState(false);
  const [avatarPreview, setAvatarPreview] = useState<string>('');
  const navigate = useNavigate();

  useEffect(() => {
    // 订阅用户状态变化
    const unsubscribe = subscribeToUserState((newUser) => {
      if (!newUser) {
        // 如果用户未登录，重定向到登录页
        navigate('/login');
        return;
      }
      setUser(newUser);
      form.setFieldsValue({
        nickname: newUser.nickname,
      });
      // 如果有头像，设置预览
      if (newUser.avatar) {
        setAvatarPreview(newUser.avatar);
      }
    });

    // 初始化时检查用户状态
    const token = localStorage.getItem('token');
    if (!token) {
      navigate('/login');
      return;
    }

    getProfile().catch(() => {
      navigate('/login');
    });

    return () => {
      unsubscribe();
    };
  }, [form, navigate]);

  const handleSubmit = async (values: { nickname: string }) => {
    try {
      setLoading(true);
      const formData = new FormData();
      formData.append('nickname', values.nickname);
      if (avatarFile) {
        formData.append('avatar', avatarFile);
      }
      await updateProfile(formData);
      message.success('个人资料更新成功');
    } catch (error) {
      message.error('更新失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  const beforeUpload = (file: RcFile) => {
    const isImage = file.type.startsWith('image/');
    if (!isImage) {
      message.error('只能上传图片文件！');
      return false;
    }
    const isLt2M = file.size / 1024 / 1024 < 2;
    if (!isLt2M) {
      message.error('图片必须小于 2MB！');
      return false;
    }

    // 创建预览URL
    const previewUrl = URL.createObjectURL(file);
    setAvatarPreview(previewUrl);
    setAvatarFile(file);

    // 清理预览URL
    return false;
  };

  const handleAvatarChange = (info: any) => {
    if (info.file.status === 'uploading') {
      setAvatarLoading(true);
      return;
    }
    if (info.file.status === 'done') {
      setAvatarLoading(false);
    }
  };

  if (!user) {
    return <Spin size="large" />;
  }

  return (
    <Card title="个人资料" style={{ maxWidth: 800, margin: '0 auto', marginTop: 24 }}>
      <Row gutter={24}>
        <Col span={8}>
          <div style={{ textAlign: 'center' }}>
            <div style={{ marginBottom: 16 }}>
              <Avatar
                size={120}
                icon={<UserOutlined />}
                src={avatarPreview || user.avatar}
                style={{ border: '1px solid #f0f0f0' }}
              />
            </div>
            <Upload
              name="avatar"
              showUploadList={false}
              beforeUpload={beforeUpload}
              onChange={handleAvatarChange}
            >
              <Button 
                icon={avatarLoading ? <LoadingOutlined /> : <UploadOutlined />}
                disabled={avatarLoading}
              >
                {avatarLoading ? '上传中...' : '更换头像'}
              </Button>
            </Upload>
            <div style={{ marginTop: 8, color: '#888', fontSize: 12 }}>
              支持 JPG、PNG 格式，文件小于 2MB
            </div>
          </div>
        </Col>
        <Col span={16}>
          <Form
            form={form}
            layout="vertical"
            onFinish={handleSubmit}
            initialValues={{
              nickname: user.nickname,
            }}
          >
            <Form.Item label="用户名">
              <Input value={user.username} disabled />
            </Form.Item>
            <Form.Item label="邮箱">
              <Input value={user.email} disabled />
            </Form.Item>
            <Form.Item
              label="昵称"
              name="nickname"
              rules={[
                { required: true, message: '请输入昵称' },
                { min: 2, message: '昵称至少2个字符' },
                { max: 20, message: '昵称最多20个字符' }
              ]}
              extra="昵称将显示在您的个人主页和评论中"
            >
              <Input placeholder="请输入您的昵称" maxLength={20} showCount />
            </Form.Item>
            <Form.Item>
              <Button 
                type="primary" 
                htmlType="submit" 
                loading={loading} 
                block
              >
                保存修改
              </Button>
            </Form.Item>
          </Form>
        </Col>
      </Row>
    </Card>
  );
};

export default Profile; 