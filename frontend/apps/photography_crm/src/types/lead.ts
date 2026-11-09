export type LeadSource =
  | 'All Leads'
  | 'Messenger'
  | 'Instagram'
  | 'Instagram Post'
  | 'Instagram Ads'
  | 'Whatsapp'
  | 'Leads Ads'
  | 'Google Search'
  | 'Google Ads'
  | 'Friend Referral'
  | 'Family Referral'
  | 'Studio Shop'
  | 'Form'
  | 'Manual Leads'

// Primary Pipeline Stages
export type PipelineStage = 'new_leads' | 'contacted' | 'qualified' | 'converted'

// Sub-stages per Pipeline Stage
export type ContactedSubStage = 'Contacted' | 'Scheduled' | 'Untouched'
export type QualifiedSubStage = 'Negotiation' | 'Ready to Book' | 'Interested'
export type ConvertedSubStage = 'Booked' | 'Lost Lead'

export type LeadStage =
  | 'New Lead'
  | ContactedSubStage
  | QualifiedSubStage
  | ConvertedSubStage

export type LeadPipeline = 'Sales Pipeline' | 'Wedding Inquiries' | 'Commercial Shoots' | 'Portraits'
export type ProjectType = 'Wedding' | 'Pre-Wedding' | 'Engagement' | 'Corporate' | 'Maternity' | 'Fashion' | 'Birthday'
export type QuoteType = 'Standard Package' | 'Custom Estimate' | 'Retainer' | 'Hourly Shoot'

// Navigation Perspectives
export type PrimaryPerspective =
  | 'all'
  | 'new_leads'
  | 'contacted'
  | 'qualified'
  | 'converted'
  | 'followups'

export type FollowupSubFilter =
  | 'all_followup'
  | 'my_followup'
  | 'due_followup'
  | 'today_followup'
  | 'upcoming_followup'

export interface ProjectSegment {
  id: string
  name: string
  type: ProjectType
  startDate: string
  endDate: string
}

export interface ServiceItem {
  id: string
  name: string
  groupName: string
  qty: number
  cost: number
  subTotal: number
  sessionPlanned?: boolean
}

export interface LeadDeliverable {
  id: string
  name: string
  category: 'Wedding' | 'Complimentary' | 'Deliverables'
  selected: boolean
}

export interface Lead {
  id: string
  name: string
  mobileNo: string
  countryCode: string
  email: string
  address?: string
  createdDate: string
  leadSource: LeadSource
  generate: string
  assignee: {
    name: string
    avatar?: string
    initials: string
  }
  leadPipeline: LeadPipeline
  pipelineStage: PipelineStage
  leadStage: LeadStage
  quoteAmount: number
  followupDate: string // Format: YYYY-MM-DD or YYYY-MM-DDTHH:mm
  projectType: ProjectType
  quoteType: QuoteType
  projectStartDate: string
  projectEndDate: string
  projectName?: string
  projectAddress?: string
  projectSegments: ProjectSegment[]
  packageName?: string
  selectedServices: ServiceItem[]
  selectedDeliverables: string[]
  quotationTemplate?: string
  discountAmount?: number
  taxAmount?: number
}
