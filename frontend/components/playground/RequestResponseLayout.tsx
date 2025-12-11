import { ReactNode } from 'react'

interface RequestResponseLayoutProps {
  requestPanel: ReactNode
  responsePanel: ReactNode
  actionButton: ReactNode
  isLoading?: boolean
}

export function RequestResponseLayout({
  requestPanel,
  responsePanel,
  actionButton,
  isLoading = false,
}: RequestResponseLayoutProps) {
  return (
    <div className="h-full flex flex-col">
      {/* Two Panel Layout */}
      <div className="flex-1 grid grid-cols-1 lg:grid-cols-2 gap-6 min-h-0">
        {/* Left Panel - Request Configuration */}
        <div className="flex flex-col min-h-0 bg-slate-900/50 rounded-lg border border-slate-700">
          {/* Scrollable Content Area */}
          <div className="flex-1 overflow-y-auto p-6">
            <h3 className="text-lg font-semibold text-white mb-4">Request</h3>
            {requestPanel}
          </div>

          {/* Fixed Action Button at Bottom */}
          <div className="border-t border-slate-700 p-4 bg-slate-900/80">
            {actionButton}
          </div>
        </div>

        {/* Right Panel - Response */}
        <div className="flex flex-col min-h-0 bg-slate-900/50 rounded-lg border border-slate-700">
          <div className="flex-1 overflow-y-auto p-4">
            <h3 className="text-base font-semibold text-white mb-3">Response</h3>
            {isLoading ? (
              <div className="flex items-center justify-center h-32">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
              </div>
            ) : (
              responsePanel
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
