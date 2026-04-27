import type { Route } from "./+types/subscribe"
import { resolveTenant } from "~/lib/tenant.server"
import { apiClient } from "~/lib/api.server"
import type { WizardData } from "~/components/subscription/use-wizard"
import type { Product, SubscriptionConfig } from "~/components/subscription/types"

type UserProfile = {
  street: string; complement: string; neighborhood: string
  city: string; state: string; zip: string; phone: string
}

type OtpActionResult =
  | { intent: "send-otp"; ok: true }
  | { intent: "verify-otp"; ok: true; accessToken: string; refreshToken: string }
  | { intent: "refresh"; ok: true; accessToken: string; refreshToken: string }
  | { intent: "refresh"; ok: false }
  | { intent: "get-profile"; ok: true; profile: UserProfile }
  | { intent: "save-address"; ok: true }
  | { intent: string; error: string }

type ConfirmActionResult = { ok: true } | { error: string }

function parseAuthCookie(request: Request): string | null {
  const cookie = request.headers.get("Cookie") ?? ""
  return cookie.match(/auth_token=([^;]+)/)?.[1] ?? null
}

export async function loader({ request }: Route.LoaderArgs) {
  const tenant = await resolveTenant(request)
  if (!tenant) throw new Response("Not found", { status: 404 })

  const api = apiClient(tenant.id)
  const token = parseAuthCookie(request)

  const [productsRes, configRes, profileRes] = await Promise.all([
    api.get<{ data: readonly Product[] }>("/product/", { isActive: "1" }),
    api.get<SubscriptionConfig>("/subscription/config"),
    token ? api.get<UserProfile>("/user/me", undefined, token) : Promise.resolve(null),
  ])

  if (!productsRes.ok) throw new Response("Not found", { status: 404 })
  if (!configRes.ok) throw new Response("Not found", { status: 404 })

  return {
    products: productsRes.data.data,
    config: configRes.data,
    token,
    profile: profileRes?.ok ? profileRes.data : null,
  }
}

export async function action({ request }: Route.ActionArgs): Promise<OtpActionResult | ConfirmActionResult | Response> {
  const tenant = await resolveTenant(request)
  if (!tenant) return { intent: "unknown", error: "UNKNOWN" }

  const api = apiClient(tenant.id)
  const body = await request.json()
  const intent = body.intent as string | undefined

  if (intent === "send-otp") {
    const result = await api.post("/auth/otp", { email: body.email, name: body.name, phone: body.phone })
    if (!result.ok) return { intent: "send-otp", error: result.code }
    return { intent: "send-otp", ok: true }
  }

  if (intent === "verify-otp") {
    const result = await api.post<{ accessToken: string; refreshToken: string }>("/auth/otp/verify", {
      email: body.email,
      code: body.code,
    })
    if (!result.ok) return { intent: "verify-otp", error: result.code }
    const payload: OtpActionResult = { intent: "verify-otp", ok: true, accessToken: result.data.accessToken, refreshToken: result.data.refreshToken }
    return new Response(JSON.stringify(payload), {
      headers: {
        "Content-Type": "application/json",
        "Set-Cookie": `auth_token=${result.data.accessToken}; Path=/; HttpOnly; SameSite=Lax; Max-Age=3600`,
      },
    })
  }

  if (intent === "refresh") {
    const result = await api.post<{ accessToken: string; refreshToken: string }>("/auth/refresh", {
      refreshToken: body.refreshToken,
    })
    if (!result.ok) return { intent: "refresh", ok: false }
    return { intent: "refresh", ok: true, accessToken: result.data.accessToken, refreshToken: result.data.refreshToken }
  }

  if (intent === "get-profile") {
    const result = await api.get<UserProfile>("/user/me", undefined, body.token)
    const profile = result.ok
      ? result.data
      : { street: "", complement: "", neighborhood: "", city: "", state: "", zip: "", phone: "" }
    return { intent: "get-profile", ok: true, profile }
  }

  if (intent === "save-address") {
    const result = await api.patch("/user/me", {
      phone: body.phone,
      street: body.street,
      complement: body.complement,
      neighborhood: body.neighborhood,
      city: body.city,
      state: body.state,
      zip: body.zip,
    }, body.token)
    if (!result.ok) return { intent: "save-address", error: result.code }
    return { intent: "save-address", ok: true }
  }

  // intent === "confirm": create subscription
  const data = body as WizardData
  const result = await api.post("/subscription/", {
    items: data.items,
    deliveryDays: data.deliveryDays,
    totalCycles: data.totalCycles,
  }, data.user.token)
  if (!result.ok) return { error: result.code }
  return { ok: true }
}
