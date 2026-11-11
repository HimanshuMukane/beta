import * as React from 'react'
import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@repo/ui/components/ui/dialog'
import { Label } from '@repo/ui/components/ui/label'
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
  IconCheck,
  IconPlus,
} from "@repo/ui/icons/tabler-icons"
import {
  ProductionTaskDrawer,
  type TaskDetailData,
} from './ProductionTaskDrawer'

export interface ProductionCategory {
  name: string
  tasks: string[]
}

const PRE_PRODUCTION_CATEGORIES: ProductionCategory[] = [
  {
    name: 'Client Information',
    tasks: ['Collect details', 'Event details', 'Notes'],
  },
  {
    name: 'Creative Planning',
    tasks: ['Discuss style & theme', 'Create mood board', 'Prepare shot list'],
  },
  {
    name: 'Schedule & Crew',
    tasks: ['Create Sessions', 'Assign Crew', 'Schedule pre-shoot meeting'],
  },
  {
    name: 'Location & Permissions',
    tasks: [
      'Visit venue',
      'Identify shooting spots & lighting conditions',
      'Secure permits',
    ],
  },
  {
    name: 'Equipment Prep',
    tasks: [
      'Make list of required gear',
      'Charge batteries & pack spares',
      'Format memory cards',
      'Pack backup',
    ],
  },
  {
    name: 'Logistics',
    tasks: [
      'Arrange crew transport',
      'Refreshments for long events',
      'Prepare emergency kit',
    ],
  },
  {
    name: 'Team Briefing',
    tasks: [
      'Share event flow & shot list',
      'Assign roles & responsibilities',
      'Double-check equipment readiness',
    ],
  },
]

const POST_PRODUCTION_CATEGORIES: ProductionCategory[] = [
  {
    name: 'File Management',
    tasks: ['Transfer all File', 'Upload to cloud storage', 'Organize files'],
  },
  {
    name: 'Culling & Selection',
    tasks: [
      'Review all photos',
      'Remove unwanted shots',
      'Shortlist best images',
      'Tag highlight photos',
    ],
  },
  {
    name: 'Editing & Retouching',
    tasks: [
      'Color correction',
      'Exposure adjustments',
      'Retouch faces',
      'Apply styles',
    ],
  },
  {
    name: 'Video Editing',
    tasks: [
      'Trim & arrange clips',
      'Add effects',
      'Sync music/audio',
      'Color grade video',
    ],
  },
  {
    name: 'Client Selection',
    tasks: [
      'Create Gallery & Upload',
      'Share',
      'Download selected images',
    ],
  },
  {
    name: 'Album Design',
    tasks: [
      'Design album layout',
      'Send preview to client',
      'Collect feedback',
      'Album Approved',
      'Share Album Preview',
    ],
  },
  {
    name: 'Digital Album',
    tasks: [],
  },
  {
    name: 'Printing & Delivery',
    tasks: [
      'Link Lab',
      'Share Approved Sources to lab',
      'Select Design from lab catalog',
      'Place Order',
      'Source Downloaded',
      'Printing Completed',
      'Items received from lab',
      'Delivered Album',
    ],
  },
  {
    name: 'Lab Progress',
    tasks: [],
  },
  {
    name: 'Closure',
    tasks: [
      'Collect pending payment',
      'Request testimonials',
      'Archive final files',
    ],
  },
]

interface ProductionMatrixViewProps {
  projectId: string
  projectName: string
}

export function ProductionMatrixView({
  projectId: _projectId,
  projectName,
}: ProductionMatrixViewProps) {
  const [activePhase, setActivePhase] = React.useState<
    'Pre-production' | 'Post-production'
  >('Pre-production')

  const [activeDrawerTask, setActiveDrawerTask] =
    React.useState<TaskDetailData | null>(null)
  const [isDrawerOpen, setIsDrawerOpen] = React.useState(false)

  // Track task states in local state keyed by `${category}_${taskName}`
  const [taskStateMap, setTaskStateMap] = React.useState<
    Record<string, TaskDetailData>
  >(() => {
    const map: Record<string, TaskDetailData> = {}
    const initCategory = (cats: ProductionCategory[], phase: 'Pre-production' | 'Post-production') => {
      cats.forEach((cat) => {
        cat.tasks.forEach((tName) => {
          const key = `${phase}_${cat.name}_${tName}`
          map[key] = {
            taskId: key,
            taskName: tName,
            categoryName: cat.name,
            phase,
            status: 'Pending',
            priority: 'Medium',
            dueDate: '2026-11-20',
            assignee: 'Lead Photographer',
            notes: '',
            checklist: [
              { id: '1', text: 'Verify requirements', done: false },
              { id: '2', text: 'Sign-off with client/lead', done: false },
            ],
          }
        })
      })
    }
    initCategory(PRE_PRODUCTION_CATEGORIES, 'Pre-production')
    initCategory(POST_PRODUCTION_CATEGORIES, 'Post-production')
    return map
  })

  const [isAddProjectModalOpen, setIsAddProjectModalOpen] = React.useState(false)
  const [newProjectPhase, setNewProjectPhase] = React.useState('Pre-production')

  const currentCategories =
    activePhase === 'Pre-production'
      ? PRE_PRODUCTION_CATEGORIES
      : POST_PRODUCTION_CATEGORIES

  const handleCellClick = (catName: string, taskName: string) => {
    const key = `${activePhase}_${catName}_${taskName}`
    const existing =
      taskStateMap[key] || {
        taskId: key,
        taskName,
        categoryName: catName,
        phase: activePhase,
        status: 'Pending',
        priority: 'Medium',
        dueDate: '2026-11-20',
        assignee: 'Lead Photographer',
        notes: '',
        checklist: [],
      }

    setActiveDrawerTask(existing)
    setIsDrawerOpen(true)
  }

  const handleSaveDrawerTask = (updated: TaskDetailData) => {
    setTaskStateMap((prev) => ({
      ...prev,
      [updated.taskId]: updated,
    }))
  }

  // Calculate phase summary statistics
  const phaseStats = React.useMemo(() => {
    let total = 0
    let completed = 0
    let inProgress = 0

    currentCategories.forEach((cat) => {
      cat.tasks.forEach((t) => {
        total++
        const key = `${activePhase}_${cat.name}_${t}`
        const st = taskStateMap[key]?.status
        if (st === 'Completed') completed++
        if (st === 'In Progress') inProgress++
      })
    })

    return { total, completed, inProgress }
  }, [currentCategories, activePhase, taskStateMap])

  return (
    <div className="space-y-4">
      {/* Matrix Controls & Progress Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-muted/20 p-3 rounded-xl border">
        <div className="flex items-center gap-3">
          <div className="flex items-center rounded-lg bg-card border p-1 shadow-xs">
            <button
              type="button"
              onClick={() => setActivePhase('Pre-production')}
              className={cn(
                'px-4 py-1.5 text-xs font-bold rounded-md transition cursor-pointer',
                activePhase === 'Pre-production'
                  ? 'bg-foreground text-background shadow-xs font-extrabold'
                  : 'text-muted-foreground hover:text-foreground'
              )}
            >
              Pre Production Tasks
            </button>
            <button
              type="button"
              onClick={() => setActivePhase('Post-production')}
              className={cn(
                'px-4 py-1.5 text-xs font-bold rounded-md transition cursor-pointer',
                activePhase === 'Post-production'
                  ? 'bg-foreground text-background shadow-xs font-extrabold'
                  : 'text-muted-foreground hover:text-foreground'
              )}
            >
              Post-Production Tasks
            </button>
          </div>

          <div className="hidden md:flex items-center gap-2 text-xs text-muted-foreground">
            <span className="font-semibold text-foreground">
              {phaseStats.completed}/{phaseStats.total}
            </span>{' '}
            done
            <span className="h-1.5 w-16 bg-muted rounded-full overflow-hidden inline-block">
              <span
                className="h-full bg-emerald-500 block transition-all"
                style={{
                  width: `${
                    phaseStats.total > 0
                      ? (phaseStats.completed / phaseStats.total) * 100
                      : 0
                  }%`,
                }}
              />
            </span>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <Button
            size="sm"
            onClick={() => setIsAddProjectModalOpen(true)}
            className="bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs h-8 shadow-xs"
          >
            <IconPlus className="h-3.5 w-3.5 mr-1" />
            Add Projects
          </Button>
        </div>
      </div>

      {/* Production Matrix Table Matching Image 1 with Horizontal Scroll & Sticky Top Category Header */}
      <div className="rounded-xl border border-border/80 bg-card shadow-xs overflow-hidden">
        <div className="overflow-x-auto scrollbar-thin">
          <table className="w-full border-collapse text-left min-w-[1200px]">
            {/* Top Categories Row (Header Level 1) */}
            <thead>
              <tr className="bg-slate-700/80 dark:bg-slate-800 text-white border-b border-border">
                {currentCategories.map((cat) => {
                  const colSpan = Math.max(1, cat.tasks.length)
                  return (
                    <th
                      key={cat.name}
                      colSpan={colSpan}
                      className="px-3 py-2 text-center text-[11px] font-bold tracking-wider uppercase border-r border-slate-600/60 last:border-r-0 select-none"
                    >
                      {cat.name}
                    </th>
                  )
                })}
              </tr>

              {/* Sub-Tasks Header Row (Header Level 2) */}
              <tr className="bg-card border-b border-border/70 text-muted-foreground text-xs font-semibold">
                {currentCategories.flatMap((cat) => {
                  if (cat.tasks.length === 0) {
                    return (
                      <th
                        key={`${cat.name}-empty`}
                        className="px-3 py-2.5 text-center text-xs text-muted-foreground/60 border-r last:border-r-0 font-normal italic min-w-[130px]"
                      >
                        No tasks listed
                      </th>
                    )
                  }
                  return cat.tasks.map((taskName) => (
                    <th
                      key={`${cat.name}_${taskName}`}
                      className="px-3 py-2.5 text-center text-xs text-foreground font-semibold border-r last:border-r-0 min-w-[145px] hover:bg-muted/40 transition"
                    >
                      <span className="line-clamp-2">{taskName}</span>
                    </th>
                  ))
                })}
              </tr>
            </thead>

            {/* Matrix Data Rows (Clicking any cell opens the interactive task configuration drawer) */}
            <tbody>
              {/* Project Primary Execution Row */}
              <tr className="hover:bg-muted/10 transition group">
                {currentCategories.flatMap((cat) => {
                  if (cat.tasks.length === 0) {
                    return (
                      <td
                        key={`${cat.name}-cell-empty`}
                        className="px-3 py-6 text-center text-xs text-muted-foreground/40 border-r last:border-r-0 border-b bg-muted/5 italic"
                      >
                        —
                      </td>
                    )
                  }

                  return cat.tasks.map((taskName) => {
                    const key = `${activePhase}_${cat.name}_${taskName}`
                    const taskData = taskStateMap[key]
                    const status = taskData?.status || 'Pending'

                    const statusStyles: Record<string, string> = {
                      Completed:
                        'bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-950/30 dark:text-emerald-400 dark:border-emerald-800',
                      'In Progress':
                        'bg-blue-50 text-blue-700 border-blue-200 dark:bg-blue-950/30 dark:text-blue-400 dark:border-blue-800',
                      Blocked:
                        'bg-rose-50 text-rose-700 border-rose-200 dark:bg-rose-950/30 dark:text-rose-400 dark:border-rose-800',
                      Pending:
                        'bg-muted/30 text-muted-foreground border-transparent hover:border-border',
                    }

                    return (
                      <td
                        key={`${cat.name}_${taskName}_cell`}
                        onClick={() => handleCellClick(cat.name, taskName)}
                        className="p-2 border-r last:border-r-0 border-b align-top cursor-pointer hover:bg-primary/5 transition"
                      >
                        <div
                          className={cn(
                            'rounded-lg border p-2 flex flex-col justify-between min-h-[64px] transition-all hover:shadow-xs hover:border-primary/50',
                            statusStyles[status]
                          )}
                        >
                          <div className="flex items-center justify-between gap-1 mb-1">
                            <span className="text-[10px] font-bold uppercase tracking-wider">
                              {status}
                            </span>
                            {status === 'Completed' ? (
                              <IconCheck className="h-3.5 w-3.5 text-emerald-600" />
                            ) : (
                              <span className="h-2 w-2 rounded-full bg-primary/40" />
                            )}
                          </div>

                          <div className="text-[11px] font-medium truncate text-foreground/80">
                            {taskData?.assignee || 'Unassigned'}
                          </div>
                        </div>
                      </td>
                    )
                  })
                })}
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      {/* Production Task Drawer */}
      <ProductionTaskDrawer
        open={isDrawerOpen}
        onOpenChange={setIsDrawerOpen}
        task={activeDrawerTask}
        onSaveTask={handleSaveDrawerTask}
      />

      {/* Add Project Dialog */}
      <Dialog
        open={isAddProjectModalOpen}
        onOpenChange={setIsAddProjectModalOpen}
      >
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Add Production Matrix Project</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-2">
            <div className="space-y-1.5">
              <Label>Project Context</Label>
              <Input value={projectName} disabled className="bg-muted text-xs" />
            </div>
            <div className="space-y-1.5">
              <Label>Starting Workflow Phase</Label>
              <Select
                value={newProjectPhase}
                onValueChange={setNewProjectPhase}
              >
                <SelectTrigger className="text-xs">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Pre-production">Pre-production Phase</SelectItem>
                  <SelectItem value="Post-production">Post-production Phase</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <p className="text-xs text-muted-foreground">
              Initializes the photography production matrix workflows for this event.
            </p>
          </div>
          <DialogFooter>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setIsAddProjectModalOpen(false)}
            >
              Cancel
            </Button>
            <Button
              size="sm"
              className="bg-emerald-600 hover:bg-emerald-700 text-white"
              onClick={() => {
                toast.success('Production matrix initialized for project')
                setIsAddProjectModalOpen(false)
              }}
            >
              Initialize Tasks
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

