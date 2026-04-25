import type { CycleDiscountRule, Product, ProductDiscountRule } from "./types"

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
  paymentAmount: number
}

export function calculatePricing(
  selectedItems: { productId: string; quantity: number }[],
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
    const discountPercent = (product.discountRules ?? []).reduce((best, rule) => {
      return quantity >= rule.quantityThreshold && rule.discount > best ? rule.discount : best
    }, 0)
    const discountAmount = Math.floor((subtotal * discountPercent) / 100)

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

  const net = itemsTotal - itemsDiscount
  const cycleDiscountPercent = cycleDiscountRules.reduce((best, rule) => {
    return totalCycles >= rule.cyclesThreshold && rule.discount > best ? rule.discount : best
  }, 0)
  const cycleDiscount = Math.floor((net * cycleDiscountPercent) / 100)

  return {
    items,
    itemsTotal,
    itemsDiscount,
    cycleDiscountPercent,
    cycleDiscount,
    paymentAmount: net - cycleDiscount,
  }
}

export function formatPrice(cents: number): string {
  return (cents / 100).toLocaleString("pt-BR", { style: "currency", currency: "BRL" })
}
