import { apiClient } from '@repo/ui/lib/api-client'

export interface ChecklistItem {
  id: string
  text: string
  done: boolean
}

export interface TaskItem {
  id: string
  taskId: string
  tenantId: string
  projectId: string
  taskName: string
  categoryName: string
  phase: string
  status: string
  priority: string
  dueDate: string
  assignee: string
  notes: string
  checklist: ChecklistItem[]
  orderIndex: number
}

export interface CreateTaskPayload {
  projectId: string
  taskName: string
  categoryName: string
  phase: string
  status?: string
  priority?: string
  dueDate?: string
  assignee?: string
  notes?: string
  checklist?: ChecklistItem[]
}

export interface ProjectDeliverableItem {
  id: string
  tenantId: string
  projectId: string
  deliverable: string
  dueDate: string
  assignee: string
  status: 'New' | 'In Progress' | 'Client Review' | 'Completed' | 'Delivered'
  note?: string
}

export interface CreateProjectDeliverablePayload {
  projectId: string
  deliverable: string
  dueDate?: string
  assignee?: string
  status?: string
  note?: string
}

export const tasksApi = {
  list: async (params?: { projectId?: string; phase?: string }): Promise<TaskItem[]> => {
    const query = new URLSearchParams()
    if (params?.projectId) query.append('projectId', params.projectId)
    if (params?.phase) query.append('phase', params.phase)
    const queryString = query.toString()
    const endpoint = queryString ? `/api/tasks?${queryString}` : '/api/tasks'
    return apiClient<TaskItem[]>(endpoint)
  },

  get: async (taskId: string): Promise<TaskItem> => {
    return apiClient<TaskItem>(`/api/tasks/${taskId}`)
  },

  create: async (payload: CreateTaskPayload): Promise<TaskItem> => {
    return apiClient<TaskItem>(`/api/projects/${payload.projectId}/tasks`, {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  update: async (taskId: string, payload: Partial<CreateTaskPayload>): Promise<TaskItem> => {
    return apiClient<TaskItem>(`/api/tasks/${taskId}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  },

  updateStatus: async (taskId: string, status: string): Promise<{ status: string; id: string; taskStatus: string }> => {
    return apiClient<{ status: string; id: string; taskStatus: string }>(`/api/tasks/${taskId}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ status }),
    })
  },

  delete: async (taskId: string): Promise<{ status: string; id: string }> => {
    return apiClient<{ status: string; id: string }>(`/api/tasks/${taskId}`, {
      method: 'DELETE',
    })
  },
}

export const projectDeliverablesApi = {
  list: async (params?: { projectId?: string; status?: string; assignee?: string }): Promise<ProjectDeliverableItem[]> => {
    const query = new URLSearchParams()
    if (params?.projectId) query.append('projectId', params.projectId)
    if (params?.status) query.append('status', params.status)
    if (params?.assignee) query.append('assignee', params.assignee)
    const queryString = query.toString()
    const endpoint = queryString ? `/api/project-deliverables?${queryString}` : '/api/project-deliverables'
    return apiClient<ProjectDeliverableItem[]>(endpoint)
  },

  create: async (payload: CreateProjectDeliverablePayload): Promise<ProjectDeliverableItem> => {
    return apiClient<ProjectDeliverableItem>('/api/project-deliverables', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  update: async (id: string, payload: Partial<CreateProjectDeliverablePayload>): Promise<ProjectDeliverableItem> => {
    return apiClient<ProjectDeliverableItem>(`/api/project-deliverables/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  },

  updateStatus: async (id: string, status: string): Promise<{ status: string; id: string; deliverableStatus: string }> => {
    return apiClient<{ status: string; id: string; deliverableStatus: string }>(`/api/project-deliverables/${id}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ status }),
    })
  },

  delete: async (id: string): Promise<{ status: string; id: string }> => {
    return apiClient<{ status: string; id: string }>(`/api/project-deliverables/${id}`, {
      method: 'DELETE',
    })
  },
}
