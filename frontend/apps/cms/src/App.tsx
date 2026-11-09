import { useState, Suspense, lazy } from 'react'
import {
  DashboardShell,
  SplitAuthLayout,
  DataTable,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
  Button,
  Badge,
  Checkbox,
} from '@repo/ui'
import { createColumnHelper } from '@tanstack/react-table'
import type { DataTableFeatures } from '@repo/ui'
import { IconArrowsSort } from '@tabler/icons-react'

const LoginForm = lazy(() => import('@repo/ui/components/login-form').then(m => ({ default: m.LoginForm })))
const SignupForm = lazy(() => import('@repo/ui/components/signup-form').then(m => ({ default: m.SignupForm })))

// Example data model
type PaymentRecord = {
  id: string
  amount: number
  status: 'pending' | 'processing' | 'success' | 'failed'
  email: string
}

const samplePayments: PaymentRecord[] = [
  { id: '728ed52f', amount: 100, status: 'pending', email: 'm@example.com' },
  { id: '489e1d42', amount: 125, status: 'processing', email: 'client@gmail.com' },
  { id: '6a91c01b', amount: 350, status: 'success', email: 'studio@domain.com' },
  { id: '9f23d84a', amount: 75, status: 'failed', email: 'dev@agency.io' },
  { id: '2c84a51e', amount: 450, status: 'success', email: 'sarah@photography.com' },
  { id: '5d91e62f', amount: 200, status: 'pending', email: 'alex@creative.org' },
]

const columnHelper = createColumnHelper<DataTableFeatures, PaymentRecord>()

const paymentColumns = columnHelper.columns([
  columnHelper.display({
    id: 'select',
    header: ({ table }) => (
      <Checkbox
        checked={
          table.getIsAllPageRowsSelected() ||
          (table.getIsSomePageRowsSelected() && 'indeterminate')
        }
        onCheckedChange={(value: boolean | 'indeterminate') => table.toggleAllPageRowsSelected(!!value)}
        aria-label="Select all"
      />
    ),
    cell: ({ row }) => (
      <Checkbox
        checked={row.getIsSelected()}
        onCheckedChange={(value: boolean | 'indeterminate') => row.toggleSelected(!!value)}
        aria-label="Select row"
      />
    ),
    enableSorting: false,
    enableHiding: false,
  }),
  columnHelper.accessor('status', {
    header: 'Status',
    cell: ({ row }) => {
      const status = row.getValue('status') as string
      const variant =
        status === 'success'
          ? 'default'
          : status === 'processing'
          ? 'secondary'
          : status === 'failed'
          ? 'destructive'
          : 'outline'
      return <Badge variant={variant} className="capitalize">{status}</Badge>
    },
  }),
  columnHelper.accessor('email', {
    header: ({ column }) => (
      <Button
        variant="ghost"
        className="-ml-3 h-8"
        onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
      >
        Email
        <IconArrowsSort className="ml-2 h-4 w-4" />
      </Button>
    ),
  }),
  columnHelper.accessor('amount', {
    header: () => <div className="text-right">Amount</div>,
    cell: ({ row }) => {
      const amount = parseFloat(row.getValue('amount'))
      const formatted = new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD',
      }).format(amount)
      return <div className="text-right font-medium">{formatted}</div>
    },
  }),
])

export default function App() {
  const [currentView, setCurrentView] = useState<'dashboard' | 'signin' | 'signup'>(() => {
    const path = window.location.pathname
    if (path === '/signin' || path === '/login') return 'signin'
    if (path === '/signup' || path === '/register') return 'signup'
    return 'dashboard'
  })

  if (currentView === 'signin') {
    return (
      <SplitAuthLayout appName="CMS Portal">
        <div className="flex flex-col gap-6">
          <div className="flex items-center justify-between">
            <span className="text-lg font-bold tracking-tight">CMS Portal</span>
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
      <SplitAuthLayout appName="CMS Portal">
        <div className="flex flex-col gap-6">
          <div className="flex items-center justify-between">
            <span className="text-lg font-bold tracking-tight">CMS Portal</span>
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
    <DashboardShell title="CMS Portal" breadcrumb="Overview">
      <div className="flex items-center justify-between mb-2">
        <div>
          <h1 className="text-lg font-bold tracking-tight">CMS Portal</h1>
          <p className="text-sm text-muted-foreground">
            Content management, blogs and media assets.
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm" onClick={() => setCurrentView('signin')}>
            Switch to Sign In
          </Button>
          <Button variant="outline" size="sm" onClick={() => setCurrentView('signup')}>
            Switch to Sign Up
          </Button>
        </div>
      </div>

      {/* KPI Stats */}
      <div className="grid gap-2.5 md:grid-cols-3 mb-3">
        <Card className="border bg-card shadow-sm">
          <CardHeader>
            <CardTitle>Total Items</CardTitle>
            <CardDescription>Records in workspace</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="text-xl font-bold">24</div>
            <p className="text-xs text-muted-foreground mt-1">Synchronized with API</p>
          </CardContent>
        </Card>
        <Card className="border bg-card shadow-sm">
          <CardHeader>
            <CardTitle>System Status</CardTitle>
            <CardDescription>Backend connection</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="text-xl font-bold text-green-600">Connected</div>
            <p className="text-xs text-muted-foreground mt-1">Ready for requests</p>
          </CardContent>
        </Card>
        <Card className="border bg-card shadow-sm">
          <CardHeader>
            <CardTitle>Activity</CardTitle>
            <CardDescription>Recent operations</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="text-xl font-bold">12</div>
            <p className="text-xs text-muted-foreground mt-1">Actions today</p>
          </CardContent>
        </Card>
      </div>

      {/* TanStack Table v9 DataTable Integration */}
      <Card className="border bg-card shadow-sm">
        <CardHeader>
          <CardTitle>Recent Payments & Invoices</CardTitle>
          <CardDescription>
            Interactive DataTable with column sorting, filtering, row selection, and visibility controls.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <DataTable
            columns={paymentColumns}
            data={samplePayments}
            searchKey="email"
            searchPlaceholder="Filter emails..."
          />
        </CardContent>
      </Card>
    </DashboardShell>
  )
}
