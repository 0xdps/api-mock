import { redirect } from 'next/navigation'

export default function PlaygroundPage() {
  // Redirect to docs page (playground is now integrated there)
  redirect('/docs')
}
