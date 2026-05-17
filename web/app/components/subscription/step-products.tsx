import { useTranslation } from "react-i18next"
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

export function StepProducts({ products, items, onChange, onNext }: Props) {
  const { t } = useTranslation()
  const selectedIds = new Set(items.map((i) => i.productId))

  function toggle(productId: string) {
    if (selectedIds.has(productId)) {
      onChange(items.filter((i) => i.productId !== productId))
    } else {
      onChange([...items, { productId, quantity: 1 }])
    }
  }

  const hasItems = items.length > 0

  return (
    <div className="flex flex-col">
      <p className="text-sm text-muted-foreground mb-6">
        {t("subscription.products.hint")}
      </p>

      <div className="space-y-3">
        {products.map((product) => {
          const isSelected = selectedIds.has(product.id)
          return (
            <button
              key={product.id}
              type="button"
              onClick={() => toggle(product.id)}
              className={cn(
                "w-full flex items-center gap-3 rounded-xl border p-3 text-left transition-colors",
                isSelected ? "border-primary bg-primary/5" : "border-border bg-card hover:bg-muted"
              )}
            >
              {product.photo ? (
                <img
                  src={product.photo}
                  alt={product.name}
                  className="w-16 h-16 rounded-lg object-cover shrink-0"
                />
              ) : (
                <div className="w-16 h-16 rounded-lg bg-muted flex items-center justify-center shrink-0">
                  <span className="text-xs text-muted-foreground">{t("subscription.product.photoPlaceholder")}</span>
                </div>
              )}

              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium leading-snug line-clamp-1">{product.name}</p>
                <p className="text-sm font-semibold text-primary">{formatPrice(product.price)}</p>
                {product.unit && (
                  <p className="text-xs text-muted-foreground">{t("subscription.product.perUnit", { unit: product.unit })}</p>
                )}
                {product.discountRules && product.discountRules.length > 0 && (
                  <p className="text-xs text-primary mt-0.5">
                    {t("subscription.products.quantityDiscountAvailable")}
                  </p>
                )}
              </div>

              <div
                className={cn(
                  "w-6 h-6 rounded-full border-2 flex items-center justify-center shrink-0 transition-colors",
                  isSelected ? "border-primary bg-primary text-white" : "border-border"
                )}
              >
                {isSelected && <span className="text-xs font-bold">✓</span>}
              </div>
            </button>
          )
        })}
      </div>

      <div className="sticky bottom-0 border-t bg-background p-4 space-y-2">
        {hasItems && (
          <p className="text-center text-sm text-muted-foreground">
            {items.length === 1
              ? t("subscription.products.selectedOne")
              : t("subscription.products.selectedMany", { count: items.length })}
          </p>
        )}
        <Button className="w-full" disabled={!hasItems} onClick={onNext}>
          {hasItems ? t("subscription.products.continue") : t("subscription.products.selectAtLeastOne")}
        </Button>
      </div>
    </div>
  )
}
