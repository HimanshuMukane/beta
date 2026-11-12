import * as React from 'react'
import {
  IconPlus,
  IconTrash,
  IconPalette,
  IconCheck,
} from "@repo/ui/icons/tabler-icons"
import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@repo/ui/components/ui/select'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@repo/ui/components/ui/popover'
import { toast } from '@repo/ui/components/ui/sonner'
import { useWorkflows, type WorkflowItem } from '../../hooks/useWorkflows'

const COLOR_PALETTE = [
  '#94a3b8', // slate/gray
  '#e11d48', // rose/red
  '#f59e0b', // amber/yellow
  '#10b981', // emerald/green
  '#3b82f6', // blue
  '#6366f1', // indigo
  '#8b5cf6', // purple
  '#ec4899', // pink
  '#06b6d4', // cyan
  '#14b8a6', // teal
]

const PIPELINES_BY_TYPE: Record<string, string[]> = {
  lead: ['Contacted', 'Qualified', 'Converted', 'Followups'],
  project: ['Production'],
  order: ['Fulfillment'],
}

export function WorkflowPage() {
  const [activeWorkflowType, setActiveWorkflowType] = React.useState<'lead' | 'project' | 'order'>('lead')
  const [activeTab, setActiveTab] = React.useState<'editing' | 'reordering'>('editing')

  const { workflows, isLoading, saveBatchAsync, isSaving } = useWorkflows(activeWorkflowType)
  const [stages, setStages] = React.useState<WorkflowItem[]>([])

  React.useEffect(() => {
    if (workflows) {
      setStages(workflows)
    }
  }, [workflows])

  // Save batch to backend
  const handleSaveAll = async () => {
    try {
      await saveBatchAsync(
        stages.map((s, idx) => ({
          ...s,
          order_index: idx + 1,
        }))
      )
      toast.success('Workflow stages saved successfully!')
    } catch (err: any) {
      toast.error(err?.message || 'Failed to save workflow stages')
    }
  }

  // Add new stage
  const handleAddStage = () => {
    const pipelines = PIPELINES_BY_TYPE[activeWorkflowType] || ['Default']
    const newStage: WorkflowItem = {
      id: `temp-${Date.now()}`,
      workflow_type: activeWorkflowType,
      pipeline: pipelines[0] || 'Contacted',
      stage: '',
      color_code: COLOR_PALETTE[Math.floor(Math.random() * COLOR_PALETTE.length)] || '#10b981',
      order_index: stages.length + 1,
    }
    setStages((prev) => [...prev, newStage])
  }

  // Update stage field
  const handleUpdateStage = (id: string, field: keyof WorkflowItem, value: any) => {
    setStages((prev) =>
      prev.map((s) => (s.id === id ? { ...s, [field]: value } : s))
    )
  }

  // Delete stage
  const handleDeleteStage = (id: string) => {
    if (stages.length <= 1) {
      toast.error('Workflow requires at least one stage')
      return
    }
    setStages((prev) => prev.filter((s) => s.id !== id))
  }

  // Move stage up/down in re-ordering mode
  const handleMoveStage = (index: number, direction: 'up' | 'down') => {
    const targetIndex = direction === 'up' ? index - 1 : index + 1
    if (targetIndex < 0 || targetIndex >= stages.length) return
    const updated = [...stages]
    const current = updated[index]
    const target = updated[targetIndex]
    if (!current || !target) return
    updated[index] = target
    updated[targetIndex] = current
    setStages(updated)
  }

  const pipelines = PIPELINES_BY_TYPE[activeWorkflowType] || ['Contacted']

  return (
    <div className="space-y-6 max-w-2xl mx-auto pb-20 pt-2">
      {/* Workflow Category Tabs (Lead / Project / Order) */}
      <div className="flex items-center gap-2 border-b border-border/80 pb-3">
        {[
          { id: 'lead', label: 'Lead Stages' },
          { id: 'project', label: 'Project Stages' },
          { id: 'order', label: 'Order Stages' },
        ].map((tab) => (
          <button
            key={tab.id}
            type="button"
            onClick={() => setActiveWorkflowType(tab.id as any)}
            className={`px-4 py-1.5 rounded-lg text-xs font-semibold transition ${
              activeWorkflowType === tab.id
                ? 'bg-foreground text-background shadow-xs'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted/60'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Header Row: Title + Save Button */}
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold tracking-tight text-foreground">
          {activeWorkflowType === 'lead'
            ? 'Leads Stages'
            : activeWorkflowType === 'project'
            ? 'Project Stages'
            : 'Order Stages'}
        </h1>

        <Button
          onClick={handleSaveAll}
          disabled={isSaving || isLoading}
          className="bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs px-5 h-9 rounded-lg shadow-xs"
        >
          {isSaving ? 'Saving...' : 'Save'}
        </Button>
      </div>

      {/* Mode Switcher Pill: Editing / Re-Ordering */}
      <div className="flex justify-center">
        <div className="inline-flex rounded-lg bg-muted/60 p-1 border border-border/60">
          <button
            type="button"
            onClick={() => setActiveTab('editing')}
            className={`px-6 py-1.5 rounded-md text-xs font-semibold transition ${
              activeTab === 'editing'
                ? 'bg-card text-foreground shadow-xs'
                : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            Editing
          </button>
          <button
            type="button"
            onClick={() => setActiveTab('reordering')}
            className={`px-6 py-1.5 rounded-md text-xs font-semibold transition ${
              activeTab === 'reordering'
                ? 'bg-card text-foreground shadow-xs'
                : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            Re-Ordering
          </button>
        </div>
      </div>

      {/* Table Header Labels */}
      <div className="grid grid-cols-12 gap-3 px-1 text-xs font-semibold text-foreground/80">
        <div className="col-span-4">Pipeline</div>
        <div className="col-span-8">Stage</div>
      </div>

      {/* List of Stages */}
      {isLoading ? (
        <div className="p-8 text-center text-xs text-muted-foreground">Loading stages...</div>
      ) : (
        <div className="space-y-2.5">
          {stages.map((stg, idx) => (
            <div
              key={stg.id}
              className="grid grid-cols-12 gap-3 items-center bg-card p-1.5 rounded-xl border border-border/60 shadow-2xs hover:border-border transition"
            >
              {/* Reorder Handle in Re-ordering mode */}
              {activeTab === 'reordering' && (
                <div className="col-span-1 flex items-center justify-center gap-1 text-muted-foreground">
                  <button
                    type="button"
                    disabled={idx === 0}
                    onClick={() => handleMoveStage(idx, 'up')}
                    className="hover:text-foreground disabled:opacity-30 p-1 text-xs font-bold"
                    title="Move Up"
                  >
                    ▲
                  </button>
                  <button
                    type="button"
                    disabled={idx === stages.length - 1}
                    onClick={() => handleMoveStage(idx, 'down')}
                    className="hover:text-foreground disabled:opacity-30 p-1 text-xs font-bold"
                    title="Move Down"
                  >
                    ▼
                  </button>
                </div>
              )}

              {/* Pipeline Dropdown */}
              <div className={activeTab === 'reordering' ? 'col-span-4' : 'col-span-4'}>
                <Select
                  value={stg.pipeline}
                  onValueChange={(val) => handleUpdateStage(stg.id, 'pipeline', val)}
                >
                  <SelectTrigger className="h-10 text-xs bg-card border-border/70 rounded-lg">
                    <SelectValue placeholder="Pipeline" />
                  </SelectTrigger>
                  <SelectContent className="text-xs">
                    {pipelines.map((p) => (
                      <SelectItem key={p} value={p}>
                        {p}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              {/* Stage Name Input */}
              <div className={activeTab === 'reordering' ? 'col-span-5' : 'col-span-6'}>
                <Input
                  value={stg.stage}
                  onChange={(e) => handleUpdateStage(stg.id, 'stage', e.target.value)}
                  placeholder="Stage Name"
                  className="h-10 text-xs bg-muted/20 border-border/70 rounded-lg font-medium text-foreground"
                />
              </div>

              {/* Color Palette Picker */}
              <div className="col-span-1 flex items-center justify-center">
                <Popover>
                  <PopoverTrigger asChild>
                    <button
                      type="button"
                      style={{ backgroundColor: stg.color_code || '#94a3b8' }}
                      className="h-8 w-8 rounded-full flex items-center justify-center text-white shadow-2xs hover:scale-105 transition"
                      title={`Color: ${stg.color_code}`}
                    >
                      <IconPalette className="h-4 w-4 drop-shadow-xs" />
                    </button>
                  </PopoverTrigger>
                  <PopoverContent align="center" className="w-48 p-2.5">
                    <div className="text-[11px] font-semibold text-muted-foreground mb-2">
                      Select Stage Color
                    </div>
                    <div className="grid grid-cols-5 gap-2">
                      {COLOR_PALETTE.map((hex) => (
                        <button
                          key={hex}
                          type="button"
                          style={{ backgroundColor: hex }}
                          onClick={() => handleUpdateStage(stg.id, 'color_code', hex)}
                          className="h-6 w-6 rounded-full flex items-center justify-center text-white shadow-xs hover:scale-110 transition"
                        >
                          {stg.color_code === hex && <IconCheck className="h-3.5 w-3.5" />}
                        </button>
                      ))}
                    </div>
                  </PopoverContent>
                </Popover>
              </div>

              {/* Delete Button */}
              <div className="col-span-1 flex items-center justify-center">
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => handleDeleteStage(stg.id)}
                  className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive rounded-full"
                  title="Delete stage"
                >
                  <IconTrash className="h-4 w-4" />
                </Button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Add Stage Action */}
      <div className="pt-2">
        <button
          type="button"
          onClick={handleAddStage}
          className="inline-flex items-center gap-1.5 text-xs font-semibold text-emerald-600 dark:text-emerald-400 hover:text-emerald-700 transition"
        >
          <IconPlus className="h-4 w-4" />
          <span>Add Stage</span>
        </button>
      </div>
    </div>
  )
}
