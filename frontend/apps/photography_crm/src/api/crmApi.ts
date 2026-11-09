import { apiClient } from '@repo/ui'
import type { Lead, PipelineStage, LeadStage } from '../types/lead'
import type { Project, ProjectStage, ProjectStatus } from '../types/project'
import type { ProjectSession, SessionStatus } from '../types/session'
import type { Order, OrderStage } from '../types/order'
import type {
  BillingPayment,
  BillingInvoice,
  BillingExpense,
  BillingCancelled,
  BillingRefunded,
} from '../types/billing'

// Leads API
export const leadsApi = {
  list: async (): Promise<Lead[]> => {
    return apiClient<Lead[]>('/api/leads')
  },
  create: async (lead: Partial<Lead>): Promise<Lead> => {
    return apiClient<Lead>('/api/leads', {
      method: 'POST',
      body: JSON.stringify(lead),
    })
  },
  update: async (id: string, lead: Partial<Lead>): Promise<Lead> => {
    return apiClient<Lead>(`/api/leads/${id}`, {
      method: 'PUT',
      body: JSON.stringify(lead),
    })
  },
  updateStage: async (
    id: string,
    pipelineStage: PipelineStage,
    leadStage: LeadStage
  ): Promise<{ status: string }> => {
    return apiClient<{ status: string }>(`/api/leads/${id}/stage`, {
      method: 'PATCH',
      body: JSON.stringify({ pipelineStage, leadStage }),
    })
  },
  convert: async (
    id: string
  ): Promise<{ status: string; clientId: string; projectId: string; projectName: string }> => {
    return apiClient(`/api/leads/${id}/convert`, {
      method: 'POST',
    })
  },
  delete: async (id: string): Promise<{ status: string }> => {
    return apiClient<{ status: string }>(`/api/leads/${id}`, {
      method: 'DELETE',
    })
  },
}

// Projects API
export const projectsApi = {
  list: async (): Promise<Project[]> => {
    return apiClient<Project[]>('/api/projects')
  },
  get: async (id: string): Promise<Project> => {
    return apiClient<Project>(`/api/projects/${id}`)
  },
  create: async (project: Partial<Project>): Promise<Project> => {
    return apiClient<Project>('/api/projects', {
      method: 'POST',
      body: JSON.stringify(project),
    })
  },
  update: async (id: string, project: Partial<Project>): Promise<Project> => {
    return apiClient<Project>(`/api/projects/${id}`, {
      method: 'PUT',
      body: JSON.stringify(project),
    })
  },
  updateStage: async (
    id: string,
    stage: ProjectStage,
    status?: ProjectStatus
  ): Promise<{ status: string }> => {
    return apiClient<{ status: string }>(`/api/projects/${id}/stage`, {
      method: 'PATCH',
      body: JSON.stringify({ stage, status }),
    })
  },
  delete: async (id: string): Promise<{ status: string }> => {
    return apiClient<{ status: string }>(`/api/projects/${id}`, {
      method: 'DELETE',
    })
  },
}

// Sessions API
export const sessionsApi = {
  list: async (): Promise<ProjectSession[]> => {
    return apiClient<ProjectSession[]>('/api/sessions')
  },
  create: async (session: Partial<ProjectSession>): Promise<ProjectSession> => {
    return apiClient<ProjectSession>('/api/sessions', {
      method: 'POST',
      body: JSON.stringify(session),
    })
  },
  update: async (id: string, session: Partial<ProjectSession>): Promise<ProjectSession> => {
    return apiClient<ProjectSession>(`/api/sessions/${id}`, {
      method: 'PUT',
      body: JSON.stringify(session),
    })
  },
  updateStatus: async (id: string, status: SessionStatus): Promise<{ status: string }> => {
    return apiClient<{ status: string }>(`/api/sessions/${id}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ status }),
    })
  },
  delete: async (id: string): Promise<{ status: string }> => {
    return apiClient<{ status: string }>(`/api/sessions/${id}`, {
      method: 'DELETE',
    })
  },
}

// Orders API
export const ordersApi = {
  list: async (): Promise<Order[]> => {
    return apiClient<Order[]>('/api/orders')
  },
  create: async (order: Partial<Order>): Promise<Order> => {
    return apiClient<Order>('/api/orders', {
      method: 'POST',
      body: JSON.stringify(order),
    })
  },
  update: async (id: string, order: Partial<Order>): Promise<Order> => {
    return apiClient<Order>(`/api/orders/${id}`, {
      method: 'PUT',
      body: JSON.stringify(order),
    })
  },
  updateStage: async (id: string, stage: OrderStage): Promise<{ status: string }> => {
    return apiClient<{ status: string }>(`/api/orders/${id}/stage`, {
      method: 'PATCH',
      body: JSON.stringify({ stage }),
    })
  },
  delete: async (id: string): Promise<{ status: string }> => {
    return apiClient<{ status: string }>(`/api/orders/${id}`, {
      method: 'DELETE',
    })
  },
}

// Billing API
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

