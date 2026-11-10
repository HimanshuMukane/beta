import * as React from 'react'
import {
  IconPlus,
  IconTrash,
  IconUserPlus,
  IconUsers,
  IconUpload,
  IconPhone,
  IconMail,
  IconBriefcase,
} from '@tabler/icons-react'
import { Card, CardContent } from '@repo/ui/components/ui/card'
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
import { Badge } from '@repo/ui/components/ui/badge'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@repo/ui/components/ui/dialog'
import { toast } from '@repo/ui/components/ui/sonner'
import { apiClient } from '@repo/ui/lib/api-client'
import type { TeamDesignation, CrewMember } from '../../types/settings'
import { settingsApi } from '../../api/settingsApi'
import { PhoneInputField } from '../../components/PhoneInputField'

export function PhotographersPage() {
  const [designations, setDesignations] = React.useState<TeamDesignation[]>([])
  const [crewMembers, setCrewMembers] = React.useState<CrewMember[]>([])

  const loadData = React.useCallback(() => {
    apiClient<CrewMember[]>('/api/crew')
      .then((users) => {
        if (Array.isArray(users)) {
          setCrewMembers(users)
        }
      })
      .catch(() => {
        settingsApi.getSettings().then((s) => {
          if (s.crewMembers) setCrewMembers(s.crewMembers)
        }).catch(console.error)
      })

    apiClient<TeamDesignation[]>('/api/roles')
      .then((roles) => {
        if (Array.isArray(roles)) {
          // Never show admin designation; no crew member can be admin
          const nonAdminRoles = roles.filter(
            (r) => r.name.toLowerCase() !== 'admin'
          )
          setDesignations(nonAdminRoles)
        }
      })
      .catch(console.error)
  }, [])

  React.useEffect(() => {
    loadData()
  }, [loadData])

  const [isDesignationsDrawerOpen, setIsDesignationsDrawerOpen] = React.useState(false)
  const [newDesignationInput, setNewDesignationInput] = React.useState('')

  const [isAddCrewDrawerOpen, setIsAddCrewDrawerOpen] = React.useState(false)
  const [fullName, setFullName] = React.useState('')
  const [countryCode, setCountryCode] = React.useState('+91')
  const [mobileNo, setMobileNo] = React.useState('')
  const [email, setEmail] = React.useState('')
  const [selectedDesignation, setSelectedDesignation] = React.useState('')
  const [note, setNote] = React.useState('')

  React.useEffect(() => {
    if (!selectedDesignation && designations.length > 0) {
      setSelectedDesignation(designations[0]?.name || '')
    }
  }, [designations, selectedDesignation])

  const handleAddDesignation = async () => {
    const name = newDesignationInput.trim()
    if (!name) return
    if (name.toLowerCase() === 'admin') {
      toast.error('Cannot create or use admin designation for crew')
      return
    }
    try {
      const created = await apiClient<TeamDesignation>('/api/roles', {
        method: 'POST',
        body: JSON.stringify({ name }),
      })
      setDesignations((prev) => [...prev, created])
      setNewDesignationInput('')
      toast.success(`Added designation "${created.name}"`)
    } catch {
      toast.error('Failed to add designation')
    }
  }

  const handleDeleteDesignation = async (id: string) => {
    try {
      await apiClient(`/api/roles/${id}`, { method: 'DELETE' })
      setDesignations((prev) => prev.filter((d) => d.id !== id))
      toast.success('Designation removed')
    } catch {
      toast.error('Failed to remove designation')
    }
  }

  // Crew Member handlers
  const handleCreateCrewMember = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!fullName.trim()) {
      toast.error('Full name is required')
      return
    }
    if (!mobileNo.trim()) {
      toast.error('Mobile number is required')
      return
    }

    const matchedRole = designations.find((d) => d.name === selectedDesignation)
    const effectiveDesignation = selectedDesignation || designations[0]?.name || 'Photographer'

    const payload = {
      fullName: fullName.trim(),
      countryCode,
      mobileNo: mobileNo.trim(),
      email: email.trim(),
      designation: effectiveDesignation,
      role_id: matchedRole?.id,
      note: note.trim() || undefined,
    }

    try {
      const created = await apiClient<CrewMember>('/api/crew', {
        method: 'POST',
        body: JSON.stringify(payload),
      })
      setCrewMembers((prev) => [created, ...prev])
      setIsAddCrewDrawerOpen(false)
      setFullName('')
      setMobileNo('')
      setEmail('')
      setNote('')
      toast.success(`Added crew member ${created.fullName}`)
    } catch {
      // Fallback to local settings store if API fails
      const fallbackCrew: CrewMember = {
        id: `crew-${Date.now()}`,
        ...payload,
      }
      const updated = [fallbackCrew, ...crewMembers]
      setCrewMembers(updated)
      settingsApi.updateSettings({ crewMembers: updated }).catch(console.error)
      setIsAddCrewDrawerOpen(false)
      setFullName('')
      setMobileNo('')
      setEmail('')
      setNote('')
      toast.success(`Added crew member ${fallbackCrew.fullName}`)
    }
  }

  const handleDeleteCrewMember = async (id: string) => {
    try {
      await apiClient(`/api/crew/${id}`, { method: 'DELETE' })
      setCrewMembers((prev) => prev.filter((c) => c.id !== id))
      toast.success('Crew member removed')
    } catch {
      const updated = crewMembers.filter((c) => c.id !== id)
      setCrewMembers(updated)
      settingsApi.updateSettings({ crewMembers: updated }).catch(console.error)
      toast.success('Crew member removed')
    }
  }

  return (
    <div className="space-y-6 max-w-5xl mx-auto pb-16">
      {/* Top Banner and Actions */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b pb-4">
        <div className="flex items-center gap-3">
          <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary border border-primary/20">
            <IconUsers className="h-5 w-5" />
          </span>
          <div>
            <h1 className="text-xl font-bold tracking-tight text-foreground">
              Photographers & Studio Crew
            </h1>
            <p className="text-xs text-muted-foreground mt-0.5">
              Manage your camera operators, drone pilots, cinematographers, and designations.
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2.5 self-start sm:self-auto">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setIsDesignationsDrawerOpen(true)}
            className="text-xs font-semibold gap-1.5 h-8"
          >
            <IconBriefcase className="h-3.5 w-3.5" />
            Manage Designations ({designations.length})
          </Button>
          <Button
            size="sm"
            onClick={() => setIsAddCrewDrawerOpen(true)}
            className="text-xs font-semibold gap-1.5 bg-primary text-primary-foreground shadow-xs h-8"
          >
            <IconUserPlus className="h-3.5 w-3.5" />
            Add Crew Member
          </Button>
        </div>
      </div>

      {/* Crew Members Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {crewMembers.map((member) => (
          <Card
            key={member.id}
            className="border border-border/80 shadow-2xs hover:shadow-md hover:border-primary/40 transition rounded-xl relative group"
          >
            <CardContent className="p-4 space-y-3">
              <div className="flex items-start justify-between gap-3">
                <div className="flex items-center gap-3">
                  <div className="h-11 w-11 rounded-full bg-primary/10 text-primary flex items-center justify-center font-bold text-sm shrink-0 border border-primary/20">
                    {member.fullName
                      .split(' ')
                      .map((n) => n[0])
                      .join('')
                      .toUpperCase()
                      .substring(0, 2)}
                  </div>
                  <div className="min-w-0">
                    <h3 className="text-sm font-bold text-foreground truncate">{member.fullName}</h3>
                    <Badge variant="outline" className="text-[10px] font-semibold bg-muted/40 mt-0.5">
                      {member.designation}
                    </Badge>
                  </div>
                </div>

                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => handleDeleteCrewMember(member.id)}
                  className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive opacity-40 group-hover:opacity-100 transition"
                >
                  <IconTrash className="h-3.5 w-3.5" />
                </Button>
              </div>

              <div className="space-y-1.5 pt-1 text-xs text-muted-foreground border-t border-border/50">
                <div className="flex items-center gap-2">
                  <IconPhone className="h-3.5 w-3.5 shrink-0 text-muted-foreground/70" />
                  <span>{member.countryCode} {member.mobileNo}</span>
                </div>
                {member.email && (
                  <div className="flex items-center gap-2">
                    <IconMail className="h-3.5 w-3.5 shrink-0 text-muted-foreground/70" />
                    <span className="truncate">{member.email}</span>
                  </div>
                )}
                {member.note && (
                  <p className="text-[11px] text-muted-foreground/80 italic pt-1 line-clamp-2">
                    "{member.note}"
                  </p>
                )}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* DRAWER 1: Manage Team Designations (Matches Screenshot 1) */}
      <Dialog open={isDesignationsDrawerOpen} onOpenChange={setIsDesignationsDrawerOpen}>
        <DialogContent className="max-w-md">
          <DialogHeader className="border-b pb-3">
            <DialogTitle className="text-lg font-bold text-foreground">
              Manage Team Designations
            </DialogTitle>
          </DialogHeader>

          <div className="py-2 space-y-4">
            <div className="rounded-xl border border-border/80 overflow-hidden">
              <div className="bg-muted/40 px-3 py-2 text-[11px] font-bold uppercase tracking-wider text-muted-foreground border-b">
                Name
              </div>

              <div className="divide-y divide-border/60 max-h-[300px] overflow-y-auto">
                {designations.map((des) => (
                  <div
                    key={des.id}
                    className="flex items-center justify-between px-3.5 py-2.5 hover:bg-muted/15 transition text-xs"
                  >
                    <span className="font-medium text-foreground">{des.name}</span>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handleDeleteDesignation(des.id)}
                      className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive"
                    >
                      <IconTrash className="h-3.5 w-3.5" />
                    </Button>
                  </div>
                ))}
              </div>

              {/* Bottom input with add button matching Screenshot 1 */}
              <div className="p-2.5 border-t bg-muted/20 flex items-center gap-2">
                <Input
                  value={newDesignationInput}
                  onChange={(e) => setNewDesignationInput(e.target.value)}
                  placeholder="Designation Name"
                  className="h-8 text-xs bg-background flex-1"
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                      e.preventDefault()
                      handleAddDesignation()
                    }
                  }}
                />
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={handleAddDesignation}
                  className="h-8 w-8 p-0 rounded-full shrink-0"
                >
                  <IconPlus className="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => setIsDesignationsDrawerOpen(false)}
            >
              Close
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* DRAWER 2: Add Crew Member (Matches Screenshot 2) */}
      <Dialog open={isAddCrewDrawerOpen} onOpenChange={setIsAddCrewDrawerOpen}>
        <DialogContent className="max-w-md">
          <form onSubmit={handleCreateCrewMember}>
            <DialogHeader className="border-b pb-3 flex flex-row items-center justify-between">
              <DialogTitle className="text-lg font-bold text-foreground">
                Add Crew Member
              </DialogTitle>
              <Button
                type="submit"
                size="sm"
                className="bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs h-8 mr-6 shadow-xs"
              >
                Create Member
              </Button>
            </DialogHeader>

            <div className="space-y-4 py-4 text-xs">
              {/* Photo Upload Box (Screenshot 2 Match) */}
              <div className="rounded-xl border border-border/80 bg-muted/20 p-4 flex items-center gap-4">
                <div className="h-12 w-12 rounded-full bg-slate-800 text-white flex items-center justify-center shrink-0">
                  <IconUsers className="h-6 w-6" />
                </div>
                <div className="space-y-1">
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    className="h-7 text-xs gap-1.5 font-medium bg-background"
                  >
                    <IconUpload className="h-3.5 w-3.5" />
                    Upload Photo
                  </Button>
                  <p className="text-[10px] text-muted-foreground">
                    JPG or PNG. Max 1 MB recommended.
                  </p>
                </div>
              </div>

              {/* Full Name */}
              <div className="space-y-1">
                <Label className="text-xs font-semibold text-foreground">Full Name</Label>
                <Input
                  value={fullName}
                  onChange={(e) => setFullName(e.target.value)}
                  placeholder="Enter full name"
                  className="h-9 text-xs"
                  required
                />
              </div>

              {/* Mobile No */}
              <div className="space-y-1">
                <Label className="text-xs font-semibold text-foreground">Mobile No</Label>
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
              <div className="space-y-1">
                <Label className="text-xs font-semibold text-foreground">Email Address</Label>
                <Input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="Enter email address"
                  className="h-9 text-xs"
                />
              </div>

              {/* Designation Select */}
              <div className="space-y-1">
                <Label className="text-xs font-semibold text-foreground">Designation</Label>
                <Select value={selectedDesignation} onValueChange={setSelectedDesignation}>
                  <SelectTrigger className="h-9 text-xs">
                    <SelectValue placeholder="Select designation" />
                  </SelectTrigger>
                  <SelectContent>
                    {designations.map((d) => (
                      <SelectItem key={d.id} value={d.name} className="text-xs">
                        {d.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              {/* Note */}
              <div className="space-y-1">
                <Label className="text-xs font-semibold text-foreground">Note</Label>
                <textarea
                  value={note}
                  onChange={(e) => setNote(e.target.value)}
                  placeholder="Optional note about this crew member"
                  rows={3}
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-xs shadow-xs focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring resize-none"
                />
              </div>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
