import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@repo/ui/components/ui/card'
import { Button } from '@repo/ui/components/ui/button'
import { Badge } from '@repo/ui/components/ui/badge'
import { toast } from '@repo/ui/components/ui/sonner'
import {
  IconCalendarEvent,
  IconClock,
  IconUsers,
  IconFileSpreadsheet,
  IconCheck,
  IconStarFilled,
} from '@tabler/icons-react'

export function SessionsReportPage() {
  const shootTypeBreakdown = [
    { type: 'Wedding', sessions: 28, totalHours: 196, crewAllocations: 84 },
    { type: 'Pre-Wedding', sessions: 16, totalHours: 96, crewAllocations: 32 },
    { type: 'Engagement', sessions: 14, totalHours: 56, crewAllocations: 28 },
    { type: 'Maternity & Newborn', sessions: 22, totalHours: 66, crewAllocations: 22 },
    { type: 'Corporate Headshots', sessions: 12, totalHours: 48, crewAllocations: 18 },
    { type: 'Birthday & Event', sessions: 18, totalHours: 72, crewAllocations: 36 },
  ]

  const crewUtilization = [
    { name: 'Rohan Mehta', role: 'Lead Photographer', shootsCompleted: 24, hours: 148, rating: '4.9' },
    { name: 'Aditya Roy', role: 'Cinematographer', shootsCompleted: 20, hours: 132, rating: '4.9' },
    { name: 'Priya Sharma', role: 'Candid Photographer', shootsCompleted: 18, hours: 110, rating: '4.8' },
    { name: 'Vikram Seth', role: 'Drone Operator', shootsCompleted: 14, hours: 56, rating: '5.0' },
    { name: 'Neha Gupta', role: 'Shoot Assistant', shootsCompleted: 22, hours: 120, rating: '4.8' },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-600 text-white shadow-2xs">
              <IconCalendarEvent className="h-4 w-4" />
            </span>
            <h1 className="text-xl font-bold tracking-tight text-foreground">
              Sessions & Crew Utilization Report
            </h1>
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            Shoot volume timeline, camera hours logged, and crew member deployments.
          </p>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => toast.success('Sessions schedule report exported')}
            className="h-9 text-xs gap-1.5"
          >
            <IconFileSpreadsheet className="h-3.5 w-3.5 text-emerald-600" />
            <span>Export Roster CSV</span>
          </Button>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Total Sessions Executed
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-foreground">
              110
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-emerald-600 font-semibold">
              <IconCheck className="h-3.5 w-3.5" />
              <span>100% on-schedule completion</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Total Shoot Hours
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-emerald-600">
              534.0 Hrs
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground font-medium">
              <IconClock className="h-3.5 w-3.5 text-emerald-600" />
              <span>Avg shoot: 4.85 hours</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Active Crew Deployed
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-primary">
              14 Members
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground font-medium">
              <IconUsers className="h-3.5 w-3.5 text-primary" />
              <span>Across lead, candid & drone</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Upcoming Booked Shoots
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-amber-600">
              18 Sessions
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-amber-600 font-medium">
              <IconCalendarEvent className="h-3.5 w-3.5" />
              <span>Next 30 days schedule</span>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Grid: Crew Utilization Roster + Project Type Distribution */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <Card className="lg:col-span-7 border bg-card shadow-2xs">
          <CardHeader className="pb-3 border-b">
            <CardTitle className="text-sm font-bold">Crew Deployment & Workload</CardTitle>
            <CardDescription className="text-xs">Individual photographer field hours and performance.</CardDescription>
          </CardHeader>
          <div className="overflow-x-auto">
            <table className="w-full text-xs">
              <thead className="bg-muted/40 border-b text-[11px] uppercase tracking-wider text-muted-foreground font-semibold">
                <tr>
                  <th className="py-2.5 px-3.5 text-left">Crew Member</th>
                  <th className="py-2.5 px-3.5 text-center">Shoots</th>
                  <th className="py-2.5 px-3.5 text-right">Hours Logged</th>
                  <th className="py-2.5 px-3.5 text-right">Client Rating</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {crewUtilization.map((c) => (
                  <tr key={c.name} className="hover:bg-muted/20 transition">
                    <td className="py-3 px-3.5">
                      <p className="font-semibold text-foreground">{c.name}</p>
                      <p className="text-[11px] text-muted-foreground">{c.role}</p>
                    </td>
                    <td className="py-3 px-3.5 text-center font-bold tabular-nums">{c.shootsCompleted}</td>
                    <td className="py-3 px-3.5 text-right font-bold tabular-nums text-foreground">{c.hours} hrs</td>
                    <td className="py-3 px-3.5 text-right font-bold text-amber-600 tabular-nums">
                      <span className="inline-flex items-center gap-1 justify-end">
                        <IconStarFilled className="h-3 w-3 text-amber-500" />
                        <span>{c.rating}</span>
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Card>

        <Card className="lg:col-span-5 border bg-card shadow-2xs">
          <CardHeader className="pb-3 border-b">
            <CardTitle className="text-sm font-bold">Project Genre Hours</CardTitle>
            <CardDescription className="text-xs">Hours spent on site by event type.</CardDescription>
          </CardHeader>
          <CardContent className="pt-4 space-y-3">
            {shootTypeBreakdown.map((st) => (
              <div key={st.type} className="flex items-center justify-between p-2.5 rounded-lg border bg-muted/20 text-xs">
                <div>
                  <p className="font-semibold text-foreground">{st.type}</p>
                  <p className="text-[11px] text-muted-foreground">{st.sessions} Sessions • {st.crewAllocations} Crew slots</p>
                </div>
                <Badge variant="secondary" className="font-bold tabular-nums text-xs">
                  {st.totalHours} hrs
                </Badge>
              </div>
            ))}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
