export type ProjectStage =
  | 'Booked'
  | 'Shoot Planning'
  | 'Shoot Completed'
  | 'Editing'
  | 'Delivery'
  | 'Closure'

export type ProjectStatus =
  | 'New'
  | 'Booked'
  | 'Active'
  | 'Completed'
  | 'On Hold'
  | 'Cancelled'

export type ProjectType =
  | 'Wedding'
  | 'Pre-Wedding'
  | 'Engagement'
  | 'Birthday'
  | 'Maternity'
  | 'Newborn'
  | 'Portrait'
  | 'Corporate'
  | 'Product'
  | 'Fashion'
  | 'Event'
  | 'Other'

export interface ProjectDocument {
  id: string
  type: 'Form' | 'Lead' | 'Invoice' | 'Contract' | 'Quote'
  title: string
  referenceNo?: string
  amount?: number
  status: 'Submitted' | 'Converted' | 'Pending' | 'Paid' | 'Signed'
  date: string
}

import type { ShootTiming, SessionStatus } from './session'

export interface ProjectSession {
  id: string
  title: string
  projectName?: string
  sessionManager?: string
  projectType?: ProjectType
  shootTiming?: ShootTiming
  status: SessionStatus
  date: string
  startTime?: string
  duration?: string
  timeRange: string
  venue: string
  city: string
  locationLink?: string
  remarks?: string
  members: string[]
  shootsActivityCount?: number
}

export interface ProjectTask {
  id: string
  title: string
  assignee: string
  dueDate: string
  completed: boolean
  priority: 'low' | 'medium' | 'high'
}

export interface ProjectPrintingItem {
  id: string
  labName: string
  galleriesName: string
  orderName: string
  orderType: string
  associate: string
  status: string
  cost?: number
}

export interface CustomerDeliverable {
  id: string
  deliverable: string
  dueDate: string
  assignee: string
  status: 'New' | 'In Progress' | 'Client Review' | 'Completed' | 'Delivered'
  note?: string
}

export interface ProjectExpenseItem {
  id: string
  category: string
  description: string
  amount: number
  date: string
  paidBy: string
  sessionName?: string
}

export interface ProjectSegment {
  id: string
  name: string
  type: ProjectType
  startDate: string
  endDate: string
}

export interface ProjectPaymentRecord {
  id: string
  createdDate: string
  paymentDate: string
  amount: number
  type: string
  referenceNo: string
  status: string
  note?: string
}

export interface ProjectFinance {
  profit: number
  currentProfit: number
  netPayableAmount: number
  paidAmount: number
  balanceAmount: number
  paymentDue: number
  expense: number
  eventDiscountAmount: number
  paymentHistory: ProjectPaymentRecord[]
}

export interface Project {
  id: string
  projectName: string
  clientName: string
  clientInitials: string
  clientMobile: string
  clientCountryCode: string
  clientEmail: string
  clientAddress?: string
  leadSource?: string
  projectManager: {
    name: string
    initials: string
    avatar?: string
  }
  projectType: ProjectType
  projectAddress?: string
  sessionsCount: number
  tasksCount: { completed: number; total: number }
  deliverablesCount: { ready: number; total: number }
  status: ProjectStatus
  stage: ProjectStage
  startDate: string
  endDate: string
  createdDate: string
  // Deep detail tabs
  segments: ProjectSegment[]
  documents: ProjectDocument[]
  sessions: ProjectSession[]
  tasks: ProjectTask[]
  printing: ProjectPrintingItem[]
  deliverables: CustomerDeliverable[]
  expenses: ProjectExpenseItem[]
  notes: string
  // Insights
  shootHours: number
  shootMembersCount: number
  shootExpense: number
  galleriesCount: number
  finance: ProjectFinance
}
