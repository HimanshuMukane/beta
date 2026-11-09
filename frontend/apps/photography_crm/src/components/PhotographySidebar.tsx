import * as React from 'react'
import {
  IconReceipt,
  IconTargetArrow,
  IconShoppingBag,
  IconCalendarEvent,
  IconChartPie,
  IconReportAnalytics,
  IconUserCheck,
  IconMessageHeart,
  IconArrowLeft,
  IconCamera,
  IconFolders,
  IconChevronRight,
  IconSettings,
  IconDatabase,
  IconUsers,
  IconId,
  IconBriefcase,
  IconBox,
  IconPackage,
  IconFileText,
  IconFileInvoice,
  IconFileCheck,
  IconSignature,
  IconMessageDots,
  IconGitFork,
} from '@tabler/icons-react'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  Button,
  cn,
} from '@repo/ui'
import type { SettingsTab } from '../types/settings'

export type ReportCategory =
  | 'billing'
  | 'leads'
  | 'orders'
  | 'sessions'
  | 'finance'
  | 'master-report'
  | 'account-analytics'
  | 'feedback'

interface PhotographySidebarProps extends React.ComponentProps<typeof Sidebar> {
  mode: 'main' | 'reports' | 'settings'
  currentRoute: string
  activeReport: ReportCategory
  activeSetting: SettingsTab
  onSelectRoute: (route: string) => void
  onSelectReport: (report: ReportCategory) => void
  onSelectSetting: (setting: SettingsTab) => void
  onExitSubmode: () => void
}

const REPORT_MENU_ITEMS: {
  id: ReportCategory
  title: string
  icon: React.ComponentType<{ className?: string }>
  badge?: string
}[] = [
  { id: 'billing', title: 'Billing', icon: IconReceipt },
  { id: 'leads', title: 'Leads', icon: IconTargetArrow },
  { id: 'orders', title: 'Orders', icon: IconShoppingBag },
  { id: 'sessions', title: 'Sessions', icon: IconCalendarEvent },
  { id: 'finance', title: 'Finance', icon: IconChartPie, badge: 'Insights' },
  { id: 'master-report', title: 'Master Report', icon: IconReportAnalytics },
  { id: 'account-analytics', title: 'Account Analytics', icon: IconUserCheck },
  { id: 'feedback', title: 'Feedback', icon: IconMessageHeart },
]

export function PhotographySidebar({
  mode,
  currentRoute,
  activeReport,
  activeSetting,
  onSelectRoute,
  onSelectReport,
  onSelectSetting,
  onExitSubmode,
  ...props
}: PhotographySidebarProps) {
  // SETTINGS SIDEBAR MODE
  if (mode === 'settings') {
    return (
      <Sidebar collapsible="icon" {...props}>
        {/* Header with Back button to Studio CRM */}
        <SidebarHeader className="border-b p-3">
          <div className="flex items-center gap-2 min-w-0">
            <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-emerald-600 text-white shadow-xs">
              <IconSettings className="h-4 w-4" />
            </span>
            <div className="truncate">
              <h3 className="text-xs font-bold leading-tight text-foreground">
                Studio Settings
              </h3>
              <p className="text-[10px] text-muted-foreground truncate">
                Presets & Configurations
              </p>
            </div>
          </div>
          <Button
            variant="outline"
            size="sm"
            onClick={onExitSubmode}
            className="w-full mt-2 h-7 text-[11px] font-semibold gap-1.5 border-dashed"
          >
            <IconArrowLeft className="h-3.5 w-3.5" />
            <span>Back to Studio CRM</span>
          </Button>
        </SidebarHeader>

        <SidebarContent>
          {/* GROUP 1: Settings */}
          <SidebarGroup>
            <SidebarGroupLabel className="text-[10px] uppercase font-bold tracking-wider text-muted-foreground">
              Settings
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {[
                  { id: 'master-data' as SettingsTab, title: 'Master Data', icon: IconDatabase },
                  { id: 'photographers' as SettingsTab, title: 'Photographers', icon: IconUsers },
                  { id: 'sequence-id' as SettingsTab, title: 'Sequence ID', icon: IconId },
                ].map((item) => {
                  const Icon = item.icon
                  const isActive = activeSetting === item.id
                  return (
                    <SidebarMenuItem key={item.id}>
                      <SidebarMenuButton
                        isActive={isActive}
                        onClick={() => onSelectSetting(item.id)}
                        className={cn(
                          'h-9 text-xs transition font-medium',
                          isActive &&
                            'bg-primary/10 text-primary font-bold shadow-2xs hover:bg-primary/15'
                        )}
                      >
                        <Icon className="h-4 w-4 shrink-0" />
                        <span className="flex-1 truncate">{item.title}</span>
                        {isActive && (
                          <IconChevronRight className="h-3.5 w-3.5 text-primary shrink-0" />
                        )}
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  )
                })}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>

          {/* GROUP 2: Preset Configuration */}
          <SidebarGroup>
            <SidebarGroupLabel className="text-[10px] uppercase font-bold tracking-wider text-muted-foreground">
              Preset Configuration
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {[
                  { id: 'preset-services' as SettingsTab, title: 'Services', icon: IconBriefcase },
                  { id: 'preset-packages' as SettingsTab, title: 'Packages', icon: IconBox },
                  { id: 'preset-orders' as SettingsTab, title: 'Orders', icon: IconShoppingBag },
                  { id: 'preset-deliverables' as SettingsTab, title: 'Deliverables', icon: IconPackage },
                ].map((item) => {
                  const Icon = item.icon
                  const isActive = activeSetting === item.id
                  return (
                    <SidebarMenuItem key={item.id}>
                      <SidebarMenuButton
                        isActive={isActive}
                        onClick={() => onSelectSetting(item.id)}
                        className={cn(
                          'h-9 text-xs transition font-medium',
                          isActive &&
                            'bg-primary/10 text-primary font-bold shadow-2xs hover:bg-primary/15'
                        )}
                      >
                        <Icon className="h-4 w-4 shrink-0" />
                        <span className="flex-1 truncate">{item.title}</span>
                        {isActive && (
                          <IconChevronRight className="h-3.5 w-3.5 text-primary shrink-0" />
                        )}
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  )
                })}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>

          {/* GROUP 3: Customize Design & Templates */}
          <SidebarGroup>
            <SidebarGroupLabel className="text-[10px] uppercase font-bold tracking-wider text-muted-foreground">
              Customize Design & Templates
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {[
                  { id: 'template-forms' as SettingsTab, title: 'Forms', icon: IconFileText },
                  { id: 'template-quotation-pdf' as SettingsTab, title: 'Quotation PDF', icon: IconFileInvoice },
                  { id: 'template-invoice-pdf' as SettingsTab, title: 'Invoice PDF', icon: IconFileCheck },
                  { id: 'template-contracts' as SettingsTab, title: 'Contracts', icon: IconSignature },
                  { id: 'template-feedbacks' as SettingsTab, title: 'Feedbacks', icon: IconMessageDots },
                ].map((item) => {
                  const Icon = item.icon
                  const isActive = activeSetting === item.id
                  return (
                    <SidebarMenuItem key={item.id}>
                      <SidebarMenuButton
                        isActive={isActive}
                        onClick={() => onSelectSetting(item.id)}
                        className={cn(
                          'h-9 text-xs transition font-medium',
                          isActive &&
                            'bg-primary/10 text-primary font-bold shadow-2xs hover:bg-primary/15'
                        )}
                      >
                        <Icon className="h-4 w-4 shrink-0" />
                        <span className="flex-1 truncate">{item.title}</span>
                        {isActive && (
                          <IconChevronRight className="h-3.5 w-3.5 text-primary shrink-0" />
                        )}
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  )
                })}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>

          {/* GROUP 4: Stages */}
          <SidebarGroup>
            <SidebarGroupLabel className="text-[10px] uppercase font-bold tracking-wider text-muted-foreground">
              Stages
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                <SidebarMenuItem>
                  <SidebarMenuButton
                    isActive={activeSetting === 'stages-workflow'}
                    onClick={() => onSelectSetting('stages-workflow')}
                    className={cn(
                      'h-9 text-xs transition font-medium',
                      activeSetting === 'stages-workflow' &&
                        'bg-primary/10 text-primary font-bold shadow-2xs hover:bg-primary/15'
                    )}
                  >
                    <IconGitFork className="h-4 w-4 shrink-0" />
                    <span className="flex-1 truncate">Workflow</span>
                    {activeSetting === 'stages-workflow' && (
                      <IconChevronRight className="h-3.5 w-3.5 text-primary shrink-0" />
                    )}
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>

        <SidebarFooter className="border-t p-3 text-[11px] text-muted-foreground">
          <p className="truncate">Studio Configuration Engine</p>
        </SidebarFooter>
      </Sidebar>
    )
  }

  // REPORTS SIDEBAR MODE
  if (mode === 'reports') {
    return (
      <Sidebar collapsible="icon" {...props}>
        {/* Header with Back button to Studio CRM */}
        <SidebarHeader className="border-b p-3">
          <div className="flex items-center justify-between gap-2">
            <div className="flex items-center gap-2 min-w-0">
              <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-emerald-600 text-white shadow-xs">
                <IconReportAnalytics className="h-4 w-4" />
              </span>
              <div className="truncate">
                <h3 className="text-xs font-bold leading-tight text-foreground">
                  Reports & BI
                </h3>
                <p className="text-[10px] text-muted-foreground truncate">
                  Studio Performance
                </p>
              </div>
            </div>
          </div>
          <Button
            variant="outline"
            size="sm"
            onClick={onExitSubmode}
            className="w-full mt-2 h-7 text-[11px] font-semibold gap-1.5 border-dashed"
          >
            <IconArrowLeft className="h-3.5 w-3.5" />
            <span>Back to Studio CRM</span>
          </Button>
        </SidebarHeader>

        {/* Reports Navigation Items */}
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel className="text-[10px] uppercase font-bold tracking-wider text-muted-foreground">
              Analytics & Statements
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {REPORT_MENU_ITEMS.map((item) => {
                  const Icon = item.icon
                  const isActive = activeReport === item.id
                  return (
                    <SidebarMenuItem key={item.id}>
                      <SidebarMenuButton
                        isActive={isActive}
                        onClick={() => onSelectReport(item.id)}
                        className={cn(
                          'h-9 text-xs transition font-medium',
                          isActive &&
                            'bg-primary/10 text-primary font-bold shadow-2xs hover:bg-primary/15'
                        )}
                      >
                        <Icon className="h-4 w-4 shrink-0" />
                        <span className="flex-1 truncate">{item.title}</span>
                        {item.badge && (
                          <span className="text-[9px] uppercase px-1.5 py-0.5 rounded-full bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 font-bold">
                            {item.badge}
                          </span>
                        )}
                        {isActive && (
                          <IconChevronRight className="h-3.5 w-3.5 text-primary shrink-0" />
                        )}
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  )
                })}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>

        <SidebarFooter className="border-t p-3 text-[11px] text-muted-foreground">
          <div className="flex items-center gap-2">
            <span className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
            <span className="font-medium">Live Studio Telemetry</span>
          </div>
        </SidebarFooter>
      </Sidebar>
    )
  }

  // Normal Main CRM Sidebar
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader className="border-b p-3">
        <div className="flex items-center gap-2.5">
          <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-primary text-primary-foreground shadow-xs">
            <IconCamera className="h-4 w-4" />
          </span>
          <div className="truncate">
            <h3 className="text-xs font-bold leading-tight text-foreground">
              Photography CRM
            </h3>
            <p className="text-[10px] text-muted-foreground">especa Studio</p>
          </div>
        </div>
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel className="text-[10px] uppercase font-bold tracking-wider text-muted-foreground">
            Core Modules
          </SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton
                  isActive={currentRoute.startsWith('project')}
                  onClick={() => onSelectRoute('projects-list')}
                  className={cn(
                    'h-9 text-xs',
                    currentRoute.startsWith('project') && 'font-bold text-primary bg-primary/10'
                  )}
                >
                  <IconFolders className="h-4 w-4" />
                  <span>Projects</span>
                </SidebarMenuButton>
              </SidebarMenuItem>

              <SidebarMenuItem>
                <SidebarMenuButton
                  isActive={currentRoute === 'sessions-list'}
                  onClick={() => onSelectRoute('sessions-list')}
                  className={cn(
                    'h-9 text-xs',
                    currentRoute === 'sessions-list' && 'font-bold text-primary bg-primary/10'
                  )}
                >
                  <IconCalendarEvent className="h-4 w-4" />
                  <span>Sessions</span>
                </SidebarMenuButton>
              </SidebarMenuItem>

              <SidebarMenuItem>
                <SidebarMenuButton
                  isActive={currentRoute.startsWith('lead')}
                  onClick={() => onSelectRoute('leads-list')}
                  className={cn(
                    'h-9 text-xs',
                    currentRoute.startsWith('lead') && 'font-bold text-primary bg-primary/10'
                  )}
                >
                  <IconTargetArrow className="h-4 w-4" />
                  <span>Leads</span>
                </SidebarMenuButton>
              </SidebarMenuItem>

              <SidebarMenuItem>
                <SidebarMenuButton
                  isActive={currentRoute === 'billings'}
                  onClick={() => onSelectRoute('billings')}
                  className={cn(
                    'h-9 text-xs',
                    currentRoute === 'billings' && 'font-bold text-primary bg-primary/10'
                  )}
                >
                  <IconReceipt className="h-4 w-4" />
                  <span>Billings</span>
                </SidebarMenuButton>
              </SidebarMenuItem>

              <SidebarMenuItem>
                <SidebarMenuButton
                  isActive={currentRoute.startsWith('order')}
                  onClick={() => onSelectRoute('orders-list')}
                  className={cn(
                    'h-9 text-xs',
                    currentRoute.startsWith('order') && 'font-bold text-primary bg-primary/10'
                  )}
                >
                  <IconShoppingBag className="h-4 w-4" />
                  <span>Orders</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter className="border-t p-3 text-[11px] text-muted-foreground">
        <p className="truncate">especa Studio Suite v2.0</p>
      </SidebarFooter>
    </Sidebar>
  )
}
