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
  IconShoppingBag,
  IconPackage,
  IconTruckDelivery,
  IconFileSpreadsheet,
  IconClock,
  IconCurrencyRupee,
} from "@repo/ui/icons/tabler-icons"

export function OrdersReportPage() {

  const orderTypes = [
    { type: 'Walk-in & Passport', count: 142, revenue: 142000, avgTurnaround: '45 mins' },
    { type: 'Print & Album Orders', count: 86, revenue: 473000, avgTurnaround: '5 days' },
    { type: 'Framing & Canvas Prints', count: 54, revenue: 216000, avgTurnaround: '3 days' },
    { type: 'Commercial / Product Shoots', count: 18, revenue: 396000, avgTurnaround: '7 days' },
    { type: 'Custom Merch & Photo Gifts', count: 35, revenue: 87500, avgTurnaround: '2 days' },
  ]

  const fulfillmentStages = [
    { stage: 'Order Confirmation', count: 12, color: 'bg-blue-500' },
    { stage: 'Product Ready / In Lab', count: 18, color: 'bg-amber-500' },
    { stage: 'Ready For Delivery', count: 15, color: 'bg-purple-500' },
    { stage: 'Product Delivered', count: 290, color: 'bg-emerald-500' },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-purple-600 text-white shadow-2xs">
              <IconShoppingBag className="h-4 w-4" />
            </span>
            <h1 className="text-xl font-bold tracking-tight text-foreground">
              Orders & Merchandise Fulfillment Report
            </h1>
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            Studio print sales, lab dispatch lifecycle, and product margin analysis.
          </p>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => toast.success('Orders report exported')}
            className="h-9 text-xs gap-1.5"
          >
            <IconFileSpreadsheet className="h-3.5 w-3.5 text-purple-600" />
            <span>Export Orders Report</span>
          </Button>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Total Orders Fulfilled
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-foreground">
              335
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-emerald-600 font-semibold">
              <IconPackage className="h-3.5 w-3.5" />
              <span>96.2% on-time client pickup</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Merchandise Revenue
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-purple-600">
              ₹13,14,500
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground font-medium">
              <IconCurrencyRupee className="h-3.5 w-3.5 text-purple-600" />
              <span>Avg order value: ₹3,920</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              In-Production Queue
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-amber-600">
              45
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-amber-600 font-medium">
              <IconClock className="h-3.5 w-3.5" />
              <span>Active in print lab or framing</span>
            </div>
          </CardContent>
        </Card>

        <Card className="border bg-card shadow-2xs">
          <CardHeader className="pb-2">
            <CardDescription className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Pending Deliveries
            </CardDescription>
            <CardTitle className="text-2xl font-black tabular-nums text-blue-600">
              15
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <div className="flex items-center gap-1.5 text-xs text-blue-600 font-medium">
              <IconTruckDelivery className="h-3.5 w-3.5" />
              <span>Ready for client collection</span>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Breakdown tables */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <Card className="lg:col-span-8 border bg-card shadow-2xs">
          <CardHeader className="pb-3 border-b">
            <CardTitle className="text-sm font-bold">Product Category Performance</CardTitle>
            <CardDescription className="text-xs">Sales distribution across photography print products.</CardDescription>
          </CardHeader>
          <div className="overflow-x-auto">
            <table className="w-full text-xs">
              <thead className="bg-muted/40 border-b text-[11px] uppercase tracking-wider text-muted-foreground font-semibold">
                <tr>
                  <th className="py-2.5 px-3.5 text-left">Category / Type</th>
                  <th className="py-2.5 px-3.5 text-right">Units Sold</th>
                  <th className="py-2.5 px-3.5 text-right">Gross Revenue</th>
                  <th className="py-2.5 px-3.5 text-center">Avg Turnaround</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {orderTypes.map((ot) => (
                  <tr key={ot.type} className="hover:bg-muted/20 transition">
                    <td className="py-3 px-3.5 font-semibold text-foreground">{ot.type}</td>
                    <td className="py-3 px-3.5 text-right tabular-nums">{ot.count}</td>
                    <td className="py-3 px-3.5 text-right font-bold tabular-nums text-foreground">
                      ₹{ot.revenue.toLocaleString('en-IN')}
                    </td>
                    <td className="py-3 px-3.5 text-center text-muted-foreground">{ot.avgTurnaround}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Card>

        <Card className="lg:col-span-4 border bg-card shadow-2xs">
          <CardHeader className="pb-3 border-b">
            <CardTitle className="text-sm font-bold">Fulfillment Pipeline Status</CardTitle>
            <CardDescription className="text-xs">Live stage distribution.</CardDescription>
          </CardHeader>
          <CardContent className="pt-4 space-y-3">
            {fulfillmentStages.map((fs) => (
              <div key={fs.stage} className="flex items-center justify-between p-2.5 rounded-lg border bg-muted/20 text-xs">
                <div className="flex items-center gap-2">
                  <span className={`h-2.5 w-2.5 rounded-full ${fs.color}`} />
                  <span className="font-semibold text-foreground">{fs.stage}</span>
                </div>
                <Badge variant="outline" className="font-bold tabular-nums text-xs">
                  {fs.count}
                </Badge>
              </div>
            ))}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
