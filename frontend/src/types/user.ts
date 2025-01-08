export interface User {
  id: string;
  username: string;
  email: string;
  nickname: string;
  avatar: string;
  role: 'user' | 'admin';
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface UpdateProfileRequest {
  nickname?: string;
  avatar?: File;
}

export interface AuthResponse {
  user: User;
  token: string;
} 