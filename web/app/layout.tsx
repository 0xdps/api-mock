import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'apimock.codes - Schema-Driven Mock API',
  description: 'Free mock API service with realistic data. Schema-driven, instantly available, perfect for prototyping and testing.',
  keywords: ['mock api', 'fake data', 'rest api', 'json api', 'testing', 'prototyping'],
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        {children}
      </body>
    </html>
  )
}
