/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  eslint: {
    ignoreDuringBuilds: true,
  },
  env: {
    NEXT_PUBLIC_API_URL: 
      process.env.NEXT_PUBLIC_API_URL || 
      (process.env.NODE_ENV === 'production' 
        ? `https://${process.env.VERCEL_URL || 'apimock02.vercel.app'}/api`
        : 'http://localhost:8080'),
  },
}

module.exports = nextConfig
