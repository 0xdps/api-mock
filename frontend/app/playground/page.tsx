import { Metadata } from 'next'
import { redirect } from 'next/navigation'

export const metadata: Metadata = {
  title: 'Playground - API Mockly',
  description: 'Interactive API testing playground with utilities for echo testing, status codes, delays, middleware, and chaos engineering.',
}

export default function PlaygroundPage() {
  // Redirect to echo by default
  redirect('/playground/echo')
}
