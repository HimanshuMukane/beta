import * as React from 'react'
import {
  IconSearch,
  IconCheck,
  IconChevronDown,
} from "@repo/ui/icons/tabler-icons"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@repo/ui/components/ui/popover'
import { Input } from '@repo/ui/components/ui/input'
import { cn } from '@repo/ui/lib/utils'
import type { LeadSource } from '../types/lead'

export const DETAILED_LEAD_SOURCES = [
  'Google Search',
  'Friend Referral',
  'Family Referral',
  'Instagram Post',
  'Instagram Ads',
  'Google Ads',
  'Studio Shop',
  'Messenger',
  'Whatsapp',
  'Leads Ads',
  'Form',
  'Manual Leads',
]

interface LeadSourceComboboxProps {
  value: string
  onChange: (value: LeadSource) => void
  placeholder?: string
  className?: string
}

export function LeadSourceCombobox({
  value,
  onChange,
  placeholder = 'Select lead source',
  className,
}: LeadSourceComboboxProps) {
  const [open, setOpen] = React.useState(false)
  const [search, setSearch] = React.useState('')

  const filteredSources = React.useMemo(() => {
    if (!search.trim()) return DETAILED_LEAD_SOURCES
    const q = search.toLowerCase()
    return DETAILED_LEAD_SOURCES.filter((s) => s.toLowerCase().includes(q))
  }, [search])

  const handleSelect = (selected: string) => {
    onChange(selected as LeadSource)
    setOpen(false)
    setSearch('')
  }

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
            {value || placeholder}
          </span>
          <IconChevronDown className="h-4 w-4 opacity-50 shrink-0" />
        </button>
      </PopoverTrigger>

      <PopoverContent
        align="start"
        className="w-[280px] sm:w-[320px] p-0 shadow-lg border rounded-xl overflow-hidden bg-background"
      >
        {/* Search input header */}
        <div className="flex items-center gap-2 border-b px-3 py-2 bg-muted/20">
          <IconSearch className="h-4 w-4 text-muted-foreground shrink-0" />
          <Input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search"
            className="h-8 border-0 bg-transparent p-0 text-sm focus-visible:ring-0 placeholder:text-muted-foreground"
            autoFocus
          />
          {search && (
            <button
              type="button"
              onClick={() => setSearch('')}
              className="text-xs text-muted-foreground hover:text-foreground"
            >
              Clear
            </button>
          )}
        </div>

        {/* Options list */}
        <div className="max-h-60 overflow-y-auto p-1 divide-y divide-border/40">
          {filteredSources.length === 0 ? (
            <div className="py-6 text-center text-xs text-muted-foreground">
              No lead source matching "{search}"
            </div>
          ) : (
            filteredSources.map((item) => {
              const isSelected = value === item
              return (
                <div
                  key={item}
                  onClick={() => handleSelect(item)}
                  className={cn(
                    'flex cursor-pointer items-center justify-between rounded-lg px-3 py-2 text-sm transition',
                    isSelected
                      ? 'bg-sky-50 dark:bg-sky-950/40 text-foreground font-semibold'
                      : 'text-muted-foreground hover:bg-muted/50 hover:text-foreground'
                  )}
                >
                  <span>{item}</span>
                  {isSelected && (
                    <IconCheck className="h-4 w-4 text-emerald-600 dark:text-emerald-400" />
                  )}
                </div>
              )
            })
          )}
        </div>
      </PopoverContent>
    </Popover>
  )
}

