import {
  ComboBox as AriaComboBox,
  Input as AriaInput,
  Button as AriaButton,
  ListBox as AriaListBox,
  ListBoxItem as AriaListBoxItem,
  Popover as AriaPopover,
  type Key,
} from "react-aria-components"
import { IconChevronDown, IconCheck, IconX } from "@tabler/icons-react"
import { cn } from "@repo/ui/lib/utils"

export interface ComboboxOption {
  id: string
  label: string
  description?: string
  disabled?: boolean
}

export interface ComboboxProps {
  items: ComboboxOption[]
  selectedKey?: string | null
  onSelectionChange?: (key: string | null) => void
  placeholder?: string
  label?: string
  disabled?: boolean
  className?: string
  allowsCustomValue?: boolean
}

export function Combobox({
  items,
  selectedKey,
  onSelectionChange,
  placeholder = "Select or search...",
  label,
  disabled = false,
  className,
  allowsCustomValue = false,
}: ComboboxProps) {
  return (
    <AriaComboBox
      aria-label={label || placeholder}
      isDisabled={disabled}
      selectedKey={selectedKey}
      onSelectionChange={(key: Key | null) => {
        onSelectionChange?.(key ? String(key) : null)
      }}
      allowsCustomValue={allowsCustomValue}
      className={cn("relative w-full text-xs", className)}
    >
      <div className="relative flex items-center w-full">
        <AriaInput
          placeholder={placeholder}
          className="flex h-8 w-full rounded-md border border-input bg-background pl-2.5 pr-8 py-1 text-xs ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
        />
        <div className="absolute right-1 flex items-center">
          {selectedKey && (
            <button
              type="button"
              onClick={() => onSelectionChange?.(null)}
              className="p-1 text-muted-foreground hover:text-foreground rounded-xs"
            >
              <IconX className="h-3 w-3" />
            </button>
          )}
          <AriaButton className="p-1 text-muted-foreground hover:text-foreground rounded-xs outline-none">
            <IconChevronDown className="h-3.5 w-3.5 opacity-60" />
          </AriaButton>
        </div>
      </div>

      <AriaPopover className="z-50 min-w-[var(--trigger-width)] max-h-60 overflow-y-auto rounded-md border bg-popover p-1 text-popover-foreground shadow-md outline-none">
        <AriaListBox items={items} className="outline-none space-y-0.5">
          {(item) => (
            <AriaListBoxItem
              id={item.id}
              textValue={item.label}
              isDisabled={item.disabled}
              className={({ isFocused, isSelected, isDisabled }) =>
                cn(
                  "relative flex cursor-pointer select-none items-center justify-between rounded-sm px-2 py-1.5 text-xs outline-none",
                  isFocused && "bg-accent text-accent-foreground",
                  isSelected && "font-medium text-primary",
                  isDisabled && "pointer-events-none opacity-50"
                )
              }
            >
              {({ isSelected }) => (
                <>
                  <div className="flex flex-col">
                    <span>{item.label}</span>
                    {item.description && (
                      <span className="text-[10px] text-muted-foreground">
                        {item.description}
                      </span>
                    )}
                  </div>
                  {isSelected && <IconCheck className="h-3.5 w-3.5 text-primary" />}
                </>
              )}
            </AriaListBoxItem>
          )}
        </AriaListBox>
      </AriaPopover>
    </AriaComboBox>
  )
}
