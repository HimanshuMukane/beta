import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiClient } from '@repo/ui/lib/api-client'
import type { OrderPresetItem } from '../pages/settings/OrdersSettingsPage'

export function useOrderPresets() {
  const queryClient = useQueryClient()

  const query = useQuery({
    queryKey: ['order-presets'],
    queryFn: async (): Promise<OrderPresetItem[]> => {
      try {
        const data = await apiClient<OrderPresetItem[]>('/api/order-presets')
        return Array.isArray(data) ? data : []
      } catch {
        return []
      }
    },
  })

  const saveBatchMutation = useMutation({
    mutationFn: async (items: OrderPresetItem[]) => {
      return apiClient('/api/order-presets/batch', {
        method: 'PUT',
        body: JSON.stringify(items),
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['order-presets'] })
    },
  })

  return {
    presets: query.data || [],
    isLoading: query.isLoading,
    refetch: query.refetch,
    savePresets: saveBatchMutation.mutateAsync,
    isSaving: saveBatchMutation.isPending,
  }
}

