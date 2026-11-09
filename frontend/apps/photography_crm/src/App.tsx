import { useState, useEffect, useMemo, Suspense, lazy } from 'react'
import {
  BrowserRouter,
  Routes,
  Route,
  Navigate,
  useNavigate,
  useLocation,
  useParams,
} from 'react-router-dom'
import {
  DashboardShell,
  Toaster,
  toast,
  AuthProvider,
  ProtectedRoute,
} from '@repo/ui'
import SigninPage from './pages/auth/SigninPage'
import SignupPage from './pages/auth/SignupPage'
import type { Lead, PipelineStage } from './types/lead'
import type { Project, ProjectStage } from './types/project'
import type { ProjectSession, SessionStatus } from './types/session'
import type { Order, OrderStage } from './types/order'
import type { SettingsTab } from './types/settings'
import { PhotographySidebar, type ReportCategory } from './components/PhotographySidebar'
import { leadsApi, projectsApi, sessionsApi, ordersApi } from './api/crmApi'

// Lazy Load Core Pages
const LeadsPage = lazy(() =>
  import('./pages/LeadsPage').then((m) => ({ default: m.LeadsPage }))
)
const LeadFormPage = lazy(() =>
  import('./pages/LeadFormPage').then((m) => ({ default: m.LeadFormPage }))
)
const ProjectsPage = lazy(() =>
  import('./pages/ProjectsPage').then((m) => ({ default: m.ProjectsPage }))
)
const ProjectCreatePage = lazy(() =>
  import('./pages/ProjectCreatePage').then((m) => ({ default: m.ProjectCreatePage }))
)
const ProjectEditPage = lazy(() =>
  import('./pages/ProjectEditPage').then((m) => ({ default: m.ProjectEditPage }))
)
const SessionsPage = lazy(() =>
  import('./pages/SessionsPage').then((m) => ({ default: m.SessionsPage }))
)
const BillingsPage = lazy(() =>
  import('./pages/BillingsPage').then((m) => ({ default: m.BillingsPage }))
)
const OrdersPage = lazy(() =>
  import('./pages/OrdersPage').then((m) => ({ default: m.OrdersPage }))
)
const OrderFormPage = lazy(() =>
  import('./pages/OrderFormPage').then((m) => ({ default: m.OrderFormPage }))
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

function AppContent() {
  const navigate = useNavigate()
  const location = useLocation()
  const pathname = location.pathname

  const [leads, setLeads] = useState<Lead[]>([])
  const [editingLead, setEditingLead] = useState<Lead | null>(null)

  const [projects, setProjects] = useState<Project[]>([])
  const [editingProject, setEditingProject] = useState<Project | null>(null)
  const [projectPrefillData, setProjectPrefillData] = useState<Partial<Project> | null>(null)

  const [orders, setOrders] = useState<Order[]>([])
  const [editingOrder, setEditingOrder] = useState<Order | null>(null)

  const [sessions, setSessions] = useState<ProjectSession[]>([])

  useEffect(() => {
    leadsApi.list().then(setLeads).catch((err) => console.error('Failed to load leads:', err))
    projectsApi.list().then(setProjects).catch((err) => console.error('Failed to load projects:', err))
    sessionsApi.list().then(setSessions).catch((err) => console.error('Failed to load sessions:', err))
    ordersApi.list().then(setOrders).catch((err) => console.error('Failed to load orders:', err))
  }, [])

  // Derive current route key for sidebar active states
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

  // Session Handlers
  const handleSaveSession = (savedSession: ProjectSession) => {
    const exists = sessions.some((s) => s.id === savedSession.id)
    if (exists) {
      sessionsApi.update(savedSession.id, savedSession).then((updated) => {
        setSessions((prev) => prev.map((s) => (s.id === updated.id ? updated : s)))
      }).catch((err) => toast.error(err.message || 'Failed to update session'))
    } else {
      sessionsApi.create(savedSession).then((created) => {
        setSessions((prev) => [created, ...prev])
      }).catch((err) => toast.error(err.message || 'Failed to create session'))
    }
  }

  const handleDeleteSession = (sessionId: string) => {
    sessionsApi.delete(sessionId).then(() => {
      setSessions((prev) => prev.filter((s) => s.id !== sessionId))
      toast.success('Session deleted')
    }).catch((err) => toast.error(err.message || 'Failed to delete session'))
  }

  const handleSessionStatusChange = (sessionId: string, newStatus: SessionStatus) => {
    sessionsApi.updateStatus(sessionId, newStatus).then(() => {
      setSessions((prev) =>
        prev.map((s) => (s.id === sessionId ? { ...s, status: newStatus } : s))
      )
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
    const exists = leads.some((l) => l.id === leadToSave.id)
    if (exists) {
      leadsApi.update(leadToSave.id, leadToSave).then((saved) => {
        setLeads((prev) => prev.map((l) => (l.id === saved.id ? saved : l)))
        toast.success(`Lead "${saved.name}" updated successfully`)
      }).catch((err) => toast.error(err.message || 'Failed to update lead'))
    } else {
      leadsApi.create(leadToSave).then((saved) => {
        setLeads((prev) => [saved, ...prev])
        toast.success(`Lead "${saved.name}" created successfully`)
      }).catch((err) => toast.error(err.message || 'Failed to create lead'))
    }
  }

  const handleDeleteLead = (leadId: string) => {
    leadsApi.delete(leadId).then(() => {
      setLeads((prev) => prev.filter((l) => l.id !== leadId))
      toast.success('Lead deleted successfully')
    }).catch((err) => toast.error(err.message || 'Failed to delete lead'))
  }

  const handleLeadStageChange = (leadId: string, newStage: PipelineStage) => {
    let newSubStage = 'New Lead'
    if (newStage === 'contacted') newSubStage = 'Contacted'
    else if (newStage === 'qualified') newSubStage = 'Negotiation'
    else if (newStage === 'converted') newSubStage = 'Booked'

    leadsApi.updateStage(leadId, newStage, newSubStage as any).then(() => {
      setLeads((prev) =>
        prev.map((lead) => {
          if (lead.id !== leadId) return lead
          return {
            ...lead,
            pipelineStage: newStage,
            leadStage: newSubStage as any,
          }
        })
      )
      toast.success(`Lead moved to ${newStage.replace('_', ' ').toUpperCase()}`)
    }).catch((err) => toast.error(err.message || 'Failed to update lead stage'))
  }

  // Project CRUD
  const handleSaveProject = (projectToSave: Project) => {
    const exists = projects.some((p) => p.id === projectToSave.id)
    if (exists) {
      projectsApi.update(projectToSave.id, projectToSave).then((saved) => {
        setProjects((prev) => prev.map((p) => (p.id === saved.id ? saved : p)))
        toast.success(`Project "${saved.projectName}" updated successfully`)
      }).catch((err) => toast.error(err.message || 'Failed to update project'))
    } else {
      projectsApi.create(projectToSave).then((saved) => {
        setProjects((prev) => [saved, ...prev])
        toast.success(`Project "${saved.projectName}" created successfully`)
      }).catch((err) => toast.error(err.message || 'Failed to create project'))
    }
  }

  const handleDeleteProject = (projectId: string) => {
    projectsApi.delete(projectId).then(() => {
      setProjects((prev) => prev.filter((p) => p.id !== projectId))
      toast.success('Project deleted successfully')
    }).catch((err) => toast.error(err.message || 'Failed to delete project'))
  }

  const handleProjectStageChange = (projectId: string, newStage: ProjectStage) => {
    projectsApi.updateStage(projectId, newStage).then(() => {
      setProjects((prev) =>
        prev.map((p) => {
          if (p.id !== projectId) return p
          return { ...p, stage: newStage }
        })
      )
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
    const exists = orders.some((o) => o.id === savedOrder.id)
    if (exists) {
      ordersApi.update(savedOrder.id, savedOrder).then((updated) => {
        setOrders((prev) => prev.map((o) => (o.id === updated.id ? updated : o)))
        toast.success('Order updated successfully')
      }).catch((err) => toast.error(err.message || 'Failed to update order'))
    } else {
      ordersApi.create(savedOrder).then((created) => {
        setOrders((prev) => [created, ...prev])
        toast.success('Order created successfully')
      }).catch((err) => toast.error(err.message || 'Failed to create order'))
    }
    navigate('/orders')
  }

  const handleDeleteOrder = (orderId: string) => {
    ordersApi.delete(orderId).then(() => {
      setOrders((prev) => prev.filter((o) => o.id !== orderId))
      toast.success('Order deleted successfully')
    }).catch((err) => toast.error(err.message || 'Failed to delete order'))
  }

  const handleOrderStageChange = (orderId: string, newStage: OrderStage) => {
    ordersApi.updateStage(orderId, newStage).then(() => {
      setOrders((prev) =>
        prev.map((o) => {
          if (o.id !== orderId) return o
          return { ...o, stage: newStage }
        })
      )
      toast.success(`Order stage updated to ${newStage}`)
    }).catch((err) => toast.error(err.message || 'Failed to update order stage'))
  }

  // Dynamic Interactive Breadcrumb Def
  const breadcrumbItems = useMemo(() => {
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
  }, [activeSetting, activeReport, pathname, editingLead, editingProject, editingOrder])

  return (
    <DashboardShell
      title="Photography CRM"
      breadcrumb={breadcrumbItems}
      sidebar={
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
          {/* Default redirect to projects */}
          <Route path="/" element={<Navigate to="/projects" replace />} />

          {/* Core Modules */}
          <Route
            path="/projects"
            element={
              <ProjectsPage
                projects={projects}
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
                projects={projects}
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
              <SessionsPage
                sessions={sessions}
                onAddNewSession={() => { }}
                onEditSession={() => { }}
                onDeleteSession={handleDeleteSession}
                onStatusChange={handleSessionStatusChange}
                onSaveSession={handleSaveSession}
                existingProjectNames={projects.map((p) => p.projectName)}
              />
            }
          />

          <Route
            path="/leads"
            element={
              <LeadsPage
                leads={leads}
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
                leads={leads}
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
              <OrdersPage
                orders={orders}
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
              <OrderFormPage
                initialOrder={editingOrder}
                onBack={handleNavigateToOrdersList}
                onSave={handleSaveOrder}
                existingProjectOptions={projects.map((p) => ({ id: p.id, name: p.projectName }))}
              />
            }
          />
          <Route
            path="/orders/edit/:id"
            element={
              <OrderEditPageWrapper
                orders={orders}
                editingOrder={editingOrder}
                onBack={handleNavigateToOrdersList}
                onSave={handleSaveOrder}
                existingProjectOptions={projects.map((p) => ({ id: p.id, name: p.projectName }))}
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

          {/* Catch-all redirect */}
          <Route path="*" element={<Navigate to="/projects" replace />} />
        </Routes>
      </Suspense>
    </DashboardShell>
  )
}

function LeadEditPageWrapper({
  leads,
  editingLead,
  onBack,
  onSave,
}: {
  leads: Lead[]
  editingLead: Lead | null
  onBack: () => void
  onSave: (l: Lead) => void
}) {
  const { id } = useParams<{ id: string }>()
  const targetLead = editingLead || leads.find((l) => l.id === id) || null
  return <LeadFormPage initialLead={targetLead} onBack={onBack} onSave={onSave} />
}

function ProjectEditPageWrapper({
  projects,
  editingProject,
  onBack,
  onSave,
  onDelete,
}: {
  projects: Project[]
  editingProject: Project | null
  onBack: () => void
  onSave: (p: Project) => void
  onDelete: (id: string) => void
}) {
  const { id } = useParams<{ id: string }>()
  const targetProject = editingProject || projects.find((p) => p.id === id) || null

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
  orders,
  editingOrder,
  onBack,
  onSave,
  existingProjectOptions,
}: {
  orders: Order[]
  editingOrder: Order | null
  onBack: () => void
  onSave: (o: Order) => void
  existingProjectOptions: { id: string; name: string }[]
}) {
  const { id } = useParams<{ id: string }>()
  const targetOrder = editingOrder || orders.find((o) => o.id === id) || null

  return (
    <OrderFormPage
      initialOrder={targetOrder}
      onBack={onBack}
      onSave={onSave}
      existingProjectOptions={existingProjectOptions}
    />
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
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
      </AuthProvider>
    </BrowserRouter>
  )
}
