import React, { useState } from 'react';
import { useMutation } from 'react-query';
import { newsService } from '../services/apiClient';
import { X, Upload, Link as LinkIcon, Image as ImageIcon } from 'lucide-react';

const NewsSubmissionForm = ({ onClose, onSuccess }) => {
    const [formData, setFormData] = useState({
        content: '',
        link: '',
        photoURL: '',
    });
    const [errors, setErrors] = useState({});

    const submitNewsMutation = useMutation(
        (data) => newsService.submit(data),
        {
            onSuccess: () => {
                onSuccess();
                setFormData({ content: '', link: '', photoURL: '' });
                setErrors({});
            },
            onError: (error) => {
                console.error('Failed to submit news:', error);
                setErrors({ submit: 'Failed to submit news. Please try again.' });
            },
        }
    );

    const handleSubmit = async (e) => {
        e.preventDefault();

        // Validation
        const newErrors = {};
        if (!formData.content.trim()) {
            newErrors.content = 'News content is required';
        }
        if (formData.content.trim().length < 10) {
            newErrors.content = 'News content must be at least 10 characters long';
        }

        if (Object.keys(newErrors).length > 0) {
            setErrors(newErrors);
            return;
        }

        try {
            await submitNewsMutation.mutateAsync(formData);
        } catch (error) {
            // Error is handled in onError callback
        }
    };

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));

        // Clear error when user starts typing
        if (errors[name]) {
            setErrors(prev => ({ ...prev, [name]: '' }));
        }
    };

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
                {/* Header */}
                <div className="flex items-center justify-between p-6 border-b border-gray-200">
                    <h2 className="text-xl font-semibold text-gray-900">Submit News for Verification</h2>
                    <button
                        onClick={onClose}
                        className="text-gray-400 hover:text-gray-600 transition-colors duration-200"
                    >
                        <X className="h-6 w-6" />
                    </button>
                </div>

                {/* Form */}
                <form onSubmit={handleSubmit} className="p-6 space-y-6">
                    {/* News Content */}
                    <div>
                        <label htmlFor="content" className="block text-sm font-medium text-gray-700 mb-2">
                            News Content *
                        </label>
                        <textarea
                            id="content"
                            name="content"
                            rows={6}
                            value={formData.content}
                            onChange={handleInputChange}
                            placeholder="Enter the news content you want to verify..."
                            className={`input-field ${errors.content ? 'border-danger-500' : ''}`}
                        />
                        {errors.content && (
                            <p className="mt-1 text-sm text-danger-600">{errors.content}</p>
                        )}
                        <p className="mt-1 text-sm text-gray-500">
                            Minimum 10 characters. Be as detailed as possible for better verification.
                        </p>
                    </div>

                    {/* Source Link */}
                    <div>
                        <label htmlFor="link" className="block text-sm font-medium text-gray-700 mb-2">
                            <LinkIcon className="h-4 w-4 inline mr-2" />
                            Source Link (Optional)
                        </label>
                        <input
                            type="url"
                            id="link"
                            name="link"
                            value={formData.link}
                            onChange={handleInputChange}
                            placeholder="https://example.com/news-article"
                            className="input-field"
                        />
                        <p className="mt-1 text-sm text-gray-500">
                            Provide a link to the original source if available.
                        </p>
                    </div>

                    {/* Photo URL */}
                    <div>
                        <label htmlFor="photoURL" className="block text-sm font-medium text-gray-700 mb-2">
                            <ImageIcon className="h-4 w-4 inline mr-2" />
                            Photo URL (Optional)
                        </label>
                        <input
                            type="url"
                            id="photoURL"
                            name="photoURL"
                            value={formData.photoURL}
                            onChange={handleInputChange}
                            placeholder="https://example.com/image.jpg"
                            className="input-field"
                        />
                        <p className="mt-1 text-sm text-gray-500">
                            Provide a link to an image related to the news if available.
                        </p>
                    </div>

                    {/* Submit Error */}
                    {errors.submit && (
                        <div className="bg-danger-50 border border-danger-200 rounded-lg p-4">
                            <p className="text-sm text-danger-600">{errors.submit}</p>
                        </div>
                    )}

                    {/* Actions */}
                    <div className="flex items-center justify-end space-x-3 pt-4 border-t border-gray-200">
                        <button
                            type="button"
                            onClick={onClose}
                            className="btn-secondary"
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            disabled={submitNewsMutation.isLoading}
                            className="btn-primary flex items-center space-x-2"
                        >
                            {submitNewsMutation.isLoading ? (
                                <>
                                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                                    <span>Submitting...</span>
                                </>
                            ) : (
                                <>
                                    <Upload className="h-4 w-4" />
                                    <span>Submit for Verification</span>
                                </>
                            )}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default NewsSubmissionForm;
