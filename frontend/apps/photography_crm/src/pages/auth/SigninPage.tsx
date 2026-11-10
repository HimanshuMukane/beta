import * as React from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { Button } from '@repo/ui/components/ui/button'
import { Input } from '@repo/ui/components/ui/input'
import { Label } from '@repo/ui/components/ui/label'
import { SplitAuthLayout } from '@repo/ui/components/split-auth-layout'
import { toast } from '@repo/ui/components/ui/sonner'
import { IconCamera, IconLock, IconMail, IconArrowRight } from '@tabler/icons-react'
import { useAuth } from '@repo/ui/hooks/use-auth'

export default function SigninPage() {
  const navigate = useNavigate()
  const { login, user } = useAuth()
  const [email, setEmail] = React.useState('especa')
  const [password, setPassword] = React.useState('especa')
  const [loading, setLoading] = React.useState(false)

  React.useEffect(() => {
    if (user) {
      navigate('/projects', { replace: true })
    }
  }, [user, navigate])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!email.trim() || !password.trim()) {
      toast.error('Email and password are required')
      return
    }

    try {
      setLoading(true)
      await login(email.trim(), password.trim())
      toast.success('Signed in successfully')
      navigate('/projects', { replace: true })
    } catch (err: any) {
      toast.error(err.message || 'Invalid credentials')
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
          "From client lead capture to final wedding gallery delivery, manage your entire photography studio workflow in one place."
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
            Sign in to your studio
          </h1>
          <p className="text-xs text-muted-foreground">
            Enter your credentials to access leads, shoots, and bookings.
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-1.5">
            <Label htmlFor="email" className="text-xs font-semibold text-foreground">
              Username or Email
            </Label>
            <div className="relative">
              <span className="absolute left-3 top-2.5 text-muted-foreground">
                <IconMail className="h-4 w-4" />
              </span>
              <Input
                id="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="especa or admin@especa.com"
                className="h-9 text-xs pl-9"
                required
                autoFocus
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <div className="flex items-center justify-between">
              <Label htmlFor="password" className="text-xs font-semibold text-foreground">
                Password
              </Label>
              <span className="text-[11px] text-muted-foreground">
                Default: <code className="bg-muted px-1 py-0.5 rounded text-[10px]">especa</code>
              </span>
            </div>
            <div className="relative">
              <span className="absolute left-3 top-2.5 text-muted-foreground">
                <IconLock className="h-4 w-4" />
              </span>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Enter password"
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
            {loading ? 'Signing in...' : 'Sign In'}
            <IconArrowRight className="h-3.5 w-3.5" />
          </Button>
        </form>

        <div className="text-center text-xs text-muted-foreground pt-2">
          New studio team member?{' '}
          <Link
            to="/signup"
            className="font-semibold text-primary underline underline-offset-4 hover:text-primary/80"
          >
            Create an account
          </Link>
        </div>
      </div>
    </SplitAuthLayout>
  )
}
