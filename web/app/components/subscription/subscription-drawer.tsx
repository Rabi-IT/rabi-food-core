import { useEffect, useRef, useState } from "react"
import { useFetcher } from "react-router"
// loaderFetcher removed — data is passed as props from the storefront loader
import { Sheet } from "~/components/ui/sheet"
import { cn } from "~/lib/utils"
import { isTokenExpired, getStoredRefreshToken, storeRefreshToken, clearRefreshToken } from "~/lib/token"
import { useWizard } from "./use-wizard"
import { StepProducts } from "./step-products"
import { StepBasket } from "./step-basket"
import { StepDelivery } from "./step-delivery"
import { StepPlan } from "./step-plan"
import { StepUser } from "./step-user"
import { StepAddress } from "./step-address"
import { StepSummary } from "./step-summary"
import type { Product, SubscriptionConfig } from "./types"

type UserProfile = {
  readonly street: string
  readonly complement: string
  readonly neighborhood: string
  readonly city: string
  readonly state: string
  readonly zip: string
  readonly phone: string
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

type Props = {
  readonly open: boolean
  readonly onClose: () => void
  readonly initialProductId?: string | null
  readonly products: readonly Product[]
  readonly config: SubscriptionConfig | null
  readonly serverToken: string | null
}

export default function SubscriptionDrawer({ open, onClose, initialProductId, products, config, serverToken }: Props) {
  const loaderToken = serverToken

  const otpFetcher = useFetcher<OtpActionResult>()
  const addressFetcher = useFetcher<OtpActionResult>()
  const profileFetcher = useFetcher<OtpActionResult>()
  const confirmFetcher = useFetcher<ConfirmActionResult>()

  const { step, setStep, data, setData, reset } = useWizard()
  const [otpSubStep, setOtpSubStep] = useState<"identity" | "code">("identity")
  const [forcedUserMode, setForcedUserMode] = useState<"new" | "existing" | undefined>(undefined)
  const refreshInProgressRef = useRef(false)

  useEffect(() => {
    if (!open || !initialProductId) return
    setData({ items: [{ productId: initialProductId, quantity: 1 }] })
    setStep(2)
  }, [open, initialProductId]) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (!loaderToken || data.user.token) return
    setData({ user: { ...data.user, token: loaderToken } })
  }, [loaderToken]) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (open) return
    reset()
  }, [open]) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
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
        { method: "post", action: "/subscribe", encType: "application/json" },
      )
    })
  }, [step]) // eslint-disable-line react-hooks/exhaustive-deps

  const existingProfile =
    profileFetcher.data && "ok" in profileFetcher.data && profileFetcher.data.intent === "get-profile"
      ? profileFetcher.data.profile
      : null

  useEffect(() => {
    if (
      addressFetcher.data &&
      "ok" in addressFetcher.data &&
      addressFetcher.data.intent === "save-address"
    ) {
      setStep("summary")
    }
  }, [addressFetcher.data]) // eslint-disable-line react-hooks/exhaustive-deps

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

      const res = await fetch("/subscribe", {
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
      { method: "post", action: "/subscribe", encType: "application/json" },
    )
  }

  function handleVerifyOtp(email: string, code: string) {
    otpFetcher.submit(
      { intent: "verify-otp", email, code },
      { method: "post", action: "/subscribe", encType: "application/json" },
    )
  }

  async function handleSaveAddress() {
    const token = await ensureToken()
    if (!token) return
    addressFetcher.submit(
      { intent: "save-address", token, phone: data.user.phone, ...data.address },
      { method: "post", action: "/subscribe", encType: "application/json" },
    )
  }

  function handleConfirm() {
    confirmFetcher.submit(
      { ...data, intent: "confirm" } as unknown as Record<string, string>,
      { method: "post", action: "/subscribe", encType: "application/json" },
    )
  }

  function handleBack() {
    if (step === "summary") return setStep(6)
    if (step === 2) return setStep(1)
    if (step === 3) return setStep(2)
    if (step === 4) return setStep(3)
    if (step === 6) return setStep(data.user.token && loaderToken ? 4 : 5)
  }

  function renderBackButton() {
    if (step === 1) return null
    if (step === 5) return null
    if (step === 6 && !loaderToken) return null
    return (
      <button
        type="button"
        onClick={handleBack}
        className="text-sm text-muted-foreground hover:text-foreground transition-colors"
      >
        ←
      </button>
    )
  }

  const isSubmitting = confirmFetcher.state !== "idle"
  const confirmResult = confirmFetcher.data as ConfirmActionResult | undefined
  const isSuccess = !!(confirmResult && "ok" in confirmResult)

  const otpError = otpFetcher.data && "error" in otpFetcher.data ? otpFetcher.data.error : undefined
  const addressError = addressFetcher.data && "error" in addressFetcher.data ? addressFetcher.data.error : undefined
  const currentStepNumber = step === "summary" ? null : (step as number)
  const stepLabel = STEP_LABELS[step as keyof typeof STEP_LABELS]

  return (
    <Sheet open={open} onClose={onClose} variant="side">
      {isSuccess && (
        <div className="flex flex-col items-center justify-center h-full px-6 text-center">
          <div className="text-5xl mb-4">🎉</div>
          <h1 className="text-2xl font-bold mb-2">Assinatura confirmada!</h1>
          <p className="text-muted-foreground mb-8">
            Sua cesta recorrente está ativa. Em breve você receberá um e-mail de confirmação.
          </p>
          <button
            type="button"
            onClick={onClose}
            className="text-sm text-primary font-medium hover:underline"
          >
            Fechar
          </button>
        </div>
      )}

      {!isSuccess && (
        <div className="flex flex-col h-full">
          <header className="shrink-0 border-b px-4 py-3">
            <div className="flex items-center gap-3">
              {renderBackButton()}

              <div className="flex-1 min-w-0">
                <p className="text-xs text-muted-foreground">
                  {currentStepNumber ? `Etapa ${currentStepNumber} de 6` : "Resumo"}
                </p>
                <p className="text-sm font-medium truncate">{stepLabel}</p>
              </div>

              <button
                type="button"
                onClick={onClose}
                className="shrink-0 text-muted-foreground hover:text-foreground transition-colors text-lg leading-none"
              >
                ✕
              </button>
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

          {!config?.isOpen && (
            <div className="mx-4 mt-4 shrink-0 rounded-xl border border-destructive/30 bg-destructive/5 p-4 text-center">
              <p className="text-sm text-destructive">
                As assinaturas estão temporariamente indisponíveis. Tente novamente mais tarde.
              </p>
            </div>
          )}

          <div className="flex-1 overflow-y-auto px-4 py-6">
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
                orderLeadMinutes={config?.orderLeadMinutes ?? 0}
                onChange={(deliveryDays) => setData({ deliveryDays: [...deliveryDays] })}
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
                onChange={(totalCycles) => setData({ totalCycles })}
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
                cycleDiscountRules={config?.discountRules ?? []}
                cutoffOffsetMinutes={config?.cutoffOffsetMinutes ?? 0}
                error={confirmResult && "error" in confirmResult ? confirmResult.error : undefined}
                isSubmitting={isSubmitting}
                onConfirm={handleConfirm}
                onBack={() => setStep(6)}
              />
            )}
          </div>
        </div>
      )}
    </Sheet>
  )
}
