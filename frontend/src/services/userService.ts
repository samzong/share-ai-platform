import axios from "axios";
import {
  LoginRequest,
  RegisterRequest,
  UpdateProfileRequest,
  AuthResponse,
  User,
} from "../types/user";

const BASE_URL = process.env.REACT_APP_API_URL || "http://localhost:8080";
const API_URL = `${BASE_URL}/v1`;

// 设置请求拦截器，添加 token
axios.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export const login = async (data: LoginRequest): Promise<AuthResponse> => {
  const response = await axios.post(`${API_URL}/auth/login`, data);
  // 保存 token 到本地存储
  localStorage.setItem("token", response.data.token);
  return response.data;
};

export const register = async (
  data: RegisterRequest
): Promise<AuthResponse> => {
  const response = await axios.post(`${API_URL}/auth/register`, data);
  // 保存 token 到本地存储
  localStorage.setItem("token", response.data.token);
  return response.data;
};

export const logout = async (): Promise<void> => {
  await axios.post(`${API_URL}/auth/logout`);
  // 清除本地存储的 token
  localStorage.removeItem("token");
};

export const getProfile = async (): Promise<User> => {
  const response = await axios.get(`${API_URL}/users/profile`);
  return response.data;
};

export const updateProfile = async (
  data: UpdateProfileRequest
): Promise<User> => {
  const formData = new FormData();
  if (data.nickname) {
    formData.append("nickname", data.nickname);
  }
  if (data.avatar) {
    formData.append("avatar", data.avatar);
  }

  const response = await axios.put(`${API_URL}/users/profile`, formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
  return response.data;
};
