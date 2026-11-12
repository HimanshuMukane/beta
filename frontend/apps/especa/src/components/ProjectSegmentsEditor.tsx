import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@repo/ui/components/ui/select'
import { IconCalendar, IconTrash, IconGripVertical } from "@repo/ui/icons/tabler-icons"
import type { ProjectSegment } from '../types/project'
import { PROJECT_TYPES } from '../lib/formatters'

interface ProjectSegmentsEditorProps {
  segments: ProjectSegment[]
  onChange: (segments: ProjectSegment[]) => void
}

export function ProjectSegmentsEditor({ segments, onChange }: ProjectSegmentsEditorProps) {
  const handleAddSegment = () => {
    const nextIdx = segments.length + 1
    const newSeg: ProjectSegment = {
      id: `seg-${Date.now()}`,
      name: `Segment ${nextIdx}`,
      type: 'Wedding',
      startDate: new Date().toISOString().split('T')[0],
      endDate: new Date().toISOString().split('T')[0],
    }
    onChange([...segments, newSeg])
  }

  const handleRemoveSegment = (id: string) => {
    onChange(segments.filter((s) => s.id !== id))
  }

  const handleUpdateSegment = (id: string, field: keyof ProjectSegment, val: string) => {
    onChange(
      segments.map((s) => (s.id === id ? { ...s, [field]: val } : s))
    )
  }

  return (
    <div className="pt-2">
      <div className="flex items-center justify-between mb-2">
        <div className="flex items-center gap-1.5 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
          <IconCalendar className="h-4 w-4 text-primary" />
          <span>Project Segments</span>
        </div>
        <span className="text-[11px] text-muted-foreground">
          Add one segment per date block
        </span>
      </div>

      <div className="space-y-3">
        {segments.map((seg, idx) => (
          <div
            key={seg.id}
            className="rounded-lg border bg-muted/20 p-3.5 space-y-2 relative transition hover:border-border"
          >
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <IconGripVertical className="h-4 w-4 text-muted-foreground cursor-grab" />
                <span className="text-xs font-bold text-foreground">
                  #{idx + 1}
                </span>
                <Input
                  value={seg.name}
                  onChange={(e) => handleUpdateSegment(seg.id, 'name', e.target.value)}
                  placeholder="Segment title"
                  className="h-8 text-xs font-semibold w-40 bg-background"
                />
              </div>

              {segments.length > 1 && (
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => handleRemoveSegment(seg.id)}
                  className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive"
                  title="Remove segment"
                >
                  <IconTrash className="h-3.5 w-3.5" />
                </Button>
              )}
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 pt-1">
              <div className="space-y-1">
                <span className="text-[10px] uppercase font-bold text-muted-foreground">
                  Shoot Type
                </span>
                <Select
                  value={seg.type}
                  onValueChange={(val) => handleUpdateSegment(seg.id, 'type', val)}
                >
                  <SelectTrigger className="h-8 text-xs bg-background">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {PROJECT_TYPES.map((pt) => (
                      <SelectItem key={pt} value={pt} className="text-xs">
                        {pt}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-1">
                <span className="text-[10px] uppercase font-bold text-muted-foreground">
                  Date Range
                </span>
                <div className="flex items-center gap-1.5">
                  <Input
                    type="date"
                    value={seg.startDate}
                    onChange={(e) => handleUpdateSegment(seg.id, 'startDate', e.target.value)}
                    className="h-8 text-xs bg-background px-2"
                  />
                  <span className="text-xs text-muted-foreground">→</span>
                  <Input
                    type="date"
                    value={seg.endDate}
                    onChange={(e) => handleUpdateSegment(seg.id, 'endDate', e.target.value)}
                    className="h-8 text-xs bg-background px-2"
                  />
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="pt-2">
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={handleAddSegment}
          className="text-xs text-emerald-600 dark:text-emerald-400 hover:text-emerald-700 p-0 h-auto font-medium"
        >
          + Add Projects segment
        </Button>
        <p className="text-[11px] text-muted-foreground mt-1">
          Long-duration Projects are easier to manage in calendars with segments.
        </p>
      </div>
    </div>
  )
}

