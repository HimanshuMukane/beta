import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
  Button,
  Badge,
  toast,
} from '@repo/ui'
import {
  IconTargetArrow,
  IconChartBar,
  IconFileSpreadsheet,
  IconArrowUpRight,
  IconClock,
  IconPercentage,
} from '@tabler/icons-react'

export function LeadsReportPage() {

  const leadSources = [
    { source: 'Instagram Ads & DMs', count: 48, percentage: 42, converted: 18, revenue: 980000 },
    { source: 'Client Referrals / Word of Mouth', count: 32, percentage: 28, converted: 22, revenue: 1450000 },
    { source: 'Website Inquiry Form', count: 20, percentage: 17, converted: 9, revenue: 520000 },
    { source: 'Venue Partners / Wedding Planners', count: 15, percentage: 13, converted: 8, revenue: 640000 },
  ]

  const stageBreakdown = [
    { stage: 'New Inquiries', count: 34, value: '₹14.2L' },
    { stage: 'Contacted / Discovery Call', count: 26, value: '₹11.8L' },
    { stage: 'Proposal / Quote Sent', count: 18, value: '₹9.6L' },
    { stage: 'Negotiation / Customizing', count: 12, value: '₹6.8L' },
    { stage: 'Booked / Advance Paid', count: 57, value: '₹35.9L' },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-600 text-white shadow-2xs">
              <IconTargetArrow className="h-4 w-4" />
            </span>
            <h1 className="text-xl font-bold tracking-tight text-foreground">
              Leads Pipeline & Acquisition Report
            </h1>
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            Conversion funnel analysis, lead attribution channels, and deal velocity metrics.
          </p>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => toast.success('Leads report exported')}
            className="h-9 text-xs gap-1.5"
          >
            <IconFileSpreadsheet className="h-3.5 w-3.5 text-blue-600" />
            <span>Export Analytics</span>
          </Button>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Total Inquiries
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-foreground">
              115
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-emerald-600 font-semibold">
              <IconArrowUpRight className="h-3.5 w-3.5" />
              <span>+18.4% from last month</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Converted Bookings
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-emerald-600">
              57
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground font-medium">
              <IconPercentage className="h-3.5 w-3.5 text-emerald-600" />
              <span>Overall conversion: 49.5%</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Pipeline Value
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-primary">
              ₹42,40,000
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground font-medium">
              <IconChartBar className="h-3.5 w-3.5 text-primary" />
              <span>Average deal size: ₹74,380</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Avg Turnaround Time
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-foreground">
              3.2 Days
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-emerald-600 font-medium">
              <IconClock className="h-3.5 w-3.5" />
              <span>From Inquiry to First Quote</span>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Grid: Acquisition Source Performance + Stage Velocity */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Acquisition Attribution */}
        <Card className="lg:col-span-7 border bg-card shadow-2xs">
          <CardHeader className="pb-3 border-b">
            <CardTitle className="text-sm font-bold">Acquisition Channel Performance</CardTitle>
            <CardDescription className="text-xs">
              Where high-ticket photography clients are discovering the studio.
            </CardDescription>
          </CardHeader>
          <CardContent className="pt-4 space-y-4">
            {leadSources.map((src) => (
              <div key={src.source} className="space-y-1.5">
                <div className="flex items-center justify-between text-xs">
                  <span className="font-semibold text-foreground">{src.source}</span>
                  <div className="flex items-center gap-3">
                    <span className="text-muted-foreground tabular-nums">
                      {src.converted} booked ({Math.round((src.converted / src.count) * 100)}%)
                    </span>
                    <span className="font-bold text-foreground tabular-nums">
                      ₹{src.revenue.toLocaleString('en-IN')}
                    </span>
                  </div>
                </div>
                {/* Visual bar */}
                <div className="h-2 w-full rounded-full bg-muted overflow-hidden flex">
                  <div
                    className="h-full bg-blue-600 rounded-full transition-all"
                    style={{ width: `${src.percentage}%` }}
                  />
                </div>
              </div>
            ))}
          </CardContent>
        </Card>

        {/* Funnel Stage Breakdown */}
        <Card className="lg:col-span-5 border bg-card shadow-2xs">
          <CardHeader className="pb-3 border-b">
            <CardTitle className="text-sm font-bold">Pipeline Stage Volume</CardTitle>
            <CardDescription className="text-xs">Active deals moving through sales stages.</CardDescription>
          </CardHeader>
          <CardContent className="pt-4 space-y-3">
            {stageBreakdown.map((stg) => (
              <div
                key={stg.stage}
                className="flex items-center justify-between p-2.5 rounded-lg border bg-muted/20 text-xs"
              >
                <div>
                  <p className="font-semibold text-foreground">{stg.stage}</p>
                  <p className="text-[11px] text-muted-foreground">Est. Pipeline: {stg.value}</p>
                </div>
                <Badge variant="secondary" className="font-bold tabular-nums text-xs">
                  {stg.count} Leads
                </Badge>
              </div>
            ))}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
