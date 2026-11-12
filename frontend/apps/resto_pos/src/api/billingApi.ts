import { apiClient } from '@repo/ui/lib/api-client'
import type {
  BillingPayment,
  BillingInvoice,
  BillingExpense,
  BillingCancelled,
  BillingRefunded,
} from '../types/billing'

export const billingApi = {
  getPayments: async (): Promise<BillingPayment[]> => {
    return apiClient<BillingPayment[]>('/api/billing/payments')
  },
  createPayment: async (payment: Partial<BillingPayment>): Promise<BillingPayment> => {
    return apiClient<BillingPayment>('/api/billing/payments', {
      method: 'POST',
      body: JSON.stringify(payment),
    })
  },
  deletePayment: async (id: string): Promise<{ status: string }> => {
    return apiClient<{ status: string }>(`/api/billing/payments/${id}`, {
      method: 'DELETE',
    })
  },

  getInvoices: async (): Promise<BillingInvoice[]> => {
    return apiClient<BillingInvoice[]>('/api/billing/invoices')
  },
  createInvoice: async (invoice: Partial<BillingInvoice>): Promise<BillingInvoice> => {
    return apiClient<BillingInvoice>('/api/billing/invoices', {
      method: 'POST',
      body: JSON.stringify(invoice),
    })
  },
  deleteInvoice: async (id: string): Promise<{ status: string }> => {
    return apiClient<{ status: string }>(`/api/billing/invoices/${id}`, {
      method: 'DELETE',
    })
  },

  getExpenses: async (): Promise<BillingExpense[]> => {
    return apiClient<BillingExpense[]>('/api/billing/expenses')
  },
  createExpense: async (expense: Partial<BillingExpense>): Promise<BillingExpense> => {
    return apiClient<BillingExpense>('/api/billing/expenses', {
      method: 'POST',
      body: JSON.stringify(expense),
    })
  },
  deleteExpense: async (id: string): Promise<{ status: string }> => {
    return apiClient<{ status: string }>(`/api/billing/expenses/${id}`, {
      method: 'DELETE',
    })
  },

  getCancelled: async (): Promise<BillingCancelled[]> => {
    return apiClient<BillingCancelled[]>('/api/billing/cancelled')
  },

  getRefunds: async (): Promise<BillingRefunded[]> => {
    return apiClient<BillingRefunded[]>('/api/billing/refunds')
  },
  createRefund: async (refund: Partial<BillingRefunded>): Promise<BillingRefunded> => {
    return apiClient<BillingRefunded>('/api/billing/refunds', {
      method: 'POST',
      body: JSON.stringify(refund),
    })
  },
}
