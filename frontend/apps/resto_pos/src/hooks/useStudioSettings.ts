import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { settingsApi, type StudioSettingsStore } from '../api/settingsApi'

export const SETTINGS_QUERY_KEY = ['studio-settings']

export function useStudioSettings() {
  const queryClient = useQueryClient()

  const query = useQuery({
    queryKey: SETTINGS_QUERY_KEY,
    queryFn: () => settingsApi.getSettings(),
  })

  const mutation = useMutation({
    mutationFn: (updates: Partial<StudioSettingsStore>) => settingsApi.updateSettings(updates),
    onSuccess: (newData) => {
      queryClient.setQueryData(SETTINGS_QUERY_KEY, newData)
    },
  })

  return {
    settings: query.data,
    isLoading: query.isLoading,
    isError: query.isError,
    error: query.error,
    updateSettings: mutation.mutate,
    updateSettingsAsync: mutation.mutateAsync,
    isUpdating: mutation.isPending,
  }
}

