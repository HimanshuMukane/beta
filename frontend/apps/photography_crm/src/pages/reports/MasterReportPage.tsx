import * as React from 'react'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
  Button,
  Badge,
  Input,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  toast,
} from '@repo/ui'
import {
  IconReportAnalytics,
  IconFileSpreadsheet,
  IconSearch,
  IconPrinter,
  IconTrendingUp,
  IconBriefcase,
  IconCash,
} from '@tabler/icons-react'

export function MasterReportPage() {
  const [searchTerm, setSearchTerm] = React.useState('')
  const [statusFilter, setStatusFilter] = React.useState('all')

  const masterRows = [
    {
      id: 'PRJ-101',
      project: 'Kavya Rajendran Wedding',
      client: 'Kavya Rajendran',
      type: 'Wedding',
      sessions: 2,
      leadsSource: 'Instagram',
      manager: 'Rohan Mehta',
      quoteAmount: 54000,
      paidAmount: 30000,
      balance: 24000,
      stage: 'Booked',
      startDate: '2026-11-11',
    },
    {
      id: 'PRJ-102',
      project: 'Arjun & Sneha Sangeet Night',
      client: 'Arjun Verma',
      type: 'Wedding',
      sessions: 1,
      leadsSource: 'Referral',
      manager: 'Priya Sharma',
      quoteAmount: 85000,
      paidAmount: 85000,
      balance: 0,
      stage: 'Shoot Completed',
      startDate: '2026-11-02',
    },
    {
      id: 'PRJ-103',
      project: 'Deepika Rao Maternity Session',
      client: 'Deepika Rao',
      type: 'Maternity',
      sessions: 1,
      leadsSource: 'Website',
      manager: 'Rohan Mehta',
      quoteAmount: 32000,
      paidAmount: 32000,
      balance: 0,
      stage: 'Editing',
      startDate: '2026-11-04',
    },
    {
      id: 'PRJ-104',
      project: 'Goa Destination Pre-Wedding',
      client: 'Rohan Kapoor',
      type: 'Pre-Wedding',
      sessions: 3,
      leadsSource: 'Instagram',
      manager: 'Aditya Roy',
      quoteAmount: 110000,
      paidAmount: 60000,
      balance: 50000,
      stage: 'Delivery',
      startDate: '2026-10-15',
    },
    {
      id: 'PRJ-105',
      project: 'Titan Leadership Corporate Shoot',
      client: 'Aakash Singhania',
      type: 'Corporate',
      sessions: 1,
      leadsSource: 'Venue Partner',
      manager: 'Rohan Mehta',
      quoteAmount: 45000,
      paidAmount: 0,
      balance: 45000,
      stage: 'Shoot Planning',
      startDate: '2026-11-20',
    },
  ]

  const filtered = masterRows.filter((row) => {
    if (statusFilter !== 'all' && row.stage.toLowerCase() !== statusFilter) return false
    if (searchTerm.trim()) {
      const q = searchTerm.toLowerCase()
      return (
        row.project.toLowerCase().includes(q) ||
        row.client.toLowerCase().includes(q) ||
        row.id.toLowerCase().includes(q) ||
        row.manager.toLowerCase().includes(q)
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
              <IconReportAnalytics className="h-4 w-4" />
            </span>
            <h1 className="text-xl font-bold tracking-tight text-foreground">
              Master Studio Operations & Project Ledger
            </h1>
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            Consolidated 360-degree ledger combining project dates, sessions, revenue, balance, and leads.
          </p>
        </div>

        <div className="flex items-center gap-2 flex-wrap">
          <Button
            variant="outline"
            size="sm"
            onClick={() => toast.success('Master ledger exported to CSV')}
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
            <span>Print Master</span>
          </Button>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Total Projects Logged
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-foreground">
              48 Projects
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-emerald-600 font-semibold">
              <IconBriefcase className="h-3.5 w-3.5" />
              <span>Across 6 shoot genres</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Gross Contracted Value
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-emerald-600">
              ₹44,60,000
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground font-medium">
              <IconCash className="h-3.5 w-3.5 text-emerald-600" />
              <span>Invoiced & agreed quotes</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Realized Cash Collected
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-primary">
              ₹32,70,000
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-primary font-medium">
              <IconTrendingUp className="h-3.5 w-3.5" />
              <span>73.3% realized advances</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Unsettled Balances
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-amber-600">
              ₹11,90,000
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-amber-600 font-medium">
              <span>Payable upon final photo delivery</span>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Filter and Search Bar */}
      <Card className="p-3 shadow-2xs border bg-card">
        <div className="flex flex-col sm:flex-row items-center gap-3">
          <div className="relative flex-1 w-full">
            <IconSearch className="absolute left-2.5 top-2.5 h-3.5 w-3.5 text-muted-foreground" />
            <Input
              placeholder="Search across project name, client, manager, or code..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-8 h-9 text-xs"
            />
          </div>
          <div className="w-full sm:w-[200px]">
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="h-9 text-xs">
                <SelectValue placeholder="All Stages" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all" className="text-xs">All Stages</SelectItem>
                <SelectItem value="booked" className="text-xs">Booked</SelectItem>
                <SelectItem value="shoot planning" className="text-xs">Shoot Planning</SelectItem>
                <SelectItem value="shoot completed" className="text-xs">Shoot Completed</SelectItem>
                <SelectItem value="editing" className="text-xs">Editing</SelectItem>
                <SelectItem value="delivery" className="text-xs">Delivery</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </Card>

      {/* Master Operations Table */}
      <Card className="border bg-card shadow-2xs overflow-hidden">
        <CardHeader className="pb-3 border-b">
          <CardTitle className="text-sm font-bold">Consolidated Operations Register</CardTitle>
          <CardDescription className="text-xs">
            Showing {filtered.length} matched records with complete financials.
          </CardDescription>
        </CardHeader>
        <div className="overflow-x-auto">
          <table className="w-full text-xs">
            <thead className="bg-muted/40 border-b text-[11px] uppercase tracking-wider text-muted-foreground font-semibold">
              <tr>
                <th className="py-2.5 px-3.5 text-left">Project & Client</th>
                <th className="py-2.5 px-3.5 text-left">Type</th>
                <th className="py-2.5 px-3.5 text-center">Sessions</th>
                <th className="py-2.5 px-3.5 text-left">Manager</th>
                <th className="py-2.5 px-3.5 text-right">Contract Value</th>
                <th className="py-2.5 px-3.5 text-right">Paid</th>
                <th className="py-2.5 px-3.5 text-right">Balance</th>
                <th className="py-2.5 px-3.5 text-center">Production Stage</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {filtered.map((row) => (
                <tr key={row.id} className="hover:bg-muted/20 transition">
                  <td className="py-3 px-3.5">
                    <p className="font-bold text-foreground">{row.project}</p>
                    <p className="text-[11px] text-muted-foreground">{row.client} • {row.startDate}</p>
                  </td>
                  <td className="py-3 px-3.5">
                    <Badge variant="secondary" className="text-[10px]">{row.type}</Badge>
                  </td>
                  <td className="py-3 px-3.5 text-center font-bold tabular-nums">{row.sessions}</td>
                  <td className="py-3 px-3.5 font-medium text-foreground">{row.manager}</td>
                  <td className="py-3 px-3.5 text-right font-bold tabular-nums">
                    ₹{row.quoteAmount.toLocaleString('en-IN')}
                  </td>
                  <td className="py-3 px-3.5 text-right font-bold text-emerald-600 tabular-nums">
                    ₹{row.paidAmount.toLocaleString('en-IN')}
                  </td>
                  <td className="py-3 px-3.5 text-right font-bold tabular-nums text-foreground">
                    ₹{row.balance.toLocaleString('en-IN')}
                  </td>
                  <td className="py-3 px-3.5 text-center">
                    <Badge variant="outline" className="text-[10px] font-semibold bg-muted/30">
                      {row.stage}
                    </Badge>
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
