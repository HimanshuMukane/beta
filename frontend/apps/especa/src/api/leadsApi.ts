import { apiClient } from '@repo/ui/lib/api-client'
import type { Lead, PipelineStage, LeadStage } from '../types/lead'

export interface LeadListParams {
  page?: number
  pageSize?: number
  search?: string
  status?: string
}

export const leadsApi = {
  list: async (params?: LeadListParams, signal?: AbortSignal): Promise<Lead[]> => {
    return apiClient<Lead[]>('/api/leads', { params: params as any, signal })
  },
  get: async (id: string, signal?: AbortSignal): Promise<Lead> => {
    return apiClient<Lead>(`/api/leads/${id}`, { signal })
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
