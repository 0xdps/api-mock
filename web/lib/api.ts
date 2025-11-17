export const getApiUrl = () => {
  // In browser
  if (typeof window !== 'undefined') {
    // Production: use current domain
    if (window.location.hostname !== 'localhost') {
      return `${window.location.origin}/api`;
    }
    // Local dev: use separate API server (without /api prefix)
    return 'http://localhost:8080';
  }
  
  // Server-side: use env var or defaults
  if (process.env.NEXT_PUBLIC_API_URL) {
    return process.env.NEXT_PUBLIC_API_URL;
  }
  
  return process.env.NODE_ENV === 'production' 
    ? '/api' // Relative URL for SSR in production
    : 'http://localhost:8080'; // Local dev API server (without /api prefix)
};
