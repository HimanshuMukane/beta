import * as React from "react"
import { format } from "date-fns"
import { IconCalendar, IconClock } from "../../icons/tabler-icons"
import { Calendar } from "../ui/calendar"
import { Popover, PopoverContent, PopoverTrigger } from "../ui/popover"
import { Button } from "../ui/button"
import { Input } from "../ui/input"
import { cn } from "@repo/ui/lib/utils"

export interface DateTimePickerProps {
  date?: Date
  setDate?: (date: Date | undefined) => void
  placeholder?: string
  includeTime?: boolean
  disabled?: boolean
  className?: string
}

export function DateTimePicker({
  date,
  setDate,
  placeholder = "Pick a date",
  includeTime = false,
  disabled = false,
  className,
}: DateTimePickerProps) {
  const [selectedDate, setSelectedDate] = React.useState<Date | undefined>(date)
  const [timeValue, setTimeValue] = React.useState<string>(
    date ? format(date, "HH:mm") : "12:00"
  )

  React.useEffect(() => {
    setSelectedDate(date)
    if (date) {
      setTimeValue(format(date, "HH:mm"))
    }
  }, [date])

  const handleSelectDate = (newDate: Date | undefined) => {
    if (!newDate) {
      setSelectedDate(undefined)
      setDate?.(undefined)
      return
    }

    if (includeTime && timeValue) {
      const [hours, minutes] = timeValue.split(":").map(Number)
      if (!isNaN(hours) && !isNaN(minutes)) {
        newDate.setHours(hours, minutes, 0, 0)
      }
    }

    setSelectedDate(newDate)
    setDate?.(newDate)
  }

  const handleTimeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value
    setTimeValue(val)
    if (selectedDate && val) {
      const [hours, minutes] = val.split(":").map(Number)
      if (!isNaN(hours) && !isNaN(minutes)) {
        const updated = new Date(selectedDate)
        updated.setHours(hours, minutes, 0, 0)
        setSelectedDate(updated)
        setDate?.(updated)
      }
    }
  }

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          disabled={disabled}
          className={cn(
            "h-8 justify-start text-left font-normal text-xs px-2.5",
            !selectedDate && "text-muted-foreground",
            className
          )}
        >
          <IconCalendar className="mr-1.5 h-3.5 w-3.5 opacity-70" />
          {selectedDate ? (
            includeTime ? (
              format(selectedDate, "PPP p")
            ) : (
              format(selectedDate, "PPP")
            )
          ) : (
            <span>{placeholder}</span>
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0" align="start">
        <Calendar
          mode="single"
          selected={selectedDate}
          onSelect={handleSelectDate}
          initialFocus
        />
        {includeTime && (
          <div className="flex items-center gap-2 border-t p-2.5 bg-muted/20">
            <IconClock className="h-3.5 w-3.5 text-muted-foreground" />
            <span className="text-xs text-muted-foreground">Time:</span>
            <Input
              type="time"
              value={timeValue}
              onChange={handleTimeChange}
              className="h-7 w-28 text-xs px-2"
            />
          </div>
        )}
      </PopoverContent>
    </Popover>
  )
}
