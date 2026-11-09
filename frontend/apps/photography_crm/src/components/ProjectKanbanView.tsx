import * as React from 'react'
import {
  IconCalendar,
  IconDotsVertical,
  IconEdit,
  IconTrash,
  IconPlus,
  IconEye,
  IconChecklist,
} from '@tabler/icons-react'
import {
  Badge,
  Button,
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  cn,
} from '@repo/ui'
import type { Project, ProjectStage } from '../types/project'

interface ProjectKanbanViewProps {
  projects: Project[]
  onProjectClick: (project: Project) => void
  onEditProject: (project: Project) => void
  onDeleteProject: (projectId: string) => void
  onStageChange: (projectId: string, newStage: ProjectStage) => void
  onAddNewProject: () => void
}

interface KanbanStageDef {
  id: ProjectStage
  title: string
  accentBorder: string
  badgeClass: string
}

const PROJECT_STAGES: KanbanStageDef[] = [
  { id: 'Booked', title: 'Booked', accentBorder: 'border-t-blue-500', badgeClass: 'bg-blue-500/15 text-blue-700 dark:text-blue-300 border-blue-500/30' },
  { id: 'Shoot Planning', title: 'Shoot Planning', accentBorder: 'border-t-amber-500', badgeClass: 'bg-amber-500/15 text-amber-700 dark:text-amber-300 border-amber-500/30' },
  { id: 'Shoot Completed', title: 'Shoot Completed', accentBorder: 'border-t-indigo-500', badgeClass: 'bg-indigo-500/15 text-indigo-700 dark:text-indigo-300 border-indigo-500/30' },
  { id: 'Editing', title: 'Editing', accentBorder: 'border-t-purple-500', badgeClass: 'bg-purple-500/15 text-purple-700 dark:text-purple-300 border-purple-500/30' },
  { id: 'Delivery', title: 'Delivery', accentBorder: 'border-t-teal-500', badgeClass: 'bg-teal-500/15 text-teal-700 dark:text-teal-300 border-teal-500/30' },
  { id: 'Closure', title: 'Closure', accentBorder: 'border-t-emerald-500', badgeClass: 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300 border-emerald-500/30' },
]

export function ProjectKanbanView({
  projects,
  onProjectClick,
  onEditProject,
  onDeleteProject,
  onStageChange,
  onAddNewProject,
}: ProjectKanbanViewProps) {
  const [draggedProjectId, setDraggedProjectId] = React.useState<string | null>(null)
  const [dragOverColumn, setDragOverColumn] = React.useState<ProjectStage | null>(null)

  const handleDragStart = (e: React.DragEvent, id: string) => {
    e.dataTransfer.setData('text/plain', id)
    e.dataTransfer.effectAllowed = 'move'
    setDraggedProjectId(id)
  }

  const handleDragOver = (e: React.DragEvent, colId: ProjectStage) => {
    e.preventDefault()
    e.dataTransfer.dropEffect = 'move'
    if (dragOverColumn !== colId) {
      setDragOverColumn(colId)
    }
  }

  const handleDragLeave = () => {
    setDragOverColumn(null)
  }

  const handleDrop = (e: React.DragEvent, targetStage: ProjectStage) => {
    e.preventDefault()
    const projectId = e.dataTransfer.getData('text/plain') || draggedProjectId
    if (projectId) {
      onStageChange(projectId, targetStage)
    }
    setDraggedProjectId(null)
    setDragOverColumn(null)
  }

  const handleDragEnd = () => {
    setDraggedProjectId(null)
    setDragOverColumn(null)
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-6 gap-3.5 items-start overflow-x-auto pb-4">
      {PROJECT_STAGES.map((col) => {
        const stageProjects = projects.filter((p) => p.stage === col.id)
        const isTarget = dragOverColumn === col.id

        return (
          <div
            key={col.id}
            onDragOver={(e) => handleDragOver(e, col.id)}
            onDragLeave={handleDragLeave}
            onDrop={(e) => handleDrop(e, col.id)}
            className={cn(
              'flex flex-col rounded-xl border bg-muted/20 border-t-4 transition-all duration-150 min-h-[480px]',
              col.accentBorder,
              isTarget && 'bg-primary/5 ring-2 ring-primary/40'
            )}
          >
            {/* Column Header */}
            <div className="flex items-center justify-between p-3 border-b bg-card/60 rounded-t-lg">
              <div className="flex items-center gap-1.5 min-w-0">
                <span className="font-semibold text-xs text-foreground truncate">{col.title}</span>
                <span className="flex h-4 w-4 shrink-0 items-center justify-center rounded-full bg-muted text-[10px] font-bold text-muted-foreground tabular-nums">
                  {stageProjects.length}
                </span>
              </div>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={onAddNewProject}
                className="h-5 w-5 p-0 text-muted-foreground hover:text-foreground"
                title="Add Project"
              >
                <IconPlus className="h-3 w-3" />
              </Button>
            </div>

            {/* Cards List */}
            <div className="flex-1 p-2 space-y-2.5 overflow-y-auto max-h-[calc(100vh-280px)]">
              {stageProjects.length === 0 ? (
                <div className="flex flex-col items-center justify-center py-10 text-center border border-dashed rounded-lg border-border/80">
                  <p className="text-[11px] text-muted-foreground">No projects</p>
                </div>
              ) : (
                stageProjects.map((prj) => {
                  const isDragging = draggedProjectId === prj.id

                  return (
                    <div
                      key={prj.id}
                      draggable
                      onDragStart={(e) => handleDragStart(e, prj.id)}
                      onDragEnd={handleDragEnd}
                      onClick={() => onProjectClick(prj)}
                      className={cn(
                        'group rounded-xl border bg-card p-3 shadow-xs transition-all duration-150 hover:shadow-md hover:border-primary/40 cursor-grab active:cursor-grabbing space-y-2.5 relative',
                        isDragging && 'opacity-40 scale-95 border-dashed border-primary'
                      )}
                    >
                      {/* Top: ID + Status + Menu */}
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-1.5">
                          <span className="text-[10px] font-mono text-muted-foreground font-semibold">
                            {prj.id}
                          </span>
                          <Badge variant="outline" className="text-[9px] px-1 py-0 font-medium">
                            {prj.status}
                          </Badge>
                        </div>

                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button
                              type="button"
                              variant="ghost"
                              size="sm"
                              onClick={(e) => e.stopPropagation()}
                              className="h-5 w-5 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
                            >
                              <IconDotsVertical className="h-3 w-3 text-muted-foreground" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end" className="text-xs">
                            <DropdownMenuItem
                              onClick={(e) => {
                                e.stopPropagation()
                                onEditProject(prj)
                              }}
                            >
                              <IconEdit className="mr-2 h-3.5 w-3.5" /> Edit Project
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              onClick={(e) => {
                                e.stopPropagation()
                                onProjectClick(prj)
                              }}
                            >
                              <IconEye className="mr-2 h-3.5 w-3.5" /> Open Workspace
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              onClick={(e) => {
                                e.stopPropagation()
                                onDeleteProject(prj.id)
                              }}
                              className="text-destructive focus:text-destructive"
                            >
                              <IconTrash className="mr-2 h-3.5 w-3.5" /> Delete
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </div>

                      {/* Project Title & Client */}
                      <div>
                        <h4 className="font-semibold text-xs text-foreground line-clamp-1 group-hover:text-primary transition-colors">
                          {prj.projectName}
                        </h4>
                        <p className="text-[11px] text-muted-foreground line-clamp-1">
                          {prj.clientName}
                        </p>
                      </div>

                      {/* Sessions & Deliverables progress count */}
                      <div className="flex items-center gap-2 text-[10px] text-muted-foreground pt-1 border-t">
                        <span className="flex items-center gap-1">
                          <IconCalendar className="h-3 w-3" />
                          <span>{prj.sessionsCount} sessions</span>
                        </span>
                        <span>•</span>
                        <span className="flex items-center gap-1">
                          <IconChecklist className="h-3 w-3" />
                          <span>
                            {prj.tasksCount.completed}/{prj.tasksCount.total}
                          </span>
                        </span>
                      </div>

                      {/* Footer: Manager initials + Type */}
                      <div className="flex items-center justify-between text-[10px] text-muted-foreground">
                        <Badge variant="secondary" className="text-[9px] px-1 py-0">
                          {prj.projectType}
                        </Badge>
                        <span
                          className="flex h-4 w-4 items-center justify-center rounded-full bg-primary/10 text-[9px] font-bold text-primary"
                          title={`PM: ${prj.projectManager.name}`}
                        >
                          {prj.projectManager.initials}
                        </span>
                      </div>
                    </div>
                  )
                })
              )}
            </div>
          </div>
        )
      })}
    </div>
  )
}
