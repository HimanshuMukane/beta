import * as React from 'react'
import {
  IconPlus,
  IconTrash,
  IconFolderPlus,
  IconPackage,
  IconCheck,
  IconLayersSubtract,
} from '@tabler/icons-react'
import {
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
  Tooltip,
  TooltipTrigger,
  TooltipContent,
  toast,
} from '@repo/ui'
import { useStudioSettings } from '../../hooks/useStudioSettings'

export interface DeliverablePresetItem {
  id: string
  groupName: string
  name: string
}

const DEFAULT_GROUPS = ['Wedding', 'Internal Works', 'Complimentary']

const INITIAL_DELIVERABLES: DeliverablePresetItem[] = [
  { id: 'del-1', groupName: 'Wedding', name: 'Raw Photos' },
  { id: 'del-2', groupName: 'Wedding', name: 'Album' },
  { id: 'del-3', groupName: 'Wedding', name: 'Pen drive' },
  { id: 'del-4', groupName: 'Wedding', name: 'Wedding Teaser' },
  { id: 'del-5', groupName: 'Wedding', name: 'Highlights' },
  { id: 'del-6', groupName: 'Complimentary', name: 'Photo Calendar' },
  { id: 'del-7', groupName: 'Complimentary', name: 'Custom Coffee Mug' },
  { id: 'del-8', groupName: 'Complimentary', name: 'Mini Photo Frame' },
  { id: 'del-9', groupName: 'Complimentary', name: 'Photo Keychain' },
]

export function DeliverablesSettingsPage() {
  const { settings, updateSettings } = useStudioSettings()
  const [items, setItems] = React.useState<DeliverablePresetItem[]>(
    settings?.deliverablesPresets || INITIAL_DELIVERABLES
  )
  const [groups, setGroups] = React.useState<string[]>(
    settings?.deliverableGroups || DEFAULT_GROUPS
  )
  const [isGroupModalOpen, setIsGroupModalOpen] = React.useState(false)
  const [newGroupName, setNewGroupName] = React.useState('')
  const [filterGroup, setFilterGroup] = React.useState<string>('all')

  React.useEffect(() => {
    if (settings?.deliverablesPresets) {
      setItems(settings.deliverablesPresets)
    }
    if (settings?.deliverableGroups) {
      setGroups(settings.deliverableGroups)
    }
  }, [settings?.deliverablesPresets, settings?.deliverableGroups])

  const handleUpdate = (id: string, field: keyof DeliverablePresetItem, value: string) => {
    setItems((prev) =>
      prev.map((item) => (item.id === id ? { ...item, [field]: value } : item))
    )
  }

  const handleAddItem = () => {
    const defaultGroup = filterGroup !== 'all' ? filterGroup : groups[0] || 'Wedding'
    const newItem: DeliverablePresetItem = {
      id: `del-${Date.now()}`,
      groupName: defaultGroup,
      name: '',
    }
    setItems((prev) => [...prev, newItem])
  }

  const handleDeleteItem = (id: string) => {
    const updated = items.filter((item) => item.id !== id)
    setItems(updated)
    updateSettings({ deliverablesPresets: updated })
    toast.success('Deliverable removed')
  }

  const handleAddGroup = () => {
    if (!newGroupName.trim()) return
    const trimmed = newGroupName.trim()
    if (groups.includes(trimmed)) {
      toast.error('Group already exists')
      return
    }
    const updatedGroups = [...groups, trimmed]
    setGroups(updatedGroups)
    updateSettings({ deliverableGroups: updatedGroups })
    setNewGroupName('')
    setIsGroupModalOpen(false)
    toast.success(`Group "${trimmed}" created`)
  }

  const handleDeleteGroup = (groupToDelete: string) => {
    if (groups.length <= 1) {
      toast.error('At least one group must remain')
      return
    }
    const updatedGroups = groups.filter((g) => g !== groupToDelete)
    setGroups(updatedGroups)
    const fallback = updatedGroups[0] || 'Wedding'
    const updatedItems = items.map((i) =>
      i.groupName === groupToDelete ? { ...i, groupName: fallback } : i
    )
    setItems(updatedItems)
    updateSettings({ deliverableGroups: updatedGroups, deliverablesPresets: updatedItems })
    toast.success(`Removed group "${groupToDelete}"`)
  }

  const handleSaveAll = () => {
    const valid = items.filter((i) => i.name.trim().length > 0)
    setItems(valid)
    updateSettings({ deliverablesPresets: valid, deliverableGroups: groups })
    toast.success('Deliverables presets saved successfully')
  }

  const filteredItems = filterGroup === 'all'
    ? items
    : items.filter((i) => i.groupName === filterGroup)

  return (
    <div className="max-w-4xl mx-auto space-y-6 pb-16">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b pb-4">
        <div className="flex items-center gap-3">
          <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary border border-primary/20">
            <IconPackage className="h-5 w-5" />
          </span>
          <div>
            <h1 className="text-xl font-bold tracking-tight text-foreground">Deliverables Presets</h1>
            <p className="text-xs text-muted-foreground">
              Define standard deliverables and organize them by groups for projects and quotes.
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2 self-start sm:self-auto">
          <Button
            type="button"
            variant="outline"
            onClick={() => setIsGroupModalOpen(true)}
            className="h-8 text-xs font-medium gap-1.5 border-dashed"
          >
            <IconFolderPlus className="h-3.5 w-3.5 text-primary" />
            Manage Groups ({groups.length})
          </Button>

          <Button
            type="button"
            onClick={handleSaveAll}
            className="h-8 text-xs font-semibold gap-1.5 bg-primary text-primary-foreground shadow-xs"
          >
            <IconCheck className="h-3.5 w-3.5" />
            Save Changes
          </Button>
        </div>
      </div>

      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Select value={filterGroup} onValueChange={setFilterGroup}>
            <SelectTrigger className="h-8 text-xs w-48 bg-background">
              <SelectValue placeholder="Filter by group" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Groups ({items.length})</SelectItem>
              {groups.map((g) => (
                <SelectItem key={g} value={g}>
                  {g} ({items.filter((i) => i.groupName === g).length})
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <span className="text-xs text-muted-foreground font-medium">
          Showing <strong>{filteredItems.length}</strong> of <strong>{items.length}</strong> deliverables
        </span>
      </div>

      <div className="bg-card border border-border/80 shadow-2xs rounded-xl overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse text-xs">
            <thead>
              <tr className="border-b bg-muted/40 text-muted-foreground text-[11px] font-semibold uppercase tracking-wider">
                <th className="py-3 px-4 w-[40%]">Group Name</th>
                <th className="py-3 px-4 w-[54%]">Deliverables Name</th>
                <th className="py-3 px-3 w-[6%] text-center"></th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border/60">
              {filteredItems.map((item) => (
                <tr key={item.id} className="hover:bg-muted/10 transition-colors">
                  <td className="py-2.5 px-4">
                    <Select
                      value={item.groupName}
                      onValueChange={(val) => handleUpdate(item.id, 'groupName', val)}
                    >
                      <SelectTrigger className="h-8 text-xs bg-background">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        {groups.map((g) => (
                          <SelectItem key={g} value={g}>
                            {g}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </td>
                  <td className="py-2.5 px-4">
                    <Input
                      value={item.name}
                      onChange={(e) => handleUpdate(item.id, 'name', e.target.value)}
                      placeholder="e.g. Raw Photos, Album, Pen drive..."
                      className="h-8 text-xs bg-background"
                    />
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
                        Remove deliverable
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
            Add Deliverable
          </Button>

          <span className="text-xs text-muted-foreground flex items-center gap-1.5">
            <IconPackage className="h-3.5 w-3.5 text-primary" />
            Showing <strong>{filteredItems.length}</strong> of <strong>{items.length}</strong> deliverables
          </span>
        </div>
      </div>

      {/* Manage Groups Modal */}
      <Dialog open={isGroupModalOpen} onOpenChange={setIsGroupModalOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="text-sm font-bold flex items-center gap-2">
              <IconLayersSubtract className="h-4 w-4 text-primary" />
              Manage Deliverables Groups
            </DialogTitle>
            <DialogDescription className="text-xs text-muted-foreground">
              Define the categories used to group client deliverables (e.g., Wedding, Complimentary, Corporate).
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-2">
            <div className="flex items-center gap-2">
              <Input
                value={newGroupName}
                onChange={(e) => setNewGroupName(e.target.value)}
                placeholder="New group name (e.g. Maternity, VIP)..."
                className="h-8 text-xs flex-1"
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    e.preventDefault()
                    handleAddGroup()
                  }
                }}
              />
              <Button
                type="button"
                onClick={handleAddGroup}
                className="h-8 text-xs font-semibold gap-1 bg-primary text-primary-foreground"
              >
                <IconPlus className="h-3.5 w-3.5" />
                Add
              </Button>
            </div>

            <div className="border rounded-lg divide-y max-h-56 overflow-y-auto">
              {groups.map((group) => {
                const count = items.filter((i) => i.groupName === group).length
                return (
                  <div
                    key={group}
                    className="p-2.5 flex items-center justify-between hover:bg-muted/30 text-xs"
                  >
                    <span className="font-medium text-foreground">{group}</span>
                    <div className="flex items-center gap-3">
                      <span className="text-[11px] text-muted-foreground bg-muted px-2 py-0.5 rounded-full">
                        {count} item{count !== 1 ? 's' : ''}
                      </span>
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        onClick={() => handleDeleteGroup(group)}
                        className="h-6 w-6 text-muted-foreground hover:text-destructive"
                        title="Delete group"
                      >
                        <IconTrash className="h-3 w-3" />
                      </Button>
                    </div>
                  </div>
                )
              })}
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setIsGroupModalOpen(false)}
              className="h-8 text-xs"
            >
              Done
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

