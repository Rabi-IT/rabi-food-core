import { useEffect, useState } from "react"
import { useTranslation } from "react-i18next"
import { Button } from "~/components/ui/button"
import { cn } from "~/lib/utils"
import { calculatePricing, formatPrice } from "./pricing"
import type { WizardData } from "./use-wizard"
import type { CycleDiscountRule, Product, SavedCard } from "./types"

const FIRST_SLIDE = 0
const LAST_SLIDE = 3
const SLIDE_DOTS = [0, 1, 2, 3] as const

// Apple web easing — moderate start, gradual organic deceleration (reverse-engineered from apple.com CSS).
const TRANSITION = "transition-all duration-500 [transition-timing-function:cubic-bezier(0.4,0.01,0.165,0.99)]"

// Beat = 1000ms — fixed tempo across all slides (~60 BPM, resting heart rate).
// All elements share the same rhythm regardless of text length.
// First element at 300ms: attentional blink minimum.
const BEAT = 1000
const FIRST_BEAT = 300
const SLIDE_REVEAL_DELAYS: Record<number, readonly number[]> = {
  0: [FIRST_BEAT, FIRST_BEAT + BEAT, FIRST_BEAT + BEAT * 2, FIRST_BEAT + BEAT * 3, FIRST_BEAT + BEAT * 4],
  1: [FIRST_BEAT, FIRST_BEAT + BEAT, FIRST_BEAT + BEAT * 2, FIRST_BEAT + BEAT * 3],
  2: [FIRST_BEAT, FIRST_BEAT + BEAT, FIRST_BEAT + BEAT * 2],
  3: [FIRST_BEAT, FIRST_BEAT + BEAT, FIRST_BEAT + BEAT * 2, FIRST_BEAT + BEAT * 3],
}

type StepConfirmProps = {
  readonly data: WizardData
  readonly products: readonly Product[]
  readonly cycleDiscountRules: readonly CycleDiscountRule[]
  readonly savedCard?: SavedCard | null
  readonly error?: string
  readonly isSubmitting: boolean
  readonly onConfirm: () => void
  readonly onPayWithNewCard: () => void
  readonly onBack: () => void
}

function formatPhone(raw: string): string {
  const digits = raw.replace(/\D/g, "")
  const local = digits.startsWith("55") ? digits.slice(2) : digits
  if (local.length === 11) return `(${local.slice(0, 2)}) ${local.slice(2, 7)}-${local.slice(7)}`
  if (local.length === 10) return `(${local.slice(0, 2)}) ${local.slice(2, 6)}-${local.slice(6)}`
  return raw
}

function getFirstDeliveryDate(deliveryDays: readonly { weekday: number }[]): Date {
  const days = deliveryDays.map((d) => d.weekday).sort((a, b) => a - b)
  const cursor = new Date()
  cursor.setHours(0, 0, 0, 0)
  cursor.setDate(cursor.getDate() + 1)
  while (!days.includes(cursor.getDay())) {
    cursor.setDate(cursor.getDate() + 1)
  }
  return cursor
}

// Unrevealed: opacity-0 with NO transition (prevents flash on first render —
// browser would otherwise animate from default opacity:1 to opacity:0).
// Revealed: transition kicks in only when going from hidden → visible.
function revealCn(step: number, revealed: number, extra?: string): string {
  if (revealed >= step) return cn(TRANSITION, "opacity-100 translate-y-0", extra)
  return cn("opacity-0 translate-y-2", extra)
}

function useReveal(slide: number): number {
  const [revealed, setRevealed] = useState(0)

  useEffect(() => {
    setRevealed(0)
    const delays = SLIDE_REVEAL_DELAYS[slide] ?? []
    const timeouts = delays.map((delay, index) =>
      setTimeout(() => setRevealed(index + 1), delay),
    )
    return () => {
      timeouts.forEach(clearTimeout)
    }
  }, [slide]) // eslint-disable-line react-hooks/exhaustive-deps

  return revealed
}

export function StepConfirm({
  data,
  products,
  cycleDiscountRules,
  savedCard,
  error,
  isSubmitting,
  onConfirm,
  onPayWithNewCard,
  onBack,
}: StepConfirmProps) {
  const { t } = useTranslation()
  const [slide, setSlide] = useState(FIRST_SLIDE)
  const revealed = useReveal(slide)

  const flexibilityItems = t("subscription.confirm.flexibility.items", { returnObjects: true }) as readonly string[]

  const pricing = calculatePricing(data.items, products, data.totalCycles, cycleDiscountRules)
  const total = pricing.paymentAmount * data.totalCycles
  const firstDelivery = getFirstDeliveryDate(data.deliveryDays)
  const firstDeliveryLabel = firstDelivery.toLocaleDateString("pt-BR", {
    weekday: "long",
    day: "numeric",
    month: "long",
  })
  const deliveryLabel = data.totalCycles === 1
    ? t("subscription.confirm.deliveryOne", { count: data.totalCycles })
    : t("subscription.confirm.deliveryMany", { count: data.totalCycles })
  const confirmLabel = isSubmitting ? t("subscription.confirm.processing") : t("subscription.confirm.confirm")
  const isLastSlide = slide === LAST_SLIDE

  function handleBack() {
    if (slide === FIRST_SLIDE) {
      onBack()
      return
    }
    setSlide((s) => s - 1)
  }

  function renderButton() {
    if (isLastSlide) {
      return (
        <Button className="w-full" onClick={onConfirm} disabled={isSubmitting}>
          {confirmLabel}
        </Button>
      )
    }
    return (
      <Button className="w-full" onClick={() => setSlide((s) => s + 1)}>
        {t("subscription.confirm.continue")}
      </Button>
    )
  }

  function renderSlide3() {
    if (savedCard) {
      return (
        <div className="space-y-10">
          {data.user.phone && (
            <div className={revealCn(1, revealed)}>
              <p className="text-xs text-muted-foreground uppercase tracking-wide mb-1">
                {t("subscription.confirm.contactLabel")}
              </p>
              <p className="text-base font-medium">{formatPhone(data.user.phone)}</p>
            </div>
          )}

          <div className={revealCn(2, revealed)}>
            <p className="text-xs text-muted-foreground uppercase tracking-wide mb-3">
              {t("subscription.confirm.chargedTo")}
            </p>
            <div className="flex items-center justify-between rounded-2xl border border-border bg-muted/40 px-4 py-3">
              <div className="flex items-center gap-3">
                <span className="text-2xl">💳</span>
                <div>
                  <p className="text-sm font-semibold leading-tight">
                    {savedCard.brand.charAt(0).toUpperCase() + savedCard.brand.slice(1)} •••• {savedCard.last4}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {String(savedCard.expMonth).padStart(2, "0")}/{String(savedCard.expYear).slice(-2)}
                  </p>
                </div>
              </div>
              <button
                type="button"
                onClick={onPayWithNewCard}
                disabled={isSubmitting}
                className="text-xs text-primary hover:text-primary/80 transition-colors font-medium shrink-0 ml-4"
              >
                {t("subscription.confirm.changeCard")}
              </button>
            </div>
          </div>

          {error && (
            <div className={revealCn(3, revealed)}>
              <p className="text-sm text-destructive rounded-lg border border-destructive/30 bg-destructive/5 p-3">
                {error}
              </p>
            </div>
          )}
        </div>
      )
    }

    return (
      <div className="space-y-10">
        {data.user.phone && (
          <div className={revealCn(1, revealed)}>
            <p className="text-xs text-muted-foreground uppercase tracking-wide mb-1">
              {t("subscription.confirm.contactLabel")}
            </p>
            <p className="text-base font-medium">{formatPhone(data.user.phone)}</p>
          </div>
        )}

        <div className={revealCn(2, revealed, "space-y-1")}>
          <p className="text-xs text-muted-foreground uppercase tracking-wide mb-1">
            {t("subscription.confirm.paymentMethod")}
          </p>
          <p className="text-2xl font-semibold leading-snug">{t("subscription.confirm.newCard")}</p>
          <p className="text-base text-muted-foreground">
            {t("subscription.confirm.newCardHint")}
          </p>
        </div>

        {error && (
          <div className={revealCn(3, revealed)}>
            <p className="text-sm text-destructive rounded-lg border border-destructive/30 bg-destructive/5 p-3">
              {error}
            </p>
          </div>
        )}
      </div>
    )
  }

  function renderSlideContent() {
    if (slide === 3) return renderSlide3()

    if (slide === 2) {
      return (
        <div className="space-y-6">
          <p className={revealCn(1, revealed, "text-xs text-muted-foreground uppercase tracking-wide")}>
            {t("subscription.confirm.firstDelivery")}
          </p>
          <h2 className={revealCn(2, revealed, "text-4xl font-bold font-display leading-tight capitalize")}>
            {firstDeliveryLabel}
          </h2>
          <p className={revealCn(3, revealed, "text-base text-muted-foreground leading-relaxed")}>
            {t("subscription.confirm.weeklyHint")}
          </p>
        </div>
      )
    }

    if (slide === 1) {
      return (
        <div className="space-y-3">
          <p className={revealCn(1, revealed, "text-xs text-muted-foreground uppercase tracking-wide")}>
            {t("subscription.confirm.youAreTaking")}
          </p>
          <h2 className={revealCn(2, revealed, "text-5xl font-bold font-display text-primary leading-none")}>
            {deliveryLabel}
          </h2>
          <p className={revealCn(3, revealed, "text-4xl font-bold font-display leading-none")}>
            {formatPrice(total)}
          </p>
          <div className={revealCn(4, revealed)}>
            <p className="text-sm text-muted-foreground">
              {t("subscription.confirm.pricePerDelivery", { price: formatPrice(pricing.paymentAmount) })}
            </p>
            {pricing.totalDiscount > 0 && (
              <p className="text-sm font-medium text-primary pt-2">
                {t("subscription.confirm.savings", { amount: formatPrice(pricing.totalDiscount * data.totalCycles) })}
              </p>
            )}
          </div>
        </div>
      )
    }

    return (
      <div className="space-y-8">
        <div className={revealCn(1, revealed, "space-y-2")}>
          <h2 className="text-4xl font-bold font-display leading-tight">
            {t("subscription.confirm.flexibility.title")}
          </h2>
          <p className="text-base text-muted-foreground">{t("subscription.confirm.flexibility.subtitle")}</p>
        </div>
        <ul className="space-y-5">
          {flexibilityItems.map((item, index) => (
            <li key={item} className={revealCn(index + 2, revealed, "flex items-start gap-3")}>
              <span className="text-primary font-bold text-lg leading-none mt-0.5">✓</span>
              <span className="text-base">{item}</span>
            </li>
          ))}
        </ul>
      </div>
    )
  }

  return (
    <div className="fixed inset-0 z-50 bg-background flex flex-col animate-in fade-in-0 duration-300 [animation-timing-function:cubic-bezier(0.4,0.01,0.165,0.99)]">
      <header className="shrink-0 flex items-center justify-between px-6 pt-8 pb-4">
        <button
          type="button"
          onClick={handleBack}
          className="text-sm text-muted-foreground hover:text-foreground transition-colors"
        >
          ←
        </button>
        <div className="flex gap-2">
          {SLIDE_DOTS.map((dot) => (
            <div
              key={dot}
              className={cn(
                "w-2 h-2 rounded-full transition-colors duration-300",
                dot <= slide ? "bg-primary" : "bg-muted",
              )}
            />
          ))}
        </div>
        <div className="w-8" />
      </header>

      <div
        key={slide}
        className="flex-1 min-h-0 px-6 overflow-y-auto animate-in fade-in-0 duration-300 [animation-timing-function:cubic-bezier(0.4,0.01,0.165,0.99)]"
      >
        <div className="min-h-full flex flex-col justify-start pt-10 pb-6">
          {renderSlideContent()}
        </div>
      </div>

      <div className="shrink-0 px-6 pb-12 pt-4">
        {renderButton()}
      </div>
    </div>
  )
}
