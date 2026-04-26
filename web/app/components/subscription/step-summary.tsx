import { Button } from "~/components/ui/button"
import { calculatePricing, formatPrice } from "./pricing"
import type { WizardData } from "./use-wizard"
import type { CycleDiscountRule, Product } from "./types"

const WEEKDAYS = ["Domingo", "Segunda", "Terça", "Quarta", "Quinta", "Sexta", "Sábado"] as const

type Props = {
  readonly data: WizardData
  readonly products: readonly Product[]
  readonly cycleDiscountRules: readonly CycleDiscountRule[]
  readonly cutoffOffsetMinutes: number
  readonly error?: string
  readonly isSubmitting: boolean
  readonly onConfirm: () => void
  readonly onBack: () => void
}

function getFirstDeliveryDate(weekdays: readonly { weekday: number }[]): Date {
  const days = weekdays.map((d) => d.weekday).sort((a, b) => a - b)
  const cursor = new Date()
  cursor.setHours(0, 0, 0, 0)
  cursor.setDate(cursor.getDate() + 1)
  while (!days.includes(cursor.getDay())) {
    cursor.setDate(cursor.getDate() + 1)
  }
  return cursor
}

export function StepSummary({
  data,
  products,
  cycleDiscountRules,
  cutoffOffsetMinutes,
  error,
  isSubmitting,
  onConfirm,
  onBack,
}: Props) {
  const pricing = calculatePricing(data.items, products, data.totalCycles, cycleDiscountRules)
  const cutoffHours = Math.round(cutoffOffsetMinutes / 60)
  const selectedDays = data.deliveryDays
    .slice()
    .sort((a, b) => a.weekday - b.weekday)
    .map((d) => WEEKDAYS[d.weekday])

  const firstDelivery = getFirstDeliveryDate(data.deliveryDays)
  const firstDeliveryLabel = firstDelivery.toLocaleDateString("pt-BR", {
    weekday: "long", day: "numeric", month: "long",
  })

  return (
    <div className="flex flex-col pb-48">
      <section className="space-y-3">
        <h3 className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">Sua cesta</h3>
        {pricing.items.map((item) => (
          <div key={item.productId} className="flex items-center justify-between gap-2">
            <div>
              <p className="text-sm font-medium">{item.name}</p>
              <p className="text-xs text-muted-foreground">
                {item.quantity}× {formatPrice(item.unitPrice)}
                {item.discountPercent > 0 && (
                  <span className="ml-1 text-primary">− {item.discountPercent}%</span>
                )}
              </p>
            </div>
            <p className="text-sm font-semibold shrink-0">{formatPrice(item.total)}</p>
          </div>
        ))}
      </section>

      <div className="my-4 border-t" />

      <section className="space-y-2 text-sm">
        <div className="flex justify-between">
          <span className="text-muted-foreground">Subtotal</span>
          <span>{formatPrice(pricing.itemsTotal)}</span>
        </div>

        {pricing.itemsDiscount > 0 && (
          <div className="flex justify-between">
            <span className="text-muted-foreground">Desconto por quantidade</span>
            <span className="text-primary">− {formatPrice(pricing.itemsDiscount)}</span>
          </div>
        )}

        {pricing.cycleDiscount > 0 && (
          <>
            <div className="flex justify-between border-t pt-2">
              <span className="text-muted-foreground">
                Desconto do plano ({pricing.cycleDiscountPercent}%)
              </span>
              <span className="text-primary">− {formatPrice(pricing.cycleDiscount)}</span>
            </div>
            <p className="text-xs text-muted-foreground">
              Aplicado sobre {formatPrice(pricing.itemsTotal - pricing.itemsDiscount)}
            </p>
          </>
        )}

        <div className="flex justify-between font-semibold border-t pt-2 mt-1">
          <span>Total por entrega</span>
          <span className="text-primary text-lg">{formatPrice(pricing.paymentAmount)}</span>
        </div>

        {data.totalCycles > 1 && (
          <div className="flex justify-between text-muted-foreground">
            <span>{data.totalCycles} entregas × {formatPrice(pricing.paymentAmount)}</span>
            <span className="font-medium text-foreground">{formatPrice(pricing.paymentAmount * data.totalCycles)}</span>
          </div>
        )}

        {pricing.totalDiscount > 0 && (
          <div className="flex justify-between font-medium text-primary">
            <span>Você economiza no plano</span>
            <span>− {formatPrice(pricing.totalDiscount * data.totalCycles)}</span>
          </div>
        )}
      </section>

      <div className="my-4 border-t" />

      <section className="space-y-2">
        <h3 className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">Entregas</h3>
        <p className="text-sm">{selectedDays.join(", ")}</p>
        {cutoffHours > 0 && (
          <p className="text-xs text-muted-foreground">
            Alterações até {cutoffHours}h antes da entrega
          </p>
        )}
        <div className="mt-2 rounded-lg bg-primary/5 border border-primary/20 px-3 py-2">
          <p className="text-sm font-medium text-primary">
            Primeira entrega: <span className="font-semibold">{firstDeliveryLabel}</span>
          </p>
        </div>
      </section>

      <div className="my-4 border-t" />

      <section className="space-y-3">
        <h3 className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">Flexibilidade</h3>
        <ul className="space-y-1.5 text-sm text-foreground/80">
          <li>✓ Troque qualquer produto da cesta quando quiser</li>
          <li>✓ Adie ou pule uma entrega se não puder receber</li>
          <li>✓ Altere os dias e horários a qualquer momento</li>
          <li>✓ Cancele a renovação quando quiser — sem multa</li>
        </ul>
      </section>

      <div className="my-4 border-t" />

      <section className="space-y-3 rounded-xl border p-4">
        <p className="text-sm font-semibold">Cancelamento e reembolso</p>
        <div className="space-y-3 text-sm text-muted-foreground">
          <div>
            <p className="font-medium text-foreground">Não quer renovar?</p>
            <p>Cancele a renovação a qualquer momento. As entregas já pagas continuam chegando até o fim do ciclo.</p>
          </div>
          <div>
            <p className="font-medium text-foreground">Arrependeu dentro de 7 dias após a primeira entrega?</p>
            <p>Reembolsamos as entregas não realizadas automaticamente, sem burocracia.</p>
          </div>
          <div>
            <p className="font-medium text-foreground">Depois dos 7 dias?</p>
            <p>Sem reembolso, mas sem multa. As entregas que já pagou continuam sendo realizadas.</p>
          </div>
        </div>
      </section>

      <div className="my-4 border-t" />

      <section className="space-y-1">
        <h3 className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">Entrega para</h3>
        <p className="text-sm font-medium">{data.user.name}</p>
        <p className="text-sm text-muted-foreground">
          {data.address.street}
          {data.address.complement ? `, ${data.address.complement}` : ""}
        </p>
        <p className="text-sm text-muted-foreground">
          {data.address.neighborhood} — {data.address.city}/{data.address.state}
        </p>
      </section>

      <div className="fixed inset-x-0 bottom-0 border-t bg-background p-4 space-y-2">
        {error && (
          <p className="text-sm text-destructive text-center">{error}</p>
        )}
        <Button className="w-full" onClick={onConfirm} disabled={isSubmitting}>
          {isSubmitting ? "Processando…" : "Confirmar assinatura"}
        </Button>
        <button
          type="button"
          onClick={onBack}
          disabled={isSubmitting}
          className="w-full text-sm text-muted-foreground hover:text-foreground transition-colors py-1"
        >
          Voltar e editar
        </button>
      </div>
    </div>
  )
}
