import { HttpClient } from 'pingpong-fetch';

export const getApiUrl = () => {
  // In browser
  if (typeof window !== 'undefined') {
    // Production: use the deployed backend
    if (window.location.hostname !== 'localhost') {
      return 'https://api.mockly.codes';
    }
    // Local dev: use separate API server
    return 'http://localhost:8080';
  }
  
  // Server-side: use env var or defaults
  if (process.env.NEXT_PUBLIC_API_URL) {
    return process.env.NEXT_PUBLIC_API_URL;
  }
  
  return process.env.NODE_ENV === 'production' 
    ? 'https://api.mockly.codes' // Production backend
    : 'http://localhost:8080'; // Local dev API server
};

// Lazy singleton pattern - only create client when first accessed
let _apiClient: HttpClient | null = null;

export function getClient(): HttpClient {
  if (!_apiClient) {
    _apiClient = new HttpClient({
      baseURL: getApiUrl(),
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json'
      },
      retries: 2,
      retryDelay: 1000,
    });
  }
  return _apiClient;
}

// Simple wrapper object with methods
export const apiClient = {
  get: (url: string, options?: any) => getClient().get(url, options),
  post: (url: string, data?: any, options?: any) => getClient().post(url, data, options),
  put: (url: string, data?: any, options?: any) => getClient().put(url, data, options),
  patch: (url: string, data?: any, options?: any) => getClient().patch(url, data, options),
  delete: (url: string, options?: any) => getClient().delete(url, options),
  head: (url: string, options?: any) => getClient().head(url, options),
  options: (url: string, options?: any) => getClient().options(url, options),
};

// Helper function to create a client for server-side rendering
export const createApiClient = (baseURL?: string) => {
  return new HttpClient({
    baseURL: baseURL || getApiUrl(),
    timeout: 30000,
    headers: {
      'Content-Type': 'application/json'
    },
    retries: 2,
    retryDelay: 1000,
  });
};
