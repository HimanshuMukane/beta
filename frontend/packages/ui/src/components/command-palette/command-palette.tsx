import * as React from "react"
import { Command as CommandPrimitive } from "cmdk"
import { IconSearch } from "@tabler/icons-react"
import { Dialog, DialogContent } from "../ui/dialog"

export interface CommandItemType {
  id: string
  label: string
  shortcut?: string[]
  icon?: React.ComponentType<{ className?: string }>
  onSelect?: () => void
  group?: string
  keywords?: string[]
}

export interface CommandPaletteProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  items?: CommandItemType[]
  groups?: string[]
  placeholder?: string
}

export function CommandPalette({
  open,
  onOpenChange,
  items = [],
  groups: _groups,
  placeholder = "Type a command or search...",
}: CommandPaletteProps) {
  React.useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault()
        onOpenChange(!open)
      }
    }
    document.addEventListener("keydown", down)
    return () => document.removeEventListener("keydown", down)
  }, [open, onOpenChange])

  const groupedItems = React.useMemo(() => {
    const map = new Map<string, CommandItemType[]>()
    items.forEach((item) => {
      const g = item.group || "Actions"
      if (!map.has(g)) {
        map.set(g, [])
      }
      map.get(g)!.push(item)
    })
    return map
  }, [items])

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="overflow-hidden p-0 shadow-lg max-w-xl sm:rounded-lg border-border">
        <CommandPrimitive className="flex h-full w-full flex-col overflow-hidden rounded-md bg-popover text-popover-foreground">
          <div className="flex items-center border-b px-3">
            <IconSearch className="mr-2 h-4 w-4 shrink-0 text-muted-foreground" />
            <CommandPrimitive.Input
              placeholder={placeholder}
              className="flex h-11 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50"
            />
          </div>
          <CommandPrimitive.List className="max-h-[300px] overflow-y-auto overflow-x-hidden p-1.5">
            <CommandPrimitive.Empty className="py-6 text-center text-xs text-muted-foreground">
              No results found.
            </CommandPrimitive.Empty>

            {Array.from(groupedItems.entries()).map(([groupName, groupItems]) => (
              <CommandPrimitive.Group
                key={groupName}
                heading={groupName}
                className="overflow-hidden p-1 text-muted-foreground [&_[cmdk-group-heading]]:px-2 [&_[cmdk-group-heading]]:py-1.5 [&_[cmdk-group-heading]]:text-[10px] [&_[cmdk-group-heading]]:font-semibold [&_[cmdk-group-heading]]:uppercase [&_[cmdk-group-heading]]:tracking-wider"
              >
                {groupItems.map((item) => {
                  const Icon = item.icon
                  return (
                    <CommandPrimitive.Item
                      key={item.id}
                      value={`${item.label} ${item.keywords?.join(" ") || ""}`}
                      onSelect={() => {
                        item.onSelect?.()
                        onOpenChange(false)
                      }}
                      className="relative flex cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-xs outline-none aria-selected:bg-accent aria-selected:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
                    >
                      {Icon && <Icon className="mr-2 h-3.5 w-3.5 text-muted-foreground" />}
                      <span className="flex-1">{item.label}</span>
                      {item.shortcut && (
                        <div className="ml-auto flex items-center gap-1">
                          {item.shortcut.map((key) => (
                            <kbd
                              key={key}
                              className="pointer-events-none inline-flex h-4 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[9px] font-medium text-muted-foreground opacity-100"
                            >
                              {key}
                            </kbd>
                          ))}
                        </div>
                      )}
                    </CommandPrimitive.Item>
                  )
                })}
              </CommandPrimitive.Group>
            ))}
          </CommandPrimitive.List>
        </CommandPrimitive>
      </DialogContent>
    </Dialog>
  )
}
