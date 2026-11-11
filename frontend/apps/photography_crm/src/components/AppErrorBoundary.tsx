import React from 'react'
import { ErrorBoundary as ReactErrorBoundary, FallbackProps } from 'react-error-boundary'
import { Button } from '@repo/ui/components/ui/button'
import { IconAlertTriangle, IconRotate } from "@repo/ui/icons/tabler-icons"

// Deduplication cache to prevent flooding logs with the same error
const recentErrors = new Set<string>()

export interface FrontendErrorReport {
  message: string
  stack?: string
  url?: string
  line?: number
  column?: number
  component?: string
  errorId?: string
  requestId?: string
  extra?: Record<string, unknown>
}

export async function reportFrontendError(report: FrontendErrorReport): Promise<void> {
  const hashKey = `${report.message}_${report.line}_${report.component}`
  if (recentErrors.has(hashKey)) {
    return
  }
  recentErrors.add(hashKey)

  // Clear from set after 10 seconds to allow recurrent distinct events
  setTimeout(() => recentErrors.delete(hashKey), 10000)

  try {
    await fetch('/api/logs/frontend', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        message: report.message,
        stack: report.stack,
        url: report.url || window.location.href,
        line: report.line,
        column: report.column,
        component: report.component,
        error_id: report.errorId,
        request_id: report.requestId,
        user_agent: navigator.userAgent,
        extra: report.extra,
      }),
    })
  } catch {
    // Fail silently in frontend reporting to prevent crashing client
  }
}

// Global window unhandled error and promise rejection listeners
export function initGlobalFrontendErrorListeners(): void {
  if (typeof window === 'undefined') return

  window.addEventListener('error', (event) => {
    reportFrontendError({
      message: event.message || 'Script error',
      url: event.filename,
      line: event.lineno,
      column: event.colno,
      stack: event.error?.stack,
    })
  })

  window.addEventListener('unhandledrejection', (event) => {
    const reason = event.reason
    reportFrontendError({
      message: reason?.message || String(reason || 'Unhandled Promise Rejection'),
      stack: reason?.stack,
      extra: { type: 'unhandledrejection' },
    })
  })
}

function ErrorFallback({ error, resetErrorBoundary }: FallbackProps) {
  const errorMessage = error instanceof Error ? error.message : typeof error === 'object' && error && 'message' in error ? String((error as any).message) : ''

  return (
    <div className="min-h-[400px] flex items-center justify-center p-6 bg-background">
      <div className="max-w-md w-full p-6 text-center rounded-2xl border border-destructive/20 bg-destructive/5 space-y-4 shadow-sm">
        <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-xl bg-destructive/10 text-destructive">
          <IconAlertTriangle className="h-6 w-6" />
        </div>
        <div className="space-y-1">
          <h2 className="text-base font-bold text-foreground">Something went wrong</h2>
          <p className="text-xs text-muted-foreground">
            An unexpected error occurred in this view. Diagnostics have been recorded.
          </p>
        </div>
        {errorMessage ? (
          <div className="p-3 text-[11px] font-mono text-muted-foreground bg-background rounded-lg border text-left overflow-x-auto max-h-24">
            {errorMessage}
          </div>
        ) : null}
        <div className="pt-2 flex justify-center gap-3">
          <Button
            size="sm"
            onClick={resetErrorBoundary}
            className="text-xs font-semibold gap-1.5 bg-primary text-primary-foreground"
          >
            <IconRotate className="h-3.5 w-3.5" />
            Try Again
          </Button>
          <Button
            size="sm"
            variant="outline"
            onClick={() => window.location.reload()}
            className="text-xs font-semibold"
          >
            Reload App
          </Button>
        </div>
      </div>
    </div>
  )
}

interface AppErrorBoundaryProps {
  children: React.ReactNode
  componentName?: string
}

export function AppErrorBoundary({ children, componentName = 'Root' }: AppErrorBoundaryProps) {
  return (
    <ReactErrorBoundary
      FallbackComponent={ErrorFallback}
      onError={(err: unknown, info) => {
        const errorObj = err instanceof Error ? err : new Error(String(err))
        reportFrontendError({
          message: errorObj.message,
          stack: errorObj.stack,
          component: componentName,
          extra: { componentStack: info.componentStack },
        })
      }}
    >
      {children}
    </ReactErrorBoundary>
  )
}

