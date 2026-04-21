import { useLoaderData, useRouteLoaderData, Link, useSearchParams } from "react-router"
import type { Route } from "./+types/storefront-page"
import { resolveTenant } from "~/lib/tenant.server"
import { apiClient } from "~/lib/api.server"
import { cn } from "~/lib/utils"
import InstitutionalPage from "./institutional-page"

type Category = { id: string; name: string; description: string }
type Product = {
  id: string
  name: string
  description: string
  photo: string
  categoryId: string
  unit: string
  price: number
  isActive: boolean
}
type PaginatedCategories = { data: Category[]; maxPages: number }
type PaginatedProducts = { data: Product[]; maxPages: number }

export async function loader({ request }: Route.LoaderArgs) {
  const tenant = await resolveTenant(request)
  if (!tenant) return { tenant: null, categories: [] as Category[], products: [] as Product[], selectedCategoryId: null }

  const url = new URL(request.url)
  const categoryId = url.searchParams.get("categoryId") ?? undefined
  const name = url.searchParams.get("name") ?? undefined

  const api = apiClient(tenant.id)

  const [categories, products] = await Promise.all([
    api.get<PaginatedCategories>("/category/"),
    api.get<PaginatedProducts>("/product/", {
      isActive: "1",
      categoryId,
      name,
    }),
  ])

  return { tenant, categories: categories.data, products: products.data, selectedCategoryId: categoryId ?? null }
}

function formatPrice(cents: number): string {
  return (cents / 100).toLocaleString("pt-BR", { style: "currency", currency: "BRL" })
}

export default function StorefrontHome() {
  const { tenant, categories, products, selectedCategoryId } = useLoaderData<typeof loader>()
  if (!tenant) return <InstitutionalPage />
  const layout = useRouteLoaderData("routes/storefront/layout") as { tenant: { name: string } }
  const [searchParams, setSearchParams] = useSearchParams()

  function selectCategory(id: string | null) {
    const next = new URLSearchParams(searchParams)
    if (id) next.set("categoryId", id)
    else next.delete("categoryId")
    setSearchParams(next, { preventScrollReset: true })
  }

  return (
    <div className="min-h-screen bg-background text-foreground">
      <header className="border-b px-6 py-4">
        <h1 className="text-xl font-semibold">{layout?.tenant.name}</h1>
      </header>

      <main className="mx-auto max-w-6xl px-4 py-6 space-y-6">
        {/* Category filter */}
        <div className="flex gap-2 overflow-x-auto pb-1">
          <button
            onClick={() => selectCategory(null)}
            className={cn(
              "shrink-0 rounded-full border px-4 py-1.5 text-sm font-medium transition-colors",
              selectedCategoryId === null
                ? "bg-primary text-primary-foreground border-primary"
                : "bg-background hover:bg-muted border-border"
            )}
          >
            Todos
          </button>
          {categories.map((cat) => (
            <button
              key={cat.id}
              onClick={() => selectCategory(cat.id)}
              className={cn(
                "shrink-0 rounded-full border px-4 py-1.5 text-sm font-medium transition-colors",
                selectedCategoryId === cat.id
                  ? "bg-primary text-primary-foreground border-primary"
                  : "bg-background hover:bg-muted border-border"
              )}
            >
              {cat.name}
            </button>
          ))}
        </div>

        {/* Product grid */}
        {products.length === 0 ? (
          <p className="text-center text-muted-foreground py-16">Nenhum produto encontrado.</p>
        ) : (
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
            {products.map((product) => (
              <Link
                key={product.id}
                to={`/product/${product.id}`}
                className="group rounded-xl border bg-card overflow-hidden hover:shadow-md transition-shadow"
              >
                {product.photo ? (
                  <img
                    src={product.photo}
                    alt={product.name}
                    className="aspect-square w-full object-cover group-hover:opacity-90 transition-opacity"
                  />
                ) : (
                  <div className="aspect-square w-full bg-muted flex items-center justify-center">
                    <span className="text-muted-foreground text-xs">Sem foto</span>
                  </div>
                )}
                <div className="p-3 space-y-1">
                  <p className="text-sm font-medium leading-snug line-clamp-2">{product.name}</p>
                  <p className="text-sm font-semibold text-primary">{formatPrice(product.price)}</p>
                  {product.unit && (
                    <p className="text-xs text-muted-foreground">por {product.unit}</p>
                  )}
                </div>
              </Link>
            ))}
          </div>
        )}
      </main>
    </div>
  )
}
