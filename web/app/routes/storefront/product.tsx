import { useLoaderData, Link } from "react-router"
import type { Route } from "./+types/product"
import { apiClient } from "~/lib/api.server"

type Product = {
  id: string
  tenantId: string
  name: string
  description: string
  photo: string
  categoryId: string
  categoryName: string
  unit: string
  price: number
  isActive: boolean
}

export async function loader({ params }: Route.LoaderArgs) {
  const product = await apiClient().get<Product>(`/product/${params.id}`)
  return { product }
}

function formatPrice(cents: number): string {
  return (cents / 100).toLocaleString("pt-BR", { style: "currency", currency: "BRL" })
}

export default function ProductDetail() {
  const { product } = useLoaderData<typeof loader>()

  return (
    <div className="min-h-screen bg-background text-foreground">
      <header className="border-b px-6 py-4">
        <Link to="/" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
          ← Voltar
        </Link>
      </header>

      <main className="mx-auto max-w-2xl px-4 py-8">
        {product.photo ? (
          <img
            src={product.photo}
            alt={product.name}
            className="w-full rounded-xl aspect-video object-cover mb-6"
          />
        ) : (
          <div className="w-full rounded-xl aspect-video bg-muted flex items-center justify-center mb-6">
            <span className="text-muted-foreground text-sm">Sem foto</span>
          </div>
        )}

        <div className="space-y-3">
          {product.categoryName && (
            <p className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
              {product.categoryName}
            </p>
          )}
          <h1 className="text-2xl font-semibold">{product.name}</h1>
          <p className="text-2xl font-bold text-primary">{formatPrice(product.price)}</p>
          {product.unit && (
            <p className="text-sm text-muted-foreground">por {product.unit}</p>
          )}
          {product.description && (
            <p className="text-sm text-foreground/80 leading-relaxed pt-2">{product.description}</p>
          )}
        </div>
      </main>
    </div>
  )
}
