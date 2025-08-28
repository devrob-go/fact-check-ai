import axios from 'axios';

// Create axios instance with default configuration
export const apiClient = axios.create({
    baseURL: process.env.REACT_APP_API_BASE_URL || '',
    timeout: 60000,
    headers: {
        'Content-Type': 'application/json',
        'X-Requested-With': 'XMLHttpRequest',
    },
    withCredentials: false,
});

// Request interceptor to add auth token
apiClient.interceptors.request.use(
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

// Response interceptor to handle common errors
apiClient.interceptors.response.use(
    (response) => {
        return response.data;
    },
    (error) => {
        // Log detailed error information for debugging
        console.error('API Error:', {
            message: error.message,
            code: error.code,
            status: error.response?.status,
            url: error.config?.url,
            method: error.config?.method,
        });

        // Don't redirect for aborted requests or network errors
        if (error.code === 'ECONNABORTED' || error.code === 'ERR_NETWORK') {
            console.warn('Network request aborted or failed, retrying...');
            return Promise.reject(error);
        }

        if (error.response?.status === 401) {
            // Token expired or invalid, redirect to login
            localStorage.removeItem('token');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

// API service functions
export const newsService = {
    // Submit news for verification
    submit: async (newsData) => {
        return apiClient.post('/api/v1/news/submit', newsData);
    },

    // Verify news using AI
    verify: async (newsId) => {
        return apiClient.get(`/api/v1/news/verify/${newsId}`);
    },

    // Get user's news submissions
    getUserNews: async (userId) => {
        return apiClient.get(`/api/v1/news/user/${userId}`);
    },
};

export const authService = {
    // Get login URL
    getLoginUrl: async () => {
        return apiClient.get('/api/v1/auth/login');
    },

    // Handle OAuth callback
    handleCallback: async (code) => {
        return apiClient.get(`/api/v1/auth/callback?code=${code}`);
    },

    // Logout
    logout: async () => {
        return apiClient.post('/api/v1/auth/logout');
    },
};
