import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { newsService } from '../services/apiClient';
import { Shield, Plus, CheckCircle, XCircle, Clock, Link, Image } from 'lucide-react';
import NewsSubmissionForm from './NewsSubmissionForm';
import NewsCard from './NewsCard';

const Dashboard = () => {
    const { user } = useAuth();
    const [showSubmissionForm, setShowSubmissionForm] = useState(false);
    const queryClient = useQueryClient();

    // Fetch user's news submissions
    const { data: newsData, isLoading, error } = useQuery(
        ['userNews', user?.id],
        () => newsService.getUserNews(user.id),
        {
            enabled: !!user?.id,
            refetchInterval: 30000, // Refetch every 30 seconds
        }
    );

    // Verify news mutation
    const verifyNewsMutation = useMutation(
        (newsId) => newsService.verify(newsId),
        {
            onSuccess: () => {
                queryClient.invalidateQueries(['userNews', user?.id]);
            },
        }
    );

    const handleVerifyNews = async (newsId) => {
        try {
            await verifyNewsMutation.mutateAsync(newsId);
        } catch (error) {
            console.error('Failed to verify news:', error);
        }
    };

    const getStatusIcon = (status) => {
        switch (status) {
            case 'true':
                return <CheckCircle className="h-5 w-5 text-success-500" />;
            case 'false':
                return <XCircle className="h-5 w-5 text-danger-500" />;
            default:
                return <Clock className="h-5 w-5 text-gray-400" />;
        }
    };

    const getStatusColor = (status) => {
        switch (status) {
            case 'true':
                return 'bg-success-100 text-success-800';
            case 'false':
                return 'bg-danger-100 text-danger-800';
            default:
                return 'bg-gray-100 text-gray-800';
        }
    };

    const getStatusText = (status) => {
        switch (status) {
            case 'true':
                return 'Verified True';
            case 'false':
                return 'Verified False';
            default:
                return 'Pending Verification';
        }
    };

    if (isLoading) {
        return (
            <div className="flex items-center justify-center min-h-64">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="text-center py-12">
                <div className="text-red-600 mb-4">
                    <XCircle className="h-12 w-12 mx-auto" />
                </div>
                <h3 className="text-lg font-medium text-gray-900 mb-2">Error loading dashboard</h3>
                <p className="text-gray-600">Failed to load your news submissions. Please try again.</p>
            </div>
        );
    }

    const newsList = newsData?.news || [];

    return (
        <div className="space-y-8">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
                    <p className="text-gray-600 mt-2">Welcome back, {user?.name}!</p>
                </div>
                <button
                    onClick={() => setShowSubmissionForm(true)}
                    className="btn-primary flex items-center space-x-2"
                >
                    <Plus className="h-5 w-5" />
                    <span>Submit News</span>
                </button>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="card">
                    <div className="flex items-center">
                        <div className="flex-shrink-0">
                            <Clock className="h-8 w-8 text-gray-400" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-500">Pending</p>
                            <p className="text-2xl font-semibold text-gray-900">
                                {newsList.filter(n => n.status === 'pending').length}
                            </p>
                        </div>
                    </div>
                </div>

                <div className="card">
                    <div className="flex items-center">
                        <div className="flex-shrink-0">
                            <CheckCircle className="h-8 w-8 text-success-500" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-500">Verified True</p>
                            <p className="text-2xl font-semibold text-gray-900">
                                {newsList.filter(n => n.status === 'true').length}
                            </p>
                        </div>
                    </div>
                </div>

                <div className="card">
                    <div className="flex items-center">
                        <div className="flex-shrink-0">
                            <XCircle className="h-8 w-8 text-danger-500" />
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-500">Verified False</p>
                            <p className="text-2xl font-semibold text-gray-900">
                                {newsList.filter(n => n.status === 'false').length}
                            </p>
                        </div>
                    </div>
                </div>
            </div>

            {/* News Submissions */}
            <div className="space-y-6">
                <div className="flex items-center justify-between">
                    <h2 className="text-xl font-semibold text-gray-900">Your News Submissions</h2>
                    <span className="text-sm text-gray-500">{newsList.length} total</span>
                </div>

                {newsList.length === 0 ? (
                    <div className="card text-center py-12">
                        <Shield className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                        <h3 className="text-lg font-medium text-gray-900 mb-2">No news submitted yet</h3>
                        <p className="text-gray-600 mb-4">Start by submitting a news article for verification.</p>
                        <button
                            onClick={() => setShowSubmissionForm(true)}
                            className="btn-primary"
                        >
                            Submit Your First News
                        </button>
                    </div>
                ) : (
                    <div className="space-y-4">
                        {newsList.map((news) => (
                            <NewsCard
                                key={news.id}
                                news={news}
                                onVerify={handleVerifyNews}
                                isVerifying={verifyNewsMutation.isLoading}
                                getStatusIcon={getStatusIcon}
                                getStatusColor={getStatusColor}
                                getStatusText={getStatusText}
                            />
                        ))}
                    </div>
                )}
            </div>

            {/* News Submission Modal */}
            {showSubmissionForm && (
                <NewsSubmissionForm
                    onClose={() => setShowSubmissionForm(false)}
                    onSuccess={() => {
                        setShowSubmissionForm(false);
                        queryClient.invalidateQueries(['userNews', user?.id]);
                    }}
                />
            )}
        </div>
    );
};

export default Dashboard;
