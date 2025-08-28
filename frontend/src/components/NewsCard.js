import React, { useState } from 'react';
import { Link, Image, Clock, CheckCircle, XCircle, Play } from 'lucide-react';

const NewsCard = ({
    news,
    onVerify,
    isVerifying,
    getStatusIcon,
    getStatusColor,
    getStatusText
}) => {
    const [showExplanation, setShowExplanation] = useState(false);

    const formatDate = (dateString) => {
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    const truncateText = (text, maxLength = 150) => {
        if (text.length <= maxLength) return text;
        return text.substring(0, maxLength) + '...';
    };

    const handleVerify = () => {
        onVerify(news.id);
    };

    return (
        <div className="card hover:shadow-lg transition-shadow duration-200">
            {/* Header */}
            <div className="flex items-start justify-between mb-4">
                <div className="flex items-center space-x-3">
                    {getStatusIcon(news.status)}
                    <div>
                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(news.status)}`}>
                            {getStatusText(news.status)}
                        </span>
                        <p className="text-xs text-gray-500 mt-1">
                            Submitted {formatDate(news.created_at)}
                        </p>
                    </div>
                </div>

                {news.status === 'pending' && (
                    <button
                        onClick={handleVerify}
                        disabled={isVerifying}
                        className="btn-primary flex items-center space-x-2 text-sm"
                    >
                        {isVerifying ? (
                            <>
                                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                                <span>Verifying...</span>
                            </>
                        ) : (
                            <>
                                <Play className="h-4 w-4" />
                                <span>Verify Now</span>
                            </>
                        )}
                    </button>
                )}
            </div>

            {/* Content */}
            <div className="mb-4">
                <p className="text-gray-900 leading-relaxed">
                    {truncateText(news.content)}
                </p>
                {news.content.length > 150 && (
                    <button
                        onClick={() => setShowExplanation(!showExplanation)}
                        className="text-primary-600 hover:text-primary-700 text-sm font-medium mt-2"
                    >
                        {showExplanation ? 'Show less' : 'Read more'}
                    </button>
                )}

                {showExplanation && (
                    <div className="mt-3 p-3 bg-gray-50 rounded-lg">
                        <p className="text-gray-700 text-sm leading-relaxed">{news.content}</p>
                    </div>
                )}
            </div>

            {/* Links and Media */}
            <div className="space-y-2 mb-4">
                {news.link && (
                    <div className="flex items-center space-x-2 text-sm">
                        <Link className="h-4 w-4 text-gray-400" />
                        <a
                            href={news.link}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-primary-600 hover:text-primary-700 truncate"
                        >
                            {news.link}
                        </a>
                    </div>
                )}

                {news.photo_url && (
                    <div className="flex items-center space-x-2 text-sm">
                        <Image className="h-4 w-4 text-gray-400" />
                        <a
                            href={news.photo_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-primary-600 hover:text-primary-700 truncate"
                        >
                            View Image
                        </a>
                    </div>
                )}
            </div>

            {/* AI Explanation */}
            {news.explanation && news.status !== 'pending' && (
                <div className="border-t border-gray-200 pt-4">
                    <div className="flex items-center space-x-2 mb-2">
                        <div className="w-2 h-2 rounded-full bg-primary-500"></div>
                        <span className="text-sm font-medium text-gray-700">AI Analysis</span>
                    </div>
                    <div className="bg-blue-50 border border-blue-200 rounded-lg p-3">
                        <p className="text-sm text-gray-700 leading-relaxed">
                            {news.explanation}
                        </p>
                    </div>
                </div>
            )}

            {/* Footer */}
            <div className="flex items-center justify-between text-xs text-gray-500 mt-4 pt-4 border-t border-gray-200">
                <span>ID: {news.id.substring(0, 8)}...</span>
                {news.updated_at !== news.created_at && (
                    <span>Updated {formatDate(news.updated_at)}</span>
                )}
            </div>
        </div>
    );
};

export default NewsCard;
