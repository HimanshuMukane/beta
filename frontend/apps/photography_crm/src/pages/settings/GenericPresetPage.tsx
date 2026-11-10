import * as React from 'react'
import {
  IconPlus,
  IconTrash,
} from '@tabler/icons-react'
import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import { toast } from '@repo/ui/components/ui/sonner'

interface PresetItem {
  id: string
  name: string
  rate?: number
  unit?: string
  description?: string
}

interface GenericPresetPageProps {
  title: string
  subtitle: string
  icon: React.ComponentType<{ className?: string }>
  initialItems: PresetItem[]
  hasPrice?: boolean
}

export function GenericPresetPage({
  title,
  subtitle,
  icon: Icon,
  initialItems,
  hasPrice = true,
}: GenericPresetPageProps) {
  const [items, setItems] = React.useState<PresetItem[]>(initialItems)
  const [newName, setNewName] = React.useState('')
  const [newRate, setNewRate] = React.useState<number>(0)
  const [newUnit, setNewUnit] = React.useState('per unit')

  const handleAddItem = (e: React.FormEvent) => {
    e.preventDefault()
    if (!newName.trim()) return
    const newItem: PresetItem = {
      id: `preset-${Date.now()}`,
      name: newName.trim(),
      rate: hasPrice ? newRate : undefined,
      unit: hasPrice ? newUnit : undefined,
    }
    setItems((prev) => [...prev, newItem])
    setNewName('')
    setNewRate(0)
    toast.success(`Added ${newItem.name}`)
  }

  const handleDeleteItem = (id: string) => {
    setItems((prev) => prev.filter((i) => i.id !== id))
    toast.success('Preset item removed')
  }

  return (
    <div className="space-y-6 max-w-4xl mx-auto pb-16">
      <div className="flex items-center justify-between border-b pb-4">
        <div className="flex items-center gap-3">
          <span className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 text-primary border border-primary/20">
            <Icon className="h-5 w-5" />
          </span>
          <div>
            <h1 className="text-xl font-bold tracking-tight text-foreground">{title}</h1>
            <p className="text-xs text-muted-foreground">{subtitle}</p>
          </div>
        </div>
      </div>

      {/* Add New Preset Row */}
      <form
        onSubmit={handleAddItem}
        className="bg-card p-4 rounded-xl border border-border/80 shadow-2xs flex flex-wrap items-center gap-3"
      >
        <Input
          value={newName}
          onChange={(e) => setNewName(e.target.value)}
          placeholder={`Add new ${title.toLowerCase()} item...`}
          className="h-9 text-xs flex-1 min-w-[200px]"
          required
        />
        {hasPrice && (
          <>
            <Input
              type="number"
              min="0"
              value={newRate}
              onChange={(e) => setNewRate(parseFloat(e.target.value) || 0)}
              placeholder="Base Price (₹)"
              className="h-9 text-xs w-28 text-right font-mono"
            />
            <Input
              value={newUnit}
              onChange={(e) => setNewUnit(e.target.value)}
              placeholder="Unit (e.g. per hr)"
              className="h-9 text-xs w-28"
            />
          </>
        )}
        <Button
          type="submit"
          size="sm"
          className="bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs h-9 px-4 gap-1.5 shadow-xs"
        >
          <IconPlus className="h-4 w-4" />
          Add Preset
        </Button>
      </form>

      {/* List of presets */}
      <div className="rounded-xl border border-border/80 overflow-hidden bg-card shadow-2xs">
        <div className="divide-y divide-border/60">
          {items.map((item) => (
            <div
              key={item.id}
              className="flex items-center justify-between p-3.5 hover:bg-muted/15 transition text-xs"
            >
              <div className="flex items-center gap-3">
                <span className="font-semibold text-foreground">{item.name}</span>
                {item.unit && (
                  <span className="text-[10px] text-muted-foreground">({item.unit})</span>
                )}
              </div>

              <div className="flex items-center gap-4">
                {item.rate !== undefined && (
                  <span className="font-bold text-xs text-foreground font-mono">
                    ₹{item.rate.toLocaleString('en-IN')}
                  </span>
                )}
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => handleDeleteItem(item.id)}
                  className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive"
                >
                  <IconTrash className="h-3.5 w-3.5" />
                </Button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
