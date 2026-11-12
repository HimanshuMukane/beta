import * as React from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import { Label } from '@repo/ui/components/ui/label'
import { SplitAuthLayout } from '@repo/ui/components/split-auth-layout'
import { toast } from '@repo/ui/components/ui/sonner'
import { IconCamera, IconLock, IconMail, IconUser, IconArrowRight } from "@repo/ui/icons/tabler-icons"
import { useAuth } from '@repo/ui/hooks/use-auth'

export default function SignupPage() {
  const navigate = useNavigate()
  const { signup, user } = useAuth()
  const [firstName, setFirstName] = React.useState('')
  const [lastName, setLastName] = React.useState('')
  const [email, setEmail] = React.useState('')
  const [password, setPassword] = React.useState('')
  const [confirmPassword, setConfirmPassword] = React.useState('')
  const [loading, setLoading] = React.useState(false)

  React.useEffect(() => {
    if (user) {
      navigate('/projects', { replace: true })
    }
  }, [user, navigate])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!firstName.trim() || !lastName.trim()) {
      toast.error('First and last name are required')
      return
    }
    if (!email.includes('@')) {
      toast.error('Please enter a valid email address')
      return
    }
    if (password.length < 8) {
      toast.error('Password must be at least 8 characters long')
      return
    }
    if (password !== confirmPassword) {
      toast.error('Passwords do not match')
      return
    }

    try {
      setLoading(true)
      await signup(firstName.trim(), lastName.trim(), email.trim(), password)
      toast.success('Account created successfully')
      navigate('/projects', { replace: true })
    } catch (err: any) {
      toast.error(err.message || 'Signup failed')
    } finally {
      setLoading(false)
    }
  }

  const leftHero = (
    <div className="flex h-full w-full flex-col justify-between p-12 bg-zinc-950 text-white relative">
      <div className="absolute inset-0 bg-radial from-primary/10 via-transparent to-transparent opacity-40 pointer-events-none" />
      <div className="flex items-center gap-3 z-10">
        <span className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary text-primary-foreground shadow-sm">
          <IconCamera className="h-5 w-5" />
        </span>
        <div>
          <h2 className="text-sm font-bold tracking-tight">especa Studio Suite</h2>
          <p className="text-[11px] text-zinc-400">Photography CRM & Operations</p>
        </div>
      </div>

      <div className="space-y-3 z-10 max-w-md">
        <blockquote className="text-lg font-medium leading-snug text-zinc-100">
          "Join the photography CRM platform built for high-touch portrait and wedding cinema production."
        </blockquote>
        <p className="text-xs text-zinc-400 font-mono">
          especa Photography Studio Suite v2.0
        </p>
      </div>

      <div className="text-[11px] text-zinc-500 z-10">
        Secure Enterprise Studio Portal
      </div>
    </div>
  )

  return (
    <SplitAuthLayout imageSlot={leftHero} appName="especa Photography CRM">
      <div className="space-y-6">
        <div className="space-y-2 text-center sm:text-left">
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            Create your studio account
          </h1>
          <p className="text-xs text-muted-foreground">
            Get started by entering your studio credentials below.
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <Label htmlFor="firstName" className="text-xs font-semibold text-foreground">
                First Name
              </Label>
              <div className="relative">
                <span className="absolute left-3 top-2.5 text-muted-foreground">
                  <IconUser className="h-4 w-4" />
                </span>
                <Input
                  id="firstName"
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                  placeholder="Rahul"
                  className="h-9 text-xs pl-9"
                  required
                />
              </div>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="lastName" className="text-xs font-semibold text-foreground">
                Last Name
              </Label>
              <Input
                id="lastName"
                value={lastName}
                onChange={(e) => setLastName(e.target.value)}
                placeholder="Sharma"
                className="h-9 text-xs"
                required
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="email" className="text-xs font-semibold text-foreground">
              Work Email
            </Label>
            <div className="relative">
              <span className="absolute left-3 top-2.5 text-muted-foreground">
                <IconMail className="h-4 w-4" />
              </span>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="rahul@especa.com"
                className="h-9 text-xs pl-9"
                required
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="password" className="text-xs font-semibold text-foreground">
              Password (min. 8 characters)
            </Label>
            <div className="relative">
              <span className="absolute left-3 top-2.5 text-muted-foreground">
                <IconLock className="h-4 w-4" />
              </span>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Create password"
                className="h-9 text-xs pl-9"
                required
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="confirmPassword" className="text-xs font-semibold text-foreground">
              Confirm Password
            </Label>
            <div className="relative">
              <span className="absolute left-3 top-2.5 text-muted-foreground">
                <IconLock className="h-4 w-4" />
              </span>
              <Input
                id="confirmPassword"
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder="Re-enter password"
                className="h-9 text-xs pl-9"
                required
              />
            </div>
          </div>

          <Button
            type="submit"
            disabled={loading}
            className="w-full h-9 text-xs font-semibold gap-1.5 bg-primary text-primary-foreground shadow-xs"
          >
            {loading ? 'Creating Account...' : 'Register Studio Account'}
            <IconArrowRight className="h-3.5 w-3.5" />
          </Button>
        </form>

        <div className="text-center text-xs text-muted-foreground pt-2">
          Already registered with especa?{' '}
          <Link
            to="/signin"
            className="font-semibold text-primary underline underline-offset-4 hover:text-primary/80"
          >
            Sign in
          </Link>
        </div>
      </div>
    </SplitAuthLayout>
  )
}
