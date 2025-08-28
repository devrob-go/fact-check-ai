import React, { createContext, useContext, useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { API_ENDPOINTS } from '../config';
import { apiClient } from '../services/apiClient';

const AuthContext = createContext();

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [token, setToken] = useState(localStorage.getItem('token'));
    const [isLoading, setIsLoading] = useState(true);
    const navigate = useNavigate();

    // Track ongoing requests to prevent cancellation
    const abortControllers = useRef(new Set());

    const isAuthenticated = !!token && !!user;

    useEffect(() => {
        const initializeAuth = async () => {
            if (token) {
                try {
                    console.log('AuthContext: Validating token, current headers:', apiClient.defaults.headers.common);
                    // Verify token and get user info
                    const userData = await apiClient.get('/api/v1/auth/me');
                    setUser(userData);
                } catch (error) {
                    console.warn('Token validation failed, clearing invalid token:', error.message);
                    // Clear invalid token without calling logout endpoint
                    setToken(null);
                    setUser(null);
                    localStorage.removeItem('token');
                    delete apiClient.defaults.headers.common['Authorization'];
                }
            }
            setIsLoading(false);
        };

        initializeAuth();
    }, [token]);

    // Cleanup function to abort ongoing requests
    useEffect(() => {
        return () => {
            // Abort all ongoing requests when component unmounts
            abortControllers.current.forEach(controller => controller.abort());
            abortControllers.current.clear();
        };
    }, []);

    const login = async () => {
        try {
            const response = await apiClient.get(API_ENDPOINTS.LOGIN);
            const { auth_url } = response;

            // Redirect to Google OAuth
            window.location.href = auth_url;
        } catch (error) {
            console.error('Login failed:', error);
            throw error;
        }
    };

    const handleCallback = async (code) => {
        try {
            console.log('AuthContext: Starting OAuth callback processing');

            const response = await apiClient.get(`${API_ENDPOINTS.CALLBACK}?code=${code}`);
            const { token: newToken, user: userData } = response;

            console.log('AuthContext: Received token and user data, setting authentication state');

            setToken(newToken);
            setUser(userData);
            localStorage.setItem('token', newToken);

            // Set default authorization header
            apiClient.defaults.headers.common['Authorization'] = `Bearer ${newToken}`;

            console.log('AuthContext: Authorization header set:', apiClient.defaults.headers.common['Authorization']);

            console.log('AuthContext: Authentication state set, navigating to dashboard');
            navigate('/dashboard');

            return Promise.resolve();
        } catch (error) {
            console.error('AuthContext: OAuth callback failed:', error);
            throw error;
        }
    };

    const logout = async () => {
        try {
            // Only try to call logout endpoint if we have a valid token
            if (token && apiClient.defaults.headers.common['Authorization']) {
                await apiClient.post(API_ENDPOINTS.LOGOUT);
            }
        } catch (error) {
            // If logout fails (e.g., 401), just log it but continue with cleanup
            console.warn('Backend logout failed, continuing with frontend cleanup:', error.message);
        } finally {
            // Always clean up frontend state regardless of backend response
            setToken(null);
            setUser(null);
            localStorage.removeItem('token');
            delete apiClient.defaults.headers.common['Authorization'];
            navigate('/login');
        }
    };

    const value = {
        user,
        token,
        isAuthenticated,
        isLoading,
        login,
        handleCallback,
        logout,
    };

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
};
