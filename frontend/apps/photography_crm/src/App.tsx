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
import {
  DashboardShell,
  Toaster,
  Button,
  toast,
} from '@repo/ui'
import {
  IconReportAnalytics,
  IconSettings,
  IconFolders,
  IconCalendar,
  IconTarget,
  IconCreditCard,
  IconShoppingBag,
} from '@tabler/icons-react'
import { INITIAL_LEADS } from './data/mock-leads'
import { INITIAL_PROJECTS } from './data/mock-projects'
import type { Lead, PipelineStage } from './types/lead'
import type { Project, ProjectStage } from './types/project'
import type { ProjectSession, SessionStatus } from './types/session'
import type { Order, OrderStage } from './types/order'
import type { SettingsTab } from './types/settings'
import { INITIAL_ORDERS } from './data/mock-orders'
import { PhotographySidebar, type ReportCategory } from './components/PhotographySidebar'

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
const GenericPresetPage = lazy(() =>
  import('./pages/settings/GenericPresetPage').then((m) => ({ default: m.GenericPresetPage }))
)
const ServicesSettingsPage = lazy(() =>
  import('./pages/settings/ServicesSettingsPage').then((m) => ({ default: m.ServicesSettingsPage }))
)
const PackagesSettingsPage = lazy(() =>
  import('./pages/settings/PackagesSettingsPage').then((m) => ({ default: m.PackagesSettingsPage }))
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

  const [leads, setLeads] = useState<Lead[]>(INITIAL_LEADS)
  const [editingLead, setEditingLead] = useState<Lead | null>(null)

  const [projects, setProjects] = useState<Project[]>(INITIAL_PROJECTS)
  const [editingProject, setEditingProject] = useState<Project | null>(null)

  const [orders, setOrders] = useState<Order[]>(INITIAL_ORDERS)
  const [editingOrder, setEditingOrder] = useState<Order | null>(null)

  const [sessions, setSessions] = useState<ProjectSession[]>(() => {
    return INITIAL_PROJECTS.flatMap((p) =>
      p.sessions.map((s) => ({
        ...s,
        projectName: s.projectName || p.projectName,
        projectType: s.projectType || p.projectType,
        sessionManager: s.sessionManager || 'especa',
        shootTiming: s.shootTiming || 'Evening',
        duration: s.duration || '4 hours',
      }))
    )
  })

  // Derive sidebar mode and active subtabs from pathname
  const sidebarMode: 'main' | 'reports' | 'settings' = useMemo(() => {
    if (pathname.startsWith('/settings')) return 'settings'
    if (pathname.startsWith('/reports')) return 'reports'
    return 'main'
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

  // Derive currentRoute key for main sidebar active state
  const currentMainRoute = useMemo(() => {
    if (pathname.startsWith('/projects')) return 'projects-list'
    if (pathname.startsWith('/sessions')) return 'sessions-list'
    if (pathname.startsWith('/leads')) return 'leads-list'
    if (pathname.startsWith('/billings')) return 'billings'
    if (pathname.startsWith('/orders')) return 'orders-list'
    return 'projects-list'
  }, [pathname])

  // Session Handlers
  const handleSaveSession = (savedSession: ProjectSession) => {
    setSessions((prev) => {
      const exists = prev.some((s) => s.id === savedSession.id)
      if (exists) {
        return prev.map((s) => (s.id === savedSession.id ? savedSession : s))
      }
      return [savedSession, ...prev]
    })

    if (savedSession.projectName) {
      setProjects((prev) =>
        prev.map((proj) => {
          if (proj.projectName === savedSession.projectName) {
            const sesExists = proj.sessions.some((s) => s.id === savedSession.id)
            const updatedSessions = sesExists
              ? proj.sessions.map((s) => (s.id === savedSession.id ? savedSession : s))
              : [...proj.sessions, savedSession]
            return {
              ...proj,
              sessions: updatedSessions,
              sessionsCount: updatedSessions.length,
            }
          }
          return proj
        })
      )
    }
  }

  const handleDeleteSession = (sessionId: string) => {
    setSessions((prev) => prev.filter((s) => s.id !== sessionId))
    setProjects((prev) =>
      prev.map((proj) => ({
        ...proj,
        sessions: proj.sessions.filter((s) => s.id !== sessionId),
        sessionsCount: proj.sessions.filter((s) => s.id !== sessionId).length,
      }))
    )
  }

  const handleSessionStatusChange = (sessionId: string, newStatus: SessionStatus) => {
    setSessions((prev) =>
      prev.map((s) => (s.id === sessionId ? { ...s, status: newStatus } : s))
    )
    toast.success(`Session status updated to ${newStatus}`)
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
  const handleNavigateToCreateProject = () => {
    setEditingProject(null)
    navigate('/projects/create')
  }

  const handleNavigateToEditProject = (project: Project) => {
    setEditingProject(project)
    navigate(`/projects/edit/${project.id}`)
  }

  const handleNavigateToProjectsList = () => {
    setEditingProject(null)
    navigate('/projects')
  }

  // Lead CRUD
  const handleSaveLead = (leadToSave: Lead) => {
    setLeads((prev) => {
      const exists = prev.some((l) => l.id === leadToSave.id)
      if (exists) {
        return prev.map((l) => (l.id === leadToSave.id ? leadToSave : l))
      }
      return [leadToSave, ...prev]
    })
  }

  const handleDeleteLead = (leadId: string) => {
    setLeads((prev) => prev.filter((l) => l.id !== leadId))
  }

  const handleLeadStageChange = (leadId: string, newStage: PipelineStage) => {
    setLeads((prev) =>
      prev.map((lead) => {
        if (lead.id !== leadId) return lead
        let newSubStage = lead.leadStage
        if (newStage === 'new_leads') newSubStage = 'New Lead'
        else if (newStage === 'contacted') newSubStage = 'Contacted'
        else if (newStage === 'qualified') newSubStage = 'Negotiation'
        else if (newStage === 'converted') newSubStage = 'Booked'

        return {
          ...lead,
          pipelineStage: newStage,
          leadStage: newSubStage,
        }
      })
    )
    toast.success(`Lead moved to ${newStage.replace('_', ' ').toUpperCase()}`)
  }

  // Project CRUD
  const handleSaveProject = (projectToSave: Project) => {
    setProjects((prev) => {
      const exists = prev.some((p) => p.id === projectToSave.id)
      if (exists) {
        return prev.map((p) => (p.id === projectToSave.id ? projectToSave : p))
      }
      return [projectToSave, ...prev]
    })
  }

  const handleDeleteProject = (projectId: string) => {
    setProjects((prev) => prev.filter((p) => p.id !== projectId))
  }

  const handleProjectStageChange = (projectId: string, newStage: ProjectStage) => {
    setProjects((prev) =>
      prev.map((p) => {
        if (p.id !== projectId) return p
        return { ...p, stage: newStage }
      })
    )
    toast.success(`Project moved to stage ${newStage}`)
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
    setOrders((prev) => {
      const exists = prev.some((o) => o.id === savedOrder.id)
      if (exists) {
        return prev.map((o) => (o.id === savedOrder.id ? savedOrder : o))
      }
      return [savedOrder, ...prev]
    })
    navigate('/orders')
  }

  const handleDeleteOrder = (orderId: string) => {
    setOrders((prev) => prev.filter((o) => o.id !== orderId))
    toast.success('Order deleted successfully')
  }

  const handleOrderStageChange = (orderId: string, newStage: OrderStage) => {
    setOrders((prev) =>
      prev.map((o) => {
        if (o.id !== orderId) return o
        return { ...o, stage: newStage }
      })
    )
    toast.success(`Order stage updated to ${newStage}`)
  }

  // Dynamic Breadcrumb
  const breadcrumbTitle = useMemo(() => {
    if (sidebarMode === 'settings') {
      const titles: Record<SettingsTab, string> = {
        'master-data': 'Settings: Master Data',
        photographers: 'Settings: Photographers & Crew',
        'sequence-id': 'Settings: Sequence ID & Visibility',
        'preset-services': 'Preset Configuration: Services',
        'preset-packages': 'Preset Configuration: Packages',
        'preset-orders': 'Preset Configuration: Orders',
        'preset-deliverables': 'Preset Configuration: Deliverables',
        'template-forms': 'Design: Forms Template',
        'template-quotation-pdf': 'Design: Quotation PDF Template',
        'template-invoice-pdf': 'Design: Invoice PDF Template',
        'template-contracts': 'Design: Contracts Template',
        'template-feedbacks': 'Design: Feedbacks Template',
        'stages-workflow': 'Stages: Production Workflow',
      }
      return titles[activeSetting] || 'Studio Settings'
    }

    if (sidebarMode === 'reports') {
      const titles: Record<ReportCategory, string> = {
        billing: 'Reports: Billing Statement',
        leads: 'Reports: Leads Acquisition Funnel',
        orders: 'Reports: Orders & Merchandise',
        sessions: 'Reports: Sessions & Crew Deployment',
        finance: 'Reports: Financial P&L & Profitability',
        'master-report': 'Reports: Master Operations Register',
        'account-analytics': 'Reports: Client Account Analytics & LTV',
        feedback: 'Reports: Client Feedback & NPS',
      }
      return titles[activeReport] || 'Studio Reports'
    }

    if (pathname === '/leads/create') return 'Create Lead'
    if (pathname.startsWith('/leads/edit')) return `Edit Lead (${editingLead?.name || ''})`
    if (pathname.startsWith('/leads')) return 'Photography Leads'
    if (pathname === '/projects/create') return 'Create Project'
    if (pathname.startsWith('/projects/edit')) return `Project: ${editingProject?.projectName || ''}`
    if (pathname.startsWith('/sessions')) return 'Sessions'
    if (pathname.startsWith('/billings')) return 'Billings'
    if (pathname === '/orders/create') return 'New Order'
    if (pathname.startsWith('/orders/edit')) return `Order: ${editingOrder?.orderNumber || ''}`
    if (pathname.startsWith('/orders')) return 'Orders'
    return 'Projects'
  }, [sidebarMode, activeSetting, activeReport, pathname, editingLead, editingProject, editingOrder])

  return (
    <DashboardShell
      title="Photography CRM"
      breadcrumb={breadcrumbTitle}
      sidebar={
        <PhotographySidebar
          mode={sidebarMode}
          currentRoute={currentMainRoute}
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
          onExitSubmode={() => {
            navigate('/projects')
          }}
        />
      }
      headerActions={
        <div className="flex items-center gap-2">
          <Button
            variant={sidebarMode === 'reports' ? 'default' : 'outline'}
            size="sm"
            onClick={() => {
              if (sidebarMode === 'reports') {
                navigate('/projects')
              } else {
                navigate(`/reports/${activeReport || 'billing'}`)
              }
            }}
            className={`h-7 text-xs font-semibold gap-1.5 transition ${
              sidebarMode === 'reports'
                ? 'bg-emerald-600 hover:bg-emerald-700 text-white shadow-xs'
                : 'border-border text-foreground hover:bg-muted'
            }`}
          >
            <IconReportAnalytics className="h-3.5 w-3.5" />
            <span>Reports</span>
          </Button>

          <Button
            variant={sidebarMode === 'settings' ? 'default' : 'outline'}
            size="sm"
            onClick={() => {
              if (sidebarMode === 'settings') {
                navigate('/projects')
              } else {
                navigate(`/settings/${activeSetting || 'master-data'}`)
              }
            }}
            className={`h-7 text-xs font-semibold gap-1.5 transition ${
              sidebarMode === 'settings'
                ? 'bg-emerald-600 hover:bg-emerald-700 text-white shadow-xs'
                : 'border-border text-foreground hover:bg-muted'
            }`}
          >
            <IconSettings className="h-3.5 w-3.5" />
            <span>Settings</span>
          </Button>
        </div>
      }
    >
      <Toaster position="top-right" richColors />

      {/* Main Module Tabs Navigation Bar (only shown in core CRM mode) */}
      {sidebarMode === 'main' && (
        <div className="flex items-center gap-2 border-b bg-muted/10 px-1 py-1 mb-4 rounded-lg">
          <button
            type="button"
            onClick={() => navigate('/projects')}
            className={`flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-md transition ${
              pathname.startsWith('/projects')
                ? 'bg-background text-foreground shadow-xs font-bold border'
                : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            <IconFolders className="h-3.5 w-3.5" />
            <span>Projects ({projects.length})</span>
          </button>
          <button
            type="button"
            onClick={() => navigate('/sessions')}
            className={`flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-md transition ${
              pathname.startsWith('/sessions')
                ? 'bg-background text-foreground shadow-xs font-bold border'
                : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            <IconCalendar className="h-3.5 w-3.5" />
            <span>Sessions ({sessions.length})</span>
          </button>
          <button
            type="button"
            onClick={() => navigate('/leads')}
            className={`flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-md transition ${
              pathname.startsWith('/leads')
                ? 'bg-background text-foreground shadow-xs font-bold border'
                : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            <IconTarget className="h-3.5 w-3.5" />
            <span>Leads ({leads.length})</span>
          </button>
          <button
            type="button"
            onClick={() => navigate('/billings')}
            className={`flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-md transition ${
              pathname.startsWith('/billings')
                ? 'bg-background text-foreground shadow-xs font-bold border'
                : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            <IconCreditCard className="h-3.5 w-3.5" />
            <span>Billings</span>
          </button>
          <button
            type="button"
            onClick={() => navigate('/orders')}
            className={`flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-md transition ${
              pathname.startsWith('/orders')
                ? 'bg-background text-foreground shadow-xs font-bold border'
                : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            <IconShoppingBag className="h-3.5 w-3.5" />
            <span>Orders ({orders.length})</span>
          </button>
        </div>
      )}

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
                onAddNewSession={() => {}}
                onEditSession={() => {}}
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
          <Route
            path="/settings/preset-orders"
            element={
              <GenericPresetPage
                title="Orders & Merchandise Presets"
                subtitle="Default pricing for albums, frames, and prints"
                icon={IconSettings}
                initialItems={[
                  { id: 'ordp-1', name: '12x36 Photobook (40 Pgs)', rate: 12500, unit: 'album' },
                  { id: 'ordp-2', name: 'Acrylic Wall Frame (24x36)', rate: 6500, unit: 'frame' },
                  { id: 'ordp-3', name: 'Canvas Gallery Wrap (20x30)', rate: 4200, unit: 'canvas' },
                  { id: 'ordp-4', name: 'Premium Wooden Keepsake Box', rate: 3500, unit: 'box' },
                ]}
              />
            }
          />
          <Route
            path="/settings/preset-deliverables"
            element={
              <GenericPresetPage
                title="Deliverables Presets"
                subtitle="Standard deliverables checklist for project milestones"
                icon={IconSettings}
                hasPrice={false}
                initialItems={[
                  { id: 'del-1', name: 'High-Res Edited Photos (Soft Copy)' },
                  { id: 'del-2', name: 'Cinematic 4K Wedding Teaser (3-5 mins)' },
                  { id: 'del-3', name: 'Full Length Documentary Film (45 mins)' },
                  { id: 'del-4', name: 'Custom Engraved Wooden USB Drive' },
                  { id: 'del-5', name: 'Handcrafted Flushmount Coffee Table Album' },
                ]}
              />
            }
          />
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
          <Route
            path="/settings/template-contracts"
            element={
              <TemplateCustomizerPage
                title="Photography Contract Agreement"
                subtitle="Copyright ownership, cancellation clauses, and client liability agreements"
                type="contract"
              />
            }
          />
          <Route
            path="/settings/template-feedbacks"
            element={
              <TemplateCustomizerPage
                title="Client Feedback & NPS Survey"
                subtitle="Post-shoot satisfaction survey questions and review requests"
                type="feedback"
              />
            }
          />
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
      <AppContent />
    </BrowserRouter>
  )
}
