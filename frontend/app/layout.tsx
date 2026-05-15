import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import { UmamiAnalytics } from '@/components/UmamiAnalytics'
import { AuthProvider } from '@/components/AuthProvider'

const inter = Inter({ subsets: ['latin'] })

export const metadata = {
  title: 'Mockly - Schema-Driven Mock API',
  description: 'Free mock API service with realistic data. Schema-driven, instantly available, perfect for prototyping and testing.',
  keywords: ['mock api', 'fake data', 'rest api', 'json api', 'testing', 'prototyping'],
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={inter.className}>
        <AuthProvider>
          {children}
        </AuthProvider>
        <UmamiAnalytics />
      </body>
    </html>
  )
}
