import React from 'react';
import { useAuth } from '../contexts/AuthContext';
import { Shield, LogOut, User } from 'lucide-react';

const Navbar = () => {
    const { user, logout } = useAuth();

    const handleLogout = async () => {
        try {
            await logout();
        } catch (error) {
            console.error('Logout failed:', error);
        }
    };

    return (
        <nav className="bg-white shadow-sm border-b border-gray-200">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex justify-between items-center h-16">
                    {/* Logo and Brand */}
                    <div className="flex items-center space-x-3">
                        <div className="flex items-center justify-center w-8 h-8 rounded-lg bg-primary-600">
                            <Shield className="h-5 w-5 text-white" />
                        </div>
                        <span className="text-xl font-bold text-gray-900">Fact-Check</span>
                    </div>

                    {/* Navigation Links */}
                    <div className="hidden md:flex items-center space-x-8">
                        <a
                            href="/dashboard"
                            className="text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium transition-colors duration-200"
                        >
                            Dashboard
                        </a>
                        <a
                            href="/submit"
                            className="text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium transition-colors duration-200"
                        >
                            Submit News
                        </a>
                        <a
                            href="/history"
                            className="text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md text-sm font-medium transition-colors duration-200"
                        >
                            History
                        </a>
                    </div>

                    {/* User Menu */}
                    <div className="flex items-center space-x-4">
                        <div className="flex items-center space-x-3">
                            {user?.picture && (
                                <img
                                    src={user.picture}
                                    alt={user.name}
                                    className="h-8 w-8 rounded-full"
                                />
                            )}
                            <div className="hidden md:block">
                                <p className="text-sm font-medium text-gray-900">{user?.name}</p>
                                <p className="text-xs text-gray-500">{user?.email}</p>
                            </div>
                        </div>

                        <button
                            onClick={handleLogout}
                            className="flex items-center space-x-2 text-gray-700 hover:text-red-600 px-3 py-2 rounded-md text-sm font-medium transition-colors duration-200"
                        >
                            <LogOut className="h-4 w-4" />
                            <span className="hidden md:inline">Logout</span>
                        </button>
                    </div>
                </div>
            </div>
        </nav>
    );
};

export default Navbar;
