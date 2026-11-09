import * as React from 'react'
import {
  IconPlus,
  IconMinus,
  IconBriefcase,
  IconFolderPlus,
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

export interface ServiceItem {
  id: string
  groupName: string
  serviceName: string
  quantity: number
  amount: number
}

const INITIAL_SERVICE_GROUPS = [
  'Basic Service',
  'Advance Service',
  'Other Service',
  'Delivery',
  'Offers',
  'Add-ons',
]

const INITIAL_SERVICES: ServiceItem[] = [
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
  const [serviceGroups, setServiceGroups] = React.useState<string[]>(INITIAL_SERVICE_GROUPS)
  const [services, setServices] = React.useState<ServiceItem[]>(INITIAL_SERVICES)

  // Manage Service Groups Modal State (matching Image 2)
  const [isManageGroupsOpen, setIsManageGroupsOpen] = React.useState(false)
  const [editingGroups, setEditingGroups] = React.useState<string[]>([])
  const [newGroupNameInput, setNewGroupNameInput] = React.useState('')

  // Quick Inline Add Service Row State
  const [newRowGroup, setNewRowGroup] = React.useState<string>('Basic Service')
  const [newRowName, setNewRowName] = React.useState('')
  const [newRowQty, setNewRowQty] = React.useState<number>(1)
  const [newRowAmount, setNewRowAmount] = React.useState<number>(0)

  // Open modal handler
  const handleOpenManageGroups = () => {
    setEditingGroups([...serviceGroups])
    setNewGroupNameInput('')
    setIsManageGroupsOpen(true)
  }

  // Add group inside modal
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

  // Remove group inside modal
  const handleRemoveGroupFromModal = (indexToRemove: number) => {
    setEditingGroups((prev) => prev.filter((_, idx) => idx !== indexToRemove))
  }

  // Save groups from modal
  const handleSaveServiceGroups = () => {
    if (editingGroups.length === 0) {
      toast.error('At least one service group is required')
      return
    }
    setServiceGroups(editingGroups)
    setIsManageGroupsOpen(false)
    toast.success('Service groups updated successfully')
  }

  // Inline row updates
  const handleUpdateServiceField = (
    id: string,
    field: keyof ServiceItem,
    val: any
  ) => {
    setServices((prev) =>
      prev.map((s) => (s.id === id ? { ...s, [field]: val } : s))
    )
  }

  // Remove service row
  const handleRemoveService = (id: string) => {
    setServices((prev) => prev.filter((s) => s.id !== id))
    toast.info('Service removed')
  }

  // Add new service row
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

    setServices((prev) => [...prev, newService])
    setNewRowName('')
    setNewRowQty(1)
    setNewRowAmount(0)
    toast.success(`Service "${newService.serviceName}" added`)
  }

  return (
    <div className="space-y-6 max-w-5xl mx-auto">
      {/* Title & Description matching Image 1 */}
      <div className="text-center space-y-1.5 pt-2">
        <h1 className="text-2xl font-bold tracking-tight text-foreground">Services</h1>
        <p className="text-xs text-muted-foreground max-w-xl mx-auto">
          Preset event services names with service cost and reduce your quote, invoice generation time.
        </p>

        {/* Manage Service Group Action Button */}
        <div className="pt-2">
          <Button
            variant="outline"
            size="sm"
            onClick={handleOpenManageGroups}
            className="text-xs font-semibold text-emerald-600 dark:text-emerald-400 hover:text-emerald-700 hover:bg-emerald-50 dark:hover:bg-emerald-950/30 border-emerald-500/30 gap-1.5 h-8 px-3.5 shadow-2xs"
          >
            <IconFolderPlus className="h-3.5 w-3.5" />
            <span>Manage Service Group</span>
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

