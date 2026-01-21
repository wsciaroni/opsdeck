import axios from 'axios';
import toast from 'react-hot-toast';

export const client = axios.create({
  baseURL: '/api',
  withCredentials: true,
});

client.interceptors.response.use(
  (response) => response,
  (error) => {
    // If the error is 401 from /me endpoint, suppress toast and redirect
    // This prevents "Unauthorized" alerts and reloads during initial auth checks
    if (error.response?.status === 401 && error.config?.url?.endsWith('/me')) {
      return Promise.reject(error);
    }

    const message = error.response?.data?.error || "Something went wrong";
    toast.error(message);

    if (error.response?.status === 401) {
      // Prevent infinite redirect loop if already on login page
      if (window.location.pathname !== '/login') {
        window.location.href = '/login';
      }
    }

    return Promise.reject(error);
  }
);
