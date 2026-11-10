import * as React from 'react'
import {
  Button,
  Badge,
  DataTable,
  cn,
  toast,
} from '@repo/ui'
import { createColumnHelper } from '@tanstack/react-table'
import type {
  DataTableFeatures,
  DataTableFilterConfig,
} from '@repo/ui'
import {
  IconDownload,
  IconPlus,
} from '@tabler/icons-react'
import { isWithinInterval, parse } from 'date-fns'

import type {
  BillingTab,
  BillingPayment,
  BillingInvoice,
  BillingExpense,
  BillingCancelled,
  BillingRefunded,
} from '../types/billing'
import { billingApi } from '../api/crmApi'
import { BillingHero } from '../components/BillingHero'
import { BillingDateRangePicker, type DateRange } from '../components/BillingDateRangePicker'
import { AddPaymentModal } from '../components/AddPaymentModal'
import { AddInvoiceModal } from '../components/AddInvoiceModal'
import { AddExpenseModal } from '../components/AddExpenseModal'
import { AddRefundModal } from '../components/AddRefundModal'

// Column Helpers for each table
const paymentColHelper = createColumnHelper<DataTableFeatures, BillingPayment>()
const invoiceColHelper = createColumnHelper<DataTableFeatures, BillingInvoice>()
const expenseColHelper = createColumnHelper<DataTableFeatures, BillingExpense>()
const cancelledColHelper = createColumnHelper<DataTableFeatures, BillingCancelled>()
const refundedColHelper = createColumnHelper<DataTableFeatures, BillingRefunded>()

// Safe date parser for varied formats like "Nov 01, 2026 12:00 PM" or "Nov 01, 2026"
function parseSafeDate(dateStr: string): Date | null {
  if (!dateStr) return null
  const d = new Date(dateStr)
  if (!isNaN(d.getTime())) return d
  try {
    const parsed = parse(dateStr, 'MMM dd, yyyy hh:mm a', new Date())
    if (!isNaN(parsed.getTime())) return parsed
  } catch {
    // continue
  }
  return null
}

export function BillingsPage() {
  const [activeTab, setActiveTab] = React.useState<BillingTab>('payments')

  // Date Range filter state
  const [dateRange, setDateRange] = React.useState<DateRange>({
    startDate: new Date(2000, 0, 1),
    endDate: new Date(2026, 11, 31),
  })

  // Data states
  const [payments, setPayments] = React.useState<BillingPayment[]>([])
  const [invoices, setInvoices] = React.useState<BillingInvoice[]>([])
  const [expenses, setExpenses] = React.useState<BillingExpense[]>([])
  const [cancelled, setCancelled] = React.useState<BillingCancelled[]>([])
  const [refunded, setRefunded] = React.useState<BillingRefunded[]>([])

  const loadBillingData = React.useCallback(() => {
    billingApi.getPayments().then(setPayments).catch((err) => console.error(err))
    billingApi.getInvoices().then(setInvoices).catch((err) => console.error(err))
    billingApi.getExpenses().then(setExpenses).catch((err) => console.error(err))
    billingApi.getCancelled().then(setCancelled).catch((err) => console.error(err))
    billingApi.getRefunds().then(setRefunded).catch((err) => console.error(err))
  }, [])

  React.useEffect(() => {
    loadBillingData()
  }, [loadBillingData])

  // Modal dialog states
  const [isPaymentModalOpen, setIsPaymentModalOpen] = React.useState(false)
  const [isInvoiceModalOpen, setIsInvoiceModalOpen] = React.useState(false)
  const [isExpenseModalOpen, setIsExpenseModalOpen] = React.useState(false)
  const [isRefundModalOpen, setIsRefundModalOpen] = React.useState(false)

  // Filtering helpers by active date range
  const isDateInRange = React.useCallback(
    (dateStr: string) => {
      const d = parseSafeDate(dateStr)
      if (!d) return true
      return isWithinInterval(d, {
        start: dateRange.startDate,
        end: dateRange.endDate,
      })
    },
    [dateRange]
  )

  const filteredPayments = React.useMemo(() => {
    return payments.filter((p) => isDateInRange(p.paymentDate))
  }, [payments, isDateInRange])

  const filteredInvoices = React.useMemo(() => {
    return invoices.filter((i) => isDateInRange(i.invoiceDate))
  }, [invoices, isDateInRange])

  const filteredExpenses = React.useMemo(() => {
    return expenses.filter((e) => isDateInRange(e.date))
  }, [expenses, isDateInRange])

  const filteredCancelled = React.useMemo(() => {
    return cancelled.filter((c) => isDateInRange(c.cancellationDate))
  }, [cancelled, isDateInRange])

  const filteredRefunded = React.useMemo(() => {
    return refunded.filter((r) => isDateInRange(r.refundDate))
  }, [refunded, isDateInRange])

  // Top Hero Configuration depending on Active Tab
  const heroConfig = React.useMemo(() => {
    switch (activeTab) {
      case 'payments': {
        const total = filteredPayments.reduce((acc, p) => acc + p.paymentAmount, 0)
        return {
          cards: [
            {
              label: 'Total Payment Amount',
              value: `₹${total.toLocaleString('en-IN')}`,
              subtitle: `from ${filteredPayments.length} Payment${filteredPayments.length === 1 ? '' : 's'}`,
              colorClass: 'text-emerald-600 dark:text-emerald-400',
            },
          ],
          chartTitle: 'All Payments',
          chartData: {
            labels: ['Nov 2026'],
            values: [total],
          },
          chartColor: '#2563EB',
        }
      }
      case 'invoices': {
        const totalInvoice = filteredInvoices.reduce((acc, i) => acc + i.totalAmount, 0)
        const totalPaid = filteredInvoices.reduce((acc, i) => acc + i.paidAmount, 0)
        const totalDue = filteredInvoices.reduce((acc, i) => acc + i.pendingAmount, 0)
        return {
          cards: [
            {
              label: 'Total Invoice Amount',
              value: `₹${totalInvoice.toLocaleString('en-IN')}`,
              subtitle: `${filteredInvoices.length} Invoices issued`,
              colorClass: 'text-blue-600 dark:text-blue-400',
            },
            {
              label: 'Paid Amount',
              value: `₹${totalPaid.toLocaleString('en-IN')}`,
              subtitle: 'Received payments',
              colorClass: 'text-emerald-600 dark:text-emerald-400',
            },
            {
              label: 'Due Amount',
              value: `₹${totalDue.toLocaleString('en-IN')}`,
              subtitle: 'Pending collection',
              colorClass: 'text-amber-600 dark:text-amber-400',
            },
          ],
          chartTitle: 'Invoices & Realization',
          chartData: {
            labels: ['Total Billed', 'Paid', 'Due'],
            values: [totalInvoice, totalPaid, totalDue],
          },
          chartColor: '#3B82F6',
        }
      }
      case 'expenses': {
        const totalExpense = filteredExpenses.reduce((acc, e) => acc + e.amount, 0)
        return {
          cards: [
            {
              label: 'Expense Amount',
              value: `₹${totalExpense.toLocaleString('en-IN')}`,
              subtitle: `from ${filteredExpenses.length} Expense item${filteredExpenses.length === 1 ? '' : 's'}`,
              colorClass: 'text-rose-600 dark:text-rose-400',
            },
          ],
          chartTitle: 'All Expenses',
          chartData: {
            labels: ['Nov 2026'],
            values: [totalExpense],
          },
          chartColor: '#E11D48',
        }
      }
      case 'cancelled': {
        const totalCancelled = filteredCancelled.reduce((acc, c) => acc + c.paymentAmount, 0)
        return {
          cards: [
            {
              label: 'Cancelled Amount',
              value: `₹${totalCancelled.toLocaleString('en-IN')}`,
              subtitle: `from ${filteredCancelled.length} Cancelled booking${filteredCancelled.length === 1 ? '' : 's'}`,
              colorClass: 'text-orange-600 dark:text-orange-400',
            },
          ],
          chartTitle: 'Cancelled Value',
          chartData: {
            labels: ['Nov 2026'],
            values: [totalCancelled],
          },
          chartColor: '#EA580C',
        }
      }
      case 'refunded': {
        const totalRefunded = filteredRefunded.reduce((acc, r) => acc + r.refundAmount, 0)
        return {
          cards: [
            {
              label: 'Refunded Amount',
              value: `₹${totalRefunded.toLocaleString('en-IN')}`,
              subtitle: `from ${filteredRefunded.length} Processed refund${filteredRefunded.length === 1 ? '' : 's'}`,
              colorClass: 'text-violet-600 dark:text-violet-400',
            },
          ],
          chartTitle: 'Refunded Outflow',
          chartData: {
            labels: ['Nov 2026'],
            values: [totalRefunded],
          },
          chartColor: '#8B5CF6',
        }
      }
    }
  }, [
    activeTab,
    filteredPayments,
    filteredInvoices,
    filteredExpenses,
    filteredCancelled,
    filteredRefunded,
  ])

  // --- 1. PAYMENTS COLUMNS ---
  // Table column headers: ID, Payment Date, Client Name, Type, Payment Amount, Payment Status, Download, Created On
  const paymentColumns = React.useMemo(
    () =>
      paymentColHelper.columns([
      paymentColHelper.accessor('id', {
        header: 'ID',
        meta: { align: 'left', headerClassName: 'text-xs font-bold' },
        cell: (info) => (
          <span className="font-semibold text-xs text-foreground tracking-wide">
            {info.getValue()}
          </span>
        ),
      }),
      paymentColHelper.accessor('paymentDate', {
        header: 'Payment Date',
        meta: { align: 'left' },
        cell: (info) => (
          <span className="text-xs text-muted-foreground whitespace-nowrap">
            {info.getValue()}
          </span>
        ),
      }),
      paymentColHelper.accessor('clientName', {
        header: 'Client Name',
        meta: { align: 'left' },
        cell: (info) => {
          const row = info.row.original
          return (
            <div className="flex items-center gap-2">
              <div className="h-7 w-7 rounded-full bg-slate-900 text-white dark:bg-slate-100 dark:text-slate-900 flex items-center justify-center text-[10px] font-bold shrink-0 shadow-xs">
                {row.clientInitials || row.clientName.slice(0, 2).toUpperCase()}
              </div>
              <span className="font-semibold text-xs text-foreground">
                {info.getValue()}
              </span>
            </div>
          )
        },
      }),
      paymentColHelper.accessor('type', {
        header: 'Type',
        meta: { align: 'left' },
        cell: (info) => (
          <Badge variant="outline" className="text-xs font-medium border-border/80">
            {info.getValue()}
          </Badge>
        ),
      }),
      paymentColHelper.accessor('paymentAmount', {
        header: 'Payment Amount',
        meta: { align: 'right' },
        cell: (info) => (
          <span className="font-bold text-xs text-foreground">
            ₹{info.getValue().toLocaleString('en-IN')}
          </span>
        ),
      }),
      paymentColHelper.accessor('paymentStatus', {
        header: 'Payment Status',
        meta: { align: 'center' },
        cell: (info) => {
          const st = info.getValue()
          const colorMap: Record<string, string> = {
            Completed: 'bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-950/40 dark:text-emerald-400 dark:border-emerald-800',
            'Payment Due': 'bg-amber-50 text-amber-700 border-amber-200 dark:bg-amber-950/40 dark:text-amber-400 dark:border-amber-800',
            Pending: 'bg-blue-50 text-blue-700 border-blue-200 dark:bg-blue-950/40 dark:text-blue-400 dark:border-blue-800',
            Failed: 'bg-rose-50 text-rose-700 border-rose-200 dark:bg-rose-950/40 dark:text-rose-400 dark:border-rose-800',
          }
          return (
            <Badge
              variant="outline"
              className={cn('text-[11px] font-semibold px-2 py-0.5 rounded-full', colorMap[st] || 'bg-muted')}
            >
              {st}
            </Badge>
          )
        },
      }),
      paymentColHelper.display({
        id: 'download',
        header: 'Download',
        meta: { align: 'center' },
        cell: (info) => (
          <Button
            variant="ghost"
            size="icon"
            className="h-7 w-7 text-muted-foreground hover:text-foreground"
            onClick={() => toast.success(`Receipt for ${info.row.original.id} downloaded`)}
            title="Download Payment Receipt"
          >
            <IconDownload className="h-4 w-4" />
          </Button>
        ),
      }),
      paymentColHelper.accessor('createdOn', {
        header: 'Created On',
        meta: { align: 'right' },
        cell: (info) => (
          <span className="text-xs text-muted-foreground whitespace-nowrap">
            {info.getValue()}
          </span>
        ),
      }),
    ]),
    []
  )

  // --- 2. INVOICE COLUMNS ---
  // Table column headers: ID, Invoice Name, Client Name, Projects Name, Total Amount, Payment Status, Paid Amount, Pending Amount, Invoice Date, Created Date
  const invoiceColumns = React.useMemo(
    () =>
      invoiceColHelper.columns([
      invoiceColHelper.accessor('id', {
        header: 'ID',
        meta: { align: 'left', headerClassName: 'text-xs font-bold' },
        cell: (info) => (
          <span className="font-semibold text-xs text-foreground tracking-wide">
            {info.getValue()}
          </span>
        ),
      }),
      invoiceColHelper.accessor('invoiceName', {
        header: 'Invoice Name',
        meta: { align: 'left' },
        cell: (info) => (
          <span className="font-medium text-xs text-foreground max-w-[200px] truncate block" title={info.getValue()}>
            {info.getValue()}
          </span>
        ),
      }),
      invoiceColHelper.accessor('clientName', {
        header: 'Client Name',
        meta: { align: 'left' },
        cell: (info) => {
          const row = info.row.original
          return (
            <div className="flex items-center gap-2">
              <div className="h-7 w-7 rounded-full bg-slate-900 text-white dark:bg-slate-100 dark:text-slate-900 flex items-center justify-center text-[10px] font-bold shrink-0">
                {row.clientInitials || row.clientName.slice(0, 2).toUpperCase()}
              </div>
              <span className="font-semibold text-xs text-foreground">
                {info.getValue()}
              </span>
            </div>
          )
        },
      }),
      invoiceColHelper.accessor('projectsName', {
        header: 'Projects Name',
        meta: { align: 'left' },
        cell: (info) => (
          <span className="text-xs text-muted-foreground font-medium">
            {info.getValue()}
          </span>
        ),
      }),
      invoiceColHelper.accessor('totalAmount', {
        header: 'Total Amount',
        meta: { align: 'right' },
        cell: (info) => (
          <span className="font-bold text-xs text-foreground">
            ₹{info.getValue().toLocaleString('en-IN')}
          </span>
        ),
      }),
      invoiceColHelper.accessor('paymentStatus', {
        header: 'Payment Status',
        meta: { align: 'center' },
        cell: (info) => {
          const st = info.getValue()
          const colorMap: Record<string, string> = {
            Paid: 'bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-950/40 dark:text-emerald-400 dark:border-emerald-800',
            'Partially Paid': 'bg-amber-50 text-amber-700 border-amber-200 dark:bg-amber-950/40 dark:text-amber-400 dark:border-amber-800',
            Unpaid: 'bg-rose-50 text-rose-700 border-rose-200 dark:bg-rose-950/40 dark:text-rose-400 dark:border-rose-800',
            Overdue: 'bg-red-50 text-red-700 border-red-200 dark:bg-red-950/40 dark:text-red-400 dark:border-red-800',
            Draft: 'bg-muted text-muted-foreground border-border',
          }
          return (
            <Badge
              variant="outline"
              className={cn('text-[11px] font-semibold px-2 py-0.5 rounded-full', colorMap[st] || 'bg-muted')}
            >
              {st}
            </Badge>
          )
        },
      }),
      invoiceColHelper.accessor('paidAmount', {
        header: 'Paid Amount',
        meta: { align: 'right' },
        cell: (info) => (
          <span className="font-semibold text-xs text-emerald-600 dark:text-emerald-400">
            ₹{info.getValue().toLocaleString('en-IN')}
          </span>
        ),
      }),
      invoiceColHelper.accessor('pendingAmount', {
        header: 'Pending Amount',
        meta: { align: 'right' },
        cell: (info) => (
          <span className="font-semibold text-xs text-amber-600 dark:text-amber-400">
            ₹{info.getValue().toLocaleString('en-IN')}
          </span>
        ),
      }),
      invoiceColHelper.accessor('invoiceDate', {
        header: 'Invoice Date',
        meta: { align: 'left' },
        cell: (info) => (
          <span className="text-xs text-muted-foreground whitespace-nowrap">
            {info.getValue()}
          </span>
        ),
      }),
      invoiceColHelper.accessor('createdDate', {
        header: 'Created Date',
        meta: { align: 'right' },
        cell: (info) => (
          <span className="text-xs text-muted-foreground whitespace-nowrap">
            {info.getValue()}
          </span>
        ),
      }),
    ]),
    []
  )

  // --- 3. EXPENSE COLUMNS ---
  // Table column headers: ID, Expense Name, Date, Amount, Type, Created Date
  const expenseColumns = React.useMemo(
    () =>
      expenseColHelper.columns([
      expenseColHelper.accessor('id', {
        header: 'ID',
        meta: { align: 'left', headerClassName: 'text-xs font-bold' },
        cell: (info) => (
          <span className="font-semibold text-xs text-foreground tracking-wide">
            {info.getValue()}
          </span>
        ),
      }),
      expenseColHelper.accessor('expenseName', {
        header: 'Expense Name',
        meta: { align: 'left' },
        cell: (info) => (
          <div className="flex flex-col">
            <span className="font-semibold text-xs text-foreground">
              {info.getValue()}
            </span>
            {info.row.original.vendorOrClient && (
              <span className="text-[11px] text-muted-foreground">
                Paid to: {info.row.original.vendorOrClient}
              </span>
            )}
          </div>
        ),
      }),
      expenseColHelper.accessor('date', {
        header: 'Date',
        meta: { align: 'left' },
        cell: (info) => (
          <span className="text-xs text-muted-foreground whitespace-nowrap">
            {info.getValue()}
          </span>
        ),
      }),
      expenseColHelper.accessor('amount', {
        header: 'Amount',
        meta: { align: 'right' },
        cell: (info) => (
          <span className="font-bold text-xs text-rose-600 dark:text-rose-400">
            ₹{info.getValue().toLocaleString('en-IN')}
          </span>
        ),
      }),
      expenseColHelper.accessor('type', {
        header: 'Type',
        meta: { align: 'left' },
        cell: (info) => (
          <Badge variant="outline" className="text-xs font-medium border-border/80">
            {info.getValue()}
          </Badge>
        ),
      }),
      expenseColHelper.accessor('createdDate', {
        header: 'Created Date',
        meta: { align: 'right' },
        cell: (info) => (
          <span className="text-xs text-muted-foreground whitespace-nowrap">
            {info.getValue()}
          </span>
        ),
      }),
    ]),
    []
  )

  // --- 4. CANCELLED COLUMNS ---
  // Table column headers: Action, Cancellation Date, Payment Date, Type, Client Name, Projects Name, Payment Amount, Created On
  const cancelledColumns = React.useMemo(
    () =>
      cancelledColHelper.columns([
      cancelledColHelper.accessor('action', {
        header: 'Action',
        meta: { align: 'left' },
        cell: (info) => (
          <Badge
            variant="outline"
            className="text-[11px] font-semibold bg-rose-50 text-rose-700 border-rose-200 dark:bg-rose-950/40 dark:text-rose-400 dark:border-rose-800"
          >
            {info.getValue()}
          </Badge>
        ),
      }),
      cancelledColHelper.accessor('cancellationDate', {
        header: 'Cancellation Date',
        meta: { align: 'left' },
        cell: (info) => (
          <span className="text-xs text-muted-foreground whitespace-nowrap">
            {info.getValue()}
          </span>
        ),
      }),
      cancelledColHelper.accessor('paymentDate', {
        header: 'Payment Date',
        meta: { align: 'left' },
        cell: (info) => (
          <span className="text-xs text-muted-foreground whitespace-nowrap">
            {info.getValue()}
          </span>
        ),
      }),
      cancelledColHelper.accessor('type', {
        header: 'Type',
        meta: { align: 'left' },
        cell: (info) => (
          <Badge variant="outline" className="text-xs font-medium border-border/80">
            {info.getValue()}
          </Badge>
        ),
      }),
      cancelledColHelper.accessor('clientName', {
        header: 'Client Name',
        meta: { align: 'left' },
        cell: (info) => {
          const row = info.row.original
          return (
            <div className="flex items-center gap-2">
              <div className="h-7 w-7 rounded-full bg-slate-900 text-white dark:bg-slate-100 dark:text-slate-900 flex items-center justify-center text-[10px] font-bold shrink-0">
                {row.clientInitials || row.clientName.slice(0, 2).toUpperCase()}
              </div>
              <span className="font-semibold text-xs text-foreground">
                {info.getValue()}
              </span>
            </div>
          )
        },
      }),
      cancelledColHelper.accessor('projectsName', {
        header: 'Projects Name',
        meta: { align: 'left' },
        cell: (info) => (
          <span className="text-xs text-muted-foreground font-medium">
            {info.getValue()}
          </span>
        ),
      }),
      cancelledColHelper.accessor('paymentAmount', {
        header: 'Payment Amount',
        meta: { align: 'right' },
        cell: (info) => (
          <span className="font-bold text-xs text-foreground">
            ₹{info.getValue().toLocaleString('en-IN')}
          </span>
        ),
      }),
      cancelledColHelper.accessor('createdOn', {
        header: 'Created On',
        meta: { align: 'right' },
        cell: (info) => (
          <span className="text-xs text-muted-foreground whitespace-nowrap">
            {info.getValue()}
          </span>
        ),
      }),
    ]),
    []
  )

  // --- 5. REFUNDED COLUMNS ---
  // Table column headers: ID, Invoice Name, Refund Date, Client Name, Projects Name, Payment Amount, Cancellation Charges, Refund Amount, Note, Status, Created On
  const refundedColumns = React.useMemo(
    () =>
      refundedColHelper.columns([
      refundedColHelper.accessor('id', {
        header: 'ID',
        meta: { align: 'left', headerClassName: 'text-xs font-bold' },
        cell: (info) => (
          <span className="font-semibold text-xs text-foreground tracking-wide">
            {info.getValue()}
          </span>
        ),
      }),
      refundedColHelper.accessor('invoiceName', {
        header: 'Invoice Name',
        meta: { align: 'left' },
        cell: (info) => (
          <span className="font-medium text-xs text-foreground max-w-[180px] truncate block" title={info.getValue()}>
            {info.getValue()}
          </span>
        ),
      }),
      refundedColHelper.accessor('refundDate', {
        header: 'Refund Date',
        meta: { align: 'left' },
        cell: (info) => (
          <span className="text-xs text-muted-foreground whitespace-nowrap">
            {info.getValue()}
          </span>
        ),
      }),
      refundedColHelper.accessor('clientName', {
        header: 'Client Name',
        meta: { align: 'left' },
        cell: (info) => {
          const row = info.row.original
          return (
            <div className="flex items-center gap-2">
              <div className="h-7 w-7 rounded-full bg-slate-900 text-white dark:bg-slate-100 dark:text-slate-900 flex items-center justify-center text-[10px] font-bold shrink-0">
                {row.clientInitials || row.clientName.slice(0, 2).toUpperCase()}
              </div>
              <span className="font-semibold text-xs text-foreground">
                {info.getValue()}
              </span>
            </div>
          )
        },
      }),
      refundedColHelper.accessor('projectsName', {
        header: 'Projects Name',
        meta: { align: 'left' },
        cell: (info) => (
          <span className="text-xs text-muted-foreground font-medium">
            {info.getValue()}
          </span>
        ),
      }),
      refundedColHelper.accessor('paymentAmount', {
        header: 'Payment Amount',
        meta: { align: 'right' },
        cell: (info) => (
          <span className="text-xs text-muted-foreground font-medium">
            ₹{info.getValue().toLocaleString('en-IN')}
          </span>
        ),
      }),
      refundedColHelper.accessor('cancellationCharges', {
        header: 'Cancellation Charges',
        meta: { align: 'right' },
        cell: (info) => (
          <span className="text-xs text-rose-600 dark:text-rose-400 font-semibold">
            ₹{info.getValue().toLocaleString('en-IN')}
          </span>
        ),
      }),
      refundedColHelper.accessor('refundAmount', {
        header: 'Refund Amount',
        meta: { align: 'right' },
        cell: (info) => (
          <span className="font-bold text-xs text-emerald-600 dark:text-emerald-400">
            ₹{info.getValue().toLocaleString('en-IN')}
          </span>
        ),
      }),
      refundedColHelper.accessor('note', {
        header: 'Note',
        meta: { align: 'left' },
        cell: (info) => (
          <span className="text-xs text-muted-foreground italic max-w-[200px] truncate block" title={info.getValue()}>
            {info.getValue() || '—'}
          </span>
        ),
      }),
      refundedColHelper.accessor('status', {
        header: 'Status',
        meta: { align: 'center' },
        cell: (info) => {
          const st = info.getValue()
          const colorMap: Record<string, string> = {
            Completed: 'bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-950/40 dark:text-emerald-400 dark:border-emerald-800',
            Processing: 'bg-blue-50 text-blue-700 border-blue-200 dark:bg-blue-950/40 dark:text-blue-400 dark:border-blue-800',
            Failed: 'bg-rose-50 text-rose-700 border-rose-200 dark:bg-rose-950/40 dark:text-rose-400 dark:border-rose-800',
          }
          return (
            <Badge
              variant="outline"
              className={cn('text-[11px] font-semibold px-2 py-0.5 rounded-full', colorMap[st] || 'bg-muted')}
            >
              {st}
            </Badge>
          )
        },
      }),
      refundedColHelper.accessor('createdOn', {
        header: 'Created On',
        meta: { align: 'right' },
        cell: (info) => (
          <span className="text-xs text-muted-foreground whitespace-nowrap">
            {info.getValue()}
          </span>
        ),
      }),
    ]),
    []
  )

  // Filters configurations for DataTable
  const paymentFilterConfigs: DataTableFilterConfig[] = [
    {
      columnId: 'paymentStatus',
      title: 'Payment Status',
      type: 'select',
      options: [
        { label: 'Completed', value: 'Completed' },
        { label: 'Payment Due', value: 'Payment Due' },
        { label: 'Pending', value: 'Pending' },
        { label: 'Failed', value: 'Failed' },
      ],
    },
    {
      columnId: 'type',
      title: 'Payment Type',
      type: 'select',
      options: [
        { label: 'Invoice', value: 'Invoice' },
        { label: 'UPI Transfer', value: 'UPI Transfer' },
        { label: 'Bank Transfer (NEFT)', value: 'Bank Transfer (NEFT)' },
        { label: 'Credit Card', value: 'Credit Card' },
      ],
    },
  ]

  const invoiceFilterConfigs: DataTableFilterConfig[] = [
    {
      columnId: 'paymentStatus',
      title: 'Invoice Status',
      type: 'select',
      options: [
        { label: 'Paid', value: 'Paid' },
        { label: 'Partially Paid', value: 'Partially Paid' },
        { label: 'Unpaid', value: 'Unpaid' },
        { label: 'Overdue', value: 'Overdue' },
      ],
    },
  ]

  // Add Button based on active tab
  const renderAddButton = () => {
    switch (activeTab) {
      case 'payments':
        return (
          <Button size="sm" onClick={() => setIsPaymentModalOpen(true)}>
            <IconPlus className="h-4 w-4 mr-1.5" />
            Record Payment
          </Button>
        )
      case 'invoices':
        return (
          <Button size="sm" onClick={() => setIsInvoiceModalOpen(true)}>
            <IconPlus className="h-4 w-4 mr-1.5" />
            Create Invoice
          </Button>
        )
      case 'expenses':
        return (
          <Button size="sm" onClick={() => setIsExpenseModalOpen(true)}>
            <IconPlus className="h-4 w-4 mr-1.5" />
            Add Expense
          </Button>
        )
      case 'refunded':
        return (
          <Button size="sm" onClick={() => setIsRefundModalOpen(true)}>
            <IconPlus className="h-4 w-4 mr-1.5" />
            Issue Refund
          </Button>
        )
      default:
        return null
    }
  }

  const datePickerPrefix = React.useMemo(() => {
    switch (activeTab) {
      case 'payments':
        return 'Payments Date Range'
      case 'invoices':
        return 'Invoices Date Range'
      case 'expenses':
        return 'Expenses Date Range'
      case 'cancelled':
        return 'Cancelled Date Range'
      case 'refunded':
        return 'Refunded Date Range'
    }
  }, [activeTab])

  return (
    <div className="space-y-6">
      {/* Top Hero Section: Cards on left + Chart on right + Date Range selector */}
      <BillingHero
        cards={heroConfig.cards}
        chartTitle={heroConfig.chartTitle}
        chartData={heroConfig.chartData}
        chartColor={heroConfig.chartColor}
        dateRangeFilterNode={
          <BillingDateRangePicker
            labelPrefix={datePickerPrefix}
            dateRange={dateRange}
            onChange={setDateRange}
            onClear={() =>
              setDateRange({
                startDate: new Date(2000, 0, 1),
                endDate: new Date(2026, 11, 31),
              })
            }
          />
        }
      />

      {/* Tabs Navigation (Payments, Invoice, Expense, Cancelled, Refunded) */}
      <div className="flex items-center justify-between border-b pb-1">
        <div className="flex items-center gap-1 sm:gap-2 overflow-x-auto">
          <button
            type="button"
            onClick={() => setActiveTab('payments')}
            className={cn(
              'px-4 py-2 text-xs font-bold transition-all relative whitespace-nowrap cursor-pointer rounded-t-md',
              activeTab === 'payments'
                ? 'text-primary border-b-2 border-primary bg-primary/5 font-extrabold'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted/30'
            )}
          >
            Payments
            <span className="ml-1.5 text-[10px] font-semibold text-muted-foreground">
              ({filteredPayments.length})
            </span>
          </button>

          <button
            type="button"
            onClick={() => setActiveTab('invoices')}
            className={cn(
              'px-4 py-2 text-xs font-bold transition-all relative whitespace-nowrap cursor-pointer rounded-t-md',
              activeTab === 'invoices'
                ? 'text-primary border-b-2 border-primary bg-primary/5 font-extrabold'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted/30'
            )}
          >
            Invoice
            <span className="ml-1.5 text-[10px] font-semibold text-muted-foreground">
              ({filteredInvoices.length})
            </span>
          </button>

          <button
            type="button"
            onClick={() => setActiveTab('expenses')}
            className={cn(
              'px-4 py-2 text-xs font-bold transition-all relative whitespace-nowrap cursor-pointer rounded-t-md',
              activeTab === 'expenses'
                ? 'text-primary border-b-2 border-primary bg-primary/5 font-extrabold'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted/30'
            )}
          >
            Expense
            <span className="ml-1.5 text-[10px] font-semibold text-muted-foreground">
              ({filteredExpenses.length})
            </span>
          </button>

          <button
            type="button"
            onClick={() => setActiveTab('cancelled')}
            className={cn(
              'px-4 py-2 text-xs font-bold transition-all relative whitespace-nowrap cursor-pointer rounded-t-md',
              activeTab === 'cancelled'
                ? 'text-primary border-b-2 border-primary bg-primary/5 font-extrabold'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted/30'
            )}
          >
            Cancelled
            <span className="ml-1.5 text-[10px] font-semibold text-muted-foreground">
              ({filteredCancelled.length})
            </span>
          </button>

          <button
            type="button"
            onClick={() => setActiveTab('refunded')}
            className={cn(
              'px-4 py-2 text-xs font-bold transition-all relative whitespace-nowrap cursor-pointer rounded-t-md',
              activeTab === 'refunded'
                ? 'text-primary border-b-2 border-primary bg-primary/5 font-extrabold'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted/30'
            )}
          >
            Refunded
            <span className="ml-1.5 text-[10px] font-semibold text-muted-foreground">
              ({filteredRefunded.length})
            </span>
          </button>
        </div>

        <div>{renderAddButton()}</div>
      </div>

      {/* Table Section */}
      <div className="rounded-xl border bg-card shadow-xs overflow-hidden">
        {activeTab === 'payments' && (
          <DataTable
            data={filteredPayments}
            columns={paymentColumns}
            filterConfigs={paymentFilterConfigs}
            emptyTitle="No payments found"
            emptyDescription="No payments found for selected date range"
          />
        )}

        {activeTab === 'invoices' && (
          <DataTable
            data={filteredInvoices}
            columns={invoiceColumns}
            filterConfigs={invoiceFilterConfigs}
            emptyTitle="No invoices found"
            emptyDescription="No invoices found for selected date range"
          />
        )}

        {activeTab === 'expenses' && (
          <DataTable
            data={filteredExpenses}
            columns={expenseColumns}
            emptyTitle="No expenses found"
            emptyDescription="No expenses found for selected date range"
          />
        )}

        {activeTab === 'cancelled' && (
          <DataTable
            data={filteredCancelled}
            columns={cancelledColumns}
            emptyTitle="No cancelled bookings found"
            emptyDescription="No cancelled bookings recorded"
          />
        )}

        {activeTab === 'refunded' && (
          <DataTable
            data={filteredRefunded}
            columns={refundedColumns}
            emptyTitle="No refunds found"
            emptyDescription="No refunds recorded"
          />
        )}
      </div>

      {/* Action Modals */}
      <AddPaymentModal
        open={isPaymentModalOpen}
        onOpenChange={setIsPaymentModalOpen}
        onAddPayment={(newPay) => {
          billingApi.createPayment(newPay).then((created) => {
            setPayments((prev) => [created, ...prev])
            toast.success('Payment recorded successfully')
          }).catch((err) => toast.error(err.message || 'Failed to record payment'))
        }}
      />
      <AddInvoiceModal
        open={isInvoiceModalOpen}
        onOpenChange={setIsInvoiceModalOpen}
        onAddInvoice={(newInv) => {
          billingApi.createInvoice(newInv).then((created) => {
            setInvoices((prev) => [created, ...prev])
            toast.success('Invoice created successfully')
          }).catch((err) => toast.error(err.message || 'Failed to create invoice'))
        }}
      />
      <AddExpenseModal
        open={isExpenseModalOpen}
        onOpenChange={setIsExpenseModalOpen}
        onAddExpense={(newExp) => {
          billingApi.createExpense(newExp).then((created) => {
            setExpenses((prev) => [created, ...prev])
            toast.success('Expense added successfully')
          }).catch((err) => toast.error(err.message || 'Failed to add expense'))
        }}
      />
      <AddRefundModal
        open={isRefundModalOpen}
        onOpenChange={setIsRefundModalOpen}
        onAddRefund={(newRef) => {
          billingApi.createRefund(newRef).then((created) => {
            setRefunded((prev) => [created, ...prev])
            toast.success('Refund issued successfully')
          }).catch((err) => toast.error(err.message || 'Failed to issue refund'))
        }}
      />
    </div>
  )
}
