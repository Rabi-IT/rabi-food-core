import { useEffect } from "react"
import { Button } from "~/components/ui/button"
import { calculatePricing, formatPrice } from "./pricing"
import { cn } from "~/lib/utils"
import type { CycleDiscountRule, Product } from "./types"

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

function getRenewalDate(weekdays: readonly number[], totalCycles: number): Date {
  const sorted = [...weekdays].sort((a, b) => a - b)
  const cursor = new Date()
  cursor.setHours(0, 0, 0, 0)
  cursor.setDate(cursor.getDate() + 1)
  let count = 0
  while (count < totalCycles) {
    if (sorted.includes(cursor.getDay())) count++
    if (count < totalCycles) cursor.setDate(cursor.getDate() + 1)
  }
  return cursor
}

export function StepPlan({ products, items, totalCycles, minCycles, deliveryWeekdays, cycleDiscountRules, onChange, onNext }: Props) {
  const label = (n: number) => `${n} ${n === 1 ? "entrega" : "entregas"}`

  const options = [
    { cycles: minCycles, label: label(minCycles), discount: 0 },
    ...cycleDiscountRules
      .slice()
      .sort((a, b) => a.cyclesThreshold - b.cyclesThreshold)
      .filter((rule) => rule.cyclesThreshold > minCycles)
      .map((rule) => ({
        cycles: rule.cyclesThreshold,
        label: label(rule.cyclesThreshold),
        discount: rule.discount,
      })),
  ]

  const basePricing = calculatePricing(items, products, minCycles, [])

  useEffect(() => {
    if (options.length === 1) onChange(options[0].cycles)
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  const renewalDate = getRenewalDate(deliveryWeekdays, totalCycles)

  return (
    <div className="flex flex-col">
      <p className="text-sm text-muted-foreground mb-1">
        Cada entrega corresponde a um dia da semana que você selecionou.
      </p>
      <p className="text-sm text-muted-foreground mb-6">
        Pague mais entregas antecipadas e ganhe desconto adicional.
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
              onClick={() => onChange(opt.cycles)}
              className={cn(
                "w-full rounded-xl border p-4 text-left transition-colors",
                isSelected ? "border-primary bg-primary/5" : "border-border bg-card hover:bg-muted"
              )}
            >
              <div className="flex items-center justify-between mb-3">
                <p className="font-semibold">{opt.label}</p>
                {isSelected && (
                  <span className="text-xs font-medium text-primary bg-primary/10 rounded-full px-2 py-0.5">
                    ✓ Selecionado
                  </span>
                )}
              </div>

              <div className="flex items-end justify-between gap-2">
                <div>
                  {opt.discount > 0 && (
                    <span className="inline-block text-xs font-medium text-primary bg-primary/10 rounded-full px-2 py-0.5 mb-1">
                      −{opt.discount}% no plano
                    </span>
                  )}
                  {savings > 0 && (
                    <p className="text-sm font-medium text-primary">
                      Economia de {formatPrice(savings * opt.cycles)} no total
                    </p>
                  )}
                </div>
                <div className="text-right shrink-0">
                  {savings > 0 && (
                    <p className="text-xs text-muted-foreground line-through">
                      {formatPrice(basePricing.paymentAmount)}/entrega
                    </p>
                  )}
                  <p className="text-xl font-bold text-primary leading-none">
                    {formatPrice(optPricing.paymentAmount)}
                  </p>
                  <p className="text-xs text-muted-foreground">por entrega</p>
                </div>
              </div>
            </button>
          )
        })}
      </div>

      <div className="mt-6 rounded-lg border p-4 space-y-2">
        <p className="text-sm font-medium">Renovação automática</p>
        <p className="text-xs text-muted-foreground">
          Ao usar todas as {totalCycles} entregas, um novo ciclo começa automaticamente.
          Você pode cancelar a renovação a qualquer momento.
        </p>
        <p className="text-xs text-primary">
          Próxima renovação estimada:{" "}
          <span className="font-semibold">
            {renewalDate.toLocaleDateString("pt-BR", { weekday: "long", day: "numeric", month: "long" })}
          </span>
        </p>
      </div>

      <div className="mt-4 rounded-lg bg-muted/50 p-4 space-y-2">
        <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Você tem controle total</p>
        <ul className="space-y-1.5 text-sm text-foreground/80">
          <li>Adie uma entrega se não puder receber</li>
          <li>Altere os dias e horários de entrega quando quiser</li>
          <li>Receba com menos frequência pulando entregas</li>
          <li>Altere os produtos da cesta a qualquer momento</li>
        </ul>
      </div>

      <div className="sticky bottom-0 border-t bg-background p-4">
        <Button className="w-full" onClick={onNext}>
          Continuar
        </Button>
      </div>
    </div>
  )
}
