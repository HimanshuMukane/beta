import * as React from 'react'
import {
  IconPlus,
  IconMinus,
  IconTrash,
  IconEdit,
  IconSparkles,
  IconBox,
  IconArrowLeft,
  IconCopy,
  IconTrendingUp,
  IconChevronDown,
} from '@tabler/icons-react'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
  Button,
  Input,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Textarea,
  Badge,
  toast,
  cn,
} from '@repo/ui'
import { DeliverablesDropdown } from '../../components/DeliverablesDropdown'

export interface PackageServiceItem {
  id: string
  name: string
  cost: number
  quantity: number
  isCustom?: boolean
}

export interface PackageServiceGroup {
  id: string
  groupName: string
  quantity: number
  services: PackageServiceItem[]
}

export interface StudioPackage {
  id: string
  name: string
  description: string
  totalCost: number
  serviceGroups: PackageServiceGroup[]
  deliverables: string[]
}

const DEFAULT_SERVICE_GROUPS = [
  'Wedding',
  'Pre-Wedding',
  'Engagement',
  'Basic Service',
  'Advance Service',
  'Delivery',
]

const SAMPLE_PACKAGES: StudioPackage[] = [
  {
    id: 'pkg-1',
    name: 'Sample Wedding Package',
    description: 'Complete 2-day traditional and candid wedding coverage with luxury albums.',
    totalCost: 60000,
    serviceGroups: [
      {
        id: 'grp-1',
        groupName: 'Wedding',
        quantity: 1,
        services: [
          { id: 's-1', name: 'Traditional Photography', cost: 200, quantity: 1 },
          { id: 's-2', name: 'Traditional Videography', cost: 200, quantity: 1, isCustom: true },
        ],
      },
    ],
    deliverables: ['Raw Photos', 'Pen Drive', 'Handcrafted Album'],
  },
]

export function PackagesSettingsPage() {
  const [packages, setPackages] = React.useState<StudioPackage[]>(SAMPLE_PACKAGES)
  const [viewMode, setViewMode] = React.useState<'list' | 'editor'>('list')
  const [editingPackageId, setEditingPackageId] = React.useState<string | null>(null)

  // Editor Form State
  const [packageName, setPackageName] = React.useState('')
  const [packageDescription, setPackageDescription] = React.useState('')
  const [serviceGroups, setServiceGroups] = React.useState<PackageServiceGroup[]>([])
  const [selectedDeliverables, setSelectedDeliverables] = React.useState<string[]>([])

  // Open editor for new package
  const handleOpenNewPackage = () => {
    setEditingPackageId(null)
    setPackageName('')
    setPackageDescription('')
    setServiceGroups([
      {
        id: `grp-${Date.now()}`,
        groupName: 'Wedding',
        quantity: 1,
        services: [
          { id: `srv-${Date.now()}-1`, name: 'Traditional Photography', cost: 200, quantity: 1 },
        ],
      },
    ])
    setSelectedDeliverables(['Raw Photos', 'Pen Drive'])
    setViewMode('editor')
  }

  // Open editor for existing package
  const handleManagePackage = (pkg: StudioPackage) => {
    setEditingPackageId(pkg.id)
    setPackageName(pkg.name)
    setPackageDescription(pkg.description)
    setServiceGroups(pkg.serviceGroups || [])
    setSelectedDeliverables(pkg.deliverables || [])
    setViewMode('editor')
  }

  // Delete package
  const handleDeletePackage = (id: string) => {
    setPackages((prev) => prev.filter((p) => p.id !== id))
    toast.success('Package deleted')
  }

  // Duplicate package
  const handleDuplicatePackage = (pkg: StudioPackage) => {
    const duplicated: StudioPackage = {
      ...pkg,
      id: `pkg-${Date.now()}`,
      name: `${pkg.name} (Copy)`,
    }
    setPackages((prev) => [duplicated, ...prev])
    toast.success('Package duplicated')
  }

  // Add new Service Group in editor
  const handleAddServiceGroup = () => {
    const newGroup: PackageServiceGroup = {
      id: `grp-${Date.now()}`,
      groupName: 'Wedding',
      quantity: 1,
      services: [],
    }
    setServiceGroups((prev) => [...prev, newGroup])
  }

  // Remove Service Group
  const handleRemoveServiceGroup = (groupId: string) => {
    setServiceGroups((prev) => prev.filter((g) => g.id !== groupId))
  }

  // Add service inside group
  const handleAddServiceToGroup = (groupId: string, isCustom = false) => {
    const newService: PackageServiceItem = {
      id: `srv-${Date.now()}`,
      name: isCustom ? 'Custom Service' : 'Candid Photography',
      cost: 200,
      quantity: 1,
      isCustom,
    }
    setServiceGroups((prev) =>
      prev.map((g) =>
        g.id === groupId ? { ...g, services: [...g.services, newService] } : g
      )
    )
  }

  // Update service quantity
  const handleUpdateServiceQty = (groupId: string, serviceId: string, delta: number) => {
    setServiceGroups((prev) =>
      prev.map((g) => {
        if (g.id !== groupId) return g
        return {
          ...g,
          services: g.services.map((s) => {
            if (s.id !== serviceId) return s
            const nextQty = Math.max(1, s.quantity + delta)
            return { ...s, quantity: nextQty }
          }),
        }
      })
    )
  }

  // Update service group fields
  const handleUpdateGroup = (groupId: string, field: keyof PackageServiceGroup, val: any) => {
    setServiceGroups((prev) =>
      prev.map((g) => (g.id === groupId ? { ...g, [field]: val } : g))
    )
  }

  // Calculate total package cost dynamically
  const totalCalculatedCost = React.useMemo(() => {
    return serviceGroups.reduce((acc, grp) => {
      const groupSubtotal = grp.services.reduce(
        (sAcc, s) => sAcc + (s.cost || 0) * (s.quantity || 1),
        0
      )
      return acc + groupSubtotal * (grp.quantity || 1)
    }, 0)
  }, [serviceGroups])

  // Save package handler
  const handleSavePackage = () => {
    if (!packageName.trim()) {
      toast.error('Package Name is required')
      return
    }

    const saved: StudioPackage = {
      id: editingPackageId || `pkg-${Date.now()}`,
      name: packageName.trim(),
      description: packageDescription.trim(),
      totalCost: totalCalculatedCost || 60000,
      serviceGroups,
      deliverables: selectedDeliverables,
    }

    setPackages((prev) => {
      const exists = prev.some((p) => p.id === saved.id)
      return exists ? prev.map((p) => (p.id === saved.id ? saved : p)) : [saved, ...prev]
    })

    toast.success(editingPackageId ? 'Package updated!' : 'Package created!')
    setViewMode('list')
  }

  // 1. LIST VIEW (matching Image 4)
  if (viewMode === 'list') {
    return (
      <div className="space-y-6 max-w-4xl mx-auto pt-2">
        {/* Title & Subtitle */}
        <div className="text-center space-y-1.5">
          <h1 className="text-2xl font-bold tracking-tight text-foreground">Packages</h1>
          <p className="text-xs text-muted-foreground max-w-xl mx-auto">
            Add complete package details, including the package name and service descriptions. Easily use the package name in leads and invoices.
          </p>
        </div>

        {/* Limit badge & New Package button */}
        <div className="flex flex-col items-center justify-center gap-2 pt-2">
          <Badge
            variant="secondary"
            className="text-[11px] font-semibold bg-purple-500/15 text-purple-700 dark:text-purple-300 px-3 py-1"
          >
            {packages.length} Created / 3 Limit
          </Badge>
          <button
            type="button"
            onClick={handleOpenNewPackage}
            className="text-xs font-bold text-emerald-600 dark:text-emerald-400 hover:underline inline-flex items-center gap-1 transition"
          >
            <IconPlus className="h-3.5 w-3.5" />
            <span>New Package</span>
          </button>
        </div>

        {/* Package Cards List */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-4">
          {packages.map((pkg) => (
            <Card
              key={pkg.id}
              className="border border-border bg-card p-5 rounded-2xl shadow-2xs hover:shadow-xs transition space-y-4"
            >
              <div className="flex items-start justify-between gap-2">
                <div>
                  <h3 className="text-base font-bold text-foreground">{pkg.name}</h3>
                  <p className="text-xs text-muted-foreground mt-0.5 line-clamp-2">
                    {pkg.description || 'Custom bundled studio offering.'}
                  </p>
                </div>
                <IconTrendingUp className="h-5 w-5 text-emerald-600 shrink-0" />
              </div>

              <div className="flex items-center justify-between pt-2 border-t border-dashed">
                <span className="text-xl font-bold tabular-nums text-foreground">
                  ₹{pkg.totalCost.toLocaleString('en-IN')}
                </span>

                <div className="flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handleManagePackage(pkg)}
                    className="h-8 px-4 text-xs font-semibold rounded-full border-border hover:bg-muted"
                  >
                    Manage
                  </Button>
                  <button
                    type="button"
                    onClick={() => handleDeletePackage(pkg.id)}
                    className="h-8 w-8 rounded-full border border-rose-300 text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/30 flex items-center justify-center transition"
                    title="Delete Package"
                  >
                    <IconTrash className="h-3.5 w-3.5" />
                  </button>
                </div>
              </div>
            </Card>
          ))}
        </div>
      </div>
    )
  }

  // 2. EDITOR VIEW (matching Image 3)
  return (
    <div className="space-y-6 max-w-4xl mx-auto pt-2 pb-12">
      {/* Back button */}
      <div className="flex items-center justify-between">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setViewMode('list')}
          className="text-xs font-semibold gap-1 text-muted-foreground hover:text-foreground"
        >
          <IconArrowLeft className="h-4 w-4" />
          <span>Back to Packages</span>
        </Button>
        <Button
          size="sm"
          onClick={handleSavePackage}
          className="bg-emerald-600 hover:bg-emerald-700 text-white text-xs h-8 px-5 font-semibold shadow-xs"
        >
          Save Package
        </Button>
      </div>

      <Card className="border bg-card p-6 rounded-2xl shadow-2xs space-y-6">
        {/* Package Name */}
        <div className="space-y-1.5">
          <label className="text-xs font-bold text-foreground">Package Name</label>
          <Input
            value={packageName}
            onChange={(e) => setPackageName(e.target.value)}
            placeholder="Enter package name (e.g. Royal Wedding Bundle)"
            className="h-10 text-xs bg-background/50 border-input"
          />
        </div>

        {/* Package Description */}
        <div className="space-y-1.5">
          <label className="text-xs font-bold text-foreground">Package Description</label>
          <Textarea
            value={packageDescription}
            onChange={(e) => setPackageDescription(e.target.value)}
            placeholder="Enter package details and inclusions..."
            className="min-h-[85px] text-xs bg-background/50 border-input resize-y"
          />
        </div>

        {/* Services Grouping Section */}
        <div className="space-y-3 pt-2 border-t">
          <label className="text-xs font-bold text-foreground">Services</label>

          {/* Table Header Row */}
          <div className="grid grid-cols-12 gap-3 px-3 py-2 bg-muted/40 rounded-lg text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">
            <div className="col-span-5">Service Group & Services</div>
            <div className="col-span-2 text-center">Qty</div>
            <div className="col-span-2">Service Cost</div>
            <div className="col-span-2">Sub Total</div>
            <div className="col-span-1 text-center">Actions</div>
          </div>

          {/* Group Blocks */}
          <div className="space-y-4">
            {serviceGroups.map((grp) => {
              const groupTotal = grp.services.reduce(
                (sum, s) => sum + (s.cost || 0) * (s.quantity || 1),
                0
              )
              return (
                <div
                  key={grp.id}
                  className="rounded-xl border border-border bg-muted/10 p-3 space-y-3"
                >
                  {/* Top row of Group */}
                  <div className="grid grid-cols-12 gap-3 items-center">
                    {/* Group Name Select dropdown */}
                    <div className="col-span-5">
                      <Select
                        value={grp.groupName}
                        onValueChange={(val) => handleUpdateGroup(grp.id, 'groupName', val)}
                      >
                        <SelectTrigger className="h-9 text-xs bg-background border-input shadow-2xs">
                          <SelectValue placeholder="Group Name (e.g. Wedding)" />
                        </SelectTrigger>
                        <SelectContent>
                          {DEFAULT_SERVICE_GROUPS.map((g) => (
                            <SelectItem key={g} value={g} className="text-xs">
                              {g}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>

                    {/* Group Quantity */}
                    <div className="col-span-2">
                      <Input
                        type="number"
                        min={1}
                        value={grp.quantity}
                        onChange={(e) =>
                          handleUpdateGroup(grp.id, 'quantity', parseInt(e.target.value, 10) || 1)
                        }
                        className="h-9 text-xs bg-background border-input text-center shadow-2xs tabular-nums"
                      />
                    </div>

                    {/* Service Cost */}
                    <div className="col-span-2">
                      <div className="relative flex items-center">
                        <span className="absolute left-2 text-[10px] text-muted-foreground font-semibold">
                          INR
                        </span>
                        <Input
                          readOnly
                          value={groupTotal}
                          className="pl-8 h-9 text-xs bg-muted/40 border-input font-bold tabular-nums"
                        />
                      </div>
                    </div>

                    {/* Sub Total */}
                    <div className="col-span-2">
                      <div className="relative flex items-center">
                        <span className="absolute left-2 text-[10px] text-muted-foreground font-semibold">
                          INR
                        </span>
                        <Input
                          readOnly
                          value={groupTotal * (grp.quantity || 1)}
                          className="pl-8 h-9 text-xs bg-muted/40 border-input font-bold tabular-nums text-foreground"
                        />
                      </div>
                    </div>

                    {/* Actions: Delete Group & Duplicate Group */}
                    <div className="col-span-1 flex items-center justify-center gap-1">
                      <button
                        type="button"
                        onClick={() => handleRemoveServiceGroup(grp.id)}
                        className="text-muted-foreground hover:text-destructive p-1 rounded transition"
                        title="Delete Group"
                      >
                        <IconTrash className="h-4 w-4" />
                      </button>
                    </div>
                  </div>

                  {/* Sub-services inside this group */}
                  <div className="pl-3 border-l-2 border-primary/20 space-y-2">
                    {grp.services.map((srv) => (
                      <div
                        key={srv.id}
                        className="flex items-center justify-between p-2 rounded-lg bg-card border text-xs"
                      >
                        <div className="flex items-center gap-2">
                          <span className="font-semibold text-foreground">{srv.name}</span>
                          {srv.isCustom && (
                            <Badge variant="outline" className="text-[9px] h-4 px-1">
                              Custom
                            </Badge>
                          )}
                          <span className="text-[11px] text-muted-foreground">₹{srv.cost}</span>
                        </div>

                        {/* Quantity Counter (+ 1 -) matching Image 3 */}
                        <div className="flex items-center gap-1.5 bg-muted/50 border rounded-md px-1.5 py-0.5">
                          <button
                            type="button"
                            onClick={() => handleUpdateServiceQty(grp.id, srv.id, 1)}
                            className="text-muted-foreground hover:text-foreground text-xs font-bold"
                          >
                            +
                          </button>
                          <span className="text-xs font-bold tabular-nums px-1.5">{srv.quantity}</span>
                          <button
                            type="button"
                            onClick={() => handleUpdateServiceQty(grp.id, srv.id, -1)}
                            className="text-muted-foreground hover:text-foreground text-xs font-bold"
                          >
                            -
                          </button>
                        </div>
                      </div>
                    ))}

                    {/* Add Service & Add Custom buttons inside group (matching Image 3) */}
                    <div className="flex items-center gap-3 pt-1 text-xs">
                      <button
                        type="button"
                        onClick={() => handleAddServiceToGroup(grp.id, false)}
                        className="font-bold text-emerald-600 hover:underline inline-flex items-center gap-1"
                      >
                        <IconPlus className="h-3 w-3" />
                        <span>Add Service</span>
                      </button>
                      <button
                        type="button"
                        onClick={() => handleAddServiceToGroup(grp.id, true)}
                        className="font-bold text-emerald-600 hover:underline inline-flex items-center gap-1"
                      >
                        <IconPlus className="h-3 w-3" />
                        <span>Add Custom</span>
                      </button>
                    </div>
                  </div>
                </div>
              )
            })}
          </div>

          {/* Add Service Group dashed button */}
          <div className="pt-2">
            <button
              type="button"
              onClick={handleAddServiceGroup}
              className="w-full py-2.5 rounded-xl border border-dashed border-border hover:border-emerald-600 text-xs font-semibold text-emerald-600 dark:text-emerald-400 hover:bg-emerald-50/50 dark:hover:bg-emerald-950/20 flex items-center justify-center gap-1.5 transition"
            >
              <IconPlus className="h-3.5 w-3.5" />
              <span>Add Service Group</span>
            </button>
          </div>
        </div>

        {/* Deliverables Section */}
        <div className="space-y-1.5 pt-4 border-t">
          <label className="text-[11px] font-bold text-muted-foreground uppercase tracking-wider">
            Deliverables
          </label>
          <DeliverablesDropdown
            selected={selectedDeliverables}
            onChange={setSelectedDeliverables}
          />
        </div>

        {/* Total Cost Display matching Image 3 bottom right */}
        <div className="flex justify-end pt-4 border-t">
          <div className="text-right">
            <span className="text-xl font-black text-foreground">
              Total Cost ₹{totalCalculatedCost.toLocaleString('en-IN')}
            </span>
            <p className="text-[11px] text-muted-foreground">Standard base rate before event taxes</p>
          </div>
        </div>
      </Card>
    </div>
  )
}

