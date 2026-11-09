import { useState, Suspense, lazy } from 'react'
import {
  DashboardShell,
  SplitAuthLayout,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
  Button,
  Badge,
  Checkbox,
  Toaster,
  toast,
  cn,
} from '@repo/ui'
import type { KanbanTask, KanbanColumnDef } from '@repo/ui/components/kanban/kanban'
import { createColumnHelper } from '@tanstack/react-table'
import type { DataTableFeatures, DataTableFilterConfig, DataTableBulkAction, DataTableRowAction, DataTableProps } from '@repo/ui'
import {
  IconCopy,
  IconTrash,
  IconArchive,
  IconCheck,
  IconCalendar,
} from '@tabler/icons-react'

const DataTable = lazy(() => import('@repo/ui/components/data-table').then(m => ({ default: m.DataTable }))) as <TData extends Record<string, any>>(props: DataTableProps<TData>) => React.JSX.Element
const Kanban = lazy(() => import('@repo/ui/components/kanban/kanban').then(m => ({ default: m.Kanban })))
const LoginForm = lazy(() => import('@repo/ui/components/login-form').then(m => ({ default: m.LoginForm })))
const SignupForm = lazy(() => import('@repo/ui/components/signup-form').then(m => ({ default: m.SignupForm })))

type PaymentRecord = {
  id: string
  clientName: string
  email: string
  package: string
  shootDate: string
  amount: number
  status: 'pending' | 'processing' | 'success' | 'failed'
}

const initialPayments: PaymentRecord[] = [
  { id: 'inv-101', clientName: 'Eleanor Vance', email: 'eleanor@vanceweddings.com', package: 'Full Wedding Coverage', shootDate: '2026-12-14', amount: 3200, status: 'success' },
  { id: 'inv-102', clientName: 'Liam Chen', email: 'liam.chen@architects.org', package: 'Commercial Interior', shootDate: '2026-11-20', amount: 1450, status: 'processing' },
  { id: 'inv-103', clientName: 'Sophia Rodriguez', email: 'sophia@rodriguezstudio.com', package: 'Editorial Portrait', shootDate: '2026-11-28', amount: 850, status: 'pending' },
  { id: 'inv-104', clientName: 'Marcus Sterling', email: 'marcus@sterlingcapital.co', package: 'Corporate Headshots', shootDate: '2026-10-15', amount: 2100, status: 'success' },
  { id: 'inv-105', clientName: 'Chloe Dubois', email: 'chloe@atelierdubois.fr', package: 'Fashion Lookbook', shootDate: '2026-11-05', amount: 4800, status: 'failed' },
  { id: 'inv-106', clientName: 'David Kalu', email: 'david@kalucreative.net', package: 'Engagement Session', shootDate: '2026-12-02', amount: 950, status: 'pending' },
  { id: 'inv-107', clientName: 'Amara Patel', email: 'amara.p@voguecraft.io', package: 'Product Catalog', shootDate: '2026-11-12', amount: 2600, status: 'success' },
  { id: 'inv-108', clientName: 'Lucas Bennett', email: 'lucas@bennettevents.com', package: 'Conference Keynote', shootDate: '2026-11-08', amount: 1750, status: 'processing' },
]

const kanbanColumns: KanbanColumnDef[] = [
  { id: 'lead', title: 'New Leads', color: 'blue' },
  { id: 'consultation', title: 'Consultation Scheduled', color: 'purple' },
  { id: 'proposal', title: 'Proposal Sent', color: 'yellow' },
  { id: 'booked', title: 'Booked & Deposit Paid', color: 'green' },
  { id: 'completed', title: 'Delivered / Completed', color: 'gray' },
]

const initialKanbanTasks: KanbanTask[] = [
  {
    id: 'opp-1',
    title: 'Autumn Destination Wedding',
    status: 'consultation',
    customerName: 'Sarah & Jack Reynolds',
    dealValue: 4800,
    priority: 'high',
    dueDate: '2026-11-25',
    assignee: { name: 'Elena Rostova', initials: 'ER' },
    tags: ['Wedding', 'Destination'],
    commentsCount: 3,
    attachmentsCount: 1,
  },
  {
    id: 'opp-2',
    title: 'Brand Launch Lookbook',
    status: 'proposal',
    customerName: 'Aura Athletics',
    dealValue: 3200,
    priority: 'urgent',
    dueDate: '2026-11-18',
    assignee: { name: 'Marcus Sterling', initials: 'MS' },
    tags: ['Commercial', 'Fashion'],
    commentsCount: 5,
    attachmentsCount: 2,
  },
  {
    id: 'opp-3',
    title: 'Executive Team Portraits',
    status: 'booked',
    customerName: 'Apex Financial Group',
    dealValue: 2400,
    priority: 'medium',
    dueDate: '2026-12-04',
    assignee: { name: 'Elena Rostova', initials: 'ER' },
    tags: ['Corporate', 'Studio'],
    commentsCount: 2,
  },
  {
    id: 'opp-4',
    title: 'Culinary Menu Campaign',
    status: 'lead',
    customerName: 'Bistro Lumiere',
    dealValue: 1600,
    priority: 'low',
    dueDate: '2026-11-30',
    tags: ['Food', 'Editorial'],
  },
  {
    id: 'opp-5',
    title: 'Architecture & Interior Feature',
    status: 'completed',
    customerName: 'Studio ArchModern',
    dealValue: 3900,
    priority: 'medium',
    dueDate: '2026-11-02',
    assignee: { name: 'Marcus Sterling', initials: 'MS' },
    tags: ['Architecture'],
  },
]

const columnHelper = createColumnHelper<DataTableFeatures, PaymentRecord>()

const paymentColumns = columnHelper.columns([
  columnHelper.display({
    id: 'select',
    header: ({ table }) => (
      <div className="flex items-center justify-center">
        <Checkbox
          checked={
            table.getIsAllPageRowsSelected() ||
            (table.getIsSomePageRowsSelected() && 'indeterminate')
          }
          onCheckedChange={(value: boolean | 'indeterminate') => table.toggleAllPageRowsSelected(!!value)}
          aria-label="Select all"
        />
      </div>
    ),
    cell: ({ row }) => {
      const isSelected = row.getIsSelected()
      return (
        <div className="relative flex items-center justify-center h-5 w-6 mx-auto">
          {/* Serial number shown when row is not hovered and not selected */}
          <span
            className={cn(
              "text-xs font-mono text-muted-foreground tabular-nums transition-opacity",
              isSelected ? "hidden" : "block group-hover/row:hidden"
            )}
          >
            {row.index + 1}
          </span>
          {/* Checkbox shown on row hover or when row is selected */}
          <div
            className={cn(
              "transition-opacity",
              isSelected ? "block" : "hidden group-hover/row:block"
            )}
            onClick={(e) => e.stopPropagation()}
          >
            <Checkbox
              checked={isSelected}
              onCheckedChange={(value: boolean | 'indeterminate') => row.toggleSelected(!!value)}
              aria-label={`Select row ${row.index + 1}`}
            />
          </div>
        </div>
      )
    },
    enableSorting: false,
    enableHiding: false,
  }),
  columnHelper.accessor('id', {
    header: 'Invoice ID',
    cell: ({ row }) => <span className="font-mono text-xs text-muted-foreground">{row.getValue('id')}</span>,
  }),
  columnHelper.accessor('clientName', {
    header: 'Client',
    cell: ({ row }) => (
      <div>
        <div className="font-medium text-foreground">{row.getValue('clientName')}</div>
        <div className="text-xs text-muted-foreground">{row.original.email}</div>
      </div>
    ),
  }),
  columnHelper.accessor('package', {
    header: 'Shoot Package',
    cell: ({ row }) => <span className="text-sm">{row.getValue('package')}</span>,
  }),
  columnHelper.accessor('shootDate', {
    header: 'Shoot Date',
    cell: ({ row }) => (
      <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
        <IconCalendar className="h-3.5 w-3.5" />
        <span>{row.getValue('shootDate')}</span>
      </div>
    ),
  }),
  columnHelper.accessor('status', {
    header: 'Status',
    cell: ({ row }) => {
      const status = row.getValue('status') as PaymentRecord['status']
      const variant =
        status === 'success'
          ? 'success'
          : status === 'processing'
          ? 'info'
          : status === 'failed'
          ? 'destructive'
          : 'warning'
      return <Badge variant={variant} className="capitalize text-xs">{status}</Badge>
    },
  }),
  columnHelper.accessor('amount', {
    header: () => <div className="text-right">Amount</div>,
    cell: ({ row }) => {
      const amount = parseFloat(row.getValue('amount'))
      const formatted = new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD',
      }).format(amount)
      return <div className="text-right font-medium tabular-nums">{formatted}</div>
    },
  }),
])

const filterConfigs: DataTableFilterConfig[] = [
  {
    columnId: 'status',
    title: 'Status',
    type: 'select',
    options: [
      { label: 'Success', value: 'success' },
      { label: 'Pending', value: 'pending' },
      { label: 'Processing', value: 'processing' },
      { label: 'Failed', value: 'failed' },
    ],
  },
  {
    columnId: 'clientName',
    title: 'Client Name',
    type: 'text',
  },
  {
    columnId: 'package',
    title: 'Package',
    type: 'text',
  },
]

export default function App() {
  const [payments, setPayments] = useState<PaymentRecord[]>(initialPayments)
  const [kanbanTasks, setKanbanTasks] = useState<KanbanTask[]>(initialKanbanTasks)
  const [dashboardTab, setDashboardTab] = useState<'table' | 'kanban'>('kanban')
  const [currentView, setCurrentView] = useState<'dashboard' | 'signin' | 'signup'>(() => {
    const path = window.location.pathname
    if (path === '/signin' || path === '/login') return 'signin'
    if (path === '/signup' || path === '/register') return 'signup'
    return 'dashboard'
  })

  const bulkActions: DataTableBulkAction<PaymentRecord>[] = [
    {
      label: 'Archive',
      icon: IconArchive,
      action: (selected, clearSelection) => {
        toast.success(`Archived ${selected.length} record(s)`)
        clearSelection()
      },
    },
    {
      label: 'Delete',
      icon: IconTrash,
      variant: 'destructive',
      requiresConfirmation: true,
      confirmationTitle: 'Delete Selected Invoices',
      confirmationDescription: (count) =>
        `Are you sure you want to delete ${count} selected invoice(s)? This action will remove records from the CRM and cannot be undone.`,
      action: (selected, clearSelection) => {
        const selectedIds = new Set(selected.map((s) => s.id))
        setPayments((prev) => prev.filter((p) => !selectedIds.has(p.id)))
        toast.success(`Deleted ${selected.length} invoice(s)`)
        clearSelection()
      },
    },
  ]

  const rowActions: DataTableRowAction<PaymentRecord>[] = [
    {
      label: 'Copy Invoice ID',
      icon: IconCopy,
      action: (record) => {
        navigator.clipboard?.writeText(record.id)
        toast.success(`Copied ID ${record.id} to clipboard`)
      },
    },
    {
      label: 'Mark as Paid',
      icon: IconCheck,
      action: (record) => {
        setPayments((prev) =>
          prev.map((p) => (p.id === record.id ? { ...p, status: 'success' } : p))
        )
        toast.success(`Invoice ${record.id} marked as Paid`)
      },
    },
    {
      label: 'Delete Invoice',
      icon: IconTrash,
      variant: 'destructive',
      separatorBefore: true,
      action: (record) => {
        setPayments((prev) => prev.filter((p) => p.id !== record.id))
        toast.success(`Invoice ${record.id} deleted`)
      },
    },
  ]

  if (currentView === 'signin') {
    return (
      <SplitAuthLayout appName="Photography CRM">
        <div className="flex flex-col gap-6">
          <div className="flex items-center justify-between">
            <span className="text-lg font-bold tracking-tight">Photography CRM</span>
            <Button variant="ghost" size="sm" onClick={() => setCurrentView('dashboard')}>
              Back to Dashboard
            </Button>
          </div>
          <Card className="border bg-card shadow-sm">
            <CardContent className="pt-6">
              <Suspense fallback={<div className="h-48 flex items-center justify-center text-sm text-muted-foreground">Loading form...</div>}>
                <LoginForm />
              </Suspense>
            </CardContent>
          </Card>
          <div className="text-center text-sm text-muted-foreground">
            Don't have an account?{' '}
            <button
              onClick={() => setCurrentView('signup')}
              className="text-primary underline underline-offset-4 font-medium"
            >
              Sign up
            </button>
          </div>
        </div>
      </SplitAuthLayout>
    )
  }

  if (currentView === 'signup') {
    return (
      <SplitAuthLayout appName="Photography CRM">
        <div className="flex flex-col gap-6">
          <div className="flex items-center justify-between">
            <span className="text-lg font-bold tracking-tight">Photography CRM</span>
            <Button variant="ghost" size="sm" onClick={() => setCurrentView('dashboard')}>
              Back to Dashboard
            </Button>
          </div>
          <Card className="border bg-card shadow-sm">
            <CardContent className="pt-6">
              <Suspense fallback={<div className="h-64 flex items-center justify-center text-sm text-muted-foreground">Loading form...</div>}>
                <SignupForm />
              </Suspense>
            </CardContent>
          </Card>
          <div className="text-center text-sm text-muted-foreground">
            Already have an account?{' '}
            <button
              onClick={() => setCurrentView('signin')}
              className="text-primary underline underline-offset-4 font-medium"
            >
              Sign in
            </button>
          </div>
        </div>
      </SplitAuthLayout>
    )
  }

  return (
    <DashboardShell title="Photography CRM" breadcrumb="Overview">
      <Toaster position="top-right" />
      <div className="flex items-center justify-between mb-2">
        <div>
          <h1 className="text-lg font-bold tracking-tight">Photography CRM</h1>
          <p className="text-sm text-muted-foreground">
            Shoots, client galleries, bookings and invoices.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <div className="flex items-center rounded-md border bg-muted/40 p-0.5">
            <Button
              variant={dashboardTab === 'kanban' ? 'secondary' : 'ghost'}
              size="sm"
              onClick={() => setDashboardTab('kanban')}
              className="h-7 px-2.5 text-xs font-medium"
            >
              Pipeline (Kanban)
            </Button>
            <Button
              variant={dashboardTab === 'table' ? 'secondary' : 'ghost'}
              size="sm"
              onClick={() => setDashboardTab('table')}
              className="h-7 px-2.5 text-xs font-medium"
            >
              Invoices (Table)
            </Button>
          </div>
          <Button variant="outline" size="sm" onClick={() => setCurrentView('signin')}>
            Switch to Sign In
          </Button>
          <Button variant="outline" size="sm" onClick={() => setCurrentView('signup')}>
            Switch to Sign Up
          </Button>
        </div>
      </div>

      <div className="grid gap-2.5 md:grid-cols-3 mb-2.5">
        <Card className="border bg-card shadow-sm">
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Total Invoices</CardTitle>
            <CardDescription>Records in workspace</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="text-xl font-bold">{payments.length}</div>
            <p className="text-xs text-muted-foreground mt-1">Real-time status</p>
          </CardContent>
        </Card>
        <Card className="border bg-card shadow-sm">
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">System Status</CardTitle>
            <CardDescription>Backend connection</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="text-xl font-bold text-green-600">Connected</div>
            <p className="text-xs text-muted-foreground mt-1">Ready for requests</p>
          </CardContent>
        </Card>
        <Card className="border bg-card shadow-sm">
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">Revenue Billed</CardTitle>
            <CardDescription>Gross invoice total</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="text-xl font-bold tabular-nums">
              ${payments.reduce((acc, p) => acc + p.amount, 0).toLocaleString()}
            </div>
            <p className="text-xs text-muted-foreground mt-1">Active shoot billing</p>
          </CardContent>
        </Card>
      </div>

      {dashboardTab === 'kanban' ? (
        <Card className="border bg-card shadow-xs">
          <CardContent className="p-3">
            <Suspense fallback={<div className="h-96 flex items-center justify-center text-sm text-muted-foreground">Loading pipeline board...</div>}>
              <Kanban
                title="Sales Opportunities & Bookings"
                description="Drag and drop deals across stages, track values, and manage clients."
                columns={kanbanColumns}
                tasks={kanbanTasks}
                onTasksChange={setKanbanTasks}
                onTaskMove={(taskId: string, targetStatus: string) => {
                  const updated = kanbanTasks.map(t => t.id === taskId ? { ...t, status: targetStatus } : t)
                  setKanbanTasks(updated)
                  return true
                }}
                viewMode={dashboardTab}
                onViewModeChange={(mode: 'table' | 'kanban') => setDashboardTab(mode)}
              />
            </Suspense>
          </CardContent>
        </Card>
      ) : (
        <Card className="border bg-card shadow-sm">
          <CardHeader>
            <CardTitle>Recent Payments & Invoices</CardTitle>
            <CardDescription>
              Production-grade SaaS data table with multi-column search, structured filters, density settings, CSV export, and row drawer.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Suspense fallback={<div className="h-96 flex items-center justify-center text-sm text-muted-foreground">Loading invoice table...</div>}>
              <DataTable
                columns={paymentColumns}
                data={payments}
                searchableColumns={['clientName', 'email', 'package', 'id']}
                searchPlaceholder="Search invoices... (Press ⌘K)"
                filterConfigs={filterConfigs}
                bulkActions={bulkActions}
                rowActions={rowActions}
                renderDetailDrawer={(row, close) => (
                  <div className="space-y-6">
                    <div className="rounded-lg border bg-muted/30 p-4 space-y-3">
                      <div className="flex items-center justify-between">
                        <span className="text-xs font-medium text-muted-foreground">Status</span>
                        <Badge
                          variant={
                            row.status === 'success'
                              ? 'default'
                              : row.status === 'failed'
                              ? 'destructive'
                              : 'secondary'
                          }
                          className="capitalize text-xs"
                        >
                          {row.status}
                        </Badge>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-xs font-medium text-muted-foreground">Invoice Amount</span>
                        <span className="text-lg font-bold tabular-nums">${row.amount.toLocaleString()}</span>
                      </div>
                    </div>

                    <div className="space-y-3">
                      <h4 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                        Client Details
                      </h4>
                      <div className="space-y-2 text-sm">
                        <div className="flex justify-between py-1 border-b">
                          <span className="text-muted-foreground">Client Name</span>
                          <span className="font-medium">{row.clientName}</span>
                        </div>
                        <div className="flex justify-between py-1 border-b">
                          <span className="text-muted-foreground">Email</span>
                          <span className="font-mono text-xs">{row.email}</span>
                        </div>
                        <div className="flex justify-between py-1 border-b">
                          <span className="text-muted-foreground">Package</span>
                          <span>{row.package}</span>
                        </div>
                        <div className="flex justify-between py-1 border-b">
                          <span className="text-muted-foreground">Shoot Date</span>
                          <span>{row.shootDate}</span>
                        </div>
                      </div>
                    </div>

                    <div className="flex gap-2 pt-4">
                      <Button
                        className="flex-1"
                        onClick={() => {
                          toast.success(`Sent receipt to ${row.email}`)
                          close()
                        }}
                      >
                        Send Receipt
                      </Button>
                      <Button variant="outline" onClick={close}>
                        Close
                      </Button>
                    </div>
                  </div>
                )}
                detailDrawerTitle={(row) => `Invoice ${row.id}`}
                detailDrawerDescription={(row) => `Review details and billing status for ${row.clientName}.`}
                enableExport={true}
                exportFilename="crm-invoices.csv"
              />
            </Suspense>
          </CardContent>
        </Card>
      )}
    </DashboardShell>
  )
}
