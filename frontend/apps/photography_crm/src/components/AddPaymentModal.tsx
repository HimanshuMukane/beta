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
import type { BillingPayment, PaymentStatus } from '../types/billing'

interface AddPaymentModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onAddPayment: (payment: BillingPayment) => void
}

export function AddPaymentModal({
  open,
  onOpenChange,
  onAddPayment,
}: AddPaymentModalProps) {
  const [clientName, setClientName] = React.useState('')
  const [projectName, setProjectName] = React.useState('')
  const [amount, setAmount] = React.useState('')
  const [type, setType] = React.useState('UPI Transfer')
  const [status, setStatus] = React.useState<PaymentStatus>('Completed')
  const [paymentDate, setPaymentDate] = React.useState(format(new Date(), 'yyyy-MM-dd'))

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!clientName.trim() || !amount) {
      toast.error('Please enter client name and amount')
      return
    }

    const initials = clientName
      .trim()
      .split(' ')
      .slice(0, 2)
      .map((n) => n[0]?.toUpperCase() || '')
      .join('')

    const newPayment: BillingPayment = {
      id: `ESPE/PAY/${Math.floor(100 + Math.random() * 900)}`,
      clientName: clientName.trim(),
      clientInitials: initials || 'CL',
      projectName: projectName.trim() || 'General Photography Booking',
      paymentAmount: parseFloat(amount) || 0,
      type,
      paymentStatus: status,
      paymentDate: `${format(new Date(paymentDate), 'MMM dd, yyyy')} 12:00 PM`,
      createdOn: `${format(new Date(), 'MMM dd, yyyy hh:mm a')}`,
      downloadUrl: '#',
    }

    onAddPayment(newPayment)
    toast.success('Payment recorded successfully')
    onOpenChange(false)
    setClientName('')
    setProjectName('')
    setAmount('')
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Record New Payment</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4 py-2">
          <div className="space-y-1.5">
            <Label htmlFor="pay-client">Client Name *</Label>
            <Input
              id="pay-client"
              placeholder="e.g. Aditi Sharma"
              value={clientName}
              onChange={(e) => setClientName(e.target.value)}
              required
            />
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="pay-project">Project / Event Name</Label>
            <Input
              id="pay-project"
              placeholder="e.g. Aditi & Kunal Wedding"
              value={projectName}
              onChange={(e) => setProjectName(e.target.value)}
            />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <Label htmlFor="pay-amount">Amount (₹) *</Label>
              <Input
                id="pay-amount"
                type="number"
                placeholder="25000"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                required
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="pay-date">Payment Date</Label>
              <Input
                id="pay-date"
                type="date"
                value={paymentDate}
                onChange={(e) => setPaymentDate(e.target.value)}
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <Label>Payment Mode / Type</Label>
              <Select value={type} onValueChange={setType}>
                <SelectTrigger>
                  <SelectValue placeholder="Select type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="UPI Transfer">UPI Transfer</SelectItem>
                  <SelectItem value="Invoice">Invoice</SelectItem>
                  <SelectItem value="Bank Transfer (NEFT)">Bank Transfer (NEFT)</SelectItem>
                  <SelectItem value="Credit Card">Credit Card</SelectItem>
                  <SelectItem value="Cash">Cash</SelectItem>
                  <SelectItem value="Cheque">Cheque</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1.5">
              <Label>Status</Label>
              <Select
                value={status}
                onValueChange={(val) => setStatus(val as PaymentStatus)}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Completed">Completed</SelectItem>
                  <SelectItem value="Payment Due">Payment Due</SelectItem>
                  <SelectItem value="Pending">Pending</SelectItem>
                  <SelectItem value="Failed">Failed</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <DialogFooter className="pt-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              Cancel
            </Button>
            <Button type="submit">Save Payment</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

