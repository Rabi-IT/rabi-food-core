import { useEffect } from "react"
import { Button } from "~/components/ui/button"
import { calculatePricing, formatPrice } from "./pricing"
import { cn } from "~/lib/utils"
import type { CycleDiscountRule, Product } from "./types"

type Props = {
  products: readonly Product[]
  items: { productId: string; quantity: number }[]
  totalCycles: number
  minCycles: number
  deliveryWeekdays: number[]
  cycleDiscountRules: readonly CycleDiscountRule[]
  autoRenew: boolean
  onChange: (totalCycles: number, autoRenew: boolean) => void
  onNext: () => void
}

function getRenewalDate(weekdays: number[], totalCycles: number): Date {
  const sorted = [...weekdays].sort((a, b) => a - b)
  const cursor = new Date()
  cursor.setHours(0, 0, 0, 0)
  cursor.setDate(cursor.getDate() + 1) // start from tomorrow
  let count = 0
  while (count < totalCycles) {
    if (sorted.includes(cursor.getDay())) count++
    if (count < totalCycles) cursor.setDate(cursor.getDate() + 1)
  }
  return cursor
}

export function StepPlan({ products, items, totalCycles, minCycles, deliveryWeekdays, cycleDiscountRules, autoRenew, onChange, onNext }: Props) {
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
    if (options.length === 1) onChange(options[0].cycles, autoRenew)
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div className="flex flex-col pb-32">
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
              onClick={() => onChange(opt.cycles, autoRenew)}
              className={cn(
                "w-full rounded-xl border p-4 text-left transition-colors",
                isSelected ? "border-primary bg-primary/5" : "border-border bg-card hover:bg-muted"
              )}
            >
              {/* Header: label + selected */}
              <div className="flex items-center justify-between mb-3">
                <p className="font-semibold">{opt.label}</p>
                {isSelected && (
                  <span className="text-xs font-medium text-primary bg-primary/10 rounded-full px-2 py-0.5">
                    ✓ Selecionado
                  </span>
                )}
              </div>

              {/* Price row */}
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

      <div className="mt-6 rounded-lg border p-4 space-y-3">
        <div className="flex items-center justify-between gap-3">
          <div>
            <p className="text-sm font-medium">Renovação automática</p>
            <p className="text-xs text-muted-foreground">
              Ao usar todas as {totalCycles} entregas, um novo ciclo é gerado automaticamente.
            </p>
          </div>
          <button
            type="button"
            role="switch"
            aria-checked={autoRenew}
            onClick={() => onChange(totalCycles, !autoRenew)}
            className={cn(
              "relative w-11 h-6 rounded-full transition-colors shrink-0",
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

        {autoRenew && (
          <p className="text-xs text-primary border-t pt-3">
            Sua próxima cobrança será em{" "}
            <span className="font-semibold">
              {getRenewalDate(deliveryWeekdays, totalCycles).toLocaleDateString("pt-BR", {
                weekday: "long", day: "numeric", month: "long",
              })}
            </span>.
          </p>
        )}
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

      <div className="fixed inset-x-0 bottom-0 border-t bg-background p-4">
        <Button className="w-full" onClick={onNext}>
          Continuar
        </Button>
      </div>
    </div>
  )
}
