import { apiClient } from '@repo/ui/lib/api-client'
import type { Project, ProjectStage, ProjectStatus } from '../types/project'

export interface ProjectListParams {
  page?: number
  pageSize?: number
  search?: string
  status?: string
}

export const projectsApi = {
  list: async (_params?: ProjectListParams): Promise<Project[]> => {
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
