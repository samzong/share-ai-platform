import api from './api';
import { User } from '../types/user';

// 创建一个事件总线来处理用户状态变化
const userStateCallbacks: ((user: User | null) => void)[] = [];

export const subscribeToUserState = (callback: (user: User | null) => void) => {
  userStateCallbacks.push(callback);
  return () => {
    const index = userStateCallbacks.indexOf(callback);
    if (index > -1) {
      userStateCallbacks.splice(index, 1);
    }
  };
};

export const notifyUserStateChange = (user: User | null) => {
  console.log('notifyUserStateChange called with user:', user);
  userStateCallbacks.forEach(callback => callback(user));
};

export const login = async (data: { username: string; password: string }) => {
  console.log('Login started');
  const response = await api.post('/v1/auth/login', data);
  console.log('Login response:', response.data);
  if (response.data.token) {
    localStorage.setItem('token', response.data.token);
    // 登录成功后立即获取用户信息
    try {
      console.log('Fetching user profile after login');
      const userResponse = await api.get('/v1/users/profile');
      const userData = userResponse.data;
      console.log('User profile fetched:', userData);
      notifyUserStateChange(userData);
      return {
        token: response.data.token,
        user: userData
      };
    } catch (error) {
      console.error('Failed to fetch user profile after login:', error);
      throw error;
    }
  }
  return response.data;
};

export const register = async (data: { username: string; email: string; password: string }) => {
  const response = await api.post('/v1/auth/register', data);
  return response.data;
};

export const logout = async () => {
  const response = await api.post('/v1/auth/logout');
  localStorage.removeItem('token');
  notifyUserStateChange(null);
  return response.data;
};

export const getProfile = async (): Promise<User> => {
  const token = localStorage.getItem('token');
  if (!token) {
    throw new Error('No token found');
  }
  const response = await api.get('/v1/users/profile');
  return response.data;
};

export const updateProfile = async (data: FormData): Promise<User> => {
  const response = await api.put('/v1/users/profile', data, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  notifyUserStateChange(response.data);
  return response.data;
};
