import * as React from 'react'
import { Card } from '@repo/ui/components/ui/card'
import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@repo/ui/components/ui/select'
import { toast } from '@repo/ui/components/ui/sonner'
import { cn } from '@repo/ui/lib/utils'
import {
  IconPackage,
  IconClock,
  IconFilter,
  IconSearch,
  IconCheck,
  IconUsers,
} from '@tabler/icons-react'
import { useAuth } from '@repo/ui/hooks/use-auth'
import type { Project, CustomerDeliverable } from '../../types/project'
import { AddDeliverableModal } from '../../components/AddDeliverableModal'

interface DeliverableRowItem {
  projectId: string
  projectName: string
  projectDates: string
  clientInitials: string
  deliverableId: string
  deliverableName: string
  assignee: string
  dueDate: string
  status: 'New' | 'In Progress' | 'Client Review' | 'Completed' | 'Delivered'
  readyCount: number
  totalCount: number
  progressPct: number
}

interface UserDeliverablesPageProps {
  projects: Project[]
  onUpdateProject?: (project: Project) => void
}

type TabType = 'Upcoming' | 'TBD' | 'Overdue' | 'Delivered'

export function UserDeliverablesPage({
  projects,
  onUpdateProject,
}: UserDeliverablesPageProps) {
  const { user } = useAuth()
  const [activeTab, setActiveTab] = React.useState<TabType>('Upcoming')
  const [yourDeliverablesOnly, setYourDeliverablesOnly] = React.useState(false)
  const [searchQuery, setSearchQuery] = React.useState('')
  const [isFilterOpen, setIsFilterOpen] = React.useState(false)
  const [sortOption, setSortOption] = React.useState<'due_asc' | 'due_desc' | 'project'>('due_asc')

  const [addModalOpen, setAddModalOpen] = React.useState(false)
  const [targetProjectForAdd, setTargetProjectForAdd] = React.useState<Project | null>(null)

  const isAssignedToUser = (assigneeName: string) => {
    if (!user) return false
    const userName = `${user.first_name || ''} ${user.last_name || ''}`.trim().toLowerCase()
    const userEmail = (user.email || '').toLowerCase()
    const aName = (assigneeName || '').toLowerCase()
    return (userName && aName.includes(userName)) || (userEmail && aName.includes(userEmail))
  }

  // Flatten deliverables across projects into items
  const allDeliverables = React.useMemo(() => {
    const list: DeliverableRowItem[] = []

    projects.forEach((p) => {
      const pDels = p.deliverables && p.deliverables.length > 0
        ? p.deliverables
        : [
          {
            id: `${p.id}-del-1`,
            deliverable: 'Raw Photos & Highlights',
            dueDate: p.endDate || '2026-12-01',
            assignee: p.projectManager?.name || 'Rohan Mehta',
            status: 'New' as const,
          },
          {
            id: `${p.id}-del-2`,
            deliverable: 'Handcrafted Album',
            dueDate: '2026-12-15',
            assignee: 'Unassigned',
            status: 'New' as const,
          },
        ]

      const ready = pDels.filter((d) => d.status === 'Completed' || d.status === 'Delivered').length
      const total = pDels.length
      const pct = total > 0 ? Math.round((ready / total) * 100) : 0

      pDels.forEach((del) => {
        list.push({
          projectId: p.id,
          projectName: p.projectName,
          projectDates: `${p.startDate ? p.startDate.slice(5) : '08-11'} - ${p.endDate ? p.endDate.slice(5) : '11-11'}`,
          clientInitials: p.clientInitials || (p.projectName ? p.projectName.slice(0, 2).toUpperCase() : 'PR'),
          deliverableId: del.id,
          deliverableName: del.deliverable,
          assignee: del.assignee || 'Unassigned',
          dueDate: del.dueDate || '2026-12-01',
          status: del.status || 'New',
          readyCount: ready,
          totalCount: total,
          progressPct: pct,
        })
      })
    })

    return list
  }, [projects])

  // Filter deliverables
  const filteredDeliverables = React.useMemo(() => {
    const todayStr = new Date().toISOString().split('T')[0]

    return allDeliverables.filter((item) => {
      if (searchQuery.trim()) {
        const q = searchQuery.toLowerCase()
        const matchProj = item.projectName.toLowerCase().includes(q)
        const matchDel = item.deliverableName.toLowerCase().includes(q)
        const matchAssignee = item.assignee.toLowerCase().includes(q)
        if (!matchProj && !matchDel && !matchAssignee) return false
      }

      if (yourDeliverablesOnly && !isAssignedToUser(item.assignee)) {
        return false
      }

      if (activeTab === 'Delivered') {
        return item.status === 'Delivered'
      }

      if (activeTab === 'TBD') {
        return !item.dueDate || item.dueDate === 'TBD'
      }

      if (activeTab === 'Overdue') {
        return item.status !== 'Delivered' && item.status !== 'Completed' && item.dueDate < todayStr
      }

      if (activeTab === 'Upcoming') {
        return item.status !== 'Delivered' && item.dueDate >= todayStr
      }

      return true
    }).sort((a, b) => {
      if (sortOption === 'due_asc') return a.dueDate.localeCompare(b.dueDate)
      if (sortOption === 'due_desc') return b.dueDate.localeCompare(a.dueDate)
      return a.projectName.localeCompare(b.projectName)
    })
  }, [allDeliverables, activeTab, yourDeliverablesOnly, searchQuery, sortOption, user])

  // Group by Month (Screenshot 3: "December, 2026")
  const groupedDeliverables = React.useMemo(() => {
    const groups: Record<string, DeliverableRowItem[]> = {}

    filteredDeliverables.forEach((item) => {
      let groupKey = 'December, 2026'
      if (item.dueDate && item.dueDate.includes('-')) {
        const [year, month] = item.dueDate.split('-')
        const monthNames = [
          'January', 'February', 'March', 'April', 'May', 'June',
          'July', 'August', 'September', 'October', 'November', 'December'
        ]
        const mIdx = parseInt(month, 10) - 1
        groupKey = `${monthNames[mIdx] || 'December'}, ${year || '2026'}`
      }
      if (!groups[groupKey]) groups[groupKey] = []
      groups[groupKey].push(item)
    })

    return groups
  }, [filteredDeliverables])

  const handleStatusChange = (item: DeliverableRowItem, newStatus: CustomerDeliverable['status']) => {
    const targetProj = projects.find((p) => p.id === item.projectId)
    if (!targetProj || !onUpdateProject) return

    const updatedDeliverables = (targetProj.deliverables || []).map((d) => {
      if (d.id === item.deliverableId) {
        return { ...d, status: newStatus }
      }
      return d
    })

    onUpdateProject({
      ...targetProj,
      deliverables: updatedDeliverables,
    })
    toast.success(`Updated status to ${newStatus}`)
  }

  const formatDisplayDate = (dStr: string) => {
    if (!dStr || dStr === 'TBD') return 'No Due Date'
    try {
      const d = new Date(dStr)
      if (!isNaN(d.getTime())) {
        return d.toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' })
      }
    } catch { }
    return dStr
  }

  const upcomingCount = React.useMemo(() => {
    const today = new Date().toISOString().split('T')[0]
    return allDeliverables.filter((d) => d.status !== 'Delivered' && d.dueDate >= today).length
  }, [allDeliverables])

  const tbdCount = React.useMemo(() => {
    return allDeliverables.filter((d) => !d.dueDate || d.dueDate === 'TBD').length
  }, [allDeliverables])

  const overdueCount = React.useMemo(() => {
    const today = new Date().toISOString().split('T')[0]
    return allDeliverables.filter((d) => d.status !== 'Delivered' && d.status !== 'Completed' && d.dueDate < today).length
  }, [allDeliverables])

  const deliveredCount = React.useMemo(() => {
    return allDeliverables.filter((d) => d.status === 'Delivered').length
  }, [allDeliverables])

  return (
    <div className="space-y-6">
      {/* Top Filter and Controls Bar */}
      <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
        {/* Left Pills: Upcoming | TBD - No Due Date | Overdue | Delivered */}
        <div className="inline-flex items-center rounded-full bg-muted/60 p-1 border shadow-2xs flex-wrap">
          <button
            type="button"
            onClick={() => setActiveTab('Upcoming')}
            className={cn(
              'inline-flex items-center gap-1.5 rounded-full px-4 py-1.5 text-xs font-semibold transition-all',
              activeTab === 'Upcoming'
                ? 'bg-foreground text-background shadow-xs font-bold'
                : 'text-muted-foreground hover:text-foreground'
            )}
          >
            <span>Upcoming</span>
            <span className={cn(
              'text-[10px] px-1.5 py-0.2 rounded-full font-bold',
              activeTab === 'Upcoming' ? 'bg-background/20 text-background' : 'bg-muted text-muted-foreground'
            )}>
              {upcomingCount}
            </span>
          </button>

          <button
            type="button"
            onClick={() => setActiveTab('TBD')}
            className={cn(
              'inline-flex items-center gap-1.5 rounded-full px-4 py-1.5 text-xs font-semibold transition-all',
              activeTab === 'TBD'
                ? 'bg-foreground text-background shadow-xs font-bold'
                : 'text-muted-foreground hover:text-foreground'
            )}
          >
            <span>TBD - No Due Date</span>
            {tbdCount > 0 && (
              <span className={cn(
                'text-[10px] px-1.5 py-0.2 rounded-full font-bold',
                activeTab === 'TBD' ? 'bg-background/20 text-background' : 'bg-muted text-muted-foreground'
              )}>
                {tbdCount}
              </span>
            )}
          </button>

          <button
            type="button"
            onClick={() => setActiveTab('Overdue')}
            className={cn(
              'inline-flex items-center gap-1.5 rounded-full px-4 py-1.5 text-xs font-semibold transition-all',
              activeTab === 'Overdue'
                ? 'bg-foreground text-background shadow-xs font-bold'
                : 'text-muted-foreground hover:text-foreground'
            )}
          >
            <span>Overdue</span>
            {overdueCount > 0 && (
              <span className={cn(
                'text-[10px] px-1.5 py-0.2 rounded-full font-bold',
                activeTab === 'Overdue' ? 'bg-background/20 text-background' : 'bg-destructive/20 text-destructive'
              )}>
                {overdueCount}
              </span>
            )}
          </button>

          <button
            type="button"
            onClick={() => setActiveTab('Delivered')}
            className={cn(
              'inline-flex items-center gap-1.5 rounded-full px-4 py-1.5 text-xs font-semibold transition-all',
              activeTab === 'Delivered'
                ? 'bg-foreground text-background shadow-xs font-bold'
                : 'text-muted-foreground hover:text-foreground'
            )}
          >
            <span>Delivered</span>
            <span className={cn(
              'text-[10px] px-1.5 py-0.2 rounded-full font-bold',
              activeTab === 'Delivered' ? 'bg-background/20 text-background' : 'bg-muted text-muted-foreground'
            )}>
              {deliveredCount}
            </span>
          </button>
        </div>

        {/* Right Controls: Your Deliverables | Filters | Sort */}
        <div className="flex items-center gap-2 flex-wrap">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setYourDeliverablesOnly(!yourDeliverablesOnly)}
            className={cn(
              'h-9 rounded-full text-xs gap-1.5 px-3.5 border-border shadow-2xs transition-all font-medium',
              yourDeliverablesOnly && 'bg-primary/10 border-primary/40 text-primary font-bold'
            )}
          >
            <IconCheck className={cn('h-3.5 w-3.5', yourDeliverablesOnly ? 'opacity-100 text-primary' : 'opacity-40')} />
            <span>Your Deliverables</span>
          </Button>

          <Button
            variant="outline"
            size="sm"
            onClick={() => setIsFilterOpen(!isFilterOpen)}
            className={cn(
              'h-9 rounded-full text-xs gap-1.5 px-3.5 border-border shadow-2xs transition-all font-medium',
              isFilterOpen && 'bg-accent text-accent-foreground border-accent'
            )}
          >
            <IconFilter className="h-3.5 w-3.5" />
            <span>Filters</span>
          </Button>

          <Select
            value={sortOption}
            onValueChange={(val) => setSortOption(val as any)}
          >
            <SelectTrigger className="h-9 rounded-full text-xs bg-card border-border shadow-2xs w-[130px]">
              <div className="flex items-center gap-1.5 text-muted-foreground truncate">
                <IconClock className="h-3.5 w-3.5 shrink-0" />
                <span className="font-semibold text-foreground text-xs truncate">
                  {sortOption === 'due_asc' ? 'Due Date' : sortOption === 'due_desc' ? 'Latest' : 'Project'}
                </span>
              </div>
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="due_asc" className="text-xs">
                Due Date (Earliest)
              </SelectItem>
              <SelectItem value="due_desc" className="text-xs">
                Due Date (Latest)
              </SelectItem>
              <SelectItem value="project" className="text-xs">
                Project Name
              </SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Filter Drawer */}
      {isFilterOpen && (
        <Card className="p-3.5 border-border/80 bg-muted/20">
          <div className="flex items-center gap-3">
            <div className="relative flex-1">
              <IconSearch className="absolute left-3 top-2.5 h-3.5 w-3.5 text-muted-foreground" />
              <Input
                placeholder="Search deliverables, projects or assignees..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-9 h-8 text-xs bg-card"
              />
            </div>
            {(searchQuery || yourDeliverablesOnly) && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => {
                  setSearchQuery('')
                  setYourDeliverablesOnly(false)
                }}
                className="text-xs h-8"
              >
                Reset
              </Button>
            )}
          </div>
        </Card>
      )}

      {/* Column Headers as in Screenshot 3 */}
      <div className="hidden lg:grid grid-cols-12 gap-4 px-5 text-[11px] font-bold uppercase tracking-wider text-muted-foreground select-none">
        <div className="col-span-5">Delivery Details</div>
        <div className="col-span-2 text-center">Assignee</div>
        <div className="col-span-2 text-center">Delivery Due Date</div>
        <div className="col-span-2 text-center">Progress</div>
        <div className="col-span-1 text-right">Status</div>
      </div>

      {/* Grouped by Month Cards */}
      <div className="space-y-6">
        {Object.keys(groupedDeliverables).length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-center border border-dashed rounded-2xl bg-card">
            <IconPackage className="h-10 w-10 text-muted-foreground/30 mb-2" />
            <h4 className="text-sm font-semibold text-foreground">No deliverables found</h4>
            <p className="text-xs text-muted-foreground max-w-sm mt-1">
              {searchQuery || yourDeliverablesOnly
                ? 'No deliverables match your filter.'
                : 'There are no items currently in this view.'}
            </p>
          </div>
        ) : (
          Object.entries(groupedDeliverables).map(([monthGroup, items]) => (
            <div key={monthGroup} className="space-y-3">
              <h3 className="text-sm font-bold text-foreground tracking-tight">
                {monthGroup}
              </h3>

              <div className="space-y-3">
                {items.map((item) => {
                  return (
                    <div
                      key={item.deliverableId}
                      className="group rounded-2xl border border-border/80 bg-blue-50/20 dark:bg-card/40 p-4 transition-all hover:border-primary/40 hover:shadow-xs"
                    >
                      <div className="grid grid-cols-1 lg:grid-cols-12 gap-4 items-center">
                        {/* Delivery Details */}
                        <div className="lg:col-span-5 flex items-center gap-3.5 min-w-0">
                          {/* Monogram Box */}
                          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-blue-100/70 text-blue-800 dark:bg-blue-950/60 dark:text-blue-300 font-bold text-sm tracking-wide shadow-2xs border border-blue-200/50 dark:border-blue-800/40">
                            {item.clientInitials}
                          </div>

                          {/* Titles */}
                          <div className="min-w-0 space-y-0.5">
                            <h4 className="text-sm font-bold text-foreground truncate group-hover:text-primary transition">
                              {item.projectName}
                            </h4>
                            <p className="text-xs font-medium text-muted-foreground truncate">
                              {item.deliverableName}
                            </p>
                            <p className="text-[10px] text-muted-foreground">
                              {item.projectDates}
                            </p>
                          </div>
                        </div>

                        {/* Assignee */}
                        <div className="lg:col-span-2 flex justify-start lg:justify-center">
                          <span className="inline-flex items-center gap-1.5 rounded-full bg-muted/60 px-3 py-1 text-xs font-medium text-foreground border border-border/60">
                            <IconUsers className="h-3 w-3 text-muted-foreground" />
                            {item.assignee === 'Unassigned' ? 'No Team Members' : item.assignee}
                          </span>
                        </div>

                        {/* Delivery Due Date */}
                        <div className="lg:col-span-2 flex justify-start lg:justify-center">
                          <span className="text-xs font-semibold text-primary">
                            {formatDisplayDate(item.dueDate)}
                          </span>
                        </div>

                        {/* Progress Bar: 0 | [====] | 4 */}
                        <div className="lg:col-span-2 flex items-center justify-start lg:justify-center gap-2">
                          <span className="text-xs font-semibold text-muted-foreground tabular-nums">
                            {item.readyCount}
                          </span>
                          <div className="w-24 bg-red-100 dark:bg-muted/60 rounded-full h-2.5 overflow-hidden flex items-center">
                            <div
                              className="bg-emerald-500 h-full rounded-full transition-all"
                              style={{ width: `${item.progressPct}%` }}
                            />
                          </div>
                          <span className="text-xs font-semibold text-muted-foreground tabular-nums">
                            {item.totalCount}
                          </span>
                          <span className="text-[10px] font-bold text-muted-foreground">
                            {item.progressPct}%
                          </span>
                        </div>

                        {/* Status Dropdown */}
                        <div className="lg:col-span-1 flex justify-start lg:justify-end">
                          <Select
                            value={item.status}
                            onValueChange={(val) => handleStatusChange(item, val as any)}
                          >
                            <SelectTrigger className="h-8 text-xs font-bold rounded-lg border-primary/20 bg-primary/10 text-primary hover:bg-primary/20 transition px-3 gap-1.5 w-[110px]">
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent align="end">
                              <SelectItem value="New" className="text-xs font-medium">New</SelectItem>
                              <SelectItem value="In Progress" className="text-xs font-medium">In Progress</SelectItem>
                              <SelectItem value="Client Review" className="text-xs font-medium">Client Review</SelectItem>
                              <SelectItem value="Completed" className="text-xs font-medium">Completed</SelectItem>
                              <SelectItem value="Delivered" className="text-xs font-medium">Delivered</SelectItem>
                            </SelectContent>
                          </Select>
                        </div>
                      </div>
                    </div>
                  )
                })}
              </div>
            </div>
          ))
        )}
      </div>

      {/* Add Deliverable Modal */}
      {targetProjectForAdd && (
        <AddDeliverableModal
          open={addModalOpen}
          onOpenChange={(op) => {
            setAddModalOpen(op)
            if (!op) setTargetProjectForAdd(null)
          }}
          onAdd={(newDel) => {
            if (targetProjectForAdd && onUpdateProject) {
              const updated = [...(targetProjectForAdd.deliverables || []), newDel]
              onUpdateProject({ ...targetProjectForAdd, deliverables: updated })
              toast.success(`Deliverable added to ${targetProjectForAdd.projectName}`)
            }
            setAddModalOpen(false)
            setTargetProjectForAdd(null)
          }}
        />
      )}
    </div>
  )
}
