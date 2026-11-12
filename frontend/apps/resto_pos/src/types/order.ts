export type OrderStage =
  | 'Order Confirmation'
  | 'Product Ready'
  | 'Ready For Delivery'
  | 'Product Delivered'

export type OrderPaymentStatus =
  | 'Paid'
  | 'Partially Paid'
  | 'Pending'
  | 'Overdue'

export type OrderType =
  | 'Walk-in'
  | 'Online'
  | 'Studio Session'
  | 'Print & Album'
  | 'Framing'
  | 'Commercial'
  | 'Passport/Visa Photo'
  | 'Custom Merch'

export interface OrderLineItem {
  id: string
  name: string
  quantity: number
  itemCost: number
  subTotal: number
}

export interface OrderExpense {
  id: string
  date: string
  name: string
  amount: number
}

export interface OrderPaymentRecord {
  id: string
  amount: number
  date: string
  method: 'Cash' | 'Card' | 'UPI' | 'Bank Transfer' | 'Other'
  reference?: string
}

export interface Order {
  id: string
  orderNumber: string
  clientName: string
  mobileNo: string
  countryCode?: string
  email?: string
  address?: string
  leadSource?: string

  referenceNo?: string
  orderDate: string
  orderType: OrderType | string
  stage: OrderStage
  assignee: string

  items: OrderLineItem[]
  expenses: OrderExpense[]
  payments?: OrderPaymentRecord[]

  notes?: string
  discount: number
  taxPercent: number
  taxAmount: number
  subTotal: number
  totalAmount: number
  paidAmount: number
  balanceAmount: number
  paymentStatus: OrderPaymentStatus

  createdDate: string
  relatedProjectId?: string
}

