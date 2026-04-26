import { cn } from "~/lib/utils"

type Category = { readonly id: string; readonly name: string }

type Props = {
  readonly ref: React.RefObject<HTMLDivElement | null>
  readonly categories: readonly Category[]
  readonly activeCategoryId: string | null
  readonly onCategoryClick: (id: string) => void
}

export function CategoryBar({ ref, categories, activeCategoryId, onCategoryClick }: Props) {
  return (
    <div ref={ref} className="sticky top-0 z-10 bg-background border-b flex gap-2 overflow-x-auto px-4 py-2">
      {categories.map((cat) => (
        <button
          key={cat.id}
          type="button"
          data-category={cat.id}
          onClick={() => onCategoryClick(cat.id)}
          className={cn(
            "shrink-0 rounded-full border px-4 py-1.5 text-sm font-medium transition-colors whitespace-nowrap",
            activeCategoryId === cat.id
              ? "bg-primary text-primary-foreground border-primary"
              : "bg-background hover:bg-muted border-border",
          )}
        >
          {cat.name}
        </button>
      ))}
    </div>
  )
}
