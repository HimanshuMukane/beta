import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiClient } from '@repo/ui/lib/api-client'

export interface WorkflowItem {
  id: string
  tenant_id?: string
  workflow_type: string
  pipeline: string
  stage: string
  color_code: string
  order_index: number
}

export const workflowsApi = {
  list: async (type = 'lead'): Promise<WorkflowItem[]> => {
    return apiClient<WorkflowItem[]>(`/api/workflows?type=${type}`)
  },
  create: async (item: Partial<WorkflowItem>): Promise<WorkflowItem> => {
    return apiClient<WorkflowItem>('/api/workflows', {
      method: 'POST',
      body: JSON.stringify(item),
    })
  },
  update: async (id: string, item: Partial<WorkflowItem>): Promise<{ status: string }> => {
    return apiClient<{ status: string }>(`/api/workflows/${id}`, {
      method: 'PUT',
      body: JSON.stringify(item),
    })
  },
  saveBatch: async (workflow_type: string, stages: WorkflowItem[]): Promise<{ status: string }> => {
    return apiClient<{ status: string }>('/api/workflows/batch', {
      method: 'PUT',
      body: JSON.stringify({ workflow_type, stages }),
    })
  },
  delete: async (id: string): Promise<{ status: string }> => {
    return apiClient<{ status: string }>(`/api/workflows/${id}`, {
      method: 'DELETE',
    })
  },
}

export const WORKFLOWS_QUERY_KEY = ['workflows']

export function useWorkflows(workflowType = 'lead') {
  const queryClient = useQueryClient()

  const query = useQuery({
    queryKey: [...WORKFLOWS_QUERY_KEY, workflowType],
    queryFn: () => workflowsApi.list(workflowType),
  })

  const saveBatchMutation = useMutation({
    mutationFn: (stages: WorkflowItem[]) => workflowsApi.saveBatch(workflowType, stages),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [...WORKFLOWS_QUERY_KEY, workflowType] })
    },
  })

  return {
    workflows: query.data || [],
    isLoading: query.isLoading,
    isError: query.isError,
    saveBatch: saveBatchMutation.mutate,
    saveBatchAsync: saveBatchMutation.mutateAsync,
    isSaving: saveBatchMutation.isPending,
    refetch: query.refetch,
  }
}

