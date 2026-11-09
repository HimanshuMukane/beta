import * as React from 'react'
import {
  IconSearch,
  IconCheck,
  IconArrowsUpDown,
  IconChevronDown,
  IconX,
} from '@tabler/icons-react'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
  Input,
  Checkbox,
  cn,
} from '@repo/ui'
import { AVAILABLE_DELIVERABLES } from '../data/mock-leads'

interface DeliverablesDropdownProps {
  selectedDeliverables: string[]
  onChange: (deliverables: string[]) => void
}

export function DeliverablesDropdown({
  selectedDeliverables,
  onChange,
}: DeliverablesDropdownProps) {
  const [open, setOpen] = React.useState(false)
  const [search, setSearch] = React.useState('')

  const categories = React.useMemo(() => {
    const cats = Array.from(new Set(AVAILABLE_DELIVERABLES.map((d) => d.category)))
    return cats
  }, [])

  const filteredItems = React.useMemo(() => {
    if (!search.trim()) return AVAILABLE_DELIVERABLES
    const q = search.toLowerCase()
    return AVAILABLE_DELIVERABLES.filter((d) => d.name.toLowerCase().includes(q))
  }, [search])

  const toggleItem = (name: string) => {
    if (selectedDeliverables.includes(name)) {
      onChange(selectedDeliverables.filter((n) => n !== name))
    } else {
      onChange([...selectedDeliverables, name])
    }
  }

  const toggleCategory = (category: string) => {
    const itemsInCat = AVAILABLE_DELIVERABLES.filter((d) => d.category === category).map(
      (d) => d.name
    )
    const allSelected = itemsInCat.every((n) => selectedDeliverables.includes(n))
    if (allSelected) {
      onChange(selectedDeliverables.filter((n) => !itemsInCat.includes(n)))
    } else {
      const merged = Array.from(new Set([...selectedDeliverables, ...itemsInCat]))
      onChange(merged)
    }
  }

  return (
    <div className="space-y-2">
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <div
            role="button"
            tabIndex={0}
            className="flex min-h-[46px] w-full cursor-pointer items-center justify-between rounded-lg border border-input bg-background px-3.5 py-2 text-sm shadow-xs transition hover:border-primary/50 focus:outline-none focus:ring-2 focus:ring-ring"
          >
            <div className="flex flex-wrap items-center gap-1.5 flex-1 pr-2">
              {selectedDeliverables.length === 0 ? (
                <span className="text-muted-foreground">Choose Deliverables</span>
              ) : (
                selectedDeliverables.map((item) => (
                  <span
                    key={item}
                    className="inline-flex items-center gap-1 rounded-md border border-border bg-muted/60 px-2 py-0.5 text-xs font-medium text-foreground transition hover:bg-muted"
                  >
                    {item}
                    <button
                      type="button"
                      onClick={(e) => {
                        e.stopPropagation()
                        toggleItem(item)
                      }}
                      className="rounded-sm hover:text-destructive focus:outline-none"
                    >
                      <IconX className="h-3 w-3" />
                    </button>
                  </span>
                ))
              )}
            </div>
            <div className="flex items-center gap-1 text-muted-foreground shrink-0">
              <IconArrowsUpDown className="h-4 w-4 opacity-70" />
              <IconChevronDown className="h-4 w-4 opacity-50" />
            </div>
          </div>
        </PopoverTrigger>

        <PopoverContent
          align="start"
          className="w-[340px] sm:w-[400px] p-0 shadow-lg border rounded-xl overflow-hidden"
        >
          {/* Search Bar */}
          <div className="flex items-center gap-2 border-b px-3 py-2 bg-muted/20">
            <IconSearch className="h-4 w-4 text-muted-foreground shrink-0" />
            <Input
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Search deliverables..."
              className="h-8 border-0 bg-transparent p-0 text-sm focus-visible:ring-0 placeholder:text-muted-foreground"
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

          {/* Grouped Deliverables List */}
          <div className="max-h-72 overflow-y-auto divide-y divide-border/60 p-1">
            {categories.map((cat) => {
              const catItems = filteredItems.filter((i) => i.category === cat)
              if (catItems.length === 0) return null

              const allChecked = catItems.every((i) => selectedDeliverables.includes(i.name))
              const someChecked =
                !allChecked && catItems.some((i) => selectedDeliverables.includes(i.name))

              return (
                <div key={cat} className="py-2 first:pt-1 last:pb-1">
                  <div className="flex items-center justify-between px-3 py-1">
                    <span className="text-xs font-semibold tracking-wide text-muted-foreground uppercase">
                      {cat}
                    </span>
                    <Checkbox
                      checked={allChecked ? true : someChecked ? 'indeterminate' : false}
                      onCheckedChange={() => toggleCategory(cat)}
                      className="h-4 w-4"
                    />
                  </div>

                  <div className="mt-1 space-y-0.5">
                    {catItems.map((item) => {
                      const isSelected = selectedDeliverables.includes(item.name)
                      return (
                        <div
                          key={item.id}
                          onClick={() => toggleItem(item.name)}
                          className={cn(
                            'flex cursor-pointer items-center justify-between rounded-md px-3 py-1.5 text-sm transition',
                            isSelected
                              ? 'bg-accent/40 font-medium text-foreground'
                              : 'text-muted-foreground hover:bg-muted/50 hover:text-foreground'
                          )}
                        >
                          <span className={cn(isSelected && 'font-semibold text-foreground')}>
                            {item.name}
                          </span>
                          {isSelected && <IconCheck className="h-4 w-4 text-emerald-600 dark:text-emerald-400" />}
                        </div>
                      )
                    })}
                  </div>
                </div>
              )
            })}
          </div>

          {/* Quick count footer */}
          <div className="flex items-center justify-between border-t bg-muted/30 px-3 py-2 text-xs text-muted-foreground">
            <span>{selectedDeliverables.length} deliverable(s) chosen</span>
            {selectedDeliverables.length > 0 && (
              <button
                type="button"
                onClick={() => onChange([])}
                className="text-xs font-medium text-destructive hover:underline"
              >
                Reset all
              </button>
            )}
          </div>
        </PopoverContent>
      </Popover>
    </div>
  )
}
