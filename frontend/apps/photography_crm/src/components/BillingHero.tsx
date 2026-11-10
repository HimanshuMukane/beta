import * as React from 'react'
import { Card, CardContent } from '@repo/ui/components/ui/card'
import type { EChartsOption } from 'echarts'

const Chart = React.lazy(() =>
  import('@repo/ui/components/charts/chart').then((m) => ({ default: m.Chart }))
)

interface BillingHeroProps {
  cards: {
    label: string
    value: string
    subtitle?: string
    colorClass?: string
    badge?: string
  }[]
  chartTitle: string
  chartData: {
    labels: string[]
    values: number[]
  }
  chartColor?: string
  dateRangeFilterNode: React.ReactNode
}

export function BillingHero({
  cards,
  chartTitle,
  chartData,
  chartColor = '#2563EB',
  dateRangeFilterNode,
}: BillingHeroProps) {
  const chartOption: EChartsOption = React.useMemo(() => {
    return {
      tooltip: {
        trigger: 'axis',
        axisPointer: { type: 'shadow' },
        formatter: (params: unknown) => {
          if (Array.isArray(params) && params.length > 0) {
            const p = params[0] as { name: string; value: number }
            return `<div class="p-1"><div class="text-xs font-semibold text-slate-500 mb-1">${p.name}</div><div class="text-sm font-bold text-foreground">₹${p.value?.toLocaleString('en-IN') ?? 0}</div></div>`
          }
          return ''
        },
      },
      grid: {
        top: 25,
        right: 20,
        bottom: 30,
        left: 45,
        containLabel: true,
      },
      xAxis: {
        type: 'category',
        data: chartData.labels,
        axisLine: { lineStyle: { color: 'hsl(var(--border))' } },
        axisLabel: {
          color: 'hsl(var(--muted-foreground))',
          fontSize: 11,
          fontWeight: 500,
        },
        axisTick: { show: false },
      },
      yAxis: {
        type: 'value',
        axisLine: { show: false },
        splitLine: {
          lineStyle: {
            color: 'hsl(var(--border))',
            type: 'dashed',
            opacity: 0.6,
          },
        },
        axisLabel: {
          color: 'hsl(var(--muted-foreground))',
          fontSize: 11,
          formatter: (val: number) => {
            if (val >= 100000) return `₹${val / 100000}L`
            if (val >= 1000) return `₹${val / 1000}k`
            return `₹${val}`
          },
        },
      },
      series: [
        {
          name: chartTitle,
          type: 'bar',
          data: chartData.values,
          barMaxWidth: 110,
          barMinWidth: 24,
          itemStyle: {
            color: chartColor,
            borderRadius: [6, 6, 0, 0],
          },
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowColor: 'rgba(0, 0, 0, 0.15)',
            },
          },
        },
      ],
    }
  }, [chartData, chartColor, chartTitle])

  return (
    <div className="rounded-2xl border bg-card/60 p-4 sm:p-5 shadow-xs mb-6 backdrop-blur-xs">
      {/* Top Bar with Date Range Filter aligned to the right */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pb-4">
        <div>
          <h2 className="text-sm font-bold tracking-tight text-foreground uppercase">
            Overview Summary
          </h2>
          <p className="text-xs text-muted-foreground">
            Financial analytics and period breakdown
          </p>
        </div>
        <div className="self-end sm:self-auto">{dateRangeFilterNode}</div>
      </div>

      {/* Main Grid: Cards on the left + Chart on the right */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-4 items-stretch">
        {/* Left Column: Metric Cards */}
        <div
          className={`lg:col-span-4 flex flex-col justify-center gap-3 ${
            cards.length > 2 ? 'lg:col-span-4' : 'lg:col-span-3'
          }`}
        >
          {cards.map((c, i) => (
            <Card
              key={i}
              className="border border-border/80 bg-background/80 shadow-xs hover:border-primary/40 transition-colors"
            >
              <CardContent className="p-4 sm:p-5">
                <div className="flex items-center justify-between gap-2 mb-1.5">
                  <span className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                    {c.label}
                  </span>
                  {c.badge && (
                    <span className="text-[10px] font-bold px-1.5 py-0.5 rounded-full bg-primary/10 text-primary">
                      {c.badge}
                    </span>
                  )}
                </div>
                <div
                  className={`text-2xl sm:text-3xl font-extrabold tracking-tight ${
                    c.colorClass || 'text-foreground'
                  }`}
                >
                  {c.value}
                </div>
                {c.subtitle && (
                  <p className="text-xs text-muted-foreground mt-1 font-medium">
                    {c.subtitle}
                  </p>
                )}
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Right Column: Chart */}
        <div
          className={`rounded-xl border border-border/70 bg-background/90 p-4 sm:p-5 shadow-xs flex flex-col justify-between ${
            cards.length > 2 ? 'lg:col-span-8' : 'lg:col-span-9'
          }`}
        >
          <div className="flex items-center justify-between border-b pb-2.5 mb-2">
            <h3 className="text-sm font-bold text-foreground">{chartTitle}</h3>
            <span className="text-[11px] font-semibold text-muted-foreground">
              Period Volume
            </span>
          </div>

          <div className="w-full flex-1 min-h-[190px]">
            <React.Suspense fallback={<div className="h-[190px] w-full animate-pulse rounded-lg bg-muted/30" />}>
              <Chart option={chartOption} height={190} />
            </React.Suspense>
          </div>
        </div>
      </div>
    </div>
  )
}

