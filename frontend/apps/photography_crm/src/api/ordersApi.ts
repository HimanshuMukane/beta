import { apiClient } from '@repo/ui/lib/api-client'
import type { Order, OrderStage } from '../types/order'

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
