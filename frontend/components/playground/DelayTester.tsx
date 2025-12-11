'use client'

import { useState, useEffect, useRef } from 'react'
import { getApiUrl } from '@/lib/api'

const API_URL = getApiUrl()

export function DelayTester() {
  const [delay, setDelay] = useState(2500)
  const [testType, setTestType] = useState<'param' | 'endpoint'>('param')
  const [isRunning, setIsRunning] = useState(false)
  const [progress, setProgress] = useState(0)
  const [results, setResults] = useState<any>(null)
  const [error, setError] = useState<string | null>(null)
  const abortControllerRef = useRef<AbortController | null>(null)
  const startTimeRef = useRef<number>(0)
  const progressIntervalRef = useRef<NodeJS.Timeout | null>(null)

  useEffect(() => {
    return () => {
      if (progressIntervalRef.current) {
        clearInterval(progressIntervalRef.current)
      }
      if (abortControllerRef.current) {
        abortControllerRef.current.abort()
      }
    }
  }, [])

  const handleStart = async () => {
    setIsRunning(true)
    setProgress(0)
    setResults(null)
    setError(null)
    
    abortControllerRef.current = new AbortController()
    startTimeRef.current = performance.now()

    // Update progress bar
    progressIntervalRef.current = setInterval(() => {
      const elapsed = performance.now() - startTimeRef.current
      const currentProgress = Math.min((elapsed / delay) * 100, 100)
      setProgress(currentProgress)
    }, 50)

    try {
      const url = testType === 'param' 
        ? `${API_URL}/echo?delay=${delay}`
        : `${API_URL}/delay/${delay}`

      const res = await fetch(url, { signal: abortControllerRef.current.signal })
      const endTime = performance.now()
      const actualTime = Math.round(endTime - startTimeRef.current)

      if (!res.ok) {
        throw new Error(`HTTP ${res.status}: ${res.statusText}`)
      }

      const data = await res.json()
      
      setResults({
        requested: delay,
        actual: actualTime,
        difference: actualTime - delay,
        data,
      })
      setProgress(100)
    } catch (err: any) {
      if (err.name !== 'AbortError') {
        setError(err.message)
      }
    } finally {
      setIsRunning(false)
      if (progressIntervalRef.current) {
        clearInterval(progressIntervalRef.current)
        progressIntervalRef.current = null
      }
    }
  }

  const handleStop = () => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort()
    }
    if (progressIntervalRef.current) {
      clearInterval(progressIntervalRef.current)
      progressIntervalRef.current = null
    }
    setIsRunning(false)
    setProgress(0)
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-white mb-2">
          Delay & Timeout Testing
        </h2>
        <p className="text-slate-400">
          Test API delays and timeout handling with visual progress indicators
        </p>
      </div>

      {/* Delay Slider */}
      <div>
        <label className="block text-slate-300 mb-2 font-medium text-sm">
          Delay Duration: {delay}ms {delay >= 30000 && '⚠️ Capped at 30s'}
        </label>
        <input
          type="range"
          min="0"
          max="30000"
          step="500"
          value={delay}
          onChange={(e) => setDelay(parseInt(e.target.value))}
          disabled={isRunning}
          className="w-full disabled:opacity-50"
        />
        <div className="flex justify-between text-xs text-slate-500 mt-1">
          <span>0ms</span>
          <span>15s</span>
          <span>30s (max)</span>
        </div>
        {delay >= 30000 && (
          <p className="text-xs text-yellow-400 mt-2">
            ⚠️ Maximum delay is capped at 30 seconds by the server
          </p>
        )}
      </div>

      {/* Test Type */}
      <div>
        <label className="block text-slate-300 mb-2 font-medium text-sm">
          Test Type
        </label>
        <div className="space-y-2">
          <label className="flex items-center gap-3 p-3 bg-slate-900/50 rounded border border-slate-700 cursor-pointer hover:bg-slate-900/70 transition">
            <input
              type="radio"
              checked={testType === 'param'}
              onChange={() => setTestType('param')}
              disabled={isRunning}
              className="text-blue-500"
            />
            <div className="flex-1">
              <span className="text-slate-200 font-medium">Query Parameter</span>
              <p className="text-xs text-slate-400 mt-1">
                <code className="text-blue-400">GET /echo?delay={delay}</code>
              </p>
            </div>
          </label>
          <label className="flex items-center gap-3 p-3 bg-slate-900/50 rounded border border-slate-700 cursor-pointer hover:bg-slate-900/70 transition">
            <input
              type="radio"
              checked={testType === 'endpoint'}
              onChange={() => setTestType('endpoint')}
              disabled={isRunning}
              className="text-blue-500"
            />
            <div className="flex-1">
              <span className="text-slate-200 font-medium">Dedicated Endpoint</span>
              <p className="text-xs text-slate-400 mt-1">
                <code className="text-blue-400">GET /delay/{delay}</code>
              </p>
            </div>
          </label>
        </div>
      </div>

      {/* Control Buttons */}
      <div className="flex gap-3">
        <button
          onClick={handleStart}
          disabled={isRunning}
          className="flex-1 bg-blue-500 hover:bg-blue-600 disabled:bg-slate-700 disabled:cursor-not-allowed text-white px-6 py-3 rounded font-semibold transition"
        >
          {isRunning ? '⏳ Running...' : '▶ Start Request'}
        </button>
        {isRunning && (
          <button
            onClick={handleStop}
            className="px-6 py-3 bg-red-900/30 hover:bg-red-900/50 text-red-400 rounded font-semibold transition"
          >
            ⏹ Stop
          </button>
        )}
      </div>

      {/* Progress Bar */}
      {isRunning && (
        <div className="space-y-2">
          <div className="flex items-center justify-between text-sm">
            <span className="text-slate-300">
              ⏳ Waiting... ({Math.round((progress / 100) * delay)}ms / {delay}ms)
            </span>
            <span className="text-blue-400 font-semibold">{Math.round(progress)}%</span>
          </div>
          <div className="w-full bg-slate-700 rounded-full h-4 overflow-hidden">
            <div
              className="bg-blue-500 h-full transition-all duration-100 ease-linear"
              style={{ width: `${progress}%` }}
            />
          </div>
        </div>
      )}

      {/* Results */}
      {results && (
        <div className="border-t border-slate-700 pt-6 space-y-4">
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div className="bg-slate-900/50 p-4 rounded border border-slate-700">
              <p className="text-xs text-slate-400 mb-1">Requested Delay</p>
              <p className="text-2xl font-bold text-white">{results.requested}ms</p>
            </div>
            <div className="bg-slate-900/50 p-4 rounded border border-slate-700">
              <p className="text-xs text-slate-400 mb-1">Actual Time</p>
              <p className="text-2xl font-bold text-green-400">{results.actual}ms</p>
            </div>
            <div className="bg-slate-900/50 p-4 rounded border border-slate-700">
              <p className="text-xs text-slate-400 mb-1">Difference</p>
              <p className={`text-2xl font-bold ${
                Math.abs(results.difference) < 50 ? 'text-green-400' :
                Math.abs(results.difference) < 200 ? 'text-yellow-400' :
                'text-red-400'
              }`}>
                {results.difference > 0 ? '+' : ''}{results.difference}ms
              </p>
            </div>
          </div>

          <div className="bg-green-900/20 border border-green-500/30 p-4 rounded">
            <div className="flex items-center gap-2 mb-2">
              <span className="text-green-400 font-semibold">✓ Success</span>
              <span className="text-slate-400 text-sm">
                Delay accuracy: {((1 - Math.abs(results.difference) / results.requested) * 100).toFixed(1)}%
              </span>
            </div>
            <p className="text-xs text-slate-400">
              {Math.abs(results.difference) < 50 
                ? '🎯 Excellent accuracy! Delay was within 50ms of requested time.'
                : Math.abs(results.difference) < 200
                ? '✓ Good accuracy. Delay was within 200ms of requested time.'
                : '⚠️ Timing variance detected. This is normal for network requests.'}
            </p>
          </div>

          <div>
            <h3 className="text-lg font-semibold text-white mb-2">Response Data</h3>
            <div className="bg-slate-900 p-4 rounded border border-slate-600 overflow-auto max-h-64">
              <pre className="text-green-400 text-sm">
                <code>{JSON.stringify(results.data, null, 2)}</code>
              </pre>
            </div>
          </div>
        </div>
      )}

      {/* Error */}
      {error && (
        <div className="bg-red-900/20 border border-red-500/50 p-4 rounded">
          <div className="flex items-center gap-2 mb-2">
            <span className="text-red-400 font-semibold">❌ Error</span>
          </div>
          <p className="text-red-300 text-sm">{error}</p>
        </div>
      )}
    </div>
  )
}
