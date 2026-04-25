import { useState } from "react"
import { Button } from "~/components/ui/button"
import { ProductSheet } from "./product-sheet"
import { calculatePricing, formatPrice } from "./pricing"
import { cn } from "~/lib/utils"
import type { CycleDiscountRule, Product, ProductDiscountRule } from "./types"

type Props = {
  products: readonly Product[]
  items: { productId: string; quantity: number }[]
  totalCycles: number
  cycleDiscountRules: readonly CycleDiscountRule[]
  onChange: (items: { productId: string; quantity: number }[]) => void
  onNext: () => void
}

export function StepBasket({ products, items, totalCycles, cycleDiscountRules, onChange, onNext }: Props) {
  console.log("[StepBasket]")
  const [detailProduct, setDetailProduct] = useState<Product | null>(null)

  function getQuantity(productId: string) {
    return items.find((i) => i.productId === productId)?.quantity ?? 0
  }

  function setQuantity(productId: string, qty: number) {
    console.log("[StepBasket] setQuantity", productId, qty)
    if (qty <= 0) {
      onChange(items.filter((i) => i.productId !== productId))
    } else {
      const exists = items.find((i) => i.productId === productId)
      if (exists) {
        onChange(items.map((i) => (i.productId === productId ? { ...i, quantity: qty } : i)))
      } else {
        onChange([...items, { productId, quantity: qty }])
      }
    }
  }

  const pricing = calculatePricing(items, products, totalCycles, cycleDiscountRules)
  const hasItems = items.length > 0

  return (
    <>
      <div className="flex flex-col min-h-0">
        <div className="space-y-3 pb-36">
          {products.map((product) => {
            const qty = getQuantity(product.id)
            const rules = product.discountRules ?? []
            const activeRule = rules
              .filter((r) => qty >= r.quantityThreshold)
              .reduce<ProductDiscountRule | null>((best, r) => (!best || r.discount > best.discount ? r : best), null)
            const nextRule = rules
              .filter((r) => qty < r.quantityThreshold)
              .reduce<ProductDiscountRule | null>((closest, r) => (!closest || r.quantityThreshold < closest.quantityThreshold ? r : closest), null)

            return (
              <div key={product.id} className="flex items-center gap-3 rounded-xl border bg-card p-3">
                <button
                  type="button"
                  onClick={() => setDetailProduct(product)}
                  className="shrink-0"
                >
                  {product.photo ? (
                    <img
                      src={product.photo}
                      alt={product.name}
                      className="w-16 h-16 rounded-lg object-cover"
                    />
                  ) : (
                    <div className="w-16 h-16 rounded-lg bg-muted flex items-center justify-center">
                      <span className="text-xs text-muted-foreground">Foto</span>
                    </div>
                  )}
                </button>

                <div className="flex-1 min-w-0">
                  <button
                    type="button"
                    onClick={() => setDetailProduct(product)}
                    className="text-left"
                  >
                    <p className="text-sm font-medium leading-snug line-clamp-1">{product.name}</p>
                    <p className="text-sm font-semibold text-primary">{formatPrice(product.price)}</p>
                    {product.unit && (
                      <p className="text-xs text-muted-foreground">por {product.unit}</p>
                    )}
                  </button>

                  {activeRule && (
                    <span className="inline-block mt-1 text-xs font-medium text-primary bg-primary/10 rounded-full px-2 py-0.5">
                      {activeRule.discount}% off aplicado
                    </span>
                  )}
                  {!activeRule && nextRule && (
                    <p className="mt-1 text-xs text-muted-foreground">
                      Leve {nextRule.quantityThreshold - qty} mais → {nextRule.discount}% off
                    </p>
                  )}
                </div>

                <div className="flex items-center gap-2 shrink-0">
                  <button
                    type="button"
                    onClick={() => setQuantity(product.id, qty - 1)}
                    className="w-8 h-8 rounded-full border flex items-center justify-center text-lg font-medium hover:bg-muted transition-colors"
                  >
                    −
                  </button>
                  <span className="w-5 text-center text-sm font-semibold">{qty}</span>
                  <button
                    type="button"
                    onClick={() => setQuantity(product.id, qty + 1)}
                    className="w-8 h-8 rounded-full border flex items-center justify-center text-lg font-medium hover:bg-muted transition-colors"
                  >
                    +
                  </button>
                </div>
              </div>
            )
          })}
        </div>
      </div>

      {/* Sticky footer */}
      {/* <div className="fixed inset-x-0 bottom-0 border-t bg-background p-4 space-y-3">
        {hasItems && (
          <div className="flex justify-between items-center text-sm">
            <span className="text-muted-foreground">
              {items.reduce((s, i) => s + i.quantity, 0)} itens
            </span>
            <div className="text-right">
              {pricing.itemsDiscount > 0 && (
                <p className="text-xs text-muted-foreground line-through">
                  {formatPrice(pricing.itemsTotal)}
                </p>
              )}
              <p className="font-semibold text-primary">
                {formatPrice(pricing.itemsTotal - pricing.itemsDiscount)}
                <span className="text-xs font-normal text-muted-foreground"> /entrega</span>
              </p>
            </div>
          </div>
        )}
        <Button
          className="w-full"
          disabled={!hasItems}
          onClick={onNext}
        >
          {hasItems ? "Continuar" : "Adicione ao menos um produto"}
        </Button>
        <p className="text-center text-xs text-muted-foreground">
          Você pode alterar os produtos a qualquer momento
        </p>
      </div> */}

      {/* <ProductSheet product={detailProduct} onClose={() => setDetailProduct(null)} /> */}
    </>
  )
}
