import * as React from "react"
import ReactECharts from "echarts-for-react"
import type { EChartsOption } from "echarts"
import { useTheme } from "next-themes"
import { cn } from "@repo/ui/lib/utils"

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
      backgroundColor: "transparent",
      textStyle: {
        fontFamily: "inherit",
        color: isDark ? "#94a3b8" : "#64748b",
      },
      tooltip: {
        backgroundColor: isDark ? "rgba(15, 23, 42, 0.95)" : "rgba(255, 255, 255, 0.95)",
        borderColor: isDark ? "#334155" : "#e2e8f0",
        borderWidth: 1,
        textStyle: {
          color: isDark ? "#f8fafc" : "#0f172a",
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
