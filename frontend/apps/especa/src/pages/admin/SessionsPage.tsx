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
  IconPlus,
  IconSearch,
  IconEdit,
  IconTrash,
  IconExternalLink,
  IconUsers,
} from "@repo/ui/icons/tabler-icons"
import type { ProjectSession, SessionStatus } from '../../types/session'
import { PROJECT_TYPES } from '../../types/session'
import { CreateSessionDrawer } from '../../components/CreateSessionDrawer'

interface SessionsPageProps {
  sessions: ProjectSession[]
  onAddNewSession?: () => void
  onEditSession?: (session: ProjectSession) => void
  onDeleteSession: (sessionId: string) => void
  onStatusChange?: (sessionId: string, newStatus: SessionStatus) => void
  onSaveSession: (session: ProjectSession) => void
  existingProjectNames?: string[]
}

type TabType = 'Upcoming' | 'Past' | 'TBD'

export function SessionsPage({
  sessions,
  onDeleteSession,
  onSaveSession,
  existingProjectNames = [],
}: SessionsPageProps) {
  // Pill tab filter (Upcoming / Past / TBD) as seen in Screenshot 1
  const [activeTab, setActiveTab] = React.useState<TabType>('Upcoming')
  const [sortOption, setSortOption] = React.useState<'recently_created' | 'date_asc' | 'date_desc'>('recently_created')
  const [searchQuery, setSearchQuery] = React.useState('')
  const [isFilterOpen, setIsFilterOpen] = React.useState(false)
  const [typeFilter, setTypeFilter] = React.useState<string>('All')

  // Create/Edit Drawer state
  const [isDrawerOpen, setIsDrawerOpen] = React.useState(false)
  const [selectedSessionToEdit, setSelectedSessionToEdit] = React.useState<ProjectSession | null>(null)

  const handleOpenCreateDrawer = () => {
    setSelectedSessionToEdit(null)
    setIsDrawerOpen(true)
  }

  const handleOpenEditDrawer = (session: ProjectSession) => {
    setSelectedSessionToEdit(session)
    setIsDrawerOpen(true)
  }

  // Filter sessions by tab, search, type
  const filteredSessions = React.useMemo(() => {
    const todayStr = '2026-11-09'

    return sessions.filter((s) => {
      // Search
      if (searchQuery.trim()) {
        const q = searchQuery.toLowerCase()
        const matchesTitle = s.title.toLowerCase().includes(q)
        const matchesVenue = (s.venue || '').toLowerCase().includes(q)
        const matchesCity = (s.city || '').toLowerCase().includes(q)
        const matchesCrew = (s.members || []).some((m) => m.toLowerCase().includes(q))
        if (!matchesTitle && !matchesVenue && !matchesCity && !matchesCrew) return false
      }

      // Project Type filter
      if (typeFilter !== 'All' && s.projectType !== typeFilter) {
        return false
      }

      // Tab filter
      if (activeTab === 'TBD') {
        return !s.date || s.date === 'TBD'
      }

      // Extract comparable date
      const sessionDateStr = s.date.includes('-')
        ? s.date
        : s.date.includes('Nov 2026')
          ? `2026-11-${(s.date.match(/\d+/) || ['11'])[0].padStart(2, '0')}`
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
      return b.id.localeCompare(a.id) // recently created
    })
  }, [sessions, activeTab, sortOption, searchQuery, typeFilter])

  // Group sessions by month/year heading (e.g. "November, 2026")
  const groupedSessions = React.useMemo(() => {
    const groups: { [key: string]: ProjectSession[] } = {}

    filteredSessions.forEach((ses) => {
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
  }, [filteredSessions])

  // Parse day abbreviation and date number for left tile
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

  return (
    <div className="space-y-6">
      {/* Header bar matches Screenshot 1 with modern UX additions */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        {/* Pill-shaped Tabs: Upcoming | Past | TBD */}
        <div className="inline-flex items-center rounded-full bg-muted/50 p-1 border shadow-2xs">
          <button
            type="button"
            onClick={() => setActiveTab('Upcoming')}
            className={cn(
              'inline-flex items-center gap-1.5 rounded-full px-4 py-1.5 text-xs font-semibold transition',
              activeTab === 'Upcoming'
                ? 'bg-card text-foreground shadow-2xs'
                : 'text-muted-foreground hover:text-foreground'
            )}
          >
            <IconCalendar className="h-3.5 w-3.5" />
            <span>Upcoming</span>
          </button>
          <button
            type="button"
            onClick={() => setActiveTab('Past')}
            className={cn(
              'inline-flex items-center gap-1.5 rounded-full px-4 py-1.5 text-xs font-semibold transition',
              activeTab === 'Past'
                ? 'bg-card text-foreground shadow-2xs'
                : 'text-muted-foreground hover:text-foreground'
            )}
          >
            <IconClock className="h-3.5 w-3.5" />
            <span>Past</span>
          </button>
          <button
            type="button"
            onClick={() => setActiveTab('TBD')}
            className={cn(
              'inline-flex items-center gap-1.5 rounded-full px-4 py-1.5 text-xs font-semibold transition',
              activeTab === 'TBD'
                ? 'bg-card text-foreground shadow-2xs'
                : 'text-muted-foreground hover:text-foreground'
            )}
          >
            <IconCalendar className="h-3.5 w-3.5" />
            <span>TBD</span>
          </button>
        </div>

        {/* Right Action Tools: Sort/Recently Created | Filter | New Session Button */}
        <div className="flex items-center gap-2.5 flex-wrap">
          {/* Recently Created dropdown sort */}
          <Select
            value={sortOption}
            onValueChange={(val) => setSortOption(val as any)}
          >
            <SelectTrigger className="h-9 text-xs bg-card border-border shadow-2xs w-[155px]">
              <div className="flex items-center gap-1.5 text-muted-foreground">
                <IconClock className="h-3.5 w-3.5" />
                <span className="font-semibold text-foreground text-xs">
                  {sortOption === 'recently_created'
                    ? 'Recently Created'
                    : sortOption === 'date_asc'
                      ? 'Date (Earliest)'
                      : 'Date (Latest)'}
                </span>
              </div>
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="recently_created" className="text-xs">
                Recently Created
              </SelectItem>
              <SelectItem value="date_asc" className="text-xs">
                Date (Earliest first)
              </SelectItem>
              <SelectItem value="date_desc" className="text-xs">
                Date (Latest first)
              </SelectItem>
            </SelectContent>
          </Select>

          {/* Filters Button */}
          <Button
            variant="outline"
            size="sm"
            onClick={() => setIsFilterOpen(!isFilterOpen)}
            className={cn(
              'h-9 text-xs gap-1.5 shadow-2xs transition',
              isFilterOpen && 'bg-accent text-accent-foreground border-accent'
            )}
          >
            <IconFilter className="h-3.5 w-3.5 text-muted-foreground" />
            <span>Filters</span>
            {typeFilter !== 'All' && (
              <span className="h-2 w-2 rounded-full bg-emerald-500" />
            )}
          </Button>

          {/* New Session Button (Green accent matching Screenshot 1 & 2) */}
          <Button
            onClick={handleOpenCreateDrawer}
            className="h-9 text-xs font-semibold bg-emerald-600 hover:bg-emerald-700 text-white gap-1.5 shadow-2xs px-4"
          >
            <IconPlus className="h-4 w-4" />
            <span>New Session</span>
          </Button>
        </div>
      </div>

      {/* Expandable Filters Bar */}
      {isFilterOpen && (
        <Card className="p-3 shadow-2xs border bg-card animate-in fade-in-50 duration-200">
          <div className="flex flex-col sm:flex-row items-center gap-3">
            <div className="relative flex-1 w-full">
              <IconSearch className="absolute left-2.5 top-2.5 h-3.5 w-3.5 text-muted-foreground" />
              <Input
                placeholder="Search sessions by project, client, venue, or crew..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-8 h-9 text-xs"
              />
            </div>
            <div className="w-full sm:w-[220px]">
              <Select value={typeFilter} onValueChange={setTypeFilter}>
                <SelectTrigger className="h-9 text-xs">
                  <SelectValue placeholder="All Project Types" />
                </SelectTrigger>
                <SelectContent className="max-h-60">
                  <SelectItem value="All" className="text-xs">
                    All Project Types
                  </SelectItem>
                  {PROJECT_TYPES.map((pt) => (
                    <SelectItem key={pt} value={pt} className="text-xs">
                      {pt}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            {(searchQuery || typeFilter !== 'All') && (
              <Button
                variant="ghost"
                size="sm"
                className="text-xs h-9 px-3"
                onClick={() => {
                  setSearchQuery('')
                  setTypeFilter('All')
                }}
              >
                Reset
              </Button>
            )}
          </div>
        </Card>
      )}

      {/* Sessions Timeline List by Group */}
      <div className="space-y-6">
        {Object.keys(groupedSessions).length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-center border border-dashed rounded-2xl bg-card">
            <IconCalendar className="h-10 w-10 text-muted-foreground/30 mb-3" />
            <h3 className="text-sm font-semibold text-foreground">No sessions found</h3>
            <p className="text-xs text-muted-foreground max-w-sm mt-1 mb-4">
              {searchQuery || typeFilter !== 'All'
                ? 'No sessions match your filter criteria.'
                : 'There are no sessions in this view. Schedule your first shoot session.'}
            </p>
            <Button
              onClick={handleOpenCreateDrawer}
              size="sm"
              className="bg-emerald-600 hover:bg-emerald-700 text-white text-xs h-8 gap-1.5"
            >
              <IconPlus className="h-3.5 w-3.5" />
              <span>Schedule Session</span>
            </Button>
          </div>
        ) : (
          Object.entries(groupedSessions).map(([monthGroup, monthSessions]) => (
            <div key={monthGroup} className="space-y-3">
              {/* Month Header (Screenshot 1: "November, 2026") */}
              <h3 className="text-sm font-bold text-foreground tracking-tight">
                {monthGroup}
              </h3>

              {/* Session Cards */}
              <div className="space-y-3">
                {monthSessions.map((session) => {
                  const { day, dayNum } = getDayInfo(session.date)
                  return (
                    <div
                      key={session.id}
                      className="group relative rounded-xl border border-border bg-card p-4 transition-all hover:border-primary/40 hover:shadow-xs"
                    >
                      <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
                        {/* Left & Center section */}
                        <div className="flex items-start sm:items-center gap-4 sm:gap-6 flex-1">
                          {/* Calendar Day & Date Box (Screenshot 1: "Wed 11", "Sat 14") */}
                          <div className="flex flex-col items-center justify-center rounded-xl bg-muted/40 border border-border/80 px-3.5 py-2.5 min-w-[76px] text-center shadow-2xs group-hover:bg-primary/5 transition">
                            <span className="text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                              {day}
                            </span>
                            <span className="text-2xl font-black tabular-nums text-foreground leading-tight mt-0.5">
                              {dayNum}
                            </span>
                          </div>

                          {/* Session Details */}
                          <div className="space-y-1.5 flex-1">
                            {/* Time & Venue & Location row */}
                            <div className="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-muted-foreground font-medium">
                              <span className="flex items-center gap-1.5 text-foreground/80">
                                <IconClock className="h-3.5 w-3.5 text-muted-foreground" />
                                {session.timeRange || session.startTime || 'Time TBD'}
                              </span>
                              {session.venue && (
                                <span className="flex items-center gap-1.5">
                                  <IconMapPin className="h-3.5 w-3.5 text-muted-foreground" />
                                  {session.venue}
                                </span>
                              )}
                              {session.city && (
                                <span className="flex items-center gap-1 text-muted-foreground/80">
                                  <span>•</span>
                                  <span>{session.city}</span>
                                </span>
                              )}
                              {session.projectType && (
                                <Badge variant="secondary" className="text-[10px] h-4.5 px-1.5">
                                  {session.projectType}
                                </Badge>
                              )}
                              {session.shootTiming && (
                                <span className="text-[11px] text-muted-foreground/70">
                                  • {session.shootTiming}
                                </span>
                              )}
                            </div>

                            {/* Title (e.g. Kavya Rajendran Wedding - Reception | Sunset) */}
                            <h4
                              onClick={() => handleOpenEditDrawer(session)}
                              className="font-bold text-sm text-foreground hover:text-primary transition cursor-pointer"
                            >
                              {session.title || session.projectName}
                            </h4>

                            {/* Assigned Crew or "No Members Added" */}
                            <div className="flex items-center gap-2">
                              {session.members && session.members.length > 0 ? (
                                <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                                  <IconUsers className="h-3.5 w-3.5 text-emerald-600" />
                                  <span className="font-medium text-foreground">
                                    Crew: {session.members.join(', ')}
                                  </span>
                                </div>
                              ) : (
                                <span className="text-xs text-muted-foreground/70 font-medium">
                                  No Members Added
                                </span>
                              )}
                              {session.locationLink && (
                                <a
                                  href={session.locationLink}
                                  target="_blank"
                                  rel="noopener noreferrer"
                                  className="text-[11px] text-primary hover:underline flex items-center gap-0.5 ml-2"
                                  onClick={(e) => e.stopPropagation()}
                                >
                                  <span>Map</span>
                                  <IconExternalLink className="h-2.5 w-2.5" />
                                </a>
                              )}
                            </div>
                          </div>
                        </div>

                        {/* Right: Status Pill & Action Buttons */}
                        <div className="flex items-center gap-3 justify-between lg:justify-end border-t lg:border-t-0 pt-2 lg:pt-0">
                          {/* Status Badge */}
                          <Badge
                            variant="outline"
                            className={cn(
                              'text-xs font-semibold px-3 py-1 rounded-md border',
                              session.status === 'New' && 'bg-blue-50 text-blue-700 border-blue-200 dark:bg-blue-950/40 dark:text-blue-300 dark:border-blue-800',
                              session.status === 'Scheduled' && 'bg-amber-50 text-amber-700 border-amber-200 dark:bg-amber-950/40 dark:text-amber-300 dark:border-amber-800',
                              session.status === 'Completed' && 'bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-950/40 dark:text-emerald-300 dark:border-emerald-800',
                              session.status === 'Delayed' && 'bg-rose-50 text-rose-700 border-rose-200 dark:bg-rose-950/40 dark:text-rose-300 dark:border-rose-800'
                            )}
                          >
                            {session.status}
                          </Badge>

                          {/* Quick Edit button */}
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleOpenEditDrawer(session)}
                            className="h-8 w-8 p-0 text-muted-foreground hover:text-foreground"
                            title="Edit Session"
                          >
                            <IconEdit className="h-4 w-4" />
                          </Button>

                          {/* Delete */}
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => {
                              if (confirm(`Delete session "${session.title}"?`)) {
                                onDeleteSession(session.id)
                                toast.success('Session deleted')
                              }
                            }}
                            className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive"
                            title="Delete Session"
                          >
                            <IconTrash className="h-4 w-4" />
                          </Button>
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

      {/* Drawer Snap Points for Creating & Editing Sessions */}
      <CreateSessionDrawer
        open={isDrawerOpen}
        onOpenChange={setIsDrawerOpen}
        onSave={(newOrUpdatedSession) => {
          onSaveSession(newOrUpdatedSession)
        }}
        initialSession={selectedSessionToEdit}
        existingProjectNames={existingProjectNames}
      />
    </div>
  )
}
