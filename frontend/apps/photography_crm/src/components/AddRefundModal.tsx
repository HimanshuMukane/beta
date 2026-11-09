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
import type { BillingRefunded, RefundStatus } from '../types/billing'

interface AddRefundModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onAddRefund: (refund: BillingRefunded) => void
}

export function AddRefundModal({
  open,
  onOpenChange,
  onAddRefund,
}: AddRefundModalProps) {
  const [invoiceName, setInvoiceName] = React.useState('')
  const [clientName, setClientName] = React.useState('')
  const [projectsName, setProjectsName] = React.useState('')
  const [paymentAmount, setPaymentAmount] = React.useState('')
  const [cancellationCharges, setCancellationCharges] = React.useState('0')
  const [note, setNote] = React.useState('')
  const [status, setStatus] = React.useState<RefundStatus>('Completed')
  const [refundDate, setRefundDate] = React.useState(format(new Date(), 'yyyy-MM-dd'))

  const calculatedRefund = React.useMemo(() => {
    const pay = parseFloat(paymentAmount) || 0
    const charge = parseFloat(cancellationCharges) || 0
    return Math.max(0, pay - charge)
  }, [paymentAmount, cancellationCharges])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!invoiceName.trim() || !clientName.trim() || !paymentAmount) {
      toast.error('Please fill in required fields')
      return
    }

    const pay = parseFloat(paymentAmount) || 0
    const charge = parseFloat(cancellationCharges) || 0

    const initials = clientName
      .trim()
      .split(' ')
      .slice(0, 2)
      .map((n) => n[0]?.toUpperCase() || '')
      .join('')

    const newRefund: BillingRefunded = {
      id: `REF-${Math.floor(100 + Math.random() * 900)}`,
      invoiceName: invoiceName.trim(),
      clientName: clientName.trim(),
      clientInitials: initials || 'CL',
      projectsName: projectsName.trim() || 'Custom Photography Package',
      paymentAmount: pay,
      cancellationCharges: charge,
      refundAmount: Math.max(0, pay - charge),
      note: note.trim() || 'Refund processed per client request',
      status,
      refundDate: `${format(new Date(refundDate), 'MMM dd, yyyy')} 12:00 PM`,
      createdOn: `${format(new Date(), 'MMM dd, yyyy hh:mm a')}`,
    }

    onAddRefund(newRefund)
    toast.success('Refund recorded successfully')
    onOpenChange(false)
    setInvoiceName('')
    setClientName('')
    setProjectsName('')
    setPaymentAmount('')
    setCancellationCharges('0')
    setNote('')
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Issue / Record Refund</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4 py-2">
          <div className="space-y-1.5">
            <Label htmlFor="ref-inv">Invoice / Item Reference *</Label>
            <Input
              id="ref-inv"
              placeholder="e.g. Pre-Wedding Advance Retainer"
              value={invoiceName}
              onChange={(e) => setInvoiceName(e.target.value)}
              required
            />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <Label htmlFor="ref-client">Client Name *</Label>
              <Input
                id="ref-client"
                placeholder="e.g. Ritu Varma"
                value={clientName}
                onChange={(e) => setClientName(e.target.value)}
                required
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="ref-proj">Project Name</Label>
              <Input
                id="ref-proj"
                placeholder="e.g. Ritu & Dev Engagement"
                value={projectsName}
                onChange={(e) => setProjectsName(e.target.value)}
              />
            </div>
          </div>

          <div className="grid grid-cols-3 gap-3">
            <div className="space-y-1.5">
              <Label htmlFor="ref-paid">Payment Amount (₹) *</Label>
              <Input
                id="ref-paid"
                type="number"
                placeholder="20000"
                value={paymentAmount}
                onChange={(e) => setPaymentAmount(e.target.value)}
                required
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="ref-charges">Cancel Charges (₹)</Label>
              <Input
                id="ref-charges"
                type="number"
                placeholder="2000"
                value={cancellationCharges}
                onChange={(e) => setCancellationCharges(e.target.value)}
              />
            </div>
            <div className="space-y-1.5">
              <Label>Refund Amount (₹)</Label>
              <div className="h-9 px-3 rounded-md bg-muted/50 border flex items-center font-bold text-emerald-600 text-sm">
                ₹{calculatedRefund.toLocaleString('en-IN')}
              </div>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <Label htmlFor="ref-date">Refund Date</Label>
              <Input
                id="ref-date"
                type="date"
                value={refundDate}
                onChange={(e) => setRefundDate(e.target.value)}
              />
            </div>
            <div className="space-y-1.5">
              <Label>Status</Label>
              <Select
                value={status}
                onValueChange={(val) => setStatus(val as RefundStatus)}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Completed">Completed</SelectItem>
                  <SelectItem value="Processing">Processing</SelectItem>
                  <SelectItem value="Failed">Failed</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="ref-note">Reason / Note</Label>
            <Input
              id="ref-note"
              placeholder="e.g. Partial refund due to weather rescheduling"
              value={note}
              onChange={(e) => setNote(e.target.value)}
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
            <Button type="submit">Process Refund</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

