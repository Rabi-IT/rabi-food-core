import { useState, useCallback, useEffect } from "react"

export type WizardStep = 1 | 2 | 3 | 4 | 5 | 6 | "summary"

export type WizardData = {
  items: { productId: string; quantity: number }[]
  deliveryDays: { weekday: number; startHour: number; endHour: number }[]
  totalCycles: number
  autoRenew: boolean
  user: { name: string; email: string; phone: string; token: string }
  address: {
    street: string
    complement: string
    neighborhood: string
    city: string
    state: string
    zip: string
  }
}

const defaultData: WizardData = {
  items: [],
  deliveryDays: [],
  totalCycles: 1,
  autoRenew: false,
  user: { name: "", email: "", phone: "", token: "" },
  address: { street: "", complement: "", neighborhood: "", city: "", state: "", zip: "" },
}

export function useWizard() {
  const [step, setStep] = useState<WizardStep>(1)
  const [data, setDataRaw] = useState<WizardData>(defaultData)

  const setData = useCallback((update: Partial<WizardData>) => {
    setDataRaw((prev) => {
      const next = { ...prev, ...update }
      return next
    })
  }, [])

  const reset = useCallback(() => {
    setDataRaw(defaultData)
    setStep(1)
  }, [])

  return { step, setStep, data, setData, reset }
}
