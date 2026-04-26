import { useLoaderData, useActionData, useNavigation, useSubmit, useFetcher, Link } from "react-router"
import type { Route } from "./+types/subscribe"
import { resolveTenant } from "~/lib/tenant.server"
import { apiClient } from "~/lib/api.server"
import { isTokenExpired, getStoredRefreshToken, storeRefreshToken, clearRefreshToken } from "~/lib/token"
import { useWizard, type WizardData } from "~/components/subscription/use-wizard"
import { StepProducts } from "~/components/subscription/step-products"
import { StepBasket } from "~/components/subscription/step-basket"
import { StepDelivery } from "~/components/subscription/step-delivery"
import { StepPlan } from "~/components/subscription/step-plan"
import { StepUser } from "~/components/subscription/step-user"
import { StepAddress } from "~/components/subscription/step-address"
import { StepSummary } from "~/components/subscription/step-summary"
import { useEffect, useRef, useState } from "react"
import { cn } from "~/lib/utils"
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
    autoRenew: data.autoRenew,
  }, data.user.token)
  if (!result.ok) return { error: result.code }
  return { ok: true }
}

const STEP_LABELS = {
  1: "Escolha os produtos",
  2: "Quantidades",
  3: "Dias de entrega",
  4: "Seu plano",
  5: "Identificação",
  6: "Endereço de entrega",
  summary: "Resumo",
} as const

const NUMERIC_STEPS = [1, 2, 3, 4, 5, 6] as const

export default function Subscribe() {
  const { products, config, token: loaderToken, profile: loaderProfile } = useLoaderData<typeof loader>()
  const actionData = useActionData<typeof action>()
  const navigation = useNavigation()
  const submit = useSubmit()
  const otpFetcher = useFetcher<OtpActionResult>()
  const addressFetcher = useFetcher<OtpActionResult>()
  const { step, setStep, data, setData, reset } = useWizard()

  const isSubmitting = navigation.state !== "idle"
  const confirmResult = actionData as ConfirmActionResult | undefined

  const [otpSubStep, setOtpSubStep] = useState<"identity" | "code">("identity")
  const [forcedUserMode, setForcedUserMode] = useState<"new" | "existing" | undefined>(undefined)
  const refreshInProgressRef = useRef(false)

  // Pre-populate token and profile from loader (already authenticated)
  useEffect(() => {
    if (loaderToken && !data.user.token) {
      setData({ user: { ...data.user, token: loaderToken } })
    }
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    console.log("[otp] fetcher.data:", otpFetcher.data)
    if (!otpFetcher.data || !("ok" in otpFetcher.data)) return

    if (otpFetcher.data.intent === "send-otp") {
      setOtpSubStep("code")
    }

    if (otpFetcher.data.intent === "verify-otp" && "accessToken" in otpFetcher.data) {
      storeRefreshToken(otpFetcher.data.refreshToken)
      setForcedUserMode(undefined)
      setData({ user: { ...data.user, token: otpFetcher.data.accessToken } })
      setStep(6)
    }
  }, [otpFetcher.data]) // eslint-disable-line react-hooks/exhaustive-deps

  // Fetch profile when entering step 6
  const profileFetcher = useFetcher<OtpActionResult>()
  const prevStep6 = useRef(false)
  useEffect(() => {
    if (step !== 6 || prevStep6.current) {
      if (step !== 6) prevStep6.current = false
      return
    }
    prevStep6.current = true
    ensureToken().then((token) => {
      if (!token) return
      profileFetcher.submit(
        { intent: "get-profile", token },
        { method: "post", encType: "application/json" },
      )
    })
  }, [step]) // eslint-disable-line react-hooks/exhaustive-deps

  const existingProfile =
    profileFetcher.data && "ok" in profileFetcher.data && profileFetcher.data.intent === "get-profile"
      ? profileFetcher.data.profile
      : null

  // Address saved → advance to summary
  useEffect(() => {
    if (
      addressFetcher.data &&
      "ok" in addressFetcher.data &&
      addressFetcher.data.intent === "save-address"
    ) {
      setStep("summary")
    }
  }, [addressFetcher.data]) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (confirmResult && "ok" in confirmResult) reset()
  }, [confirmResult]) // eslint-disable-line react-hooks/exhaustive-deps

  async function ensureToken(): Promise<string | null> {
    const { token } = data.user
    if (!token) return null
    if (!isTokenExpired(token)) return token
    if (refreshInProgressRef.current) return null

    refreshInProgressRef.current = true
    try {
      const refreshToken = getStoredRefreshToken()
      if (!refreshToken) {
        setForcedUserMode("existing")
        setOtpSubStep("identity")
        setStep(5)
        return null
      }

      const res = await fetch(window.location.pathname, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ intent: "refresh", refreshToken }),
      })
      const result = await res.json() as { intent: "refresh"; ok: boolean; accessToken?: string; refreshToken?: string }
      if (result.ok && result.accessToken && result.refreshToken) {
        storeRefreshToken(result.refreshToken)
        setData({ user: { ...data.user, token: result.accessToken } })
        return result.accessToken
      }

      clearRefreshToken()
      setForcedUserMode("existing")
      setOtpSubStep("identity")
      setStep(5)
      return null
    } finally {
      refreshInProgressRef.current = false
    }
  }

  function handleSendOtp(name: string, email: string, phone: string) {
    otpFetcher.submit(
      { intent: "send-otp", name, email, phone },
      { method: "post", encType: "application/json" },
    )
  }

  function handleVerifyOtp(email: string, code: string) {
    otpFetcher.submit(
      { intent: "verify-otp", email, code },
      { method: "post", encType: "application/json" },
    )
  }

  async function handleSaveAddress() {
    const token = await ensureToken()
    if (!token) return
    addressFetcher.submit(
      { intent: "save-address", token, phone: data.user.phone, ...data.address },
      { method: "post", encType: "application/json" },
    )
  }

  function handleConfirm() {
    submit(
      { ...data, intent: "confirm" } as unknown as Record<string, string>,
      { method: "post", encType: "application/json" },
    )
  }

  if (confirmResult && "ok" in confirmResult) {
    return (
      <div className="min-h-screen bg-background flex flex-col items-center justify-center px-4 text-center">
        <div className="text-5xl mb-4">🎉</div>
        <h1 className="text-2xl font-bold mb-2">Assinatura confirmada!</h1>
        <p className="text-muted-foreground mb-8">
          Sua cesta recorrente está ativa. Em breve você receberá um e-mail de confirmação.
        </p>
        <Link to="/" className="text-sm text-primary font-medium hover:underline">
          Voltar à loja
        </Link>
      </div>
    )
  }

  const currentStepNumber = step === "summary" ? null : (step as number)
  const stepLabel = STEP_LABELS[step as keyof typeof STEP_LABELS]

  const otpError =
    otpFetcher.data && "error" in otpFetcher.data ? otpFetcher.data.error : undefined
  const addressError =
    addressFetcher.data && "error" in addressFetcher.data ? addressFetcher.data.error : undefined

  return (
    <div className="min-h-screen bg-background text-foreground">
      <header className="sticky top-0 z-10 bg-background border-b px-4 py-3">
        <div className="flex items-center gap-3">
          {step === 1 ? (
            <Link to="/" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
              ←
            </Link>
          ) : step !== 5 && (step !== 6 || !!loaderToken) ? (
            <button
              type="button"
              onClick={() => {
                if (step === "summary") setStep(6)
                else if (step === 2) setStep(1)
                else if (step === 3) setStep(2)
                else if (step === 4) setStep(3)
                else if (step === 6) setStep(data.user.token && loaderToken ? 4 : 5)
              }}
              className="text-sm text-muted-foreground hover:text-foreground transition-colors"
            >
              ←
            </button>
          ) : null}

          <div className="flex-1">
            <p className="text-xs text-muted-foreground">
              {currentStepNumber ? `Etapa ${currentStepNumber} de 6` : "Resumo"}
            </p>
            <p className="text-sm font-medium">{stepLabel}</p>
          </div>
        </div>

        {currentStepNumber && (
          <div className="mt-2 flex gap-1">
            {NUMERIC_STEPS.map((s) => (
              <div
                key={s}
                className={cn(
                  "h-1 flex-1 rounded-full transition-colors",
                  s <= currentStepNumber ? "bg-primary" : "bg-muted"
                )}
              />
            ))}
          </div>
        )}
      </header>

      <main className="mx-auto max-w-lg px-4 py-6">
        {!config?.isOpen && (
          <div className="rounded-xl border border-destructive/30 bg-destructive/5 p-4 mb-6 text-center">
            <p className="text-sm text-destructive">
              As assinaturas estão temporariamente indisponíveis. Tente novamente mais tarde.
            </p>
          </div>
        )}

        {step === 1 && (
          <StepProducts
            products={products}
            items={data.items}
            onChange={(items) => setData({ items })}
            onNext={() => setStep(2)}
          />
        )}

        {step === 2 && (
          <StepBasket
            products={products}
            items={data.items}
            onChange={(items) => setData({ items })}
            onNext={() => setStep(3)}
          />
        )}

        {step === 3 && (
          <StepDelivery
            deliveryDays={data.deliveryDays}
            onChange={(deliveryDays) => setData({ deliveryDays })}
            onNext={() => setStep(4)}
          />
        )}

        {step === 4 && (
          <StepPlan
            products={products}
            items={data.items}
            totalCycles={Math.max(data.totalCycles, data.deliveryDays.length)}
            minCycles={data.deliveryDays.length}
            deliveryWeekdays={data.deliveryDays.map((d) => d.weekday)}
            cycleDiscountRules={config?.discountRules ?? []}
            autoRenew={data.autoRenew}
            onChange={(totalCycles, autoRenew) => setData({ totalCycles, autoRenew })}
            onNext={() => setStep(data.user.token ? 6 : 5)}
          />
        )}

        {step === 5 && (
          <StepUser
            user={data.user}
            subStep={otpSubStep}
            initialMode={forcedUserMode}
            onChange={(user) => setData({ user })}
            onBack={() => setStep(4)}
            onSendOtp={handleSendOtp}
            onVerifyOtp={handleVerifyOtp}
            onResetOtp={() => setOtpSubStep("identity")}
            isSending={otpFetcher.state !== "idle"}
            isVerifying={false}
            otpError={otpError}
          />
        )}

        {step === 6 && (
          <StepAddress
            phone={data.user.phone}
            address={data.address}
            existingProfile={existingProfile}
            isLoadingProfile={profileFetcher.state !== "idle"}
            onChangePhone={(phone) => setData({ user: { ...data.user, phone } })}
            onChangeAddress={(address) => setData({ address })}
            onNext={handleSaveAddress}
            isSaving={addressFetcher.state !== "idle"}
            saveError={addressError}
          />
        )}

        {step === "summary" && (
          <StepSummary
            data={data}
            products={products}
            cycleDiscountRules={config.discountRules}
            cutoffOffsetMinutes={config?.cutoffOffsetMinutes ?? 0}
            error={confirmResult && "error" in confirmResult ? confirmResult.error : undefined}
            isSubmitting={isSubmitting}
            onConfirm={handleConfirm}
            onBack={() => setStep(6)}
          />
        )}
      </main>
    </div>
  )
}
