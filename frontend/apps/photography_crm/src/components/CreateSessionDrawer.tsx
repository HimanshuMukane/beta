import * as React from 'react'
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerDescription,
  DrawerFooter,
  DrawerClose,
} from '@repo/ui/components/ui/drawer'
import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import { Label } from '@repo/ui/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@repo/ui/components/ui/select'
import { Textarea } from '@repo/ui/components/ui/textarea'
import { Badge } from '@repo/ui/components/ui/badge'
import { toast } from '@repo/ui/components/ui/sonner'
import { cn } from '@repo/ui/lib/utils'
import {
  IconClock,
  IconMapPin,
  IconNotes,
  IconUsers,
  IconPlus,
  IconX,
  IconSparkles,
  IconUserCheck,
  IconArrowUpRight,
} from '@tabler/icons-react'
import type {
  ProjectSession,
  ProjectType,
  ShootTiming,
  SessionStatus,
} from '../types/session'
import {
  PROJECT_TYPES,
  SHOOT_TIMINGS,
  DURATION_OPTIONS,
  DEFAULT_AVAILABLE_CREW,
} from '../types/session'

interface CreateSessionDrawerProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSave: (session: ProjectSession) => void
  initialSession?: ProjectSession | null
  existingProjectNames?: string[]
}

const SNAP_POINTS = ['34rem', 1]

export function CreateSessionDrawer({
  open,
  onOpenChange,
  onSave,
  initialSession,
  existingProjectNames = [],
}: CreateSessionDrawerProps) {
  // Form fields
  const [sessionManager, setSessionManager] = React.useState('especa')
  const [projectType, setProjectType] = React.useState<ProjectType>('Wedding')
  const [shootTiming, setShootTiming] = React.useState<ShootTiming>('Evening')
  const [status, setStatus] = React.useState<SessionStatus>('New')

  const [projectName, setProjectName] = React.useState('')
  const [date, setDate] = React.useState('2026-11-11')
  const [startTime, setStartTime] = React.useState('18:00')
  const [duration, setDuration] = React.useState('4 hours')

  const [venue, setVenue] = React.useState('')
  const [locationLink, setLocationLink] = React.useState('')
  const [remarks, setRemarks] = React.useState('')

  // Crew Management
  const [assignedCrew, setAssignedCrew] = React.useState<string[]>([])
  const [isManagingCrew, setIsManagingCrew] = React.useState(false)

  // Populate form if editing or opening
  React.useEffect(() => {
    if (initialSession) {
      setSessionManager(initialSession.sessionManager || 'especa')
      setProjectType(initialSession.projectType || 'Wedding')
      setShootTiming(initialSession.shootTiming || 'Evening')
      setStatus(initialSession.status || 'New')
      setProjectName(initialSession.projectName || initialSession.title || '')
      setDate(initialSession.date || '2026-11-11')
      setStartTime(initialSession.startTime || '18:00')
      setDuration(initialSession.duration || '4 hours')
      setVenue(initialSession.venue || '')
      setLocationLink(initialSession.locationLink || '')
      setRemarks(initialSession.remarks || '')
      setAssignedCrew(initialSession.members || [])
    } else {
      setSessionManager('especa')
      setProjectType('Wedding')
      setShootTiming('Evening')
      setStatus('New')
      setProjectName('')
      setDate('2026-11-11')
      setStartTime('18:00')
      setDuration('4 hours')
      setVenue('')
      setLocationLink('')
      setRemarks('')
      setAssignedCrew([])
    }
    setIsManagingCrew(false)
  }, [initialSession, open])

  // Calculated End Time display for UX delight
  const formattedEndTime = React.useMemo(() => {
    if (!startTime) return ''
    try {
      const [hStr, mStr] = startTime.split(':')
      const startH = parseInt(hStr, 10)
      const startM = parseInt(mStr, 10)

      let totalMinutesToAdd = 240
      const hourMatch = duration.match(/(\d+)\s*hour/)
      const minMatch = duration.match(/(\d+)\s*min/)
      if (hourMatch) totalMinutesToAdd = parseInt(hourMatch[1], 10) * 60
      if (minMatch) totalMinutesToAdd += parseInt(minMatch[1], 10)

      const startTotal = startH * 60 + startM
      const endTotal = (startTotal + totalMinutesToAdd) % (24 * 60)
      const endH = Math.floor(endTotal / 60)
      const endM = endTotal % 60

      const format12 = (h: number, m: number) => {
        const period = h >= 12 ? 'PM' : 'AM'
        const h12 = h % 12 || 12
        return `${String(h12).padStart(2, '0')}:${String(m).padStart(2, '0')} ${period}`
      }

      return `${format12(startH, startM)} - ${format12(endH, endM)}`
    } catch {
      return startTime
    }
  }, [startTime, duration])

  const toggleCrewMember = (memberName: string) => {
    setAssignedCrew((prev) =>
      prev.includes(memberName)
        ? prev.filter((m) => m !== memberName)
        : [...prev, memberName]
    )
  }

  const handleSubmit = (e?: React.FormEvent) => {
    if (e) e.preventDefault()

    if (!projectName.trim()) {
      toast.error('Project Name is required')
      return
    }

    if (!date) {
      toast.error('Date is required')
      return
    }

    const newSession: ProjectSession = {
      id: initialSession?.id || `ses-${Date.now()}`,
      title: projectName,
      projectName: projectName,
      sessionManager,
      projectType,
      shootTiming,
      status,
      date,
      startTime,
      duration,
      timeRange: formattedEndTime,
      venue: venue || 'Main Location',
      city: venue ? venue.split(',')[1]?.trim() || venue : 'Local',
      locationLink,
      remarks,
      members: assignedCrew,
      shootsActivityCount: initialSession?.shootsActivityCount || 0,
    }

    onSave(newSession)
    toast.success(
      initialSession
        ? 'Session updated successfully!'
        : 'Session created successfully!'
    )
    onOpenChange(false)
  }

  return (
    <Drawer
      open={open}
      onOpenChange={onOpenChange}
      snapPoints={SNAP_POINTS}
    >
      <DrawerContent className="max-h-[94vh] flex flex-col bg-background border-border">
        {/* Top Controls Bar: Session Manager | Project Type | Shoot Timing | Status */}
        <div className="border-b bg-muted/20 px-4 sm:px-6 py-3">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-3 max-w-4xl mx-auto">
            {/* Session Manager */}
            <div className="space-y-1">
              <label className="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">
                Session Manager
              </label>
              <Select value={sessionManager} onValueChange={setSessionManager}>
                <SelectTrigger className="h-9 bg-card border-border shadow-xs text-xs">
                  <div className="flex items-center gap-2 truncate">
                    <span className="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-primary/10 text-[10px] font-bold text-primary">
                      {sessionManager.charAt(0).toUpperCase()}
                    </span>
                    <span className="truncate">{sessionManager}</span>
                  </div>
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="especa">especa</SelectItem>
                  <SelectItem value="Rohan Mehta">Rohan Mehta</SelectItem>
                  <SelectItem value="Priya Sharma">Priya Sharma</SelectItem>
                  <SelectItem value="Aditya Roy">Aditya Roy</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Project Type */}
            <div className="space-y-1">
              <label className="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">
                Project Type
              </label>
              <Select
                value={projectType}
                onValueChange={(val) => setProjectType(val as ProjectType)}
              >
                <SelectTrigger className="h-9 bg-card border-border shadow-xs text-xs">
                  <SelectValue placeholder="Project Type" />
                </SelectTrigger>
                <SelectContent className="max-h-60">
                  {PROJECT_TYPES.map((pt) => (
                    <SelectItem key={pt} value={pt} className="text-xs">
                      {pt}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            {/* Shoot Timing */}
            <div className="space-y-1">
              <label className="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">
                Shoot Timing
              </label>
              <Select
                value={shootTiming}
                onValueChange={(val) => setShootTiming(val as ShootTiming)}
              >
                <SelectTrigger className="h-9 bg-card border-border shadow-xs text-xs">
                  <SelectValue placeholder="Shoot Timing" />
                </SelectTrigger>
                <SelectContent>
                  {SHOOT_TIMINGS.map((st) => (
                    <SelectItem key={st} value={st} className="text-xs">
                      {st}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            {/* Status with dynamic accent color */}
            <div className="space-y-1">
              <label className="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">
                Status
              </label>
              <Select
                value={status}
                onValueChange={(val) => setStatus(val as SessionStatus)}
              >
                <SelectTrigger
                  className={cn(
                    'h-9 font-semibold text-xs border shadow-xs transition',
                    status === 'New' && 'bg-blue-600 text-white hover:bg-blue-700 border-blue-600',
                    status === 'Scheduled' && 'bg-amber-600 text-white hover:bg-amber-700 border-amber-600',
                    status === 'Completed' && 'bg-emerald-600 text-white hover:bg-emerald-700 border-emerald-600',
                    status === 'Delayed' && 'bg-rose-600 text-white hover:bg-rose-700 border-rose-600'
                  )}
                >
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="New">New</SelectItem>
                  <SelectItem value="Scheduled">Scheduled</SelectItem>
                  <SelectItem value="Completed">Completed</SelectItem>
                  <SelectItem value="Delayed">Delayed</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </div>

        {/* Drawer Scrollable Body Form */}
        <div className="flex-1 overflow-y-auto px-4 sm:px-6 py-4">
          <div className="max-w-4xl mx-auto space-y-5">
            <DrawerHeader className="p-0 text-left">
              <div className="flex items-center justify-between">
                <div>
                  <DrawerTitle className="text-lg font-bold text-foreground flex items-center gap-2">
                    <IconSparkles className="h-5 w-5 text-emerald-600" />
                    {initialSession ? 'Edit Photography Session' : 'Create Session'}
                  </DrawerTitle>
                  <DrawerDescription className="text-xs text-muted-foreground mt-0.5">
                    Plan session timing, assign crew members, set venues, and capture shooting instructions.
                  </DrawerDescription>
                </div>
                {formattedEndTime && (
                  <Badge variant="outline" className="hidden sm:inline-flex items-center gap-1.5 py-1 px-2.5 text-xs font-medium bg-card">
                    <IconClock className="h-3.5 w-3.5 text-emerald-600" />
                    {formattedEndTime}
                  </Badge>
                )}
              </div>
            </DrawerHeader>

            {/* Inner Content Card (matches Screenshot 2) */}
            <div className="rounded-2xl border border-border bg-card p-5 sm:p-6 shadow-xs space-y-5">
              {/* Project Name */}
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="projectName" className="text-xs font-semibold text-foreground">
                    Project Name <span className="text-destructive font-bold">*</span>
                  </Label>
                  {existingProjectNames.length > 0 && (
                    <span className="text-[11px] text-muted-foreground">
                      Or select existing project below
                    </span>
                  )}
                </div>
                <Input
                  id="projectName"
                  value={projectName}
                  onChange={(e) => setProjectName(e.target.value)}
                  placeholder="Enter Project Name (e.g. Kavya Rajendran Wedding - Reception | Sunset)"
                  className="h-10 text-xs bg-background/50 border-input focus:bg-background"
                  list="existing-projects-datalist"
                  required
                />
                {existingProjectNames.length > 0 && (
                  <datalist id="existing-projects-datalist">
                    {existingProjectNames.map((pName) => (
                      <option key={pName} value={pName} />
                    ))}
                  </datalist>
                )}
              </div>

              {/* Date, Start Time, Duration */}
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                {/* Date */}
                <div className="space-y-1.5">
                  <Label htmlFor="sessionDate" className="text-xs font-semibold text-foreground">
                    Date <span className="text-destructive font-bold">*</span>
                  </Label>
                  <div className="relative">
                    <Input
                      id="sessionDate"
                      type="date"
                      value={date}
                      onChange={(e) => setDate(e.target.value)}
                      className="h-10 text-xs bg-background/50 border-input focus:bg-background"
                      required
                    />
                  </div>
                </div>

                {/* Start Time */}
                <div className="space-y-1.5">
                  <Label htmlFor="startTime" className="text-xs font-semibold text-foreground">
                    Start Time <span className="text-destructive font-bold">*</span>
                  </Label>
                  <div className="relative">
                    <Input
                      id="startTime"
                      type="time"
                      value={startTime}
                      onChange={(e) => setStartTime(e.target.value)}
                      className="h-10 text-xs bg-background/50 border-input focus:bg-background"
                      required
                    />
                  </div>
                </div>

                {/* Duration dropdown */}
                <div className="space-y-1.5">
                  <Label className="text-xs font-semibold text-foreground">
                    Duration <span className="text-destructive font-bold">*</span>
                  </Label>
                  <Select value={duration} onValueChange={setDuration}>
                    <SelectTrigger className="h-10 text-xs bg-background/50 border-input focus:bg-background">
                      <SelectValue placeholder="Duration" />
                    </SelectTrigger>
                    <SelectContent className="max-h-60">
                      {DURATION_OPTIONS.map((opt) => (
                        <SelectItem key={opt} value={opt} className="text-xs">
                          {opt}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>

              {/* Assigned Crew / Shoot Team */}
              <div className="space-y-2 pt-1 border-t border-dashed">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Label className="text-xs font-semibold text-foreground flex items-center gap-1.5">
                      <IconUsers className="h-4 w-4 text-muted-foreground" />
                      Assigned Crew / Shoot Team
                    </Label>
                    {assignedCrew.length > 0 && (
                      <Badge variant="secondary" className="text-[10px] h-5 px-1.5">
                        {assignedCrew.length} Assigned
                      </Badge>
                    )}
                  </div>
                  <button
                    type="button"
                    onClick={() => setIsManagingCrew(!isManagingCrew)}
                    className="text-xs font-semibold text-primary hover:underline flex items-center gap-1 transition"
                  >
                    <span>{isManagingCrew ? 'Done Managing' : 'Manage Crew'}</span>
                    <IconArrowUpRight className="h-3 w-3" />
                  </button>
                </div>

                {/* Assigned members pill list or Empty placeholder */}
                {assignedCrew.length === 0 && !isManagingCrew ? (
                  <div className="py-7 text-center rounded-xl bg-muted/15 border border-dashed border-border">
                    <p className="text-xs italic text-muted-foreground font-medium">
                      No Members Added
                    </p>
                    <button
                      type="button"
                      onClick={() => setIsManagingCrew(true)}
                      className="text-[11px] font-semibold text-emerald-600 hover:underline mt-1 inline-flex items-center gap-1"
                    >
                      <IconPlus className="h-3 w-3" />
                      Add photographers or assistants
                    </button>
                  </div>
                ) : (
                  <div className="flex flex-wrap gap-2 pt-1">
                    {assignedCrew.map((name) => (
                      <span
                        key={name}
                        className="inline-flex items-center gap-1.5 rounded-full bg-secondary/80 border px-3 py-1 text-xs font-medium text-foreground shadow-2xs"
                      >
                        <span className="h-2 w-2 rounded-full bg-emerald-500" />
                        {name}
                        <button
                          type="button"
                          onClick={() => toggleCrewMember(name)}
                          className="ml-0.5 rounded-full p-0.5 text-muted-foreground hover:bg-muted hover:text-foreground"
                        >
                          <IconX className="h-3 w-3" />
                        </button>
                      </span>
                    ))}
                  </div>
                )}

                {/* Crew selection drawer/dropdown panel when 'Manage Crew' is opened */}
                {isManagingCrew && (
                  <div className="rounded-xl border bg-muted/20 p-3 mt-2 space-y-2.5 animate-in fade-in-50 duration-200">
                    <div className="flex items-center justify-between">
                      <span className="text-xs font-semibold text-foreground">
                        Select Crew Members for this Shoot
                      </span>
                      <span className="text-[11px] text-muted-foreground">
                        Click member to toggle assignment
                      </span>
                    </div>
                    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-2">
                      {DEFAULT_AVAILABLE_CREW.map((crew) => {
                        const isSelected = assignedCrew.includes(crew.name)
                        return (
                          <button
                            type="button"
                            key={crew.id}
                            onClick={() => toggleCrewMember(crew.name)}
                            className={cn(
                              'flex items-center justify-between p-2 rounded-lg border text-xs text-left transition',
                              isSelected
                                ? 'bg-primary/10 border-primary text-primary font-semibold shadow-2xs'
                                : 'bg-card border-border hover:border-muted-foreground/40 text-foreground'
                            )}
                          >
                            <div className="flex items-center gap-2 truncate">
                              <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-muted font-bold text-[10px]">
                                {crew.name.split(' ').map((n) => n[0]).join('')}
                              </span>
                              <div className="truncate">
                                <p className="truncate leading-tight">{crew.name}</p>
                                <p className="text-[10px] text-muted-foreground font-normal">
                                  {crew.role}
                                </p>
                              </div>
                            </div>
                            {isSelected && <IconUserCheck className="h-4 w-4 shrink-0 text-primary" />}
                          </button>
                        )
                      })}
                    </div>
                  </div>
                )}
              </div>

              {/* Venue & Shoot Location Link */}
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 pt-1 border-t border-dashed">
                {/* Venue */}
                <div className="space-y-1.5">
                  <Label htmlFor="venueInput" className="text-xs font-semibold text-foreground flex items-center gap-1.5">
                    <IconMapPin className="h-3.5 w-3.5 text-muted-foreground" />
                    Venue
                  </Label>
                  <Input
                    id="venueInput"
                    value={venue}
                    onChange={(e) => setVenue(e.target.value)}
                    placeholder="Enter Venue (e.g. Leela Palace)"
                    className="h-10 text-xs bg-background/50 border-input focus:bg-background"
                  />
                </div>

                {/* Shoot Location / Location Link */}
                <div className="space-y-1.5">
                  <Label htmlFor="locationLink" className="text-xs font-semibold text-foreground">
                    Shoot Location / Location Link
                  </Label>
                  <Input
                    id="locationLink"
                    value={locationLink}
                    onChange={(e) => setLocationLink(e.target.value)}
                    placeholder="Enter Location / Google Maps Link (e.g. Chennai or https://maps...)"
                    className="h-10 text-xs bg-background/50 border-input focus:bg-background"
                  />
                </div>
              </div>

              {/* Remarks */}
              <div className="space-y-1.5 pt-1 border-t border-dashed">
                <Label htmlFor="remarksInput" className="text-xs font-semibold text-foreground flex items-center gap-1.5">
                  <IconNotes className="h-3.5 w-3.5 text-muted-foreground" />
                  Remarks
                </Label>
                <Textarea
                  id="remarksInput"
                  value={remarks}
                  onChange={(e) => setRemarks(e.target.value)}
                  placeholder="Enter Remarks (e.g. Bring extra 85mm prime lens, drone clearance granted, sunset schedule)"
                  className="min-h-[90px] text-xs bg-background/50 border-input focus:bg-background resize-y"
                />
              </div>
            </div>
          </div>
        </div>

        {/* Drawer Bottom Sticky Footer */}
        <DrawerFooter className="border-t bg-card px-4 sm:px-6 py-3.5">
          <div className="flex flex-col sm:flex-row items-center justify-between gap-3 max-w-4xl mx-auto w-full">
            <p className="text-xs text-muted-foreground hidden sm:block">
              All timings sync automatically with shoot crew calendar & tasks.
            </p>
            <div className="flex items-center gap-3 w-full sm:w-auto justify-end">
              <DrawerClose asChild>
                <Button variant="outline" className="h-10 px-5 text-xs font-semibold">
                  Cancel
                </Button>
              </DrawerClose>
              <Button
                onClick={() => handleSubmit()}
                className="h-10 px-6 text-xs font-semibold bg-emerald-600 hover:bg-emerald-700 text-white shadow-xs"
              >
                {initialSession ? 'Update Session' : 'Create Session'}
              </Button>
            </div>
          </div>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  )
}
