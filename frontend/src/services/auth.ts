import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface UserResponse {
  id: string;
  username: string;
  email: string;
  role: string;
  token: string;
}

export const register = async (data: RegisterRequest): Promise<UserResponse> => {
  const response = await axios.post(`${API_URL}/api/v1/auth/register`, data);
  return response.data;
};

export const setToken = (token: string) => {
  localStorage.setItem('token', token);
};

export const getToken = (): string | null => {
  return localStorage.getItem('token');
};

export const removeToken = () => {
  localStorage.removeItem('token');
}; 