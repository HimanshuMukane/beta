import * as React from 'react'
import {
  IconHome,
  IconCamera,
  IconChecklist,
  IconPackage,
  IconLogout,
  IconUser,
} from "@repo/ui/icons/tabler-icons"
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
} from '@repo/ui/components/ui/sidebar'
import { Button } from '@repo/ui/components/ui/button'
import { useAuth } from '@repo/ui/hooks/use-auth'
import { cn } from '@repo/ui/lib/utils'

export type UserRouteKey = 'home' | 'shoots' | 'tasks' | 'deliverables'

interface UserSidebarProps extends React.ComponentProps<typeof Sidebar> {
  currentRoute: UserRouteKey
  onSelectRoute: (route: UserRouteKey) => void
}

const USER_NAV_ITEMS: {
  id: UserRouteKey
  title: string
  icon: React.ComponentType<{ className?: string }>
}[] = [
  { id: 'home', title: 'Home', icon: IconHome },
  { id: 'shoots', title: 'Shoots', icon: IconCamera },
  { id: 'tasks', title: 'Tasks', icon: IconChecklist },
  { id: 'deliverables', title: 'Deliverables', icon: IconPackage },
]

export function UserSidebar({
  currentRoute,
  onSelectRoute,
  ...props
}: UserSidebarProps) {
  const { user, logout } = useAuth()

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader className="border-b p-3">
        <div className="flex items-center gap-2.5">
          <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-emerald-600 text-white shadow-xs">
            <IconCamera className="h-4 w-4" />
          </span>
          <div className="truncate">
            <h3 className="text-xs font-bold leading-tight text-foreground">
              Photography Portal
            </h3>
            <p className="text-[10px] text-muted-foreground">Crew Workspace</p>
          </div>
        </div>
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel className="text-[10px] uppercase font-bold tracking-wider text-muted-foreground">
            Navigation
          </SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {USER_NAV_ITEMS.map((item) => {
                const Icon = item.icon
                const isActive = currentRoute === item.id
                return (
                  <SidebarMenuItem key={item.id}>
                    <SidebarMenuButton
                      isActive={isActive}
                      onClick={() => onSelectRoute(item.id)}
                      className={cn(
                        'h-9 text-xs transition-colors font-medium',
                        isActive && 'font-bold text-emerald-600 dark:text-emerald-400 bg-emerald-500/10'
                      )}
                    >
                      <Icon className={cn('h-4 w-4', isActive && 'text-emerald-600 dark:text-emerald-400')} />
                      <span>{item.title}</span>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                )
              })}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter className="border-t p-3 space-y-2">
        {user && (
          <div className="flex items-center justify-between gap-2 rounded-lg bg-muted/40 p-2 border border-border/50">
            <div className="flex items-center gap-2 min-w-0">
              <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-emerald-500/10 text-emerald-600">
                <IconUser className="h-4 w-4" />
              </span>
              <div className="min-w-0 flex-1">
                <p className="text-xs font-semibold leading-none truncate text-foreground">
                  {user.first_name ? `${user.first_name} ${user.last_name || ''}`.trim() : user.email}
                </p>
                <p className="text-[10px] text-muted-foreground truncate leading-none mt-1">
                  {user.role_name || 'Team Member'}
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

