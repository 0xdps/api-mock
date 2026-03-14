'use client'

import { useState, useEffect } from 'react'
import { 
  Play, 
  Square, 
  RefreshCw, 
  CheckCircle2, 
  AlertCircle, 
  AlertTriangle, 
  Zap, 
  Dices,
  Check,
  X
} from 'lucide-react'
import { getApiUrl, apiClient } from '@/lib/api'
import { trackUmamiEvent } from '@/lib/analytics'
import { RequestResponseLayout } from './RequestResponseLayout'

const API_URL = getApiUrl()

interface RequestResult {
  index: number
  success: boolean
  status?: number
  time: number
  error?: string
}

export function ChaosTester() {
  const [totalRequests, setTotalRequests] = useState(50)
  const [concurrentRequests, setConcurrentRequests] = useState(5)
  const [flakyRate, setFlakyRate] = useState(60)
  const [maxDelay, setMaxDelay] = useState(2000)
  const [resource, setResource] = useState('products')
  
  const [isRunning, setIsRunning] = useState(false)
  const [progress, setProgress] = useState(0)
  const [results, setResults] = useState<RequestResult[]>([])
  const [stats, setStats] = useState({
    completed: 0,
    successful: 0,
    failed: 0,
    avgTime: 0,
    minTime: Infinity,
    maxTime: 0,
  })

  useEffect(() => {
    if (results.length === 0) return
    
    const successful = results.filter(r => r.success).length
    const failed = results.length - successful
    const times = results.map(r => r.time)
    const avgTime = Math.round(times.reduce((a, b) => a + b, 0) / times.length)
    const minTime = Math.min(...times)
    const maxTime = Math.max(...times)
    
    setStats({
      completed: results.length,
      successful,
      failed,
      avgTime,
      minTime,
      maxTime,
    })
  }, [results])

  const handleStart = async () => {
    trackUmamiEvent('c_s', {
      totalRequests,
      concurrentRequests,
      flakyRate,
      maxDelay,
      resource,
    })

    setIsRunning(true)
    setResults([])
    setProgress(0)
    
    const allResults: RequestResult[] = []
    let completed = 0

    // Create batches of concurrent requests
    for (let i = 0; i < totalRequests; i += concurrentRequests) {
      const batchSize = Math.min(concurrentRequests, totalRequests - i)
      const batch = []
      
      for (let j = 0; j < batchSize; j++) {
        const requestIndex = i + j
        const randomDelay = Math.floor(Math.random() * maxDelay)
        
        batch.push(
          (async () => {
            const startTime = performance.now()
            try {
              const url = `${API_URL}/${resource}?count=5&delay=${randomDelay}&flakyRate=${flakyRate / 100}`
              const res = await apiClient.get(url)
              const endTime = performance.now()
              
              return {
                index: requestIndex + 1,
                success: res.ok(),
                status: res.status,
                time: Math.round(endTime - startTime),
              }
            } catch (err: any) {
              const endTime = performance.now()
              return {
                index: requestIndex + 1,
                success: false,
                time: Math.round(endTime - startTime),
                error: err.message,
              }
            }
          })()
        )
      }
      
      const batchResults = await Promise.all(batch)
      allResults.push(...batchResults)
      completed += batchResults.length
      
      setResults([...allResults])
      setProgress((completed / totalRequests) * 100)
    }
    
    const successfulCount = allResults.filter(result => result.success).length
    const failedCount = allResults.length - successfulCount
    const averageTime = allResults.length > 0
      ? Math.round(allResults.reduce((sum, result) => sum + result.time, 0) / allResults.length)
      : 0

    trackUmamiEvent('c_ok', {
      totalRequests,
      completed: allResults.length,
      successful: successfulCount,
      failed: failedCount,
      avgTime: averageTime,
      resource,
    })

    setIsRunning(false)
  }

  const handleStop = () => {
    trackUmamiEvent('c_x', {
      completed: stats.completed,
      progress: Math.round(progress),
      totalRequests,
      resource,
    })

    setIsRunning(false)
  }

  const getTimeBucket = (time: number): number => {
    return Math.floor(time / 500) * 500
  }

  const getHistogramData = () => {
    if (results.length === 0) return []
    
    const buckets: Record<number, number> = {}
    results.forEach(r => {
      const bucket = getTimeBucket(r.time)
      buckets[bucket] = (buckets[bucket] || 0) + 1
    })
    
    return Object.entries(buckets)
      .sort(([a], [b]) => parseInt(a) - parseInt(b))
      .map(([bucket, count]) => ({ bucket: parseInt(bucket), count }))
  }

  const histogramData = getHistogramData()
  const maxCount = Math.max(...histogramData.map(d => d.count), 1)

  const requestPanel = (
    <>
      {/* Configuration */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-slate-300 mb-2 text-sm">
            Total Requests: {totalRequests}
          </label>
          <input
            type="range"
            min="1"
            max="100"
            value={totalRequests}
            onChange={(e) => setTotalRequests(parseInt(e.target.value))}
            disabled={isRunning}
            className="w-full disabled:opacity-50"
          />
        </div>
        
        <div>
          <label className="block text-slate-300 mb-2 text-sm">
            Concurrent: {concurrentRequests}
          </label>
          <input
            type="range"
            min="1"
            max="10"
            value={concurrentRequests}
            onChange={(e) => setConcurrentRequests(parseInt(e.target.value))}
            disabled={isRunning}
            className="w-full disabled:opacity-50"
          />
        </div>
        
        <div>
          <label className="block text-slate-300 mb-2 text-sm">
            Flaky Rate: {flakyRate}% {flakyRate > 50 && <AlertTriangle className="w-3.5 h-3.5 text-amber-500 inline-block ml-1 mb-0.5" />}
          </label>
          <input
            type="range"
            min="0"
            max="100"
            step="5"
            value={flakyRate}
            onChange={(e) => setFlakyRate(parseInt(e.target.value))}
            disabled={isRunning}
            className="w-full disabled:opacity-50"
          />
        </div>
        
        <div>
          <label className="block text-slate-300 mb-2 text-sm">
            Max Random Delay: {maxDelay}ms
          </label>
          <input
            type="range"
            min="0"
            max="5000"
            step="100"
            value={maxDelay}
            onChange={(e) => setMaxDelay(parseInt(e.target.value))}
            disabled={isRunning}
            className="w-full disabled:opacity-50"
          />
        </div>
      </div>

      {/* Resource Selection */}
      <div>
        <label className="block text-slate-300 mb-2 font-medium text-sm">
          Test Resource
        </label>
        <select
          value={resource}
          onChange={(e) => setResource(e.target.value)}
          disabled={isRunning}
          className="w-full bg-slate-700 text-white px-3 py-2 rounded border border-slate-600 focus:border-blue-500 focus:outline-none text-sm disabled:opacity-50"
        >
          <option value="products">Products</option>
          <option value="users">Users</option>
          <option value="orders">Orders</option>
          <option value="articles">Articles</option>
        </select>
      </div>
    </>
  )

  const actionButton = (
    <div className="flex gap-3">
      <button
        onClick={handleStart}
        disabled={isRunning}
        className="flex-1 bg-blue-500 hover:bg-blue-600 disabled:bg-slate-700 disabled:cursor-not-allowed text-white px-6 py-3 rounded-xl font-semibold text-base transition flex items-center justify-center gap-3"
      >
        {isRunning ? <Dices className="w-5 h-5 animate-bounce" /> : <Play className="w-5 h-5 fill-current" />}
        {isRunning ? 'Running Chaos Test...' : 'Start Chaos Test'}
      </button>
      {isRunning && (
        <button
          onClick={handleStop}
          className="px-6 py-3 bg-red-900/30 hover:bg-red-900/50 text-red-400 rounded-xl font-bold transition flex items-center gap-2"
        >
          <Square className="w-5 h-5 fill-current" /> Stop
        </button>
      )}
    </div>
  )

  const responsePanel = (
    <>
      {(isRunning || results.length > 0) ? (
        <div className="space-y-4">
          <div className="space-y-2">
            <div className="flex items-center justify-between text-sm">
              <span className="text-slate-300">
                Progress: {stats.completed}/{totalRequests}
              </span>
              <span className="text-blue-400 font-semibold">{Math.round(progress)}%</span>
            </div>
            <div className="w-full bg-slate-700 rounded-full h-3 overflow-hidden">
              <div
                className="bg-blue-500 h-full transition-all duration-300"
                style={{ width: `${progress}%` }}
              />
            </div>
          </div>

          {/* Live Statistics */}
          <div className="grid grid-cols-2 sm:grid-cols-5 gap-3">
            <div className="bg-slate-900/50 p-3 rounded border border-slate-700">
              <p className="text-xs text-slate-400 mb-1">Completed</p>
              <p className="text-xl font-bold text-white">{stats.completed}</p>
            </div>
            <div className="bg-green-900/20 p-3 rounded border border-green-500/30">
              <p className="text-xs text-green-400 mb-1">Success</p>
              <p className="text-xl font-bold text-green-400">
                {stats.successful} ({stats.completed > 0 ? Math.round((stats.successful / stats.completed) * 100) : 0}%)
              </p>
            </div>
            <div className="bg-red-900/20 p-3 rounded border border-red-500/30">
              <p className="text-xs text-red-400 mb-1">Failed</p>
              <p className="text-xl font-bold text-red-400">
                {stats.failed} ({stats.completed > 0 ? Math.round((stats.failed / stats.completed) * 100) : 0}%)
              </p>
            </div>
            <div className="bg-slate-900/50 p-3 rounded border border-slate-700">
              <p className="text-xs text-slate-400 mb-1">Avg Time</p>
              <p className="text-xl font-bold text-white">{stats.avgTime}ms</p>
            </div>
            <div className="bg-slate-900/50 p-3 rounded border border-slate-700">
              <p className="text-xs text-slate-400 mb-1">Min / Max</p>
              <p className="text-sm font-bold text-white">
                {stats.minTime === Infinity ? '-' : stats.minTime}ms / {stats.maxTime}ms
              </p>
            </div>
          </div>

          {/* Response Time Distribution */}
          <div className="bg-slate-900/50 p-4 rounded border border-slate-700">
            <h3 className="text-sm font-semibold text-white mb-3">Response Time Distribution</h3>
            {isRunning ? (
              <div className="flex items-center justify-center h-32 text-slate-400">
                <div className="text-center">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto mb-2"></div>
                  <p className="text-sm">Collecting data...</p>
                </div>
              </div>
            ) : histogramData.length > 0 ? (
              <div className="flex items-end justify-between gap-2 h-32 px-2">
                {histogramData.map(({ bucket, count }) => {
                  const heightPercent = maxCount > 0 ? (count / maxCount) * 100 : 0
                  return (
                    <div key={bucket} className="flex flex-col items-center gap-2 flex-1 max-w-[60px]">
                      <div className="relative w-full flex items-end justify-center" style={{ height: '100px' }}>
                        <div
                          className="w-full bg-blue-500 hover:bg-blue-400 rounded-t transition-all cursor-pointer"
                          style={{ 
                            height: `${Math.max(heightPercent, 5)}%`,
                          }}
                          title={`${bucket}-${bucket + 500}ms: ${count} requests`}
                        >
                          <div className="text-xs text-white font-semibold text-center pt-1">
                            {count}
                          </div>
                        </div>
                      </div>
                      <span className="text-xs text-slate-400 whitespace-nowrap">
                        {bucket}ms
                      </span>
                    </div>
                  )
                })}
              </div>
            ) : (
              <div className="flex items-center justify-center h-32 text-slate-400">
                <p className="text-sm">No data yet</p>
              </div>
            )}
          </div>

          {/* Individual Results */}
          {results.length > 0 && (
            <div className="bg-slate-900/50 p-4 rounded border border-slate-700">
              <h3 className="text-sm font-semibold text-white mb-3">
                Individual Results {isRunning && <span className="text-blue-400 text-xs ml-2">(Live)</span>}
              </h3>
              <div className="max-h-64 overflow-y-auto space-y-1">
                {results.map((result) => (
                  <div
                    key={result.index}
                    className={`flex items-center justify-between p-2 rounded text-xs ${
                      result.success
                        ? 'bg-green-900/20 text-green-400'
                        : 'bg-red-900/20 text-red-400'
                    }`}
                  >
                    <span>
                      #{result.index} {result.success ? <Check className="w-3 h-3 inline ml-1" /> : <X className="w-3 h-3 inline ml-1" />}
                    </span>
                    <span>
                      {result.status && `${result.status} - `}{result.time}ms
                    </span>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      ) : (
        <div className="text-center py-12 text-slate-400">
          <p>Start a chaos test to see results</p>
        </div>
      )}
    </>
  )

  return (
    <RequestResponseLayout
      requestPanel={requestPanel}
      responsePanel={responsePanel}
      actionButton={actionButton}
    />
  )
}
