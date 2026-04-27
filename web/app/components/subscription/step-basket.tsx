import { useState } from "react"
import { Button } from "~/components/ui/button"
import { formatPrice } from "./pricing"
import { cn } from "~/lib/utils"
import type { Product } from "./types"

type Props = {
  products: readonly Product[]
  items: { productId: string; quantity: number }[]
  onChange: (items: { productId: string; quantity: number }[]) => void
  onNext: () => void
}

export function StepBasket({ products, items, onChange, onNext }: Props) {
  const [productIndex, setProductIndex] = useState(0)

  const selectedProducts = items
    .map((item) => ({
      item,
      product: products.find((p) => p.id === item.productId),
    }))
    .filter((x): x is { item: typeof items[0]; product: Product } => x.product !== undefined)

  const current = selectedProducts[productIndex]
  if (!current) return null

  const { item, product } = current
  const qty = item.quantity
  const rules = (product.discountRules ?? []).slice().sort((a, b) => a.quantityThreshold - b.quantityThreshold)

  function setQuantity(qty: number) {
    if (qty < 1) return
    onChange(items.map((i) => (i.productId === product.id ? { ...i, quantity: qty } : i)))
  }

  function handleNext() {
    if (productIndex < selectedProducts.length - 1) {
      setProductIndex(productIndex + 1)
    } else {
      onNext()
    }
  }

  const isLast = productIndex === selectedProducts.length - 1
  const activeRule = rules.filter((r) => qty >= r.quantityThreshold).at(-1) ?? null
  const subtotal = product.price * qty
  const discountAmount = activeRule ? Math.floor((subtotal * activeRule.discount) / 100) : 0

  return (
    <div className="flex flex-col">
      {/* Progress within this step */}
      <div className="flex gap-1 mb-6">
        {selectedProducts.map((_, i) => (
          <div
            key={i}
            className={cn(
              "h-1 flex-1 rounded-full transition-colors",
              i <= productIndex ? "bg-primary" : "bg-muted"
            )}
          />
        ))}
      </div>

      <p className="text-xs text-muted-foreground mb-1">
        Produto {productIndex + 1} de {selectedProducts.length}
      </p>

      {/* Product card */}
      <div className="rounded-xl border bg-card p-4 mb-6">
        <div className="flex gap-4">
          {product.photo ? (
            <img src={product.photo} alt={product.name} className="w-20 h-20 rounded-lg object-cover shrink-0" />
          ) : (
            <div className="w-20 h-20 rounded-lg bg-muted flex items-center justify-center shrink-0">
              <span className="text-xs text-muted-foreground">Foto</span>
            </div>
          )}
          <div className="flex-1 min-w-0">
            <p className="font-semibold leading-snug">{product.name}</p>
            {product.description && (
              <p className="text-xs text-muted-foreground mt-1 line-clamp-2">{product.description}</p>
            )}
            <p className="text-sm font-semibold text-primary mt-1">{formatPrice(product.price)}</p>
            {product.unit && <p className="text-xs text-muted-foreground">por {product.unit}</p>}
          </div>
        </div>

        {/* Quantity selector */}
        <div className="mt-4 flex items-center justify-between">
          <p className="text-sm text-muted-foreground">Quantidade</p>
          <div className="flex items-center gap-4">
            <button
              type="button"
              onClick={() => setQuantity(qty - 1)}
              disabled={qty <= 1}
              className="w-10 h-10 rounded-full border flex items-center justify-center text-xl font-medium hover:bg-muted transition-colors disabled:opacity-30"
            >
              −
            </button>
            <span className="w-6 text-center text-lg font-bold">{qty}</span>
            <button
              type="button"
              onClick={() => setQuantity(qty + 1)}
              className="w-10 h-10 rounded-full border flex items-center justify-center text-xl font-medium hover:bg-muted transition-colors"
            >
              +
            </button>
          </div>
        </div>

        {/* Subtotal */}
        <div className="mt-3 flex items-baseline justify-between">
          <p className="text-sm text-muted-foreground">Subtotal</p>
          <div className="text-right">
            {discountAmount > 0 && (
              <p className="text-xs text-muted-foreground line-through">{formatPrice(subtotal)}</p>
            )}
            <p className="font-bold text-primary">{formatPrice(subtotal - discountAmount)}</p>
          </div>
        </div>
      </div>

      {/* Discount tiers */}
      {rules.length > 0 && (
        <div className="space-y-2">
          <p className="text-sm font-medium">Descontos por quantidade</p>
          {rules.map((rule) => {
            const isActive = activeRule?.quantityThreshold === rule.quantityThreshold
            const nextUnlocked = rules.find((r) => r.quantityThreshold > qty)

            return (
              <button
                key={rule.quantityThreshold}
                type="button"
                onClick={() => setQuantity(rule.quantityThreshold)}
                className={cn(
                  "w-full flex items-center justify-between rounded-lg border px-4 py-3 text-sm transition-colors text-left",
                  isActive ? "border-primary bg-primary/5" : "border-border bg-card hover:bg-muted"
                )}
              >
                <div className="flex items-center gap-2">
                  {isActive ? (
                    <span className="w-5 h-5 rounded-full bg-primary flex items-center justify-center text-white text-xs shrink-0">✓</span>
                  ) : (
                    <span className="w-5 h-5 rounded-full border-2 border-muted-foreground/30 shrink-0" />
                  )}
                  <span className={cn("font-medium", isActive ? "text-primary" : "text-foreground")}>
                    {rule.quantityThreshold}+ {product.unit ?? "unidades"}
                  </span>
                  {nextUnlocked?.quantityThreshold === rule.quantityThreshold && (
                    <span className="text-xs text-muted-foreground">
                      (mais {rule.quantityThreshold - qty} para desbloquear)
                    </span>
                  )}
                </div>
                <span className={cn("font-semibold shrink-0", isActive ? "text-primary" : "text-muted-foreground")}>
                  {rule.discount}% off
                </span>
              </button>
            )
          })}
        </div>
      )}

      <div className="sticky bottom-0 border-t bg-background p-4">
        <Button className="w-full" onClick={handleNext}>
          {isLast ? "Continuar" : "Próximo produto"}
        </Button>
      </div>
    </div>
  )
}
