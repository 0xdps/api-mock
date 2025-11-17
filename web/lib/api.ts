export const getApiUrl = () => {
  // If explicitly set, use that
  if (process.env.NEXT_PUBLIC_API_URL) {
    return process.env.NEXT_PUBLIC_API_URL;
  }
  
  // In browser, use current domain
  if (typeof window !== 'undefined') {
    return `${window.location.origin}/api`;
  }
  
  // Server-side fallback
  return process.env.NODE_ENV === 'production' 
    ? `https://${process.env.VERCEL_URL || 'apimock02.vercel.app'}/api`
    : 'http://localhost:8080';
};
