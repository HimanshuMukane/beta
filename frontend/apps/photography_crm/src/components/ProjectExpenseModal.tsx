import * as React from 'react'
import {
  Button,
  Input,
  Label,
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  toast,
} from '@repo/ui'
import type { ProjectExpenseItem } from '../types/project'

interface ProjectExpenseModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSaveExpense: (expense: ProjectExpenseItem) => void
  initialExpense?: ProjectExpenseItem | null
  existingSessions?: string[]
}

export function ProjectExpenseModal({
  open,
  onOpenChange,
  onSaveExpense,
  initialExpense,
  existingSessions = ['Reception | Sunset', 'Wedding | Golden Hour (Evening)'],
}: ProjectExpenseModalProps) {
  const [sessionName, setSessionName] = React.useState('')
  const [spentBy, setSpentBy] = React.useState('E')
  const [category, setCategory] = React.useState('Travel Allowance')
  const [amount, setAmount] = React.useState('')
  const [notes, setNotes] = React.useState('')

  React.useEffect(() => {
    if (initialExpense) {
      setSessionName(initialExpense.sessionName || existingSessions[0] || '')
      setSpentBy(initialExpense.paidBy || 'E')
      setCategory(initialExpense.category || 'Travel Allowance')
      setAmount(initialExpense.amount?.toString() || '')
      setNotes(initialExpense.description || '')
    } else {
      setSessionName(existingSessions[0] || 'Reception | Sunset')
      setSpentBy('E')
      setCategory('Travel Allowance')
      setAmount('')
      setNotes('')
    }
  }, [initialExpense, open, existingSessions])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!amount) {
      toast.error('Please enter expense amount')
      return
    }

    const item: ProjectExpenseItem = {
      id: initialExpense ? initialExpense.id : `exp-${Date.now()}`,
      category,
      description: notes.trim() || category,
      amount: parseFloat(amount) || 0,
      date: initialExpense ? initialExpense.date : new Date().toISOString().split('T')[0],
      paidBy: spentBy,
      sessionName: sessionName || existingSessions[0],
    }

    onSaveExpense(item)
    toast.success(initialExpense ? 'Expense updated' : 'Expense recorded')
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>
            {initialExpense ? 'Edit Project Expense' : 'Add Project Expense'}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4 py-2">
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <Label>Session Name</Label>
              <Select value={sessionName} onValueChange={setSessionName}>
                <SelectTrigger className="text-xs">
                  <SelectValue placeholder="Select session" />
                </SelectTrigger>
                <SelectContent>
                  {existingSessions.map((s) => (
                    <SelectItem key={s} value={s}>
                      {s}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1.5">
              <Label>Spent By</Label>
              <Input
                placeholder="e.g. E (Crew lead)"
                value={spentBy}
                onChange={(e) => setSpentBy(e.target.value)}
                className="text-xs"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <Label>Expense Type</Label>
              <Select value={category} onValueChange={setCategory}>
                <SelectTrigger className="text-xs">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Travel Allowance">Travel Allowance</SelectItem>
                  <SelectItem value="Rental">Rental</SelectItem>
                  <SelectItem value="Crew Food & Refreshments">Crew Food & Refreshments</SelectItem>
                  <SelectItem value="Permits & Entry">Permits & Entry</SelectItem>
                  <SelectItem value="Other">Other</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1.5">
              <Label>Amount (₹) *</Label>
              <Input
                type="number"
                placeholder="1500"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                required
                className="text-xs"
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <Label>Notes</Label>
            <Input
              placeholder="e.g. Tolls and cab fare for team"
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              className="text-xs"
            />
          </div>

          <DialogFooter className="pt-2">
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => onOpenChange(false)}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              size="sm"
              className="bg-emerald-600 hover:bg-emerald-700 text-white"
            >
              Save Expense
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

