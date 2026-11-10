import { apiClient } from '@repo/ui/lib/api-client'
import type { ProjectSession, SessionStatus } from '../types/session'

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
