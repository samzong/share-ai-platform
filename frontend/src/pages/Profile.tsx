import React, { useState, useEffect } from 'react';
import { Card, Form, Input, Button, Upload, message, Avatar } from 'antd';
import { UploadOutlined, UserOutlined } from '@ant-design/icons';
import { getProfile, updateProfile } from '../services/userService';
import { User } from '../types/user';
import { RcFile } from 'antd/lib/upload/interface';

const Profile: React.FC = () => {
  const [form] = Form.useForm();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(false);
  const [avatarFile, setAvatarFile] = useState<File | null>(null);

  useEffect(() => {
    fetchProfile();
  }, []);

  const fetchProfile = async () => {
    try {
      const data = await getProfile();
      setUser(data);
      form.setFieldsValue({
        nickname: data.nickname,
      });
    } catch (error) {
      message.error('获取个人资料失败');
    }
  };

  const handleSubmit = async (values: { nickname: string }) => {
    try {
      setLoading(true);
      await updateProfile({
        nickname: values.nickname,
        avatar: avatarFile || undefined,
      });
      message.success('更新成功');
      fetchProfile();
    } catch (error) {
      message.error('更新失败');
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
    setAvatarFile(file);
    return false;
  };

  return (
    <Card title="个人资料" style={{ maxWidth: 600, margin: '0 auto', marginTop: 24 }}>
      <div style={{ textAlign: 'center', marginBottom: 24 }}>
        <Avatar
          size={100}
          icon={<UserOutlined />}
          src={user?.avatar}
        />
      </div>
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
      >
        <Form.Item label="用户名">
          <Input value={user?.username} disabled />
        </Form.Item>
        <Form.Item label="邮箱">
          <Input value={user?.email} disabled />
        </Form.Item>
        <Form.Item
          label="昵称"
          name="nickname"
          rules={[{ required: true, message: '请输入昵称' }]}
        >
          <Input />
        </Form.Item>
        <Form.Item label="头像">
          <Upload
            beforeUpload={beforeUpload}
            maxCount={1}
            showUploadList={false}
          >
            <Button icon={<UploadOutlined />}>选择图片</Button>
          </Upload>
          {avatarFile && <div style={{ marginTop: 8 }}>已选择: {avatarFile.name}</div>}
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" loading={loading} block>
            保存修改
          </Button>
        </Form.Item>
      </Form>
    </Card>
  );
};

export default Profile; 