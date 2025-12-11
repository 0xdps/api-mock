import { Header } from '@/components/Header'

export default function PlaygroundLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 flex flex-col">
      <Header />
      
      <div className="flex-1 overflow-y-auto">
        {children}
      </div>
    </div>
  )
}
