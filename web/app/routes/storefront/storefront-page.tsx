import { useLoaderData, useRouteLoaderData, Link } from "react-router"
import { useEffect, useRef, useState } from "react"
import type { Route } from "./+types/storefront-page"
import { resolveTenant } from "~/lib/tenant.server"
import { apiClient } from "~/lib/api.server"
import { cn } from "~/lib/utils"
import { ProductSheet } from "~/components/subscription/product-sheet"
import InstitutionalPage from "./institutional-page"

type Product = {
  readonly id: string
  readonly name: string
  readonly description: string
  readonly photo: string
  readonly categoryId: string
  readonly unit: string
  readonly price: number
  readonly discountRules: readonly { quantityThreshold: number; discount: number }[]
}

type Category = { readonly id: string; readonly name: string }

// Fraction of a category's own height that is considered "dead zone" —
// when only this much of the current category remains visible above the next one,
// we already consider the next category active.
const DEAD_ZONE_RATIO = 0.3

type CatalogResponse = {
  readonly products: readonly Product[]
  readonly categories: readonly Category[]
}

export async function loader({ request }: Route.LoaderArgs) {
  const tenant = await resolveTenant(request)
  if (!tenant) return { tenant: null, categories: [] as Category[], products: [] as Product[] }

  const api = apiClient(tenant.id)
  const res = await api.get<CatalogResponse>("/product/catalog")

  return {
    tenant,
    products: res.ok ? res.data.products : [],
    categories: res.ok ? res.data.categories : [],
  }
}

function formatPrice(cents: number): string {
  return (cents / 100).toLocaleString("pt-BR", { style: "currency", currency: "BRL" })
}

export default function StorefrontHome() {
  const { tenant, categories, products } = useLoaderData<typeof loader>()
  if (!tenant) return <InstitutionalPage />

  const layout = useRouteLoaderData("routes/storefront/layout") as { tenant: { name: string } }
  const [activeCategoryId, setActiveCategoryId] = useState<string | null>(categories[0]?.id ?? null)
  const [detailProduct, setDetailProduct] = useState<Product | null>(null)
  const [stickyHeight, setStickyHeight] = useState(0)
  const [titleGap, setTitleGap] = useState(0)
  const [bottomPadding, setBottomPadding] = useState(0)

  const categoryBarRef = useRef<HTMLDivElement>(null)
  const mainRef = useRef<HTMLDivElement>(null)
  const firstTitleRef = useRef<HTMLHeadingElement>(null)

  const productsByCategory = categories.map((cat) => ({
    category: cat,
    products: products.filter((p) => p.categoryId === cat.id),
  }))

  useEffect(() => {
    function calculate() {
      const bar = categoryBarRef.current
      const main = mainRef.current
      if (!bar || !main) return

      const barHeight = bar.getBoundingClientRect().height
      setStickyHeight(barHeight)

      if (firstTitleRef.current) {
        const mb = parseFloat(getComputedStyle(firstTitleRef.current).marginBottom)
        setTitleGap(mb)
      }

      const sections = main.querySelectorAll("section")
      const lastSection = sections[sections.length - 1]
      if (!lastSection) return

      const mb = firstTitleRef.current
        ? parseFloat(getComputedStyle(firstTitleRef.current).marginBottom)
        : 0
      const available = window.innerHeight - barHeight - mb
      const lastHeight = lastSection.getBoundingClientRect().height
      setBottomPadding(Math.max(0, available - lastHeight))
    }

    calculate()
    window.addEventListener("resize", calculate)
    return () => window.removeEventListener("resize", calculate)
  }, [productsByCategory])

  useEffect(() => {
    if (categories.length === 0) return

    const offset = stickyHeight + titleGap

    function handleScroll() {
      let active = categories[0]?.id ?? null

      for (let i = 0; i < categories.length; i++) {
        const el = document.getElementById(`category-${categories[i].id}`)
        if (!el) continue
        if (el.getBoundingClientRect().top > offset) continue

        active = categories[i].id

        const next = categories[i + 1]
        if (!next) continue
        const nextEl = document.getElementById(`category-${next.id}`)
        if (!nextEl) continue

        const visibleHeight = nextEl.getBoundingClientRect().top - offset
        const sectionHeight = el.getBoundingClientRect().height
        if (visibleHeight < sectionHeight * DEAD_ZONE_RATIO) {
          active = next.id
        }
      }

      setActiveCategoryId(active)
    }

    window.addEventListener("scroll", handleScroll, { passive: true })
    return () => window.removeEventListener("scroll", handleScroll)
  }, [categories, stickyHeight, titleGap])

  useEffect(() => {
    if (!activeCategoryId || !categoryBarRef.current) return
    const btn = categoryBarRef.current.querySelector(`[data-category="${activeCategoryId}"]`)
    btn?.scrollIntoView({ behavior: "smooth", block: "nearest", inline: "center" })
  }, [activeCategoryId])

  function scrollToCategory(id: string) {
    document.getElementById(`category-${id}`)?.scrollIntoView({ behavior: "smooth", block: "start" })
  }

  return (
    <div className="min-h-screen bg-background text-foreground">
      <header className="border-b px-6 py-4">
        <h1 className="text-xl font-semibold">{layout?.tenant.name}</h1>
      </header>

      <div className="bg-primary/5 border-b border-primary/20 px-6 py-4 flex items-center justify-between gap-4">
        <div>
          <p className="text-sm font-semibold">Receba sempre, sem esquecer</p>
          <p className="text-xs text-muted-foreground">Monte sua cesta e receba nos dias que escolher</p>
        </div>
        <Link
          to="/subscribe"
          className="shrink-0 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors"
        >
          Assinar
        </Link>
      </div>

      <div ref={categoryBarRef} className="sticky top-0 z-10 bg-background border-b flex gap-2 overflow-x-auto px-4 py-2">
        {categories.map((cat) => (
          <button
            key={cat.id}
            type="button"
            data-category={cat.id}
            onClick={() => scrollToCategory(cat.id)}
            className={cn(
              "shrink-0 rounded-full border px-4 py-1.5 text-sm font-medium transition-colors whitespace-nowrap",
              activeCategoryId === cat.id
                ? "bg-primary text-primary-foreground border-primary"
                : "bg-background hover:bg-muted border-border"
            )}
          >
            {cat.name}
          </button>
        ))}
      </div>

      <main
        ref={mainRef}
        className="mx-auto max-w-6xl px-4 py-6 space-y-10"
        style={{ paddingBottom: bottomPadding > 0 ? `${bottomPadding}px` : undefined }}
      >
        {productsByCategory.map(({ category, products: catProducts }, index) => {
          if (catProducts.length === 0) return null
          return (
            <section
              key={category.id}
              id={`category-${category.id}`}
              style={{ scrollMarginTop: `${stickyHeight + titleGap}px` }}
            >
              <h2
                ref={index === 0 ? firstTitleRef : undefined}
                className="text-base font-semibold mb-4"
              >
                {category.name}
              </h2>
              <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
                {catProducts.map((product) => (
                  <button
                    key={product.id}
                    type="button"
                    onClick={() => setDetailProduct(product)}
                    className="group rounded-xl border bg-card overflow-hidden hover:shadow-md transition-shadow text-left"
                  >
                    {product.photo && (
                      <img
                        src={product.photo}
                        alt={product.name}
                        className="aspect-square w-full object-cover group-hover:opacity-90 transition-opacity"
                      />
                    )}
                    {!product.photo && (
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
                  </button>
                ))}
              </div>
            </section>
          )
        })}
      </main>

      <ProductSheet product={detailProduct} onClose={() => setDetailProduct(null)} />
    </div>
  )
}
