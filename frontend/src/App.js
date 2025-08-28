import React, { useEffect, useState } from 'react';
import { Routes, Route, Navigate, useLocation, useNavigate } from 'react-router-dom';
import { useQuery } from 'react-query';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import { API_BASE_URL } from './config';
import Login from './components/Login';
import Dashboard from './components/Dashboard';
import Navbar from './components/Navbar';
import LoadingSpinner from './components/LoadingSpinner';

// Protected route component
const ProtectedRoute = ({ children }) => {
    const { isAuthenticated, isLoading } = useAuth();

    if (isLoading) {
        return <LoadingSpinner />;
    }

    if (!isAuthenticated) {
        return <Navigate to="/login" replace />;
    }

    return children;
};

// OAuth callback component
const OAuthCallback = () => {
    const { handleCallback } = useAuth();
    const location = useLocation();
    const [isProcessing, setIsProcessing] = useState(false);
    const navigate = useNavigate();

    useEffect(() => {
        const urlParams = new URLSearchParams(location.search);
        const code = urlParams.get('code');
        const error = urlParams.get('error');

        if (error) {
            console.error('OAuth error:', error);
            navigate('/login');
            return;
        }

        if (code && !isProcessing) {
            setIsProcessing(true);
            console.log('Processing OAuth callback with code:', code.substring(0, 10) + '...');

            handleCallback(code)
                .then(() => {
                    console.log('OAuth callback successful');
                    // Navigate to dashboard will happen in handleCallback
                })
                .catch((error) => {
                    console.error('OAuth callback failed:', error);
                    navigate('/login');
                })
                .finally(() => {
                    setIsProcessing(false);
                });
        }
    }, [location.search, handleCallback, isProcessing, navigate]);

    return <LoadingSpinner />;
};

// Main app content
const AppContent = () => {
    const { isAuthenticated } = useAuth();

    return (
        <div className="min-h-screen bg-gray-50">
            {isAuthenticated && <Navbar />}
            <main className="container mx-auto px-4 py-8">
                <Routes>
                    <Route path="/login" element={<Login />} />
                    <Route path="/auth/callback" element={<OAuthCallback />} />
                    <Route
                        path="/dashboard"
                        element={
                            <ProtectedRoute>
                                <Dashboard />
                            </ProtectedRoute>
                        }
                    />
                    <Route
                        path="/"
                        element={
                            isAuthenticated ? <Navigate to="/dashboard" replace /> : <Navigate to="/login" replace />
                        }
                    />
                </Routes>
            </main>
        </div>
    );
};

// Main App component
const App = () => {
    return (
        <AuthProvider>
            <AppContent />
        </AuthProvider>
    );
};

export default App;
