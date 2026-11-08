import React, { createContext, useContext, useEffect, useState, useCallback } from "react"
import { apiClient } from "../lib/api-client"

export interface User {
  id: string
  first_name: string
  last_name: string
  email: string
  role_name?: string
  status?: string
}

export interface AuthContextType {
  user: User | null
  loading: boolean
  error: string | null
  login: (email: string, password: string) => Promise<void>
  signup: (firstName: string, lastName: string, email: string, password: string) => Promise<void>
  logout: () => Promise<void>
  refreshUser: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({
  children,
  onUnauthenticated,
}: {
  children: React.ReactNode
  onUnauthenticated?: () => void
}) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState<boolean>(true)
  const [error, setError] = useState<string | null>(null)

  const refreshUser = useCallback(async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await apiClient<{ user: User }>("/api/auth/me")
      setUser(data.user)
    } catch (err: any) {
      setUser(null)
      if (err.status !== 401) {
        setError(err.message || "Failed to load session")
      }
      if (err.status === 401 && onUnauthenticated) {
        onUnauthenticated()
      }
    } finally {
      setLoading(false)
    }
  }, [onUnauthenticated])

  useEffect(() => {
    refreshUser()
  }, [refreshUser])

  const login = async (email: string, password: string) => {
    setError(null)
    const data = await apiClient<{ message: string; user: User }>("/api/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    })
    setUser(data.user)
  }

  const signup = async (firstName: string, lastName: string, email: string, password: string) => {
    setError(null)
    const data = await apiClient<{ message: string; user: User }>("/api/auth/signup", {
      method: "POST",
      body: JSON.stringify({
        first_name: firstName,
        last_name: lastName,
        email,
        password,
      }),
    })
    setUser(data.user)
  }

  const logout = async () => {
    try {
      await apiClient("/api/auth/logout", { method: "POST" })
    } finally {
      setUser(null)
      if (onUnauthenticated) {
        onUnauthenticated()
      }
    }
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        error,
        login,
        signup,
        logout,
        refreshUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider")
  }
  return context
}

