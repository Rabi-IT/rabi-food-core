import { apiClient } from "~/lib/api.server"

const DASHBOARD_SLUGS = new Set(["app", "www"])

export type Tenant = {
  id: string
  name: string
  slug: string
  language: string
}

export function getTenantSlug(request: Request): string | null {
  const host = request.headers.get("host") ?? ""
  const hostname = host.split(":")[0] // strip port
  const parts = hostname.split(".")

  // bare localhost or app.* subdomain → dashboard
  if (parts.length === 1 || DASHBOARD_SLUGS.has(parts[0])) {
    return null
  }

  return parts[0]
}

export function isDashboard(request: Request): boolean {
  return getTenantSlug(request) === null
}

export async function resolveTenant(request: Request): Promise<Tenant | null> {
  const slug = getTenantSlug(request)
  if (!slug) return null
  const result = await apiClient().get<Tenant>(`/tenant/slug/${slug}`)
  return result.ok ? result.data : null
}
