import * as React from 'react'
import { Card, CardContent } from '@repo/ui/components/ui/card'
import { Button } from '@repo/ui/components/ui/button'
import { Badge } from '@repo/ui/components/ui/badge'
import { Checkbox } from '@repo/ui/components/ui/checkbox'
import {
  DataTable,
  type DataTableFeatures,
  type DataTableFilterConfig,
  type DataTableBulkAction,
  type DataTableRowAction,
} from '@repo/ui/components/data-table'
import { toast } from '@repo/ui/components/ui/sonner'
import { cn } from '@repo/ui/lib/utils'
import { createColumnHelper } from '@tanstack/react-table'
import {
  IconPlus,
  IconTable,
  IconLayoutKanban,
  IconEdit,
  IconTrash,
  IconCamera,
  IconChecklist,
  IconPhoto,
  IconBriefcase,
} from '@tabler/icons-react'
import type { Project, ProjectStage, ProjectStatus } from '../../types/project'
import { ProjectKanbanView } from '../../components/ProjectKanbanView'
import { useWorkflows } from '../../hooks/useWorkflows'

interface ProjectsPageProps {
  projects: Project[]
  onAddNewProject: () => void
  onEditProject: (project: Project) => void
  onDeleteProject: (projectId: string) => void
  onStageChange: (projectId: string, newStage: ProjectStage) => void
}

const columnHelper = createColumnHelper<DataTableFeatures, Project>()

const STATUSES: (ProjectStatus | 'All')[] = [
  'All',
  'New',
  'Booked',
  'Active',
  'Completed',
  'On Hold',
  'Cancelled',
]

export function ProjectsPage({
  projects,
  onAddNewProject,
  onEditProject,
  onDeleteProject,
  onStageChange,
}: ProjectsPageProps) {
  const { workflows } = useWorkflows('project')

  const dynamicStages: string[] = React.useMemo(() => {
    if (workflows && workflows.length > 0) {
      return workflows.map((w) => w.stage)
    }
    return ['Booked', 'Shoot Planning', 'Shoot Completed', 'Editing', 'Delivery', 'Closure']
  }, [workflows])

  const stageColorMap = React.useMemo(() => {
    const map = new Map<string, string>()
    ;(workflows || []).forEach((w) => {
      map.set(w.stage.toLowerCase(), w.color_code)
    })
    return map
  }, [workflows])

  const [viewMode, setViewMode] = React.useState<'table' | 'kanban'>('table')
  const [stageFilter, setStageFilter] = React.useState<string>('All')
  const [statusFilter, setStatusFilter] = React.useState<ProjectStatus | 'All'>('All')

  // Filtered projects
  const filteredProjects = React.useMemo(() => {
    return projects.filter((p) => {
      if (stageFilter !== 'All' && p.stage !== stageFilter) return false
      if (statusFilter !== 'All' && p.status !== statusFilter) return false
      return true
    })
  }, [projects, stageFilter, statusFilter])

  // KPIs
  const metrics = React.useMemo(() => {
    const total = projects.length
    const activeShoots = projects.filter(
      (p) => p.stage === 'Shoot Planning' || p.stage === 'Shoot Completed'
    ).length
    const inEditing = projects.filter((p) => p.stage === 'Editing').length
    const delivered = projects.filter(
      (p) => p.stage === 'Delivery' || p.stage === 'Closure'
    ).length
    return { total, activeShoots, inEditing, delivered }
  }, [projects])

  // 13 Required Columns:
  // ID, Project Name, Client Name, Project Manager, Project Type, Sessions, Tasks, Deliverables, Status, Stage, Start Date, End Date, Created Date
  const columns = React.useMemo(
    () =>
      columnHelper.columns([
        columnHelper.display({
          id: 'select',
          header: ({ table }) => (
            <Checkbox
              checked={
                table.getIsAllPageRowsSelected() ||
                (table.getIsSomePageRowsSelected() && 'indeterminate')
              }
              onCheckedChange={(value: boolean | 'indeterminate') =>
                table.toggleAllPageRowsSelected(!!value)
              }
              aria-label="Select all"
            />
          ),
          cell: ({ row }) => (
            <Checkbox
              checked={row.getIsSelected()}
              onCheckedChange={(value: boolean | 'indeterminate') =>
                row.toggleSelected(!!value)
              }
              aria-label="Select row"
            />
          ),
          size: 40,
        }),
        // 1. ID
        columnHelper.accessor('id', {
          header: 'ID',
          cell: (info) => (
            <span className="font-mono text-xs text-muted-foreground font-semibold">
              {info.getValue()}
            </span>
          ),
        }),
        // 2. Project Name
        columnHelper.accessor('projectName', {
          header: 'Project Name',
          cell: (info) => {
            const row = info.row.original
            return (
              <div className="flex flex-col min-w-[170px]">
                <span
                  onClick={() => onEditProject(row)}
                  className="font-semibold text-foreground hover:text-primary hover:underline cursor-pointer text-xs"
                >
                  {row.projectName}
                </span>
                <span className="text-[11px] text-muted-foreground">{row.projectAddress || '—'}</span>
              </div>
            )
          },
        }),
        // 3. Client Name
        columnHelper.accessor('clientName', {
          header: 'Client Name',
          cell: (info) => (
            <div className="text-xs font-medium text-foreground whitespace-nowrap">
              {info.getValue()}
            </div>
          ),
        }),
        // 4. Project Manager
        columnHelper.accessor('projectManager', {
          header: 'Project Manager',
          cell: (info) => {
            const pm = info.getValue()
            return (
              <div className="flex items-center gap-1.5 whitespace-nowrap">
                <span className="flex h-5 w-5 items-center justify-center rounded-full bg-primary/10 text-[10px] font-bold text-primary">
                  {pm?.initials || 'PM'}
                </span>
                <span className="text-xs text-foreground">{pm?.name}</span>
              </div>
            )
          },
        }),
        // 5. Project Type
        columnHelper.accessor('projectType', {
          header: 'Project Type',
          cell: (info) => (
            <Badge variant="secondary" className="text-[11px] font-medium whitespace-nowrap">
              {info.getValue()}
            </Badge>
          ),
        }),
        // 6. Sessions
        columnHelper.accessor('sessionsCount', {
          header: 'Sessions',
          cell: (info) => (
            <span className="text-xs font-semibold tabular-nums text-foreground">
              {info.getValue()} Session(s)
            </span>
          ),
        }),
        // 7. Tasks
        columnHelper.accessor('tasksCount', {
          header: 'Tasks',
          cell: (info) => {
            const val = info.getValue()
            return (
              <span className="text-xs tabular-nums text-muted-foreground">
                {val.completed}/{val.total} done
              </span>
            )
          },
        }),
        // 8. Deliverables
        columnHelper.accessor('deliverablesCount', {
          header: 'Deliverables',
          cell: (info) => {
            const val = info.getValue()
            return (
              <span className="text-xs tabular-nums text-muted-foreground">
                {val.ready}/{val.total} ready
              </span>
            )
          },
        }),
        // 9. Status
        columnHelper.accessor('status', {
          header: 'Status',
          cell: (info) => {
            const val = info.getValue() as ProjectStatus
            return (
              <Badge variant="outline" className="text-[11px] font-medium whitespace-nowrap">
                {val}
              </Badge>
            )
          },
        }),
        // 10. Stage
        columnHelper.accessor('stage', {
          header: 'Stage',
          cell: (info) => {
            const val = info.getValue() as ProjectStage
            const color = stageColorMap.get(val.toLowerCase())
            return (
              <Badge
                variant="outline"
                className="text-[11px] font-semibold whitespace-nowrap"
                style={
                  color
                    ? {
                        backgroundColor: `${color}20`,
                        borderColor: `${color}60`,
                        color: color,
                      }
                    : undefined
                }
              >
                {val}
              </Badge>
            )
          },
        }),
        // 11. Start Date
        columnHelper.accessor('startDate', {
          header: 'Start Date',
          cell: (info) => (
            <span className="text-xs text-muted-foreground tabular-nums whitespace-nowrap">
              {info.getValue() ? info.getValue().replace('T', ' ') : '—'}
            </span>
          ),
        }),
        // 12. End Date
        columnHelper.accessor('endDate', {
          header: 'End Date',
          cell: (info) => (
            <span className="text-xs text-muted-foreground tabular-nums whitespace-nowrap">
              {info.getValue() ? info.getValue().replace('T', ' ') : '—'}
            </span>
          ),
        }),
        // 13. Created Date
        columnHelper.accessor('createdDate', {
          header: 'Created Date',
          cell: (info) => (
            <span className="text-xs text-muted-foreground tabular-nums whitespace-nowrap">
              {info.getValue()}
            </span>
          ),
        }),
      ]),
    [onEditProject]
  )

  // Filters for tanstack table
  const filterConfigs: DataTableFilterConfig[] = React.useMemo(
    () => [
      {
        columnId: 'stage',
        title: 'Stage',
        type: 'select',
        options: dynamicStages.map((s) => ({ label: s, value: s })),
      },
      {
        columnId: 'status',
        title: 'Status',
        type: 'select',
        options: STATUSES.filter((s) => s !== 'All').map((s) => ({ label: s, value: s })),
      },
    ],
    [dynamicStages]
  )

  // Bulk actions
  const bulkActions: DataTableBulkAction<Project>[] = [
    {
      label: 'Delete Projects',
      variant: 'destructive',
      action: (selectedRows: Project[], clearSelection: () => void) => {
        selectedRows.forEach((r) => onDeleteProject(r.id))
        clearSelection()
        toast.success(`Deleted ${selectedRows.length} project(s)`)
      },
    },
    {
      label: 'Mark as Completed',
      action: (selectedRows: Project[], clearSelection: () => void) => {
        selectedRows.forEach((r) => onStageChange(r.id, 'Closure'))
        clearSelection()
        toast.success(`Moved ${selectedRows.length} project(s) to Closure`)
      },
    },
  ]

  // Row actions
  const rowActions: DataTableRowAction<Project>[] = [
    {
      label: 'Open Project Workspace',
      icon: IconEdit,
      action: (row: Project) => onEditProject(row),
    },
    {
      label: 'Delete Project',
      icon: IconTrash,
      variant: 'destructive',
      action: (row: Project) => {
        onDeleteProject(row.id)
        toast.success(`Project ${row.projectName} deleted`)
      },
    },
  ]

  return (
    <div className="space-y-6">
      {/* Header section with KPIs and View Toggle */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">Projects Workspace</h1>
          <p className="text-xs sm:text-sm text-muted-foreground mt-0.5">
            Production schedules, shoot management, post-production editing, and deliveries.
          </p>
        </div>

        <div className="flex items-center gap-2.5">
          {/* Table / Kanban Toggle */}
          <div className="inline-flex rounded-lg border bg-muted/30 p-1">
            <button
              type="button"
              onClick={() => setViewMode('table')}
              className={cn(
                'flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-md transition shadow-xs',
                viewMode === 'table'
                  ? 'bg-background text-foreground shadow-xs'
                  : 'text-muted-foreground hover:text-foreground'
              )}
            >
              <IconTable className="h-3.5 w-3.5" />
              <span>Table</span>
            </button>
            <button
              type="button"
              onClick={() => setViewMode('kanban')}
              className={cn(
                'flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-md transition',
                viewMode === 'kanban'
                  ? 'bg-background text-foreground shadow-xs'
                  : 'text-muted-foreground hover:text-foreground'
              )}
            >
              <IconLayoutKanban className="h-3.5 w-3.5" />
              <span>Kanban</span>
            </button>
          </div>

          <Button
            onClick={onAddNewProject}
            className="bg-emerald-600 hover:bg-emerald-700 text-white font-medium gap-1.5 shadow-xs"
          >
            <IconPlus className="h-4 w-4" />
            New Project
          </Button>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3 sm:gap-4">
        <Card className="border shadow-xs bg-card">
          <CardContent className="p-4 flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-muted-foreground">Total Projects</p>
              <h3 className="text-xl sm:text-2xl font-bold mt-1 tabular-nums">{metrics.total}</h3>
            </div>
            <div className="h-10 w-10 rounded-xl bg-primary/10 text-primary flex items-center justify-center">
              <IconBriefcase className="h-5 w-5" />
            </div>
          </CardContent>
        </Card>

        <Card className="border shadow-xs bg-card">
          <CardContent className="p-4 flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-muted-foreground">Active Shoots</p>
              <h3 className="text-xl sm:text-2xl font-bold mt-1 tabular-nums text-amber-600 dark:text-amber-400">
                {metrics.activeShoots}
              </h3>
            </div>
            <div className="h-10 w-10 rounded-xl bg-amber-500/10 text-amber-600 dark:text-amber-400 flex items-center justify-center">
              <IconCamera className="h-5 w-5" />
            </div>
          </CardContent>
        </Card>

        <Card className="border shadow-xs bg-card">
          <CardContent className="p-4 flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-muted-foreground">Post-Production (Editing)</p>
              <h3 className="text-xl sm:text-2xl font-bold mt-1 tabular-nums text-purple-600 dark:text-purple-400">
                {metrics.inEditing}
              </h3>
            </div>
            <div className="h-10 w-10 rounded-xl bg-purple-500/10 text-purple-600 dark:text-purple-400 flex items-center justify-center">
              <IconPhoto className="h-5 w-5" />
            </div>
          </CardContent>
        </Card>

        <Card className="border shadow-xs bg-card">
          <CardContent className="p-4 flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-muted-foreground">Delivered / Closed</p>
              <h3 className="text-xl sm:text-2xl font-bold mt-1 tabular-nums text-emerald-600 dark:text-emerald-400">
                {metrics.delivered}
              </h3>
            </div>
            <div className="h-10 w-10 rounded-xl bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 flex items-center justify-center">
              <IconChecklist className="h-5 w-5" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Stage Navigation Filter Bar */}
      <div className="border-b">
        <div className="flex items-center gap-1 overflow-x-auto scrollbar-none -mb-px">
          {['All', ...dynamicStages].map((stg) => {
            const count =
              stg === 'All'
                ? projects.length
                : projects.filter((p) => p.stage === stg).length
            const isActive = stageFilter === stg

            return (
              <button
                key={stg}
                type="button"
                onClick={() => setStageFilter(stg)}
                className={cn(
                  'flex items-center gap-2 px-4 py-2.5 text-xs font-semibold rounded-t-lg transition border-b-2 whitespace-nowrap',
                  isActive
                    ? 'border-primary text-primary bg-accent/40'
                    : 'border-transparent text-muted-foreground hover:text-foreground hover:bg-muted/40'
                )}
              >
                <span>{stg === 'All' ? 'All Stages' : stg}</span>
                <span
                  className={cn(
                    'rounded-full px-2 py-0.2 text-[10px] tabular-nums font-bold',
                    isActive ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground'
                  )}
                >
                  {count}
                </span>
              </button>
            )
          })}
        </div>
      </div>

      {/* Secondary Status Filters */}
      <div className="flex items-center gap-2 overflow-x-auto py-1 scrollbar-none">
        <span className="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider shrink-0 mr-1">
          Status:
        </span>
        {STATUSES.map((st) => (
          <button
            key={st}
            type="button"
            onClick={() => setStatusFilter(st)}
            className={cn(
              'rounded-full px-3 py-1 text-xs font-medium transition whitespace-nowrap border',
              statusFilter === st
                ? 'bg-foreground text-background border-foreground font-semibold'
                : 'bg-card text-muted-foreground border-border hover:bg-muted hover:text-foreground'
            )}
          >
            {st === 'All' ? 'All Status' : st}
          </button>
        ))}
      </div>

      {/* Table vs Kanban View */}
      {viewMode === 'table' ? (
        <Card className="border bg-card shadow-xs">
          <CardContent className="p-3 sm:p-5">
            <DataTable
              columns={columns}
              data={filteredProjects}
              searchableColumns={['projectName', 'clientName', 'id', 'projectType']}
              searchPlaceholder="Search projects by client, ID, or project name..."
              filterConfigs={filterConfigs}
              bulkActions={bulkActions}
              rowActions={rowActions}
              renderDetailDrawer={(row, close) => (
                <div className="space-y-6">
                  <div className="rounded-lg border bg-muted/30 p-4 space-y-3">
                    <div className="flex items-center justify-between">
                      <span className="text-xs font-medium text-muted-foreground">Project Stage</span>
                      <Badge variant="outline" className="font-semibold text-xs">
                        {row.stage}
                      </Badge>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-xs font-medium text-muted-foreground">Status</span>
                      <Badge variant="secondary" className="text-xs">
                        {row.status}
                      </Badge>
                    </div>
                  </div>

                  <div className="space-y-3">
                    <h4 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                      Overview
                    </h4>
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between py-1 border-b">
                        <span className="text-muted-foreground">Project</span>
                        <span className="font-medium">{row.projectName}</span>
                      </div>
                      <div className="flex justify-between py-1 border-b">
                        <span className="text-muted-foreground">Client</span>
                        <span>{row.clientName}</span>
                      </div>
                      <div className="flex justify-between py-1 border-b">
                        <span className="text-muted-foreground">Type</span>
                        <span>{row.projectType}</span>
                      </div>
                      <div className="flex justify-between py-1 border-b">
                        <span className="text-muted-foreground">Manager</span>
                        <span>{row.projectManager.name}</span>
                      </div>
                      <div className="flex justify-between py-1 border-b">
                        <span className="text-muted-foreground">Timeline</span>
                        <span className="tabular-nums text-xs">
                          {row.startDate.replace('T', ' ')} → {row.endDate.replace('T', ' ')}
                        </span>
                      </div>
                      <div className="flex justify-between py-1 border-b">
                        <span className="text-muted-foreground">Sessions</span>
                        <span className="tabular-nums">{row.sessionsCount} sessions</span>
                      </div>
                    </div>
                  </div>

                  <div className="flex gap-2 pt-4">
                    <Button
                      className="flex-1 bg-emerald-600 hover:bg-emerald-700 text-white"
                      onClick={() => {
                        close()
                        onEditProject(row)
                      }}
                    >
                      Open Full Project Workspace
                    </Button>
                    <Button variant="outline" onClick={close}>
                      Close
                    </Button>
                  </div>
                </div>
              )}
              detailDrawerTitle={(row) => `Project - ${row.projectName}`}
              detailDrawerDescription={(row) =>
                `Client: ${row.clientName} • Stage: ${row.stage} • Status: ${row.status}`
              }
              enableExport={true}
              exportFilename="photography-projects.csv"
            />
          </CardContent>
        </Card>
      ) : (
        <ProjectKanbanView
          projects={filteredProjects}
          onProjectClick={(p) => onEditProject(p)}
          onEditProject={onEditProject}
          onDeleteProject={onDeleteProject}
          onStageChange={onStageChange}
          onAddNewProject={onAddNewProject}
        />
      )}
    </div>
  )
}
