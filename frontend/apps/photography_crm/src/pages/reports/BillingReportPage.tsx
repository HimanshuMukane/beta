import * as React from 'react'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@repo/ui/components/ui/card'
import { Button } from '@repo/ui/components/ui/button'
import { Badge } from '@repo/ui/components/ui/badge'
import { Input } from '@repo/ui/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@repo/ui/components/ui/select'
import { toast } from '@repo/ui/components/ui/sonner'
import {
  IconReceipt,
  IconDownload,
  IconFileSpreadsheet,
  IconSearch,
  IconTrendingUp,
  IconCash,
  IconClockHour4,
  IconAlertCircle,
  IconPrinter,
} from '@tabler/icons-react'

export function BillingReportPage() {
  const [dateRange, setDateRange] = React.useState('this-year')
  const [statusFilter, setStatusFilter] = React.useState('all')
  const [search, setSearch] = React.useState('')

  const summary = {
    grossBilled: 1845000,
    collected: 1420000,
    outstanding: 375000,
    overdue: 50000,
    invoicesCount: 42,
    collectionRate: 77.0,
  }

  const statementRows = [
    {
      invoiceNo: 'INV-2026-101',
      date: '2026-11-01',
      client: 'Kavya Rajendran',
      project: 'Kavya Rajendran Wedding',
      amount: 54000,
      paid: 30000,
      balance: 24000,
      status: 'Partially Paid',
      dueDate: '2026-11-20',
    },
    {
      invoiceNo: 'INV-2026-102',
      date: '2026-11-02',
      client: 'Arjun Verma',
      project: 'Verma Sangeet Night',
      amount: 85000,
      paid: 85000,
      balance: 0,
      status: 'Paid',
      dueDate: '2026-11-02',
    },
    {
      invoiceNo: 'INV-2026-103',
      date: '2026-11-04',
      client: 'Deepika Rao',
      project: 'Maternity Studio Portrait',
      amount: 32000,
      paid: 32000,
      balance: 0,
      status: 'Paid',
      dueDate: '2026-11-04',
    },
    {
      invoiceNo: 'INV-2026-104',
      date: '2026-10-15',
      client: 'Rohan Kapoor',
      project: 'Pre-Wedding Shoot Goa',
      amount: 110000,
      paid: 60000,
      balance: 50000,
      status: 'Overdue',
      dueDate: '2026-11-01',
    },
    {
      invoiceNo: 'INV-2026-105',
      date: '2026-11-06',
      client: 'Aakash Singhania',
      project: 'Corporate Leadership Headshots',
      amount: 45000,
      paid: 0,
      balance: 45000,
      status: 'Unpaid',
      dueDate: '2026-11-25',
    },
  ]

  const filteredRows = statementRows.filter((r) => {
    if (statusFilter !== 'all' && r.status.toLowerCase() !== statusFilter) return false
    if (search.trim()) {
      const q = search.toLowerCase()
      return (
        r.client.toLowerCase().includes(q) ||
        r.project.toLowerCase().includes(q) ||
        r.invoiceNo.toLowerCase().includes(q)
      )
    }
    return true
  })

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-600 text-white shadow-2xs">
              <IconReceipt className="h-4 w-4" />
            </span>
            <h1 className="text-xl font-bold tracking-tight text-foreground">
              Billing Statement & Invoicing Report
            </h1>
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            Accounts receivable aging, client collection rate, and invoice breakdown.
          </p>
        </div>

        <div className="flex items-center gap-2 flex-wrap">
          <Select value={dateRange} onValueChange={setDateRange}>
            <SelectTrigger className="h-9 text-xs bg-card w-[150px]">
              <SelectValue placeholder="Date range" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="this-month" className="text-xs">This Month</SelectItem>
              <SelectItem value="last-quarter" className="text-xs">Last Quarter</SelectItem>
              <SelectItem value="this-year" className="text-xs">FY 2026-2027</SelectItem>
            </SelectContent>
          </Select>

          <Button
            variant="outline"
            size="sm"
            onClick={() => toast.success('Billing CSV exported')}
            className="h-9 text-xs gap-1.5"
          >
            <IconFileSpreadsheet className="h-3.5 w-3.5 text-emerald-600" />
            <span>Export CSV</span>
          </Button>

          <Button
            size="sm"
            onClick={() => window.print()}
            className="h-9 text-xs gap-1.5 bg-primary text-primary-foreground"
          >
            <IconPrinter className="h-3.5 w-3.5" />
            <span>Print Report</span>
          </Button>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Total Invoiced
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-foreground">
              ₹{summary.grossBilled.toLocaleString('en-IN')}
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-emerald-600 font-semibold">
              <IconTrendingUp className="h-3.5 w-3.5" />
              <span>42 Invoices generated</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Total Collected
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-emerald-600">
              ₹{summary.collected.toLocaleString('en-IN')}
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground font-medium">
              <IconCash className="h-3.5 w-3.5 text-emerald-600" />
              <span>Collection rate: {summary.collectionRate}%</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Outstanding Balance
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-amber-600">
              ₹{summary.outstanding.toLocaleString('en-IN')}
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-amber-600 font-medium">
              <IconClockHour4 className="h-3.5 w-3.5" />
              <span>Pending from client payments</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Overdue Receivables
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-rose-600">
              ₹{summary.overdue.toLocaleString('en-IN')}
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-rose-600 font-medium">
              <IconAlertCircle className="h-3.5 w-3.5" />
              <span>1 invoice past grace period</span>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Filter Bar */}
      <Card className="p-3 shadow-2xs border bg-card">
        <div className="flex flex-col sm:flex-row items-center gap-3">
          <div className="relative flex-1 w-full">
            <IconSearch className="absolute left-2.5 top-2.5 h-3.5 w-3.5 text-muted-foreground" />
            <Input
              placeholder="Search by invoice no, client name, or project..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-8 h-9 text-xs"
            />
          </div>
          <div className="w-full sm:w-[180px]">
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="h-9 text-xs">
                <SelectValue placeholder="All Statuses" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all" className="text-xs">All Statuses</SelectItem>
                <SelectItem value="paid" className="text-xs">Paid</SelectItem>
                <SelectItem value="partially paid" className="text-xs">Partially Paid</SelectItem>
                <SelectItem value="unpaid" className="text-xs">Unpaid</SelectItem>
                <SelectItem value="overdue" className="text-xs">Overdue</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </Card>

      {/* Statement Table */}
      <Card className="border bg-card shadow-2xs overflow-hidden">
        <CardHeader className="pb-3 border-b">
          <CardTitle className="text-sm font-bold">Billing Records</CardTitle>
          <CardDescription className="text-xs">
            Showing {filteredRows.length} invoice records matching selected filters.
          </CardDescription>
        </CardHeader>
        <div className="overflow-x-auto">
          <table className="w-full text-xs">
            <thead className="bg-muted/40 border-b text-[11px] uppercase tracking-wider text-muted-foreground font-semibold">
              <tr>
                <th className="py-2.5 px-3.5 text-left">Invoice No</th>
                <th className="py-2.5 px-3.5 text-left">Date</th>
                <th className="py-2.5 px-3.5 text-left">Client & Project</th>
                <th className="py-2.5 px-3.5 text-right">Total Amount</th>
                <th className="py-2.5 px-3.5 text-right">Paid</th>
                <th className="py-2.5 px-3.5 text-right">Balance</th>
                <th className="py-2.5 px-3.5 text-center">Status</th>
                <th className="py-2.5 px-3.5 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {filteredRows.map((row) => (
                <tr key={row.invoiceNo} className="hover:bg-muted/20 transition">
                  <td className="py-3 px-3.5 font-mono font-bold text-foreground">
                    {row.invoiceNo}
                  </td>
                  <td className="py-3 px-3.5 text-muted-foreground tabular-nums">
                    {row.date}
                  </td>
                  <td className="py-3 px-3.5">
                    <p className="font-semibold text-foreground">{row.client}</p>
                    <p className="text-[11px] text-muted-foreground">{row.project}</p>
                  </td>
                  <td className="py-3 px-3.5 text-right font-bold tabular-nums">
                    ₹{row.amount.toLocaleString('en-IN')}
                  </td>
                  <td className="py-3 px-3.5 text-right font-bold text-emerald-600 tabular-nums">
                    ₹{row.paid.toLocaleString('en-IN')}
                  </td>
                  <td className="py-3 px-3.5 text-right font-bold tabular-nums text-foreground">
                    ₹{row.balance.toLocaleString('en-IN')}
                  </td>
                  <td className="py-3 px-3.5 text-center">
                    <Badge
                      variant="outline"
                      className={`text-[10px] font-semibold ${
                        row.status === 'Paid'
                          ? 'bg-emerald-50 text-emerald-700 border-emerald-200'
                          : row.status === 'Partially Paid'
                          ? 'bg-amber-50 text-amber-700 border-amber-200'
                          : row.status === 'Overdue'
                          ? 'bg-rose-50 text-rose-700 border-rose-200'
                          : 'bg-muted text-muted-foreground'
                      }`}
                    >
                      {row.status}
                    </Badge>
                  </td>
                  <td className="py-3 px-3.5 text-right">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => toast.success(`Downloaded ${row.invoiceNo}`)}
                      className="h-7 px-2 text-xs gap-1"
                    >
                      <IconDownload className="h-3 w-3" />
                      <span>PDF</span>
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  )
}
