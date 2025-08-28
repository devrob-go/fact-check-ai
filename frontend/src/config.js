// API Configuration
export const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || '';

// API Endpoints
export const API_ENDPOINTS = {
    LOGIN: `${API_BASE_URL}/api/v1/auth/login`,
    CALLBACK: `${API_BASE_URL}/api/v1/auth/callback`,
    LOGOUT: `${API_BASE_URL}/api/v1/auth/logout`,
    SUBMIT_NEWS: `${API_BASE_URL}/api/v1/news/submit`,
    VERIFY_NEWS: (id) => `${API_BASE_URL}/api/v1/news/verify/${id}`,
    USER_NEWS: (id) => `${API_BASE_URL}/api/v1/news/user/${id}`,
};

// Google OAuth Configuration
export const GOOGLE_OAUTH_CONFIG = {
    CLIENT_ID: process.env.REACT_APP_GOOGLE_CLIENT_ID || '',
    REDIRECT_URI: process.env.REACT_APP_REDIRECT_URI || 'http://localhost:3000/auth/callback',
    SCOPE: 'email profile',
};

// App Configuration
export const APP_CONFIG = {
    NAME: 'Fact-Check',
    DESCRIPTION: 'AI-powered news verification platform',
    VERSION: '1.0.0',
};
