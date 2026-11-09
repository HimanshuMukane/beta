import * as React from 'react'
import {
  IconClock,
  IconPlus,
  IconDotsVertical,
  IconEdit,
  IconTrash,
  IconEye,
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
import type { Lead, PipelineStage } from '../types/lead'

interface LeadKanbanViewProps {
  leads: Lead[]
  onLeadClick: (lead: Lead) => void
  onEditLead: (lead: Lead) => void
  onDeleteLead: (leadId: string) => void
  onStageChange: (leadId: string, newStage: PipelineStage) => void
  onAddNewLead: (stage?: PipelineStage) => void
}

interface KanbanColumn {
  id: PipelineStage
  title: string
  colorBadge: string
  accentBorder: string
}

const KANBAN_COLUMNS: KanbanColumn[] = [
  {
    id: 'new_leads',
    title: 'New Leads',
    colorBadge: 'bg-blue-500/15 text-blue-700 dark:text-blue-300 border-blue-500/30',
    accentBorder: 'border-t-blue-500',
  },
  {
    id: 'contacted',
    title: 'Contacted',
    colorBadge: 'bg-amber-500/15 text-amber-700 dark:text-amber-300 border-amber-500/30',
    accentBorder: 'border-t-amber-500',
  },
  {
    id: 'qualified',
    title: 'Qualified',
    colorBadge: 'bg-indigo-500/15 text-indigo-700 dark:text-indigo-300 border-indigo-500/30',
    accentBorder: 'border-t-indigo-500',
  },
  {
    id: 'converted',
    title: 'Converted',
    colorBadge: 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300 border-emerald-500/30',
    accentBorder: 'border-t-emerald-500',
  },
]

export function LeadKanbanView({
  leads,
  onLeadClick,
  onEditLead,
  onDeleteLead,
  onStageChange,
  onAddNewLead,
}: LeadKanbanViewProps) {
  const [draggedLeadId, setDraggedLeadId] = React.useState<string | null>(null)
  const [dragOverColumn, setDragOverColumn] = React.useState<PipelineStage | null>(null)

  // Drag and drop events
  const handleDragStart = (e: React.DragEvent, leadId: string) => {
    e.dataTransfer.setData('text/plain', leadId)
    e.dataTransfer.effectAllowed = 'move'
    setDraggedLeadId(leadId)
  }

  const handleDragOver = (e: React.DragEvent, colId: PipelineStage) => {
    e.preventDefault()
    e.dataTransfer.dropEffect = 'move'
    if (dragOverColumn !== colId) {
      setDragOverColumn(colId)
    }
  }

  const handleDragLeave = () => {
    setDragOverColumn(null)
  }

  const handleDrop = (e: React.DragEvent, targetStage: PipelineStage) => {
    e.preventDefault()
    const leadId = e.dataTransfer.getData('text/plain') || draggedLeadId
    if (leadId) {
      onStageChange(leadId, targetStage)
    }
    setDraggedLeadId(null)
    setDragOverColumn(null)
  }

  const handleDragEnd = () => {
    setDraggedLeadId(null)
    setDragOverColumn(null)
  }

  // Follow-up status helper
  const getFollowupUrgency = (dateStr: string) => {
    if (!dateStr) return { label: 'No Followup', variant: 'muted' as const }
    const target = new Date(dateStr)
    const now = new Date()
    const isToday =
      target.getFullYear() === now.getFullYear() &&
      target.getMonth() === now.getMonth() &&
      target.getDate() === now.getDate()

    if (isToday) {
      return { label: 'Today', variant: 'today' as const }
    }
    if (target < now) {
      return { label: 'Due / Overdue', variant: 'due' as const }
    }
    return { label: 'Upcoming', variant: 'upcoming' as const }
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 items-start">
      {KANBAN_COLUMNS.map((column) => {
        const columnLeads = leads.filter((l) => l.pipelineStage === column.id)
        const totalValue = columnLeads.reduce((acc, l) => acc + (l.quoteAmount || 0), 0)
        const isTarget = dragOverColumn === column.id

        return (
          <div
            key={column.id}
            onDragOver={(e) => handleDragOver(e, column.id)}
            onDragLeave={handleDragLeave}
            onDrop={(e) => handleDrop(e, column.id)}
            className={cn(
              'flex flex-col rounded-xl border bg-muted/20 border-t-4 transition-all duration-150 min-h-[500px]',
              column.accentBorder,
              isTarget && 'bg-primary/5 ring-2 ring-primary/40'
            )}
          >
            {/* Column Header */}
            <div className="flex items-center justify-between p-3.5 border-b bg-card/60 rounded-t-lg">
              <div className="flex items-center gap-2">
                <span className="font-semibold text-sm text-foreground">{column.title}</span>
                <span className="flex h-5 w-5 items-center justify-center rounded-full bg-muted text-[11px] font-bold text-muted-foreground tabular-nums">
                  {columnLeads.length}
                </span>
              </div>
              <div className="flex items-center gap-1.5">
                <span className="text-xs font-bold text-foreground tabular-nums">
                  ₹{(totalValue / 1000).toFixed(0)}k
                </span>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => onAddNewLead(column.id)}
                  className="h-6 w-6 p-0 text-muted-foreground hover:text-foreground"
                  title={`Add lead in ${column.title}`}
                >
                  <IconPlus className="h-3.5 w-3.5" />
                </Button>
              </div>
            </div>

            {/* Column Body / Cards List */}
            <div className="flex-1 p-2.5 space-y-2.5 overflow-y-auto max-h-[calc(100vh-280px)]">
              {columnLeads.length === 0 ? (
                <div className="flex flex-col items-center justify-center py-12 text-center border border-dashed rounded-lg border-border/80">
                  <p className="text-xs text-muted-foreground">No leads in this stage</p>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => onAddNewLead(column.id)}
                    className="mt-2 text-xs text-primary h-7 px-2"
                  >
                    + Add Lead
                  </Button>
                </div>
              ) : (
                columnLeads.map((lead) => {
                  const urgency = getFollowupUrgency(lead.followupDate)
                  const isDragging = draggedLeadId === lead.id

                  return (
                    <div
                      key={lead.id}
                      draggable
                      onDragStart={(e) => handleDragStart(e, lead.id)}
                      onDragEnd={handleDragEnd}
                      onClick={() => onLeadClick(lead)}
                      className={cn(
                        'group rounded-xl border bg-card p-3.5 shadow-xs transition-all duration-150 hover:shadow-md hover:border-primary/40 cursor-grab active:cursor-grabbing space-y-3 relative',
                        isDragging && 'opacity-40 scale-95 border-dashed border-primary'
                      )}
                    >
                      {/* Top Row: Lead ID + Channel Badge + Context Menu */}
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-1.5">
                          <Badge variant="outline" className="text-[10px] font-mono px-1.5 py-0">
                            {lead.id}
                          </Badge>
                          <Badge
                            variant="secondary"
                            className="text-[10px] font-medium px-1.5 py-0 truncate max-w-[90px]"
                          >
                            {lead.leadSource}
                          </Badge>
                        </div>

                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button
                              type="button"
                              variant="ghost"
                              size="sm"
                              onClick={(e) => e.stopPropagation()}
                              className="h-6 w-6 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
                            >
                              <IconDotsVertical className="h-3.5 w-3.5 text-muted-foreground" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end" className="text-xs">
                            <DropdownMenuItem
                              onClick={(e) => {
                                e.stopPropagation()
                                onEditLead(lead)
                              }}
                            >
                              <IconEdit className="mr-2 h-3.5 w-3.5" /> Edit Lead
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              onClick={(e) => {
                                e.stopPropagation()
                                onLeadClick(lead)
                              }}
                            >
                              <IconEye className="mr-2 h-3.5 w-3.5" /> Quick View
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              onClick={(e) => {
                                e.stopPropagation()
                                onDeleteLead(lead.id)
                              }}
                              className="text-destructive focus:text-destructive"
                            >
                              <IconTrash className="mr-2 h-3.5 w-3.5" /> Delete
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </div>

                      {/* Lead Name & Project Type */}
                      <div>
                        <h4 className="font-semibold text-sm text-foreground line-clamp-1 group-hover:text-primary transition-colors">
                          {lead.name}
                        </h4>
                        <p className="text-xs text-muted-foreground line-clamp-1">
                          {lead.projectName || `${lead.projectType} Project`}
                        </p>
                      </div>

                      {/* Sub-Stage Pill */}
                      <div className="flex items-center justify-between pt-0.5">
                        <span className="text-[11px] font-medium text-foreground bg-muted/60 px-2 py-0.5 rounded-md">
                          {lead.leadStage}
                        </span>
                        <span className="text-xs font-bold text-foreground tabular-nums">
                          ₹{lead.quoteAmount.toLocaleString('en-IN')}
                        </span>
                      </div>

                      {/* Footer Row: Followup Pill + Assignee Initials */}
                      <div className="flex items-center justify-between border-t pt-2.5 text-[11px]">
                        <div
                          className={cn(
                            'flex items-center gap-1 font-medium px-1.5 py-0.5 rounded text-[10px]',
                            urgency.variant === 'due' &&
                              'bg-destructive/10 text-destructive border border-destructive/20',
                            urgency.variant === 'today' &&
                              'bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/20',
                            urgency.variant === 'upcoming' &&
                              'bg-muted text-muted-foreground'
                          )}
                          title={`Followup scheduled: ${lead.followupDate}`}
                        >
                          <IconClock className="h-3 w-3 shrink-0" />
                          <span>{urgency.label}</span>
                        </div>

                        <div className="flex items-center gap-1">
                          <span
                            className="flex h-5 w-5 items-center justify-center rounded-full bg-primary/10 text-[10px] font-bold text-primary"
                            title={`Assigned to: ${lead.assignee.name}`}
                          >
                            {lead.assignee.initials}
                          </span>
                        </div>
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
