import { useEffect, useMemo, useRef, useState } from "react"
import type { SavedCard } from "./types"
import { useFetcher } from "react-router"
import { loadStripe } from "@stripe/stripe-js"
import { Elements } from "@stripe/react-stripe-js"
import { useTranslation } from "react-i18next"
import { Sheet, SheetContent } from "~/components/ui/sheet"
import { cn } from "~/lib/utils"
import { isTokenExpired, getStoredRefreshToken, storeRefreshToken, clearRefreshToken } from "~/lib/token"
import { useWizard } from "./use-wizard"
import { StepProducts } from "./step-products"
import { StepBasket } from "./step-basket"
import { StepDelivery } from "./step-delivery"
import { StepPlan } from "./step-plan"
import { StepUser } from "./step-user"
import { StepAddress } from "./step-address"
import { StepConfirm } from "./step-confirm"
import { StepPayment } from "./step-payment"
import type { Product, SubscriptionConfig, UserProfile } from "./types"

type OtpActionResult =
  | { intent: "send-otp"; ok: true }
  | { intent: "verify-otp"; ok: true; accessToken: string; refreshToken: string }
  | { intent: "refresh"; ok: true; accessToken: string; refreshToken: string }
  | { intent: "refresh"; ok: false }
  | { intent: "get-profile"; ok: true; profile: UserProfile }
  | { intent: "get-payment-profile"; ok: true; hasPaymentMethod: boolean; card?: SavedCard }
  | { intent: "save-address"; ok: true }
  | { intent: string; error: string }

type CreateIntentResult =
  | { intent: "create-payment-intent"; ok: true; clientSecret: string }
  | { intent: "create-payment-intent"; error: string }

type ConfirmActionResult =
  | { ok: true; subscriptionIds: readonly string[] }
  | { ok: true; requiresAction: true; clientSecret: string }
  | { error: string }

const STEP_LABELS = {
  1: "subscription.drawer.steps.products",
  2: "subscription.drawer.steps.basket",
  3: "subscription.drawer.steps.delivery",
  4: "subscription.drawer.steps.plan",
  5: "subscription.drawer.steps.identity",
  6: "subscription.drawer.steps.address",
  summary: "subscription.drawer.steps.summary",
  payment: "subscription.drawer.steps.payment",
} as const

const NUMERIC_STEPS = [1, 2, 3, 4, 5, 6] as const

type Props = {
  readonly open: boolean
  readonly onClose: () => void
  readonly initialProductId?: string | null
  readonly products: readonly Product[]
  readonly config: SubscriptionConfig | null
  readonly serverToken: string | null
  readonly stripePublishableKey: string
}

export default function SubscriptionDrawer({ open, onClose, initialProductId, products, config, serverToken, stripePublishableKey }: Props) {
  const { t } = useTranslation()
  const loaderToken = serverToken

  const otpFetcher = useFetcher<OtpActionResult>()
  const addressFetcher = useFetcher<OtpActionResult>()
  const profileFetcher = useFetcher<OtpActionResult>()
  const paymentProfileFetcher = useFetcher<OtpActionResult>()
  const intentFetcher = useFetcher<CreateIntentResult>()
  const confirmFetcher = useFetcher<ConfirmActionResult>()

  const { step, setStep, data, setData, reset } = useWizard()
  const [otpSubStep, setOtpSubStep] = useState<"identity" | "code">("identity")
  const [forcedUserMode, setForcedUserMode] = useState<"new" | "existing" | undefined>(undefined)
  const [successVisible, setSuccessVisible] = useState(false)
  const [hasPaymentMethod, setHasPaymentMethod] = useState<boolean | null>(null)
  const [savedCard, setSavedCard] = useState<SavedCard | null>(null)
  const [paymentClientSecret, setPaymentClientSecret] = useState<string | null>(null)
  const refreshInProgressRef = useRef(false)

  const stripePromise = useMemo(
    () => (stripePublishableKey ? loadStripe(stripePublishableKey) : null),
    [stripePublishableKey],
  )

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
    setSuccessVisible(false)
    setOtpSubStep("identity")
    setHasPaymentMethod(null)
    setSavedCard(null)
    setPaymentClientSecret(null)
  }, [open]) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (!open || data.user.token || loaderToken) return
    const storedRefresh = getStoredRefreshToken()
    if (!storedRefresh) return

    fetch("/subscribe", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ intent: "refresh", refreshToken: storedRefresh }),
    })
      .then((res) => res.json())
      .then((result: { intent: "refresh"; ok: boolean; accessToken?: string; refreshToken?: string }) => {
        if (result.ok && result.accessToken && result.refreshToken) {
          storeRefreshToken(result.refreshToken)
          setData({ user: { ...data.user, token: result.accessToken } })
        } else {
          clearRefreshToken()
        }
      })
      .catch(() => clearRefreshToken())
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
    if (!existingProfile?.street.trim()) return
    setData({
      user: { ...data.user, phone: existingProfile.phone },
      address: {
        street: existingProfile.street,
        complement: existingProfile.complement,
        neighborhood: existingProfile.neighborhood,
        city: existingProfile.city,
        state: existingProfile.state,
        zip: existingProfile.zip,
      },
    })
  }, [profileFetcher.data]) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (
      addressFetcher.data &&
      "ok" in addressFetcher.data &&
      addressFetcher.data.intent === "save-address"
    ) {
      setStep("summary")
      ensureToken().then((token) => {
        if (!token) return
        paymentProfileFetcher.submit(
          { intent: "get-payment-profile", token },
          { method: "post", action: "/subscribe", encType: "application/json" },
        )
      })
    }
  }, [addressFetcher.data]) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (
      paymentProfileFetcher.data &&
      "ok" in paymentProfileFetcher.data &&
      paymentProfileFetcher.data.intent === "get-payment-profile"
    ) {
      setHasPaymentMethod(paymentProfileFetcher.data.hasPaymentMethod)
      setSavedCard(paymentProfileFetcher.data.card ?? null)
    }
  }, [paymentProfileFetcher.data]) // eslint-disable-line react-hooks/exhaustive-deps

  // Handle create-payment-intent response
  useEffect(() => {
    if (!intentFetcher.data) return
    if ("error" in intentFetcher.data) return
    if (intentFetcher.data.intent === "create-payment-intent" && "clientSecret" in intentFetcher.data) {
      if (!stripePromise) return
      setPaymentClientSecret(intentFetcher.data.clientSecret)
      setStep("payment")
    }
  }, [intentFetcher.data, stripePromise]) // eslint-disable-line react-hooks/exhaustive-deps

  // Handle confirm response
  const confirmResult = confirmFetcher.data as ConfirmActionResult | undefined
  useEffect(() => {
    if (!confirmResult) return
    if ("ok" in confirmResult && confirmResult.ok) {
      if ("requiresAction" in confirmResult && confirmResult.requiresAction) {
        // 3DS required on saved card — switch to payment step with new clientSecret
        setPaymentClientSecret(confirmResult.clientSecret)
        setStep("payment")
        return
      }
      setSuccessVisible(true)
    }
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

  async function handleConfirm(paymentIntentId?: string) {
    const token = await ensureToken()
    if (!token) return
    confirmFetcher.submit(
      {
        intent: "confirm",
        paymentIntentId: paymentIntentId ?? "",
        items: data.items,
        deliveryDays: data.deliveryDays,
        totalCycles: data.totalCycles,
        user: { ...data.user, token },
      // useFetcher.submit requires Record<string, string> but we need nested objects for JSON — cast is safe here
      } as unknown as Record<string, string>,
      { method: "post", action: "/subscribe", encType: "application/json" },
    )
  }

  async function handleConfirmFromSummary() {
    if (hasPaymentMethod) {
      handleConfirm()
      return
    }
    handlePayWithNewCard()
  }

  async function handlePayWithNewCard() {
    const token = await ensureToken()
    if (!token) return
    intentFetcher.submit(
      {
        intent: "create-payment-intent",
        token,
        items: data.items,
        totalCycles: data.totalCycles,
      },
      { method: "post", action: "/subscribe", encType: "application/json" },
    )
  }

  function renderStepCaption() {
    if (!currentStepNumber) return t("subscription.drawer.steps.summary")
    return t("subscription.drawer.stepProgress", { step: currentStepNumber, total: NUMERIC_STEPS.length })
  }

  function handleBack() {
    if (step === "summary") return setStep(6)
    if (step === "payment") return setStep("summary")
    if (step === 2) return setStep(1)
    if (step === 3) return setStep(2)
    if (step === 4) return setStep(3)
    if (step === 5 && otpSubStep === "code") return setOtpSubStep("identity")
    if (step === 5) return setStep(4)
    if (step === 6) return setStep(4)
  }

  function renderBackButton() {
    if (step === 1) return null
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
  const isLoadingIntent = intentFetcher.state !== "idle"
  const confirmError = confirmResult && "error" in confirmResult ? confirmResult.error : undefined
  const intentError = intentFetcher.data && "error" in intentFetcher.data ? intentFetcher.data.error : undefined
  const summaryError = confirmError ?? intentError

  const isConfirmStep = step === "summary"
  const isPaymentStep = step === "payment"
  const isSheetHidden = isConfirmStep || isPaymentStep || successVisible

  const currentStepNumber = step === "summary" || step === "payment" ? null : (step as number)
  const stepLabel = STEP_LABELS[step as keyof typeof STEP_LABELS]

  return (
    <>
      <Sheet open={open && !isSheetHidden} onOpenChange={(isOpen) => { if (!isOpen) onClose() }}>
        <SheetContent side="right" showCloseButton={false} className="w-full sm:max-w-md p-0 overflow-hidden">
          <div className="flex-1 flex flex-col min-h-0">
            <header className="shrink-0 border-b px-4 py-3">
              <div className="flex items-center gap-3">
                {renderBackButton()}

                <div className="flex-1 min-w-0">
                  <p className="text-xs text-muted-foreground">{renderStepCaption()}</p>
                  <p className="text-sm font-medium truncate">{t(stepLabel)}</p>
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
                <p className="text-sm text-destructive">{t("subscription.drawer.closed")}</p>
              </div>
            )}

            <div className="flex-1 min-h-0 overflow-y-auto px-4 py-6">
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
                  onNext={async () => {
                    const token = await ensureToken()
                    setStep(token ? 6 : 5)
                  }}
                />
              )}

              {step === 5 && (
                <StepUser
                  user={data.user}
                  subStep={otpSubStep}
                  initialMode={forcedUserMode}
                  existingIdentity={data.user.name || data.user.email || undefined}
                  onContinueAsExisting={data.user.token ? () => setStep(6) : undefined}
                  onChange={(user) => setData({ user })}
                  onSendOtp={handleSendOtp}
                  onVerifyOtp={handleVerifyOtp}
                  onResetOtp={() => setOtpSubStep("identity")}
                  isSending={otpFetcher.state !== "idle"}
                  isVerifying={false}
                  otpError={otpFetcher.data && "error" in otpFetcher.data ? otpFetcher.data.error : undefined}
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
                  onChangeIdentity={() => setStep(5)}
                  isSaving={addressFetcher.state !== "idle"}
                  saveError={addressFetcher.data && "error" in addressFetcher.data ? addressFetcher.data.error : undefined}
                />
              )}
            </div>
          </div>
        </SheetContent>
      </Sheet>

      {open && isConfirmStep && !successVisible && (
        <StepConfirm
          data={data}
          products={products}
          cycleDiscountRules={config?.discountRules ?? []}
          savedCard={savedCard}
          error={summaryError}
          isSubmitting={isSubmitting || isLoadingIntent}
          onConfirm={handleConfirmFromSummary}
          onPayWithNewCard={handlePayWithNewCard}
          onBack={() => setStep(6)}
        />
      )}

      {open && isPaymentStep && !successVisible && paymentClientSecret && stripePromise && (
        <Elements stripe={stripePromise} options={{ clientSecret: paymentClientSecret, locale: "pt-BR", appearance: { theme: "flat", variables: { borderRadius: "12px", fontSizeBase: "16px" } } }}>
          <StepPayment
            clientSecret={paymentClientSecret}
            isConfirming={isSubmitting}
            error={confirmError}
            onSuccess={(paymentIntentId) => handleConfirm(paymentIntentId)}
            onBack={() => setStep("summary")}
          />
        </Elements>
      )}

      {open && isSubmitting && !successVisible && (isConfirmStep || !paymentClientSecret) && (
        <div className="fixed inset-0 z-[60] bg-background flex flex-col items-center justify-center px-6 text-center animate-in fade-in-0 duration-200">
          <div className="text-4xl mb-6 animate-pulse">⏳</div>
          <p className="text-lg font-medium mb-2">{t("payment.processing")}</p>
          <p className="text-sm text-muted-foreground">{t("payment.processingHint")}</p>
        </div>
      )}

      {open && successVisible && (
        <div className="fixed inset-0 z-50 bg-background flex flex-col items-center justify-center px-6 text-center animate-in fade-in-0 duration-200">
          <div className="text-5xl mb-4">🎉</div>
          <h1 className="text-2xl font-bold mb-2">{t("subscription.drawer.success.title")}</h1>
          <p className="text-muted-foreground mb-8">{t("subscription.drawer.success.description")}</p>
          <button
            type="button"
            onClick={onClose}
            className="text-sm text-primary font-medium hover:underline"
          >
            {t("subscription.drawer.success.close")}
          </button>
        </div>
      )}
    </>
  )
}
