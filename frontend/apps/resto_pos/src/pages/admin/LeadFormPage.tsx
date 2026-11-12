import * as React from 'react'
import {
  IconArrowLeft,
  IconDeviceFloppy,
  IconEdit,
} from "@repo/ui/icons/tabler-icons"
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@repo/ui/components/ui/select'
import { Badge } from '@repo/ui/components/ui/badge'
import { toast } from '@repo/ui/components/ui/sonner'
import type { Lead, ProjectSegment, ServiceItem, LeadSource } from '../../types/lead'
import { EditServiceDrawer } from '../../components/EditServiceDrawer'
import { DeliverablesDropdown } from '../../components/DeliverablesDropdown'
import { DateTimePickerCustom } from '@repo/ui/components/common/DateTimePickerCustom'
import { LeadSourceCombobox } from '../../components/LeadSourceCombobox'
import { PhoneInputField } from '@repo/ui/components/common/PhoneInputField'
import { ProjectSegmentsEditor } from '../../components/ProjectSegmentsEditor'

interface LeadFormPageProps {
  initialLead?: Lead | null
  onBack: () => void
  onSave: (lead: Lead) => void
}

const QUOTATION_TEMPLATES = [
  'Choose Quotation Design Templates',
  'Modern Regal PDF',
  'Minimal Elegant Template',
  'Luxury Gold Palette',
  'Corporate Clean Template',
]

export function LeadFormPage({ initialLead, onBack, onSave }: LeadFormPageProps) {
  const isEditing = Boolean(initialLead)

  // Form State
  const [name, setName] = React.useState(initialLead?.name || '')
  const [countryCode, setCountryCode] = React.useState(initialLead?.countryCode || '+91')
  const [mobileNo, setMobileNo] = React.useState(initialLead?.mobileNo || '')
  const [email, setEmail] = React.useState(initialLead?.email || '')
  const [address, setAddress] = React.useState(initialLead?.address || '')
  const [leadSource, setLeadSource] = React.useState<LeadSource>(
    initialLead?.leadSource || 'Instagram'
  )
  const [followupDate, setFollowupDate] = React.useState(
    initialLead?.followupDate || '2026-11-11T16:28'
  )

  // Project Details
  const [projectName, setProjectName] = React.useState(
    initialLead?.projectName || (initialLead ? `${initialLead.name} Shoot` : '')
  )
  const [projectAddress, setProjectAddress] = React.useState(initialLead?.projectAddress || '')
  const [projectSegments, setProjectSegments] = React.useState<ProjectSegment[]>(
    initialLead?.projectSegments?.length
      ? initialLead.projectSegments
      : [
          {
            id: 'seg-1',
            name: 'Main Shoot',
            type: initialLead?.projectType || 'Wedding',
            startDate: initialLead?.projectStartDate || '2026-12-01',
            endDate: initialLead?.projectEndDate || '2026-12-02',
          },
        ]
  )

  // Package & Services
  const [packageName] = React.useState(
    initialLead?.packageName || 'Sample Wedding Package'
  )
  const [services, setServices] = React.useState<ServiceItem[]>(
    initialLead?.selectedServices?.length
      ? initialLead.selectedServices
      : [
          {
            id: 'srv-base',
            name: 'Full Day Photo & Video Coverage',
            groupName: 'Wedding',
            qty: 1,
            cost: 60000,
            subTotal: 60000,
            sessionPlanned: true,
          },
        ]
  )

  // Quotation Design & Deliverables
  const [template, setTemplate] = React.useState(
    initialLead?.quotationTemplate || 'Modern Regal PDF'
  )
  const [selectedDeliverables, setSelectedDeliverables] = React.useState<string[]>(
    initialLead?.selectedDeliverables || ['Raw Photos', 'Pen Drive']
  )

  // Financials
  const [discountAmount, setDiscountAmount] = React.useState<number>(
    initialLead?.discountAmount || 0
  )
  const [showDiscountInput, setShowDiscountInput] = React.useState(
    Boolean(initialLead?.discountAmount && initialLead.discountAmount > 0)
  )
  const [taxPercent, setTaxPercent] = React.useState<number>(18) // 18% GST standard in photo industry
  const [showTax, setShowTax] = React.useState(true)

  // Drawer
  const [isServiceDrawerOpen, setIsServiceDrawerOpen] = React.useState(false)

  // Cost Calculations
  const servicesSubTotal = React.useMemo(() => {
    return services.reduce((acc, s) => acc + (s.subTotal || 0), 0)
  }, [services])

  const calculatedTax = React.useMemo(() => {
    if (!showTax) return 0
    const taxable = Math.max(0, servicesSubTotal - discountAmount)
    return Math.round((taxable * taxPercent) / 100)
  }, [showTax, servicesSubTotal, discountAmount, taxPercent])

  const grandTotal = React.useMemo(() => {
    const afterDiscount = Math.max(0, servicesSubTotal - discountAmount)
    return afterDiscount + calculatedTax
  }, [servicesSubTotal, discountAmount, calculatedTax])

  // Save Lead
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!name.trim()) {
      toast.error('Lead name is required')
      return
    }
    if (!mobileNo.trim()) {
      toast.error('Mobile number is required')
      return
    }

    const firstSegment = projectSegments[0]

    const leadToSave: Lead = {
      id: initialLead?.id || `LD-${Math.floor(1000 + Math.random() * 9000)}`,
      name: name.trim(),
      mobileNo: mobileNo.trim(),
      countryCode,
      email: email.trim() || 'contact@client.com',
      address,
      createdDate: initialLead?.createdDate || new Date().toISOString().split('T')[0],
      leadSource,
      generate: initialLead?.generate || 'Direct Entry',
      assignee: initialLead?.assignee || { name: 'Rohan Mehta', initials: 'RM' },
      leadPipeline: initialLead?.leadPipeline || 'Wedding Inquiries',
      pipelineStage: initialLead?.pipelineStage || 'new_leads',
      leadStage: initialLead?.leadStage || 'New Lead',
      quoteAmount: grandTotal,
      followupDate,
      projectType: firstSegment?.type || 'Wedding',
      quoteType: initialLead?.quoteType || 'Standard Package',
      projectStartDate: firstSegment?.startDate || new Date().toISOString().split('T')[0],
      projectEndDate: firstSegment?.endDate || new Date().toISOString().split('T')[0],
      projectName: projectName.trim() || `${name} Project`,
      projectAddress,
      packageName,
      projectSegments,
      selectedServices: services,
      selectedDeliverables,
      quotationTemplate: template,
      discountAmount,
      taxAmount: calculatedTax,
    }

    onSave(leadToSave)
    toast.success(isEditing ? 'Lead updated successfully!' : 'New lead created!')
    onBack()
  }

  return (
    <div className="space-y-6 max-w-6xl mx-auto pb-16">
      {/* Top action bar */}
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
            <h1 className="text-xl font-bold tracking-tight text-foreground flex items-center gap-2">
              {isEditing && initialLead ? `Edit Lead: ${initialLead.name}` : 'Create New Lead'}
              {isEditing && initialLead && (
                <Badge variant="outline" className="text-xs font-mono">
                  {initialLead.id}
                </Badge>
              )}
            </h1>
            <p className="text-xs text-muted-foreground">
              Fill in client requirements, schedule segments, and quote breakdown.
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
            {isEditing ? 'Update Lead' : 'Save & Generate Quote'}
          </Button>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* ROW 1: Lead Details + Project Details */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Card 1: Lead Details */}
          <Card className="shadow-xs border bg-card">
            <CardHeader className="pb-4">
              <CardTitle className="text-base font-semibold">Lead Details</CardTitle>
              <CardDescription className="text-xs">
                Contact information and initial channel origin.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Name */}
              <div className="space-y-1.5">
                <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                  Name <span className="text-destructive">*</span>
                </Label>
                <Input
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="e.g. Aarav Sharma"
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
                  placeholder="98765 43210"
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
                  placeholder="client@domain.com"
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
                  placeholder="City, locality or hotel venue"
                  className="h-10 text-sm"
                />
              </div>

              {/* Lead Source with Search */}
              <div className="space-y-1.5">
                <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                  Lead Source <span className="text-destructive">*</span>
                </Label>
                <LeadSourceCombobox
                  value={leadSource}
                  onChange={(val) => setLeadSource(val)}
                  placeholder="Select or search lead source"
                />
              </div>

              {/* Followup Date & Time Picker */}
              <div className="space-y-1.5">
                <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                  Followup Date <span className="text-destructive">*</span>
                </Label>
                <DateTimePickerCustom
                  value={followupDate}
                  onChange={setFollowupDate}
                  placeholder="Select followup date and time"
                />
              </div>
            </CardContent>
          </Card>

          {/* Card 2: Project Details */}
          <Card className="shadow-xs border bg-card flex flex-col">
            <CardHeader className="pb-4">
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle className="text-base font-semibold">Project Details</CardTitle>
                  <CardDescription className="text-xs">
                    Scope of work, venue addresses, and date segments.
                  </CardDescription>
                </div>
              </div>
            </CardHeader>

            <CardContent className="space-y-4 flex-1 flex flex-col">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div className="space-y-1.5">
                  <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                    Project Name
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
                    Project Address
                  </Label>
                  <Input
                    value={projectAddress}
                    onChange={(e) => setProjectAddress(e.target.value)}
                    placeholder="Venue or address"
                    className="h-10 text-sm"
                  />
                </div>
              </div>

              {/* Project Segments Editor */}
              <ProjectSegmentsEditor
                segments={projectSegments}
                onChange={setProjectSegments}
              />
            </CardContent>
          </Card>
        </div>

        {/* ROW 2: Package Details */}
        <Card className="shadow-xs border bg-card">
          <CardHeader className="pb-3 flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-base font-semibold">Package Details</CardTitle>
              <CardDescription className="text-xs">
                Selected packages and pricing inclusions.
              </CardDescription>
            </div>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => setIsServiceDrawerOpen(true)}
              className="text-emerald-600 dark:text-emerald-400 border-emerald-600/30 hover:bg-emerald-50 dark:hover:bg-emerald-950/20 text-xs font-semibold gap-1.5"
            >
              <IconEdit className="h-3.5 w-3.5" />
              Edit Service
            </Button>
          </CardHeader>

          <CardContent>
            <div className="flex flex-wrap gap-4">
              {/* Active package card */}
              <div className="w-56 rounded-xl border border-sky-200 dark:border-sky-800 bg-sky-50/70 dark:bg-sky-950/30 p-4 relative shadow-xs">
                <div className="space-y-1">
                  <p className="text-sm font-semibold text-foreground">{packageName}</p>
                  <p className="text-xl font-bold tracking-tight text-foreground">
                    ₹{servicesSubTotal.toLocaleString('en-IN')}
                  </p>
                </div>
                <div className="mt-3 flex items-center justify-between">
                  <span className="text-[11px] text-muted-foreground">
                    {services.length} services included
                  </span>
                  <div className="h-2 w-2 rounded-full bg-sky-500 animate-pulse" />
                </div>
              </div>

              {/* Extra summary pill of services */}
              <div className="flex-1 min-w-[280px] rounded-xl border bg-muted/20 p-4 flex flex-col justify-between">
                <div className="space-y-2">
                  <span className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                    Configured Services Breakdown
                  </span>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                    {services.map((srv) => (
                      <div
                        key={srv.id}
                        className="flex items-center justify-between text-xs py-1 px-2.5 rounded-md bg-background border border-border/60"
                      >
                        <span className="font-medium truncate max-w-[180px]">{srv.name}</span>
                        <span className="text-muted-foreground tabular-nums">
                          ₹{srv.subTotal.toLocaleString('en-IN')}
                        </span>
                      </div>
                    ))}
                  </div>
                </div>
                <div className="pt-2 flex justify-end">
                  <Button
                    type="button"
                    variant="link"
                    size="sm"
                    onClick={() => setIsServiceDrawerOpen(true)}
                    className="text-xs text-primary p-0 h-auto"
                  >
                    Open service breakdown drawer →
                  </Button>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* ROW 3: Quotation Design & Deliverables + Financial Summary */}
        <Card className="shadow-xs border bg-card">
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-semibold">
              Quotation Design & Deliverables
            </CardTitle>
            <CardDescription className="text-xs">
              PDF styling, deliverable inclusions, and final quote amount.
            </CardDescription>
          </CardHeader>

          <CardContent>
            <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
              {/* Left Column: Template & Deliverables Selector */}
              <div className="lg:col-span-7 space-y-5">
                {/* Template Selection */}
                <div className="space-y-1.5">
                  <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                    Template <span className="text-destructive">*</span>
                  </Label>
                  <Select value={template} onValueChange={setTemplate}>
                    <SelectTrigger className="h-10 text-sm">
                      <SelectValue placeholder="Choose Quotation Design Templates" />
                    </SelectTrigger>
                    <SelectContent>
                      {QUOTATION_TEMPLATES.map((tpl) => (
                        <SelectItem key={tpl} value={tpl} className="text-sm">
                          {tpl}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                {/* Choose Deliverables */}
                <div className="space-y-1.5">
                  <Label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                    Choose Deliverables
                  </Label>
                  <DeliverablesDropdown
                    selectedDeliverables={selectedDeliverables}
                    onChange={setSelectedDeliverables}
                  />
                  <p className="text-[11px] text-muted-foreground">
                    Deliverables selected here will be listed in the client contract and timeline checklist.
                  </p>
                </div>
              </div>

              {/* Right Column: Financial Calculation Block */}
              <div className="lg:col-span-5 rounded-xl border bg-muted/20 p-5 space-y-3.5">
                <div className="flex items-center justify-between text-sm">
                  <span className="text-muted-foreground">Sub Total :</span>
                  <span className="font-semibold tabular-nums text-foreground">
                    ₹{servicesSubTotal.toLocaleString('en-IN')}
                  </span>
                </div>

                {/* Discount */}
                <div className="space-y-1.5">
                  <div className="flex items-center justify-between">
                    <button
                      type="button"
                      onClick={() => setShowDiscountInput(!showDiscountInput)}
                      className="text-xs font-semibold text-emerald-600 dark:text-emerald-400 hover:underline flex items-center gap-1"
                    >
                      {showDiscountInput ? 'Remove Discount' : '+ Add Discount'}
                    </button>
                    {showDiscountInput && discountAmount > 0 && (
                      <span className="text-xs text-destructive tabular-nums">
                        - ₹{discountAmount.toLocaleString('en-IN')}
                      </span>
                    )}
                  </div>
                  {showDiscountInput && (
                    <div className="flex items-center gap-2 pt-1">
                      <span className="text-xs text-muted-foreground">₹</span>
                      <Input
                        type="number"
                        min="0"
                        step="500"
                        value={discountAmount}
                        onChange={(e) => setDiscountAmount(parseFloat(e.target.value) || 0)}
                        className="h-8 text-xs tabular-nums"
                        placeholder="Discount amount"
                      />
                    </div>
                  )}
                </div>

                {/* Tax */}
                <div className="space-y-1.5">
                  <div className="flex items-center justify-between">
                    <button
                      type="button"
                      onClick={() => setShowTax(!showTax)}
                      className="text-xs font-semibold text-emerald-600 dark:text-emerald-400 hover:underline"
                    >
                      {showTax ? 'Remove GST / Tax' : '+ Add Tax (GST)'}
                    </button>
                    {showTax && (
                      <span className="text-xs text-muted-foreground tabular-nums">
                        ₹{calculatedTax.toLocaleString('en-IN')} ({taxPercent}%)
                      </span>
                    )}
                  </div>
                  {showTax && (
                    <div className="flex items-center gap-2 pt-1">
                      <span className="text-xs text-muted-foreground">GST %:</span>
                      <Input
                        type="number"
                        min="0"
                        max="30"
                        value={taxPercent}
                        onChange={(e) => setTaxPercent(parseInt(e.target.value) || 0)}
                        className="h-8 w-20 text-xs tabular-nums"
                      />
                    </div>
                  )}
                </div>

                <div className="border-t pt-3 flex items-center justify-between">
                  <span className="text-base font-bold text-foreground">Total</span>
                  <span className="text-2xl font-bold tracking-tight text-foreground">
                    ₹{grandTotal.toLocaleString('en-IN')}
                  </span>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </form>

      {/* Edit Lead Service Drawer */}
      <EditServiceDrawer
        open={isServiceDrawerOpen}
        onOpenChange={setIsServiceDrawerOpen}
        services={services}
        onSave={(updatedServices: ServiceItem[]) => {
          setServices(updatedServices)
          toast.success('Services updated!')
        }}
      />
    </div>
  )
}
