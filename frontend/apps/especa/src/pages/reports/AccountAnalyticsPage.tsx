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
  IconUserCheck,
  IconFileSpreadsheet,
  IconStarFilled,
} from "@repo/ui/icons/tabler-icons"

export function AccountAnalyticsPage() {
  const topAccounts = [
    { client: 'Kavya Rajendran', totalSpend: 139000, projectsCount: 2, ltvTier: 'VIP Platinum', loyaltyPoints: 1200 },
    { client: 'Rohan Kapoor', totalSpend: 110000, projectsCount: 1, ltvTier: 'Gold', loyaltyPoints: 950 },
    { client: 'Arjun Verma', totalSpend: 85000, projectsCount: 1, ltvTier: 'Gold', loyaltyPoints: 750 },
    { client: 'Aakash Singhania (Titan)', totalSpend: 145000, projectsCount: 3, ltvTier: 'Enterprise VIP', loyaltyPoints: 2100 },
    { client: 'Deepika Rao', totalSpend: 62000, projectsCount: 2, ltvTier: 'Silver', loyaltyPoints: 500 },
  ]

  const retentionMetrics = [
    { metric: 'Repeat Client Ratio', value: '31.4%', desc: 'Booked subsequent shoots (Maternity, Anniversary, Baby)' },
    { metric: 'Client Referral Rate', value: '28.2%', desc: 'Brought friends or relatives through referral links' },
    { metric: 'Average Lifetime Value (LTV)', value: '₹92,800', desc: 'Net billings per retained client account' },
    { metric: 'Average Sales Cycle', value: '14 Days', desc: 'From first phone inquiry to formal agreement sign-off' },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-600 text-white shadow-2xs">
              <IconUserCheck className="h-4 w-4" />
            </span>
            <h1 className="text-xl font-bold tracking-tight text-foreground">
              Client Accounts Analytics & Lifetime Value (LTV)
            </h1>
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            Retention cohorts, recurring client loyalty, and high-net-worth account metrics.
          </p>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => toast.success('Account analytics exported')}
            className="h-9 text-xs gap-1.5"
          >
            <IconFileSpreadsheet className="h-3.5 w-3.5 text-blue-600" />
            <span>Export Account Analytics</span>
          </Button>
        </div>
      </div>

      {/* Retention KPI Blocks */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {retentionMetrics.map((rm) => (
          <Card key={rm.metric} className="border bg-card shadow-2xs">
            <CardHeader className="pb-2">
              <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                {rm.metric}
              </CardDescription>
              <CardTitle className="text-2xl font-black tabular-nums text-foreground">
                {rm.value}
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-0">
              <p className="text-xs text-muted-foreground">{rm.desc}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Top Client Accounts Leaderboard */}
      <Card className="border bg-card shadow-2xs overflow-hidden">
        <CardHeader className="pb-3 border-b">
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="text-sm font-bold">Top High-LTV Client Accounts</CardTitle>
              <CardDescription className="text-xs">Highest cumulative billing clients across multiple sessions.</CardDescription>
            </div>
            <Badge variant="outline" className="text-xs font-semibold">
              VIP Tier Active
            </Badge>
          </div>
        </CardHeader>
        <div className="overflow-x-auto">
          <table className="w-full text-xs">
            <thead className="bg-muted/40 border-b text-[11px] uppercase tracking-wider text-muted-foreground font-semibold">
              <tr>
                <th className="py-2.5 px-3.5 text-left">Client Account</th>
                <th className="py-2.5 px-3.5 text-center">Shoots Completed</th>
                <th className="py-2.5 px-3.5 text-right">Lifetime Spend</th>
                <th className="py-2.5 px-3.5 text-center">Account Tier</th>
                <th className="py-2.5 px-3.5 text-right">Reward Points</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {topAccounts.map((acc) => (
                <tr key={acc.client} className="hover:bg-muted/20 transition">
                  <td className="py-3 px-3.5">
                    <p className="font-bold text-foreground">{acc.client}</p>
                  </td>
                  <td className="py-3 px-3.5 text-center font-bold tabular-nums">{acc.projectsCount}</td>
                  <td className="py-3 px-3.5 text-right font-black tabular-nums text-foreground">
                    ₹{acc.totalSpend.toLocaleString('en-IN')}
                  </td>
                  <td className="py-3 px-3.5 text-center">
                    <Badge
                      variant="secondary"
                      className="inline-flex items-center gap-1 text-[10px] font-semibold bg-amber-500/15 text-amber-700 dark:text-amber-400"
                    >
                      <IconStarFilled className="h-2.5 w-2.5 text-amber-500" />
                      <span>{acc.ltvTier}</span>
                    </Badge>
                  </td>
                  <td className="py-3 px-3.5 text-right font-mono tabular-nums text-muted-foreground">
                    {acc.loyaltyPoints} pts
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
