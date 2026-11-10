import * as React from 'react'
import { Card } from '@repo/ui/components/ui/card'
import { Button } from '@repo/ui/components/ui/button'
import { Badge } from '@repo/ui/components/ui/badge'
import { Input } from '@repo/ui/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@repo/ui/components/ui/select'
import { toast } from '@repo/ui/components/ui/sonner'
import { cn } from '@repo/ui/lib/utils'
import {
  IconCalendar,
  IconClock,
  IconMapPin,
  IconFilter,
  IconSearch,
  IconUsers,
  IconCheck,
  IconRotateClockwise2,
} from '@tabler/icons-react'
import { useAuth } from '@repo/ui/hooks/use-auth'
import type { ProjectSession, SessionStatus } from '../../types/session'
import { PROJECT_TYPES } from '../../types/session'
import { CreateSessionDrawer } from '../../components/CreateSessionDrawer'

interface UserShootsPageProps {
  sessions: ProjectSession[]
  onSaveSession: (session: ProjectSession) => void
  onStatusChange?: (sessionId: string, newStatus: SessionStatus) => void
  existingProjectNames?: string[]
}

type TabType = 'Upcoming' | 'Past'

export function UserShootsPage({
  sessions,
  onSaveSession,
  onStatusChange,
  existingProjectNames = [],
}: UserShootsPageProps) {
  const { user } = useAuth()
  const [activeTab, setActiveTab] = React.useState<TabType>('Upcoming')
  const [sortOption, setSortOption] = React.useState<'date_asc' | 'date_desc' | 'recently_created'>('date_asc')
  const [searchQuery, setSearchQuery] = React.useState('')
  const [isFilterOpen, setIsFilterOpen] = React.useState(false)
  const [typeFilter, setTypeFilter] = React.useState<string>('All')
  const [crewUnassignedOnly, setCrewUnassignedOnly] = React.useState(false)
  const [yourShootsOnly, setYourShootsOnly] = React.useState(false)

  const [isDrawerOpen, setIsDrawerOpen] = React.useState(false)
  const [selectedSessionToEdit, setSelectedSessionToEdit] = React.useState<ProjectSession | null>(null)

  const handleOpenAssignDrawer = (session: ProjectSession) => {
    setSelectedSessionToEdit(session)
    setIsDrawerOpen(true)
  }

  const isAssignedToUser = (session: ProjectSession) => {
    if (!user) return false
    const userName = `${user.first_name || ''} ${user.last_name || ''}`.trim().toLowerCase()
    const userEmail = (user.email || '').toLowerCase()
    return (session.members || []).some((m) => {
      const mName = m.toLowerCase()
      return (userName && mName.includes(userName)) || (userEmail && mName.includes(userEmail))
    })
  }

  const filteredShoots = React.useMemo(() => {
    const todayStr = new Date().toISOString().split('T')[0]

    return sessions.filter((s) => {
      if (searchQuery.trim()) {
        const q = searchQuery.toLowerCase()
        const matchTitle = s.title?.toLowerCase().includes(q)
        const matchProj = s.projectName?.toLowerCase().includes(q)
        const matchCity = s.city?.toLowerCase().includes(q)
        const matchVenue = s.venue?.toLowerCase().includes(q)
        if (!matchTitle && !matchProj && !matchCity && !matchVenue) return false
      }

      if (typeFilter !== 'All' && s.projectType !== typeFilter) {
        return false
      }

      if (crewUnassignedOnly && s.members && s.members.length > 0) {
        return false
      }

      if (yourShootsOnly && !isAssignedToUser(s)) {
        return false
      }

      const sessionDateStr = s.date.includes(' ')
        ? `${s.date.split(' ')[0]}`
        : s.date

      if (activeTab === 'Upcoming') {
        return sessionDateStr >= todayStr
      }
      if (activeTab === 'Past') {
        return sessionDateStr < todayStr
      }

      return true
    }).sort((a, b) => {
      if (sortOption === 'date_asc') return a.date.localeCompare(b.date)
      if (sortOption === 'date_desc') return b.date.localeCompare(a.date)
      return b.id.localeCompare(a.id)
    })
  }, [sessions, activeTab, sortOption, searchQuery, typeFilter, crewUnassignedOnly, yourShootsOnly, user])

  const groupedShoots = React.useMemo(() => {
    const groups: { [key: string]: ProjectSession[] } = {}

    filteredShoots.forEach((ses) => {
      let groupKey = 'November, 2026'
      if (ses.date && ses.date.includes('-')) {
        const [year, month] = ses.date.split('-')
        const monthNames = [
          'January', 'February', 'March', 'April', 'May', 'June',
          'July', 'August', 'September', 'October', 'November', 'December'
        ]
        const mIdx = parseInt(month, 10) - 1
        groupKey = `${monthNames[mIdx] || 'November'}, ${year || '2026'}`
      }
      if (!groups[groupKey]) {
        groups[groupKey] = []
      }
      groups[groupKey].push(ses)
    })

    return groups
  }, [filteredShoots])

  const getDayInfo = (dateStr: string) => {
    if (!dateStr || dateStr === 'TBD') return { day: 'TBD', dayNum: '—' }

    if (dateStr.includes('-')) {
      const d = new Date(dateStr)
      if (!isNaN(d.getTime())) {
        const days = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
        return {
          day: days[d.getDay()],
          dayNum: String(d.getDate()).padStart(2, '0'),
        }
      }
    }

    const match = dateStr.match(/(\w{3})\s*(\d+)/)
    if (match) {
      return { day: match[1], dayNum: match[2].padStart(2, '0') }
    }

    return { day: 'Wed', dayNum: '11' }
  }

  const upcomingCount = React.useMemo(() => {
    const todayStr = new Date().toISOString().split('T')[0]
    return sessions.filter((s) => s.date >= todayStr).length
  }, [sessions])

  const pastCount = React.useMemo(() => {
    const todayStr = new Date().toISOString().split('T')[0]
    return sessions.filter((s) => s.date < todayStr).length
  }, [sessions])

  return (
    <div className="space-y-6">
      {/* Top Filter and Controls Bar */}
      <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
        {/* Left Tabs: Upcoming | Past */}
        <div className="inline-flex items-center rounded-full bg-muted/60 p-1 border shadow-2xs">
          <button
            type="button"
            onClick={() => setActiveTab('Upcoming')}
            className={cn(
              'inline-flex items-center gap-2 rounded-full px-5 py-1.5 text-xs font-semibold transition-all',
              activeTab === 'Upcoming'
                ? 'bg-foreground text-background shadow-xs font-bold'
                : 'text-muted-foreground hover:text-foreground'
            )}
          >
            <span>Upcoming</span>
            <span className={cn(
              'text-[10px] px-1.5 py-0.2 rounded-full font-bold',
              activeTab === 'Upcoming' ? 'bg-background/20 text-background' : 'bg-muted text-muted-foreground'
            )}>
              {upcomingCount}
            </span>
          </button>
          <button
            type="button"
            onClick={() => setActiveTab('Past')}
            className={cn(
              'inline-flex items-center gap-2 rounded-full px-5 py-1.5 text-xs font-semibold transition-all',
              activeTab === 'Past'
                ? 'bg-foreground text-background shadow-xs font-bold'
                : 'text-muted-foreground hover:text-foreground'
            )}
          >
            <span>Past</span>
            <span className={cn(
              'text-[10px] px-1.5 py-0.2 rounded-full font-bold',
              activeTab === 'Past' ? 'bg-background/20 text-background' : 'bg-muted text-muted-foreground'
            )}>
              {pastCount}
            </span>
          </button>
        </div>

        {/* Right Controls: Crew unassigned | Your Shoots | Filters | Sort */}
        <div className="flex items-center gap-2 flex-wrap">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setCrewUnassignedOnly(!crewUnassignedOnly)}
            className={cn(
              'h-9 rounded-full text-xs gap-1.5 px-3.5 border-border shadow-2xs transition-all font-medium',
              crewUnassignedOnly && 'bg-primary/10 border-primary/40 text-primary font-bold'
            )}
          >
            <IconCheck className={cn('h-3.5 w-3.5', crewUnassignedOnly ? 'opacity-100 text-primary' : 'opacity-40')} />
            <span>Crew unassigned</span>
          </Button>

          <Button
            variant="outline"
            size="sm"
            onClick={() => setYourShootsOnly(!yourShootsOnly)}
            className={cn(
              'h-9 rounded-full text-xs gap-1.5 px-3.5 border-border shadow-2xs transition-all font-medium',
              yourShootsOnly && 'bg-primary/10 border-primary/40 text-primary font-bold'
            )}
          >
            <IconCheck className={cn('h-3.5 w-3.5', yourShootsOnly ? 'opacity-100 text-primary' : 'opacity-40')} />
            <span>Your Shoots</span>
          </Button>

          <Button
            variant="outline"
            size="sm"
            onClick={() => setIsFilterOpen(!isFilterOpen)}
            className={cn(
              'h-9 rounded-full text-xs gap-1.5 px-3.5 border-border shadow-2xs transition-all font-medium',
              isFilterOpen && 'bg-accent text-accent-foreground border-accent'
            )}
          >
            <IconFilter className="h-3.5 w-3.5" />
            <span>Filters</span>
          </Button>

          <Select
            value={sortOption}
            onValueChange={(val) => setSortOption(val as any)}
          >
            <SelectTrigger className="h-9 rounded-full text-xs bg-card border-border shadow-2xs w-[130px]">
              <div className="flex items-center gap-1.5 text-muted-foreground truncate">
                <IconClock className="h-3.5 w-3.5 shrink-0" />
                <span className="font-semibold text-foreground text-xs truncate">
                  {sortOption === 'date_asc'
                    ? 'Earliest'
                    : sortOption === 'date_desc'
                      ? 'Latest'
                      : 'Recent'}
                </span>
              </div>
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="date_asc" className="text-xs">
                Date (Earliest first)
              </SelectItem>
              <SelectItem value="date_desc" className="text-xs">
                Date (Latest first)
              </SelectItem>
              <SelectItem value="recently_created" className="text-xs">
                Recently Created
              </SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Collapsible Search and Filters drawer */}
      {isFilterOpen && (
        <Card className="p-4 border-border/80 bg-muted/20 animate-in fade-in-50 duration-200">
          <div className="flex flex-col sm:flex-row items-center gap-3">
            <div className="relative flex-1 w-full">
              <IconSearch className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search shoots by title, project, city, venue..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-9 h-9 text-xs bg-background"
              />
            </div>
            <div className="w-full sm:w-[200px]">
              <Select value={typeFilter} onValueChange={setTypeFilter}>
                <SelectTrigger className="h-9 text-xs bg-background">
                  <SelectValue placeholder="Event Type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="All" className="text-xs">
                    All Shoot Types
                  </SelectItem>
                  {PROJECT_TYPES.map((t) => (
                    <SelectItem key={t} value={t} className="text-xs">
                      {t}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            {(searchQuery || typeFilter !== 'All' || crewUnassignedOnly || yourShootsOnly) && (
              <Button
                variant="ghost"
                size="sm"
                className="text-xs h-9 px-3"
                onClick={() => {
                  setSearchQuery('')
                  setTypeFilter('All')
                  setCrewUnassignedOnly(false)
                  setYourShootsOnly(false)
                }}
              >
                Reset
              </Button>
            )}
          </div>
        </Card>
      )}

      {/* Shoots Grouped by Month */}
      <div className="space-y-6">
        {Object.keys(groupedShoots).length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-center border border-dashed rounded-2xl bg-card">
            <IconCalendar className="h-10 w-10 text-muted-foreground/30 mb-3" />
            <h3 className="text-sm font-semibold text-foreground">No shoots found</h3>
            <p className="text-xs text-muted-foreground max-w-sm mt-1 mb-2">
              {searchQuery || typeFilter !== 'All' || crewUnassignedOnly || yourShootsOnly
                ? 'No shoots match your filter criteria.'
                : 'There are no shoots scheduled in this timeline.'}
            </p>
          </div>
        ) : (
          Object.entries(groupedShoots).map(([monthGroup, monthShoots]) => (
            <div key={monthGroup} className="space-y-3">
              <h3 className="text-sm font-bold text-foreground tracking-tight">
                {monthGroup}
              </h3>

              <div className="space-y-3">
                {monthShoots.map((shoot) => {
                  const { day, dayNum } = getDayInfo(shoot.date)
                  const crewAssigned = shoot.members && shoot.members.length > 0
                  return (
                    <div
                      key={shoot.id}
                      className="group relative rounded-xl border border-border bg-card p-4 transition-all hover:border-primary/40 hover:shadow-xs"
                    >
                      <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
                        {/* Left & Center section */}
                        <div className="flex items-start sm:items-center gap-4 sm:gap-6 flex-1 min-w-0">
                          {/* Calendar Day & Date Box */}
                          <div className="flex flex-col items-center justify-center rounded-xl bg-muted/40 border border-border/80 px-3.5 py-2.5 min-w-[76px] text-center shadow-2xs group-hover:bg-primary/5 transition">
                            <span className="text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                              {day}
                            </span>
                            <span className="text-2xl font-black tabular-nums text-foreground leading-tight mt-0.5">
                              {dayNum}
                            </span>
                          </div>

                          {/* Shoot Meta Information */}
                          <div className="space-y-1.5 flex-1 min-w-0">
                            {/* Time & Venue & Location row */}
                            <div className="flex items-center gap-4 text-xs text-muted-foreground flex-wrap">
                              <span className="inline-flex items-center gap-1 font-medium text-foreground">
                                <IconClock className="h-3.5 w-3.5 text-muted-foreground" />
                                {shoot.timeRange || shoot.startTime || 'Time TBD'}
                              </span>

                              <span className="inline-flex items-center gap-1">
                                <IconMapPin className="h-3.5 w-3.5 text-muted-foreground" />
                                {shoot.venue || 'To be decided'}
                              </span>

                              {shoot.city && (
                                <span className="inline-flex items-center gap-1 font-medium text-foreground">
                                  <IconMapPin className="h-3.5 w-3.5 text-primary" />
                                  {shoot.city}
                                </span>
                              )}
                            </div>

                            {/* Title & Project Name */}
                            <div className="flex items-center gap-2 flex-wrap">
                              <h4 className="text-sm font-bold text-foreground tracking-tight hover:text-primary transition cursor-pointer">
                                {shoot.projectName || 'Shoot'}
                                {shoot.title ? ` - ${shoot.title}` : ''}
                              </h4>
                              {shoot.projectType && (
                                <Badge variant="outline" className="text-[10px] font-semibold">
                                  {shoot.projectType}
                                </Badge>
                              )}
                            </div>

                            {/* Assign Members & Crew chips */}
                            <div className="flex items-center gap-2 pt-1 flex-wrap">
                              <Button
                                variant="secondary"
                                size="sm"
                                onClick={() => handleOpenAssignDrawer(shoot)}
                                className="h-7 text-xs px-3 rounded-md bg-muted/80 hover:bg-muted font-medium text-foreground border border-border/60"
                              >
                                <IconUsers className="h-3 w-3 mr-1 text-muted-foreground" />
                                {crewAssigned ? 'Manage Crew' : 'Assign Members'}
                              </Button>

                              {crewAssigned ? (
                                <div className="flex items-center gap-1.5 flex-wrap">
                                  {shoot.members.map((m, idx) => (
                                    <Badge
                                      key={`${m}-${idx}`}
                                      variant="outline"
                                      className="text-[10px] font-medium bg-background/80 text-foreground border-border/80"
                                    >
                                      {m}
                                    </Badge>
                                  ))}
                                </div>
                              ) : (
                                <span className="text-[11px] text-amber-600 font-medium">
                                  No crew assigned yet
                                </span>
                              )}
                            </div>
                          </div>
                        </div>

                        {/* Right Section: Status Pill & Activity link */}
                        <div className="flex lg:flex-col items-center lg:items-end justify-between gap-3 shrink-0">
                          <Select
                            value={shoot.status || 'New'}
                            onValueChange={(val) => onStatusChange?.(shoot.id, val as SessionStatus)}
                          >
                            <SelectTrigger className="h-8 text-xs font-semibold rounded-md border-primary/20 bg-primary/5 text-primary hover:bg-primary/10 transition px-3 gap-2 w-[110px]">
                              <SelectValue placeholder="Status" />
                            </SelectTrigger>
                            <SelectContent align="end">
                              <SelectItem value="New" className="text-xs">New</SelectItem>
                              <SelectItem value="Scheduled" className="text-xs">Scheduled</SelectItem>
                              <SelectItem value="Completed" className="text-xs">Completed</SelectItem>
                              <SelectItem value="Delayed" className="text-xs">Delayed</SelectItem>
                            </SelectContent>
                          </Select>

                          <button
                            type="button"
                            onClick={() => toast.info(`Shoots Activity count: ${shoot.shootsActivityCount || 0}`)}
                            className="inline-flex items-center gap-1.5 text-xs text-primary font-semibold hover:underline cursor-pointer"
                          >
                            <IconRotateClockwise2 className="h-3.5 w-3.5 text-primary" />
                            <span>Shoots Activity</span>
                          </button>
                        </div>
                      </div>
                    </div>
                  )
                })}
              </div>
            </div>
          ))
        )}
      </div>

      {/* Crew & Session Drawer */}
      <CreateSessionDrawer
        open={isDrawerOpen}
        onOpenChange={(op) => {
          setIsDrawerOpen(op)
          if (!op) setSelectedSessionToEdit(null)
        }}
        onSave={(saved: ProjectSession) => {
          onSaveSession(saved)
          setIsDrawerOpen(false)
          setSelectedSessionToEdit(null)
          toast.success('Shoot crew updated successfully')
        }}
        initialSession={selectedSessionToEdit}
        existingProjectNames={existingProjectNames}
      />
    </div>
  )
}
