import * as React from 'react'
import {
  IconCalendar,
  IconChevronLeft,
  IconChevronRight,
  IconChevronsLeft,
  IconChevronsRight,
} from "../../icons/tabler-icons"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '../ui/popover'
import { Button } from '../ui/button'
import { cn } from '../../lib/utils'

interface DateTimePickerCustomProps {
  value: string // ISO string or "YYYY-MM-DDTHH:mm"
  onChange: (value: string) => void
  placeholder?: string
  className?: string
}

const MONTH_NAMES = [
  'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
  'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec',
]
const DAY_NAMES = ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa']

export function DateTimePickerCustom({
  value,
  onChange,
  placeholder = 'Select date & time',
  className,
}: DateTimePickerCustomProps) {
  const [open, setOpen] = React.useState(false)

  // Parse initial date or default to now
  const parsedDate = React.useMemo(() => {
    if (!value) return new Date()
    const d = new Date(value)
    return isNaN(d.getTime()) ? new Date() : d
  }, [value])

  // Working state when popup is open
  const [viewYear, setViewYear] = React.useState(parsedDate.getFullYear())
  const [viewMonth, setViewMonth] = React.useState(parsedDate.getMonth())
  const [selectedDay, setSelectedDay] = React.useState<number>(parsedDate.getDate())
  const [selectedHour, setSelectedHour] = React.useState<number>(parsedDate.getHours())
  const [selectedMinute, setSelectedMinute] = React.useState<number>(parsedDate.getMinutes())

  // Keep in sync when value changes from outside
  React.useEffect(() => {
    if (value) {
      const d = new Date(value)
      if (!isNaN(d.getTime())) {
        setViewYear(d.getFullYear())
        setViewMonth(d.getMonth())
        setSelectedDay(d.getDate())
        setSelectedHour(d.getHours())
        setSelectedMinute(d.getMinutes())
      }
    }
  }, [value, open])

  // Month navigation
  const prevMonth = () => {
    if (viewMonth === 0) {
      setViewMonth(11)
      setViewYear((y) => y - 1)
    } else {
      setViewMonth((m) => m - 1)
    }
  }

  const nextMonth = () => {
    if (viewMonth === 11) {
      setViewMonth(0)
      setViewYear((y) => y + 1)
    } else {
      setViewMonth((m) => m + 1)
    }
  }

  const prevYear = () => setViewYear((y) => y - 1)
  const nextYear = () => setViewYear((y) => y + 1)

  // Calendar matrix calculation
  const calendarDays = React.useMemo(() => {
    const firstDayIndex = new Date(viewYear, viewMonth, 1).getDay()
    const totalDaysInMonth = new Date(viewYear, viewMonth + 1, 0).getDate()
    const prevMonthDays = new Date(viewYear, viewMonth, 0).getDate()

    const days: { day: number; currentMonth: boolean; dateObj: Date }[] = []

    // Previous month padding
    for (let i = firstDayIndex - 1; i >= 0; i--) {
      const dayNum = prevMonthDays - i
      days.push({
        day: dayNum,
        currentMonth: false,
        dateObj: new Date(viewYear, viewMonth - 1, dayNum),
      })
    }

    // Current month days
    for (let i = 1; i <= totalDaysInMonth; i++) {
      days.push({
        day: i,
        currentMonth: true,
        dateObj: new Date(viewYear, viewMonth, i),
      })
    }

    // Next month padding to fill rows
    const remaining = (7 - (days.length % 7)) % 7
    for (let i = 1; i <= remaining; i++) {
      days.push({
        day: i,
        currentMonth: false,
        dateObj: new Date(viewYear, viewMonth + 1, i),
      })
    }

    return days
  }, [viewYear, viewMonth])

  // Today indicators
  const today = new Date()
  const isToday = (dayNum: number, currentMonth: boolean) => {
    return (
      currentMonth &&
      today.getFullYear() === viewYear &&
      today.getMonth() === viewMonth &&
      today.getDate() === dayNum
    )
  }

  const isSelected = (dayNum: number, currentMonth: boolean) => {
    return (
      currentMonth &&
      parsedDate.getFullYear() === viewYear &&
      parsedDate.getMonth() === viewMonth &&
      selectedDay === dayNum
    )
  }

  // Hours: 0 to 23
  const hours = Array.from({ length: 24 }, (_, i) => i)
  // Minutes: 0 to 59
  const minutes = Array.from({ length: 60 }, (_, i) => i)

  // Set to current time (Now button)
  const handleSetNow = () => {
    const n = new Date()
    setViewYear(n.getFullYear())
    setViewMonth(n.getMonth())
    setSelectedDay(n.getDate())
    setSelectedHour(n.getHours())
    setSelectedMinute(n.getMinutes())

    const pad = (num: number) => String(num).padStart(2, '0')
    const formatted = `${n.getFullYear()}-${pad(n.getMonth() + 1)}-${pad(n.getDate())}T${pad(
      n.getHours()
    )}:${pad(n.getMinutes())}`
    onChange(formatted)
    setOpen(false)
  }

  // Confirm selection (OK button)
  const handleConfirm = () => {
    const pad = (num: number) => String(num).padStart(2, '0')
    const formatted = `${viewYear}-${pad(viewMonth + 1)}-${pad(selectedDay)}T${pad(
      selectedHour
    )}:${pad(selectedMinute)}`
    onChange(formatted)
    setOpen(false)
  }

  // Format display text on trigger
  const displayValue = React.useMemo(() => {
    if (!value) return ''
    const d = new Date(value)
    if (isNaN(d.getTime())) return value
    const pad = (num: number) => String(num).padStart(2, '0')
    const day = pad(d.getDate())
    const month = MONTH_NAMES[d.getMonth()]
    const year = d.getFullYear()
    const hoursStr = pad(d.getHours())
    const minsStr = pad(d.getMinutes())
    return `${day} ${month} ${year}, ${hoursStr}:${minsStr}`
  }, [value])

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <button
          type="button"
          className={cn(
            'flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm shadow-xs transition hover:border-primary/50 focus:outline-none focus:ring-2 focus:ring-ring text-left',
            !value && 'text-muted-foreground',
            className
          )}
        >
          <span className={cn('truncate', !value && 'text-muted-foreground')}>
            {displayValue || placeholder}
          </span>
          <IconCalendar className="h-4 w-4 text-muted-foreground shrink-0 opacity-70" />
        </button>
      </PopoverTrigger>

      <PopoverContent
        align="start"
        className="w-auto p-0 shadow-xl border rounded-xl overflow-hidden bg-background"
      >
        <div className="flex divide-x divide-border">
          {/* Calendar Left Section */}
          <div className="p-3 w-[280px]">
            {/* Header: Prev Year/Month, Month Year title, Next Month/Year */}
            <div className="flex items-center justify-between pb-2 mb-1 border-b border-border/60">
              <div className="flex items-center gap-0.5">
                <button
                  type="button"
                  onClick={prevYear}
                  title="Previous Year"
                  className="h-6 w-6 rounded flex items-center justify-center text-muted-foreground hover:bg-muted hover:text-foreground"
                >
                  <IconChevronsLeft className="h-3.5 w-3.5" />
                </button>
                <button
                  type="button"
                  onClick={prevMonth}
                  title="Previous Month"
                  className="h-6 w-6 rounded flex items-center justify-center text-muted-foreground hover:bg-muted hover:text-foreground"
                >
                  <IconChevronLeft className="h-3.5 w-3.5" />
                </button>
              </div>

              <span className="text-xs font-semibold text-foreground tracking-wide">
                {MONTH_NAMES[viewMonth]} {viewYear}
              </span>

              <div className="flex items-center gap-0.5">
                <button
                  type="button"
                  onClick={nextMonth}
                  title="Next Month"
                  className="h-6 w-6 rounded flex items-center justify-center text-muted-foreground hover:bg-muted hover:text-foreground"
                >
                  <IconChevronRight className="h-3.5 w-3.5" />
                </button>
                <button
                  type="button"
                  onClick={nextYear}
                  title="Next Year"
                  className="h-6 w-6 rounded flex items-center justify-center text-muted-foreground hover:bg-muted hover:text-foreground"
                >
                  <IconChevronsRight className="h-3.5 w-3.5" />
                </button>
              </div>
            </div>

            {/* Days Header */}
            <div className="grid grid-cols-7 text-center text-[11px] font-medium text-muted-foreground py-1">
              {DAY_NAMES.map((d) => (
                <div key={d}>{d}</div>
              ))}
            </div>

            {/* Calendar Grid */}
            <div className="grid grid-cols-7 gap-1 text-center text-xs mt-1">
              {calendarDays.map((item, idx) => {
                const selected = isSelected(item.day, item.currentMonth)
                const currentDayToday = isToday(item.day, item.currentMonth)

                return (
                  <button
                    key={idx}
                    type="button"
                    disabled={!item.currentMonth}
                    onClick={() => {
                      if (item.currentMonth) setSelectedDay(item.day)
                    }}
                    className={cn(
                      'h-7 w-7 mx-auto flex items-center justify-center rounded-full text-xs transition font-normal',
                      !item.currentMonth && 'text-muted-foreground/30 pointer-events-none',
                      item.currentMonth && !selected && !currentDayToday && 'text-foreground hover:bg-muted',
                      currentDayToday && !selected && 'border border-emerald-500 text-emerald-600 dark:text-emerald-400 font-semibold',
                      selected && 'bg-emerald-600 text-white font-semibold hover:bg-emerald-700'
                    )}
                  >
                    {item.day}
                  </button>
                )
              })}
            </div>
          </div>

          {/* Time Picker Right Section */}
          <div className="w-[125px] flex flex-col bg-muted/10">
            {/* Header: Currently selected time */}
            <div className="h-9 border-b border-border/60 flex items-center justify-center text-xs font-bold text-foreground tracking-wider tabular-nums">
              {String(selectedHour).padStart(2, '0')}:{String(selectedMinute).padStart(2, '0')}
            </div>

            {/* Two columns: Hours & Minutes scroll lists */}
            <div className="flex-1 flex divide-x divide-border/60 max-h-56">
              {/* Hours Column */}
              <div className="flex-1 overflow-y-auto py-1 scrollbar-thin text-center text-xs">
                {hours.map((h) => {
                  const active = selectedHour === h
                  return (
                    <button
                      key={h}
                      type="button"
                      onClick={() => setSelectedHour(h)}
                      className={cn(
                        'w-full py-1 tabular-nums transition text-xs block',
                        active
                          ? 'text-emerald-600 dark:text-emerald-400 font-bold bg-emerald-500/10'
                          : 'text-muted-foreground hover:text-foreground hover:bg-muted/40'
                      )}
                    >
                      {String(h).padStart(2, '0')}
                    </button>
                  )
                })}
              </div>

              {/* Minutes Column */}
              <div className="flex-1 overflow-y-auto py-1 scrollbar-thin text-center text-xs">
                {minutes.map((m) => {
                  const active = selectedMinute === m
                  return (
                    <button
                      key={m}
                      type="button"
                      onClick={() => setSelectedMinute(m)}
                      className={cn(
                        'w-full py-1 tabular-nums transition text-xs block',
                        active
                          ? 'text-emerald-600 dark:text-emerald-400 font-bold bg-emerald-500/10'
                          : 'text-muted-foreground hover:text-foreground hover:bg-muted/40'
                      )}
                    >
                      {String(m).padStart(2, '0')}
                    </button>
                  )
                })}
              </div>
            </div>
          </div>
        </div>

        {/* Footer Actions */}
        <div className="flex items-center justify-between border-t border-border/80 bg-muted/30 px-3 py-2">
          <button
            type="button"
            onClick={handleSetNow}
            className="text-xs font-semibold text-emerald-600 dark:text-emerald-400 hover:underline"
          >
            Now
          </button>
          <div className="flex items-center gap-1.5">
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={() => setOpen(false)}
              className="h-7 px-2.5 text-xs text-muted-foreground"
            >
              Cancel
            </Button>
            <Button
              type="button"
              size="sm"
              onClick={handleConfirm}
              className="h-7 px-4 text-xs font-semibold bg-[#113e3b] dark:bg-emerald-600 hover:bg-[#0c2e2c] dark:hover:bg-emerald-700 text-white"
            >
              OK
            </Button>
          </div>
        </div>
      </PopoverContent>
    </Popover>
  )
}

