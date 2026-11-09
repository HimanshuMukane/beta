import * as React from 'react'
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerDescription,
  DrawerFooter,
  DrawerClose,
  Button,
  Input,
  Label,
  Badge,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Textarea,
  toast,
  cn,
} from '@repo/ui'
import {
  IconCheck,
  IconCalendar,
  IconUser,
} from '@tabler/icons-react'

export interface TaskDetailData {
  taskId: string
  taskName: string
  categoryName: string
  phase: 'Pre-production' | 'Post-production'
  status: 'Pending' | 'In Progress' | 'Completed' | 'Blocked'
  priority: 'Low' | 'Medium' | 'High'
  dueDate: string
  assignee: string
  notes: string
  checklist: { id: string; text: string; done: boolean }[]
}

interface ProductionTaskDrawerProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  task: TaskDetailData | null
  onSaveTask: (updatedTask: TaskDetailData) => void
}

export function ProductionTaskDrawer({
  open,
  onOpenChange,
  task,
  onSaveTask,
}: ProductionTaskDrawerProps) {
  const [formData, setFormData] = React.useState<TaskDetailData | null>(null)
  const [newChecklistItem, setNewChecklistItem] = React.useState('')

  React.useEffect(() => {
    if (task) {
      setFormData({ ...task })
    }
  }, [task])

  if (!formData) return null

  const handleToggleChecklist = (id: string) => {
    setFormData((prev) => {
      if (!prev) return null
      return {
        ...prev,
        checklist: prev.checklist.map((item) =>
          item.id === id ? { ...item, done: !item.done } : item
        ),
      }
    })
  }

  const handleAddChecklistItem = (e: React.FormEvent) => {
    e.preventDefault()
    if (!newChecklistItem.trim() || !formData) return
    const newItem = {
      id: `chk-${Date.now()}`,
      text: newChecklistItem.trim(),
      done: false,
    }
    setFormData({
      ...formData,
      checklist: [...formData.checklist, newItem],
    })
    setNewChecklistItem('')
  }

  const handleSave = () => {
    if (!formData) return
    onSaveTask(formData)
    toast.success(`Updated "${formData.taskName}"`)
    onOpenChange(false)
  }

  const completedCount = formData.checklist.filter((c) => c.done).length

  return (
    <Drawer open={open} onOpenChange={onOpenChange}>
      <DrawerContent className="max-h-[92vh] flex flex-col sm:max-w-xl mx-auto rounded-t-2xl bg-card border border-border shadow-2xl">
        <DrawerHeader className="border-b pb-3 px-6 text-left">
          <div className="flex items-center gap-2 mb-1">
            <Badge
              variant="outline"
              className={cn(
                'text-[10px] font-bold uppercase tracking-wider',
                formData.phase === 'Pre-production'
                  ? 'bg-blue-500/10 text-blue-600 border-blue-500/20 dark:text-blue-400'
                  : 'bg-purple-500/10 text-purple-600 border-purple-500/20 dark:text-purple-400'
              )}
            >
              {formData.phase}
            </Badge>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs font-semibold text-muted-foreground uppercase">
              {formData.categoryName}
            </span>
          </div>

          <DrawerTitle className="text-lg font-bold text-foreground">
            {formData.taskName}
          </DrawerTitle>
          <DrawerDescription className="text-xs text-muted-foreground">
            Configure stage execution, checklist verification, and assignee responsibilities.
          </DrawerDescription>
        </DrawerHeader>

        <div className="flex-1 overflow-y-auto px-6 py-4 space-y-5">
          {/* Status & Priority Row */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <Label className="text-xs">Status</Label>
              <Select
                value={formData.status}
                onValueChange={(val) =>
                  setFormData({
                    ...formData,
                    status: val as TaskDetailData['status'],
                  })
                }
              >
                <SelectTrigger className="h-9 text-xs">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Pending">Pending</SelectItem>
                  <SelectItem value="In Progress">In Progress</SelectItem>
                  <SelectItem value="Completed">Completed</SelectItem>
                  <SelectItem value="Blocked">Blocked</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1.5">
              <Label className="text-xs">Priority</Label>
              <Select
                value={formData.priority}
                onValueChange={(val) =>
                  setFormData({
                    ...formData,
                    priority: val as TaskDetailData['priority'],
                  })
                }
              >
                <SelectTrigger className="h-9 text-xs">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Low">Low</SelectItem>
                  <SelectItem value="Medium">Medium</SelectItem>
                  <SelectItem value="High">High</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          {/* Assignee & Due Date Row */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <Label className="text-xs">Responsible Crew / Lead</Label>
              <div className="relative">
                <Input
                  className="h-9 text-xs pl-8"
                  value={formData.assignee}
                  onChange={(e) =>
                    setFormData({ ...formData, assignee: e.target.value })
                  }
                  placeholder="e.g. Lead Photographer"
                />
                <IconUser className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              </div>
            </div>

            <div className="space-y-1.5">
              <Label className="text-xs">Target Due Date</Label>
              <div className="relative">
                <Input
                  type="date"
                  className="h-9 text-xs pl-8"
                  value={formData.dueDate}
                  onChange={(e) =>
                    setFormData({ ...formData, dueDate: e.target.value })
                  }
                />
                <IconCalendar className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              </div>
            </div>
          </div>

          {/* Step Verification Checklist */}
          <div className="space-y-3 pt-2">
            <div className="flex items-center justify-between">
              <Label className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
                Action Checklist ({completedCount}/{formData.checklist.length})
              </Label>
              {formData.checklist.length > 0 && (
                <span className="text-[11px] font-semibold text-primary">
                  {Math.round((completedCount / formData.checklist.length) * 100)}%
                </span>
              )}
            </div>

            <div className="space-y-1.5 max-h-48 overflow-y-auto rounded-lg border bg-muted/20 p-2.5">
              {formData.checklist.length === 0 ? (
                <p className="text-xs text-muted-foreground text-center py-2">
                  No sub-items added yet.
                </p>
              ) : (
                formData.checklist.map((item) => (
                  <label
                    key={item.id}
                    className="flex items-center gap-2.5 p-1.5 rounded-md hover:bg-muted/40 cursor-pointer text-xs transition"
                  >
                    <input
                      type="checkbox"
                      checked={item.done}
                      onChange={() => handleToggleChecklist(item.id)}
                      className="h-4 w-4 rounded border-input accent-primary"
                    />
                    <span
                      className={cn(
                        'flex-1 text-foreground font-medium',
                        item.done && 'line-through text-muted-foreground'
                      )}
                    >
                      {item.text}
                    </span>
                  </label>
                ))
              )}
            </div>

            <form onSubmit={handleAddChecklistItem} className="flex gap-2">
              <Input
                placeholder="Add checklist item..."
                value={newChecklistItem}
                onChange={(e) => setNewChecklistItem(e.target.value)}
                className="h-8 text-xs flex-1"
              />
              <Button type="submit" size="sm" variant="secondary" className="h-8 text-xs">
                Add
              </Button>
            </form>
          </div>

          {/* Special Notes & Execution Details */}
          <div className="space-y-1.5 pt-1">
            <Label className="text-xs">Instructions & Feedback Notes</Label>
            <Textarea
              rows={3}
              placeholder="Provide camera gear instructions, client mood board reference, or editing feedback..."
              value={formData.notes}
              onChange={(e) =>
                setFormData({ ...formData, notes: e.target.value })
              }
              className="text-xs resize-none"
            />
          </div>
        </div>

        <DrawerFooter className="border-t pt-3 pb-4 px-6 flex flex-row items-center justify-end gap-2 bg-muted/10">
          <DrawerClose asChild>
            <Button variant="outline" size="sm" className="h-8 text-xs">
              Cancel
            </Button>
          </DrawerClose>
          <Button
            size="sm"
            onClick={handleSave}
            className="h-8 text-xs bg-emerald-600 hover:bg-emerald-700 text-white"
          >
            <IconCheck className="h-3.5 w-3.5 mr-1.5" />
            Save Changes
          </Button>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  )
}

