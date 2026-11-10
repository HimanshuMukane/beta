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
  IconLogout,
  IconUser,
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
  SidebarMenuSub,
  SidebarMenuSubItem,
  SidebarMenuSubButton,
} from '@repo/ui/components/ui/sidebar'
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@repo/ui/components/ui/collapsible'
import { Button } from '@repo/ui/components/ui/button'
import { useAuth } from '@repo/ui/hooks/use-auth'
import { cn } from '@repo/ui/lib/utils'
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
  currentRoute: string
  activeReport: ReportCategory
  activeSetting: SettingsTab
  onSelectRoute: (route: string) => void
  onSelectReport: (report: ReportCategory) => void
  onSelectSetting: (setting: SettingsTab) => void
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

const SETTINGS_SECTIONS: {
  title: string
  items: { id: SettingsTab; title: string; icon: React.ComponentType<{ className?: string }> }[]
}[] = [
    {
      title: 'General',
      items: [
        { id: 'master-data', title: 'Master Data', icon: IconDatabase },
        { id: 'photographers', title: 'Photographers', icon: IconUsers },
        { id: 'sequence-id', title: 'Sequence ID', icon: IconId },
      ],
    },
    {
      title: 'Preset Configuration',
      items: [
        { id: 'preset-services', title: 'Services', icon: IconBriefcase },
        { id: 'preset-packages', title: 'Packages', icon: IconBox },
        { id: 'preset-orders', title: 'Orders', icon: IconShoppingBag },
        { id: 'preset-deliverables', title: 'Deliverables', icon: IconPackage },
      ],
    },
    {
      title: 'Customize Design & Templates',
      items: [
        { id: 'template-forms', title: 'Forms', icon: IconFileText },
        { id: 'template-quotation-pdf', title: 'Quotation PDF', icon: IconFileInvoice },
        { id: 'template-invoice-pdf', title: 'Invoice PDF', icon: IconFileCheck },
        { id: 'template-contracts', title: 'Contracts', icon: IconSignature },
        { id: 'template-feedbacks', title: 'Feedbacks', icon: IconMessageDots },
      ],
    },
    {
      title: 'Stages',
      items: [
        { id: 'stages-workflow', title: 'Workflow', icon: IconGitFork },
      ],
    },
  ]

export function PhotographySidebar({
  currentRoute,
  activeReport,
  activeSetting,
  onSelectRoute,
  onSelectReport,
  onSelectSetting,
  ...props
}: PhotographySidebarProps) {
  const isReportsActive = currentRoute.startsWith('report')
  const isSettingsActive = currentRoute.startsWith('setting')
  const { user, logout } = useAuth()

  const [reportsOpen, setReportsOpen] = React.useState(isReportsActive)
  const [settingsOpen, setSettingsOpen] = React.useState(isSettingsActive)

  React.useEffect(() => {
    if (isReportsActive) setReportsOpen(true)
  }, [isReportsActive])

  React.useEffect(() => {
    if (isSettingsActive) setSettingsOpen(true)
  }, [isSettingsActive])

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
        {/* Core Modules */}
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

        {/* Analytics & Reports Section */}
        <SidebarGroup>
          <SidebarGroupLabel className="text-[10px] uppercase font-bold tracking-wider text-muted-foreground">
            Analytics
          </SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <Collapsible
                open={reportsOpen}
                onOpenChange={setReportsOpen}
                className="group/collapsible"
              >
                <SidebarMenuItem>
                  <CollapsibleTrigger asChild>
                    <SidebarMenuButton
                      isActive={isReportsActive}
                      className={cn(
                        'h-9 text-xs font-medium',
                        isReportsActive && 'font-bold text-primary bg-primary/10'
                      )}
                    >
                      <IconReportAnalytics className="h-4 w-4" />
                      <span className="flex-1 text-left">Reports & BI</span>
                      <IconChevronRight className="h-3.5 w-3.5 ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
                    </SidebarMenuButton>
                  </CollapsibleTrigger>
                  <CollapsibleContent>
                    <SidebarMenuSub className="py-1">
                      {REPORT_MENU_ITEMS.map((item) => {
                        const Icon = item.icon
                        const isSubActive = isReportsActive && activeReport === item.id
                        return (
                          <SidebarMenuSubItem key={item.id}>
                            <SidebarMenuSubButton
                              isActive={isSubActive}
                              onClick={() => onSelectReport(item.id)}
                              className={cn(
                                'h-8 text-xs cursor-pointer',
                                isSubActive && 'font-bold text-primary bg-primary/10'
                              )}
                            >
                              <Icon className="h-3.5 w-3.5 shrink-0" />
                              <span className="flex-1 truncate">{item.title}</span>
                              {item.badge && (
                                <span className="text-[8px] uppercase px-1 py-0.2 rounded-full bg-emerald-500/15 text-emerald-600 font-bold">
                                  {item.badge}
                                </span>
                              )}
                            </SidebarMenuSubButton>
                          </SidebarMenuSubItem>
                        )
                      })}
                    </SidebarMenuSub>
                  </CollapsibleContent>
                </SidebarMenuItem>
              </Collapsible>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        {/* Settings & Configurations Section */}
        <SidebarGroup>
          <SidebarGroupLabel className="text-[10px] uppercase font-bold tracking-wider text-muted-foreground">
            Configuration
          </SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <Collapsible
                open={settingsOpen}
                onOpenChange={setSettingsOpen}
                className="group/collapsible"
              >
                <SidebarMenuItem>
                  <CollapsibleTrigger asChild>
                    <SidebarMenuButton
                      isActive={isSettingsActive}
                      className={cn(
                        'h-9 text-xs font-medium',
                        isSettingsActive && 'font-bold text-primary bg-primary/10'
                      )}
                    >
                      <IconSettings className="h-4 w-4" />
                      <span className="flex-1 text-left">Studio Settings</span>
                      <IconChevronRight className="h-3.5 w-3.5 ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
                    </SidebarMenuButton>
                  </CollapsibleTrigger>
                  <CollapsibleContent>
                    <SidebarMenuSub className="py-1">
                      {SETTINGS_SECTIONS.map((section, sIdx) => (
                        <React.Fragment key={section.title}>
                          <div className={cn(
                            "px-2 pb-1 text-[10px] font-bold tracking-wider text-muted-foreground uppercase select-none",
                            sIdx === 0 ? "pt-1" : "pt-3 border-t border-border/40 mt-1.5"
                          )}>
                            {section.title}
                          </div>
                          {section.items.map((item) => {
                            const Icon = item.icon
                            const isSubActive = isSettingsActive && activeSetting === item.id
                            return (
                              <SidebarMenuSubItem key={item.id}>
                                <SidebarMenuSubButton
                                  isActive={isSubActive}
                                  onClick={() => onSelectSetting(item.id)}
                                  className={cn(
                                    'h-8 text-xs cursor-pointer rounded-md transition-all font-medium',
                                    isSubActive
                                      ? 'font-bold text-primary bg-primary/10 shadow-2xs'
                                      : 'text-muted-foreground hover:text-foreground hover:bg-muted/50'
                                  )}
                                >
                                  <Icon className={cn("h-3.5 w-3.5 shrink-0", isSubActive ? "text-primary" : "text-muted-foreground")} />
                                  <span className="flex-1 truncate">{item.title}</span>
                                </SidebarMenuSubButton>
                              </SidebarMenuSubItem>
                            )
                          })}
                        </React.Fragment>
                      ))}
                    </SidebarMenuSub>
                  </CollapsibleContent>
                </SidebarMenuItem>
              </Collapsible>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter className="border-t p-3 space-y-2">
        {user && (
          <div className="flex items-center justify-between gap-2 rounded-lg bg-muted/40 p-2 border border-border/50">
            <div className="flex items-center gap-2 min-w-0">
              <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary">
                <IconUser className="h-4 w-4" />
              </span>
              <div className="min-w-0 flex-1">
                <p className="text-xs font-semibold leading-none truncate text-foreground">
                  {user.first_name ? `${user.first_name} ${user.last_name || ''}`.trim() : user.email}
                </p>
                <p className="text-[10px] text-muted-foreground truncate leading-none mt-1">
                  {user.role_name || user.email}
                </p>
              </div>
            </div>
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 text-muted-foreground hover:text-destructive hover:bg-destructive/10 shrink-0"
              title="Sign Out"
              onClick={() => logout()}
            >
              <IconLogout className="h-3.5 w-3.5" />
            </Button>
          </div>
        )}
        <p className="text-[11px] text-muted-foreground truncate">especa Studio Suite v2.0</p>
      </SidebarFooter>
    </Sidebar>
  )
}