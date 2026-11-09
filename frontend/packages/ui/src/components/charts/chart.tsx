import * as React from "react"
import ReactECharts from "echarts-for-react"
import type { EChartsOption } from "echarts"
import { useTheme } from "next-themes"
import { cn } from "@repo/ui/lib/utils"

export const CRM_CHART_COLORS = [
  "#2563EB", // Blue
  "#16A34A", // Green
  "#D97706", // Amber
  "#7C3AED", // Purple
  "#EC4899", // Pink
  "#06B6D4", // Cyan
  "#EA580C", // Orange
]

export interface ChartProps {
  option: EChartsOption
  height?: number | string
  loading?: boolean
  className?: string
  onEvents?: Record<string, (params: unknown) => void>
}

export function Chart({
  option,
  height = 300,
  loading = false,
  className,
  onEvents,
}: ChartProps) {
  const { theme } = useTheme()
  const isDark = theme === "dark"

  const themedOption = React.useMemo(() => {
    return {
      color: CRM_CHART_COLORS,
      backgroundColor: "transparent",
      textStyle: {
        fontFamily: "inherit",
        color: isDark ? "#A1A1AA" : "#64748B",
      },
      tooltip: {
        backgroundColor: isDark ? "#18181B" : "#FFFFFF",
        borderColor: isDark ? "#27272A" : "#E2E8F0",
        borderWidth: 1,
        textStyle: {
          color: isDark ? "#FAFAFA" : "#0F172A",
          fontSize: 12,
        },
        padding: [6, 10],
        extraCssText: "box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1); border-radius: 6px;",
      },
      ...option,
    }
  }, [option, isDark])

  return (
    <div className={cn("w-full relative", className)} style={{ height }}>
      <ReactECharts
        option={themedOption}
        style={{ height: "100%", width: "100%" }}
        showLoading={loading}
        loadingOption={{
          text: "",
          color: "#3b82f6",
          textColor: "#64748b",
          maskColor: isDark ? "rgba(15, 23, 42, 0.6)" : "rgba(255, 255, 255, 0.6)",
        }}
        onEvents={onEvents}
        notMerge={true}
        lazyUpdate={true}
      />
    </div>
  )
}
