'use client'

import { useState } from 'react'

interface CodeExampleProps {
  title: string
  code: string
  language?: string
}

export function CodeExample({ title, code, language = 'javascript' }: CodeExampleProps) {
  const [copied, setCopied] = useState(false)
  
  const handleCopy = async () => {
    await navigator.clipboard.writeText(code)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }
  
  return (
    <div className="bg-slate-800 rounded-lg overflow-hidden border border-slate-700">
      <div className="flex items-center justify-between px-4 py-2 bg-slate-900 border-b border-slate-700">
        <span className="text-slate-300 text-sm font-medium">{title}</span>
        <button
          onClick={handleCopy}
          className="text-xs bg-slate-700 hover:bg-slate-600 text-slate-300 px-3 py-1 rounded transition"
        >
          {copied ? '✓ Copied!' : 'Copy'}
        </button>
      </div>
      <pre className="p-4 overflow-x-auto">
        <code className="text-sm text-slate-300">{code}</code>
      </pre>
    </div>
  )
}
