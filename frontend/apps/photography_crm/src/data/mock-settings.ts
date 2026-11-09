import type {
  MasterDataCategory,
  TeamDesignation,
  CrewMember,
  SequenceIdSettings,
} from '../types/settings'

export const INITIAL_MASTER_DATA: Record<string, MasterDataCategory> = {
  projects: {
    key: 'projects',
    title: 'Projects',
    subtitle: 'Types of shoots and projects handled (Wedding, Pre-Wedding, Corporate, etc.)',
    items: [
      { id: 'p-1', name: 'Wedding' },
      { id: 'p-2', name: 'Pre-Wedding' },
      { id: 'p-3', name: 'Engagement' },
      { id: 'p-4', name: 'Corporate' },
      { id: 'p-5', name: 'Maternity' },
      { id: 'p-6', name: 'Fashion & Editorial' },
      { id: 'p-7', name: 'Birthday' },
    ],
  },
  expense: {
    key: 'expense',
    title: 'Expense',
    subtitle: 'Expense categories for equipment, travel, lab processing, rentals',
    items: [
      { id: 'e-1', name: 'Equipment Rental' },
      { id: 'e-2', name: 'Travel & Conveyance' },
      { id: 'e-3', name: 'Lab Printing Cost' },
      { id: 'e-4', name: 'Styling & Props' },
      { id: 'e-5', name: 'Food & Refreshments' },
      { id: 'e-6', name: 'Courier & Packaging' },
    ],
  },
  tax: {
    key: 'tax',
    title: 'Tax',
    subtitle: 'Standard studio GST and tax rate slabs',
    items: [
      { id: 't-1', name: 'GST 18% (Standard Services)' },
      { id: 't-2', name: 'GST 12% (Goods & Albums)' },
      { id: 't-3', name: 'GST 5% (Job Work)' },
      { id: 't-4', name: 'Tax Exempt / 0%' },
    ],
  },
  orders: {
    key: 'orders',
    title: 'Orders',
    subtitle: 'Types of studio merchandise, prints, and walk-in sales',
    items: [
      { id: 'o-1', name: 'Walk-in' },
      { id: 'o-2', name: 'Studio Session' },
      { id: 'o-3', name: 'Print & Album' },
      { id: 'o-4', name: 'Framing' },
      { id: 'o-5', name: 'Commercial Shoot' },
      { id: 'o-6', name: 'Passport/Visa Photo' },
      { id: 'o-7', name: 'Custom Merch' },
    ],
  },
  'session-timings': {
    key: 'session-timings',
    title: 'Session Timings',
    subtitle: 'Pre-configured shoot shifts and lighting windows',
    items: [
      { id: 's-1', name: 'Morning Golden Hour (06:00 AM - 10:00 AM)' },
      { id: 's-2', name: 'Midday Studio (11:00 AM - 03:00 PM)' },
      { id: 's-3', name: 'Evening Sunset (04:00 PM - 08:00 PM)' },
      { id: 's-4', name: 'Night / Reception (07:00 PM - Midnight)' },
      { id: 's-5', name: 'Full Day Event (10 Hours)' },
    ],
  },
  'projects-status': {
    key: 'projects-status',
    title: 'Projects Status',
    subtitle: 'Workflow stages and execution statuses for shoots',
    items: [
      { id: 'ps-1', name: 'Booked' },
      { id: 'ps-2', name: 'Shoot Planning' },
      { id: 'ps-3', name: 'Shoot Completed' },
      { id: 'ps-4', name: 'Editing / Post-Production' },
      { id: 'ps-5', name: 'Delivery Pending' },
      { id: 'ps-6', name: 'Closure' },
    ],
  },
}

export const INITIAL_DESIGNATIONS: TeamDesignation[] = [
  { id: 'des-1', name: 'Traditional Photographer' },
  { id: 'des-2', name: 'Traditional Videographer' },
  { id: 'des-3', name: 'Candid Photographer' },
  { id: 'des-4', name: 'Candid Videographer' },
  { id: 'des-5', name: 'Drone Operator' },
  { id: 'des-6', name: 'Cinematic Videographer' },
  { id: 'des-7', name: 'Lead Colorist & Editor' },
]

export const INITIAL_CREW_MEMBERS: CrewMember[] = [
  {
    id: 'crew-1',
    fullName: 'Rahul Sharma',
    countryCode: '+91',
    mobileNo: '9876543210',
    email: 'rahul.s@especa.com',
    designation: 'Candid Photographer',
    note: 'Specializes in low-light wedding receptions and portraiture.',
  },
  {
    id: 'crew-2',
    fullName: 'Pooja Hegde',
    countryCode: '+91',
    mobileNo: '9820011223',
    email: 'pooja.h@especa.com',
    designation: 'Cinematic Videographer',
    note: 'Gimbal specialist and lead wedding film director.',
  },
  {
    id: 'crew-3',
    fullName: 'Vikram Joshi',
    countryCode: '+91',
    mobileNo: '9711223344',
    email: 'vikram.j@especa.com',
    designation: 'Drone Operator',
    note: 'DGCA certified aerial cinema pilot.',
  },
  {
    id: 'crew-4',
    fullName: 'Amit Verma',
    countryCode: '+91',
    mobileNo: '9833445566',
    email: 'amit.v@especa.com',
    designation: 'Traditional Photographer',
    note: 'Stage group photography and family formal portraiture.',
  },
]

export const INITIAL_SEQUENCE_SETTINGS: SequenceIdSettings = {
  leads: { prefix: 'ESPE/LQuote/', nextNumber: '101' },
  projects: { prefix: 'ESPE/EVENT/', nextNumber: '101' },
  invoice: { prefix: 'ESPE/INVOICE/', nextNumber: '101' },
  orders: { prefix: 'ESPE/WALKIN/', nextNumber: '100' },
  expense: { prefix: 'ESPE/EX/', nextNumber: '100' },
  cancellation: { prefix: 'ESPE/CANCEL/', nextNumber: '100' },
  client: { prefix: 'ESPE/CLIENT/', nextNumber: '101' },
  // Visibility tab
  showClientTaxNo: true,
  studioTaxNo: '27AABCE1234F1Z5',
}

