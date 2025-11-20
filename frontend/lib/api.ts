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
