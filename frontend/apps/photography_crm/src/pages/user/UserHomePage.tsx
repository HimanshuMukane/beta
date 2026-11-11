import * as React from 'react'
import { Card } from '@repo/ui/components/ui/card'
import { Button } from '@repo/ui/components/ui/button'
import { Badge } from '@repo/ui/components/ui/badge'
import {
  IconCamera,
  IconChecklist,
  IconPackage,
  IconClock,
  IconArrowRight,
  IconMapPin,
} from "@repo/ui/icons/tabler-icons"
import { useNavigate } from 'react-router-dom'
import { useAuth } from '@repo/ui/hooks/use-auth'
import type { ProjectSession } from '../../types/session'
import type { Project } from '../../types/project'

interface UserHomePageProps {
  sessions: ProjectSession[]
  projects: Project[]
}

export function UserHomePage({ sessions, projects }: UserHomePageProps) {
  const navigate = useNavigate()
  const { user } = useAuth()

  const todayStr = new Date().toISOString().split('T')[0]

  const upcomingShoots = React.useMemo(() => {
    return sessions
      .filter((s) => s.date >= todayStr)
      .sort((a, b) => a.date.localeCompare(b.date))
      .slice(0, 3)
  }, [sessions, todayStr])

  const pendingDeliverables = React.useMemo(() => {
    const list: {
      projectId: string
      projectName: string
      deliverableName: string
      dueDate: string
      status: string
    }[] = []

    projects.forEach((p) => {
      ;(p.deliverables || []).forEach((d) => {
        if (d.status !== 'Delivered' && d.status !== 'Completed') {
          list.push({
            projectId: p.id,
            projectName: p.projectName,
            deliverableName: d.deliverable,
            dueDate: d.dueDate,
            status: d.status,
          })
        }
      })
    })

    return list.slice(0, 4)
  }, [projects])

  const firstName = user?.first_name || 'Crew Member'

  return (
    <div className="space-y-6">
      {/* Welcome Hero Banner */}
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-primary/10 via-primary/5 to-muted/40 p-6 border border-border/60">
        <div className="max-w-2xl space-y-2">
          <Badge variant="outline" className="text-[11px] font-semibold bg-background/80 text-primary border-primary/30">
            {user?.role_name || 'Studio Team Workspace'}
          </Badge>
          <h2 className="text-xl sm:text-2xl font-bold tracking-tight text-foreground">
            Welcome back, {firstName}
          </h2>
          <p className="text-xs sm:text-sm text-muted-foreground leading-relaxed">
            Here is your daily snapshot: assigned photography shoots, ongoing pre/post-production matrix tasks, and client delivery milestones.
          </p>
        </div>
      </div>

      {/* 3 Metric Summary Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <Card
          onClick={() => navigate('/shoots')}
          className="p-4 cursor-pointer hover:border-primary/40 transition hover:shadow-xs border-border/80"
        >
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-muted-foreground">Upcoming Shoots</span>
            <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-500/10 text-emerald-600">
              <IconCamera className="h-4 w-4" />
            </span>
          </div>
          <div className="mt-2 flex items-baseline gap-2">
            <span className="text-2xl font-bold text-foreground">
              {sessions.filter((s) => s.date >= todayStr).length}
            </span>
            <span className="text-xs text-muted-foreground">scheduled</span>
          </div>
        </Card>

        <Card
          onClick={() => navigate('/tasks')}
          className="p-4 cursor-pointer hover:border-primary/40 transition hover:shadow-xs border-border/80"
        >
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-muted-foreground">Active Production Projects</span>
            <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-500/10 text-blue-600">
              <IconChecklist className="h-4 w-4" />
            </span>
          </div>
          <div className="mt-2 flex items-baseline gap-2">
            <span className="text-2xl font-bold text-foreground">{projects.length}</span>
            <span className="text-xs text-muted-foreground">in pipeline</span>
          </div>
        </Card>

        <Card
          onClick={() => navigate('/deliverables')}
          className="p-4 cursor-pointer hover:border-primary/40 transition hover:shadow-xs border-border/80"
        >
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-muted-foreground">Pending Deliverables</span>
            <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-amber-500/10 text-amber-600">
              <IconPackage className="h-4 w-4" />
            </span>
          </div>
          <div className="mt-2 flex items-baseline gap-2">
            <span className="text-2xl font-bold text-foreground">
              {pendingDeliverables.length}
            </span>
            <span className="text-xs text-muted-foreground">due for review</span>
          </div>
        </Card>
      </div>

      {/* Two Columns: Next Upcoming Shoots & Priority Deliverables */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Next Shoots */}
        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <h3 className="text-sm font-bold text-foreground">Next Upcoming Shoots</h3>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => navigate('/shoots')}
              className="text-xs text-primary font-semibold h-8 gap-1"
            >
              <span>View all shoots</span>
              <IconArrowRight className="h-3.5 w-3.5" />
            </Button>
          </div>

          <div className="space-y-2.5">
            {upcomingShoots.length === 0 ? (
              <div className="rounded-xl border border-dashed p-6 text-center text-xs text-muted-foreground bg-card">
                No upcoming shoots scheduled.
              </div>
            ) : (
              upcomingShoots.map((s) => (
                <div
                  key={s.id}
                  onClick={() => navigate('/shoots')}
                  className="rounded-xl border border-border bg-card p-3.5 hover:border-primary/40 hover:shadow-2xs transition cursor-pointer flex items-center justify-between gap-3"
                >
                  <div className="flex items-center gap-3 min-w-0">
                    <div className="flex flex-col items-center justify-center rounded-lg bg-muted/60 border border-border/80 px-2.5 py-1 text-center min-w-[50px]">
                      <span className="text-[10px] font-semibold text-muted-foreground uppercase">
                        {s.date ? new Date(s.date).toLocaleDateString('en-US', { weekday: 'short' }) : 'DAY'}
                      </span>
                      <span className="text-base font-bold text-foreground leading-tight">
                        {s.date ? s.date.slice(8, 10) : '01'}
                      </span>
                    </div>

                    <div className="min-w-0 space-y-0.5">
                      <p className="text-xs font-bold text-foreground truncate">
                        {s.projectName} {s.title ? `- ${s.title}` : ''}
                      </p>
                      <div className="flex items-center gap-2 text-[11px] text-muted-foreground">
                        <span className="inline-flex items-center gap-1">
                          <IconClock className="h-3 w-3" />
                          {s.timeRange || s.startTime || 'TBD'}
                        </span>
                        {s.city && (
                          <span className="inline-flex items-center gap-1">
                            <IconMapPin className="h-3 w-3" />
                            {s.city}
                          </span>
                        )}
                      </div>
                    </div>
                  </div>

                  <Badge variant="outline" className="text-[10px] shrink-0 font-medium bg-primary/5 text-primary border-primary/20">
                    {s.status || 'New'}
                  </Badge>
                </div>
              ))
            )}
          </div>
        </div>

        {/* Priority Deliverables */}
        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <h3 className="text-sm font-bold text-foreground">Priority Deliverables</h3>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => navigate('/deliverables')}
              className="text-xs text-primary font-semibold h-8 gap-1"
            >
              <span>View all deliverables</span>
              <IconArrowRight className="h-3.5 w-3.5" />
            </Button>
          </div>

          <div className="space-y-2.5">
            {pendingDeliverables.length === 0 ? (
              <div className="rounded-xl border border-dashed p-6 text-center text-xs text-muted-foreground bg-card">
                All client deliverables are up to date.
              </div>
            ) : (
              pendingDeliverables.map((d, idx) => (
                <div
                  key={`${d.projectId}-${idx}`}
                  onClick={() => navigate('/deliverables')}
                  className="rounded-xl border border-border bg-card p-3.5 hover:border-primary/40 hover:shadow-2xs transition cursor-pointer flex items-center justify-between gap-3"
                >
                  <div className="flex items-center gap-3 min-w-0">
                    <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-blue-500/10 text-blue-600 font-bold text-xs">
                      {d.projectName.slice(0, 2).toUpperCase()}
                    </span>

                    <div className="min-w-0 space-y-0.5">
                      <p className="text-xs font-bold text-foreground truncate">
                        {d.projectName}
                      </p>
                      <p className="text-[11px] text-muted-foreground truncate">
                        {d.deliverableName} • Due {d.dueDate || 'TBD'}
                      </p>
                    </div>
                  </div>

                  <Badge variant="outline" className="text-[10px] shrink-0 font-medium bg-amber-500/10 text-amber-700 border-amber-500/20">
                    {d.status}
                  </Badge>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
