import React from "react"
import { useAuth } from "../hooks/use-auth"

export interface ProtectedRouteProps {
  children: React.ReactNode
  fallback?: React.ReactNode
  requiredRole?: string
}

export function ProtectedRoute({
  children,
  fallback,
  requiredRole,
}: ProtectedRouteProps) {
  const { user, loading } = useAuth()

  if (loading) {
    return (
      <div className="flex min-h-[50vh] items-center justify-center">
        <div className="h-6 w-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
      </div>
    )
  }

  if (!user) {
    return fallback ? <>{fallback}</> : null
  }

  if (requiredRole && user.role_name !== requiredRole) {
    return (
      <div className="flex min-h-[50vh] flex-col items-center justify-center gap-2">
        <h2 className="text-xl font-bold">Access Denied</h2>
        <p className="text-muted-foreground">You do not have permission to view this resource.</p>
      </div>
    )
  }

  return <>{children}</>
}

