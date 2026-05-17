import { useEffect } from "react"
import { Button } from "~/components/ui/button"
import { calculatePricing, formatPrice } from "./pricing"
import { cn } from "~/lib/utils"
import type { CycleDiscountRule, Product } from "./types"

type DeliveryOption = {
  readonly cycles: number
  readonly label: string
  readonly discount: number
}

type Props = {
  readonly products: readonly Product[]
  readonly items: readonly { productId: string; quantity: number }[]
  readonly totalCycles: number
  readonly minCycles: number
  readonly deliveryWeekdays: readonly number[]
  readonly cycleDiscountRules: readonly CycleDiscountRule[]
  readonly onChange: (totalCycles: number) => void
  readonly onNext: () => void
}

function buildOptions(minCycles: number, rules: readonly CycleDiscountRule[]): readonly DeliveryOption[] {
  const cycleLabel = (n: number) => `${n} ${n === 1 ? "entrega" : "entregas"}`
  const extras = rules
    .slice()
    .sort((a, b) => a.cyclesThreshold - b.cyclesThreshold)
    .filter((rule) => rule.cyclesThreshold > minCycles)
    .map((rule) => ({ cycles: rule.cyclesThreshold, label: cycleLabel(rule.cyclesThreshold), discount: rule.discount }))
  return [{ cycles: minCycles, label: cycleLabel(minCycles), discount: 0 }, ...extras]
}

export function StepPlan({ products, items, totalCycles, minCycles, deliveryWeekdays: _deliveryWeekdays, cycleDiscountRules, onChange, onNext }: Props) {
  const options = buildOptions(minCycles, cycleDiscountRules)
  const selectedPricing = calculatePricing(items, products, totalCycles, cycleDiscountRules)
  const basePricing = calculatePricing(items, products, minCycles, [])
  const selectedTotal = selectedPricing.paymentAmount * totalCycles
  const totalSavings = basePricing.paymentAmount * totalCycles - selectedTotal

  useEffect(() => {
    if (options.length === 1) onChange(options[0].cycles)
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div className="flex flex-col gap-6">
      <p className="text-sm text-muted-foreground leading-relaxed">
        🛒 Cada entrega é uma cesta no dia combinado. Mais entregas, mais desconto — escolha o que faz sentido pra você.
      </p>

      <div className="space-y-3">
        {options.map((opt) => {
          const optPricing = calculatePricing(items, products, opt.cycles, cycleDiscountRules)
          const optTotal = optPricing.paymentAmount * opt.cycles
          const isSelected = totalCycles === opt.cycles

          return (
            <button
              key={opt.cycles}
              type="button"
              onClick={() => onChange(opt.cycles)}
              className={cn(
                "w-full rounded-xl border p-4 text-left transition-colors",
                isSelected ? "border-primary bg-primary/5" : "border-border bg-card hover:bg-muted"
              )}
            >
              <div className="flex items-start justify-between gap-2">
                <div className="flex flex-col items-start gap-1.5">
                  <p className="font-semibold font-display">{opt.label}</p>
                  {opt.discount > 0 && (
                    <span className="text-xs font-medium text-primary bg-primary/10 rounded-full px-2 py-0.5">
                      −{opt.discount}% de desconto
                    </span>
                  )}
                  <span className={cn(
                    "text-xs font-medium rounded-full px-2 py-0.5",
                    isSelected ? "text-primary bg-primary/10" : "invisible"
                  )}>
                    ✓ Selecionado
                  </span>
                </div>
                <div className="text-right shrink-0">
                  <p className="text-xl font-bold text-primary leading-none">
                    {formatPrice(optTotal)}
                  </p>
                  <p className="text-xs text-muted-foreground mt-0.5">
                    {formatPrice(optPricing.paymentAmount)}/entrega
                  </p>
                </div>
              </div>
            </button>
          )
        })}
      </div>

      <div className="rounded-xl border border-primary/20 bg-primary/5 p-4 space-y-1">
        <p className="text-xs font-semibold text-primary uppercase tracking-wide mb-3">
          📦 Você está levando
        </p>
        <p className="text-sm text-foreground/80">
          {totalCycles} {totalCycles === 1 ? "entrega" : "entregas"} da sua cesta
        </p>
        <p className="text-2xl font-bold text-foreground pt-1">
          {formatPrice(selectedTotal)}
        </p>
        <p className="text-xs text-muted-foreground">
          {formatPrice(selectedPricing.paymentAmount)} por entrega
        </p>
        <p className={cn("text-sm font-medium pt-1", totalSavings > 0 ? "text-primary" : "invisible")}>
          💰 Você economiza {formatPrice(totalSavings)} no total
        </p>
      </div>

      <div className="sticky bottom-0 border-t bg-background pt-4 -mx-4 px-4 pb-2">
        <Button className="w-full" onClick={onNext}>
          Continuar
        </Button>
      </div>
    </div>
  )
}
