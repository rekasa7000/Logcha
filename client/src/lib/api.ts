const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:4000/api/v1';

export interface User {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  role: 'company_admin' | 'trainee';
  phone: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  first_name: string;
  last_name: string;
  email: string;
  password: string;
  phone: string;
  role: 'company_admin' | 'trainee';
}

export interface TimeRecord {
  id: number;
  trainee_id: number;
  date: string;
  am_time_in: string | null;
  am_time_out: string | null;
  pm_time_in: string | null;
  pm_time_out: string | null;
  am_hours: number;
  pm_hours: number;
  total_hours: number;
  status: 'present' | 'half_day_am' | 'half_day_pm' | 'absent';
  notes: string;
}

export interface TimeInRequest {
  trainee_id: number;
  session: 'am' | 'pm';
}

export interface TimeOutRequest {
  trainee_id: number;
  session: 'am' | 'pm';
}

export interface ApiResponse<T = any> {
  message?: string;
  data?: T;
  error?: string;
}

class ApiClient {
  private baseUrl: string;
  private token: string | null = null;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
    this.token = localStorage.getItem('token');
  }

  setToken(token: string | null) {
    this.token = token;
    if (token) {
      localStorage.setItem('token', token);
    } else {
      localStorage.removeItem('token');
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  // Health check
  async health(): Promise<{ status: string }> {
    return this.request('/health');
  }

  // Auth endpoints
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    return this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
  }

  async register(userData: RegisterRequest): Promise<AuthResponse> {
    return this.request('/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
  }

  async me(): Promise<User> {
    return this.request('/me');
  }

  // Time tracking endpoints
  async timeIn(data: TimeInRequest): Promise<ApiResponse<TimeRecord>> {
    return this.request('/time/in', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async timeOut(data: TimeOutRequest): Promise<ApiResponse<TimeRecord>> {
    return this.request('/time/out', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async getTimeRecords(
    traineeId: number,
    startDate?: string,
    endDate?: string
  ): Promise<ApiResponse<TimeRecord[]>> {
    const params = new URLSearchParams();
    if (startDate) params.append('start_date', startDate);
    if (endDate) params.append('end_date', endDate);
    
    const query = params.toString();
    return this.request(`/time/records/${traineeId}${query ? `?${query}` : ''}`);
  }

  async getTodayRecord(traineeId: number): Promise<ApiResponse<TimeRecord>> {
    return this.request(`/time/today/${traineeId}`);
  }
}

export const apiClient = new ApiClient(API_BASE_URL);