import { Button } from "~/components/ui/button"
import { cn } from "~/lib/utils"

const WEEKDAYS = ["Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"]
const BUSINESS_HOURS = { startHour: 8, endHour: 18 }

type DeliveryDay = { weekday: number; startHour: number; endHour: number }

type Props = {
  deliveryDays: DeliveryDay[]
  onChange: (days: DeliveryDay[]) => void
  onNext: () => void
}

export function StepDelivery({ deliveryDays, onChange, onNext }: Props) {
  function toggle(weekday: number) {
    const exists = deliveryDays.find((d) => d.weekday === weekday)
    if (exists) {
      onChange(deliveryDays.filter((d) => d.weekday !== weekday))
    } else {
      onChange([...deliveryDays, { weekday, ...BUSINESS_HOURS }])
    }
  }

  const selected = new Set(deliveryDays.map((d) => d.weekday))
  const hasSelection = deliveryDays.length > 0

  return (
    <div className="flex flex-col pb-32">
      <p className="text-sm text-muted-foreground mb-6">
        Escolha os dias da semana em que deseja receber suas entregas. As entregas acontecem em horário comercial (8h–18h).
      </p>

      <div className="grid grid-cols-7 gap-2">
        {WEEKDAYS.map((label, weekday) => (
          <button
            key={weekday}
            type="button"
            onClick={() => toggle(weekday)}
            className={cn(
              "flex flex-col items-center gap-1 rounded-xl py-3 border text-xs font-medium transition-colors",
              selected.has(weekday)
                ? "bg-primary text-primary-foreground border-primary"
                : "bg-background border-border hover:bg-muted"
            )}
          >
            <span>{label}</span>
          </button>
        ))}
      </div>

      {hasSelection && (
        <div className="mt-6 rounded-lg bg-muted p-3">
          <p className="text-xs font-medium text-muted-foreground mb-1">Dias selecionados</p>
          <p className="text-sm font-medium">
            {deliveryDays
              .slice()
              .sort((a, b) => a.weekday - b.weekday)
              .map((d) => WEEKDAYS[d.weekday])
              .join(", ")}
          </p>
        </div>
      )}

      <div className="fixed inset-x-0 bottom-0 border-t bg-background p-4">
        <Button className="w-full" disabled={!hasSelection} onClick={onNext}>
          {hasSelection ? "Continuar" : "Selecione ao menos um dia"}
        </Button>
      </div>
    </div>
  )
}
