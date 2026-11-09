import * as React from 'react'
import {
  IconPlus,
  IconTrash,
  IconShoppingBag,
  IconCheck,
  IconReceipt2,
  IconCoins,
} from '@tabler/icons-react'
import {
  Button,
  Input,
  Tooltip,
  TooltipTrigger,
  TooltipContent,
  toast,
} from '@repo/ui'
import { useStudioSettings } from '../../hooks/useStudioSettings'

export interface OrderPresetItem {
  id: string
  name: string
  quantity: number
  amount: number
  expense: number
}

const INITIAL_ORDERS_DATA: OrderPresetItem[] = [
  { id: 'ord-1', name: 'Passport Photo', quantity: 1, amount: 200, expense: 300 },
  { id: 'ord-2', name: 'Frame 10 * 12', quantity: 1, amount: 1200, expense: 500 },
  { id: 'ord-3', name: 'Gifts', quantity: 1, amount: 2500, expense: 800 },
  { id: 'ord-4', name: '12x36 Photobook (40 Pgs)', quantity: 1, amount: 12500, expense: 4500 },
  { id: 'ord-5', name: 'Acrylic Wall Frame (24x36)', quantity: 1, amount: 6500, expense: 2200 },
]

export function OrdersSettingsPage() {
  const { settings, updateSettings } = useStudioSettings()
  const [items, setItems] = React.useState<OrderPresetItem[]>(
    settings?.ordersPresets || INITIAL_ORDERS_DATA
  )

  React.useEffect(() => {
    if (settings?.ordersPresets) {
      setItems(settings.ordersPresets)
    }
  }, [settings?.ordersPresets])

  const handleUpdate = (id: string, field: keyof OrderPresetItem, value: any) => {
    setItems((prev) =>
      prev.map((item) => (item.id === id ? { ...item, [field]: value } : item))
    )
  }

  const handleAddItem = () => {
    const newItem: OrderPresetItem = {
      id: `ord-${Date.now()}`,
      name: '',
      quantity: 1,
      amount: 0,
      expense: 0,
    }
    setItems((prev) => [...prev, newItem])
  }

  const handleDeleteItem = (id: string) => {
    if (items.length <= 1) {
      toast.error('At least one line item is required')
      return
    }
    const updated = items.filter((item) => item.id !== id)
    setItems(updated)
    updateSettings({ ordersPresets: updated })
    toast.success('Line item removed')
  }

  const handleSaveAll = () => {
    const valid = items.filter((i) => i.name.trim().length > 0)
    if (valid.length === 0) {
      toast.error('Please enter at least one item name')
      return
    }
    setItems(valid)
    updateSettings({ ordersPresets: valid })
    toast.success('Orders preset configuration saved')
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6 pb-16">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b pb-4">
        <div className="flex items-center gap-3">
          <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary border border-primary/20">
            <IconShoppingBag className="h-5 w-5" />
          </span>
          <div>
            <h1 className="text-xl font-bold tracking-tight text-foreground">Orders Presets</h1>
            <p className="text-xs text-muted-foreground">
              Configure studio orders items with pricing and cost expenses for easy billing.
            </p>
          </div>
        </div>

        <Button
          type="button"
          onClick={handleSaveAll}
          className="h-8 text-xs font-semibold gap-1.5 bg-primary text-primary-foreground self-start sm:self-auto shadow-xs"
        >
          <IconCheck className="h-3.5 w-3.5" />
          Save Changes
        </Button>
      </div>

      <div className="bg-card border border-border/80 shadow-2xs rounded-xl overflow-hidden">
        <div className="p-3.5 border-b bg-muted/20 flex items-center justify-between">
          <span className="text-xs font-bold text-foreground">Catalog & Pricing Matrix</span>
          <span className="text-[11px] text-muted-foreground">
            {items.length} items configured
          </span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse text-xs">
            <thead>
              <tr className="border-b bg-muted/40 text-muted-foreground text-[11px] font-semibold uppercase tracking-wider">
                <th className="py-3 px-4 w-[38%]">Item Name</th>
                <th className="py-3 px-3 w-[14%] text-center">Quantity</th>
                <th className="py-3 px-4 w-[22%]">Amount (Billing)</th>
                <th className="py-3 px-4 w-[22%]">Expense (Cost)</th>
                <th className="py-3 px-3 w-[4%] text-center"></th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border/60">
              {items.map((item) => (
                <tr key={item.id} className="hover:bg-muted/10 transition-colors">
                  <td className="py-2.5 px-4">
                    <Input
                      value={item.name}
                      onChange={(e) => handleUpdate(item.id, 'name', e.target.value)}
                      placeholder="e.g. Passport Photo, Frame 10*12..."
                      className="h-8 text-xs bg-background"
                    />
                  </td>
                  <td className="py-2.5 px-3">
                    <Input
                      type="number"
                      min={1}
                      value={item.quantity}
                      onChange={(e) =>
                        handleUpdate(item.id, 'quantity', parseInt(e.target.value, 10) || 1)
                      }
                      className="h-8 text-xs text-center bg-background"
                    />
                  </td>
                  <td className="py-2.5 px-4">
                    <div className="relative flex items-center">
                      <span className="absolute left-2.5 text-[11px] font-semibold text-muted-foreground">
                        INR
                      </span>
                      <Input
                        type="number"
                        min={0}
                        value={item.amount || ''}
                        onChange={(e) =>
                          handleUpdate(item.id, 'amount', parseFloat(e.target.value) || 0)
                        }
                        placeholder="Amount"
                        className="h-8 text-xs pl-12 bg-background font-mono"
                      />
                    </div>
                  </td>
                  <td className="py-2.5 px-4">
                    <div className="relative flex items-center">
                      <span className="absolute left-2.5 text-[11px] font-semibold text-muted-foreground">
                        INR
                      </span>
                      <Input
                        type="number"
                        min={0}
                        value={item.expense || ''}
                        onChange={(e) =>
                          handleUpdate(item.id, 'expense', parseFloat(e.target.value) || 0)
                        }
                        placeholder="Item Expense"
                        className="h-8 text-xs pl-12 bg-background font-mono"
                      />
                    </div>
                  </td>
                  <td className="py-2.5 px-3 text-center">
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          onClick={() => handleDeleteItem(item.id)}
                          className="h-7 w-7 text-muted-foreground hover:text-destructive hover:bg-destructive/10"
                        >
                          <IconTrash className="h-3.5 w-3.5" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent side="left">
                        Remove line item
                      </TooltipContent>
                    </Tooltip>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div className="p-3 border-t bg-muted/10 flex items-center justify-between">
          <Button
            type="button"
            variant="ghost"
            onClick={handleAddItem}
            className="text-xs text-primary font-semibold hover:bg-primary/10 gap-1.5 h-8 px-2"
          >
            <IconPlus className="h-3.5 w-3.5" />
            Add Item
          </Button>

          <div className="flex items-center gap-4 text-xs text-muted-foreground font-medium">
            <span className="flex items-center gap-1">
              <IconReceipt2 className="h-3.5 w-3.5 text-primary" />
              Total Items: <strong className="text-foreground">{items.length}</strong>
            </span>
            <Tooltip>
              <TooltipTrigger asChild>
                <span className="flex items-center gap-1 cursor-help underline decoration-dotted decoration-muted-foreground/50 underline-offset-4">
                  <IconCoins className="h-3.5 w-3.5 text-emerald-600" />
                  Avg Margin:{' '}
                  <strong className="text-foreground">
                    {items.length > 0
                      ? Math.round(
                          (items.reduce((acc, i) => acc + (i.amount - i.expense), 0) /
                            (items.reduce((acc, i) => acc + (i.amount || 1), 0) || 1)) *
                            100
                        )
                      : 0}
                    %
                  </strong>
                </span>
              </TooltipTrigger>
              <TooltipContent side="top">
                Net profit margin = (Billing Amount - Expense Cost) / Billing Amount
              </TooltipContent>
            </Tooltip>
          </div>
        </div>
      </div>
    </div>
  )
}

