import * as React from 'react'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@repo/ui/components/ui/dialog'
import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import { Label } from '@repo/ui/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@repo/ui/components/ui/select'
import { Textarea } from '@repo/ui/components/ui/textarea'
import { toast } from '@repo/ui/components/ui/sonner'
import type { CustomerDeliverable } from '../types/project'

interface AddDeliverableModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onAdd: (deliverable: CustomerDeliverable) => void
}

const COMMON_DELIVERABLES = [
  'Raw Photos',
  'Album',
  'Pen drive',
  'Wedding Teaser',
  'Cinematic Highlights Film',
  'Traditional Video Edit',
  'Drone Reel',
  'Fine Art Portrait Prints',
  'Online Client Gallery',
]

const TEAM_MEMBERS = [
  'Rohan Mehta (Lead Photographer)',
  'Tanvi Shah (Cinematographer)',
  'Aditya Roy (Candid Photographer)',
  'Vikram Sen (Video Editor)',
  'Kavita Roy (Album Designer)',
]

export function AddDeliverableModal({
  open,
  onOpenChange,
  onAdd,
}: AddDeliverableModalProps) {
  const [deliverable, setDeliverable] = React.useState('')
  const [dueDate, setDueDate] = React.useState(
    new Date(Date.now() + 14 * 24 * 60 * 60 * 1000).toISOString().split('T')[0]
  )
  const [assignee, setAssignee] = React.useState('')
  const [status, setStatus] = React.useState<CustomerDeliverable['status']>('New')
  const [note, setNote] = React.useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!deliverable.trim()) {
      toast.error('Please enter or select a deliverable name')
      return
    }

    const newItem: CustomerDeliverable = {
      id: `del-${Date.now()}`,
      deliverable: deliverable.trim(),
      dueDate,
      assignee: assignee || 'Unassigned',
      status,
      note: note.trim(),
    }

    onAdd(newItem)
    toast.success('Deliverable added!')
    onOpenChange(false)
    // Reset form
    setDeliverable('')
    setAssignee('')
    setStatus('New')
    setNote('')
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-lg p-0 overflow-hidden bg-background">
        <DialogHeader className="px-6 pt-5 pb-3 border-b bg-muted/20">
          <DialogTitle className="text-base font-bold text-foreground">
            Add Customer Deliverable
          </DialogTitle>
          <DialogDescription className="text-xs text-muted-foreground">
            Schedule shoot outputs, assign editors, and add delivery notes.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="px-6 py-5 space-y-4">
          {/* Deliverables Input with suggestion chips */}
          <div className="space-y-1.5">
            <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Deliverables <span className="text-destructive">*</span>
            </Label>
            <Input
              value={deliverable}
              onChange={(e) => setDeliverable(e.target.value)}
              placeholder="e.g. Raw Photos, Album, Wedding Teaser"
              className="h-10 text-sm"
              required
              autoFocus
            />
            {/* Quick suggestion tags */}
            <div className="flex flex-wrap gap-1.5 pt-1">
              {COMMON_DELIVERABLES.slice(0, 5).map((item) => (
                <button
                  key={item}
                  type="button"
                  onClick={() => setDeliverable(item)}
                  className="rounded-md border bg-muted/40 px-2 py-0.5 text-[11px] text-muted-foreground hover:bg-muted hover:text-foreground transition"
                >
                  + {item}
                </button>
              ))}
            </div>
          </div>

          {/* Assign Team Member */}
          <div className="space-y-1.5">
            <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Assign Team Member
            </Label>
            <Select value={assignee} onValueChange={setAssignee}>
              <SelectTrigger className="h-10 text-sm">
                <SelectValue placeholder="Select Assigned Team Member" />
              </SelectTrigger>
              <SelectContent>
                {TEAM_MEMBERS.map((tm) => (
                  <SelectItem key={tm} value={tm} className="text-xs">
                    {tm}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {/* Due Date & Status (2 columns) */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                Due Date
              </Label>
              <Input
                type="date"
                value={dueDate}
                onChange={(e) => setDueDate(e.target.value)}
                className="h-10 text-sm"
              />
            </div>

            <div className="space-y-1.5">
              <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                Status
              </Label>
              <Select
                value={status}
                onValueChange={(val) => setStatus(val as CustomerDeliverable['status'])}
              >
                <SelectTrigger className="h-10 text-sm">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="New" className="text-xs">New</SelectItem>
                  <SelectItem value="In Progress" className="text-xs">In Progress</SelectItem>
                  <SelectItem value="Client Review" className="text-xs">Client Review</SelectItem>
                  <SelectItem value="Completed" className="text-xs">Completed</SelectItem>
                  <SelectItem value="Delivered" className="text-xs">Delivered</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          {/* Note */}
          <div className="space-y-1.5">
            <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Note
            </Label>
            <Textarea
              value={note}
              onChange={(e) => setNote(e.target.value)}
              placeholder="Enter Your Note Here..."
              className="min-h-20 text-xs resize-none"
            />
          </div>

          {/* Footer Actions */}
          <DialogFooter className="pt-3 border-t">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              className="bg-emerald-600 hover:bg-emerald-700 text-white font-medium px-5"
            >
              Add Deliverable
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

