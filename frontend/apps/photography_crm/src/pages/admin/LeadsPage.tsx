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
import { Sheet, SheetContent } from '@repo/ui/components/ui/sheet'
import { toast } from '@repo/ui/components/ui/sonner'
import { cn } from '@repo/ui/lib/utils'
import { createColumnHelper } from '@tanstack/react-table'
import {
  IconPlus,
  IconTable,
  IconLayoutKanban,
  IconForms,
  IconEdit,
  IconTrash,
  IconPhone,
  IconMail,
  IconCurrencyRupee,
  IconFlame,
  IconCircleCheck,
  IconAlertTriangle,
  IconClock,
  IconFolderPlus,
} from "@repo/ui/icons/tabler-icons"
import type {
  Lead,
  LeadSource,
  PipelineStage,
  PrimaryPerspective,
  FollowupSubFilter,
} from '../../types/lead'
import { LeadKanbanView } from '../../components/LeadKanbanView'
import { useWorkflows, type WorkflowItem } from '../../hooks/useWorkflows'

interface LeadsPageProps {
  leads: Lead[]
  onAddNewLead: (stage?: PipelineStage) => void
  onEditLead: (lead: Lead) => void
  onDeleteLead: (leadId: string) => void
  onStageChange: (leadId: string, newStage: PipelineStage) => void
  onConvertToProject?: (lead: Lead) => void
}

const columnHelper = createColumnHelper<DataTableFeatures, Lead>()

export function LeadsPage({
  leads,
  onAddNewLead,
  onEditLead,
  onDeleteLead,
  onStageChange,
  onConvertToProject,
}: LeadsPageProps) {
  const { workflows } = useWorkflows('lead')

  // Stages grouped by pipeline
  const pipelineStagesMap = React.useMemo(() => {
    const map: Record<string, WorkflowItem[]> = {
      Contacted: [],
      Qualified: [],
      Converted: [],
      Followups: [],
    }
    ;(workflows || []).forEach((w) => {
      const p = w.pipeline
      if (!map[p]) map[p] = []
      map[p].push(w)
    })
    return map
  }, [workflows])

  // Color lookup helper for stages
  const stageColorMap = React.useMemo(() => {
    const map = new Map<string, string>()
    ;(workflows || []).forEach((w) => {
      map.set(w.stage.toLowerCase(), w.color_code)
    })
    return map
  }, [workflows])

  // View mode: Table vs Kanban
  const [viewMode, setViewMode] = React.useState<'table' | 'kanban'>('table')

  // Level 1: Primary Perspective Tabs
  const [activePerspective, setActivePerspective] = React.useState<PrimaryPerspective>('all')

  // Level 2 Sub-Filters
  const [newLeadsSourceFilter, setNewLeadsSourceFilter] = React.useState<string>('all')
  const [contactedSubFilter, setContactedSubFilter] = React.useState<string>('all')
  const [qualifiedSubFilter, setQualifiedSubFilter] = React.useState<string>('all')
  const [convertedSubFilter, setConvertedSubFilter] = React.useState<string>('all')
  const [followupFilter, setFollowupFilter] = React.useState<FollowupSubFilter>('all_followup')

  // Detail Drawer for quick view
  const [drawerLead, setDrawerLead] = React.useState<Lead | null>(null)

  // Follow-up helper
  const getFollowupCategory = (dateStr: string) => {
    if (!dateStr) return 'upcoming'
    const target = new Date(dateStr)
    const now = new Date()
    const isToday =
      target.getFullYear() === now.getFullYear() &&
      target.getMonth() === now.getMonth() &&
      target.getDate() === now.getDate()

    if (isToday) return 'today'
    if (target < now) return 'due'
    return 'upcoming'
  }

  // Filter leads based on selected primary perspective and secondary sub-filter
  const filteredLeads = React.useMemo(() => {
    return leads.filter((lead) => {
      // Perspective 1: All
      if (activePerspective === 'all') return true

      // Perspective 2: New Leads
      if (activePerspective === 'new_leads') {
        if (lead.pipelineStage !== 'new_leads') return false
        if (newLeadsSourceFilter === 'all') return true
        return lead.leadSource.toLowerCase() === newLeadsSourceFilter.toLowerCase()
      }

      // Perspective 3: Contacted
      if (activePerspective === 'contacted') {
        if (lead.pipelineStage !== 'contacted') return false
        if (contactedSubFilter === 'all') return true
        return lead.leadStage === contactedSubFilter
      }

      // Perspective 4: Qualified
      if (activePerspective === 'qualified') {
        if (lead.pipelineStage !== 'qualified') return false
        if (qualifiedSubFilter === 'all') return true
        return lead.leadStage === qualifiedSubFilter
      }

      // Perspective 5: Converted
      if (activePerspective === 'converted') {
        if (lead.pipelineStage !== 'converted') return false
        if (convertedSubFilter === 'all') return true
        return lead.leadStage === convertedSubFilter
      }

      // Perspective 6: Followups (Cross-cutting)
      if (activePerspective === 'followups') {
        const cat = getFollowupCategory(lead.followupDate)
        if (followupFilter === 'all_followup') return true
        if (followupFilter === 'my_followup') return lead.assignee?.initials === 'RM' // Current user context
        if (followupFilter === 'due_followup') return cat === 'due'
        if (followupFilter === 'today_followup') return cat === 'today'
        if (followupFilter === 'upcoming_followup') return cat === 'upcoming'
      }

      return true
    })
  }, [
    leads,
    activePerspective,
    newLeadsSourceFilter,
    contactedSubFilter,
    qualifiedSubFilter,
    convertedSubFilter,
    followupFilter,
  ])

  // Top metrics
  const metrics = React.useMemo(() => {
    const total = leads.length
    const totalPipelineValue = leads.reduce((acc, l) => acc + (l.quoteAmount || 0), 0)
    const wonCount = leads.filter((l) => l.leadStage === 'Booked').length
    const todayFollowups = leads.filter((l) => getFollowupCategory(l.followupDate) === 'today').length
    const dueFollowups = leads.filter((l) => getFollowupCategory(l.followupDate) === 'due').length
    return { total, totalPipelineValue, wonCount, todayFollowups, dueFollowups }
  }, [leads])

  // Count badges for Level 1 primary tabs
  const perspectiveCounts = React.useMemo(() => {
    return {
      all: leads.length,
      new_leads: leads.filter((l) => l.pipelineStage === 'new_leads').length,
      contacted: leads.filter((l) => l.pipelineStage === 'contacted').length,
      qualified: leads.filter((l) => l.pipelineStage === 'qualified').length,
      converted: leads.filter((l) => l.pipelineStage === 'converted').length,
      followups: leads.filter((l) => {
        const c = getFollowupCategory(l.followupDate)
        return c === 'today' || c === 'due'
      }).length,
    }
  }, [leads])

  // Table Columns (all 15 requested columns)
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
        columnHelper.accessor('name', {
          header: 'Leads Name',
          cell: (info) => {
            const row = info.row.original
            return (
              <div className="flex flex-col min-w-[140px]">
                <span
                  onClick={() => onEditLead(row)}
                  className="font-semibold text-foreground hover:text-primary hover:underline cursor-pointer"
                >
                  {row.name}
                </span>
                <span className="text-[11px] text-muted-foreground font-mono">{row.id}</span>
              </div>
            )
          },
        }),
        columnHelper.accessor('mobileNo', {
          header: 'Mobile No',
          cell: (info) => {
            const row = info.row.original
            return (
              <div className="flex items-center gap-1.5 font-mono text-xs tabular-nums text-foreground whitespace-nowrap">
                <IconPhone className="h-3 w-3 text-muted-foreground" />
                <span>
                  {row.countryCode} {row.mobileNo}
                </span>
              </div>
            )
          },
        }),
        columnHelper.accessor('email', {
          header: 'Email Address',
          cell: (info) => (
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground max-w-[170px] truncate">
              <IconMail className="h-3 w-3 shrink-0" />
              <span className="truncate">{info.getValue()}</span>
            </div>
          ),
        }),
        columnHelper.accessor('createdDate', {
          header: 'Created Date',
          cell: (info) => (
            <span className="text-xs text-muted-foreground tabular-nums whitespace-nowrap">
              {info.getValue()}
            </span>
          ),
        }),
        columnHelper.accessor('leadSource', {
          header: 'Leads Source',
          cell: (info) => {
            const val = info.getValue() as LeadSource
            return (
              <Badge variant="outline" className="text-[11px] font-medium whitespace-nowrap">
                {val}
              </Badge>
            )
          },
        }),
        columnHelper.accessor('generate', {
          header: 'Generate',
          cell: (info) => (
            <span className="text-xs text-muted-foreground truncate max-w-[130px] block">
              {info.getValue() || 'Direct'}
            </span>
          ),
        }),
        columnHelper.accessor('assignee', {
          header: 'Assignee',
          cell: (info) => {
            const val = info.getValue()
            return (
              <div className="flex items-center gap-1.5">
                <span className="flex h-5 w-5 items-center justify-center rounded-full bg-primary/10 text-[10px] font-bold text-primary">
                  {val?.initials || 'ME'}
                </span>
                <span className="text-xs font-medium text-foreground whitespace-nowrap">
                  {val?.name}
                </span>
              </div>
            )
          },
        }),
        columnHelper.accessor('leadPipeline', {
          header: 'Leads Pipeline',
          cell: (info) => (
            <span className="text-xs font-medium text-foreground whitespace-nowrap">
              {info.getValue()}
            </span>
          ),
        }),
        columnHelper.accessor('leadStage', {
          header: 'Leads Stage',
          cell: (info) => {
            const row = info.row.original
            const stage = row.leadStage
            const color = stageColorMap.get(stage.toLowerCase())
            return (
              <Badge
                variant="outline"
                className="text-[11px] font-medium whitespace-nowrap"
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
                {stage}
              </Badge>
            )
          },
        }),
        columnHelper.accessor('quoteAmount', {
          header: 'Quote Amount',
          cell: (info) => (
            <span className="text-xs font-bold text-foreground tabular-nums whitespace-nowrap">
              ₹{info.getValue().toLocaleString('en-IN')}
            </span>
          ),
        }),
        columnHelper.accessor('followupDate', {
          header: 'Followup Date',
          cell: (info) => {
            const val = info.getValue()
            const cat = getFollowupCategory(val)
            return (
              <div className="flex items-center gap-1 tabular-nums whitespace-nowrap">
                <span
                  className={cn(
                    'h-1.5 w-1.5 rounded-full',
                    cat === 'due' && 'bg-destructive',
                    cat === 'today' && 'bg-amber-500',
                    cat === 'upcoming' && 'bg-muted-foreground/50'
                  )}
                />
                <span className="text-xs text-muted-foreground">{val ? val.replace('T', ' ') : '—'}</span>
              </div>
            )
          },
        }),
        columnHelper.accessor('projectType', {
          header: 'Projects Type',
          cell: (info) => (
            <Badge variant="secondary" className="text-[11px] font-medium whitespace-nowrap">
              {info.getValue()}
            </Badge>
          ),
        }),
        columnHelper.accessor('quoteType', {
          header: 'Quote Type',
          cell: (info) => (
            <span className="text-xs text-muted-foreground whitespace-nowrap">{info.getValue()}</span>
          ),
        }),
        columnHelper.accessor('projectStartDate', {
          header: 'Projects Start Date',
          cell: (info) => (
            <span className="text-xs text-muted-foreground tabular-nums whitespace-nowrap">
              {info.getValue()}
            </span>
          ),
        }),
        columnHelper.accessor('projectEndDate', {
          header: 'Projects End Date',
          cell: (info) => (
            <span className="text-xs text-muted-foreground tabular-nums whitespace-nowrap">
              {info.getValue()}
            </span>
          ),
        }),
      ]),
    [onEditLead]
  )

  // Filter configs for tanstack table
  const filterConfigs: DataTableFilterConfig[] = [
    {
      columnId: 'pipelineStage',
      title: 'Pipeline Stage',
      type: 'select',
      options: [
        { label: 'New Leads', value: 'new_leads' },
        { label: 'Contacted', value: 'contacted' },
        { label: 'Qualified', value: 'qualified' },
        { label: 'Converted', value: 'converted' },
      ],
    },
    {
      columnId: 'projectType',
      title: 'Project Type',
      type: 'select',
      options: [
        { label: 'Wedding', value: 'Wedding' },
        { label: 'Pre-Wedding', value: 'Pre-Wedding' },
        { label: 'Engagement', value: 'Engagement' },
        { label: 'Corporate', value: 'Corporate' },
        { label: 'Maternity', value: 'Maternity' },
      ],
    },
  ]

  // Bulk actions
  const bulkActions: DataTableBulkAction<Lead>[] = [
    {
      label: 'Delete Selected',
      variant: 'destructive',
      action: (selectedRows: Lead[], clearSelection: () => void) => {
        selectedRows.forEach((r) => onDeleteLead(r.id))
        clearSelection()
        toast.success(`Deleted ${selectedRows.length} lead(s)`)
      },
    },
    {
      label: 'Mark as Won / Booked',
      action: (selectedRows: Lead[], clearSelection: () => void) => {
        selectedRows.forEach((r) => onStageChange(r.id, 'converted'))
        clearSelection()
        toast.success(`Moved ${selectedRows.length} lead(s) to Converted`)
      },
    },
  ]

  // Row actions
  const rowActions: DataTableRowAction<Lead>[] = [
    {
      label: 'Edit Lead',
      icon: IconEdit,
      action: (row: Lead) => onEditLead(row),
    },
    ...(onConvertToProject
      ? [
          {
            label: 'Convert to Project',
            icon: IconFolderPlus,
            action: (row: Lead) => onConvertToProject(row),
          },
        ]
      : []),
    {
      label: 'Delete Lead',
      icon: IconTrash,
      variant: 'destructive',
      action: (row: Lead) => {
        onDeleteLead(row.id)
        toast.success(`Lead ${row.name} deleted`)
      },
    },
  ]

  return (
    <div className="space-y-6">
      {/* Top Header: Title, KPIs, View Toggle & Add Lead Button */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">Photography Leads</h1>
          <p className="text-xs sm:text-sm text-muted-foreground mt-0.5">
            Full-lifecycle pipeline tracking: New leads, qualification, quotation, and confirmed shoots.
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
            onClick={() => onAddNewLead()}
            className="bg-emerald-600 hover:bg-emerald-700 text-white font-medium gap-1.5 shadow-xs"
          >
            <IconPlus className="h-4 w-4" />
            Add New Lead
          </Button>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3 sm:gap-4">
        <Card className="border shadow-xs bg-card">
          <CardContent className="p-4 flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-muted-foreground">Total Inquiries</p>
              <h3 className="text-xl sm:text-2xl font-bold mt-1 tabular-nums">{metrics.total}</h3>
            </div>
            <div className="h-10 w-10 rounded-xl bg-primary/10 text-primary flex items-center justify-center">
              <IconForms className="h-5 w-5" />
            </div>
          </CardContent>
        </Card>

        <Card className="border shadow-xs bg-card">
          <CardContent className="p-4 flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-muted-foreground">Pipeline Value</p>
              <h3 className="text-xl sm:text-2xl font-bold mt-1 tabular-nums">
                ₹{(metrics.totalPipelineValue / 100000).toFixed(2)}L
              </h3>
            </div>
            <div className="h-10 w-10 rounded-xl bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 flex items-center justify-center">
              <IconCurrencyRupee className="h-5 w-5" />
            </div>
          </CardContent>
        </Card>

        <Card className="border shadow-xs bg-card">
          <CardContent className="p-4 flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-muted-foreground">Confirmed Shoots</p>
              <h3 className="text-xl sm:text-2xl font-bold mt-1 tabular-nums">{metrics.wonCount}</h3>
            </div>
            <div className="h-10 w-10 rounded-xl bg-blue-500/10 text-blue-600 dark:text-blue-400 flex items-center justify-center">
              <IconCircleCheck className="h-5 w-5" />
            </div>
          </CardContent>
        </Card>

        <Card className="border shadow-xs bg-card">
          <CardContent className="p-4 flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-muted-foreground">Action Needed</p>
              <div className="flex items-center gap-2 mt-1">
                <span className="text-xl sm:text-2xl font-bold tabular-nums text-amber-600 dark:text-amber-400">
                  {metrics.todayFollowups}
                </span>
                <span className="text-xs text-muted-foreground">today</span>
                {metrics.dueFollowups > 0 && (
                  <Badge variant="destructive" className="text-[10px] py-0 px-1.5 tabular-nums">
                    {metrics.dueFollowups} due
                  </Badge>
                )}
              </div>
            </div>
            <div className="h-10 w-10 rounded-xl bg-amber-500/10 text-amber-600 dark:text-amber-400 flex items-center justify-center">
              <IconFlame className="h-5 w-5" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* LEVEL 1: Primary Pipeline Perspective Tabs */}
      <div className="border-b">
        <div className="flex items-center gap-1 overflow-x-auto scrollbar-none -mb-px">
          {[
            { id: 'all', label: 'All Leads', count: perspectiveCounts.all },
            { id: 'new_leads', label: 'New Leads', count: perspectiveCounts.new_leads },
            { id: 'contacted', label: 'Contacted', count: perspectiveCounts.contacted },
            { id: 'qualified', label: 'Qualified', count: perspectiveCounts.qualified },
            { id: 'converted', label: 'Converted', count: perspectiveCounts.converted },
            { id: 'followups', label: 'Followups', count: perspectiveCounts.followups, isHot: true },
          ].map((tab) => {
            const isActive = activePerspective === tab.id
            return (
              <button
                key={tab.id}
                type="button"
                onClick={() => setActivePerspective(tab.id as PrimaryPerspective)}
                className={cn(
                  'flex items-center gap-2 px-4 py-2.5 text-xs font-semibold rounded-t-lg transition border-b-2 whitespace-nowrap',
                  isActive
                    ? 'border-primary text-primary bg-accent/40'
                    : 'border-transparent text-muted-foreground hover:text-foreground hover:bg-muted/40',
                  tab.isHot && !isActive && 'text-amber-600 dark:text-amber-400'
                )}
              >
                {tab.isHot && <IconFlame className="h-3.5 w-3.5 text-amber-500" />}
                <span>{tab.label}</span>
                <span
                  className={cn(
                    'rounded-full px-2 py-0.2 text-[10px] tabular-nums font-bold',
                    isActive
                      ? 'bg-primary text-primary-foreground'
                      : 'bg-muted text-muted-foreground',
                    tab.isHot && 'bg-amber-500/20 text-amber-600 dark:text-amber-300'
                  )}
                >
                  {tab.count}
                </span>
              </button>
            )
          })}
        </div>
      </div>

      {/* LEVEL 2: Dynamic Sub-Filters depending on active Level 1 tab */}
      <div className="flex items-center gap-2 overflow-x-auto py-1 scrollbar-none">
        <span className="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider shrink-0 mr-1">
          Filter:
        </span>

        {/* Level 2 for 'new_leads' -> Channels */}
        {activePerspective === 'new_leads' && (
          <>
            {[
              { id: 'all', label: 'All Leads' },
              { id: 'messenger', label: 'Messenger' },
              { id: 'instagram', label: 'Instagram' },
              { id: 'whatsapp', label: 'Whatsapp' },
              { id: 'leads ads', label: 'Leads Ads' },
              { id: 'form', label: 'Form' },
              { id: 'manual leads', label: 'Manual Leads' },
            ].map((sub) => (
              <button
                key={sub.id}
                type="button"
                onClick={() => setNewLeadsSourceFilter(sub.id)}
                className={cn(
                  'rounded-full px-3 py-1 text-xs font-medium transition whitespace-nowrap border',
                  newLeadsSourceFilter === sub.id
                    ? 'bg-foreground text-background border-foreground font-semibold'
                    : 'bg-card text-muted-foreground border-border hover:bg-muted hover:text-foreground'
                )}
              >
                {sub.label}
              </button>
            ))}
          </>
        )}

        {/* Level 2 for 'contacted' */}
        {activePerspective === 'contacted' && (
          <>
            {[
              { id: 'all', label: 'All Contacted' },
              ...(pipelineStagesMap.Contacted || []).map((w) => ({
                id: w.stage,
                label: w.stage,
                color: w.color_code,
              })),
            ].map((sub) => (
              <button
                key={sub.id}
                type="button"
                onClick={() => setContactedSubFilter(sub.id)}
                className={cn(
                  'rounded-full px-3 py-1 text-xs font-medium transition whitespace-nowrap border',
                  contactedSubFilter === sub.id
                    ? 'bg-foreground text-background border-foreground font-semibold'
                    : 'bg-card text-muted-foreground border-border hover:bg-muted hover:text-foreground'
                )}
              >
                {sub.label}
              </button>
            ))}
          </>
        )}

        {/* Level 2 for 'qualified' */}
        {activePerspective === 'qualified' && (
          <>
            {[
              { id: 'all', label: 'All Qualified' },
              ...(pipelineStagesMap.Qualified || []).map((w) => ({
                id: w.stage,
                label: w.stage,
                color: w.color_code,
              })),
            ].map((sub) => (
              <button
                key={sub.id}
                type="button"
                onClick={() => setQualifiedSubFilter(sub.id)}
                className={cn(
                  'rounded-full px-3 py-1 text-xs font-medium transition whitespace-nowrap border',
                  qualifiedSubFilter === sub.id
                    ? 'bg-foreground text-background border-foreground font-semibold'
                    : 'bg-card text-muted-foreground border-border hover:bg-muted hover:text-foreground'
                )}
              >
                {sub.label}
              </button>
            ))}
          </>
        )}

        {/* Level 2 for 'converted' */}
        {activePerspective === 'converted' && (
          <>
            {[
              { id: 'all', label: 'All Converted' },
              ...(pipelineStagesMap.Converted || []).map((w) => ({
                id: w.stage,
                label: w.stage,
                color: w.color_code,
              })),
            ].map((sub) => (
              <button
                key={sub.id}
                type="button"
                onClick={() => setConvertedSubFilter(sub.id)}
                className={cn(
                  'rounded-full px-3 py-1 text-xs font-medium transition whitespace-nowrap border',
                  convertedSubFilter === sub.id
                    ? 'bg-foreground text-background border-foreground font-semibold'
                    : 'bg-card text-muted-foreground border-border hover:bg-muted hover:text-foreground'
                )}
              >
                {sub.label}
              </button>
            ))}
          </>
        )}

        {/* Level 2 for 'followups' */}
        {activePerspective === 'followups' && (
          <>
            {[
              { id: 'all_followup', label: 'All Followup', icon: undefined },
              ...(pipelineStagesMap.Followups && pipelineStagesMap.Followups.length > 0
                ? pipelineStagesMap.Followups.filter((w) => w.stage !== 'All Followup').map((w) => ({
                    id: w.stage.toLowerCase().replace(/\s+/g, '_'),
                    label: w.stage,
                    icon:
                      w.stage.toLowerCase().includes('due')
                        ? IconAlertTriangle
                        : w.stage.toLowerCase().includes('today')
                        ? IconClock
                        : undefined,
                  }))
                : [
                    { id: 'my_followup', label: 'My Followup', icon: undefined },
                    { id: 'due_followup', label: 'Due Followup', icon: IconAlertTriangle },
                    { id: 'today_followup', label: 'Today Followup', icon: IconClock },
                    { id: 'upcoming_followup', label: 'Upcoming Followup', icon: undefined },
                  ]),
            ].map((sub) => {
              const SubIcon = sub.icon
              return (
                <button
                  key={sub.id}
                  type="button"
                  onClick={() => setFollowupFilter(sub.id as FollowupSubFilter)}
                  className={cn(
                    'inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-medium transition whitespace-nowrap border',
                    followupFilter === sub.id
                      ? 'bg-foreground text-background border-foreground font-semibold'
                      : 'bg-card text-muted-foreground border-border hover:bg-muted hover:text-foreground'
                  )}
                >
                  {SubIcon && <SubIcon className="h-3 w-3" />}
                  <span>{sub.label}</span>
                </button>
              )
            })}
          </>
        )}

        {/* Level 2 for 'all' */}
        {activePerspective === 'all' && (
          <span className="text-xs text-muted-foreground">Showing full pipeline records</span>
        )}
      </div>

      {/* View Switch: Table vs Kanban */}
      {viewMode === 'table' ? (
        <Card className="border bg-card shadow-xs">
          <CardContent className="p-3 sm:p-5">
            <DataTable
              columns={columns}
              data={filteredLeads}
              searchableColumns={['name', 'email', 'mobileNo', 'leadSource', 'projectType', 'id']}
              searchPlaceholder="Search leads by name, email, phone... (Press ⌘K)"
              filterConfigs={filterConfigs}
              bulkActions={bulkActions}
              rowActions={rowActions}
              renderDetailDrawer={(row, close) => (
                <div className="space-y-6">
                  <div className="rounded-lg border bg-muted/30 p-4 space-y-3">
                    <div className="flex items-center justify-between">
                      <span className="text-xs font-medium text-muted-foreground">Pipeline Stage</span>
                      <Badge variant="outline" className="font-semibold text-xs">
                        {row.leadStage}
                      </Badge>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-xs font-medium text-muted-foreground">Quote Total</span>
                      <span className="text-lg font-bold tabular-nums">
                        ₹{row.quoteAmount.toLocaleString('en-IN')}
                      </span>
                    </div>
                  </div>

                  <div className="space-y-3">
                    <h4 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                      Lead & Event Details
                    </h4>
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between py-1 border-b">
                        <span className="text-muted-foreground">Lead Name</span>
                        <span className="font-medium">{row.name}</span>
                      </div>
                      <div className="flex justify-between py-1 border-b">
                        <span className="text-muted-foreground">Mobile</span>
                        <span className="font-mono text-xs">
                          {row.countryCode} {row.mobileNo}
                        </span>
                      </div>
                      <div className="flex justify-between py-1 border-b">
                        <span className="text-muted-foreground">Email</span>
                        <span className="font-mono text-xs">{row.email}</span>
                      </div>
                      <div className="flex justify-between py-1 border-b">
                        <span className="text-muted-foreground">Lead Source</span>
                        <Badge variant="outline" className="text-xs">
                          {row.leadSource}
                        </Badge>
                      </div>
                      <div className="flex justify-between py-1 border-b">
                        <span className="text-muted-foreground">Project Type</span>
                        <span>{row.projectType}</span>
                      </div>
                      <div className="flex justify-between py-1 border-b">
                        <span className="text-muted-foreground">Shoot Dates</span>
                        <span className="tabular-nums">
                          {row.projectStartDate} → {row.projectEndDate}
                        </span>
                      </div>
                      <div className="flex justify-between py-1 border-b">
                        <span className="text-muted-foreground">Followup Scheduled</span>
                        <span className="tabular-nums font-mono text-xs">
                          {row.followupDate ? row.followupDate.replace('T', ' ') : '—'}
                        </span>
                      </div>
                      <div className="flex justify-between py-1 border-b">
                        <span className="text-muted-foreground">Deliverables</span>
                        <span className="text-xs text-right max-w-[200px]">
                          {row.selectedDeliverables?.join(', ') || 'Standard package'}
                        </span>
                      </div>
                    </div>
                  </div>

                  <div className="flex gap-2 pt-4">
                    <Button
                      className="flex-1 bg-emerald-600 hover:bg-emerald-700 text-white"
                      onClick={() => {
                        close()
                        onEditLead(row)
                      }}
                    >
                      Edit Complete Lead
                    </Button>
                    <Button variant="outline" onClick={close}>
                      Close
                    </Button>
                  </div>
                </div>
              )}
              detailDrawerTitle={(row) => `Lead Details - ${row.name}`}
              detailDrawerDescription={(row) =>
                `Created on ${row.createdDate} via ${row.leadSource} • Pipeline: ${row.leadPipeline}`
              }
              enableExport={true}
              exportFilename="photography-leads.csv"
            />
          </CardContent>
        </Card>
      ) : (
        /* Kanban View */
        <LeadKanbanView
          leads={filteredLeads}
          onLeadClick={(lead) => setDrawerLead(lead)}
          onEditLead={onEditLead}
          onDeleteLead={onDeleteLead}
          onStageChange={onStageChange}
          onAddNewLead={(stage) => onAddNewLead(stage)}
        />
      )}

      {/* Quick view sheet for Kanban card clicks */}
      {drawerLead && (
        <Sheet open={Boolean(drawerLead)} onOpenChange={(open) => !open && setDrawerLead(null)}>
          <SheetContent side="right" className="w-full sm:max-w-md p-6">
            <div className="space-y-6">
              <div>
                <h3 className="text-lg font-bold text-foreground">{drawerLead.name}</h3>
                <p className="text-xs text-muted-foreground mt-0.5">
                  ID: {drawerLead.id} • {drawerLead.leadSource}
                </p>
              </div>

              <div className="rounded-xl border bg-muted/30 p-4 space-y-2">
                <div className="flex justify-between items-center text-xs">
                  <span className="text-muted-foreground">Stage:</span>
                  <Badge variant="outline" className="font-semibold">
                    {drawerLead.leadStage}
                  </Badge>
                </div>
                <div className="flex justify-between items-center text-xs">
                  <span className="text-muted-foreground">Deal Amount:</span>
                  <span className="text-base font-bold text-foreground tabular-nums">
                    ₹{drawerLead.quoteAmount.toLocaleString('en-IN')}
                  </span>
                </div>
              </div>

              <div className="space-y-2 text-xs">
                <div className="flex justify-between py-1.5 border-b">
                  <span className="text-muted-foreground">Phone:</span>
                  <span className="font-mono">
                    {drawerLead.countryCode} {drawerLead.mobileNo}
                  </span>
                </div>
                <div className="flex justify-between py-1.5 border-b">
                  <span className="text-muted-foreground">Email:</span>
                  <span className="font-mono">{drawerLead.email}</span>
                </div>
                <div className="flex justify-between py-1.5 border-b">
                  <span className="text-muted-foreground">Shoot Type:</span>
                  <span>{drawerLead.projectType}</span>
                </div>
                <div className="flex justify-between py-1.5 border-b">
                  <span className="text-muted-foreground">Followup:</span>
                  <span className="font-mono">{drawerLead.followupDate?.replace('T', ' ')}</span>
                </div>
              </div>

              <div className="flex gap-2 pt-4">
                <Button
                  className="flex-1 bg-emerald-600 hover:bg-emerald-700 text-white"
                  onClick={() => {
                    const l = drawerLead
                    setDrawerLead(null)
                    onEditLead(l)
                  }}
                >
                  Edit Lead
                </Button>
                {onConvertToProject && (
                  <Button
                    variant="secondary"
                    className="gap-1.5"
                    onClick={() => {
                      const l = drawerLead
                      setDrawerLead(null)
                      onConvertToProject(l)
                    }}
                  >
                    <IconFolderPlus className="h-4 w-4" />
                    Convert to Project
                  </Button>
                )}
                <Button variant="outline" onClick={() => setDrawerLead(null)}>
                  Close
                </Button>
              </div>
            </div>
          </SheetContent>
        </Sheet>
      )}
    </div>
  )
}
