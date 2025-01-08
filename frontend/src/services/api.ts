import axios from 'axios';
import { message } from 'antd';

const baseURL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error) => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          // 未授权，清除token并跳转到登录页
          localStorage.removeItem('token');
          window.location.href = '/login';
          break;
        case 403:
          message.error('没有权限访问该资源');
          break;
        case 404:
          message.error('请求的资源不存在');
          break;
        case 500:
          message.error('服务器错误，请稍后重试');
          break;
        default:
          message.error(error.response.data.message || '请求失败');
      }
    } else {
      message.error('网络错误，请检查网络连接');
    }
    return Promise.reject(error);
  }
);

// 用户相关接口
export const userApi = {
  register: (data: any) => api.post('/users/register', data),
  login: (data: any) => api.post('/users/login', data),
  getCurrentUser: () => api.get('/users/me'),
  updateProfile: (data: any) => api.put('/users/me', data),
};

// 镜像相关接口
export const imageApi = {
  getImages: (params: any) => api.get('/images', { params }),
  getImageById: (id: string) => api.get(`/images/${id}`),
  collectImage: (id: string) => api.post(`/images/${id}/collect`),
  uncollectImage: (id: string) => api.delete(`/images/${id}/collect`),
};

// 部署相关接口
export const deployApi = {
  getDeployInfo: (imageId: string) => api.get(`/deploy/${imageId}`),
  deploy: (imageId: string, data: any) => api.post(`/deploy/${imageId}`, data),
};

export default api; 