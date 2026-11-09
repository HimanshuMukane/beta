export type BillingTab = 'payments' | 'invoices' | 'expenses' | 'cancelled' | 'refunded'

export type PaymentStatus = 'Completed' | 'Pending' | 'Failed' | 'Payment Due' | 'Partially Paid'
export type InvoiceStatus = 'Paid' | 'Partially Paid' | 'Unpaid' | 'Overdue' | 'Draft'
export type RefundStatus = 'Completed' | 'Processing' | 'Failed'

export interface BillingPayment {
  id: string
  paymentDate: string
  clientName: string
  clientAvatar?: string
  clientInitials?: string
  type: string
  paymentAmount: number
  paymentStatus: PaymentStatus
  downloadUrl?: string
  createdOn: string
  projectName?: string
}

export interface BillingInvoice {
  id: string
  invoiceName: string
  clientName: string
  clientAvatar?: string
  clientInitials?: string
  projectsName: string
  totalAmount: number
  paymentStatus: InvoiceStatus
  paidAmount: number
  pendingAmount: number
  invoiceDate: string
  createdDate: string
  downloadUrl?: string
}

export interface BillingExpense {
  id: string
  expenseName: string
  date: string
  amount: number
  type: string
  createdDate: string
  vendorOrClient?: string
  paymentMethod?: string
}

export interface BillingCancelled {
  id: string
  action: string
  cancellationDate: string
  paymentDate: string
  type: string
  clientName: string
  clientAvatar?: string
  clientInitials?: string
  projectsName: string
  paymentAmount: number
  createdOn: string
  reason?: string
}

export interface BillingRefunded {
  id: string
  invoiceName: string
  refundDate: string
  clientName: string
  clientAvatar?: string
  clientInitials?: string
  projectsName: string
  paymentAmount: number
  cancellationCharges: number
  refundAmount: number
  note: string
  status: RefundStatus
  createdOn: string
}

export interface DateRangePreset {
  label: string
  startDate: Date
  endDate: Date
}

