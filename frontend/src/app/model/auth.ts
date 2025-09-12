export interface LoginRequest {
  identifier: string; // email ili username
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface RegisterRequest {
  firstname: string;
  lastname: string;
  username: string;
  email: string;
  password: string;
}

export interface User {
  id: string;
  firstname: string;
  lastname: string;
  username: string;
  email: string;
  is_active: boolean;
  role: string;
}