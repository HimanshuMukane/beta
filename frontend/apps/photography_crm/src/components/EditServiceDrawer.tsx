import * as React from 'react'
import {
  IconPlus,
  IconTrash,
  IconCheck,
  IconInfoCircle,
  IconClock,
} from '@tabler/icons-react'
import { Sheet, SheetContent } from '@repo/ui/components/ui/sheet'
import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@repo/ui/components/ui/select'
import { cn } from '@repo/ui/lib/utils'
import type { ServiceItem } from '../types/lead'

interface EditServiceDrawerProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  services: ServiceItem[]
  onSave: (updated: ServiceItem[]) => void
}

const COMMON_SERVICE_GROUPS = ['Wedding', 'Pre-Wedding', 'Reception', 'Engagement', 'Haldi/Mehendi', 'Corporate']

export function EditServiceDrawer({
  open,
  onOpenChange,
  services,
  onSave,
}: EditServiceDrawerProps) {
  const [draftServices, setDraftServices] = React.useState<ServiceItem[]>(services)

  React.useEffect(() => {
    if (open) {
      setDraftServices(JSON.parse(JSON.stringify(services)))
    }
  }, [open, services])

  const subTotal = React.useMemo(() => {
    return draftServices.reduce((acc, s) => acc + (Number(s.subTotal) || 0), 0)
  }, [draftServices])

  const handleUpdate = (id: string, field: keyof ServiceItem, value: any) => {
    setDraftServices((prev) =>
      prev.map((s) => {
        if (s.id !== id) return s
        const next = { ...s, [field]: value }
        if (field === 'qty' || field === 'cost') {
          const qty = field === 'qty' ? Number(value) || 0 : s.qty
          const cost = field === 'cost' ? Number(value) || 0 : s.cost
          next.subTotal = qty * cost
        }
        return next
      })
    )
  }

  const handleRemove = (id: string) => {
    setDraftServices((prev) => prev.filter((s) => s.id !== id))
  }

  const handleAddService = (groupName = 'Wedding') => {
    const newItem: ServiceItem = {
      id: `srv-${Date.now()}-${Math.random().toString(36).substring(2, 5)}`,
      name: 'Photography Coverage',
      groupName,
      qty: 1,
      cost: 25000,
      subTotal: 25000,
      sessionPlanned: false,
    }
    setDraftServices((prev) => [...prev, newItem])
  }

  const handleAddCustom = (groupName = 'Wedding') => {
    const newItem: ServiceItem = {
      id: `srv-${Date.now()}-${Math.random().toString(36).substring(2, 5)}`,
      name: '',
      groupName,
      qty: 1,
      cost: 0,
      subTotal: 0,
      sessionPlanned: false,
    }
    setDraftServices((prev) => [...prev, newItem])
  }

  const handleAddGroup = () => {
    handleAddService('Custom Service Group')
  }

  const handleSave = () => {
    onSave(draftServices)
    onOpenChange(false)
  }

  // Group services by groupName
  const grouped = React.useMemo(() => {
    const map = new Map<string, ServiceItem[]>()
    draftServices.forEach((s) => {
      const g = s.groupName || 'General Services'
      if (!map.has(g)) map.set(g, [])
      map.get(g)!.push(s)
    })
    return Array.from(map.entries())
  }, [draftServices])

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent
        side="right"
        className="w-full sm:max-w-2xl lg:max-w-3xl flex flex-col p-0 bg-background"
      >
        {/* Header */}
        <div className="flex items-center justify-between border-b px-6 py-4 bg-muted/20">
          <div>
            <h2 className="text-lg font-semibold tracking-tight text-foreground">
              Edit Lead Service
            </h2>
            <p className="text-xs text-muted-foreground mt-0.5">
              Customize photography packages, hourly rates, and planning sessions.
            </p>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-xs text-muted-foreground font-medium">Sub Total:</span>
            <span className="text-lg font-bold text-foreground">
              ₹{subTotal.toLocaleString('en-IN')}
            </span>
          </div>
        </div>

        {/* Scrollable Content */}
        <div className="flex-1 overflow-y-auto px-6 py-5 space-y-6">
          {grouped.length === 0 ? (
            <div className="flex flex-col items-center justify-center rounded-xl border border-dashed py-12 text-center">
              <p className="text-sm font-medium text-muted-foreground">
                No services added yet.
              </p>
              <Button
                size="sm"
                variant="outline"
                onClick={() => handleAddService('Wedding')}
                className="mt-3"
              >
                <IconPlus className="mr-1.5 h-4 w-4" /> Add First Service
              </Button>
            </div>
          ) : (
            grouped.map(([groupName, items]) => (
              <div
                key={groupName}
                className="rounded-xl border bg-card/60 shadow-xs overflow-hidden"
              >
                {/* Table Header inside group */}
                <div className="grid grid-cols-12 gap-3 bg-muted/40 px-4 py-2.5 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                  <div className="col-span-4">Service Group & Services</div>
                  <div className="col-span-2 text-center">Qty</div>
                  <div className="col-span-2 text-right">Service Cost</div>
                  <div className="col-span-2 text-right">Sub Total</div>
                  <div className="col-span-2 text-center">Actions</div>
                </div>

                {/* Service rows */}
                <div className="divide-y divide-border/60">
                  {items.map((service) => (
                    <div
                      key={service.id}
                      className="grid grid-cols-12 gap-3 items-center px-4 py-3 hover:bg-muted/10 transition"
                    >
                      {/* Name & Group */}
                      <div className="col-span-4 space-y-1">
                        <Input
                          value={service.name}
                          onChange={(e) => handleUpdate(service.id, 'name', e.target.value)}
                          placeholder="Service name (e.g. Candid Team)"
                          className="h-8 text-xs font-medium"
                        />
                        <div className="flex items-center gap-1.5">
                          <Select
                            value={service.groupName}
                            onValueChange={(val) =>
                              handleUpdate(service.id, 'groupName', val)
                            }
                          >
                            <SelectTrigger className="h-6 text-[11px] py-0 px-2 text-muted-foreground bg-muted/50 border-0">
                              <SelectValue placeholder="Group" />
                            </SelectTrigger>
                            <SelectContent>
                              {COMMON_SERVICE_GROUPS.map((cg) => (
                                <SelectItem key={cg} value={cg} className="text-xs">
                                  {cg}
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                        </div>
                      </div>

                      {/* Qty */}
                      <div className="col-span-2">
                        <Input
                          type="number"
                          min="1"
                          value={service.qty}
                          onChange={(e) =>
                            handleUpdate(service.id, 'qty', parseInt(e.target.value) || 0)
                          }
                          className="h-8 text-center text-xs tabular-nums"
                        />
                      </div>

                      {/* Service Cost */}
                      <div className="col-span-2">
                        <div className="relative">
                          <span className="absolute left-2 top-2 text-[10px] text-muted-foreground">
                            INR
                          </span>
                          <Input
                            type="number"
                            min="0"
                            step="500"
                            value={service.cost}
                            onChange={(e) =>
                              handleUpdate(service.id, 'cost', parseFloat(e.target.value) || 0)
                            }
                            className="h-8 pl-8 text-right text-xs tabular-nums"
                          />
                        </div>
                      </div>

                      {/* Sub Total */}
                      <div className="col-span-2 text-right">
                        <span className="text-xs font-semibold tabular-nums text-foreground">
                          ₹{service.subTotal.toLocaleString('en-IN')}
                        </span>
                      </div>

                      {/* Actions */}
                      <div className="col-span-2 flex items-center justify-center gap-1">
                        <Button
                          type="button"
                          variant={service.sessionPlanned ? 'default' : 'outline'}
                          size="sm"
                          onClick={() =>
                            handleUpdate(service.id, 'sessionPlanned', !service.sessionPlanned)
                          }
                          className={cn(
                            'h-7 px-2 text-[11px] gap-1',
                            service.sessionPlanned
                              ? 'bg-emerald-600 hover:bg-emerald-700 text-white'
                              : 'text-muted-foreground'
                          )}
                          title="Toggle Session Planning"
                        >
                          <IconCheck className="h-3 w-3" />
                          <span className="hidden sm:inline">Session</span>
                        </Button>
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => handleRemove(service.id)}
                          className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive"
                        >
                          <IconTrash className="h-3.5 w-3.5" />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>

                {/* Add Service within this group */}
                <div className="flex items-center gap-2 border-t bg-muted/20 px-4 py-2">
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => handleAddService(groupName)}
                    className="h-7 text-xs text-primary hover:text-primary gap-1"
                  >
                    <IconPlus className="h-3.5 w-3.5" /> Add Service
                  </Button>
                  <span className="text-muted-foreground text-xs">•</span>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => handleAddCustom(groupName)}
                    className="h-7 text-xs text-muted-foreground hover:text-foreground gap-1"
                  >
                    <IconPlus className="h-3.5 w-3.5" /> Add Custom
                  </Button>
                </div>
              </div>
            ))
          )}

          {/* Add Service Group Button */}
          <Button
            type="button"
            variant="outline"
            onClick={handleAddGroup}
            className="w-full border-dashed py-5 text-xs text-muted-foreground hover:border-primary hover:text-primary"
          >
            <IconPlus className="mr-1.5 h-4 w-4" /> Add Service Group
          </Button>

          {/* Informational Callouts as per UI specs */}
          <div className="space-y-2.5 pt-2">
            <div className="flex items-start gap-2.5 rounded-lg border border-info/30 bg-info/10 p-3 text-xs text-foreground">
              <IconInfoCircle className="h-4 w-4 text-info shrink-0 mt-0.5" />
              <div>
                <span className="font-semibold text-info">Service Info: </span>
                You can hide service-level pricing (Qty & Amount) in the client PDF or share link using your quotation template settings.
              </div>
            </div>

            <div className="flex items-start gap-2.5 rounded-lg border border-blue-500/30 bg-blue-500/10 p-3 text-xs text-foreground">
              <IconClock className="h-4 w-4 text-blue-600 dark:text-blue-400 shrink-0 mt-0.5" />
              <div>
                <span className="font-semibold text-blue-600 dark:text-blue-400">Session Planning Tip: </span>
                Plan your session while sharing the quotation and avoid entering the same details again during booking. If you leave the date empty, the session will be created as "TBD" (To Be Decided) in the calendar.
              </div>
            </div>
          </div>
        </div>

        {/* Footer Actions */}
        <div className="flex items-center justify-end gap-3 border-t bg-muted/20 px-6 py-4">
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            onClick={handleSave}
            className="bg-emerald-600 hover:bg-emerald-700 text-white font-medium px-5"
          >
            Update
          </Button>
        </div>
      </SheetContent>
    </Sheet>
  )
}
