import { useState } from "react"
import { useTranslation } from "react-i18next"
import { PaymentElement, useStripe, useElements } from "@stripe/react-stripe-js"
import { Button } from "~/components/ui/button"

type StepPaymentProps = {
  readonly clientSecret: string
  readonly isConfirming: boolean
  readonly error?: string
  readonly onSuccess: (paymentIntentId: string) => void
  readonly onBack: () => void
}

export function StepPayment({ clientSecret, isConfirming, error, onSuccess, onBack }: StepPaymentProps) {
  const { t } = useTranslation()
  const stripe = useStripe()
  const elements = useElements()
  const [cardError, setCardError] = useState<string | undefined>(undefined)
  const [isProcessing, setIsProcessing] = useState(false)
  const [isReady, setIsReady] = useState(false)

  const isBusy = isProcessing || isConfirming
  const displayError = cardError ?? error

  async function handleConfirm() {
    if (!stripe || !elements) return
    setCardError(undefined)
    setIsProcessing(true)

    const { error: stripeError, paymentIntent } = await stripe.confirmPayment({
      elements,
      confirmParams: { return_url: window.location.href },
      redirect: "if_required",
    })

    if (stripeError) {
      setCardError(stripeError.message ?? t("payment.errorGeneric"))
      setIsProcessing(false)
      return
    }

    if (paymentIntent?.status === "succeeded") {
      onSuccess(paymentIntent.id)
      return
    }

    setCardError(t("payment.errorGeneric"))
    setIsProcessing(false)
  }

  return (
    <div className="fixed inset-0 z-50 bg-background flex flex-col animate-in fade-in-0 duration-300 [animation-timing-function:cubic-bezier(0.4,0.01,0.165,0.99)]">
      <header className="shrink-0 flex items-center justify-between px-6 pt-8 pb-4">
        <button
          type="button"
          onClick={onBack}
          disabled={isBusy}
          className="text-sm text-muted-foreground hover:text-foreground transition-colors disabled:opacity-40"
        >
          ←
        </button>
        <div className="w-8" />
      </header>

      <div className="flex-1 min-h-0 px-6 overflow-y-auto">
        <div className="pt-10 pb-6 space-y-6">
          <h2 className="text-4xl font-bold font-display leading-tight">
            {t("payment.title")}
          </h2>

          <PaymentElement
            onReady={() => setIsReady(true)}
            options={{
              defaultValues: { billingDetails: { address: { country: "BR" } } },
            }}
          />

          {displayError && (
            <p className="text-sm text-destructive rounded-lg border border-destructive/30 bg-destructive/5 p-3">
              {displayError}
            </p>
          )}
        </div>
      </div>

      <div className="shrink-0 px-6 pb-12 pt-4">
        <Button
          className="w-full"
          onClick={handleConfirm}
          disabled={isBusy || !stripe || !isReady}
        >
          {isBusy ? t("payment.confirming") : t("payment.confirm")}
        </Button>
      </div>
    </div>
  )
}
