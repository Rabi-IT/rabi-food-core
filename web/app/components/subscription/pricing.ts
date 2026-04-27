import type { CycleDiscountRule, Product } from "./types"

export type ItemBreakdown = {
  productId: string
  name: string
  quantity: number
  unitPrice: number
  subtotal: number
  discountPercent: number
  discountAmount: number
  total: number
}

export type PricingResult = {
  items: ItemBreakdown[]
  itemsTotal: number
  itemsDiscount: number
  cycleDiscountPercent: number
  cycleDiscount: number
  totalDiscount: number
  paymentAmount: number
}

// Mirrors buildSubscriptionPricing in api/features/subscription/usecases/build_subscription_pricing.go
export function calculatePricing(
  selectedItems: readonly { productId: string; quantity: number }[],
  products: readonly Product[],
  totalCycles: number,
  cycleDiscountRules: readonly CycleDiscountRule[],
): PricingResult {
  const productMap = new Map(products.map((p) => [p.id, p]))

  let itemsTotal = 0
  let itemsDiscount = 0
  const items: ItemBreakdown[] = []

  for (const { productId, quantity } of selectedItems) {
    const product = productMap.get(productId)
    if (!product) continue

    const subtotal = product.price * quantity

    let discountPercent = 0
    for (const rule of product.discountRules ?? []) {
      if (quantity >= rule.quantityThreshold && rule.discount > discountPercent) {
        discountPercent = rule.discount
      }
    }

    const discountAmount = Math.trunc((subtotal * discountPercent) / 100)

    items.push({
      productId,
      name: product.name,
      quantity,
      unitPrice: product.price,
      subtotal,
      discountPercent,
      discountAmount,
      total: subtotal - discountAmount,
    })

    itemsTotal += subtotal
    itemsDiscount += discountAmount
  }

  let cycleDiscountPercent = 0
  for (const rule of cycleDiscountRules) {
    if (totalCycles >= rule.cyclesThreshold && rule.discount > cycleDiscountPercent) {
      cycleDiscountPercent = rule.discount
    }
  }

  const net = itemsTotal - itemsDiscount
  const cycleDiscount = Math.trunc((net * cycleDiscountPercent) / 100)

  return {
    items,
    itemsTotal,
    itemsDiscount,
    cycleDiscountPercent,
    cycleDiscount,
    totalDiscount: itemsDiscount + cycleDiscount,
    paymentAmount: net - cycleDiscount,
  }
}

export function formatPrice(cents: number): string {
  return (cents / 100).toLocaleString("pt-BR", { style: "currency", currency: "BRL" })
}
