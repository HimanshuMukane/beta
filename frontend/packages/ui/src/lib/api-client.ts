// Cookie helper to read non-HttpOnly cookies such as app_csrf
export function getCookie(name: string): string | null {
  if (typeof document === "undefined") return null
  const match = document.cookie.match(new RegExp(`(?:^|; )${name}=([^;]*)`))
  return match ? decodeURIComponent(match[1]) : null
}

export function getCsrfToken(): string | null {
  return getCookie("app_csrf")
}

export interface ApiRequestOptions extends RequestInit {
  params?: Record<string, string | number | boolean | undefined>
}

// apiClient performs fetch requests with credentials="include" and automatic CSRF header injection for mutations
export async function apiClient<T = any>(
  endpoint: string,
  options: ApiRequestOptions = {}
): Promise<T> {
  const { params, headers: customHeaders, ...customOptions } = options

  let url = endpoint
  if (params) {
    const searchParams = new URLSearchParams()
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined) {
        searchParams.append(key, String(value))
      }
    })
    const query = searchParams.toString()
    if (query) {
      url += (url.includes("?") ? "&" : "?") + query
    }
  }

  const headers = new Headers(customHeaders || {})

  // Set Content-Type to application/json by default if a body is present and not FormData
  if (customOptions.body && !(customOptions.body instanceof FormData)) {
    if (!headers.has("Content-Type")) {
      headers.set("Content-Type", "application/json")
    }
  }

  // Double-submit CSRF protection on state-changing HTTP methods
  const method = (customOptions.method || "GET").toUpperCase()
  if (["POST", "PUT", "PATCH", "DELETE"].includes(method)) {
    const csrfToken = getCsrfToken()
    if (csrfToken && !headers.has("X-CSRF-Token")) {
      headers.set("X-CSRF-Token", csrfToken)
    }
  }

  const response = await fetch(url, {
    ...customOptions,
    headers,
    credentials: "include",
  })

  // Parse JSON response body if present
  let data: any = null
  const contentType = response.headers.get("content-type")
  if (contentType && contentType.includes("application/json")) {
    try {
      data = await response.json()
    } catch {
      data = null
    }
  }

  if (!response.ok) {
    let errorMessage = `Request failed with status ${response.status}`
    if (data?.error) {
      if (typeof data.error === "string") {
        errorMessage = data.error
      } else if (typeof data.error === "object" && data.error.message) {
        errorMessage = data.error.message
      }
    } else if (data?.message && typeof data.message === "string") {
      errorMessage = data.message
    }
    const error = new Error(errorMessage) as Error & { status?: number; data?: any }
    error.status = response.status
    error.data = data
    throw error
  }

  return data as T
}

