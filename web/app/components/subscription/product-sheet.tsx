import { Sheet } from "~/components/ui/sheet"
import { formatPrice } from "./pricing"
import type { Product } from "./types"

type Props = {
  product: Product | null
  onClose: () => void
}

export function ProductSheet({ product, onClose }: Props) {
  return (
    <Sheet open={product !== null} onClose={onClose}>
      {product && (
        <div className="pb-8">
          {product.photo ? (
            <img
              src={product.photo}
              alt={product.name}
              className="w-full aspect-video object-cover"
            />
          ) : (
            <div className="w-full aspect-video bg-muted flex items-center justify-center">
              <span className="text-muted-foreground text-sm">Sem foto</span>
            </div>
          )}

          <div className="p-4 space-y-3">
            <h2 className="text-xl font-semibold">{product.name}</h2>
            <p className="text-2xl font-bold text-primary">{formatPrice(product.price)}</p>
            {product.unit && (
              <p className="text-sm text-muted-foreground">por {product.unit}</p>
            )}
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
                      A partir de {rule.quantityThreshold} unidades →{" "}
                      <span className="font-medium text-primary">{rule.discount}% off</span>
                    </p>
                  ))}
              </div>
            )}
          </div>
        </div>
      )}
    </Sheet>
  )
}
