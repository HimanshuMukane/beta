import * as React from 'react'
import {
  IconCalendar,
  IconChevronLeft,
  IconChevronRight,
  IconChevronsLeft,
  IconChevronsRight,
  IconX,
  IconCheck,
} from "@repo/ui/icons/tabler-icons"
import {
  format,
  subDays,
  startOfMonth,
  endOfMonth,
  subMonths,
  startOfQuarter,
  endOfQuarter,
  subQuarters,
  startOfYear,
  endOfYear,
  subYears,
  isSameDay,
  isWithinInterval,
  isBefore,
  addMonths,
  subMonths as subMonthsFn,
  addYears,
  subYears as subYearsFn,
} from 'date-fns'
import { Button } from '@repo/ui/components/ui/button'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@repo/ui/components/ui/popover'
import { cn } from '@repo/ui/lib/utils'

export interface DateRange {
  startDate: Date
  endDate: Date
}

interface BillingDateRangePickerProps {
  labelPrefix?: string
  dateRange: DateRange
  onChange: (range: DateRange) => void
  onClear?: () => void
  className?: string
}

const PRESETS = [
  {
    label: 'Today',
    getRange: () => {
      const today = new Date(2026, 10, 9)
      return { startDate: today, endDate: today }
    },
  },
  {
    label: 'Yesterday',
    getRange: () => {
      const base = new Date(2026, 10, 9)
      const yesterday = subDays(base, 1)
      return { startDate: yesterday, endDate: yesterday }
    },
  },
  {
    label: 'Last 7 Days',
    getRange: () => {
      const base = new Date(2026, 10, 9)
      return { startDate: subDays(base, 6), endDate: base }
    },
  },
  {
    label: 'This Month',
    getRange: () => {
      const base = new Date(2026, 10, 9)
      return { startDate: startOfMonth(base), endDate: endOfMonth(base) }
    },
  },
  {
    label: 'Last Month',
    getRange: () => {
      const base = new Date(2026, 10, 9)
      const prevMonth = subMonths(base, 1)
      return { startDate: startOfMonth(prevMonth), endDate: endOfMonth(prevMonth) }
    },
  },
  {
    label: 'This Quarter',
    getRange: () => {
      const base = new Date(2026, 10, 9)
      return { startDate: startOfQuarter(base), endDate: endOfQuarter(base) }
    },
  },
  {
    label: 'Last Quarter',
    getRange: () => {
      const base = new Date(2026, 10, 9)
      const prevQ = subQuarters(base, 1)
      return { startDate: startOfQuarter(prevQ), endDate: endOfQuarter(prevQ) }
    },
  },
  {
    label: 'This Year',
    getRange: () => {
      const base = new Date(2026, 10, 9)
      return { startDate: startOfYear(base), endDate: endOfYear(base) }
    },
  },
  {
    label: 'Last Year',
    getRange: () => {
      const base = new Date(2026, 10, 9)
      const prevY = subYears(base, 1)
      return { startDate: startOfYear(prevY), endDate: endOfYear(prevY) }
    },
  },
  {
    label: 'All Time',
    getRange: () => ({
      startDate: new Date(2000, 0, 1),
      endDate: new Date(2026, 11, 31),
    }),
  },
]

const DAYS_HEADER = ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa']

export function BillingDateRangePicker({
  labelPrefix = 'Payments Date Range',
  dateRange,
  onChange,
  onClear,
  className,
}: BillingDateRangePickerProps) {
  const [open, setOpen] = React.useState(false)
  const [leftMonth, setLeftMonth] = React.useState<Date>(() => startOfMonth(dateRange.startDate))
  const [tempRange, setTempRange] = React.useState<DateRange>(dateRange)
  const [hoverDate, setHoverDate] = React.useState<Date | null>(null)
  const [selectingStart, setSelectingStart] = React.useState(false)

  // Keep synced when dateRange props change
  React.useEffect(() => {
    setTempRange(dateRange)
  }, [dateRange])

  const rightMonth = React.useMemo(() => addMonths(leftMonth, 1), [leftMonth])

  const handlePrevYear = () => setLeftMonth((prev) => subYearsFn(prev, 1))
  const handlePrevMonth = () => setLeftMonth((prev) => subMonthsFn(prev, 1))
  const handleNextMonth = () => setLeftMonth((prev) => addMonths(prev, 1))
  const handleNextYear = () => setLeftMonth((prev) => addYears(prev, 1))

  const handleDayClick = (day: Date) => {
    if (!selectingStart) {
      setTempRange({ startDate: day, endDate: day })
      setSelectingStart(true)
    } else {
      let nextStart = tempRange.startDate
      let nextEnd = day
      if (isBefore(day, nextStart)) {
        nextEnd = nextStart
        nextStart = day
      }
      const newRange = { startDate: nextStart, endDate: nextEnd }
      setTempRange(newRange)
      setSelectingStart(false)
      onChange(newRange)
      setOpen(false)
    }
  }

  const handlePresetClick = (presetFn: () => DateRange) => {
    const range = presetFn()
    setTempRange(range)
    setLeftMonth(startOfMonth(range.startDate))
    setSelectingStart(false)
    onChange(range)
    setOpen(false)
  }

  const renderMonthCalendar = (monthDate: Date) => {
    const monthStart = startOfMonth(monthDate)
    const startDayOfWeek = monthStart.getDay()
    const daysInMonth = new Date(monthDate.getFullYear(), monthDate.getMonth() + 1, 0).getDate()
    const prevMonthDays = new Date(monthDate.getFullYear(), monthDate.getMonth(), 0).getDate()

    const calendarCells: { date: Date; currentMonth: boolean }[] = []

    // Previous month filler days
    for (let i = startDayOfWeek - 1; i >= 0; i--) {
      calendarCells.push({
        date: new Date(monthDate.getFullYear(), monthDate.getMonth() - 1, prevMonthDays - i),
        currentMonth: false,
      })
    }

    // Current month days
    for (let d = 1; d <= daysInMonth; d++) {
      calendarCells.push({
        date: new Date(monthDate.getFullYear(), monthDate.getMonth(), d),
        currentMonth: true,
      })
    }

    // Next month filler to reach 42 cells (6 rows x 7 cols) or 35
    const totalCells = calendarCells.length <= 35 ? 35 : 42
    const remaining = totalCells - calendarCells.length
    for (let d = 1; d <= remaining; d++) {
      calendarCells.push({
        date: new Date(monthDate.getFullYear(), monthDate.getMonth() + 1, d),
        currentMonth: false,
      })
    }

    const effectiveEnd = selectingStart && hoverDate
      ? (isBefore(hoverDate, tempRange.startDate) ? tempRange.startDate : hoverDate)
      : tempRange.endDate
    const effectiveStart = selectingStart && hoverDate && isBefore(hoverDate, tempRange.startDate)
      ? hoverDate
      : tempRange.startDate

    return (
      <div className="w-[230px] select-none">
        <div className="text-center font-semibold text-xs text-foreground mb-3">
          {format(monthDate, 'MMM yyyy')}
        </div>
        <div className="grid grid-cols-7 gap-1 text-[11px] font-medium text-muted-foreground text-center mb-1">
          {DAYS_HEADER.map((d) => (
            <div key={d} className="py-1">
              {d}
            </div>
          ))}
        </div>
        <div className="grid grid-cols-7 gap-y-1">
          {calendarCells.map((cell, idx) => {
            const isStart = isSameDay(cell.date, effectiveStart)
            const isEnd = isSameDay(cell.date, effectiveEnd)
            const inInterval = isWithinInterval(cell.date, {
              start: effectiveStart,
              end: effectiveEnd,
            })

            return (
              <button
                key={idx}
                type="button"
                onClick={() => handleDayClick(cell.date)}
                onMouseEnter={() => selectingStart && setHoverDate(cell.date)}
                className={cn(
                  'h-7 w-7 text-xs rounded-full flex items-center justify-center transition-colors mx-auto',
                  !cell.currentMonth && 'text-muted-foreground/40 font-normal',
                  cell.currentMonth && 'text-foreground font-medium hover:bg-muted',
                  inInterval && !isStart && !isEnd && 'bg-primary/10 text-primary rounded-none first:rounded-l-full last:rounded-r-full',
                  (isStart || isEnd) &&
                    'bg-emerald-600 text-white font-semibold hover:bg-emerald-700 shadow-xs'
                )}
              >
                {cell.date.getDate()}
              </button>
            )
          })}
        </div>
      </div>
    )
  }

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <button
          type="button"
          className={cn(
            'inline-flex items-center gap-2 rounded-lg border border-input bg-card px-3 py-1.5 text-xs text-muted-foreground shadow-xs transition hover:border-primary/50 hover:text-foreground focus:outline-hidden',
            className
          )}
        >
          <span className="font-semibold text-foreground">{labelPrefix}:</span>
          <span className="font-medium text-foreground">
            {format(dateRange.startDate, 'MMM dd, yyyy')}
          </span>
          <span className="text-muted-foreground">→</span>
          <span className="font-medium text-foreground">
            {format(dateRange.endDate, 'MMM dd, yyyy')}
          </span>
          <IconCalendar className="h-4 w-4 text-muted-foreground ml-1" />
          {onClear && (
            <span
              role="button"
              tabIndex={0}
              onClick={(e) => {
                e.stopPropagation()
                onClear()
              }}
              className="ml-1 rounded-full p-0.5 text-muted-foreground hover:bg-muted hover:text-foreground"
            >
              <IconX className="h-3 w-3" />
            </span>
          )}
        </button>
      </PopoverTrigger>

      <PopoverContent
        align="end"
        sideOffset={6}
        className="w-auto p-0 shadow-xl border border-border/80 rounded-xl overflow-hidden bg-card"
      >
        <div className="flex divide-x divide-border">
          {/* Quick Select Sidebar */}
          <div className="w-36 bg-muted/20 p-3 space-y-1">
            <div className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground px-2 pb-1.5">
              Quick Select
            </div>
            {PRESETS.map((p) => (
              <button
                key={p.label}
                type="button"
                onClick={() => handlePresetClick(p.getRange)}
                className="w-full text-left px-2 py-1.5 text-xs rounded-md font-medium text-muted-foreground hover:bg-primary/10 hover:text-primary transition"
              >
                {p.label}
              </button>
            ))}
          </div>

          {/* Dual Calendar Views */}
          <div className="p-4">
            {/* Header Controls */}
            <div className="flex items-center justify-between pb-3 mb-2 border-b">
              <div className="flex items-center gap-1">
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-7 w-7 p-0 text-muted-foreground"
                  onClick={handlePrevYear}
                  title="Previous Year"
                >
                  <IconChevronsLeft className="h-3.5 w-3.5" />
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-7 w-7 p-0 text-muted-foreground"
                  onClick={handlePrevMonth}
                  title="Previous Month"
                >
                  <IconChevronLeft className="h-3.5 w-3.5" />
                </Button>
              </div>

              <div className="text-xs font-semibold text-muted-foreground">
                Select Date Range
              </div>

              <div className="flex items-center gap-1">
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-7 w-7 p-0 text-muted-foreground"
                  onClick={handleNextMonth}
                  title="Next Month"
                >
                  <IconChevronRight className="h-3.5 w-3.5" />
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-7 w-7 p-0 text-muted-foreground"
                  onClick={handleNextYear}
                  title="Next Year"
                >
                  <IconChevronsRight className="h-3.5 w-3.5" />
                </Button>
              </div>
            </div>

            <div className="flex items-start gap-6">
              {renderMonthCalendar(leftMonth)}
              {renderMonthCalendar(rightMonth)}
            </div>

            {/* Footer with actions */}
            <div className="mt-4 pt-3 border-t flex items-center justify-between">
              <div className="text-[11px] text-muted-foreground">
                Selected:{' '}
                <span className="font-semibold text-foreground">
                  {format(tempRange.startDate, 'MMM dd, yyyy')}
                </span>{' '}
                to{' '}
                <span className="font-semibold text-foreground">
                  {format(tempRange.endDate, 'MMM dd, yyyy')}
                </span>
              </div>
              <div className="flex items-center gap-2">
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => setOpen(false)}
                >
                  Cancel
                </Button>
                <Button
                  size="sm"
                  variant="default"
                  onClick={() => {
                    onChange(tempRange)
                    setOpen(false)
                  }}
                >
                  <IconCheck className="h-3.5 w-3.5 mr-1" />
                  Apply Range
                </Button>
              </div>
            </div>
          </div>
        </div>
      </PopoverContent>
    </Popover>
  )
}
