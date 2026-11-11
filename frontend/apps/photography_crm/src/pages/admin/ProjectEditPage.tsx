import * as React from 'react'
import {
  IconArrowLeft,
  IconEdit,
  IconTrash,
  IconDotsVertical,
  IconPlus,
  IconFileText,
  IconCalendar,
  IconChecklist,
  IconPrinter,
  IconPhoto,
  IconReceipt2,
  IconNote,
  IconChartBar,
  IconClock,
  IconCheck,
  IconInbox,
} from "@repo/ui/icons/tabler-icons"
import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import { Badge } from '@repo/ui/components/ui/badge'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@repo/ui/components/ui/card'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@repo/ui/components/ui/select'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@repo/ui/components/ui/dropdown-menu'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@repo/ui/components/ui/popover'
import { toast } from '@repo/ui/components/ui/sonner'
import { cn } from '@repo/ui/lib/utils'
import { useWorkflows } from '../../hooks/useWorkflows'

const RichTextEditor = React.lazy(() =>
  import('@repo/ui/components/rich-text-editor/rich-text-editor').then((m) => ({
    default: m.RichTextEditor,
  }))
)
import { CreateSessionDrawer } from '../../components/CreateSessionDrawer'
import { ProductionMatrixView } from '../../components/ProductionMatrixView'
import { ProjectExpenseModal } from '../../components/ProjectExpenseModal'
import { AddDeliverableModal } from '../../components/AddDeliverableModal'
import type {
  Project,
  ProjectStage,
  ProjectStatus,
  ProjectDocument,
  ProjectSession,
  ProjectPrintingItem,
  CustomerDeliverable,
  ProjectExpenseItem,
} from '../../types/project'

interface ProjectEditPageProps {
  project: Project
  onBack: () => void
  onUpdate: (updated: Project) => void
  onDelete: (id: string) => void
}

type ProjectTab =
  | 'Documents'
  | 'Insights'
  | 'Sessions'
  | 'Tasks'
  | 'Printing'
  | 'Deliverables'
  | 'Expense'
  | 'Notes'

const TABS: { id: ProjectTab; label: string; icon: any }[] = [
  { id: 'Documents', label: 'Documents', icon: IconFileText },
  { id: 'Insights', label: 'Insights', icon: IconChartBar },
  { id: 'Sessions', label: 'Sessions', icon: IconCalendar },
  { id: 'Tasks', label: 'Tasks', icon: IconChecklist },
  { id: 'Printing', label: 'Printing', icon: IconPrinter },
  { id: 'Deliverables', label: 'Deliverables', icon: IconPhoto },
  { id: 'Expense', label: 'Expense', icon: IconReceipt2 },
  { id: 'Notes', label: 'Notes', icon: IconNote },
]

const STAGES: ProjectStage[] = [
  'Booked',
  'Shoot Planning',
  'Shoot Completed',
  'Editing',
  'Delivery',
  'Closure',
]

const STATUSES: ProjectStatus[] = [
  'New',
  'Booked',
  'Active',
  'Completed',
  'On Hold',
  'Cancelled',
]

export function ProjectEditPage({
  project,
  onBack,
  onUpdate,
  onDelete,
}: ProjectEditPageProps) {
  const { workflows } = useWorkflows('project')
  const [activeTab, setActiveTab] = React.useState<ProjectTab>('Documents')
  const [status, setStatus] = React.useState<ProjectStatus>(project.status)
  const [stage, setStage] = React.useState<ProjectStage>(project.stage)
  const [projectName, setProjectName] = React.useState(project.projectName)
  const [isEditingTitle, setIsEditingTitle] = React.useState(false)

  const dynamicStages = React.useMemo(() => {
    if (workflows && workflows.length > 0) {
      return workflows.map((w) => w.stage)
    }
    return STAGES
  }, [workflows])

  const stageColor = React.useMemo(() => {
    const item = workflows?.find((w) => w.stage.toLowerCase() === stage.toLowerCase())
    return item?.color_code || '#2563eb'
  }, [workflows, stage])

  // Sub-items states
  const [documents] = React.useState<ProjectDocument[]>(project.documents || [])
  const [sessions, setSessions] = React.useState<ProjectSession[]>(project.sessions || [])
  const [printing] = React.useState<ProjectPrintingItem[]>(project.printing || [])
  const [deliverables, setDeliverables] = React.useState<CustomerDeliverable[]>(
    project.deliverables || []
  )
  const [expenses, setExpenses] = React.useState<ProjectExpenseItem[]>(project.expenses || [])
  const [notes, setNotes] = React.useState<string>(project.notes || '')

  // Insights Sub-tab
  const [insightSubTab, setInsightSubTab] = React.useState<'Overview' | 'Finance' | 'Payment History'>(
    'Overview'
  )

  // Shoot Expense Inline Edit State in Overview
  const [isEditingShootExpense, setIsEditingShootExpense] = React.useState(false)
  const [shootExpenseInput, setShootExpenseInput] = React.useState<number>(project.shootExpense ?? 0)

  // Add Deliverable Modal State
  const [isAddDeliverableModalOpen, setIsAddDeliverableModalOpen] = React.useState(false)

  // Quick Inline Note State for Deliverables
  const [activeNoteDelId, setActiveNoteDelId] = React.useState<string | null>(null)
  const [noteDraftText, setNoteDraftText] = React.useState<string>('')

  // Status/Stage changes trigger automatic update
  const handleStatusChange = (newStatus: ProjectStatus) => {
    setStatus(newStatus)
    const updated = { ...project, status: newStatus }
    onUpdate(updated)
    toast.success(`Project status updated to ${newStatus}`)
  }

  const handleStageChange = (newStage: ProjectStage) => {
    setStage(newStage)
    const updated = { ...project, stage: newStage }
    onUpdate(updated)
    toast.success(`Project stage moved to ${newStage}`)
  }

  const handleSaveNotes = (newNotesHtml: string) => {
    setNotes(newNotesHtml)
    onUpdate({ ...project, notes: newNotesHtml })
  }

  // Add/Edit session drawer
  const [isSessionDrawerOpen, setIsSessionDrawerOpen] = React.useState(false)
  const [editingSessionItem, setEditingSessionItem] = React.useState<ProjectSession | null>(null)

  // Add/Edit expense modal
  const [isExpenseModalOpen, setIsExpenseModalOpen] = React.useState(false)
  const [editingExpenseItem, setEditingExpenseItem] = React.useState<ProjectExpenseItem | null>(null)

  const handleOpenAddExpense = () => {
    setEditingExpenseItem(null)
    setIsExpenseModalOpen(true)
  }

  const handleOpenEditExpense = (exp: ProjectExpenseItem) => {
    setEditingExpenseItem(exp)
    setIsExpenseModalOpen(true)
  }

  const handleDeleteExpense = (id: string) => {
    const updated = expenses.filter((e) => e.id !== id)
    setExpenses(updated)
    const totalExp = updated.reduce((a, b) => a + b.amount, 0)
    onUpdate({ ...project, expenses: updated, shootExpense: totalExp })
    toast.success('Expense deleted')
  }

  const handleSaveExpense = (savedExpense: ProjectExpenseItem) => {
    let updated: ProjectExpenseItem[]
    const exists = expenses.some((e) => e.id === savedExpense.id)
    if (exists) {
      updated = expenses.map((e) => (e.id === savedExpense.id ? savedExpense : e))
    } else {
      updated = [savedExpense, ...expenses]
    }
    setExpenses(updated)
    const totalExp = updated.reduce((a, b) => a + b.amount, 0)
    onUpdate({ ...project, expenses: updated, shootExpense: totalExp })
  }

  const handleOpenAddSession = () => {
    setEditingSessionItem(null)
    setIsSessionDrawerOpen(true)
  }

  const handleOpenEditSession = (ses: ProjectSession) => {
    setEditingSessionItem(ses)
    setIsSessionDrawerOpen(true)
  }

  const handleSaveSessionFromDrawer = (savedSession: ProjectSession) => {
    setSessions((prev) => {
      const exists = prev.some((s) => s.id === savedSession.id)
      const next = exists
        ? prev.map((s) => (s.id === savedSession.id ? savedSession : s))
        : [...prev, savedSession]
      onUpdate({ ...project, sessions: next, sessionsCount: next.length })
      return next
    })
  }

  return (
    <div className="space-y-6 w-full max-w-7xl mx-auto min-w-0 pb-16">
      {/* Top Bar: Back, KW Avatar, Title, Status & Stage Pill Selectors, Actions (matches image) */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b pb-4">
        <div className="flex items-center gap-3">
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={onBack}
            className="h-9 w-9 p-0 rounded-full"
          >
            <IconArrowLeft className="h-4 w-4" />
          </Button>

          {/* Client Initials Badge Box */}
          <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-sky-100 text-sky-700 dark:bg-sky-950 dark:text-sky-300 font-bold text-sm">
            {project.clientInitials || 'KW'}
          </div>

          <div>
            <div className="flex items-center gap-2">
              {isEditingTitle ? (
                <Input
                  value={projectName}
                  onChange={(e) => setProjectName(e.target.value)}
                  onBlur={() => {
                    setIsEditingTitle(false)
                    onUpdate({ ...project, projectName })
                  }}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                      setIsEditingTitle(false)
                      onUpdate({ ...project, projectName })
                    }
                  }}
                  className="h-8 text-base font-bold w-64"
                  autoFocus
                />
              ) : (
                <h1 className="text-xl font-bold tracking-tight text-foreground">
                  {projectName}
                </h1>
              )}
              <button
                type="button"
                onClick={() => setIsEditingTitle(true)}
                className="text-muted-foreground hover:text-foreground"
                title="Edit title"
              >
                <IconEdit className="h-4 w-4" />
              </button>
            </div>
            <p className="text-xs text-muted-foreground">{project.clientName}</p>
          </div>
        </div>

        {/* Status, Stage, and Action Menu */}
        <div className="flex items-center gap-3 flex-wrap">
          {/* Status Select */}
          <div className="flex flex-col gap-0.5">
            <span className="text-[10px] uppercase font-bold text-muted-foreground">Status</span>
            <Select value={status} onValueChange={(val) => handleStatusChange(val as ProjectStatus)}>
              <SelectTrigger className="h-8 text-xs font-semibold rounded-full bg-muted/60 px-3 border-0 min-w-[100px]">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {STATUSES.map((st) => (
                  <SelectItem key={st} value={st} className="text-xs">
                    {st}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {/* Stage Select with Pill */}
          <div className="flex flex-col gap-0.5">
            <span className="text-[10px] uppercase font-bold text-muted-foreground">Stage</span>
            <Select value={stage} onValueChange={(val) => handleStageChange(val as ProjectStage)}>
              <SelectTrigger
                className="h-8 text-xs font-semibold rounded-full text-white px-3 border-0 min-w-[120px]"
                style={{ backgroundColor: stageColor }}
              >
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {dynamicStages.map((sg) => (
                  <SelectItem key={sg} value={sg} className="text-xs font-medium">
                    {sg}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {/* Actions Dropdown */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="h-8 text-xs font-medium rounded-lg gap-1 mt-auto">
                <span>Action</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="text-xs">
              <DropdownMenuItem onClick={() => toast.success('Contract generated')}>
                Generate Contract PDF
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => toast.success('Share link copied to clipboard')}>
                Copy Client Share Link
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => toast.success('Exported to Google Calendar')}>
                Sync Sessions to Calendar
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => {
              if (confirm('Are you sure you want to delete this project?')) {
                onDelete(project.id)
                toast.success('Project deleted')
                onBack()
              }
            }}
            className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive mt-auto"
            title="Delete project"
          >
            <IconTrash className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {/* Tabs Row: Documents, Insights, Sessions, Tasks, Printing, Deliverables, Expense, Notes */}
      <div className="border-b">
        <div className="flex items-center gap-1 overflow-x-auto scrollbar-none -mb-px">
          {TABS.map((tab) => {
            const Icon = tab.icon
            const isActive = activeTab === tab.id
            return (
              <button
                key={tab.id}
                type="button"
                onClick={() => setActiveTab(tab.id)}
                className={cn(
                  'flex items-center gap-2 px-4 py-2.5 text-xs font-semibold rounded-t-lg transition border-b-2 whitespace-nowrap',
                  isActive
                    ? 'border-primary text-primary font-bold'
                    : 'border-transparent text-muted-foreground hover:text-foreground hover:bg-muted/40'
                )}
              >
                <Icon className="h-4 w-4" />
                <span>{tab.label}</span>
              </button>
            )
          })}
        </div>
      </div>

      {/* TAB 1: DOCUMENTS (matches image 2) */}
      {activeTab === 'Documents' && (
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
          {/* Left: Documents Table */}
          <Card className="lg:col-span-8 shadow-xs border bg-card">
            <CardHeader className="pb-3 flex flex-row items-center justify-between">
              <div>
                <CardTitle className="text-base font-semibold">Documents</CardTitle>
                <CardDescription className="text-xs">
                  Invoices, signed contracts, lead capture forms, and quotes.
                </CardDescription>
              </div>
              <Button
                size="sm"
                className="bg-emerald-600 hover:bg-emerald-700 text-white text-xs h-8"
                onClick={() => toast.success('New document created')}
              >
                Create
              </Button>
            </CardHeader>
            <CardContent>
              <div className="divide-y divide-border/60">
                {documents.map((doc) => (
                  <div
                    key={doc.id}
                    className="flex items-center justify-between py-3 px-2 text-xs hover:bg-muted/30 rounded-lg transition"
                  >
                    <div className="flex items-center gap-3">
                      <Badge variant="outline" className="text-[11px] font-mono px-2 py-0.5 rounded-full">
                        {doc.type}
                      </Badge>
                      <span className="font-semibold text-foreground">{doc.title}</span>
                    </div>

                    <div className="flex items-center gap-6">
                      <span className="text-muted-foreground font-mono text-[11px]">
                        {doc.referenceNo}
                      </span>
                      {doc.amount !== undefined && (
                        <span className="font-bold tabular-nums text-foreground">
                          INR {doc.amount.toLocaleString('en-IN')}
                        </span>
                      )}
                      <div className="flex items-center gap-1.5">
                        <span
                          className={cn(
                            'h-1.5 w-1.5 rounded-full',
                            doc.status === 'Submitted' && 'bg-blue-500',
                            doc.status === 'Converted' || doc.status === 'Signed' || doc.status === 'Paid'
                              ? 'bg-emerald-500'
                              : 'bg-amber-500'
                          )}
                        />
                        <span className="text-[11px] font-medium">{doc.status}</span>
                      </div>
                      <span className="text-muted-foreground tabular-nums text-[11px]">
                        {doc.date}
                      </span>
                      <IconDotsVertical className="h-3.5 w-3.5 text-muted-foreground/60 cursor-pointer" />
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Right: Gallery Box (matches image 2) */}
          <Card className="lg:col-span-4 shadow-xs border bg-card">
            <CardHeader className="pb-3 flex flex-row items-center justify-between">
              <CardTitle className="text-base font-semibold">Gallery</CardTitle>
              <Button
                variant="outline"
                size="sm"
                className="text-xs h-7 px-2.5"
                onClick={() => toast.success('Upload gallery dialog opened')}
              >
                Add
              </Button>
            </CardHeader>
            <CardContent>
              <div className="flex flex-col items-center justify-center rounded-xl border border-dashed py-14 text-center bg-muted/10">
                <IconPhoto className="h-8 w-8 text-muted-foreground/40 mb-2" />
                <p className="text-xs text-muted-foreground font-medium">No Data</p>
                <p className="text-[11px] text-muted-foreground/80 mt-1">
                  Upload client proofing galleries or share links.
                </p>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* TAB 2: INSIGHTS (matches image 3) */}
      {activeTab === 'Insights' && (
        <div className="grid grid-cols-1 md:grid-cols-12 gap-6 items-start">
          {/* Left Vertical Nav: Overview, Finance, Payment History */}
          <div className="md:col-span-3 space-y-1 rounded-xl border bg-card p-3 shadow-xs">
            {(['Overview', 'Finance', 'Payment History'] as const).map((sub) => (
              <button
                key={sub}
                type="button"
                onClick={() => setInsightSubTab(sub)}
                className={cn(
                  'w-full text-left px-3 py-2 rounded-lg text-xs font-semibold transition',
                  insightSubTab === sub
                    ? 'bg-emerald-50 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-300'
                    : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                )}
              >
                {sub}
              </button>
            ))}
          </div>

          {/* Right Insight Panels */}
          <div className="md:col-span-9 space-y-6">
            {insightSubTab === 'Overview' && (
              <div className="space-y-6">
                <h3 className="text-base font-bold text-foreground">Overview</h3>

                {/* Projects Time Block */}
                <Card className="shadow-xs border bg-card">
                  <CardHeader className="pb-3">
                    <CardTitle className="text-sm font-semibold">Projects</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                      <div>
                        <span className="text-xs text-muted-foreground font-medium">
                          Projects Start Date
                        </span>
                        <p className="text-sm font-bold text-foreground mt-1 tabular-nums">
                          {project.startDate.replace('T', ' ')}
                        </p>
                      </div>
                      <div>
                        <span className="text-xs text-muted-foreground font-medium">
                          Projects End Date
                        </span>
                        <p className="text-sm font-bold text-foreground mt-1 tabular-nums">
                          {project.endDate.replace('T', ' ')}
                        </p>
                      </div>
                      <div>
                        <span className="text-xs text-muted-foreground font-medium">
                          Projects Duration
                        </span>
                        <p className="text-sm font-bold text-foreground mt-1 tabular-nums text-primary">
                          3 Days
                        </p>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                {/* Shoot Details Block */}
                <Card className="shadow-xs border bg-card">
                  <CardHeader className="pb-3">
                    <CardTitle className="text-sm font-semibold">Shoot Details</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 items-center">
                      <div>
                        <span className="text-xs text-muted-foreground font-medium">Shoot Hours</span>
                        <p className="text-sm font-bold text-foreground mt-1 tabular-nums">
                          {project.shootHours.toFixed(1)} Hours
                        </p>
                      </div>
                      <div>
                        <span className="text-xs text-muted-foreground font-medium">Shoot Members</span>
                        <div className="flex items-center gap-1.5 mt-1">
                          <span className="text-sm font-bold tabular-nums">
                            {project.shootMembersCount}
                          </span>
                          <span className="flex h-5 w-5 items-center justify-center rounded-full bg-primary/10 text-[10px] font-bold text-primary">
                            E
                          </span>
                        </div>
                      </div>
                      <div>
                        <span className="text-xs text-muted-foreground font-medium">Shoot Expense</span>
                        <div className="flex items-center gap-2 mt-1">
                          {isEditingShootExpense ? (
                            <div className="flex items-center gap-1.5">
                              <span className="text-xs font-bold text-muted-foreground">INR</span>
                              <Input
                                type="number"
                                min="0"
                                value={shootExpenseInput}
                                onChange={(e) => setShootExpenseInput(Number(e.target.value) || 0)}
                                onKeyDown={(e) => {
                                  if (e.key === 'Enter') {
                                    setIsEditingShootExpense(false)
                                    onUpdate({ ...project, shootExpense: shootExpenseInput })
                                    toast.success('Shoot expense updated')
                                  } else if (e.key === 'Escape') {
                                    setIsEditingShootExpense(false)
                                    setShootExpenseInput(project.shootExpense ?? 0)
                                  }
                                }}
                                className="h-7 w-24 text-xs font-bold px-2 tabular-nums"
                                autoFocus
                              />
                              <button
                                type="button"
                                onClick={() => {
                                  setIsEditingShootExpense(false)
                                  onUpdate({ ...project, shootExpense: shootExpenseInput })
                                  toast.success('Shoot expense updated')
                                }}
                                className="h-6 w-6 rounded flex items-center justify-center bg-emerald-600 hover:bg-emerald-700 text-white"
                                title="Save shoot expense"
                              >
                                <IconCheck className="h-3.5 w-3.5" />
                              </button>
                            </div>
                          ) : (
                            <>
                              <p className="text-sm font-bold text-foreground tabular-nums">
                                INR {project.shootExpense}
                              </p>
                              <button
                                type="button"
                                onClick={() => {
                                  setShootExpenseInput(project.shootExpense ?? 0)
                                  setIsEditingShootExpense(true)
                                }}
                                className="text-muted-foreground hover:text-foreground p-0.5 rounded transition"
                                title="Edit shoot expense"
                              >
                                <IconEdit className="h-3.5 w-3.5" />
                              </button>
                            </>
                          )}
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            )}

            {insightSubTab === 'Finance' && (
              <Card className="shadow-xs border bg-card p-6 space-y-6">
                <div>
                  <h3 className="text-base font-bold text-foreground">Finance</h3>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-y-6 gap-x-8">
                  {/* Row 1 */}
                  <div>
                    <span className="text-xs font-semibold text-foreground">Profit</span>
                    <p className="text-base font-bold text-foreground mt-1 tabular-nums">
                      <span className="text-xs font-normal text-muted-foreground mr-1">INR</span>
                      {project.finance?.profit?.toLocaleString('en-IN') ?? '49,500'}
                    </p>
                  </div>
                  <div>
                    <span className="text-xs font-semibold text-emerald-600 dark:text-emerald-400">Current Profit</span>
                    <p className="text-base font-bold text-foreground mt-1 tabular-nums">
                      <span className="text-xs font-normal text-muted-foreground mr-1">INR</span>
                      {project.finance?.currentProfit?.toLocaleString('en-IN') ?? '25,500'}
                    </p>
                  </div>
                  <div></div>

                  {/* Row 2 */}
                  <div>
                    <span className="text-xs font-semibold text-foreground">Net Payable Amount</span>
                    <p className="text-base font-bold text-foreground mt-1 tabular-nums">
                      <span className="text-xs font-normal text-muted-foreground mr-1">INR</span>
                      {project.finance?.netPayableAmount?.toLocaleString('en-IN') ?? '54,000'}
                    </p>
                  </div>
                  <div>
                    <span className="text-xs font-semibold text-foreground">Paid Amount</span>
                    <p className="text-base font-bold text-foreground mt-1 tabular-nums">
                      <span className="text-xs font-normal text-muted-foreground mr-1">INR</span>
                      {project.finance?.paidAmount?.toLocaleString('en-IN') ?? '30,000'}
                    </p>
                  </div>
                  <div>
                    <span className="text-xs font-semibold text-foreground">Balance Amount</span>
                    <p className="text-base font-bold text-foreground mt-1 tabular-nums">
                      <span className="text-xs font-normal text-muted-foreground mr-1">INR</span>
                      {project.finance?.balanceAmount?.toLocaleString('en-IN') ?? '24,000'}
                    </p>
                  </div>

                  {/* Row 3 */}
                  <div>
                    <span className="text-xs font-semibold text-foreground">Payment Due</span>
                    <p className="text-base font-bold text-foreground mt-1 tabular-nums">
                      <span className="text-xs font-normal text-muted-foreground mr-1">INR</span>
                      {project.finance?.paymentDue?.toLocaleString('en-IN') ?? '0'}
                    </p>
                  </div>
                  <div>
                    <span className="text-xs font-semibold text-foreground">Expense</span>
                    <p className="text-base font-bold text-foreground mt-1 tabular-nums">
                      <span className="text-xs font-normal text-muted-foreground mr-1">INR</span>
                      {project.finance?.expense?.toLocaleString('en-IN') ?? '4,500'}
                    </p>
                  </div>
                  <div>
                    <span className="text-xs font-semibold text-foreground">Event Discount Amount</span>
                    <p className="text-base font-bold text-foreground mt-1 tabular-nums">
                      <span className="text-xs font-normal text-muted-foreground mr-1">INR</span>
                      {project.finance?.eventDiscountAmount?.toLocaleString('en-IN') ?? '6,000'}
                    </p>
                  </div>
                </div>
              </Card>
            )}

            {insightSubTab === 'Payment History' && (
              <Card className="shadow-xs border bg-card p-6 space-y-4">
                <div>
                  <h3 className="text-base font-bold text-foreground">Payment History</h3>
                </div>

                <div className="rounded-xl border border-border/80 overflow-hidden">
                  <table className="w-full text-left text-xs">
                    <thead className="bg-muted/30 border-b text-[10px] font-bold text-muted-foreground uppercase tracking-wider">
                      <tr>
                        <th className="px-4 py-3">CREATED DATE</th>
                        <th className="px-4 py-3">PAYMENT DATE</th>
                        <th className="px-4 py-3">AMOUNT</th>
                        <th className="px-4 py-3">TYPE</th>
                        <th className="px-4 py-3">REFERENCENO</th>
                        <th className="px-4 py-3">STATUS</th>
                        <th className="px-4 py-3">NOTE</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-border/60">
                      {(project.finance?.paymentHistory && project.finance.paymentHistory.length > 0
                        ? project.finance.paymentHistory
                        : [
                            {
                              id: 'pay-1',
                              createdDate: 'Nov 01, 2026 10:41 PM',
                              paymentDate: 'Nov 01, 2026 12:00 PM',
                              amount: 30000,
                              type: 'UPI Payment',
                              referenceNo: '123456789012',
                              status: 'Payment Due',
                              note: '',
                            },
                          ]
                      ).map((pay) => (
                        <tr key={pay.id} className="hover:bg-muted/10 transition">
                          <td className="px-4 py-3 tabular-nums text-foreground">{pay.createdDate}</td>
                          <td className="px-4 py-3 tabular-nums text-foreground">{pay.paymentDate}</td>
                          <td className="px-4 py-3 font-bold tabular-nums text-foreground">
                            ₹{pay.amount.toLocaleString('en-IN')}
                          </td>
                          <td className="px-4 py-3 text-muted-foreground">{pay.type}</td>
                          <td className="px-4 py-3 font-mono text-[11px] text-muted-foreground">
                            {pay.referenceNo}
                          </td>
                          <td className="px-4 py-3 font-medium text-foreground">{pay.status}</td>
                          <td className="px-4 py-3 text-muted-foreground">{pay.note || '—'}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </Card>
            )}
          </div>
        </div>
      )}

      {/* TAB 3: SESSIONS (matches image 4) */}
      {activeTab === 'Sessions' && (
        <Card className="shadow-xs border bg-card">
          <CardHeader className="pb-3 flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-base font-semibold">Sessions</CardTitle>
              <CardDescription className="text-xs">
                Timeline breakdown for shoot days, timings, and crew members.
              </CardDescription>
            </div>
            <Button
              size="sm"
              onClick={handleOpenAddSession}
              className="bg-emerald-600 hover:bg-emerald-700 text-white text-xs h-8 gap-1"
            >
              <IconPlus className="h-3.5 w-3.5" />
              New Session
            </Button>
          </CardHeader>
          <CardContent className="space-y-4">
            <span className="text-xs font-semibold text-muted-foreground">November, 2026</span>

            {sessions.length === 0 ? (
              <div className="flex flex-col items-center justify-center p-8 text-center rounded-lg border border-dashed">
                <IconInbox className="h-10 w-10 text-muted-foreground/50 mb-2" />
                <h4 className="text-sm font-semibold text-foreground">No shoot sessions scheduled yet</h4>
                <p className="text-xs text-muted-foreground mt-0.5 mb-3">
                  Schedule sessions for ceremony days, portraits, and crew assignments.
                </p>
                <Button
                  size="sm"
                  onClick={handleOpenAddSession}
                  className="bg-emerald-600 hover:bg-emerald-700 text-white text-xs gap-1.5"
                >
                  <IconPlus className="h-3.5 w-3.5" />
                  Schedule First Session
                </Button>
              </div>
            ) : (
              <div className="space-y-3">
                {sessions.map((ses) => (
                  <div
                    key={ses.id}
                    className="rounded-xl border bg-muted/10 p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-4 transition hover:border-border"
                  >
                    {/* Left: Day & Date badge */}
                    <div className="flex items-center gap-4">
                      <div className="flex flex-col items-center justify-center rounded-lg bg-card border px-3 py-2 min-w-[70px] shadow-xs">
                        <span className="text-[11px] font-semibold text-muted-foreground uppercase">
                          {ses.date.split(',')[0]}
                        </span>
                        <span className="text-2xl font-bold tabular-nums text-foreground">
                          {ses.date.match(/\d+/)?.[0] || '11'}
                        </span>
                      </div>

                      {/* Middle: Details */}
                      <div className="space-y-1">
                        <div className="flex items-center gap-2 text-xs text-muted-foreground">
                          <IconClock className="h-3.5 w-3.5" />
                          <span>{ses.timeRange}</span>
                        </div>
                        <h4
                          onClick={() => handleOpenEditSession(ses)}
                          className="font-bold text-sm text-foreground hover:text-primary transition cursor-pointer"
                        >
                          {ses.title}
                        </h4>
                        <p className="text-xs text-muted-foreground">
                          {ses.venue} • {ses.city}
                        </p>
                        <p className="text-[11px] text-muted-foreground">
                          {ses.members.length > 0
                            ? `Crew: ${ses.members.join(', ')}`
                            : 'No Members Added'}
                        </p>
                      </div>
                    </div>

                    {/* Right: Status badge & Shoots Activity button */}
                    <div className="flex sm:flex-col items-end gap-2 justify-between">
                      <Badge variant="outline" className="text-xs font-medium">
                        {ses.status}
                      </Badge>
                      <button
                        type="button"
                        onClick={() => toast.info(`Viewing activity for ${ses.title}`)}
                        className="text-xs text-primary hover:underline flex items-center gap-1"
                      >
                        <span>Shoots Activity</span>
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* TAB 4: TASKS - Full Production Workflow Matrix */}
      {activeTab === 'Tasks' && (
        <div className="w-full min-w-0 overflow-hidden">
          <ProductionMatrixView
            projectId={project.id}
            projectName={project.projectName}
          />
        </div>
      )}

      {/* TAB 5: PRINTING (matches Image 4) */}
      {activeTab === 'Printing' && (
        <Card className="shadow-xs border bg-card">
          <CardContent className="p-4 sm:p-6 space-y-4">
            {/* Top Toolbar */}
            <div className="flex flex-wrap items-center justify-between gap-4">
              <div className="flex items-center gap-4">
                <div className="flex items-center gap-1.5 text-xs font-semibold text-muted-foreground">
                  <span className="flex h-4 w-4 items-center justify-center rounded-full border border-muted-foreground/60">
                    <IconCheck className="h-2.5 w-2.5" />
                  </span>
                  <span>Pending Albums</span>
                </div>
                <div className="flex items-center gap-1.5 text-xs font-semibold text-emerald-600 dark:text-emerald-400">
                  <span className="flex h-4 w-4 items-center justify-center rounded-full border border-emerald-600">
                    <IconCheck className="h-2.5 w-2.5" />
                  </span>
                  <span>Pending Decors</span>
                </div>
                <button
                  type="button"
                  className="flex items-center gap-1 text-xs font-semibold text-muted-foreground hover:text-foreground"
                >
                  <IconDotsVertical className="h-3.5 w-3.5 rotate-90" />
                  <span>Filters</span>
                </button>
              </div>

              <button
                type="button"
                className="flex items-center gap-1.5 text-xs font-semibold text-muted-foreground hover:text-foreground"
              >
                <IconEdit className="h-3.5 w-3.5" />
                <span>Edit Columns</span>
              </button>
            </div>

            {/* Printing Table (Matches Image 4) */}
            <div className="rounded-xl border border-border/80 overflow-hidden">
              <table className="w-full text-left text-xs">
                <thead className="bg-muted/40 border-b text-[11px] font-bold text-muted-foreground uppercase tracking-wider">
                  <tr>
                    <th className="px-4 py-3">ID</th>
                    <th className="px-4 py-3">LAB NAME</th>
                    <th className="px-4 py-3">GALLERIES NAME</th>
                    <th className="px-4 py-3">ORDER NAME</th>
                    <th className="px-4 py-3">ORDER TYPE</th>
                    <th className="px-4 py-3">ASSOCIATE</th>
                    <th className="px-4 py-3">STATUS</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border/60">
                  {printing.length === 0 ? (
                    <tr>
                      <td colSpan={7} className="px-4 py-16 text-center">
                        <div className="flex flex-col items-center justify-center">
                          <IconInbox className="h-10 w-10 text-muted-foreground/40 mb-2" />
                          <span className="text-xs font-medium text-muted-foreground">No Data</span>
                        </div>
                      </td>
                    </tr>
                  ) : (
                    printing.map((prn) => (
                      <tr key={prn.id} className="hover:bg-muted/10 transition">
                        <td className="px-4 py-3 font-mono text-[11px] text-muted-foreground">{prn.id}</td>
                        <td className="px-4 py-3 font-semibold text-foreground">{prn.labName}</td>
                        <td className="px-4 py-3 text-muted-foreground">{prn.galleriesName}</td>
                        <td className="px-4 py-3 font-medium text-foreground">{prn.orderName}</td>
                        <td className="px-4 py-3 text-muted-foreground">{prn.orderType}</td>
                        <td className="px-4 py-3 text-muted-foreground">{prn.associate}</td>
                        <td className="px-4 py-3">
                          <Badge variant="outline" className="text-[10px]">
                            {prn.status}
                          </Badge>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>
      )}

      {/* TAB 6: DELIVERABLES (Matches Image 5: Customer Deliverables) */}
      {activeTab === 'Deliverables' && (
        <Card className="shadow-xs border bg-card">
          <CardHeader className="pb-3 flex flex-row items-center justify-between">
            <CardTitle className="text-base font-bold text-foreground">
              Customer Deliverables
            </CardTitle>
            <Button
              size="sm"
              onClick={() => setIsAddDeliverableModalOpen(true)}
              className="bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs h-8 shadow-xs"
            >
              Add Deliverables
            </Button>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="rounded-xl border border-border/80 overflow-hidden">
              <table className="w-full text-left text-xs">
                <thead className="bg-muted/40 border-b text-[11px] font-bold text-muted-foreground uppercase tracking-wider">
                  <tr>
                    <th className="px-4 py-3">Deliverable</th>
                    <th className="px-4 py-3">Due Date</th>
                    <th className="px-4 py-3">Assignee</th>
                    <th className="px-4 py-3">Status</th>
                    <th className="px-4 py-3">Note</th>
                    <th className="px-4 py-3 text-right">Action</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border/60">
                  {deliverables.length === 0 ? (
                    <tr>
                      <td colSpan={6} className="px-4 py-8 text-center text-muted-foreground">
                        No customer deliverables scheduled yet. Click "Add Deliverables" above.
                      </td>
                    </tr>
                  ) : (
                    deliverables.map((del) => (
                      <tr key={del.id} className="hover:bg-muted/10 transition">
                        <td className="px-4 py-3 font-semibold text-foreground whitespace-nowrap">
                          {del.deliverable}
                        </td>
                        <td className="px-4 py-3">
                          <div className="flex items-center gap-1.5 rounded-md border bg-background px-2.5 py-1 text-xs w-36">
                            <span className="truncate">
                              {del.dueDate ? new Date(del.dueDate).toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' }) : 'Set Date'}
                            </span>
                            <IconCalendar className="h-3.5 w-3.5 text-muted-foreground ml-auto opacity-70 shrink-0" />
                          </div>
                        </td>
                        <td className="px-4 py-3 text-muted-foreground whitespace-nowrap">
                          {del.assignee ? (
                            <span className="font-medium text-foreground">{del.assignee}</span>
                          ) : (
                            <span className="text-muted-foreground/60 italic">Assigned Team Member</span>
                          )}
                        </td>
                        <td className="px-4 py-3">
                          <Select
                            value={del.status}
                            onValueChange={(val) => {
                              const updated = deliverables.map((d) =>
                                d.id === del.id ? { ...d, status: val as CustomerDeliverable['status'] } : d
                              )
                              setDeliverables(updated)
                              onUpdate({ ...project, deliverables: updated })
                            }}
                          >
                            <SelectTrigger className="h-8 w-32 text-xs bg-background">
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="New" className="text-xs">New</SelectItem>
                              <SelectItem value="In Progress" className="text-xs">In Progress</SelectItem>
                              <SelectItem value="Client Review" className="text-xs">Client Review</SelectItem>
                              <SelectItem value="Completed" className="text-xs">Completed</SelectItem>
                              <SelectItem value="Delivered" className="text-xs">Delivered</SelectItem>
                            </SelectContent>
                          </Select>
                        </td>
                        <td className="px-4 py-3">
                          <Popover
                            open={activeNoteDelId === del.id}
                            onOpenChange={(isOpen) => {
                              if (isOpen) {
                                setActiveNoteDelId(del.id)
                                setNoteDraftText(del.note || '')
                              } else {
                                setActiveNoteDelId(null)
                              }
                            }}
                          >
                            <PopoverTrigger asChild>
                              <Button
                                size="sm"
                                className="bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs h-7 px-3"
                              >
                                Note
                              </Button>
                            </PopoverTrigger>
                            <PopoverContent align="start" className="w-64 p-3 space-y-2 shadow-lg">
                              <textarea
                                value={noteDraftText}
                                onChange={(e) => setNoteDraftText(e.target.value)}
                                placeholder="Enter Your Note Here..."
                                className="w-full h-20 text-xs border rounded-md p-2 resize-none bg-background focus:outline-none focus:ring-1 focus:ring-ring"
                                autoFocus
                              />
                              <div className="flex justify-end gap-1.5">
                                <Button
                                  size="sm"
                                  variant="ghost"
                                  className="h-6 px-2 text-[11px]"
                                  onClick={() => setActiveNoteDelId(null)}
                                >
                                  Cancel
                                </Button>
                                <Button
                                  size="sm"
                                  className="bg-emerald-600 hover:bg-emerald-700 text-white h-6 px-2.5 text-[11px]"
                                  onClick={() => {
                                    const updated = deliverables.map((d) =>
                                      d.id === del.id ? { ...d, note: noteDraftText } : d
                                    )
                                    setDeliverables(updated)
                                    onUpdate({ ...project, deliverables: updated })
                                    setActiveNoteDelId(null)
                                    toast.success('Note saved')
                                  }}
                                >
                                  Save
                                </Button>
                              </div>
                            </PopoverContent>
                          </Popover>
                        </td>
                        <td className="px-4 py-3 text-right">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => {
                              const updated = deliverables.filter((d) => d.id !== del.id)
                              setDeliverables(updated)
                              onUpdate({ ...project, deliverables: updated })
                              toast.success('Deliverable removed')
                            }}
                            className="h-7 w-7 p-0 text-destructive/70 hover:text-destructive hover:bg-destructive/10"
                            title="Delete Deliverable"
                          >
                            <IconTrash className="h-4 w-4" />
                          </Button>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>
      )}

      {/* TAB 7: EXPENSE (matches Image 2 with ₹4,500 header and detailed table) */}
      {activeTab === 'Expense' && (
        <Card className="shadow-xs border bg-card">
          <CardContent className="p-4 sm:p-6 space-y-6">
            {/* Top Metric & Add Expense Action */}
            <div className="flex items-start justify-between">
              <div>
                <div className="text-3xl font-extrabold text-foreground tracking-tight">
                  ₹{expenses.reduce((acc, e) => acc + e.amount, 0).toLocaleString('en-IN')}
                </div>
                <div className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mt-1">
                  Projects EXPENSE
                </div>
              </div>
              <Button
                size="sm"
                onClick={handleOpenAddExpense}
                className="bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs h-8 shadow-xs"
              >
                Add Expense
              </Button>
            </div>

            {/* Expense Details Table */}
            <div className="space-y-3">
              <h4 className="text-sm font-bold text-foreground">Expense Details</h4>

              <div className="rounded-xl border border-border/80 overflow-hidden">
                <table className="w-full text-left text-xs">
                  <thead className="bg-muted/40 border-b text-[11px] font-bold text-muted-foreground uppercase tracking-wider">
                    <tr>
                      <th className="px-4 py-3">SPENT BY</th>
                      <th className="px-4 py-3">SESSION NAME</th>
                      <th className="px-4 py-3">EXPENSE TYPE</th>
                      <th className="px-4 py-3">AMOUNT</th>
                      <th className="px-4 py-3">NOTES</th>
                      <th className="px-4 py-3 text-right">ACTION</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border/60">
                    {expenses.length === 0 ? (
                      <tr>
                        <td colSpan={6} className="px-4 py-8 text-center text-muted-foreground">
                          No expenses recorded for this project yet.
                        </td>
                      </tr>
                    ) : (
                      expenses.map((exp) => (
                        <tr key={exp.id} className="hover:bg-muted/10 transition">
                          <td className="px-4 py-3">
                            <div className="flex h-7 w-7 items-center justify-center rounded-full bg-slate-900 text-white dark:bg-slate-100 dark:text-slate-900 font-bold text-[11px] shadow-xs">
                              {exp.paidBy || 'E'}
                            </div>
                          </td>
                          <td className="px-4 py-3 font-semibold text-foreground">
                            {exp.sessionName || project.projectName}
                          </td>
                          <td className="px-4 py-3 text-muted-foreground">
                            {exp.category}
                          </td>
                          <td className="px-4 py-3 font-bold text-foreground tabular-nums">
                            ₹{exp.amount.toLocaleString('en-IN')}
                          </td>
                          <td className="px-4 py-3 text-muted-foreground max-w-[200px] truncate">
                            {exp.description || '—'}
                          </td>
                          <td className="px-4 py-3 text-right">
                            <DropdownMenu>
                              <DropdownMenuTrigger asChild>
                                <Button variant="ghost" size="sm" className="h-7 w-7 p-0 text-muted-foreground hover:text-foreground">
                                  <IconDotsVertical className="h-4 w-4" />
                                </Button>
                              </DropdownMenuTrigger>
                              <DropdownMenuContent align="end" className="text-xs">
                                <DropdownMenuItem onClick={() => handleOpenEditExpense(exp)}>
                                  <IconEdit className="h-3.5 w-3.5 mr-2" />
                                  Edit
                                </DropdownMenuItem>
                                <DropdownMenuItem
                                  onClick={() => handleDeleteExpense(exp.id)}
                                  className="text-destructive focus:text-destructive"
                                >
                                  <IconTrash className="h-3.5 w-3.5 mr-2" />
                                  Delete
                                </DropdownMenuItem>
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* TAB 8: NOTES (matches image 5 with RichTextEditor) */}
      {activeTab === 'Notes' && (
        <Card className="shadow-xs border bg-card">
          <CardHeader className="pb-3 flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-base font-semibold">Project Notes & Brief</CardTitle>
              <CardDescription className="text-xs">
                Client aesthetic preferences, special family requests, and shot list guidelines.
              </CardDescription>
            </div>
            <Button
              size="sm"
              onClick={() => {
                onUpdate({ ...project, notes })
                toast.success('Notes saved!')
              }}
              className="bg-emerald-600 hover:bg-emerald-700 text-white text-xs h-8"
            >
              Save Notes
            </Button>
          </CardHeader>
          <CardContent>
            <div className="min-h-[300px]">
              <React.Suspense
                fallback={
                  <div className="h-[260px] w-full animate-pulse rounded-xl border border-input bg-muted/20" />
                }
              >
                <RichTextEditor
                  content={notes}
                  onChange={handleSaveNotes}
                  placeholder="Write private studio notes, client must-have shots, or shooting constraints..."
                  className="min-h-[260px] border rounded-xl"
                />
              </React.Suspense>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Reusable Snap Points Drawer for Sessions */}
      <CreateSessionDrawer
        open={isSessionDrawerOpen}
        onOpenChange={setIsSessionDrawerOpen}
        onSave={handleSaveSessionFromDrawer}
        initialSession={editingSessionItem}
        existingProjectNames={[project.projectName]}
      />

      {/* Project Expense Modal */}
      <ProjectExpenseModal
        open={isExpenseModalOpen}
        onOpenChange={setIsExpenseModalOpen}
        onSaveExpense={handleSaveExpense}
        initialExpense={editingExpenseItem}
        existingSessions={sessions.map((s) => s.title)}
      />

      {/* Add Deliverable Modal */}
      <AddDeliverableModal
        open={isAddDeliverableModalOpen}
        onOpenChange={setIsAddDeliverableModalOpen}
        onAdd={(newDel) => {
          const updated = [...deliverables, newDel]
          setDeliverables(updated)
          onUpdate({ ...project, deliverables: updated })
        }}
      />
    </div>
  )
}
