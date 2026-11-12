export type SettingsTab =
  // Settings
  | 'master-data'
  | 'photographers'
  | 'sequence-id'
  // Preset Configuration
  | 'preset-services'
  | 'preset-packages'
  | 'preset-orders'
  | 'preset-deliverables'
  // Customize Design & Templates
  | 'template-forms'
  | 'template-quotation-pdf'
  | 'template-invoice-pdf'
  | 'template-contracts'
  | 'template-feedbacks'
  // Stages
  | 'stages-workflow'

export interface MasterDataItem {
  id: string
  name: string
  description?: string
  color?: string
}

export interface MasterDataCategory {
  key: string
  title: string
  subtitle: string
  items: MasterDataItem[]
}

export interface TeamDesignation {
  id: string
  name: string
}

export interface CrewMember {
  id: string
  fullName: string
  countryCode: string
  mobileNo: string
  email: string
  designation: string
  note?: string
  avatarUrl?: string
}

export interface SequenceConfig {
  prefix: string
  nextNumber: string
}

export interface SequenceIdSettings {
  leads: SequenceConfig
  projects: SequenceConfig
  invoice: SequenceConfig
  orders: SequenceConfig
  expense: SequenceConfig
  cancellation: SequenceConfig
  client: SequenceConfig
  // Visibility tab
  showClientTaxNo: boolean
  studioTaxNo: string
}

