import { apiClient } from '@repo/ui/lib/api-client'

export interface CrewMember {
  id: string
  tenant_id: string
  fullName: string
  firstName: string
  lastName: string
  email: string
  mobileNo: string
  countryCode: string
  role_id: string
  designation: string
}

export const crewApi = {
  list: async (): Promise<CrewMember[]> => {
    return apiClient<CrewMember[]>('/api/crew')
  },
  create: async (member: Partial<CrewMember>): Promise<CrewMember> => {
    return apiClient<CrewMember>('/api/crew', {
      method: 'POST',
      body: JSON.stringify(member),
    })
  },
  delete: async (id: string): Promise<{ status: string }> => {
    return apiClient<{ status: string }>(`/api/crew/${id}`, {
      method: 'DELETE',
    })
  },
}
