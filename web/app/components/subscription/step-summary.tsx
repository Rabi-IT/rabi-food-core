import { Button } from "~/components/ui/button"
import { calculatePricing, formatPrice } from "./pricing"
import type { WizardData } from "./use-wizard"
import type { CycleDiscountRule, Product } from "./types"

const WEEKDAYS = ["Domingo", "Segunda", "Terça", "Quarta", "Quinta", "Sexta", "Sábado"]

type Props = {
  data: WizardData
  products: readonly Product[]
  cycleDiscountRules: readonly CycleDiscountRule[]
  cutoffOffsetMinutes: number
  error?: string
  isSubmitting: boolean
  onConfirm: () => void
  onBack: () => void
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

  return (
    <div className="flex flex-col pb-48">
      {/* Products */}
      <section className="space-y-3">
        <h3 className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">
          Sua cesta
        </h3>
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

      {/* Pricing breakdown */}
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
              Aplicado sobre {formatPrice(pricing.itemsTotal - pricing.itemsDiscount)} (valor após desconto por quantidade)
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

      {/* Delivery schedule */}
      <section className="space-y-2">
        <h3 className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">
          Entregas
        </h3>
        <p className="text-sm">{selectedDays.join(", ")}</p>
        <p className="text-sm text-muted-foreground">Horário comercial (8h–18h)</p>
      </section>

      <div className="my-4 border-t" />

      {/* Benefits & policies */}
      <section className="space-y-3 rounded-xl bg-primary/5 border border-primary/20 p-4">
        <p className="text-sm font-medium text-primary">Benefícios da assinatura</p>
        <ul className="space-y-1.5 text-sm text-foreground/80">
          <li>✓ Altere os produtos da cesta quando quiser</li>
          <li>✓ Adie uma entrega se não puder receber</li>
          <li>✓ Altere os dias e horários de entrega a qualquer momento</li>
          <li>✓ Pule entregas para receber com menos frequência</li>
          {data.autoRenew && <li>✓ Renovação automática ao consumir todas as entregas</li>}
          {cutoffHours > 0 && (
            <li className="text-muted-foreground">
              Alterações até {cutoffHours}h antes da entrega sem custo
            </li>
          )}
        </ul>
      </section>

      {/* Recipient */}
      <div className="my-4 border-t" />
      <section className="space-y-1">
        <h3 className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">
          Entrega para
        </h3>
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
