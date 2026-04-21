import { Outlet, data } from "react-router"
import type { Route } from "./+types/layout"
import { getTenantSlug, resolveTenant } from "~/lib/tenant.server"
import { localeCookie } from "~/i18n.server"

export async function loader({ request }: Route.LoaderArgs) {
  const slug = getTenantSlug(request)

  // No slug means dashboard subdomain — skip storefront, let dashboard layout handle it
  if (!slug) return data({ tenant: null })

  const tenant = await resolveTenant(request)
  if (!tenant) throw new Response("Tenant not found", { status: 404 })

  const headers = new Headers()

  const hasLangCookie = await localeCookie.parse(request.headers.get("Cookie"))
  if (!hasLangCookie) {
    headers.set("Set-Cookie", await localeCookie.serialize(tenant.language))
  }

  return data({ tenant }, { headers })
}

export default function StorefrontLayout() {
  return <Outlet />
}
