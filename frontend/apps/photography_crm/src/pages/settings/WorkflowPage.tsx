import * as React from 'react'
import {
  IconPlus,
  IconTrash,
} from '@tabler/icons-react'
import {
  Button,
  Input,
  toast,
} from '@repo/ui'

interface WorkflowStageItem {
  id: string
  name: string
  color: string
  description: string
}

export function WorkflowPage() {
  const [stages, setStages] = React.useState<WorkflowStageItem[]>([
    {
      id: 'stg-1',
      name: 'Booked',
      color: 'bg-blue-500/15 text-blue-700 dark:text-blue-300 border-blue-500/30',
      description: 'Contract executed & initial token advance received.',
    },
    {
      id: 'stg-2',
      name: 'Shoot Planning',
      color: 'bg-amber-500/15 text-amber-700 dark:text-amber-300 border-amber-500/30',
      description: 'Crew assigned, call-sheet generated, moodboard ready.',
    },
    {
      id: 'stg-3',
      name: 'Shoot Completed',
      color: 'bg-indigo-500/15 text-indigo-700 dark:text-indigo-300 border-indigo-500/30',
      description: 'Media backed up onto studio RAID storage.',
    },
    {
      id: 'stg-4',
      name: 'Editing',
      color: 'bg-purple-500/15 text-purple-700 dark:text-purple-300 border-purple-500/30',
      description: 'Color grading, highlight reels, and album layout design.',
    },
    {
      id: 'stg-5',
      name: 'Delivery',
      color: 'bg-teal-500/15 text-teal-700 dark:text-teal-300 border-teal-500/30',
      description: 'Client proofing, print dispatch, and cloud album link shared.',
    },
    {
      id: 'stg-6',
      name: 'Closure',
      color: 'bg-emerald-500/15 text-emerald-700 dark:text-emerald-300 border-emerald-500/30',
      description: 'All payments settled & client testimonial captured.',
    },
  ])

  const [newStageName, setNewStageName] = React.useState('')
  const [newStageDesc, setNewStageDesc] = React.useState('')

  const handleAddStage = (e: React.FormEvent) => {
    e.preventDefault()
    if (!newStageName.trim()) return

    const newStage: WorkflowStageItem = {
      id: `stg-${Date.now()}`,
      name: newStageName.trim(),
      color: 'bg-primary/10 text-primary border-primary/20',
      description: newStageDesc.trim() || 'Custom workflow progression step.',
    }

    setStages((prev) => [...prev, newStage])
    setNewStageName('')
    setNewStageDesc('')
    toast.success(`Workflow stage "${newStage.name}" added`)
  }

  const handleDeleteStage = (id: string) => {
    if (stages.length <= 2) {
      toast.error('Workflow requires at least 2 stages')
      return
    }
    setStages((prev) => prev.filter((s) => s.id !== id))
    toast.success('Stage removed')
  }

  return (
    <div className="space-y-6 max-w-4xl mx-auto pb-16">
      <div className="flex items-center justify-between border-b pb-4">
        <div>
          <h1 className="text-xl font-bold tracking-tight text-foreground">Pipeline Stages & Workflow</h1>
          <p className="text-xs text-muted-foreground">
            Customize the step-by-step production progression across projects and orders.
          </p>
        </div>
      </div>

      {/* Add Stage Form */}
      <form
        onSubmit={handleAddStage}
        className="bg-card p-4 rounded-xl border border-border/80 shadow-2xs flex flex-wrap items-center gap-3"
      >
        <Input
          value={newStageName}
          onChange={(e) => setNewStageName(e.target.value)}
          placeholder="New stage title (e.g. Quality Check)..."
          className="h-9 text-xs flex-1 min-w-[180px]"
          required
        />
        <Input
          value={newStageDesc}
          onChange={(e) => setNewStageDesc(e.target.value)}
          placeholder="Stage purpose / description..."
          className="h-9 text-xs flex-1 min-w-[200px]"
        />
        <Button
          type="submit"
          size="sm"
          className="bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs h-9 px-4 gap-1.5 shadow-xs"
        >
          <IconPlus className="h-4 w-4" />
          Add Stage
        </Button>
      </form>

      {/* Stepper visualizer */}
      <div className="space-y-3">
        <span className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
          Current Sequence ({stages.length} Steps)
        </span>

        <div className="space-y-2.5">
          {stages.map((stage, idx) => (
            <div
              key={stage.id}
              className="flex items-center justify-between p-3.5 rounded-xl border bg-card shadow-2xs hover:border-primary/40 transition text-xs"
            >
              <div className="flex items-center gap-3">
                <span className="h-6 w-6 rounded-full bg-muted flex items-center justify-center font-bold text-[11px] text-muted-foreground">
                  {idx + 1}
                </span>
                <div>
                  <div className="flex items-center gap-2">
                    <span className="font-bold text-foreground">{stage.name}</span>
                    <span
                      className={`text-[10px] font-semibold px-2 py-0.5 rounded-full border ${stage.color}`}
                    >
                      Step {idx + 1}
                    </span>
                  </div>
                  <p className="text-[11px] text-muted-foreground mt-0.5">{stage.description}</p>
                </div>
              </div>

              <Button
                variant="ghost"
                size="sm"
                onClick={() => handleDeleteStage(stage.id)}
                className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive"
              >
                <IconTrash className="h-3.5 w-3.5" />
              </Button>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
