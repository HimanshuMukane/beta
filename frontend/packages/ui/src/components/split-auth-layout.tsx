import React from "react"
import { ModeToggle } from "./mode-toggle"

interface SplitAuthLayoutProps {
  children: React.ReactNode
  imageSlot?: React.ReactNode
  appName?: string
}

export function SplitAuthLayout({
  children,
  imageSlot,
  appName,
}: SplitAuthLayoutProps) {
  return (
    <div className="grid min-h-svh lg:grid-cols-2">
      {/* Left side: Hero / Image Area (hidden on responsive mobile/tablet, shown on lg+) */}
      <div className="relative hidden bg-muted lg:block border-r border-border overflow-hidden">
        {imageSlot ? (
          imageSlot
        ) : (
          <div className="flex h-full w-full flex-col justify-between p-10 bg-zinc-900 text-zinc-100 dark:bg-zinc-950">
            <div className="flex items-center text-lg font-semibold tracking-tight">
              {appName || "Platform Workspace"}
            </div>
            <div className="space-y-2">
              <p className="text-sm text-zinc-400">
                Seamless management, client portals, and secure enterprise workflows.
              </p>
            </div>
          </div>
        )}
      </div>

      {/* Right side: Form Container */}
      <div className="relative flex flex-col justify-between p-6 sm:p-10">
        <div className="flex items-center justify-between w-full">
          <div className="text-base font-semibold tracking-tight lg:hidden">
            {appName || "Platform Workspace"}
          </div>
          <div className="ml-auto">
            <ModeToggle />
          </div>
        </div>

        <div className="flex flex-1 items-center justify-center py-8">
          <div className="w-full max-w-sm sm:max-w-md">{children}</div>
        </div>

        <div className="text-center text-xs text-muted-foreground">
          &copy; {new Date().getFullYear()} All rights reserved.
        </div>
      </div>
    </div>
  )
}

