import * as React from 'react'
import {
  IconDeviceFloppy,
  IconBold,
  IconItalic,
  IconAlignLeft,
  IconAlignCenter,
  IconAlignRight,
  IconList,
  IconListNumbers,
  IconSeparator,
  IconFileText,
  IconSignature,
  IconVariable,
} from "@repo/ui/icons/tabler-icons"
import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import { Label } from '@repo/ui/components/ui/label'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@repo/ui/components/ui/popover'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@repo/ui/components/ui/select'
import {
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from '@repo/ui/components/ui/tooltip'
import { toast } from '@repo/ui/components/ui/sonner'
import { cn } from '@repo/ui/lib/utils'

interface MergeVariable {
  key: string
  label: string
  defaultValue: string
}

const AVAILABLE_VARIABLES: MergeVariable[] = [
  { key: 'client_name', label: 'Client Name', defaultValue: 'Priya & Rohan' },
  { key: 'client_mobile', label: 'Client Mobile', defaultValue: '+91 98765 43210' },
  { key: 'client_email', label: 'Client Email', defaultValue: 'priya.rohan@example.com' },
  { key: 'event_name', label: 'Event Name', defaultValue: 'Grand Royal Wedding & Reception' },
  { key: 'event_type', label: 'Event Type', defaultValue: 'Wedding' },
  { key: 'event_date', label: 'Event Date', defaultValue: '24 Dec 2026' },
  { key: 'event_address', label: 'Event Address', defaultValue: 'Taj Palace, Mumbai' },
  { key: 'studio_name', label: 'Studio Name', defaultValue: 'especa Studios' },
  { key: 'studio_support_no', label: 'Studio Support No.', defaultValue: '+91 91234 56789' },
  { key: 'studio_support_email', label: 'Studio Support Mail', defaultValue: 'hello@especa.com' },
  { key: 'total_amount', label: 'Total Contract Value', defaultValue: '₹2,50,000' },
  { key: 'advance_amount', label: 'Advance Deposit', defaultValue: '₹1,00,000' },
]

export function ContractTemplatePage() {
  const [contractName, setContractName] = React.useState('Birthday Contract')
  const [mySignatureRequired, setMySignatureRequired] = React.useState(true)

  // Document custom variable values
  const [variableValues, setVariableValues] = React.useState<Record<string, string>>(() => {
    const initial: Record<string, string> = {}
    AVAILABLE_VARIABLES.forEach((v) => {
      initial[v.key] = v.defaultValue
    })
    return initial
  })

  const [activeEditingKey, setActiveEditingKey] = React.useState<string | null>(null)
  const [tempValue, setTempValue] = React.useState('')

  const handleOpenEdit = (vKey: string) => {
    setActiveEditingKey(vKey)
    setTempValue(variableValues[vKey] || '')
  }

  const handleConfirmEdit = (vKey: string) => {
    if (tempValue.trim()) {
      setVariableValues((prev) => ({ ...prev, [vKey]: tempValue.trim() }))
    }
    setActiveEditingKey(null)
    toast.success('Field updated in preview')
  }

  const handleSaveContract = () => {
    toast.success(`Contract template "${contractName}" saved successfully`)
  }

  const renderDynamicPill = (vKey: string, fallbackLabel: string) => {
    const currentVal = variableValues[vKey] || fallbackLabel
    const isOpen = activeEditingKey === vKey

    return (
      <Popover
        open={isOpen}
        onOpenChange={(open) => {
          if (open) handleOpenEdit(vKey)
          else setActiveEditingKey(null)
        }}
      >
        <PopoverTrigger asChild>
          <button
            type="button"
            className="inline-flex items-center px-1.5 py-0.5 mx-0.5 text-xs font-medium text-emerald-800 bg-emerald-100 hover:bg-emerald-200 dark:text-emerald-300 dark:bg-emerald-950/60 rounded border border-emerald-300 dark:border-emerald-700/60 transition cursor-pointer"
          >
            {currentVal}
          </button>
        </PopoverTrigger>
        <PopoverContent className="w-64 p-3 space-y-2.5 shadow-lg" align="start">
          <div className="space-y-1">
            <div className="text-[11px] font-bold uppercase tracking-wider text-muted-foreground">
              {fallbackLabel}
            </div>
            <Input
              value={tempValue}
              onChange={(e) => setTempValue(e.target.value)}
              placeholder="Enter value"
              className="h-8 text-xs"
              autoFocus
              onKeyDown={(e) => {
                if (e.key === 'Enter') {
                  e.preventDefault()
                  handleConfirmEdit(vKey)
                }
              }}
            />
          </div>
          <div className="flex justify-end gap-1.5">
            <Button
              type="button"
              size="sm"
              onClick={() => handleConfirmEdit(vKey)}
              className="h-7 text-xs px-3 bg-emerald-600 hover:bg-emerald-700 text-white font-semibold"
            >
              Confirm
            </Button>
          </div>
        </PopoverContent>
      </Popover>
    )
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6 pb-16">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b pb-4">
        <div className="flex items-center gap-3">
          <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary border border-primary/20">
            <IconSignature className="h-5 w-5" />
          </span>
          <div>
            <h1 className="text-xl font-bold tracking-tight text-foreground">Contracts Template</h1>
            <p className="text-xs text-muted-foreground">
              Define standard legal terms, merge fields, and electronic signature requirements.
            </p>
          </div>
        </div>

        <Button
          type="button"
          onClick={handleSaveContract}
          className="h-8 text-xs font-semibold gap-1.5 bg-primary text-primary-foreground self-start sm:self-auto shadow-xs"
        >
          <IconDeviceFloppy className="h-3.5 w-3.5" />
          Save Contract
        </Button>
      </div>

      {/* Top Header Controls (Matches Screenshot 3) */}
      <div className="bg-card border border-border/80 shadow-2xs rounded-xl p-4 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div className="space-y-1.5 w-full sm:w-72">
          <Label className="text-xs font-semibold text-foreground">Contract Name</Label>
          <Input
            value={contractName}
            onChange={(e) => setContractName(e.target.value)}
            className="h-9 text-xs bg-background font-medium"
            placeholder="e.g. Wedding Contract Template"
          />
        </div>

        <div className="flex items-center gap-3">
          <div className="text-right">
            <div className="text-xs font-semibold text-foreground">My Signature Required</div>
            <div className="text-[10px] text-muted-foreground">
              Set whether your signature is required to complete the contract.
            </div>
          </div>
          <button
            type="button"
            role="switch"
            aria-checked={mySignatureRequired}
            onClick={() => setMySignatureRequired(!mySignatureRequired)}
            className={cn(
              'relative inline-flex h-5 w-9 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none',
              mySignatureRequired ? 'bg-primary' : 'bg-muted'
            )}
          >
            <span
              className={cn(
                'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                mySignatureRequired ? 'translate-x-4' : 'translate-x-0'
              )}
            />
          </button>
        </div>
      </div>

      {/* Contract Editor Container */}
      <div className="bg-card border border-border/80 shadow-xs rounded-xl overflow-hidden">
        {/* Rich Text Toolbar (Matches Screenshot 3) */}
        <div className="border-b bg-muted/30 p-2 flex flex-wrap items-center justify-between gap-2">
          <div className="flex flex-wrap items-center gap-1">
            <Select defaultValue="p">
              <SelectTrigger className="h-8 w-28 text-xs bg-background">
                <SelectValue placeholder="Format" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="p">Paragraph</SelectItem>
                <SelectItem value="h1">Heading 1</SelectItem>
                <SelectItem value="h2">Heading 2</SelectItem>
                <SelectItem value="h3">Heading 3</SelectItem>
              </SelectContent>
            </Select>

            <div className="h-4 w-px bg-border mx-1" />

            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground">
                  <IconBold className="h-4 w-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent side="top">Bold (Ctrl+B)</TooltipContent>
            </Tooltip>

            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground">
                  <IconItalic className="h-4 w-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent side="top">Italic (Ctrl+I)</TooltipContent>
            </Tooltip>

            <div className="h-4 w-px bg-border mx-1" />

            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground">
                  <IconAlignLeft className="h-4 w-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent side="top">Align Left</TooltipContent>
            </Tooltip>

            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground">
                  <IconAlignCenter className="h-4 w-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent side="top">Align Center</TooltipContent>
            </Tooltip>

            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground">
                  <IconAlignRight className="h-4 w-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent side="top">Align Right</TooltipContent>
            </Tooltip>

            <div className="h-4 w-px bg-border mx-1" />

            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground">
                  <IconList className="h-4 w-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent side="top">Bullet List</TooltipContent>
            </Tooltip>

            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground">
                  <IconListNumbers className="h-4 w-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent side="top">Numbered List</TooltipContent>
            </Tooltip>

            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground">
                  <IconSeparator className="h-4 w-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent side="top">Divider Line</TooltipContent>
            </Tooltip>
          </div>

          <div className="flex items-center gap-2">
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => toast.info('Click any green highlighted field in the document below to edit it')}
              className="h-8 text-xs font-semibold gap-1 bg-background"
            >
              <IconFileText className="h-3.5 w-3.5 text-primary" />
              Insert Field
            </Button>

            <Select
              onValueChange={(val) => {
                const found = AVAILABLE_VARIABLES.find((v) => v.key === val)
                if (found) {
                  toast.success(`Variable "${found.label}" inserted into preview`)
                }
              }}
            >
              <SelectTrigger className="h-8 text-xs w-36 bg-background">
                <span className="flex items-center gap-1.5">
                  <IconVariable className="h-3.5 w-3.5 text-primary" />
                  Insert Variable
                </span>
              </SelectTrigger>
              <SelectContent>
                {AVAILABLE_VARIABLES.map((v) => (
                  <SelectItem key={v.key} value={v.key}>
                    {v.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        {/* Contract Interactive Preview Canvas (Matches Screenshot 3) */}
        <div className="p-8 sm:p-12 bg-background min-h-[500px] text-xs leading-relaxed space-y-6 max-w-3xl mx-auto">
          {/* Header */}
          <div className="text-center pb-4 border-b space-y-1">
            <h2 className="text-base font-bold text-foreground uppercase tracking-wider">
              {contractName || 'Standard Photography Services Agreement'}
            </h2>
            <p className="text-[11px] text-muted-foreground">
              This contract is entered between {renderDynamicPill('studio_name', 'Studio Name')} and the Client identified below.
            </p>
          </div>

          {/* Client Information */}
          <div className="space-y-2">
            <h3 className="font-bold text-foreground text-xs uppercase tracking-wide">
              Client Information:
            </h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 text-muted-foreground">
              <div>
                Client Name: {renderDynamicPill('client_name', 'Client Name')}
              </div>
              <div>
                Client Mobile: {renderDynamicPill('client_mobile', 'Client Mobile')}
              </div>
              <div className="sm:col-span-2">
                Client Email: {renderDynamicPill('client_email', 'Client Email')}
              </div>
            </div>
          </div>

          {/* Event Details */}
          <div className="space-y-2">
            <h3 className="font-bold text-foreground text-xs uppercase tracking-wide">
              Event Details:
            </h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 text-muted-foreground">
              <div>
                Event Name: {renderDynamicPill('event_name', 'Event Name')}
              </div>
              <div>
                Event Type: {renderDynamicPill('event_type', 'Event Type')}
              </div>
              <div>
                Event Date: {renderDynamicPill('event_date', 'Event Date')}
              </div>
              <div>
                Event Venue: {renderDynamicPill('event_address', 'Event Address')}
              </div>
            </div>
          </div>

          {/* Studio Information */}
          <div className="space-y-2">
            <h3 className="font-bold text-foreground text-xs uppercase tracking-wide">
              Studio Information:
            </h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 text-muted-foreground">
              <div>
                Studio Name: {renderDynamicPill('studio_name', 'Studio Name')}
              </div>
              <div>
                Support Contact: {renderDynamicPill('studio_support_no', 'Studio Support No.')}
              </div>
              <div className="sm:col-span-2">
                Support Email: {renderDynamicPill('studio_support_email', 'Studio Support Mail')}
              </div>
            </div>
          </div>

          {/* Photography Services & Terms */}
          <div className="space-y-2">
            <h3 className="font-bold text-foreground text-xs uppercase tracking-wide">
              Photography Services & Scope:
            </h3>
            <p className="text-muted-foreground leading-normal">
              1. The Studio agrees to provide professional photography and cinematography coverage for the aforementioned event.<br />
              2. Contract Total is {renderDynamicPill('total_amount', 'Total Contract Value')}. A non-refundable advance deposit of {renderDynamicPill('advance_amount', 'Advance Deposit')} is due upon execution.<br />
              3. Color-graded master photos will be delivered via digital gallery within 21 business days following event conclusion.<br />
              4. Studio retains copyright of all original raw frames and may use selected non-private shots for artistic portfolio promotion.
            </p>
          </div>

          {/* Signature Block */}
          <div className="pt-8 border-t grid grid-cols-2 gap-8 text-center">
            <div className="space-y-2">
              <div className="h-14 border-b border-dashed border-border/80 flex items-end justify-center pb-1">
                {mySignatureRequired && (
                  <span className="font-serif italic text-sm text-foreground/80">especa Studios Authorized</span>
                )}
              </div>
              <div className="text-[11px] font-semibold text-muted-foreground">
                Studio Representative
              </div>
            </div>
            <div className="space-y-2">
              <div className="h-14 border-b border-dashed border-border/80 flex items-end justify-center pb-1">
                <span className="text-[11px] text-muted-foreground italic">Client Digital Signature</span>
              </div>
              <div className="text-[11px] font-semibold text-muted-foreground">
                Client Signatory
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

