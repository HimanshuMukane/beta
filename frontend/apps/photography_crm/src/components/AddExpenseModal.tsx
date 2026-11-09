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
import { format } from 'date-fns'
import type { BillingExpense } from '../types/billing'

interface AddExpenseModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onAddExpense: (expense: BillingExpense) => void
}

export function AddExpenseModal({
  open,
  onOpenChange,
  onAddExpense,
}: AddExpenseModalProps) {
  const [expenseName, setExpenseName] = React.useState('')
  const [amount, setAmount] = React.useState('')
  const [type, setType] = React.useState('Crew & Crew Travel')
  const [date, setDate] = React.useState(format(new Date(), 'yyyy-MM-dd'))
  const [vendorOrClient, setVendorOrClient] = React.useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!expenseName.trim() || !amount) {
      toast.error('Please fill in required fields')
      return
    }

    const newExpense: BillingExpense = {
      id: `EXP-${Math.floor(100 + Math.random() * 900)}`,
      expenseName: expenseName.trim(),
      amount: parseFloat(amount) || 0,
      type,
      date: format(new Date(date), 'MMM dd, yyyy'),
      createdDate: `${format(new Date(), 'MMM dd, yyyy hh:mm a')}`,
      vendorOrClient: vendorOrClient.trim() || undefined,
    }

    onAddExpense(newExpense)
    toast.success('Expense recorded successfully')
    onOpenChange(false)
    setExpenseName('')
    setAmount('')
    setVendorOrClient('')
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Add New Expense</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4 py-2">
          <div className="space-y-1.5">
            <Label htmlFor="exp-name">Expense Description / Name *</Label>
            <Input
              id="exp-name"
              placeholder="e.g. Drone Operator Rental & Travel"
              value={expenseName}
              onChange={(e) => setExpenseName(e.target.value)}
              required
            />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <Label htmlFor="exp-amount">Amount (₹) *</Label>
              <Input
                id="exp-amount"
                type="number"
                placeholder="4500"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                required
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="exp-date">Date</Label>
              <Input
                id="exp-date"
                type="date"
                value={date}
                onChange={(e) => setDate(e.target.value)}
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <Label>Expense Category / Type</Label>
            <Select value={type} onValueChange={setType}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="Crew & Crew Travel">Crew & Crew Travel</SelectItem>
                <SelectItem value="Equipment Rental">Equipment Rental</SelectItem>
                <SelectItem value="Location Permits">Location Permits</SelectItem>
                <SelectItem value="Print & Album Production">Print & Album Production</SelectItem>
                <SelectItem value="Props & Styling">Props & Styling</SelectItem>
                <SelectItem value="Studio Maintenance">Studio Maintenance</SelectItem>
                <SelectItem value="Software & Cloud Storage">Software & Cloud Storage</SelectItem>
                <SelectItem value="Other">Other</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="exp-vendor">Vendor / Paid To</Label>
            <Input
              id="exp-vendor"
              placeholder="e.g. Cine Equipment Rentals"
              value={vendorOrClient}
              onChange={(e) => setVendorOrClient(e.target.value)}
            />
          </div>

          <DialogFooter className="pt-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              Cancel
            </Button>
            <Button type="submit">Save Expense</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

