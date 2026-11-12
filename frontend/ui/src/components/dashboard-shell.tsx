import * as React from "react"
import { AppSidebar } from "./app-sidebar"
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "./ui/breadcrumb"
import { Separator } from "./ui/separator"
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "./ui/sidebar"
import { ModeToggle } from "./mode-toggle"

interface BreadcrumbItemDef {
  label: string
  href?: string
}

interface DashboardShellProps {
  title?: string
  breadcrumb?: string | React.ReactNode | BreadcrumbItemDef[]
  children: React.ReactNode
  sidebar?: React.ReactNode
  headerActions?: React.ReactNode
}

export function DashboardShell({
  title: _title = "Dashboard",
  breadcrumb = "Dashboard",
  children,
  sidebar,
  headerActions,
}: DashboardShellProps) {
  return (
    <SidebarProvider>
      {sidebar || <AppSidebar />}
      <SidebarInset className="h-screen overflow-hidden flex flex-col">
        <header className="sticky top-0 z-40 flex h-10 shrink-0 items-center justify-between gap-2 px-3 border-b border-border bg-card">
          <div className="flex items-center gap-1.5 min-w-0">
            <SidebarTrigger className="-ml-1 h-6 w-6" />
            <Separator orientation="vertical" className="mr-1.5 h-3.5" />
            <Breadcrumb>
              <BreadcrumbList className="text-xs">
                {Array.isArray(breadcrumb) ? (
                  breadcrumb.map((item, idx) => {
                    const isLast = idx === breadcrumb.length - 1
                    return (
                      <React.Fragment key={idx}>
                        <BreadcrumbItem>
                          {isLast || !item.href ? (
                            <BreadcrumbPage>{item.label}</BreadcrumbPage>
                          ) : (
                            <BreadcrumbLink href={item.href}>{item.label}</BreadcrumbLink>
                          )}
                        </BreadcrumbItem>
                        {!isLast && <BreadcrumbSeparator />}
                      </React.Fragment>
                    )
                  })
                ) : React.isValidElement(breadcrumb) ? (
                  breadcrumb
                ) : (
                  <>
                    <BreadcrumbItem className="hidden md:block">
                      <BreadcrumbLink href="/">Home</BreadcrumbLink>
                    </BreadcrumbItem>
                    <BreadcrumbSeparator className="hidden md:block" />
                    <BreadcrumbItem>
                      <BreadcrumbPage>{breadcrumb}</BreadcrumbPage>
                    </BreadcrumbItem>
                  </>
                )}
              </BreadcrumbList>
            </Breadcrumb>
          </div>
          <div className="flex items-center gap-2">
            {headerActions}
            <ModeToggle />
          </div>
        </header>
        <div className="flex flex-1 flex-col min-w-0 gap-2 p-2 overflow-y-auto">
          {children}
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}

