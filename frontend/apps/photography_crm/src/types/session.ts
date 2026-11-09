import type { ProjectType } from './project'
export type { ProjectType }

export type ShootTiming =
  | 'Morning'
  | 'Afternoon'
  | 'Evening'
  | 'Night'
  | 'Full Day'
  | 'Multiple Days'

export type SessionStatus = 'New' | 'Scheduled' | 'Completed' | 'Delayed'

export interface CrewMember {
  id: string
  name: string
  role: 'Lead Photographer' | 'Candid Photographer' | 'Cinematographer' | 'Drone Operator' | 'Assistant' | 'Editor'
  avatar?: string
}

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

export { PROJECT_TYPES } from '../lib/formatters'

export const SHOOT_TIMINGS: ShootTiming[] = [
  'Morning',
  'Afternoon',
  'Evening',
  'Night',
  'Full Day',
  'Multiple Days',
]

// 1 hour to 11 hours in 30 min intervals
export const DURATION_OPTIONS: string[] = [
  '1 hour',
  '1 hour 30 min',
  '2 hours',
  '2 hours 30 min',
  '3 hours',
  '3 hours 30 min',
  '4 hours',
  '4 hours 30 min',
  '5 hours',
  '5 hours 30 min',
  '6 hours',
  '6 hours 30 min',
  '7 hours',
  '7 hours 30 min',
  '8 hours',
  '8 hours 30 min',
  '9 hours',
  '9 hours 30 min',
  '10 hours',
  '10 hours 30 min',
  '11 hours',
]

export const DEFAULT_AVAILABLE_CREW: CrewMember[] = [
  { id: 'cm-1', name: 'Rohan Mehta', role: 'Lead Photographer' },
  { id: 'cm-2', name: 'Aditya Roy', role: 'Cinematographer' },
  { id: 'cm-3', name: 'Priya Sharma', role: 'Candid Photographer' },
  { id: 'cm-4', name: 'Vikram Seth', role: 'Drone Operator' },
  { id: 'cm-5', name: 'Neha Gupta', role: 'Assistant' },
  { id: 'cm-6', name: 'Kavita Roy', role: 'Editor' },
]
