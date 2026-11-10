import * as React from 'react'
import { Button } from '@repo/ui/components/ui/button'
import { Badge } from '@repo/ui/components/ui/badge'
import { Input } from '@repo/ui/components/ui/input'
import { Checkbox } from '@repo/ui/components/ui/checkbox'
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
  IconChevronDown,
  IconChevronRight,
  IconSearch,
  IconListCheck,
} from '@tabler/icons-react'
import type { Project } from '../../types/project'
import {
  ProductionTaskDrawer,
  type TaskDetailData,
} from '../../components/ProductionTaskDrawer'

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

interface UserTasksPageProps {
  projects: Project[]
  onUpdateProject?: (project: Project) => void
}

export function UserTasksPage({
  projects,
}: UserTasksPageProps) {
  const [activeTab, setActiveTab] = React.useState<'pre' | 'post'>('pre')
  const [searchQuery, setSearchQuery] = React.useState('')
  const [collapsedMonths, setCollapsedMonths] = React.useState<Record<string, boolean>>({})

  // Completed task tracker: `${projectId}-${taskName}` -> boolean
  const [completedTasks, setCompletedTasks] = React.useState<Record<string, boolean>>({})

  // Drawer task details state
  const [drawerOpen, setDrawerOpen] = React.useState(false)
  const [selectedTaskDetail, setSelectedTaskDetail] = React.useState<TaskDetailData | null>(null)
  const [activeProjectId, setActiveProjectId] = React.useState<string | null>(null)

  // Custom task store per category
  const [preCategories, setPreCategories] = React.useState<ProductionCategory[]>(PRE_PRODUCTION_CATEGORIES)
  const [postCategories, setPostCategories] = React.useState<ProductionCategory[]>(POST_PRODUCTION_CATEGORIES)

  // Add Task Modal state
  const [addTaskModalOpen, setAddTaskModalOpen] = React.useState(false)
  const [targetCategory, setTargetCategory] = React.useState<string>('')
  const [newTaskName, setNewTaskName] = React.useState('')

  const activeCategories = activeTab === 'pre' ? preCategories : postCategories

  const toggleMonth = (monthKey: string) => {
    setCollapsedMonths((prev) => ({ ...prev, [monthKey]: !prev[monthKey] }))
  }

  const handleToggleTaskStatus = (projId: string, taskName: string, e: React.MouseEvent) => {
    e.stopPropagation()
    const key = `${projId}::${taskName}`
    const nextState = !completedTasks[key]
    setCompletedTasks((prev) => ({ ...prev, [key]: nextState }))
    toast.success(`Task "${taskName}" marked as ${nextState ? 'Completed' : 'Pending'}`)
  }

  const handleOpenTaskDrawer = (proj: Project, taskName: string) => {
    setActiveProjectId(proj.id)
    setSelectedTaskDetail({
      taskId: `${proj.id}-${taskName}`,
      taskName,
      categoryName: activeTab === 'pre' ? 'Pre-production' : 'Post-production',
      phase: activeTab === 'pre' ? 'Pre-production' : 'Post-production',
      status: completedTasks[`${proj.id}::${taskName}`] ? 'Completed' : 'Pending',
      priority: 'Medium',
      dueDate: proj.startDate || '2026-11-20',
      assignee: proj.projectManager?.name || 'Rohan Mehta',
      notes: `Task for project: ${proj.projectName}. Complete and verify workflow checklist before shoot milestones.`,
      checklist: [
        { id: '1', text: 'Review project requirements', done: false },
        { id: '2', text: 'Confirm schedule & gear', done: false },
      ],
    })
    setDrawerOpen(true)
  }

  const handleSaveDrawerTask = (updated: TaskDetailData) => {
    if (activeProjectId && updated.taskName) {
      const key = `${activeProjectId}::${updated.taskName}`
      setCompletedTasks((prev) => ({
        ...prev,
        [key]: updated.status === 'Completed',
      }))
    }
    setDrawerOpen(false)
    toast.success('Task details updated')
  }

  const handleAddTask = () => {
    if (!newTaskName.trim() || !targetCategory) {
      toast.error('Please specify both a task name and target group')
      return
    }

    const updater = activeTab === 'pre' ? setPreCategories : setPostCategories
    updater((prev) =>
      prev.map((cat) =>
        cat.name === targetCategory
          ? { ...cat, tasks: [...cat.tasks, newTaskName.trim()] }
          : cat
      )
    )

    toast.success(`Added task "${newTaskName.trim()}" to ${targetCategory}`)
    setNewTaskName('')
    setAddTaskModalOpen(false)
  }

  // Filter projects by search
  const filteredProjects = React.useMemo(() => {
    return projects.filter((p) => {
      if (!searchQuery.trim()) return true
      const q = searchQuery.toLowerCase()
      return (
        p.projectName?.toLowerCase().includes(q) ||
        p.clientName?.toLowerCase().includes(q) ||
        p.projectType?.toLowerCase().includes(q)
      )
    })
  }, [projects, searchQuery])

  // Group projects by Month
  const groupedProjects = React.useMemo(() => {
    const groups: Record<string, Project[]> = {}
    filteredProjects.forEach((p) => {
      let groupKey = 'November 2026'
      const dateStr = p.startDate || p.createdDate
      if (dateStr && dateStr.includes('-')) {
        const [year, month] = dateStr.split('-')
        const monthNames = [
          'January', 'February', 'March', 'April', 'May', 'June',
          'July', 'August', 'September', 'October', 'November', 'December'
        ]
        const mIdx = parseInt(month, 10) - 1
        groupKey = `${monthNames[mIdx] || 'November'} ${year || '2026'}`
      }
      if (!groups[groupKey]) groups[groupKey] = []
      groups[groupKey].push(p)
    })
    return groups
  }, [filteredProjects])

  return (
    <div className="space-y-4">
      {/* Top Header matching Screenshot 2 with enhanced UX */}
      <div className="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3 border-b pb-3">
        {/* Two Tabs: Pre Production Tasks | Post-Production Tasks */}
        <div className="flex items-center gap-1">
          <button
            type="button"
            onClick={() => setActiveTab('pre')}
            className={cn(
              'px-4 py-2 text-xs font-semibold rounded-lg transition-all border-b-2',
              activeTab === 'pre'
                ? 'border-primary text-foreground bg-muted/40 font-bold'
                : 'border-transparent text-muted-foreground hover:text-foreground hover:bg-muted/20'
            )}
          >
            Pre Production Tasks
          </button>
          <button
            type="button"
            onClick={() => setActiveTab('post')}
            className={cn(
              'px-4 py-2 text-xs font-semibold rounded-lg transition-all border-b-2',
              activeTab === 'post'
                ? 'border-primary text-foreground bg-muted/40 font-bold'
                : 'border-transparent text-muted-foreground hover:text-foreground hover:bg-muted/20'
            )}
          >
            Post-Production Tasks
          </button>
        </div>

        {/* Right Actions: Add Project task | Filter */}
        <div className="flex items-center gap-2">
          <div className="relative w-48 sm:w-64">
            <IconSearch className="absolute left-2.5 top-2.5 h-3.5 w-3.5 text-muted-foreground" />
            <Input
              placeholder="Search projects..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-8 h-8 text-xs bg-card"
            />
          </div>

          <Button
            size="sm"
            onClick={() => {
              setTargetCategory(activeCategories[0]?.name || '')
              setAddTaskModalOpen(true)
            }}
            className="h-8 text-xs bg-emerald-600 hover:bg-emerald-700 text-white font-semibold rounded-full px-4 gap-1.5 shadow-2xs"
          >
            <IconPlus className="h-3.5 w-3.5" />
            <span>Add Task</span>
          </Button>
        </div>
      </div>

      {/* Production Matrix Table */}
      <div className="space-y-6">
        {Object.keys(groupedProjects).length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-center border border-dashed rounded-xl bg-card">
            <IconListCheck className="h-10 w-10 text-muted-foreground/30 mb-2" />
            <p className="text-sm font-semibold text-foreground">No projects found</p>
            <p className="text-xs text-muted-foreground max-w-sm mt-1">
              There are no projects or tasks matching your query.
            </p>
          </div>
        ) : (
          Object.entries(groupedProjects).map(([monthKey, projsInMonth]) => {
            const isCollapsed = collapsedMonths[monthKey]
            return (
              <div key={monthKey} className="rounded-xl border border-border/80 bg-card overflow-hidden shadow-2xs">
                {/* Month Group Header accordion */}
                <div
                  onClick={() => toggleMonth(monthKey)}
                  className="flex items-center justify-between px-4 py-2.5 bg-muted/40 hover:bg-muted/60 transition cursor-pointer border-b"
                >
                  <div className="flex items-center gap-2">
                    {isCollapsed ? (
                      <IconChevronRight className="h-4 w-4 text-muted-foreground" />
                    ) : (
                      <IconChevronDown className="h-4 w-4 text-muted-foreground" />
                    )}
                    <div>
                      <span className="text-xs font-bold text-foreground">{monthKey}</span>
                      <span className="ml-2 text-[11px] text-muted-foreground">
                        {projsInMonth.length} Project{projsInMonth.length > 1 ? 's' : ''}
                      </span>
                    </div>
                  </div>
                  <Badge variant="outline" className="text-[10px] font-semibold bg-background">
                    {activeTab === 'pre' ? 'Pre-Production' : 'Post-Production'}
                  </Badge>
                </div>

                {!isCollapsed && (
                  <div className="overflow-x-auto">
                    <table className="w-full text-left border-collapse min-w-[900px]">
                      {/* Top Header: Categories */}
                      <thead>
                        <tr className="border-b bg-muted/60">
                          <th className="w-64 p-3 text-xs font-bold text-foreground border-r sticky left-0 bg-muted/90 z-10">
                            Projects
                          </th>
                          {activeCategories.map((cat) => (
                            <th
                              key={cat.name}
                              colSpan={Math.max(cat.tasks.length, 1)}
                              className="px-3 py-2 text-center text-[10px] font-bold uppercase tracking-wider text-muted-foreground border-r border-border/60 bg-muted/40"
                            >
                              {cat.name}
                            </th>
                          ))}
                        </tr>

                        {/* Sub Header: Individual Tasks */}
                        <tr className="border-b bg-muted/20 text-xs">
                          <th className="p-2 border-r sticky left-0 bg-muted/50 z-10 text-[11px] font-medium text-muted-foreground">
                            Name & Dates
                          </th>
                          {activeCategories.flatMap((cat) =>
                            cat.tasks.length === 0 ? (
                              <th
                                key={`${cat.name}-empty`}
                                className="px-3 py-2 text-center text-[11px] font-medium text-muted-foreground/60 border-r"
                              >
                                —
                              </th>
                            ) : (
                              cat.tasks.map((task) => (
                                <th
                                  key={`${cat.name}-${task}`}
                                  className="px-3 py-2 text-center text-[11px] font-medium text-muted-foreground border-r border-border/60 whitespace-nowrap min-w-[120px]"
                                >
                                  {task}
                                </th>
                              ))
                            )
                          )}
                        </tr>
                      </thead>

                      {/* Project rows */}
                      <tbody className="divide-y divide-border/60 text-xs">
                        {projsInMonth.map((project) => (
                          <tr key={project.id} className="hover:bg-muted/10 transition">
                            {/* Project Name Cell (pinned) */}
                            <td className="p-3 border-r sticky left-0 bg-card z-10 shadow-2xs font-medium">
                              <div className="flex items-center gap-2">
                                <Checkbox id={`proj-${project.id}`} />
                                <div className="min-w-0">
                                  <label
                                    htmlFor={`proj-${project.id}`}
                                    className="font-bold text-xs text-foreground cursor-pointer block truncate hover:text-primary transition"
                                  >
                                    {project.projectName}
                                  </label>
                                  <span className="text-[10px] text-muted-foreground block">
                                    {project.startDate ? project.startDate.slice(5) : '08-11'}
                                    {project.endDate ? ` to ${project.endDate.slice(5)}` : ''}
                                  </span>
                                </div>
                              </div>
                            </td>

                            {/* Task Cells */}
                            {activeCategories.flatMap((cat) =>
                              cat.tasks.length === 0 ? (
                                <td
                                  key={`${project.id}-${cat.name}-empty`}
                                  className="p-2 border-r text-center text-muted-foreground/30"
                                >
                                  —
                                </td>
                              ) : (
                                cat.tasks.map((task) => {
                                  const key = `${project.id}::${task}`
                                  const isDone = Boolean(completedTasks[key])
                                  return (
                                    <td
                                      key={`${project.id}-${cat.name}-${task}`}
                                      onClick={() => handleOpenTaskDrawer(project, task)}
                                      className="p-2 border-r border-border/60 text-center cursor-pointer hover:bg-muted/30 transition group/cell"
                                    >
                                      <div className="flex items-center justify-center">
                                        <button
                                          type="button"
                                          onClick={(e) => handleToggleTaskStatus(project.id, task, e)}
                                          className={cn(
                                            'h-5 w-5 rounded border flex items-center justify-center transition-all',
                                            isDone
                                              ? 'bg-emerald-600 border-emerald-600 text-white shadow-2xs'
                                              : 'border-border/80 group-hover/cell:border-primary/60 bg-card'
                                          )}
                                          title={isDone ? 'Mark as Pending' : 'Mark as Done'}
                                        >
                                          {isDone && <IconCheck className="h-3.5 w-3.5" />}
                                        </button>
                                      </div>
                                    </td>
                                  )
                                })
                              )
                            )}
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                )}
              </div>
            )
          })
        )}
      </div>

      {/* Task Drawer */}
      <ProductionTaskDrawer
        open={drawerOpen}
        onOpenChange={setDrawerOpen}
        task={selectedTaskDetail}
        onSaveTask={handleSaveDrawerTask}
      />

      {/* Add Task Modal */}
      <Dialog open={addTaskModalOpen} onOpenChange={setAddTaskModalOpen}>
        <DialogContent className="sm:max-w-[420px]">
          <DialogHeader>
            <DialogTitle className="text-sm font-bold">
              Add New {activeTab === 'pre' ? 'Pre-Production' : 'Post-Production'} Task
            </DialogTitle>
          </DialogHeader>

          <div className="space-y-4 py-2">
            <div className="space-y-1.5">
              <Label className="text-xs">Category</Label>
              <Select value={targetCategory} onValueChange={setTargetCategory}>
                <SelectTrigger className="h-9 text-xs">
                  <SelectValue placeholder="Select Category" />
                </SelectTrigger>
                <SelectContent>
                  {activeCategories.map((c) => (
                    <SelectItem key={c.name} value={c.name} className="text-xs">
                      {c.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1.5">
              <Label className="text-xs">Task Name</Label>
              <Input
                placeholder="e.g. Confirm audio backup mics"
                value={newTaskName}
                onChange={(e) => setNewTaskName(e.target.value)}
                className="h-9 text-xs"
              />
            </div>
          </div>

          <DialogFooter>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setAddTaskModalOpen(false)}
              className="text-xs"
            >
              Cancel
            </Button>
            <Button
              size="sm"
              onClick={handleAddTask}
              className="text-xs bg-emerald-600 hover:bg-emerald-700 text-white"
            >
              Add Task
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
