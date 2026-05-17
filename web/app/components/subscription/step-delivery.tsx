import { useTranslation } from "react-i18next"
import { Button } from "~/components/ui/button"
import { cn } from "~/lib/utils"

const TIME_WINDOW_IDS = [
  { id: "morning", startHour: 8, endHour: 12 },
  { id: "afternoon", startHour: 12, endHour: 18 },
  { id: "allDay", startHour: 8, endHour: 18 },
] as const

type TimeWindowId = typeof TIME_WINDOW_IDS[number]["id"]
type DeliveryDay = { weekday: number; startHour: number; endHour: number }

type Props = {
  readonly deliveryDays: readonly DeliveryDay[]
  readonly orderLeadMinutes: number
  readonly onChange: (days: readonly DeliveryDay[]) => void
  readonly onNext: () => void
}

function getNextDeliveryDate(weekday: number, startHour: number, orderLeadMinutes: number): Date {
  const now = new Date()
  const cursor = new Date(now)
  cursor.setHours(0, 0, 0, 0)

  for (let i = 0; i <= 7; i++) {
    if (cursor.getDay() === weekday) {
      const deliveryTime = new Date(cursor)
      deliveryTime.setHours(startHour, 0, 0, 0)
      const leadMs = orderLeadMinutes * 60 * 1000
      if (deliveryTime.getTime() - now.getTime() >= leadMs) {
        return cursor
      }
    }
    cursor.setDate(cursor.getDate() + 1)
  }

  return cursor
}

export function StepDelivery({ deliveryDays, orderLeadMinutes, onChange, onNext }: Props) {
  const { t } = useTranslation()

  const weekdays = t("subscription.delivery.weekdays", { returnObjects: true }) as readonly string[]
  const weekdaysFull = t("subscription.delivery.weekdaysFull", { returnObjects: true }) as readonly string[]

  function formatNextDelivery(date: Date): string {
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    const tomorrow = new Date(today)
    tomorrow.setDate(today.getDate() + 1)

    if (date.getTime() === today.getTime()) return t("subscription.delivery.today")
    if (date.getTime() === tomorrow.getTime()) return t("subscription.delivery.tomorrow")
    return date.toLocaleDateString("pt-BR", { weekday: "long", day: "numeric", month: "long" })
  }

  function toggle(weekday: number) {
    const exists = deliveryDays.find((d) => d.weekday === weekday)
    if (exists) {
      onChange(deliveryDays.filter((d) => d.weekday !== weekday))
      return
    }
    onChange([...deliveryDays, { weekday, startHour: 8, endHour: 18 }])
  }

  function setWindow(weekday: number, startHour: number, endHour: number) {
    onChange(deliveryDays.map((d) => d.weekday === weekday ? { ...d, startHour, endHour } : d))
  }

  const selected = new Set(deliveryDays.map((d) => d.weekday))
  const hasSelection = deliveryDays.length > 0
  const allConfigured = deliveryDays.every((d) => d.startHour !== undefined)
  const sorted = deliveryDays.slice().sort((a, b) => a.weekday - b.weekday)

  return (
    <div className="flex flex-col space-y-6">
      <div>
        <p className="text-sm text-muted-foreground mb-4">
          {t("subscription.delivery.hint")}
        </p>
        <div className="grid grid-cols-7 gap-2">
          {weekdays.map((label, weekday) => (
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
            {t("subscription.delivery.timeHint")}
          </p>
          {sorted.map((day) => {
            const active = TIME_WINDOW_IDS.find(
              (w) => w.startHour === day.startHour && w.endHour === day.endHour
            )

            const nextDelivery = getNextDeliveryDate(day.weekday, day.startHour, orderLeadMinutes)
            const nextLabel = formatNextDelivery(nextDelivery)

            return (
              <div key={day.weekday} className="space-y-2">
                <div className="flex items-baseline justify-between">
                  <p className="text-sm font-medium">{weekdaysFull[day.weekday]}</p>
                  {active && (
                    <p className="text-xs text-muted-foreground">
                      {t("subscription.delivery.nextDelivery")} <span className="font-medium text-foreground">{nextLabel}</span>
                    </p>
                  )}
                </div>
                <div className="grid grid-cols-3 gap-2">
                  {TIME_WINDOW_IDS.map((w) => {
                    const isActive = active?.id === w.id
                    const windowId = w.id as TimeWindowId
                    return (
                      <button
                        key={w.id}
                        type="button"
                        onClick={() => setWindow(day.weekday, w.startHour, w.endHour)}
                        className={cn(
                          "rounded-lg border p-2 text-center text-xs transition-colors",
                          isActive
                            ? "border-primary bg-primary/5 text-primary font-medium"
                            : "border-border bg-card hover:bg-muted"
                        )}
                      >
                        <p className="font-medium">{t(`subscription.delivery.${windowId}.label`)}</p>
                        <p className="text-muted-foreground">{t(`subscription.delivery.${windowId}.description`)}</p>
                      </button>
                    )
                  })}
                </div>
              </div>
            )
          })}
        </div>
      )}

      <div className="sticky bottom-0 border-t bg-background p-4">
        <Button className="w-full" disabled={!hasSelection || !allConfigured} onClick={onNext}>
          {hasSelection ? t("subscription.delivery.continue") : t("subscription.delivery.selectAtLeastOne")}
        </Button>
      </div>
    </div>
  )
}
