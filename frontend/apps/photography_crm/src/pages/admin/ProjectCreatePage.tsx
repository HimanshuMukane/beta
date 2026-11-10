import * as React from 'react'
import {
  IconArrowLeft,
  IconDeviceFloppy,
} from '@tabler/icons-react'
import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import { Label } from '@repo/ui/components/ui/label'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@repo/ui/components/ui/card'
import { toast } from '@repo/ui/components/ui/sonner'
import type { Project, ProjectSegment, ProjectType } from '../../types/project'
import { LeadSourceCombobox } from '../../components/LeadSourceCombobox'
import { PhoneInputField } from '../../components/PhoneInputField'
import { ProjectSegmentsEditor } from '../../components/ProjectSegmentsEditor'

interface ProjectCreatePageProps {
  initialData?: Partial<Project>
  onBack: () => void
  onSave: (project: Project) => void
}

export function ProjectCreatePage({ initialData, onBack, onSave }: ProjectCreatePageProps) {
  // Client Details Form State (prefill if provided from converted lead)
  const [clientName, setClientName] = React.useState(initialData?.clientName || '')
  const [countryCode, setCountryCode] = React.useState('+91')
  const [mobileNo, setMobileNo] = React.useState(initialData?.clientMobile || '')
  const [email, setEmail] = React.useState(initialData?.clientEmail || '')
  const [address, setAddress] = React.useState(initialData?.projectAddress || '')
  const [leadSource, setLeadSource] = React.useState('Instagram')

  // Projects Details Form State
  const [projectName, setProjectName] = React.useState(initialData?.projectName || '')
  const [projectAddress, setProjectAddress] = React.useState(initialData?.projectAddress || '')
  const [segments, setSegments] = React.useState<ProjectSegment[]>(
    initialData?.segments && initialData.segments.length > 0
      ? initialData.segments
      : [
          {
            id: 'seg-1',
            name: 'Main Shoot',
            type: (initialData?.projectType as ProjectType) || 'Wedding',
            startDate: initialData?.startDate || new Date().toISOString().split('T')[0],
            endDate: initialData?.endDate || new Date().toISOString().split('T')[0],
          },
        ]
  )

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!clientName.trim()) {
      toast.error('Client name is required')
      return
    }
    if (!mobileNo.trim()) {
      toast.error('Mobile number is required')
      return
    }

    const firstSegment = segments[0]
    const initials = clientName
      .split(' ')
      .map((n) => n[0])
      .slice(0, 2)
      .join('')
      .toUpperCase() || 'PR'

    const newProject: Project = {
      id: `PRJ-${Math.floor(100 + Math.random() * 900)}`,
      projectName: projectName.trim() || `${clientName.trim()} Wedding`,
      clientName: clientName.trim(),
      clientInitials: initials,
      clientMobile: mobileNo.trim(),
      clientCountryCode: countryCode,
      clientEmail: email.trim() || 'contact@client.com',
      clientAddress: address,
      leadSource,
      projectManager: { name: 'Rohan Mehta', initials: 'RM' },
      projectType: firstSegment?.type || 'Wedding',
      projectAddress,
      sessionsCount: segments.length,
      tasksCount: { completed: 0, total: 3 },
      deliverablesCount: { ready: 0, total: 2 },
      status: 'New',
      stage: 'Booked',
      startDate: `${firstSegment?.startDate || '2026-11-10'}T16:00`,
      endDate: `${firstSegment?.endDate || '2026-11-12'}T18:00`,
      createdDate: new Date().toISOString().split('T')[0],
      shootHours: 8.0,
      shootMembersCount: 2,
      shootExpense: 0,
      galleriesCount: 0,
      finance: {
        profit: 0,
        currentProfit: 0,
        netPayableAmount: 0,
        paidAmount: 0,
        balanceAmount: 0,
        paymentDue: 0,
        expense: 0,
        eventDiscountAmount: 0,
        paymentHistory: [],
      },
      segments,
      documents: [
        {
          id: `doc-${Date.now()}`,
          type: 'Lead',
          title: clientName,
          referenceNo: `ESPE/PRJ/${Math.floor(100 + Math.random() * 900)}`,
          status: 'Converted',
          date: new Date().toLocaleDateString('en-GB'),
        },
      ],
      sessions: segments.map((seg, idx) => ({
        id: `ses-${Date.now()}-${idx}`,
        title: seg.name,
        date: seg.startDate,
        timeRange: '04:00 PM - 09:00 PM',
        venue: projectAddress || 'Main Venue',
        city: 'Local',
        status: 'New',
        members: ['Rohan Mehta'],
        shootsActivityCount: 1,
      })),
      tasks: [
        { id: `tsk-1`, title: 'Call client to confirm schedule', assignee: 'Rohan Mehta', dueDate: segDueDate(0), completed: false, priority: 'high' },
        { id: `tsk-2`, title: 'Prepare gear and camera sensors', assignee: 'Rohan Mehta', dueDate: segDueDate(1), completed: false, priority: 'medium' },
      ],
      printing: [],
      deliverables: [
        { id: `del-1`, deliverable: 'Raw Photos', dueDate: segDueDate(0), assignee: 'Rohan Mehta', status: 'New', note: '' },
        { id: `del-2`, deliverable: 'Album', dueDate: segDueDate(1), assignee: 'Rohan Mehta', status: 'New', note: '' },
      ],
      expenses: [],
      notes: `<p>Project created on ${new Date().toLocaleDateString()}. Initial setup ready for shoot planning.</p>`,
    }

    onSave(newProject)
    toast.success('Project created successfully!')
    onBack()
  }

  function segDueDate(offsetDays: number) {
    const d = new Date()
    d.setDate(d.getDate() + offsetDays)
    return d.toISOString().split('T')[0]
  }

  return (
    <div className="space-y-6 max-w-6xl mx-auto pb-16">
      {/* Header bar */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 pb-2 border-b">
        <div className="flex items-center gap-3">
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={onBack}
            className="h-8 w-8 p-0"
          >
            <IconArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h1 className="text-xl font-bold tracking-tight text-foreground">
              Create New Project
            </h1>
            <p className="text-xs text-muted-foreground">
              Define client information, venue, and project event schedule segments.
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2.5">
          <Button type="button" variant="outline" onClick={onBack}>
            Cancel
          </Button>
          <Button
            type="button"
            onClick={handleSubmit}
            className="bg-emerald-600 hover:bg-emerald-700 text-white gap-1.5"
          >
            <IconDeviceFloppy className="h-4 w-4" />
            Create Project
          </Button>
        </div>
      </div>

      <form onSubmit={handleSubmit}>
        {/* ROW 1: Client Details + Projects Details (matches uploaded image) */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 items-start">
          {/* Card 1: Client Details */}
          <Card className="shadow-xs border bg-card">
            <CardHeader className="pb-4">
              <CardTitle className="text-base font-semibold">Client Details</CardTitle>
              <CardDescription className="text-xs">
                Contact information and primary communications origin.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Name */}
              <div className="space-y-1.5">
                <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                  Name <span className="text-destructive">*</span>
                </Label>
                <Input
                  value={clientName}
                  onChange={(e) => setClientName(e.target.value)}
                  placeholder="e.g. Kavya Rajendran"
                  className="h-10 text-sm"
                  required
                />
              </div>

              {/* Mobile Number */}
              <div className="space-y-1.5">
                <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                  Mobile No <span className="text-destructive">*</span>
                </Label>
                <PhoneInputField
                  countryCode={countryCode}
                  mobileNo={mobileNo}
                  onCountryCodeChange={setCountryCode}
                  onMobileNoChange={setMobileNo}
                  placeholder="Mobile No"
                  required
                />
              </div>

              {/* Email Address */}
              <div className="space-y-1.5">
                <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                  Email Address
                </Label>
                <Input
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="Email Address"
                  type="email"
                  className="h-10 text-sm"
                />
              </div>

              {/* Address */}
              <div className="space-y-1.5">
                <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                  Address
                </Label>
                <Input
                  value={address}
                  onChange={(e) => setAddress(e.target.value)}
                  placeholder="Address"
                  className="h-10 text-sm"
                />
              </div>

              {/* Lead Source */}
              <div className="space-y-1.5">
                <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                  Lead Source <span className="text-destructive">*</span>
                </Label>
                <LeadSourceCombobox
                  value={leadSource}
                  onChange={setLeadSource}
                  placeholder="Lead Source"
                />
              </div>
            </CardContent>
          </Card>

          {/* Card 2: Projects Details */}
          <Card className="shadow-xs border bg-card flex flex-col">
            <CardHeader className="pb-4">
              <CardTitle className="text-base font-semibold">Projects Details</CardTitle>
              <CardDescription className="text-xs">
                Title, venue, and date block segments.
              </CardDescription>
            </CardHeader>

            <CardContent className="space-y-4 flex-1 flex flex-col">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div className="space-y-1.5">
                  <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                    Projects Name
                  </Label>
                  <Input
                    value={projectName}
                    onChange={(e) => setProjectName(e.target.value)}
                    placeholder="e.g. Johnson Wedding"
                    className="h-10 text-sm"
                  />
                </div>
                <div className="space-y-1.5">
                  <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                    Projects Address
                  </Label>
                  <Input
                    value={projectAddress}
                    onChange={(e) => setProjectAddress(e.target.value)}
                    placeholder="Venue or address"
                    className="h-10 text-sm"
                  />
                </div>
              </div>

              {/* Projects Segments Editor */}
              <ProjectSegmentsEditor
                segments={segments}
                onChange={setSegments}
              />
            </CardContent>
          </Card>
        </div>
      </form>
    </div>
  )
}
