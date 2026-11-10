import { apiClient } from '@repo/ui/lib/api-client'
import type {
  MasterDataCategory,
  TeamDesignation,
  CrewMember,
  SequenceIdSettings,
} from '../types/settings'
import type { OrderPresetItem } from '../pages/settings/OrdersSettingsPage'
import type { DeliverablePresetItem } from '../pages/settings/DeliverablesSettingsPage'
import type { FeedbackQuestion } from '../pages/settings/FeedbackTemplatePage'
import type { ServiceItem } from '../pages/settings/ServicesSettingsPage'
import type { StudioPackage } from '../pages/settings/PackagesSettingsPage'

export interface WorkflowStageItem {
  id: string
  name: string
  color: string
  description: string
}

export interface StudioSettingsStore {
  masterData: Record<string, MasterDataCategory>
  designations?: TeamDesignation[]
  crewMembers: CrewMember[]
  sequenceSettings: SequenceIdSettings
  ordersPresets: OrderPresetItem[]
  deliverablesPresets: DeliverablePresetItem[]
  deliverableGroups: string[]
  services: ServiceItem[]
  serviceGroups: string[]
  packages: StudioPackage[]
  workflowStages: WorkflowStageItem[]
  contractTemplate: {
    contractName: string
    mySignatureRequired: boolean
    variableValues: Record<string, string>
  }
  feedbackTemplate: {
    templateName: string
    questions: FeedbackQuestion[]
  }
}

export const settingsApi = {
  getSettings: async (): Promise<StudioSettingsStore> => {
    try {
      const data = await apiClient<Record<string, any>>('/api/settings')
      return {
        masterData: data.master_data || {},
        designations: data.designations || [],
        crewMembers: data.crew_members || [],
        sequenceSettings: data.sequence_settings || {},
        ordersPresets: data.order_presets || [],
        deliverablesPresets: data.deliverable_presets || [],
        deliverableGroups: data.deliverable_groups || [],
        services: data.services || [],
        serviceGroups: data.service_groups || [],
        packages: data.packages || [],
        workflowStages: data.workflow_stages || [],
        contractTemplate: data.contract_template || {
          contractName: '',
          mySignatureRequired: true,
          variableValues: {},
        },
        feedbackTemplate: data.feedback_template || {
          templateName: '',
          questions: [],
        },
      }
    } catch {
      return {
        masterData: {},
        designations: [],
        crewMembers: [],
        sequenceSettings: {} as any,
        ordersPresets: [],
        deliverablesPresets: [],
        deliverableGroups: [],
        services: [],
        serviceGroups: [],
        packages: [],
        workflowStages: [],
        contractTemplate: { contractName: '', mySignatureRequired: true, variableValues: {} },
        feedbackTemplate: { templateName: '', questions: [] },
      }
    }
  },

  updateSettings: async (
    updates: Partial<StudioSettingsStore>
  ): Promise<StudioSettingsStore> => {
    const payload: Record<string, any> = {}
    if (updates.masterData !== undefined) payload.master_data = updates.masterData
    if (updates.designations !== undefined) payload.designations = updates.designations
    if (updates.crewMembers !== undefined) payload.crew_members = updates.crewMembers
    if (updates.sequenceSettings !== undefined) payload.sequence_settings = updates.sequenceSettings
    if (updates.ordersPresets !== undefined) payload.order_presets = updates.ordersPresets
    if (updates.deliverablesPresets !== undefined) payload.deliverable_presets = updates.deliverablesPresets
    if (updates.deliverableGroups !== undefined) payload.deliverable_groups = updates.deliverableGroups
    if (updates.services !== undefined) payload.services = updates.services
    if (updates.serviceGroups !== undefined) payload.service_groups = updates.serviceGroups
    if (updates.packages !== undefined) payload.packages = updates.packages
    if (updates.workflowStages !== undefined) payload.workflow_stages = updates.workflowStages
    if (updates.contractTemplate !== undefined) payload.contract_template = updates.contractTemplate
    if (updates.feedbackTemplate !== undefined) payload.feedback_template = updates.feedbackTemplate

    await apiClient('/api/settings', {
      method: 'POST',
      body: JSON.stringify(payload),
    })

    return settingsApi.getSettings()
  },
}
