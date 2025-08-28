import React, { useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { Shield, CheckCircle, AlertTriangle } from 'lucide-react';

const Login = () => {
    const { login, handleCallback, isAuthenticated } = useAuth();
    const location = useLocation();
    const navigate = useNavigate();

    useEffect(() => {
        // If already authenticated, redirect to dashboard
        if (isAuthenticated) {
            navigate('/dashboard');
            return;
        }

        // Check for OAuth callback code
        const urlParams = new URLSearchParams(location.search);
        const code = urlParams.get('code');

        if (code) {
            handleCallback(code);
        }
    }, [isAuthenticated, navigate, location.search, handleCallback]);

    const handleGoogleLogin = async () => {
        try {
            await login();
        } catch (error) {
            console.error('Login failed:', error);
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
            <div className="max-w-md w-full space-y-8">
                <div className="text-center">
                    <div className="mx-auto h-16 w-16 flex items-center justify-center rounded-full bg-primary-100">
                        <Shield className="h-8 w-8 text-primary-600" />
                    </div>
                    <h2 className="mt-6 text-3xl font-extrabold text-gray-900">
                        Welcome to Fact-Check
                    </h2>
                    <p className="mt-2 text-sm text-gray-600">
                        AI-powered news verification platform
                    </p>
                </div>

                <div className="mt-8 space-y-6">
                    <div className="card">
                        <div className="space-y-4">
                            <div className="text-center">
                                <h3 className="text-lg font-medium text-gray-900 mb-4">
                                    Sign in to continue
                                </h3>
                                <p className="text-sm text-gray-600 mb-6">
                                    Verify news articles with the power of AI
                                </p>
                            </div>

                            <button
                                onClick={handleGoogleLogin}
                                className="w-full flex items-center justify-center px-4 py-3 border border-gray-300 rounded-lg shadow-sm bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors duration-200"
                            >
                                <svg className="w-5 h-5 mr-2" viewBox="0 0 24 24">
                                    <path
                                        fill="#4285F4"
                                        d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                                    />
                                    <path
                                        fill="#34A853"
                                        d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                                    />
                                    <path
                                        fill="#FBBC05"
                                        d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                                    />
                                    <path
                                        fill="#EA4335"
                                        d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                                    />
                                </svg>
                                Continue with Google
                            </button>
                        </div>
                    </div>

                    {/* Features */}
                    <div className="space-y-4">
                        <div className="flex items-start space-x-3">
                            <CheckCircle className="h-6 w-6 text-success-500 mt-0.5" />
                            <div>
                                <h4 className="text-sm font-medium text-gray-900">AI-Powered Verification</h4>
                                <p className="text-sm text-gray-600">Get instant fact-checking using advanced AI models</p>
                            </div>
                        </div>

                        <div className="flex items-start space-x-3">
                            <CheckCircle className="h-6 w-6 text-success-500 mt-0.5" />
                            <div>
                                <h4 className="text-sm font-medium text-gray-900">Secure Authentication</h4>
                                <p className="text-sm text-gray-600">Safe and secure Google OAuth2 integration</p>
                            </div>
                        </div>

                        <div className="flex items-start space-x-3">
                            <CheckCircle className="h-6 w-6 text-success-500 mt-0.5" />
                            <div>
                                <h4 className="text-sm font-medium text-gray-900">Real-time Results</h4>
                                <p className="text-sm text-gray-600">Get verification results instantly</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Login;
