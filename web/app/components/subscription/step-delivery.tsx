import { Button } from "~/components/ui/button"
import { cn } from "~/lib/utils"

const WEEKDAYS = ["Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"]
const WEEKDAYS_FULL = ["Domingo", "Segunda", "Terça", "Quarta", "Quinta", "Sexta", "Sábado"]

const TIME_WINDOWS = [
  { label: "Manhã", description: "8h – 12h", startHour: 8, endHour: 12 },
  { label: "Tarde", description: "12h – 18h", startHour: 12, endHour: 18 },
  { label: "Dia todo", description: "8h – 18h", startHour: 8, endHour: 18 },
]

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
      onChange([...deliveryDays, { weekday, startHour: 8, endHour: 18 }])
    }
  }

  function setWindow(weekday: number, startHour: number, endHour: number) {
    onChange(deliveryDays.map((d) => d.weekday === weekday ? { ...d, startHour, endHour } : d))
  }

  const selected = new Set(deliveryDays.map((d) => d.weekday))
  const hasSelection = deliveryDays.length > 0
  const allConfigured = deliveryDays.every((d) => d.startHour !== undefined)

  const sorted = deliveryDays.slice().sort((a, b) => a.weekday - b.weekday)

  return (
    <div className="flex flex-col pb-32 space-y-6">
      <div>
        <p className="text-sm text-muted-foreground mb-4">
          Escolha os dias da semana em que deseja receber suas entregas.
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
      </div>

      {hasSelection && (
        <div className="space-y-4">
          <p className="text-sm text-muted-foreground">
            Escolha o horário preferido para cada dia.
          </p>
          {sorted.map((day) => {
            const active = TIME_WINDOWS.find(
              (w) => w.startHour === day.startHour && w.endHour === day.endHour
            )
            return (
              <div key={day.weekday} className="space-y-2">
                <p className="text-sm font-medium">{WEEKDAYS_FULL[day.weekday]}</p>
                <div className="grid grid-cols-3 gap-2">
                  {TIME_WINDOWS.map((w) => {
                    const isActive = active?.label === w.label
                    return (
                      <button
                        key={w.label}
                        type="button"
                        onClick={() => setWindow(day.weekday, w.startHour, w.endHour)}
                        className={cn(
                          "rounded-lg border p-2 text-center text-xs transition-colors",
                          isActive
                            ? "border-primary bg-primary/5 text-primary font-medium"
                            : "border-border bg-card hover:bg-muted"
                        )}
                      >
                        <p className="font-medium">{w.label}</p>
                        <p className="text-muted-foreground">{w.description}</p>
                      </button>
                    )
                  })}
                </div>
              </div>
            )
          })}
        </div>
      )}

      <div className="fixed inset-x-0 bottom-0 border-t bg-background p-4">
        <Button className="w-full" disabled={!hasSelection || !allConfigured} onClick={onNext}>
          {!hasSelection ? "Selecione ao menos um dia" : "Continuar"}
        </Button>
      </div>
    </div>
  )
}
