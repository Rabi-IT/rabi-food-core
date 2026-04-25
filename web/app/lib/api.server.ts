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
    async get<T>(path: string, params?: SearchParams): Promise<T> {
      const url = buildUrl(path, params)
      const r = await fetch(url, { headers })
      if (!r.ok) throw new Response(r.statusText, { status: r.status })
      return r.json()
    },
    async post<T>(path: string, body: unknown, token?: string): Promise<T> {
      const h: Record<string, string> = { ...headers }
      if (token) h["Authorization"] = `Bearer ${token}`
      const r = await fetch(buildUrl(path), {
        method: "POST",
        headers: h,
        body: JSON.stringify(body),
      })
      if (!r.ok) throw new Response(r.statusText, { status: r.status })
      return r.json()
    },
  }
}
