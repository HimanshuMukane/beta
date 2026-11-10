import React from 'react'
import ReactDOM from 'react-dom/client'
import { QueryClientProvider } from '@tanstack/react-query'
import { TooltipProvider } from '@repo/ui/components/ui/tooltip'
import { queryClient } from './lib/queryClient'
import { AppErrorBoundary, initGlobalFrontendErrorListeners } from './components/AppErrorBoundary'
import App from './App'
import '@repo/ui/globals.css'

// Initialize global uncaught error & rejection reporters
initGlobalFrontendErrorListeners()

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <AppErrorBoundary componentName="PhotographyCRMApp">
      <QueryClientProvider client={queryClient}>
        <TooltipProvider delayDuration={350}>
          <App />
        </TooltipProvider>
      </QueryClientProvider>
    </AppErrorBoundary>
  </React.StrictMode>,
)
