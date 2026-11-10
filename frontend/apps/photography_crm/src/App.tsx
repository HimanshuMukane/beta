import { useState, useMemo, Suspense, lazy } from 'react'
import {
  BrowserRouter,
  Routes,
  Route,
  Navigate,
  useNavigate,
  useLocation,
  useParams,
} from 'react-router-dom'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { DashboardShell } from '@repo/ui/components/dashboard-shell'
import { Toaster, toast } from '@repo/ui/components/ui/sonner'
import { AuthProvider, useAuth } from '@repo/ui/hooks/use-auth'
import { ProtectedRoute } from '@repo/ui/components/protected-route'

const SigninPage = lazy(() => import('./pages/auth/SigninPage'))
const SignupPage = lazy(() => import('./pages/auth/SignupPage'))
import type { Lead, PipelineStage } from './types/lead'
import type { Project, ProjectStage } from './types/project'
import type { ProjectSession, SessionStatus } from './types/session'
import type { Order, OrderStage } from './types/order'
import type { SettingsTab } from './types/settings'
import { PhotographySidebar, type ReportCategory } from './components/PhotographySidebar'
import { UserSidebar, type UserRouteKey } from './components/UserSidebar'
import { leadsApi, projectsApi, sessionsApi, ordersApi } from './api/crmApi'

// Lazy Load User Portal Pages
const UserHomePage = lazy(() =>
  import('./pages/user/UserHomePage').then((m) => ({ default: m.UserHomePage }))
)
const UserShootsPage = lazy(() =>
  import('./pages/user/UserShootsPage').then((m) => ({ default: m.UserShootsPage }))
)
const UserTasksPage = lazy(() =>
  import('./pages/user/UserTasksPage').then((m) => ({ default: m.UserTasksPage }))
)
const UserDeliverablesPage = lazy(() =>
  import('./pages/user/UserDeliverablesPage').then((m) => ({ default: m.UserDeliverablesPage }))
)

// Lazy Load Core Admin Pages
const LeadsPage = lazy(() =>
  import('./pages/admin/LeadsPage').then((m) => ({ default: m.LeadsPage }))
)
const LeadFormPage = lazy(() =>
  import('./pages/admin/LeadFormPage').then((m) => ({ default: m.LeadFormPage }))
)
const ProjectsPage = lazy(() =>
  import('./pages/admin/ProjectsPage').then((m) => ({ default: m.ProjectsPage }))
)
const ProjectCreatePage = lazy(() =>
  import('./pages/admin/ProjectCreatePage').then((m) => ({ default: m.ProjectCreatePage }))
)
const ProjectEditPage = lazy(() =>
  import('./pages/admin/ProjectEditPage').then((m) => ({ default: m.ProjectEditPage }))
)
const SessionsPage = lazy(() =>
  import('./pages/admin/SessionsPage').then((m) => ({ default: m.SessionsPage }))
)
const BillingsPage = lazy(() =>
  import('./pages/admin/BillingsPage').then((m) => ({ default: m.BillingsPage }))
)
const OrdersPage = lazy(() =>
  import('./pages/admin/OrdersPage').then((m) => ({ default: m.OrdersPage }))
)
const OrderFormPage = lazy(() =>
  import('./pages/admin/OrderFormPage').then((m) => ({ default: m.OrderFormPage }))
)

// Lazy Load Reports Pages
const BillingReportPage = lazy(() =>
  import('./pages/reports/BillingReportPage').then((m) => ({ default: m.BillingReportPage }))
)
const LeadsReportPage = lazy(() =>
  import('./pages/reports/LeadsReportPage').then((m) => ({ default: m.LeadsReportPage }))
)
const OrdersReportPage = lazy(() =>
  import('./pages/reports/OrdersReportPage').then((m) => ({ default: m.OrdersReportPage }))
)
const SessionsReportPage = lazy(() =>
  import('./pages/reports/SessionsReportPage').then((m) => ({ default: m.SessionsReportPage }))
)
const FinanceReportPage = lazy(() =>
  import('./pages/reports/FinanceReportPage').then((m) => ({ default: m.FinanceReportPage }))
)
const MasterReportPage = lazy(() =>
  import('./pages/reports/MasterReportPage').then((m) => ({ default: m.MasterReportPage }))
)
const AccountAnalyticsPage = lazy(() =>
  import('./pages/reports/AccountAnalyticsPage').then((m) => ({ default: m.AccountAnalyticsPage }))
)
const FeedbackReportPage = lazy(() =>
  import('./pages/reports/FeedbackReportPage').then((m) => ({ default: m.FeedbackReportPage }))
)

// Lazy Load Settings Pages
const MasterDataPage = lazy(() =>
  import('./pages/settings/MasterDataPage').then((m) => ({ default: m.MasterDataPage }))
)
const PhotographersPage = lazy(() =>
  import('./pages/settings/PhotographersPage').then((m) => ({ default: m.PhotographersPage }))
)
const SequenceIdPage = lazy(() =>
  import('./pages/settings/SequenceIdPage').then((m) => ({ default: m.SequenceIdPage }))
)
const ServicesSettingsPage = lazy(() =>
  import('./pages/settings/ServicesSettingsPage').then((m) => ({ default: m.ServicesSettingsPage }))
)
const PackagesSettingsPage = lazy(() =>
  import('./pages/settings/PackagesSettingsPage').then((m) => ({ default: m.PackagesSettingsPage }))
)
const OrdersSettingsPage = lazy(() =>
  import('./pages/settings/OrdersSettingsPage').then((m) => ({ default: m.OrdersSettingsPage }))
)
const DeliverablesSettingsPage = lazy(() =>
  import('./pages/settings/DeliverablesSettingsPage').then((m) => ({ default: m.DeliverablesSettingsPage }))
)
const ContractTemplatePage = lazy(() =>
  import('./pages/settings/ContractTemplatePage').then((m) => ({ default: m.ContractTemplatePage }))
)
const FeedbackTemplatePage = lazy(() =>
  import('./pages/settings/FeedbackTemplatePage').then((m) => ({ default: m.FeedbackTemplatePage }))
)
const TemplateCustomizerPage = lazy(() =>
  import('./pages/settings/TemplateCustomizerPage').then((m) => ({ default: m.TemplateCustomizerPage }))
)
const WorkflowPage = lazy(() =>
  import('./pages/settings/WorkflowPage').then((m) => ({ default: m.WorkflowPage }))
)

function ProjectsPageContainer({
  onAddNewProject,
  onEditProject,
  onDeleteProject,
  onStageChange,
}: {
  onAddNewProject: (prefill?: Partial<Project>) => void
  onEditProject: (project: Project) => void
  onDeleteProject: (projectId: string) => void
  onStageChange: (projectId: string, newStage: ProjectStage) => void
}) {
  const { data: projects = [], isLoading } = useQuery({
    queryKey: ['projects'],
    queryFn: ({ signal }) => projectsApi.list(undefined, signal),
    staleTime: 60_000,
    gcTime: 5 * 60_000,
  })

  if (isLoading && projects.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center p-8 text-xs text-muted-foreground">
        Loading projects...
      </div>
    )
  }

  return (
    <ProjectsPage
      projects={projects}
      onAddNewProject={onAddNewProject}
      onEditProject={onEditProject}
      onDeleteProject={onDeleteProject}
      onStageChange={onStageChange}
    />
  )
}

function LeadsPageContainer({
  onAddNewLead,
  onEditLead,
  onDeleteLead,
  onStageChange,
  onConvertToProject,
}: {
  onAddNewLead: (stage?: PipelineStage) => void
  onEditLead: (lead: Lead) => void
  onDeleteLead: (leadId: string) => void
  onStageChange: (leadId: string, newStage: PipelineStage) => void
  onConvertToProject: (lead: Lead) => void
}) {
  const { data: leads = [], isLoading } = useQuery({
    queryKey: ['leads'],
    queryFn: ({ signal }) => leadsApi.list(undefined, signal),
    staleTime: 60_000,
    gcTime: 5 * 60_000,
  })

  if (isLoading && leads.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center p-8 text-xs text-muted-foreground">
        Loading leads...
      </div>
    )
  }

  return (
    <LeadsPage
      leads={leads}
      onAddNewLead={onAddNewLead}
      onEditLead={onEditLead}
      onDeleteLead={onDeleteLead}
      onStageChange={onStageChange}
      onConvertToProject={onConvertToProject}
    />
  )
}

function SessionsPageContainer({
  onDeleteSession,
  onStatusChange,
  onSaveSession,
}: {
  onDeleteSession: (sessionId: string) => void
  onStatusChange: (sessionId: string, newStatus: SessionStatus) => void
  onSaveSession: (session: ProjectSession) => void
}) {
  const { data: sessions = [], isLoading } = useQuery({
    queryKey: ['sessions'],
    queryFn: ({ signal }) => sessionsApi.list(signal),
    staleTime: 60_000,
    gcTime: 5 * 60_000,
  })

  const { data: projects = [] } = useQuery({
    queryKey: ['projects'],
    queryFn: ({ signal }) => projectsApi.list(undefined, signal),
    staleTime: 60_000,
    gcTime: 5 * 60_000,
  })

  if (isLoading && sessions.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center p-8 text-xs text-muted-foreground">
        Loading sessions...
      </div>
    )
  }

  return (
    <SessionsPage
      sessions={sessions}
      onAddNewSession={() => { }}
      onEditSession={() => { }}
      onDeleteSession={onDeleteSession}
      onStatusChange={onStatusChange}
      onSaveSession={onSaveSession}
      existingProjectNames={projects.map((p) => p.projectName)}
    />
  )
}

function OrdersPageContainer({
  onAddNewOrder,
  onEditOrder,
  onDeleteOrder,
  onStageChange,
}: {
  onAddNewOrder: () => void
  onEditOrder: (order: Order) => void
  onDeleteOrder: (orderId: string) => void
  onStageChange: (orderId: string, newStage: OrderStage) => void
}) {
  const { data: orders = [], isLoading } = useQuery({
    queryKey: ['orders'],
    queryFn: ({ signal }) => ordersApi.list(signal),
    staleTime: 60_000,
    gcTime: 5 * 60_000,
  })

  if (isLoading && orders.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center p-8 text-xs text-muted-foreground">
        Loading orders...
      </div>
    )
  }

  return (
    <OrdersPage
      orders={orders}
      onAddNewOrder={onAddNewOrder}
      onEditOrder={onEditOrder}
      onDeleteOrder={onDeleteOrder}
      onStageChange={onStageChange}
    />
  )
}

function OrderFormPageWrapper({
  initialOrder,
  onBack,
  onSave,
}: {
  initialOrder: Order | null
  onBack: () => void
  onSave: (o: Order) => void
}) {
  const { data: projects = [] } = useQuery({
    queryKey: ['projects'],
    queryFn: ({ signal }) => projectsApi.list(undefined, signal),
    staleTime: 60_000,
    gcTime: 5 * 60_000,
  })

  return (
    <OrderFormPage
      initialOrder={initialOrder}
      onBack={onBack}
      onSave={onSave}
      existingProjectOptions={projects.map((p) => ({ id: p.id, name: p.projectName }))}
    />
  )
}

function UserHomePageContainer() {
  const { data: sessions = [] } = useQuery({
    queryKey: ['sessions'],
    queryFn: ({ signal }) => sessionsApi.list(signal),
    staleTime: 60_000,
    gcTime: 5 * 60_000,
  })

  const { data: projects = [] } = useQuery({
    queryKey: ['projects'],
    queryFn: ({ signal }) => projectsApi.list(undefined, signal),
    staleTime: 60_000,
    gcTime: 5 * 60_000,
  })

  return <UserHomePage sessions={sessions} projects={projects} />
}

function UserShootsPageContainer({
  onSaveSession,
  onStatusChange,
}: {
  onSaveSession: (session: ProjectSession) => void
  onStatusChange: (sessionId: string, newStatus: SessionStatus) => void
}) {
  const { data: sessions = [] } = useQuery({
    queryKey: ['sessions'],
    queryFn: ({ signal }) => sessionsApi.list(signal),
    staleTime: 60_000,
    gcTime: 5 * 60_000,
  })

  const { data: projects = [] } = useQuery({
    queryKey: ['projects'],
    queryFn: ({ signal }) => projectsApi.list(undefined, signal),
    staleTime: 60_000,
    gcTime: 5 * 60_000,
  })

  return (
    <UserShootsPage
      sessions={sessions}
      onSaveSession={onSaveSession}
      onStatusChange={onStatusChange}
      existingProjectNames={projects.map((p) => p.projectName)}
    />
  )
}

function UserTasksPageContainer({
  onUpdateProject,
}: {
  onUpdateProject: (p: Project) => void
}) {
  const { data: projects = [] } = useQuery({
    queryKey: ['projects'],
    queryFn: ({ signal }) => projectsApi.list(undefined, signal),
    staleTime: 60_000,
    gcTime: 5 * 60_000,
  })

  return <UserTasksPage projects={projects} onUpdateProject={onUpdateProject} />
}

function UserDeliverablesPageContainer({
  onUpdateProject,
}: {
  onUpdateProject: (p: Project) => void
}) {
  const { data: projects = [] } = useQuery({
    queryKey: ['projects'],
    queryFn: ({ signal }) => projectsApi.list(undefined, signal),
    staleTime: 60_000,
    gcTime: 5 * 60_000,
  })

  return <UserDeliverablesPage projects={projects} onUpdateProject={onUpdateProject} />
}

function AppContent() {
  const navigate = useNavigate()
  const location = useLocation()
  const pathname = location.pathname
  const { user } = useAuth()
  const queryClient = useQueryClient()
  const isAdmin = useMemo(() => {
    if (!user) return false
    const role = (user.role_name || '').toLowerCase().trim()
    return role === 'admin'
  }, [user])

  const [editingLead, setEditingLead] = useState<Lead | null>(null)
  const [editingProject, setEditingProject] = useState<Project | null>(null)
  const [projectPrefillData, setProjectPrefillData] = useState<Partial<Project> | null>(null)
  const [editingOrder, setEditingOrder] = useState<Order | null>(null)

  const currentRoute = useMemo(() => {
    if (pathname.startsWith('/projects')) return 'projects'
    if (pathname.startsWith('/sessions')) return 'sessions-list'
    if (pathname.startsWith('/leads')) return 'leads'
    if (pathname.startsWith('/billings')) return 'billings'
    if (pathname.startsWith('/orders')) return 'orders'
    if (pathname.startsWith('/reports')) return 'reports'
    if (pathname.startsWith('/settings')) return 'settings'
    return 'projects'
  }, [pathname])

  const currentUserRoute: UserRouteKey = useMemo(() => {
    if (pathname.startsWith('/shoots')) return 'shoots'
    if (pathname.startsWith('/tasks')) return 'tasks'
    if (pathname.startsWith('/deliverables')) return 'deliverables'
    return 'home'
  }, [pathname])

  const activeReport: ReportCategory = useMemo(() => {
    if (pathname.startsWith('/reports/')) {
      return pathname.replace('/reports/', '') as ReportCategory
    }
    return 'billing'
  }, [pathname])

  const activeSetting: SettingsTab = useMemo(() => {
    if (pathname.startsWith('/settings/')) {
      return pathname.replace('/settings/', '') as SettingsTab
    }
    return 'master-data'
  }, [pathname])

  const handleSaveSession = (savedSession: ProjectSession) => {
    if (savedSession.id) {
      sessionsApi.update(savedSession.id, savedSession).then(() => {
        queryClient.invalidateQueries({ queryKey: ['sessions'] })
        toast.success('Session updated successfully')
      }).catch((err) => toast.error(err.message || 'Failed to update session'))
    } else {
      sessionsApi.create(savedSession).then(() => {
        queryClient.invalidateQueries({ queryKey: ['sessions'] })
        toast.success('Session created successfully')
      }).catch((err) => toast.error(err.message || 'Failed to create session'))
    }
  }

  const handleDeleteSession = (sessionId: string) => {
    sessionsApi.delete(sessionId).then(() => {
      queryClient.invalidateQueries({ queryKey: ['sessions'] })
      toast.success('Session deleted')
    }).catch((err) => toast.error(err.message || 'Failed to delete session'))
  }

  const handleSessionStatusChange = (sessionId: string, newStatus: SessionStatus) => {
    sessionsApi.updateStatus(sessionId, newStatus).then(() => {
      queryClient.invalidateQueries({ queryKey: ['sessions'] })
      toast.success(`Session status updated to ${newStatus}`)
    }).catch((err) => toast.error(err.message || 'Failed to update session status'))
  }

  // Lead Navigation Handlers
  const handleNavigateToCreateLead = (stage?: PipelineStage) => {
    setEditingLead(
      stage
        ? ({
          id: '',
          name: '',
          mobileNo: '',
          countryCode: '+91',
          email: '',
          createdDate: new Date().toISOString().split('T')[0],
          leadSource: 'Instagram',
          generate: 'Direct Entry',
          assignee: { name: 'Rohan Mehta', initials: 'RM' },
          leadPipeline: 'Wedding Inquiries',
          pipelineStage: stage,
          leadStage:
            stage === 'contacted'
              ? 'Contacted'
              : stage === 'qualified'
                ? 'Interested'
                : 'New Lead',
          quoteAmount: 0,
          followupDate: '2026-11-11T16:28',
          projectType: 'Wedding',
          quoteType: 'Standard Package',
          projectStartDate: '2026-12-01',
          projectEndDate: '2026-12-02',
          projectSegments: [],
          selectedServices: [],
          selectedDeliverables: ['Raw Photos', 'Pen Drive'],
        } as Lead)
        : null
    )
    navigate('/leads/create')
  }

  const handleNavigateToEditLead = (lead: Lead) => {
    setEditingLead(lead)
    navigate(`/leads/edit/${lead.id}`)
  }

  const handleNavigateToLeadsList = () => {
    setEditingLead(null)
    navigate('/leads')
  }

  // Project Navigation Handlers
  const handleNavigateToCreateProject = (prefill?: Partial<Project>) => {
    setEditingProject(null)
    setProjectPrefillData(prefill || null)
    navigate('/projects/create')
  }

  const handleConvertLeadToProject = (lead: Lead) => {
    const firstSegment = lead.projectSegments?.[0]
    const prefill: Partial<Project> = {
      projectName: lead.projectName || `${lead.name} Shoot`,
      clientName: lead.name,
      clientEmail: lead.email,
      clientMobile: lead.mobileNo,
      projectType: lead.projectType || 'Wedding',
      projectAddress: lead.projectAddress || '',
      startDate: lead.projectStartDate || firstSegment?.startDate || new Date().toISOString().split('T')[0],
      endDate: lead.projectEndDate || firstSegment?.endDate || new Date().toISOString().split('T')[0],
      stage: 'Booked',
      status: 'Booked',
      segments: lead.projectSegments || [
        {
          id: 'seg-1',
          name: lead.projectName || `${lead.name} Shoot`,
          type: lead.projectType || 'Wedding',
          startDate: lead.projectStartDate || new Date().toISOString().split('T')[0],
          endDate: lead.projectEndDate || new Date().toISOString().split('T')[0],
        },
      ],
    }
    setEditingProject(null)
    setProjectPrefillData(prefill)
    navigate('/projects/create')
    toast.success(`Pre-filled project creation from Lead "${lead.name}"`)
  }

  const handleNavigateToEditProject = (project: Project) => {
    setEditingProject(project)
    setProjectPrefillData(null)
    navigate(`/projects/edit/${project.id}`)
  }

  const handleNavigateToProjectsList = () => {
    setEditingProject(null)
    setProjectPrefillData(null)
    navigate('/projects')
  }

  // Lead CRUD
  const handleSaveLead = (leadToSave: Lead) => {
    if (leadToSave.id) {
      leadsApi.update(leadToSave.id, leadToSave).then((saved) => {
        queryClient.invalidateQueries({ queryKey: ['leads'] })
        queryClient.invalidateQueries({ queryKey: ['lead', leadToSave.id] })
        toast.success(`Lead "${saved.name}" updated successfully`)
      }).catch((err) => toast.error(err.message || 'Failed to update lead'))
    } else {
      leadsApi.create(leadToSave).then((saved) => {
        queryClient.invalidateQueries({ queryKey: ['leads'] })
        toast.success(`Lead "${saved.name}" created successfully`)
      }).catch((err) => toast.error(err.message || 'Failed to create lead'))
    }
  }

  const handleDeleteLead = (leadId: string) => {
    leadsApi.delete(leadId).then(() => {
      queryClient.invalidateQueries({ queryKey: ['leads'] })
      queryClient.invalidateQueries({ queryKey: ['lead', leadId] })
      toast.success('Lead deleted successfully')
    }).catch((err) => toast.error(err.message || 'Failed to delete lead'))
  }

  const handleLeadStageChange = (leadId: string, newStage: PipelineStage) => {
    let newSubStage = 'New Lead'
    if (newStage === 'contacted') newSubStage = 'Contacted'
    else if (newStage === 'qualified') newSubStage = 'Negotiation'
    else if (newStage === 'converted') newSubStage = 'Booked'

    leadsApi.updateStage(leadId, newStage, newSubStage as any).then(() => {
      queryClient.invalidateQueries({ queryKey: ['leads'] })
      queryClient.invalidateQueries({ queryKey: ['lead', leadId] })
      toast.success(`Lead moved to ${newStage.replace('_', ' ').toUpperCase()}`)
    }).catch((err) => toast.error(err.message || 'Failed to update lead stage'))
  }

  // Project CRUD
  const handleSaveProject = (projectToSave: Project) => {
    if (projectToSave.id) {
      projectsApi.update(projectToSave.id, projectToSave).then((saved) => {
        queryClient.invalidateQueries({ queryKey: ['projects'] })
        queryClient.invalidateQueries({ queryKey: ['project', projectToSave.id] })
        toast.success(`Project "${saved.projectName}" updated successfully`)
      }).catch((err) => toast.error(err.message || 'Failed to update project'))
    } else {
      projectsApi.create(projectToSave).then((saved) => {
        queryClient.invalidateQueries({ queryKey: ['projects'] })
        toast.success(`Project "${saved.projectName}" created successfully`)
      }).catch((err) => toast.error(err.message || 'Failed to create project'))
    }
  }

  const handleDeleteProject = (projectId: string) => {
    projectsApi.delete(projectId).then(() => {
      queryClient.invalidateQueries({ queryKey: ['projects'] })
      queryClient.invalidateQueries({ queryKey: ['project', projectId] })
      toast.success('Project deleted successfully')
    }).catch((err) => toast.error(err.message || 'Failed to delete project'))
  }

  const handleProjectStageChange = (projectId: string, newStage: ProjectStage) => {
    projectsApi.updateStage(projectId, newStage).then(() => {
      queryClient.invalidateQueries({ queryKey: ['projects'] })
      queryClient.invalidateQueries({ queryKey: ['project', projectId] })
      toast.success(`Project moved to stage ${newStage}`)
    }).catch((err) => toast.error(err.message || 'Failed to update project stage'))
  }

  // Order Handlers
  const handleNavigateToCreateOrder = () => {
    setEditingOrder(null)
    navigate('/orders/create')
  }

  const handleNavigateToEditOrder = (order: Order) => {
    setEditingOrder(order)
    navigate(`/orders/edit/${order.id}`)
  }

  const handleNavigateToOrdersList = () => {
    setEditingOrder(null)
    navigate('/orders')
  }

  const handleSaveOrder = (savedOrder: Order) => {
    if (savedOrder.id) {
      ordersApi.update(savedOrder.id, savedOrder).then(() => {
        queryClient.invalidateQueries({ queryKey: ['orders'] })
        queryClient.invalidateQueries({ queryKey: ['order', savedOrder.id] })
        toast.success('Order updated successfully')
      }).catch((err) => toast.error(err.message || 'Failed to update order'))
    } else {
      ordersApi.create(savedOrder).then(() => {
        queryClient.invalidateQueries({ queryKey: ['orders'] })
        toast.success('Order created successfully')
      }).catch((err) => toast.error(err.message || 'Failed to create order'))
    }
    navigate('/orders')
  }

  const handleDeleteOrder = (orderId: string) => {
    ordersApi.delete(orderId).then(() => {
      queryClient.invalidateQueries({ queryKey: ['orders'] })
      queryClient.invalidateQueries({ queryKey: ['order', orderId] })
      toast.success('Order deleted successfully')
    }).catch((err) => toast.error(err.message || 'Failed to delete order'))
  }

  const handleOrderStageChange = (orderId: string, newStage: OrderStage) => {
    ordersApi.updateStage(orderId, newStage).then(() => {
      queryClient.invalidateQueries({ queryKey: ['orders'] })
      queryClient.invalidateQueries({ queryKey: ['order', orderId] })
      toast.success(`Order stage updated to ${newStage}`)
    }).catch((err) => toast.error(err.message || 'Failed to update order stage'))
  }

  // Dynamic Interactive Breadcrumb Def
  const breadcrumbItems = useMemo(() => {
    if (!isAdmin) {
      if (pathname.startsWith('/shoots')) {
        return [{ label: 'Shoots' }]
      }
      if (pathname.startsWith('/tasks')) {
        return [{ label: 'Tasks' }]
      }
      if (pathname.startsWith('/deliverables')) {
        return [{ label: 'Deliverables' }]
      }
      return [{ label: 'Home' }]
    }

    if (pathname.startsWith('/settings')) {
      const titles: Record<SettingsTab, string> = {
        'master-data': 'Master Data',
        photographers: 'Photographers & Crew',
        'sequence-id': 'Sequence ID',
        'preset-services': 'Services Preset',
        'preset-packages': 'Packages Preset',
        'preset-orders': 'Orders Preset',
        'preset-deliverables': 'Deliverables Preset',
        'template-forms': 'Forms Template',
        'template-quotation-pdf': 'Quotation PDF Template',
        'template-invoice-pdf': 'Invoice PDF Template',
        'template-contracts': 'Contracts Template',
        'template-feedbacks': 'Feedbacks Template',
        'stages-workflow': 'Production Workflow',
      }
      return [
        { label: 'Studio Settings', href: '/settings/master-data' },
        { label: titles[activeSetting] || 'Settings' },
      ]
    }

    if (pathname.startsWith('/reports')) {
      const titles: Record<ReportCategory, string> = {
        billing: 'Billing Statement',
        leads: 'Leads Funnel',
        orders: 'Orders & Merchandise',
        sessions: 'Sessions & Crew',
        finance: 'Finance P&L',
        'master-report': 'Master Register',
        'account-analytics': 'Account Analytics',
        feedback: 'Client Feedback',
      }
      return [
        { label: 'Reports & BI', href: '/reports/billing' },
        { label: titles[activeReport] || 'Report' },
      ]
    }

    if (pathname === '/leads/create') {
      return [
        { label: 'Leads', href: '/leads' },
        { label: 'Create Lead' },
      ]
    }
    if (pathname.startsWith('/leads/edit')) {
      return [
        { label: 'Leads', href: '/leads' },
        { label: editingLead?.name ? `Lead: ${editingLead.name}` : 'Edit Lead' },
      ]
    }
    if (pathname.startsWith('/leads')) {
      return [{ label: 'Photography Leads' }]
    }

    if (pathname === '/projects/create') {
      return [
        { label: 'Projects', href: '/projects' },
        { label: 'Create Project' },
      ]
    }
    if (pathname.startsWith('/projects/edit')) {
      return [
        { label: 'Projects', href: '/projects' },
        { label: editingProject?.projectName ? `Project: ${editingProject.projectName}` : 'Edit Project' },
      ]
    }
    if (pathname.startsWith('/projects')) {
      return [{ label: 'Projects' }]
    }

    if (pathname.startsWith('/sessions')) {
      return [{ label: 'Sessions' }]
    }
    if (pathname.startsWith('/billings')) {
      return [{ label: 'Billings' }]
    }

    if (pathname === '/orders/create') {
      return [
        { label: 'Orders', href: '/orders' },
        { label: 'New Order' },
      ]
    }
    if (pathname.startsWith('/orders/edit')) {
      return [
        { label: 'Orders', href: '/orders' },
        { label: editingOrder?.orderNumber ? `Order: ${editingOrder.orderNumber}` : 'Edit Order' },
      ]
    }
    if (pathname.startsWith('/orders')) {
      return [{ label: 'Orders' }]
    }

    return [{ label: 'Projects' }]
  }, [isAdmin, activeSetting, activeReport, pathname, editingLead, editingProject, editingOrder])

  return (
    <DashboardShell
      title={isAdmin ? "Photography CRM" : "Photography Portal"}
      breadcrumb={breadcrumbItems}
      sidebar={
        isAdmin ? (
          <PhotographySidebar
            currentRoute={currentRoute}
            activeReport={activeReport}
            activeSetting={activeSetting}
            onSelectRoute={(r) => {
              const routeMap: Record<string, string> = {
                'projects-list': '/projects',
                'sessions-list': '/sessions',
                'leads-list': '/leads',
                'billings': '/billings',
                'orders-list': '/orders',
              }
              navigate(routeMap[r] || '/projects')
            }}
            onSelectReport={(rep) => {
              navigate(`/reports/${rep}`)
            }}
            onSelectSetting={(set) => {
              navigate(`/settings/${set}`)
            }}
          />
        ) : (
          <UserSidebar
            currentRoute={currentUserRoute}
            onSelectRoute={(r) => {
              const userRouteMap: Record<UserRouteKey, string> = {
                home: '/home',
                shoots: '/shoots',
                tasks: '/tasks',
                deliverables: '/deliverables',
              }
              navigate(userRouteMap[r] || '/home')
            }}
          />
        )
      }
    >
      <Toaster position="bottom-right" richColors />

      {/* Routes Switcher with Code-Splitting Suspense */}
      <Suspense
        fallback={
          <div className="flex flex-1 items-center justify-center p-8 text-xs text-muted-foreground">
            Loading view...
          </div>
        }
      >
        <Routes>
          {isAdmin ? (
            <>
              {/* Default redirect for Admin */}
              <Route path="/" element={<Navigate to="/projects" replace />} />

              {/* Core Modules */}
              <Route
                path="/projects"
                element={
                  <ProjectsPageContainer
                    onAddNewProject={handleNavigateToCreateProject}
                    onEditProject={handleNavigateToEditProject}
                    onDeleteProject={handleDeleteProject}
                    onStageChange={handleProjectStageChange}
                  />
                }
              />
              <Route
                path="/projects/create"
                element={
                  <ProjectCreatePage
                    initialData={projectPrefillData || undefined}
                    onBack={handleNavigateToProjectsList}
                    onSave={handleSaveProject}
                  />
                }
              />
              <Route
                path="/projects/edit/:id"
                element={
                  <ProjectEditPageWrapper
                    editingProject={editingProject}
                    onBack={handleNavigateToProjectsList}
                    onSave={handleSaveProject}
                    onDelete={handleDeleteProject}
                  />
                }
              />

              <Route
                path="/sessions"
                element={
                  <SessionsPageContainer
                    onDeleteSession={handleDeleteSession}
                    onStatusChange={handleSessionStatusChange}
                    onSaveSession={handleSaveSession}
                  />
                }
              />

              <Route
                path="/leads"
                element={
                  <LeadsPageContainer
                    onAddNewLead={handleNavigateToCreateLead}
                    onEditLead={handleNavigateToEditLead}
                    onDeleteLead={handleDeleteLead}
                    onStageChange={handleLeadStageChange}
                    onConvertToProject={handleConvertLeadToProject}
                  />
                }
              />
              <Route
                path="/leads/create"
                element={
                  <LeadFormPage
                    initialLead={editingLead}
                    onBack={handleNavigateToLeadsList}
                    onSave={handleSaveLead}
                  />
                }
              />
              <Route
                path="/leads/edit/:id"
                element={
                  <LeadEditPageWrapper
                    editingLead={editingLead}
                    onBack={handleNavigateToLeadsList}
                    onSave={handleSaveLead}
                  />
                }
              />

              <Route path="/billings" element={<BillingsPage />} />

              <Route
                path="/orders"
                element={
                  <OrdersPageContainer
                    onAddNewOrder={handleNavigateToCreateOrder}
                    onEditOrder={handleNavigateToEditOrder}
                    onDeleteOrder={handleDeleteOrder}
                    onStageChange={handleOrderStageChange}
                  />
                }
              />
              <Route
                path="/orders/create"
                element={
                  <OrderFormPageWrapper
                    initialOrder={editingOrder}
                    onBack={handleNavigateToOrdersList}
                    onSave={handleSaveOrder}
                  />
                }
              />
              <Route
                path="/orders/edit/:id"
                element={
                  <OrderEditPageWrapper
                    editingOrder={editingOrder}
                    onBack={handleNavigateToOrdersList}
                    onSave={handleSaveOrder}
                  />
                }
              />

              {/* Reports Routes */}
              <Route path="/reports/billing" element={<BillingReportPage />} />
              <Route path="/reports/leads" element={<LeadsReportPage />} />
              <Route path="/reports/orders" element={<OrdersReportPage />} />
              <Route path="/reports/sessions" element={<SessionsReportPage />} />
              <Route path="/reports/finance" element={<FinanceReportPage />} />
              <Route path="/reports/master-report" element={<MasterReportPage />} />
              <Route path="/reports/account-analytics" element={<AccountAnalyticsPage />} />
              <Route path="/reports/feedback" element={<FeedbackReportPage />} />

              {/* Settings Routes */}
              <Route path="/settings/master-data" element={<MasterDataPage />} />
              <Route path="/settings/photographers" element={<PhotographersPage />} />
              <Route path="/settings/sequence-id" element={<SequenceIdPage />} />
              <Route path="/settings/preset-services" element={<ServicesSettingsPage />} />
              <Route path="/settings/preset-packages" element={<PackagesSettingsPage />} />
              <Route path="/settings/preset-orders" element={<OrdersSettingsPage />} />
              <Route path="/settings/preset-deliverables" element={<DeliverablesSettingsPage />} />
              <Route
                path="/settings/template-forms"
                element={
                  <TemplateCustomizerPage
                    title="Lead Intake & Inquiry Forms"
                    subtitle="Customize public lead capture form questions, branding, and disclaimers"
                    type="form"
                  />
                }
              />
              <Route
                path="/settings/template-quotation-pdf"
                element={
                  <TemplateCustomizerPage
                    title="Quotation PDF Template"
                    subtitle="Visual design, payment schedules, and terms on client quotation estimates"
                    type="quotation"
                  />
                }
              />
              <Route
                path="/settings/template-invoice-pdf"
                element={
                  <TemplateCustomizerPage
                    title="Invoice PDF Template"
                    subtitle="Tax invoice layout, bank details, and studio payment remittance instructions"
                    type="invoice"
                  />
                }
              />
              <Route path="/settings/template-contracts" element={<ContractTemplatePage />} />
              <Route path="/settings/template-feedbacks" element={<FeedbackTemplatePage />} />
              <Route path="/settings/stages-workflow" element={<WorkflowPage />} />

              {/* Catch-all redirect for Admin */}
              <Route path="*" element={<Navigate to="/projects" replace />} />
            </>
          ) : (
            <>
              {/* Default redirect for non-Admin users */}
              <Route path="/" element={<Navigate to="/home" replace />} />

              {/* User Portal Pages */}
              <Route
                path="/home"
                element={<UserHomePageContainer />}
              />
              <Route
                path="/shoots"
                element={
                  <UserShootsPageContainer
                    onSaveSession={handleSaveSession}
                    onStatusChange={handleSessionStatusChange}
                  />
                }
              />
              <Route
                path="/tasks"
                element={
                  <UserTasksPageContainer
                    onUpdateProject={handleSaveProject}
                  />
                }
              />
              <Route
                path="/deliverables"
                element={
                  <UserDeliverablesPageContainer
                    onUpdateProject={handleSaveProject}
                  />
                }
              />

              {/* Catch-all redirect for non-Admin */}
              <Route path="*" element={<Navigate to="/home" replace />} />
            </>
          )}
        </Routes>
      </Suspense>
    </DashboardShell>
  )
}

function LeadEditPageWrapper({
  editingLead,
  onBack,
  onSave,
}: {
  editingLead: Lead | null
  onBack: () => void
  onSave: (l: Lead) => void
}) {
  const { id } = useParams<{ id: string }>()
  const { data: fetchedLead, isLoading } = useQuery({
    queryKey: ['lead', id],
    queryFn: ({ signal }) => (id ? leadsApi.get(id, signal) : Promise.reject('No ID')),
    enabled: !editingLead && !!id,
    staleTime: 60_000,
  })

  const targetLead = editingLead || fetchedLead || null

  if (!editingLead && isLoading) {
    return (
      <div className="flex flex-1 items-center justify-center p-8 text-xs text-muted-foreground">
        Loading lead details...
      </div>
    )
  }

  return <LeadFormPage initialLead={targetLead} onBack={onBack} onSave={onSave} />
}

function ProjectEditPageWrapper({
  editingProject,
  onBack,
  onSave,
  onDelete,
}: {
  editingProject: Project | null
  onBack: () => void
  onSave: (p: Project) => void
  onDelete: (id: string) => void
}) {
  const { id } = useParams<{ id: string }>()
  const { data: fetchedProject, isLoading } = useQuery({
    queryKey: ['project', id],
    queryFn: ({ signal }) => (id ? projectsApi.get(id, signal) : Promise.reject('No ID')),
    enabled: !editingProject && !!id,
    staleTime: 60_000,
  })

  const targetProject = editingProject || fetchedProject || null

  if (!editingProject && isLoading) {
    return (
      <div className="flex flex-1 items-center justify-center p-8 text-xs text-muted-foreground">
        Loading project details...
      </div>
    )
  }

  if (!targetProject) {
    return <Navigate to="/projects" replace />
  }

  return (
    <ProjectEditPage
      project={targetProject}
      onBack={onBack}
      onUpdate={onSave}
      onDelete={onDelete}
    />
  )
}

function OrderEditPageWrapper({
  editingOrder,
  onBack,
  onSave,
}: {
  editingOrder: Order | null
  onBack: () => void
  onSave: (o: Order) => void
}) {
  const { id } = useParams<{ id: string }>()
  const { data: fetchedOrder, isLoading } = useQuery({
    queryKey: ['order', id],
    queryFn: ({ signal }) => (id ? ordersApi.get(id, signal) : Promise.reject('No ID')),
    enabled: !editingOrder && !!id,
    staleTime: 60_000,
  })

  const { data: projects = [] } = useQuery({
    queryKey: ['projects'],
    queryFn: ({ signal }) => projectsApi.list(undefined, signal),
    staleTime: 60_000,
    gcTime: 5 * 60_000,
  })

  const targetOrder = editingOrder || fetchedOrder || null

  if (!editingOrder && isLoading) {
    return (
      <div className="flex flex-1 items-center justify-center p-8 text-xs text-muted-foreground">
        Loading order details...
      </div>
    )
  }

  return (
    <OrderFormPage
      initialOrder={targetOrder}
      onBack={onBack}
      onSave={onSave}
      existingProjectOptions={projects.map((p) => ({ id: p.id, name: p.projectName }))}
    />
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Suspense
          fallback={
            <div className="flex min-h-screen items-center justify-center text-xs text-muted-foreground">
              Loading...
            </div>
          }
        >
          <Routes>
            <Route path="/signin" element={<SigninPage />} />
            <Route path="/signup" element={<SignupPage />} />
            <Route
              path="/*"
              element={
                <ProtectedRoute fallback={<Navigate to="/signin" replace />}>
                  <AppContent />
                </ProtectedRoute>
              }
            />
          </Routes>
        </Suspense>
      </AuthProvider>
    </BrowserRouter>
  )
}
