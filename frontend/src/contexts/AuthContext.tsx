import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import axios from 'axios';
import Cookies from 'js-cookie';

// Types
interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  phone?: string;
  role: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
  lastLogin?: string;
}

interface AuthResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
  expiresAt: string;
}

interface LoginCredentials {
  email: string;
  password: string;
}

interface RegisterCredentials {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  phone?: string;
}

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (credentials: LoginCredentials) => Promise<void>;
  register: (credentials: RegisterCredentials) => Promise<void>;
  logout: () => void;
  refreshToken: () => Promise<void>;
  isAuthenticated: boolean;
  isAdmin: boolean;
}

// Create context
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// API configuration
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// Configure axios defaults
axios.defaults.baseURL = API_BASE_URL;

// Token management
const TOKEN_KEY = 'certitrack_access_token';
const REFRESH_TOKEN_KEY = 'certitrack_refresh_token';

const getAccessToken = () => Cookies.get(TOKEN_KEY);
const getRefreshToken = () => Cookies.get(REFRESH_TOKEN_KEY);

const setTokens = (accessToken: string, refreshToken: string) => {
  Cookies.set(TOKEN_KEY, accessToken, { expires: 1 }); // 1 day
  Cookies.set(REFRESH_TOKEN_KEY, refreshToken, { expires: 7 }); // 7 days
};

const clearTokens = () => {
  Cookies.remove(TOKEN_KEY);
  Cookies.remove(REFRESH_TOKEN_KEY);
};

// Axios interceptors
axios.interceptors.request.use(
  (config) => {
    const token = getAccessToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor for token refresh
axios.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = getRefreshToken();
        if (refreshToken) {
          const response = await axios.post('/auth/refresh', {
            refreshToken,
          });

          const { accessToken, refreshToken: newRefreshToken } = response.data.data;
          setTokens(accessToken, newRefreshToken);

          // Retry original request
          originalRequest.headers.Authorization = `Bearer ${accessToken}`;
          return axios(originalRequest);
        }
      } catch (refreshError) {
        // Refresh failed, clear tokens and redirect to login
        clearTokens();
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

// Auth Provider Component
interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  // Check if user is authenticated on mount
  useEffect(() => {
    const initializeAuth = async () => {
      const token = getAccessToken();
      if (token) {
        try {
          const response = await axios.get('/profile');
          setUser(response.data.data);
        } catch (error) {
          console.error('Failed to fetch user profile:', error);
          clearTokens();
        }
      }
      setLoading(false);
    };

    initializeAuth();
  }, []);

  const login = async (credentials: LoginCredentials): Promise<void> => {
    try {
      const response = await axios.post('/auth/login', credentials);
      const authData: AuthResponse = response.data.data;

      setTokens(authData.accessToken, authData.refreshToken);
      setUser(authData.user);
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Login failed');
    }
  };

  const register = async (credentials: RegisterCredentials): Promise<void> => {
    try {
      const response = await axios.post('/auth/register', credentials);
      const authData: AuthResponse = response.data.data;

      setTokens(authData.accessToken, authData.refreshToken);
      setUser(authData.user);
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Registration failed');
    }
  };

  const logout = () => {
    clearTokens();
    setUser(null);

    // Call logout endpoint (optional, for server-side cleanup)
    axios.post('/auth/logout').catch(() => {
      // Ignore errors on logout endpoint
    });
  };

  const refreshTokenFn = async (): Promise<void> => {
    const refreshToken = getRefreshToken();
    if (!refreshToken) {
      throw new Error('No refresh token available');
    }

    try {
      const response = await axios.post('/auth/refresh', {
        refreshToken,
      });

      const authData: AuthResponse = response.data.data;
      setTokens(authData.accessToken, authData.refreshToken);
      setUser(authData.user);
    } catch (error: any) {
      clearTokens();
      setUser(null);
      throw new Error(error.response?.data?.error || 'Token refresh failed');
    }
  };

  const value: AuthContextType = {
    user,
    loading,
    login,
    register,
    logout,
    refreshToken: refreshTokenFn,
    isAuthenticated: !!user,
    isAdmin: user?.role === 'admin',
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// Custom hook to use auth context
export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

// Higher-order component for protected routes
export const withAuth = <P extends object>(
  WrappedComponent: React.ComponentType<P>
) => {
  const AuthenticatedComponent = (props: P) => {
    const { isAuthenticated, loading } = useAuth();

    if (loading) {
      return (
        <div className="min-h-screen flex items-center justify-center">
          <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-primary-600"></div>
        </div>
      );
    }

    if (!isAuthenticated) {
      if (typeof window !== 'undefined') {
        window.location.href = '/login';
      }
      return null;
    }

    return <WrappedComponent {...props} />;
  };

  AuthenticatedComponent.displayName = `withAuth(${WrappedComponent.displayName || WrappedComponent.name})`;
  return AuthenticatedComponent;
};

// Higher-order component for admin-only routes
export const withAdminAuth = <P extends object>(
  WrappedComponent: React.ComponentType<P>
) => {
  const AdminAuthenticatedComponent = (props: P) => {
    const { isAuthenticated, isAdmin, loading } = useAuth();

    if (loading) {
      return (
        <div className="min-h-screen flex items-center justify-center">
          <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-primary-600"></div>
        </div>
      );
    }

    if (!isAuthenticated) {
      if (typeof window !== 'undefined') {
        window.location.href = '/login';
      }
      return null;
    }

    if (!isAdmin) {
      return (
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-center">
            <h1 className="text-2xl font-bold text-gray-900 mb-4">Access Denied</h1>
            <p className="text-gray-600">You don't have permission to access this page.</p>
          </div>
        </div>
      );
    }

    return <WrappedComponent {...props} />;
  };

  AdminAuthenticatedComponent.displayName = `withAdminAuth(${WrappedComponent.displayName || WrappedComponent.name})`;
  return AdminAuthenticatedComponent;
};