import * as React from 'react'
import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import { Label } from '@repo/ui/components/ui/label'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@repo/ui/components/ui/dialog'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@repo/ui/components/ui/select'
import { toast } from '@repo/ui/components/ui/sonner'
import { format } from 'date-fns'
import type { BillingInvoice, InvoiceStatus } from '../types/billing'

interface AddInvoiceModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onAddInvoice: (invoice: BillingInvoice) => void
}

export function AddInvoiceModal({
  open,
  onOpenChange,
  onAddInvoice,
}: AddInvoiceModalProps) {
  const [invoiceName, setInvoiceName] = React.useState('')
  const [clientName, setClientName] = React.useState('')
  const [projectsName, setProjectsName] = React.useState('')
  const [totalAmount, setTotalAmount] = React.useState('')
  const [paidAmount, setPaidAmount] = React.useState('0')
  const [status, setStatus] = React.useState<InvoiceStatus>('Unpaid')
  const [invoiceDate, setInvoiceDate] = React.useState(format(new Date(), 'yyyy-MM-dd'))

  const handleTotalChange = (val: string) => {
    setTotalAmount(val)
    const tot = parseFloat(val) || 0
    const pd = parseFloat(paidAmount) || 0
    if (pd >= tot && tot > 0) setStatus('Paid')
    else if (pd > 0) setStatus('Partially Paid')
    else setStatus('Unpaid')
  }

  const handlePaidChange = (val: string) => {
    setPaidAmount(val)
    const tot = parseFloat(totalAmount) || 0
    const pd = parseFloat(val) || 0
    if (pd >= tot && tot > 0) setStatus('Paid')
    else if (pd > 0) setStatus('Partially Paid')
    else setStatus('Unpaid')
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!invoiceName.trim() || !clientName.trim() || !totalAmount) {
      toast.error('Please fill in required fields')
      return
    }

    const tot = parseFloat(totalAmount) || 0
    const pd = parseFloat(paidAmount) || 0
    const pending = Math.max(0, tot - pd)

    const initials = clientName
      .trim()
      .split(' ')
      .slice(0, 2)
      .map((n) => n[0]?.toUpperCase() || '')
      .join('')

    const newInvoice: BillingInvoice = {
      id: `INV-2026-${Math.floor(100 + Math.random() * 900)}`,
      invoiceName: invoiceName.trim(),
      clientName: clientName.trim(),
      clientInitials: initials || 'CL',
      projectsName: projectsName.trim() || 'Wedding Coverage',
      totalAmount: tot,
      paidAmount: pd,
      pendingAmount: pending,
      paymentStatus: status,
      invoiceDate: format(new Date(invoiceDate), 'MMM dd, yyyy'),
      createdDate: `${format(new Date(), 'MMM dd, yyyy hh:mm a')}`,
      downloadUrl: '#',
    }

    onAddInvoice(newInvoice)
    toast.success('Invoice created successfully')
    onOpenChange(false)
    setInvoiceName('')
    setClientName('')
    setProjectsName('')
    setTotalAmount('')
    setPaidAmount('0')
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Create New Invoice</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4 py-2">
          <div className="space-y-1.5">
            <Label htmlFor="inv-name">Invoice Title / Description *</Label>
            <Input
              id="inv-name"
              placeholder="e.g. Wedding Cinematic Film & Album"
              value={invoiceName}
              onChange={(e) => setInvoiceName(e.target.value)}
              required
            />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <Label htmlFor="inv-client">Client Name *</Label>
              <Input
                id="inv-client"
                placeholder="e.g. Shruti Nair"
                value={clientName}
                onChange={(e) => setClientName(e.target.value)}
                required
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="inv-proj">Project Name</Label>
              <Input
                id="inv-proj"
                placeholder="e.g. Shruti Wedding"
                value={projectsName}
                onChange={(e) => setProjectsName(e.target.value)}
              />
            </div>
          </div>

          <div className="grid grid-cols-3 gap-3">
            <div className="space-y-1.5">
              <Label htmlFor="inv-total">Total (₹) *</Label>
              <Input
                id="inv-total"
                type="number"
                placeholder="50000"
                value={totalAmount}
                onChange={(e) => handleTotalChange(e.target.value)}
                required
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="inv-paid">Paid (₹)</Label>
              <Input
                id="inv-paid"
                type="number"
                placeholder="0"
                value={paidAmount}
                onChange={(e) => handlePaidChange(e.target.value)}
              />
            </div>
            <div className="space-y-1.5">
              <Label>Status</Label>
              <Select
                value={status}
                onValueChange={(val) => setStatus(val as InvoiceStatus)}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Unpaid">Unpaid</SelectItem>
                  <SelectItem value="Partially Paid">Partially Paid</SelectItem>
                  <SelectItem value="Paid">Paid</SelectItem>
                  <SelectItem value="Overdue">Overdue</SelectItem>
                  <SelectItem value="Draft">Draft</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="inv-date">Invoice Date</Label>
            <Input
              id="inv-date"
              type="date"
              value={invoiceDate}
              onChange={(e) => setInvoiceDate(e.target.value)}
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
            <Button type="submit">Create Invoice</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

