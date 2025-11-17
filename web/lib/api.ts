export const getApiUrl = () => {
  // In browser, always use current domain
  if (typeof window !== 'undefined') {
    return `${window.location.origin}/api`;
  }
  
  // Server-side: use env var or localhost for dev
  if (process.env.NEXT_PUBLIC_API_URL) {
    return process.env.NEXT_PUBLIC_API_URL;
  }
  
  return process.env.NODE_ENV === 'production' 
    ? '/api' // Relative URL for SSR
    : 'http://localhost:8080';
};
