export type ProductDiscountRule = { quantityThreshold: number; discount: number }
export type CycleDiscountRule = { cyclesThreshold: number; discount: number }


export type Product = {
  readonly id: string
  readonly name: string
  readonly description: string
  readonly photo: string
  readonly unit: string
  readonly price: number
  readonly discountRules: readonly ProductDiscountRule[]
}

export type SubscriptionConfig = {
  readonly discountRules: readonly CycleDiscountRule[]
  readonly cutoffOffsetMinutes: number
  readonly orderLeadMinutes: number
  readonly isOpen: boolean
}