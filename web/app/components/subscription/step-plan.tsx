import { Button } from "~/components/ui/button"
import { calculatePricing, formatPrice } from "./pricing"
import { cn } from "~/lib/utils"
import type { CycleDiscountRule, Product } from "./types"


type Props = {
  products: readonly Product[]
  items: { productId: string; quantity: number }[]
  totalCycles: number
  cycleDiscountRules: readonly CycleDiscountRule[]
  autoRenew: boolean
  onChange: (totalCycles: number, autoRenew: boolean) => void
  onNext: () => void
}

export function StepPlan({ products, items, totalCycles, cycleDiscountRules, autoRenew, onChange, onNext }: Props) {
  // Build plan options: baseline (1 cycle) + one option per discount rule threshold
  const options = [
    { cycles: 1, label: "1 entrega", discount: 0 },
    ...cycleDiscountRules
      .slice()
      .sort((a, b) => a.cyclesThreshold - b.cyclesThreshold)
      .map((rule) => ({
        cycles: rule.cyclesThreshold,
        label: `${rule.cyclesThreshold} entregas`,
        discount: rule.discount,
      })),
  ]

  const pricing = calculatePricing(items, products, totalCycles, cycleDiscountRules)
  const basePricing = calculatePricing(items, products, 1, [])

  return (
    <div className="flex flex-col pb-32">
      <p className="text-sm text-muted-foreground mb-6">
        Pague mais entregas antecipadas e ganhe desconto adicional sobre o valor da cesta.
      </p>

      <div className="space-y-3">
        {options.map((opt) => {
          const optPricing = calculatePricing(items, products, opt.cycles, cycleDiscountRules)
          const savings = basePricing.paymentAmount - optPricing.paymentAmount
          const isSelected = totalCycles === opt.cycles

          return (
            <button
              key={opt.cycles}
              type="button"
              onClick={() => onChange(opt.cycles, autoRenew)}
              className={cn(
                "w-full rounded-xl border p-4 text-left transition-colors",
                isSelected
                  ? "border-primary bg-primary/5"
                  : "border-border bg-card hover:bg-muted"
              )}
            >
              <div className="flex items-start justify-between gap-2">
                <div>
                  <p className="font-semibold">{opt.label}</p>
                  {opt.discount > 0 && (
                    <p className="text-sm text-primary font-medium">+{opt.discount}% de desconto extra</p>
                  )}
                  {savings > 0 && (
                    <p className="text-xs text-muted-foreground mt-0.5">
                      Economia de {formatPrice(savings)} por entrega
                    </p>
                  )}
                </div>
                <div className="text-right shrink-0">
                  {opt.discount > 0 && (
                    <p className="text-xs text-muted-foreground line-through">
                      {formatPrice(basePricing.paymentAmount)}
                    </p>
                  )}
                  <p className="font-bold text-primary">{formatPrice(optPricing.paymentAmount)}</p>
                  <p className="text-xs text-muted-foreground">/entrega</p>
                </div>
              </div>

              {isSelected && (
                <div className="mt-2 flex justify-end">
                  <span className="text-xs font-medium text-primary">✓ Selecionado</span>
                </div>
              )}
            </button>
          )
        })}
      </div>

      <div className="mt-6 flex items-center justify-between rounded-lg border p-3">
        <div>
          <p className="text-sm font-medium">Renovação automática</p>
          <p className="text-xs text-muted-foreground">Renova ao consumir todos os ciclos</p>
        </div>
        <button
          type="button"
          role="switch"
          aria-checked={autoRenew}
          onClick={() => onChange(totalCycles, !autoRenew)}
          className={cn(
            "relative w-11 h-6 rounded-full transition-colors",
            autoRenew ? "bg-primary" : "bg-muted border border-border"
          )}
        >
          <span
            className={cn(
              "absolute top-0.5 left-0.5 w-5 h-5 rounded-full bg-white shadow transition-transform",
              autoRenew && "translate-x-5"
            )}
          />
        </button>
      </div>

      <div className="fixed inset-x-0 bottom-0 border-t bg-background p-4">
        <Button className="w-full" onClick={onNext}>
          Continuar
        </Button>
      </div>
    </div>
  )
}
