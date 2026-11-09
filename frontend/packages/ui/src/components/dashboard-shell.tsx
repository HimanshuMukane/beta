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

interface DashboardShellProps {
  title?: string
  breadcrumb?: string
  children: React.ReactNode
}

export function DashboardShell({
  title: _title = "Dashboard",
  breadcrumb = "Dashboard",
  children,
}: DashboardShellProps) {
  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <header className="flex h-10 shrink-0 items-center justify-between gap-2 px-3 border-b bg-card">
          <div className="flex items-center gap-1.5">
            <SidebarTrigger className="-ml-1 h-6 w-6" />
            <Separator orientation="vertical" className="mr-1.5 h-3.5" />
            <Breadcrumb>
              <BreadcrumbList className="text-xs">
                <BreadcrumbItem className="hidden md:block">
                  <BreadcrumbLink href="/">Home</BreadcrumbLink>
                </BreadcrumbItem>
                <BreadcrumbSeparator className="hidden md:block" />
                <BreadcrumbItem>
                  <BreadcrumbPage>{breadcrumb}</BreadcrumbPage>
                </BreadcrumbItem>
              </BreadcrumbList>
            </Breadcrumb>
          </div>
          <div className="flex items-center gap-1.5">
            <ModeToggle />
          </div>
        </header>
        <div className="flex flex-1 flex-col min-w-0 gap-2 p-2">
          {children}
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}

