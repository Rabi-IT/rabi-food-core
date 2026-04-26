const API_URL = process.env.API_URL ?? "http://localhost:3000"

type SearchParams = Record<string, string | undefined>

function buildUrl(path: string, params?: SearchParams): string {
  const url = new URL(path, API_URL)
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      if (v !== undefined) url.searchParams.set(k, v)
    }
  }
  return url.toString()
}

export function apiClient(tenantId?: string) {
  const headers: Record<string, string> = { "Content-Type": "application/json" }
  if (tenantId) headers["X-Tenant-ID"] = tenantId

  return {
    async get<T>(path: string, params?: SearchParams, token?: string): Promise<T> {
      const h: Record<string, string> = { ...headers }
      if (token) h["Authorization"] = `Bearer ${token}`
      const r = await fetch(buildUrl(path, params), { headers: h })
      if (!r.ok) throw new Response(r.statusText, { status: r.status })
      return r.json()
    },
    async post<T>(path: string, body: unknown, token?: string): Promise<T> {
      return this._send<T>("POST", path, body, token)
    },
    async patch<T>(path: string, body: unknown, token?: string): Promise<T> {
      return this._send<T>("PATCH", path, body, token)
    },
    async _send<T>(method: string, path: string, body: unknown, token?: string): Promise<T> {
      const h: Record<string, string> = { ...headers }
      if (token) h["Authorization"] = `Bearer ${token}`
      const r = await fetch(buildUrl(path), {
        method,
        headers: h,
        body: JSON.stringify(body),
      })
      if (r.ok && r.headers.get("Content-Type")?.includes("application/json")) {
        return r.json()
      } else if (r.ok) {
        return {} as T
      }

      if (r.status === 400) {
        const error = await r.json()
        throw new Response(error.message || r.statusText, { status: 400 })
      }

      throw new Response(r.statusText, { status: r.status })
    },
  }
}
