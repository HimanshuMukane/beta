import * as React from 'react'
import {
  IconPlus,
  IconMinus,
  IconBriefcase,
  IconFolderPlus,
  IconDeviceFloppy,
} from '@tabler/icons-react'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  Button,
  Input,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
  toast,
} from '@repo/ui'
import { settingsApi } from '../../api/settingsApi'

export interface ServiceItem {
  id: string
  groupName: string
  serviceName: string
  quantity: number
  amount: number
}

const DEFAULT_SERVICE_GROUPS = [
  'Basic Service',
  'Advance Service',
  'Other Service',
  'Delivery',
  'Offers',
  'Add-ons',
]

const DEFAULT_SERVICES: ServiceItem[] = [
  { id: 'srv-1', groupName: 'Basic Service', serviceName: 'Traditional Photography', quantity: 1, amount: 200 },
  { id: 'srv-2', groupName: 'Basic Service', serviceName: 'Traditional Videography', quantity: 1, amount: 200 },
  { id: 'srv-3', groupName: 'Advance Service', serviceName: 'Candid Photography', quantity: 1, amount: 150 },
  { id: 'srv-4', groupName: 'Advance Service', serviceName: 'Candid Videography', quantity: 1, amount: 150 },
  { id: 'srv-5', groupName: 'Advance Service', serviceName: 'Cinematic Camera', quantity: 1, amount: 200 },
  { id: 'srv-6', groupName: 'Delivery', serviceName: 'Frame', quantity: 1, amount: 200 },
  { id: 'srv-7', groupName: 'Delivery', serviceName: 'Frame 12 X 10', quantity: 1, amount: 100 },
  { id: 'srv-8', groupName: 'Other Service', serviceName: 'Live Streaming', quantity: 1, amount: 150 },
  { id: 'srv-9', groupName: 'Delivery', serviceName: 'Album 250 Sheets', quantity: 1, amount: 100 },
  { id: 'srv-10', groupName: 'Other Service', serviceName: 'Live Streaming 4 hours', quantity: 1, amount: 500 },
]

export function ServicesSettingsPage() {
  const [serviceGroups, setServiceGroups] = React.useState<string[]>(DEFAULT_SERVICE_GROUPS)
  const [services, setServices] = React.useState<ServiceItem[]>(DEFAULT_SERVICES)

  React.useEffect(() => {
    settingsApi.getSettings().then((store) => {
      if (store.serviceGroups && store.serviceGroups.length > 0) {
        setServiceGroups(store.serviceGroups)
      }
      if (store.services && store.services.length > 0) {
        setServices(store.services)
      }
    }).catch(() => {})
  }, [])

  // Manage Service Groups Modal State (matching Image 2)
  const [isManageGroupsOpen, setIsManageGroupsOpen] = React.useState(false)
  const [editingGroups, setEditingGroups] = React.useState<string[]>([])
  const [newGroupNameInput, setNewGroupNameInput] = React.useState('')

  // Quick Inline Add Service Row State
  const [newRowGroup, setNewRowGroup] = React.useState<string>('Basic Service')
  const [newRowName, setNewRowName] = React.useState('')
  const [newRowQty, setNewRowQty] = React.useState<number>(1)
  const [newRowAmount, setNewRowAmount] = React.useState<number>(0)

  const handleOpenManageGroups = () => {
    setEditingGroups([...serviceGroups])
    setNewGroupNameInput('')
    setIsManageGroupsOpen(true)
  }

  const handleAddGroupToModal = () => {
    if (!newGroupNameInput.trim()) return
    const trimmed = newGroupNameInput.trim()
    if (editingGroups.includes(trimmed)) {
      toast.error('Group name already exists')
      return
    }
    setEditingGroups((prev) => [...prev, trimmed])
    setNewGroupNameInput('')
  }

  const handleRemoveGroupFromModal = (indexToRemove: number) => {
    setEditingGroups((prev) => prev.filter((_, idx) => idx !== indexToRemove))
  }

  const handleSaveServiceGroups = async () => {
    if (editingGroups.length === 0) {
      toast.error('At least one service group is required')
      return
    }
    setServiceGroups(editingGroups)
    setIsManageGroupsOpen(false)
    try {
      await settingsApi.updateSettings({ serviceGroups: editingGroups })
      toast.success('Service groups updated successfully')
    } catch {
      toast.error('Failed to save service groups to database')
    }
  }

  const handleSaveAllServices = async () => {
    try {
      await settingsApi.updateSettings({ services, serviceGroups })
      toast.success('Services saved successfully')
    } catch {
      toast.error('Failed to save services to database')
    }
  }

  const handleUpdateServiceField = (
    id: string,
    field: keyof ServiceItem,
    val: any
  ) => {
    setServices((prev) =>
      prev.map((s) => (s.id === id ? { ...s, [field]: val } : s))
    )
  }

  const handleRemoveService = (id: string) => {
    setServices((prev) => {
      const updated = prev.filter((s) => s.id !== id)
      settingsApi.updateSettings({ services: updated }).catch(() => {})
      return updated
    })
    toast.info('Service removed')
  }

  const handleAddServiceRow = () => {
    if (!newRowName.trim()) {
      toast.error('Please enter a Service Name')
      return
    }
    const newService: ServiceItem = {
      id: `srv-${Date.now()}`,
      groupName: newRowGroup || serviceGroups[0] || 'Basic Service',
      serviceName: newRowName.trim(),
      quantity: Math.max(1, newRowQty || 1),
      amount: Math.max(0, newRowAmount || 0),
    }

    const updated = [...services, newService]
    setServices(updated)
    settingsApi.updateSettings({ services: updated }).catch(() => {})
    setNewRowName('')
    setNewRowQty(1)
    setNewRowAmount(0)
    toast.success(`Service "${newService.serviceName}" added`)
  }

  return (
    <div className="space-y-6 max-w-5xl mx-auto pb-16">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b pb-4">
        <div className="flex items-center gap-3">
          <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary border border-primary/20">
            <IconBriefcase className="h-5 w-5" />
          </span>
          <div>
            <h1 className="text-xl font-bold tracking-tight text-foreground">Services Presets</h1>
            <p className="text-xs text-muted-foreground">
              Preset event services with standard rates to streamline quote and invoice generation.
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2.5 self-start sm:self-auto">
          <Button
            variant="outline"
            size="sm"
            onClick={handleOpenManageGroups}
            className="text-xs font-semibold gap-1.5 h-8 border-dashed"
          >
            <IconFolderPlus className="h-3.5 w-3.5 text-primary" />
            <span>Manage Service Groups ({serviceGroups.length})</span>
          </Button>
          <Button
            size="sm"
            onClick={handleSaveAllServices}
            className="text-xs font-semibold gap-1.5 h-8 bg-emerald-600 hover:bg-emerald-700 text-white shadow-xs"
          >
            <IconDeviceFloppy className="h-3.5 w-3.5" />
            <span>Save All Changes</span>
          </Button>
        </div>
      </div>

      {/* Main Services Table Card */}
      <Card className="border bg-card shadow-2xs overflow-hidden">
        <CardHeader className="py-3 px-4 border-b bg-muted/20 flex flex-row items-center justify-between">
          <div className="flex items-center gap-2">
            <IconBriefcase className="h-4 w-4 text-emerald-600" />
            <CardTitle className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
              Configured Studio Services ({services.length})
            </CardTitle>
          </div>
          <span className="text-[11px] text-muted-foreground font-medium">
            Changes auto-apply to future Quotations & Invoices
          </span>
        </CardHeader>

        <CardContent className="p-0">
          {/* Header Row */}
          <div className="grid grid-cols-12 gap-3 px-4 py-2.5 border-b bg-muted/40 text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">
            <div className="col-span-3">Group Name</div>
            <div className="col-span-4">Service Name</div>
            <div className="col-span-2 text-center">Quantity</div>
            <div className="col-span-2">Amount</div>
            <div className="col-span-1 text-center">Action</div>
          </div>

          {/* Existing Service Rows */}
          <div className="divide-y divide-border">
            {services.map((srv) => (
              <div
                key={srv.id}
                className="grid grid-cols-12 gap-3 px-4 py-2.5 items-center hover:bg-muted/10 transition"
              >
                {/* Group Name Select dropdown */}
                <div className="col-span-3">
                  <Select
                    value={srv.groupName}
                    onValueChange={(val) => handleUpdateServiceField(srv.id, 'groupName', val)}
                  >
                    <SelectTrigger className="h-9 text-xs bg-background border-border shadow-2xs">
                      <SelectValue placeholder="Select Group" />
                    </SelectTrigger>
                    <SelectContent>
                      {serviceGroups.map((grp) => (
                        <SelectItem key={grp} value={grp} className="text-xs">
                          {grp}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                {/* Service Name Input */}
                <div className="col-span-4">
                  <Input
                    value={srv.serviceName}
                    onChange={(e) => handleUpdateServiceField(srv.id, 'serviceName', e.target.value)}
                    placeholder="Enter service name"
                    className="h-9 text-xs bg-background border-border shadow-2xs font-medium"
                  />
                </div>

                {/* Quantity Input */}
                <div className="col-span-2">
                  <Input
                    type="number"
                    min={1}
                    value={srv.quantity}
                    onChange={(e) =>
                      handleUpdateServiceField(srv.id, 'quantity', parseInt(e.target.value, 10) || 1)
                    }
                    className="h-9 text-xs bg-background border-border text-center shadow-2xs tabular-nums"
                  />
                </div>

                {/* Amount with INR prefix */}
                <div className="col-span-2">
                  <div className="relative flex items-center">
                    <span className="absolute left-2.5 text-[11px] font-semibold text-muted-foreground pointer-events-none">
                      INR
                    </span>
                    <Input
                      type="number"
                      min={0}
                      value={srv.amount}
                      onChange={(e) =>
                        handleUpdateServiceField(srv.id, 'amount', parseFloat(e.target.value) || 0)
                      }
                      className="pl-10 h-9 text-xs bg-background border-border shadow-2xs font-semibold tabular-nums text-foreground"
                    />
                  </div>
                </div>

                {/* Delete button (circle minus matching Screenshot 1) */}
                <div className="col-span-1 flex justify-center">
                  <button
                    type="button"
                    onClick={() => handleRemoveService(srv.id)}
                    className="h-8 w-8 rounded-full flex items-center justify-center text-muted-foreground/60 hover:text-destructive hover:bg-destructive/10 transition"
                    title="Remove Service"
                  >
                    <IconMinus className="h-4 w-4 rounded-full border border-current p-0.5" />
                  </button>
                </div>
              </div>
            ))}

            {/* Quick Inline New Service Row (matching bottom ghost row in Screenshot 1) */}
            <div className="grid grid-cols-12 gap-3 px-4 py-3 items-center bg-muted/20 border-t border-dashed">
              <div className="col-span-3">
                <Select value={newRowGroup} onValueChange={setNewRowGroup}>
                  <SelectTrigger className="h-9 text-xs bg-background border-border shadow-2xs">
                    <SelectValue placeholder="Select Group" />
                  </SelectTrigger>
                  <SelectContent>
                    {serviceGroups.map((grp) => (
                      <SelectItem key={grp} value={grp} className="text-xs">
                        {grp}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="col-span-4">
                <Input
                  value={newRowName}
                  onChange={(e) => setNewRowName(e.target.value)}
                  placeholder="Service Name"
                  className="h-9 text-xs bg-background border-border shadow-2xs"
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') handleAddServiceRow()
                  }}
                />
              </div>

              <div className="col-span-2">
                <Input
                  type="number"
                  min={1}
                  value={newRowQty}
                  onChange={(e) => setNewRowQty(parseInt(e.target.value, 10) || 1)}
                  className="h-9 text-xs bg-background border-border text-center shadow-2xs tabular-nums"
                />
              </div>

              <div className="col-span-2">
                <div className="relative flex items-center">
                  <span className="absolute left-2.5 text-[11px] font-semibold text-muted-foreground pointer-events-none">
                    INR
                  </span>
                  <Input
                    type="number"
                    min={0}
                    value={newRowAmount}
                    onChange={(e) => setNewRowAmount(parseFloat(e.target.value) || 0)}
                    className="pl-10 h-9 text-xs bg-background border-border shadow-2xs tabular-nums"
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') handleAddServiceRow()
                    }}
                  />
                </div>
              </div>

              <div className="col-span-1 flex justify-center">
                <Button
                  size="sm"
                  onClick={handleAddServiceRow}
                  className="h-8 px-2.5 text-xs bg-emerald-600 hover:bg-emerald-700 text-white font-semibold gap-1 shadow-2xs"
                >
                  <IconPlus className="h-3.5 w-3.5" />
                  <span>Add</span>
                </Button>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Bottom helper action */}
      <div className="flex items-center justify-between text-xs text-muted-foreground px-1">
        <button
          type="button"
          onClick={handleAddServiceRow}
          className="text-xs font-semibold text-emerald-600 dark:text-emerald-400 hover:underline flex items-center gap-1"
        >
          <IconPlus className="h-3.5 w-3.5" />
          <span>Add Service</span>
        </button>
        <span className="text-[11px]">All values saved locally to active studio profile.</span>
      </div>

      {/* MANAGE SERVICE GROUPS DIALOG (matching Image 2) */}
      <Dialog open={isManageGroupsOpen} onOpenChange={setIsManageGroupsOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="text-xl font-bold text-foreground">
              Manage Service Groups
            </DialogTitle>
            <DialogDescription className="text-xs text-muted-foreground">
              Define the categories used to group your studio photography and videography services.
            </DialogDescription>
          </DialogHeader>

          {/* Group Items list */}
          <div className="space-y-2.5 py-3 max-h-[50vh] overflow-y-auto">
            {editingGroups.map((grpName, index) => (
              <div key={index} className="flex items-center gap-2">
                <Input
                  value={grpName}
                  onChange={(e) => {
                    const updated = [...editingGroups]
                    updated[index] = e.target.value
                    setEditingGroups(updated)
                  }}
                  className="h-10 text-xs bg-background border-border font-medium shadow-2xs"
                />
                <button
                  type="button"
                  onClick={() => handleRemoveGroupFromModal(index)}
                  className="h-8 w-8 shrink-0 rounded-full flex items-center justify-center text-muted-foreground/60 hover:text-destructive hover:bg-destructive/10 transition"
                  title="Remove Group"
                >
                  <IconMinus className="h-4 w-4 rounded-full border border-current p-0.5" />
                </button>
              </div>
            ))}

            {/* Add Group Name input field (matching bottom row in Image 2) */}
            <div className="flex items-center gap-2 pt-2">
              <Input
                value={newGroupNameInput}
                onChange={(e) => setNewGroupNameInput(e.target.value)}
                placeholder="Add-ons (e.g. Drone, Framing, Vanity)"
                className="h-10 text-xs bg-background/50 border-border"
                onKeyDown={(e) => {
                  if (e.key === 'Enter') handleAddGroupToModal()
                }}
              />
              <Button
                variant="outline"
                size="sm"
                onClick={handleAddGroupToModal}
                className="h-10 px-3 text-xs shrink-0 font-semibold text-emerald-600 border-emerald-500/30"
              >
                <IconPlus className="h-3.5 w-3.5 mr-1" />
                Add
              </Button>
            </div>
          </div>

          <DialogFooter className="border-t pt-3 flex items-center justify-between sm:justify-between w-full">
            <button
              type="button"
              onClick={handleAddGroupToModal}
              className="text-xs font-semibold text-emerald-600 dark:text-emerald-400 hover:underline flex items-center gap-1"
            >
              <IconPlus className="h-3.5 w-3.5" />
              <span>Add Group Name</span>
            </button>
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setIsManageGroupsOpen(false)}
                className="text-xs h-9 px-4 font-semibold"
              >
                Cancel
              </Button>
              <Button
                size="sm"
                onClick={handleSaveServiceGroups}
                className="text-xs h-9 px-5 font-semibold bg-emerald-600 hover:bg-emerald-700 text-white shadow-xs"
              >
                Save
              </Button>
            </div>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

