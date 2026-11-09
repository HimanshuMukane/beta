import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
  Button,
  toast,
} from '@repo/ui'
import {
  IconChartPie,
  IconTrendingUp,
  IconReceipt2,
  IconWallet,
  IconFileSpreadsheet,
  IconArrowUpRight,
} from '@tabler/icons-react'

export function FinanceReportPage() {
  const expenseCategories = [
    { category: 'Assistant & Freelance Crew Payouts', amount: 320000, percentage: 38 },
    { category: 'Gear Rental & Drone Permissions', amount: 185000, percentage: 22 },
    { category: 'Travel, Fuel & Logistics', amount: 145000, percentage: 17 },
    { category: 'Client Accommodation & Hospitality', amount: 110000, percentage: 13 },
    { category: 'Studio Rent & Equipment Maintenance', amount: 84000, percentage: 10 },
  ]

  const revenueStreams = [
    { stream: 'Wedding & Multi-Day Package Shoots', gross: 2480000, profit: 1640000, margin: '66.1%' },
    { stream: 'Pre-Wedding & Destination Sessions', gross: 820000, profit: 540000, margin: '65.8%' },
    { stream: 'Maternity, Baby & Family Portraits', gross: 460000, profit: 345000, margin: '75.0%' },
    { stream: 'Corporate Events & Leadership Headshots', gross: 390000, profit: 295000, margin: '75.6%' },
    { stream: 'Print Albums, Canvas & USB Merch', gross: 310000, profit: 185000, margin: '59.6%' },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-600 text-white shadow-2xs">
              <IconChartPie className="h-4 w-4" />
            </span>
            <h1 className="text-xl font-bold tracking-tight text-foreground">
              Financial P&L & Profitability Statement
            </h1>
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            Studio net margins, production cost breakdown, and event profitability analysis.
          </p>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => toast.success('P&L report exported')}
            className="h-9 text-xs gap-1.5"
          >
            <IconFileSpreadsheet className="h-3.5 w-3.5 text-emerald-600" />
            <span>Export Financial P&L</span>
          </Button>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Total Studio Gross Revenue
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-foreground">
              ₹44,60,000
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-emerald-600 font-semibold">
              <IconArrowUpRight className="h-3.5 w-3.5" />
              <span>+24.8% YoY growth</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Total Production Expense
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-rose-600">
              ₹8,44,000
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground font-medium">
              <IconReceipt2 className="h-3.5 w-3.5 text-rose-600" />
              <span>18.9% of gross revenue</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Net Profit Before Tax
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-emerald-600">
              ₹36,16,000
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-emerald-600 font-semibold">
              <IconTrendingUp className="h-3.5 w-3.5" />
              <span>Net profit margin: 81.1%</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Liquid Operating Cashflow
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-primary">
              ₹18,50,000
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground font-medium">
              <IconWallet className="h-3.5 w-3.5 text-primary" />
              <span>In studio current accounts</span>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Grid: Profit Margin by Stream + Operating Cost Breakdown */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <Card className="lg:col-span-8 border bg-card shadow-2xs">
          <CardHeader className="pb-3 border-b">
            <CardTitle className="text-sm font-bold">Revenue Stream Profitability</CardTitle>
            <CardDescription className="text-xs">Gross receipts vs direct costs by package type.</CardDescription>
          </CardHeader>
          <div className="overflow-x-auto">
            <table className="w-full text-xs">
              <thead className="bg-muted/40 border-b text-[11px] uppercase tracking-wider text-muted-foreground font-semibold">
                <tr>
                  <th className="py-2.5 px-3.5 text-left">Revenue Stream</th>
                  <th className="py-2.5 px-3.5 text-right">Gross Invoiced</th>
                  <th className="py-2.5 px-3.5 text-right">Net Profit</th>
                  <th className="py-2.5 px-3.5 text-center">Profit Margin</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {revenueStreams.map((rs) => (
                  <tr key={rs.stream} className="hover:bg-muted/20 transition">
                    <td className="py-3 px-3.5 font-semibold text-foreground">{rs.stream}</td>
                    <td className="py-3 px-3.5 text-right tabular-nums font-medium text-foreground">
                      ₹{rs.gross.toLocaleString('en-IN')}
                    </td>
                    <td className="py-3 px-3.5 text-right font-bold tabular-nums text-emerald-600">
                      ₹{rs.profit.toLocaleString('en-IN')}
                    </td>
                    <td className="py-3 px-3.5 text-center font-bold text-foreground tabular-nums">
                      {rs.margin}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Card>

        <Card className="lg:col-span-4 border bg-card shadow-2xs">
          <CardHeader className="pb-3 border-b">
            <CardTitle className="text-sm font-bold">Production Cost Breakdown</CardTitle>
            <CardDescription className="text-xs">Where shoot expenses are spent.</CardDescription>
          </CardHeader>
          <CardContent className="pt-4 space-y-3.5">
            {expenseCategories.map((ec) => (
              <div key={ec.category} className="space-y-1">
                <div className="flex items-center justify-between text-xs">
                  <span className="font-semibold text-foreground truncate pr-2">{ec.category}</span>
                  <span className="font-bold tabular-nums text-foreground">₹{ec.amount.toLocaleString('en-IN')}</span>
                </div>
                <div className="h-2 w-full rounded-full bg-muted overflow-hidden flex">
                  <div
                    className="h-full bg-rose-500 rounded-full transition-all"
                    style={{ width: `${ec.percentage}%` }}
                  />
                </div>
              </div>
            ))}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
