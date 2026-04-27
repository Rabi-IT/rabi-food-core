import { Sheet } from "~/components/ui/sheet"
import { formatPrice } from "./pricing"
import type { Product } from "./types"

type Props = {
  readonly product: Product | null
  readonly onClose: () => void
  readonly onSubscribe: (productId: string) => void
}

export function ProductSheet({ product, onClose, onSubscribe }: Props) {
  function handleSubscribe(e: React.MouseEvent) {
    e.stopPropagation()
    if (!product) return
    const id = product.id
    onClose()
    setTimeout(() => onSubscribe(id), 120)
  }

  return (
    <Sheet open={product !== null} onClose={onClose}>
      {product && (
        <div className="p-4 space-y-4">
          <div className="flex gap-4 items-start">
            {product.photo ? (
              <img
                src={product.photo}
                alt={product.name}
                className="w-24 h-24 rounded-xl object-cover shrink-0"
              />
            ) : (
              <div className="w-24 h-24 rounded-xl bg-muted flex items-center justify-center shrink-0">
                <span className="text-muted-foreground text-xs">Sem foto</span>
              </div>
            )}
            <div className="flex-1 min-w-0 pt-1">
              <h2 className="text-base font-semibold leading-snug">{product.name}</h2>
              <p className="text-xl font-bold text-primary mt-1">{formatPrice(product.price)}</p>
              {product.unit && (
                <p className="text-xs text-muted-foreground">por {product.unit}</p>
              )}
            </div>
          </div>

          {product.description && (
            <p className="text-sm text-foreground/80 leading-relaxed">{product.description}</p>
          )}

          {(product.discountRules ?? []).length > 0 && (
            <div className="rounded-lg bg-primary/5 border border-primary/20 p-3 space-y-1">
              <p className="text-xs font-semibold text-primary uppercase tracking-wide">
                Descontos por quantidade
              </p>
              {product.discountRules
                .slice()
                .sort((a, b) => a.quantityThreshold - b.quantityThreshold)
                .map((rule) => (
                  <p key={rule.quantityThreshold} className="text-sm text-foreground">
                    A partir de {rule.quantityThreshold} {product.unit ?? "unidades"} →{" "}
                    <span className="font-medium text-primary">{rule.discount}% off</span>
                  </p>
                ))}
            </div>
          )}

          <button
            type="button"
            onClick={handleSubscribe}
            className="w-full rounded-lg bg-primary px-4 py-2.5 text-sm font-semibold text-primary-foreground hover:bg-primary/90 transition-colors"
          >
            Assinar este produto
          </button>
        </div>
      )}
    </Sheet>
  )
}
